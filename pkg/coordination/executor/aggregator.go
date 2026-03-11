package executor

import (
	"encoding/json"
	"math"
	"sort"
	"time"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"go.uber.org/zap"
)

// aggregateSearchResults merges search results from multiple shards
func (qe *QueryExecutor) aggregateSearchResults(responses []*pb.SearchResponse, from, size int) *SearchResult {
	if len(responses) == 0 {
		return &SearchResult{
			TotalHits: 0,
			MaxScore:  0,
			Hits:      []*SearchHit{},
		}
	}

	// Collect all hits from all shards
	var allHits []*SearchHit
	var totalHits int64
	var maxScore float64

	for _, resp := range responses {
		// Sum total hits
		if resp.Hits != nil && resp.Hits.Total != nil {
			totalHits += resp.Hits.Total.Value
		}

		// Track max score
		if resp.Hits != nil && resp.Hits.MaxScore > maxScore {
			maxScore = resp.Hits.MaxScore
		}

		// Collect hits
		if resp.Hits != nil {
			for _, hit := range resp.Hits.Hits {
				var sourceMap map[string]interface{}
				if len(hit.SourceJson) > 0 {
					json.Unmarshal(hit.SourceJson, &sourceMap)
				}
				allHits = append(allHits, &SearchHit{
					ID:     hit.Id,
					Score:  hit.Score,
					Source: sourceMap,
				})
			}
		}
	}

	// Sort hits by score (descending)
	sort.Slice(allHits, func(i, j int) bool {
		return allHits[i].Score > allHits[j].Score
	})

	// Apply pagination (from/size)
	start := from
	if start > len(allHits) {
		start = len(allHits)
	}

	end := start + size
	if end > len(allHits) {
		end = len(allHits)
	}

	paginatedHits := allHits[start:end]

	// Merge aggregations from all shards
	aggregations := qe.mergeAggregations(responses)

	return &SearchResult{
		TotalHits:    totalHits,
		MaxScore:     maxScore,
		Hits:         paginatedHits,
		Aggregations: aggregations,
	}
}

// aggregateCountResults sums document counts from multiple shards
func aggregateCountResults(responses []*pb.CountResponse) int64 {
	var total int64
	for _, resp := range responses {
		total += resp.Count
	}
	return total
}

// mergeAggregations merges aggregations from multiple shard responses
func (qe *QueryExecutor) mergeAggregations(responses []*pb.SearchResponse) map[string]*AggregationResult {
	if len(responses) == 0 {
		return nil
	}

	// Group aggregations by name across all shards
	aggsByName := make(map[string][]*pb.AggregationResult)
	for _, resp := range responses {
		if resp.Aggregations == nil {
			continue
		}
		for name, agg := range resp.Aggregations {
			aggsByName[name] = append(aggsByName[name], agg)
		}
	}

	if len(aggsByName) == 0 {
		return nil
	}

	// Merge each aggregation type
	merged := make(map[string]*AggregationResult)
	for name, aggs := range aggsByName {
		if len(aggs) == 0 {
			continue
		}

		aggType := aggs[0].Type // Assume all shards return same type

		// Track aggregation merge time
		mergeStartTime := time.Now()
		var result *AggregationResult

		switch aggType {
		case "terms", "histogram", "date_histogram", "range", "filters",
			"auto_date_histogram", "significant_terms", "multi_terms", "composite":
			result = qe.mergeBucketAggregation(aggs)
		case "stats":
			result = qe.mergeStatsAggregation(aggs, false)
		case "extended_stats":
			result = qe.mergeStatsAggregation(aggs, true)
		case "percentiles":
			result = qe.mergePercentilesAggregation(aggs)
		case "cardinality":
			result = qe.mergeCardinalityAggregation(aggs)
		case "avg", "min", "max", "sum", "value_count":
			result = qe.mergeSimpleMetricAggregation(aggs)
		default:
			qe.logger.Warn("Unknown aggregation type, skipping merge",
				zap.String("type", aggType),
				zap.String("name", name))
		}

		if result != nil {
			merged[name] = result
			// Record merge time
			aggregationMergeTime.WithLabelValues(aggType).Observe(time.Since(mergeStartTime).Seconds())
		}
	}

	return merged
}

