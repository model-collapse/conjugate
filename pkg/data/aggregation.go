package data

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"github.com/conjugate/conjugate/pkg/data/diagon"
)

// computeDataNodeAggregations computes aggregations from search results on the data node.
// This avoids transferring all documents to the coordinator for aggregation.
func computeDataNodeAggregations(hits []*diagon.Hit, aggsMap map[string]interface{}) map[string]*pb.AggregationResult {
	results := make(map[string]*pb.AggregationResult)

	for name, aggDef := range aggsMap {
		aggDefMap, ok := aggDef.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract sub-aggregations if present
		var subAggsMap map[string]interface{}
		if sub, ok := aggDefMap["aggs"]; ok {
			subAggsMap, _ = sub.(map[string]interface{})
		}

		for aggType, aggBody := range aggDefMap {
			if aggType == "aggs" {
				continue
			}
			bodyMap, ok := aggBody.(map[string]interface{})
			if !ok {
				continue
			}

			result := computeSingleDataNodeAgg(hits, aggType, bodyMap, subAggsMap)
			if result != nil {
				results[name] = result
			}
		}
	}

	return results
}

func computeSingleDataNodeAgg(hits []*diagon.Hit, aggType string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	field, _ := body["field"].(string)

	switch aggType {
	case "terms":
		return computeTermsAggData(hits, field, body, subAggsMap)
	case "date_histogram":
		return computeDateHistogramAggData(hits, field, body, subAggsMap)
	case "histogram":
		return computeHistogramAggData(hits, field, body)
	case "avg":
		return computeAvgAggData(hits, field)
	case "sum":
		return computeSumAggData(hits, field)
	case "min":
		return computeMinAggData(hits, field)
	case "max":
		return computeMaxAggData(hits, field)
	case "stats":
		return computeStatsAggData(hits, field)
	case "extended_stats":
		return computeExtendedStatsAggData(hits, field)
	case "value_count":
		return computeValueCountAggData(hits, field)
	case "cardinality":
		return computeCardinalityAggData(hits, field)
	default:
		return nil
	}
}

func computeTermsAggData(hits []*diagon.Hit, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	} else if s, ok := body["size"].(int); ok {
		size = s
	}

	counts := make(map[string]int64)
	bucketHits := make(map[string][]*diagon.Hit)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val == nil {
			continue
		}
		key := fmt.Sprintf("%v", val)
		counts[key]++
		if len(subAggsMap) > 0 {
			bucketHits[key] = append(bucketHits[key], hit)
		}
	}

	type kv struct {
		key   string
		count int64
	}
	sorted := make([]kv, 0, len(counts))
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})
	if len(sorted) > size {
		sorted = sorted[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(sorted))
	for i, kv := range sorted {
		bucket := &pb.AggregationBucket{
			Key:      kv.key,
			DocCount: kv.count,
		}
		if len(subAggsMap) > 0 {
			bucket.SubAggregations = computeDataNodeAggregations(bucketHits[kv.key], subAggsMap)
		}
		buckets[i] = bucket
	}

	return &pb.AggregationResult{
		Type:    "terms",
		Buckets: buckets,
	}
}

func computeDateHistogramAggData(hits []*diagon.Hit, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	interval := "1h"
	if v, ok := body["calendar_interval"].(string); ok {
		interval = v
	} else if v, ok := body["fixed_interval"].(string); ok {
		interval = v
	} else if v, ok := body["interval"].(string); ok {
		interval = v
	}

	duration := parseIntervalData(interval)

	type bucketData struct {
		key   time.Time
		count int64
		hits  []*diagon.Hit
	}
	bucketMap := make(map[int64]*bucketData)

	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val == nil {
			continue
		}

		var ts time.Time
		switch v := val.(type) {
		case string:
			t, err := tryParseTimeData(v)
			if err != nil {
				continue
			}
			ts = t
		case float64:
			ts = time.UnixMilli(int64(v))
		default:
			continue
		}

		bucketTime := truncateToIntervalData(ts, duration)
		bucketKey := bucketTime.UnixMilli()

		if bd, ok := bucketMap[bucketKey]; ok {
			bd.count++
			if len(subAggsMap) > 0 {
				bd.hits = append(bd.hits, hit)
			}
		} else {
			bd := &bucketData{key: bucketTime, count: 1}
			if len(subAggsMap) > 0 {
				bd.hits = []*diagon.Hit{hit}
			}
			bucketMap[bucketKey] = bd
		}
	}

	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	buckets := make([]*pb.AggregationBucket, len(keys))
	for i, k := range keys {
		bd := bucketMap[k]
		bucket := &pb.AggregationBucket{
			Key:        bd.key.Format(time.RFC3339),
			NumericKey: float64(bd.key.UnixMilli()),
			DocCount:   bd.count,
		}
		if len(subAggsMap) > 0 {
			bucket.SubAggregations = computeDataNodeAggregations(bd.hits, subAggsMap)
		}
		buckets[i] = bucket
	}

	return &pb.AggregationResult{
		Type:    "date_histogram",
		Buckets: buckets,
	}
}

func computeHistogramAggData(hits []*diagon.Hit, field string, body map[string]interface{}) *pb.AggregationResult {
	interval := 100.0
	if v, ok := body["interval"].(float64); ok {
		interval = v
	}

	bucketCounts := make(map[float64]int64)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			bucketKey := math.Floor(f/interval) * interval
			bucketCounts[bucketKey]++
		}
	}

	keys := make([]float64, 0, len(bucketCounts))
	for k := range bucketCounts {
		keys = append(keys, k)
	}
	sort.Float64s(keys)

	buckets := make([]*pb.AggregationBucket, len(keys))
	for i, k := range keys {
		buckets[i] = &pb.AggregationBucket{
			NumericKey: k,
			DocCount:   bucketCounts[k],
		}
	}

	return &pb.AggregationResult{
		Type:    "histogram",
		Buckets: buckets,
	}
}

func computeAvgAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	sum := 0.0
	count := int64(0)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			sum += f
			count++
		}
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	}
	return &pb.AggregationResult{
		Type:  "avg",
		Avg:   avg,
		Count: count,
		Sum:   sum,
	}
}

func computeSumAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	sum := 0.0
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			sum += f
		}
	}
	return &pb.AggregationResult{
		Type: "sum",
		Sum:  sum,
	}
}

func computeMinAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	minVal := math.Inf(1)
	found := false
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			if f < minVal {
				minVal = f
			}
			found = true
		}
	}
	if !found {
		minVal = 0
	}
	return &pb.AggregationResult{
		Type: "min",
		Min:  minVal,
	}
}

func computeMaxAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	maxVal := math.Inf(-1)
	found := false
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			if f > maxVal {
				maxVal = f
			}
			found = true
		}
	}
	if !found {
		maxVal = 0
	}
	return &pb.AggregationResult{
		Type: "max",
		Max:  maxVal,
	}
}

func computeStatsAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	sum := 0.0
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	count := int64(0)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			sum += f
			if f < minVal {
				minVal = f
			}
			if f > maxVal {
				maxVal = f
			}
			count++
		}
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	} else {
		minVal = 0
		maxVal = 0
	}
	return &pb.AggregationResult{
		Type:  "stats",
		Count: count,
		Min:   minVal,
		Max:   maxVal,
		Avg:   avg,
		Sum:   sum,
	}
}

func computeExtendedStatsAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	sum := 0.0
	sumOfSquares := 0.0
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	count := int64(0)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if f, ok := toFloat64Data(val); ok {
			sum += f
			sumOfSquares += f * f
			if f < minVal {
				minVal = f
			}
			if f > maxVal {
				maxVal = f
			}
			count++
		}
	}
	avg := 0.0
	variance := 0.0
	stdDev := 0.0
	if count > 0 {
		avg = sum / float64(count)
		variance = (sumOfSquares / float64(count)) - (avg * avg)
		if variance > 0 {
			stdDev = math.Sqrt(variance)
		}
	} else {
		minVal = 0
		maxVal = 0
	}
	return &pb.AggregationResult{
		Type:                       "extended_stats",
		Count:                      count,
		Min:                        minVal,
		Max:                        maxVal,
		Avg:                        avg,
		Sum:                        sum,
		SumOfSquares:               sumOfSquares,
		Variance:                   variance,
		StdDeviation:               stdDev,
		StdDeviationBoundsUpper:    avg + 2.0*stdDev,
		StdDeviationBoundsLower:    avg - 2.0*stdDev,
	}
}

func computeValueCountAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	count := int64(0)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val != nil {
			count++
		}
	}
	return &pb.AggregationResult{
		Type:  "value_count",
		Count: count,
	}
}