// mergeBucketAggregation merges bucket-based aggregations (terms, histogram, date_histogram, range)
func (qe *QueryExecutor) mergeBucketAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	aggType := aggs[0].Type

	// For range and filters aggregations, preserve bucket order and metadata
	if aggType == "range" {
		return qe.mergeRangeAggregation(aggs)
	}
	if aggType == "filters" {
		return qe.mergeFiltersAggregation(aggs)
	}

	// Single-shard fast path: skip map+sort overhead when there's only one response
	if len(aggs) == 1 {
		single := aggs[0]
		buckets := make([]*AggregationBucket, len(single.Buckets))
		for i, b := range single.Buckets {
			bucket := &AggregationBucket{
				Key:        b.Key,
				NumericKey: b.NumericKey,
				DocCount:   b.DocCount,
			}
			if len(b.SubAggregations) > 0 {
				subAggsByName := make(map[string][]*pb.AggregationResult, len(b.SubAggregations))
				for subName, subAgg := range b.SubAggregations {
					subAggsByName[subName] = []*pb.AggregationResult{subAgg}
				}
				bucket.SubAggs = qe.mergeSubAggregations(subAggsByName)
			}
			buckets[i] = bucket
		}
		return &AggregationResult{Type: aggType, Buckets: buckets}
	}

	// Sum bucket counts across all shards, collecting sub-aggs per bucket
	type bucketData struct {
		count     int64
		stringKey string                             // preserved for numeric-keyed date_histogram
		subAggs   map[string][]*pb.AggregationResult // sub-aggs grouped by name across shards
	}
	bucketMap := make(map[string]*bucketData)         // for string keys (terms)
	numericBucketMap := make(map[float64]*bucketData) // for numeric keys (histogram, date_histogram)

	isNumeric := aggType == "histogram" || aggType == "date_histogram"

	for _, agg := range aggs {
		for _, bucket := range agg.Buckets {
			if isNumeric {
				bd, ok := numericBucketMap[bucket.NumericKey]
				if !ok {
					bd = &bucketData{stringKey: bucket.Key, subAggs: make(map[string][]*pb.AggregationResult)}
					numericBucketMap[bucket.NumericKey] = bd
				}
				bd.count += bucket.DocCount
				for subName, subAgg := range bucket.SubAggregations {
					bd.subAggs[subName] = append(bd.subAggs[subName], subAgg)
				}
			} else {
				bd, ok := bucketMap[bucket.Key]
				if !ok {
					bd = &bucketData{subAggs: make(map[string][]*pb.AggregationResult)}
					bucketMap[bucket.Key] = bd
				}
				bd.count += bucket.DocCount
				for subName, subAgg := range bucket.SubAggregations {
					bd.subAggs[subName] = append(bd.subAggs[subName], subAgg)
				}
			}
		}
	}

	// Convert to result buckets
	var buckets []*AggregationBucket

	if isNumeric {
		// Numeric buckets (histogram, date_histogram)
		for key, bd := range numericBucketMap {
			bucket := &AggregationBucket{
				Key:        bd.stringKey,
				NumericKey: key,
				DocCount:   bd.count,
			}
			if len(bd.subAggs) > 0 {
				bucket.SubAggs = qe.mergeSubAggregations(bd.subAggs)
			}
			buckets = append(buckets, bucket)
		}
		// Sort by numeric key
		sort.Slice(buckets, func(i, j int) bool {
			return buckets[i].NumericKey < buckets[j].NumericKey
		})
	} else {
		// String buckets (terms, etc.)
		for key, bd := range bucketMap {
			bucket := &AggregationBucket{
				Key:      key,
				DocCount: bd.count,
			}
			if len(bd.subAggs) > 0 {
				bucket.SubAggs = qe.mergeSubAggregations(bd.subAggs)
			}
			buckets = append(buckets, bucket)
		}
		// Sort by key for date-based and composite aggs, by doc_count for others
		switch aggType {
		case "date_histogram", "auto_date_histogram", "composite":
			sort.Slice(buckets, func(i, j int) bool {
				return buckets[i].Key < buckets[j].Key
			})
		default: // terms, significant_terms, multi_terms, etc.
			sort.Slice(buckets, func(i, j int) bool {
				return buckets[i].DocCount > buckets[j].DocCount
			})
		}
	}

	return &AggregationResult{
		Type:    aggType,
		Buckets: buckets,
	}
}