func computeCardinalityAggData(hits []*diagon.Hit, field string) *pb.AggregationResult {
	unique := make(map[string]struct{})
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val != nil {
			unique[fmt.Sprintf("%v", val)] = struct{}{}
		}
	}
	return &pb.AggregationResult{
		Type:  "cardinality",
		Value: int64(len(unique)),
	}
}

// resolveField resolves dotted field paths like "cloud.region" in nested maps.
func resolveField(doc map[string]interface{}, field string) interface{} {
	if val, ok := doc[field]; ok {
		return val
	}
	parts := strings.Split(field, ".")
	var current interface{} = doc
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}

func toFloat64Data(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func parseIntervalData(interval string) time.Duration {
	switch interval {
	case "1s", "second":
		return time.Second
	case "1m", "minute":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "10m":
		return 10 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h", "hour":
		return time.Hour
	case "3h":
		return 3 * time.Hour
	case "12h":
		return 12 * time.Hour
	case "1d", "day":
		return 24 * time.Hour
	case "1w", "week":
		return 7 * 24 * time.Hour
	case "1M", "month":
		return 30 * 24 * time.Hour
	case "1q", "quarter":
		return 90 * 24 * time.Hour
	case "1y", "year":
		return 365 * 24 * time.Hour
	default:
		return time.Hour
	}
}

func truncateToIntervalData(t time.Time, interval time.Duration) time.Time {
	if interval >= 24*time.Hour {
		y, m, d := t.UTC().Date()
		days := int(interval / (24 * time.Hour))
		if days > 1 {
			dayOfYear := t.UTC().YearDay()
			d = d - (dayOfYear-1)%days
		}
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}
	return t.UTC().Truncate(interval)
}

func tryParseTimeData(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}

// computeAggregationsFromDocValues computes aggregations from lightweight AggDocValues.
// This is faster than computeDataNodeAggregations because it avoids building
// Hit objects and their Source maps - fields are already extracted as flat string values.
func computeAggregationsFromDocValues(docs []diagon.AggDocValues, aggsMap map[string]interface{}) map[string]*pb.AggregationResult {
	results := make(map[string]*pb.AggregationResult)

	for name, aggDef := range aggsMap {
		aggDefMap, ok := aggDef.(map[string]interface{})
		if !ok {
			continue
		}

		var subAggsMap map[string]interface{}
		if sub, ok := aggDefMap["aggs"]; ok {
			subAggsMap, _ = sub.(map[string]interface{})
		}

		for aggType, aggBody := range aggDefMap {
			if aggType == "aggs" {
				continue
			}
			bodyMap, ok := aggBody.(map[string]interface{})
			if !ok {
				continue
			}

			result := computeSingleAggFromDocValues(docs, aggType, bodyMap, subAggsMap)
			if result != nil {
				results[name] = result
			}
		}
	}

	return results
}

func computeSingleAggFromDocValues(docs []diagon.AggDocValues, aggType string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	field, _ := body["field"].(string)

	switch aggType {
	case "terms":
		return computeTermsAggDocValues(docs, field, body, subAggsMap)
	case "date_histogram":
		return computeDateHistogramAggDocValues(docs, field, body, subAggsMap)
	case "avg":
		return computeNumericAggDocValues(docs, field, "avg")
	case "sum":
		return computeNumericAggDocValues(docs, field, "sum")
	case "min":
		return computeNumericAggDocValues(docs, field, "min")
	case "max":
		return computeNumericAggDocValues(docs, field, "max")
	case "stats":
		return computeNumericAggDocValues(docs, field, "stats")
	case "value_count":
		count := int64(0)
		for _, doc := range docs {
			if _, ok := doc.Fields[field]; ok {
				count++
			}
		}
		return &pb.AggregationResult{Type: "value_count", Count: count}
	case "cardinality":
		unique := make(map[string]struct{})
		for _, doc := range docs {
			if fv, ok := doc.Fields[field]; ok && fv.StringVal != "" {
				unique[fv.StringVal] = struct{}{}
			}
		}
		return &pb.AggregationResult{Type: "cardinality", Value: int64(len(unique))}
	default:
		return nil
	}
}

func computeTermsAggDocValues(docs []diagon.AggDocValues, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	}

	counts := make(map[string]int64)
	var bucketDocs map[string][]diagon.AggDocValues
	if len(subAggsMap) > 0 {
		bucketDocs = make(map[string][]diagon.AggDocValues)
	}

	for _, doc := range docs {
		fv, ok := doc.Fields[field]
		if !ok || fv.StringVal == "" {
			continue
		}
		key := fv.StringVal
		counts[key]++
		if bucketDocs != nil {
			bucketDocs[key] = append(bucketDocs[key], doc)
		}
	}

	type kv struct {
		key   string
		count int64
	}
	sorted := make([]kv, 0, len(counts))
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})
	if len(sorted) > size {
		sorted = sorted[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(sorted))
	for i, kv := range sorted {
		bucket := &pb.AggregationBucket{
			Key:      kv.key,
			DocCount: kv.count,
		}
		if bucketDocs != nil {
			bucket.SubAggregations = computeAggregationsFromDocValues(bucketDocs[kv.key], subAggsMap)
		}
		buckets[i] = bucket
	}

	return &pb.AggregationResult{
		Type:    "terms",
		Buckets: buckets,
	}
}

func computeDateHistogramAggDocValues(docs []diagon.AggDocValues, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	interval := "1h"
	if v, ok := body["calendar_interval"].(string); ok {
		interval = v
	} else if v, ok := body["fixed_interval"].(string); ok {
		interval = v
	} else if v, ok := body["interval"].(string); ok {
		interval = v
	}

	duration := parseIntervalData(interval)

	type bucketData struct {
		key   time.Time
		count int64
		docs  []diagon.AggDocValues
	}
	bucketMap := make(map[int64]*bucketData)

	for _, doc := range docs {
		fv, ok := doc.Fields[field]
		if !ok || fv.StringVal == "" {
			continue
		}

		ts, err := tryParseTimeData(fv.StringVal)
		if err != nil {
			continue
		}

		bucketTime := truncateToIntervalData(ts, duration)
		bucketKey := bucketTime.UnixMilli()

		if bd, ok := bucketMap[bucketKey]; ok {
			bd.count++
			if len(subAggsMap) > 0 {
				bd.docs = append(bd.docs, doc)
			}
		} else {
			bd := &bucketData{key: bucketTime, count: 1}
			if len(subAggsMap) > 0 {
				bd.docs = []diagon.AggDocValues{doc}
			}
			bucketMap[bucketKey] = bd
		}
	}

	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	buckets := make([]*pb.AggregationBucket, len(keys))
	for i, k := range keys {
		bd := bucketMap[k]
		bucket := &pb.AggregationBucket{
			Key:        bd.key.Format(time.RFC3339),
			NumericKey: float64(bd.key.UnixMilli()),
			DocCount:   bd.count,
		}
		if len(subAggsMap) > 0 {
			bucket.SubAggregations = computeAggregationsFromDocValues(bd.docs, subAggsMap)
		}
		buckets[i] = bucket
	}

	return &pb.AggregationResult{
		Type:    "date_histogram",
		Buckets: buckets,
	}
}

// computeNumericAggDocValues computes numeric aggregations (avg, sum, min, max, stats)
// from doc values where all values are stored as strings.
func computeNumericAggDocValues(docs []diagon.AggDocValues, field string, aggType string) *pb.AggregationResult {
	sum := 0.0
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	count := int64(0)

	for _, doc := range docs {
		fv, ok := doc.Fields[field]
		if !ok || fv.StringVal == "" {
			continue
		}
		f, err := strconv.ParseFloat(fv.StringVal, 64)
		if err != nil {
			continue
		}
		sum += f
		if f < minVal {
			minVal = f
		}
		if f > maxVal {
			maxVal = f
		}
		count++
	}

	if count == 0 {
		minVal = 0
		maxVal = 0
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	}

	switch aggType {
	case "avg":
		return &pb.AggregationResult{Type: "avg", Avg: avg, Count: count, Sum: sum}
	case "sum":
		return &pb.AggregationResult{Type: "sum", Sum: sum}
	case "min":
		return &pb.AggregationResult{Type: "min", Min: minVal}
	case "max":
		return &pb.AggregationResult{Type: "max", Max: maxVal}
	case "stats":
		return &pb.AggregationResult{Type: "stats", Count: count, Min: minVal, Max: maxVal, Avg: avg, Sum: sum}
	default:
		return nil
	}
}