// mergeSubAggregations merges sub-aggregation results collected from multiple shard buckets.
// Each sub-agg name maps to a slice of pb.AggregationResult from different shards.
func (qe *QueryExecutor) mergeSubAggregations(subAggsByName map[string][]*pb.AggregationResult) map[string]*AggregationResult {
	merged := make(map[string]*AggregationResult)
	for name, subAggs := range subAggsByName {
		if len(subAggs) == 0 {
			continue
		}
		subType := subAggs[0].Type
		var result *AggregationResult
		switch subType {
		case "terms", "histogram", "date_histogram":
			result = qe.mergeBucketAggregation(subAggs)
		case "stats":
			result = qe.mergeStatsAggregation(subAggs, false)
		case "extended_stats":
			result = qe.mergeStatsAggregation(subAggs, true)
		case "avg", "min", "max", "sum", "value_count":
			result = qe.mergeSimpleMetricAggregation(subAggs)
		case "cardinality":
			result = qe.mergeCardinalityAggregation(subAggs)
		case "percentiles":
			result = qe.mergePercentilesAggregation(subAggs)
		default:
			// Unknown sub-agg type: convert first result directly
			result = &AggregationResult{
				Type: subType,
			}
		}
		if result != nil {
			merged[name] = result
		}
	}
	return merged
}

// mergeRangeAggregation merges range aggregations preserving bucket order and metadata
func (qe *QueryExecutor) mergeRangeAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	// Use first shard's buckets as template (preserves order and range definitions)
	firstAgg := aggs[0]
	buckets := make([]*AggregationBucket, len(firstAgg.Buckets))

	// Initialize buckets from first shard
	for i, bucket := range firstAgg.Buckets {
		buckets[i] = &AggregationBucket{
			Key:      bucket.Key,
			DocCount: bucket.DocCount,
		}
		// Note: from/to fields will be added when protobuf is regenerated
	}

	// Sum counts from remaining shards (matching by key)
	for shardIdx := 1; shardIdx < len(aggs); shardIdx++ {
		for _, bucket := range aggs[shardIdx].Buckets {
			// Find matching bucket by key
			for i, resultBucket := range buckets {
				if resultBucket.Key == bucket.Key {
					buckets[i].DocCount += bucket.DocCount
					break
				}
			}
		}
	}

	return &AggregationResult{
		Type:    "range",
		Buckets: buckets,
	}
}

// mergeFiltersAggregation merges filters aggregations preserving bucket order
func (qe *QueryExecutor) mergeFiltersAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	// Use first shard's buckets as template (preserves order and filter definitions)
	firstAgg := aggs[0]
	buckets := make([]*AggregationBucket, len(firstAgg.Buckets))

	// Initialize buckets from first shard
	for i, bucket := range firstAgg.Buckets {
		buckets[i] = &AggregationBucket{
			Key:      bucket.Key,
			DocCount: bucket.DocCount,
		}
	}

	// Sum counts from remaining shards (matching by key)
	for shardIdx := 1; shardIdx < len(aggs); shardIdx++ {
		for _, bucket := range aggs[shardIdx].Buckets {
			// Find matching bucket by key
			for i, resultBucket := range buckets {
				if resultBucket.Key == bucket.Key {
					buckets[i].DocCount += bucket.DocCount
					break
				}
			}
		}
	}

	return &AggregationResult{
		Type:    "filters",
		Buckets: buckets,
	}
}

// mergeStatsAggregation merges stats and extended_stats aggregations
func (qe *QueryExecutor) mergeStatsAggregation(aggs []*pb.AggregationResult, extended bool) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	result := &AggregationResult{
		Type: aggs[0].Type,
		Min:  aggs[0].Min,
		Max:  aggs[0].Max,
	}

	var totalCount int64
	var totalSum float64
	var totalSumOfSquares float64

	for _, agg := range aggs {
		totalCount += agg.Count
		totalSum += agg.Sum

		// Track global min/max
		if agg.Min < result.Min {
			result.Min = agg.Min
		}
		if agg.Max > result.Max {
			result.Max = agg.Max
		}

		if extended {
			totalSumOfSquares += agg.SumOfSquares
		}
	}

	result.Count = totalCount
	result.Sum = totalSum
	if totalCount > 0 {
		result.Avg = totalSum / float64(totalCount)
	}

	if extended {
		result.SumOfSquares = totalSumOfSquares
		// Calculate variance: Var(X) = E[X²] - E[X]²
		if totalCount > 0 {
			result.Variance = (totalSumOfSquares / float64(totalCount)) - (result.Avg * result.Avg)
			if result.Variance > 0 {
				result.StdDeviation = math.Sqrt(result.Variance)
				result.StdDeviationBoundsUpper = result.Avg + 2.0*result.StdDeviation
				result.StdDeviationBoundsLower = result.Avg - 2.0*result.StdDeviation
			}
		}
	}

	return result
}

// mergePercentilesAggregation merges percentiles aggregations
// Note: This is approximate - collecting all values would be expensive
// For now, we average the percentile values from each shard
func (qe *QueryExecutor) mergePercentilesAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	result := &AggregationResult{
		Type:   "percentiles",
		Values: make(map[string]float64),
	}

	// Track percentile sums and counts for averaging
	percentileSums := make(map[string]float64)
	percentileCounts := make(map[string]int)

	for _, agg := range aggs {
		for percentile, value := range agg.Values {
			percentileSums[percentile] += value
			percentileCounts[percentile]++
		}
	}

	// Average the percentiles
	for percentile, sum := range percentileSums {
		count := percentileCounts[percentile]
		if count > 0 {
			result.Values[percentile] = sum / float64(count)
		}
	}

	return result
}

// mergeCardinalityAggregation merges cardinality aggregations
// Note: This is approximate - true cardinality would require HyperLogLog
// For now, we sum the cardinalities (which may overcount)
func (qe *QueryExecutor) mergeCardinalityAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	result := &AggregationResult{
		Type: "cardinality",
	}

	// Sum cardinalities from all shards
	// Note: This overcounts if same values appear on multiple shards
	var total int64
	for _, agg := range aggs {
		total += agg.Value
	}

	result.Value = total

	return result
}

// mergeSimpleMetricAggregation merges simple metric aggregations (avg, min, max, sum, value_count)
func (qe *QueryExecutor) mergeSimpleMetricAggregation(aggs []*pb.AggregationResult) *AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	aggType := aggs[0].Type
	result := &AggregationResult{
		Type: aggType,
	}

	switch aggType {
	case "avg":
		// Average: compute weighted average across shards
		// We don't have document counts per shard, so we average the averages (approximation)
		var sum float64
		for _, agg := range aggs {
			sum += agg.Avg
		}
		result.Avg = sum / float64(len(aggs))

	case "min":
		// Minimum: take global minimum
		result.Min = aggs[0].Min
		for _, agg := range aggs {
			if agg.Min < result.Min {
				result.Min = agg.Min
			}
		}

	case "max":
		// Maximum: take global maximum
		result.Max = aggs[0].Max
		for _, agg := range aggs {
			if agg.Max > result.Max {
				result.Max = agg.Max
			}
		}

	case "sum":
		// Sum: sum across all shards
		var sum float64
		for _, agg := range aggs {
			sum += agg.Sum
		}
		result.Sum = sum

	case "value_count":
		// Value count: sum across all shards
		var total int64
		for _, agg := range aggs {
			total += agg.Count
		}
		result.Count = total
	}

	return result
}
