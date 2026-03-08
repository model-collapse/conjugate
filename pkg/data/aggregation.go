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
	case "range":
		return computeRangeAggDocValues(docs, field, body)
	case "composite":
		return computeCompositeAggDocValues(docs, body)
	case "multi_terms":
		return computeMultiTermsAggDocValues(docs, body)
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

// --- Range Aggregation ---

func computeRangeAggData(hits []*diagon.Hit, field string, body map[string]interface{}) *pb.AggregationResult {
	ranges := parseRanges(body)
	if len(ranges) == 0 {
		return nil
	}

	buckets := make([]*pb.AggregationBucket, len(ranges))
	for i, r := range ranges {
		buckets[i] = &pb.AggregationBucket{Key: r.key}
		if r.hasFrom {
			from := r.from
			buckets[i].From = &from
		}
		if r.hasTo {
			to := r.to
			buckets[i].To = &to
		}
	}

	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		f, ok := toFloat64Data(val)
		if !ok {
			continue
		}
		for i, r := range ranges {
			if r.contains(f) {
				buckets[i].DocCount++
			}
		}
	}

	return &pb.AggregationResult{Type: "range", Buckets: buckets}
}

func computeRangeAggDocValues(docs []diagon.AggDocValues, field string, body map[string]interface{}) *pb.AggregationResult {
	ranges := parseRanges(body)
	if len(ranges) == 0 {
		return nil
	}

	buckets := make([]*pb.AggregationBucket, len(ranges))
	for i, r := range ranges {
		buckets[i] = &pb.AggregationBucket{Key: r.key}
		if r.hasFrom {
			from := r.from
			buckets[i].From = &from
		}
		if r.hasTo {
			to := r.to
			buckets[i].To = &to
		}
	}

	for _, doc := range docs {
		fv, ok := doc.Fields[field]
		if !ok || fv.StringVal == "" {
			continue
		}
		f, err := strconv.ParseFloat(fv.StringVal, 64)
		if err != nil {
			continue
		}
		for i, r := range ranges {
			if r.contains(f) {
				buckets[i].DocCount++
			}
		}
	}

	return &pb.AggregationResult{Type: "range", Buckets: buckets}
}

type rangeSpec struct {
	from, to       float64
	hasFrom, hasTo bool
	key            string
}

func (r *rangeSpec) contains(v float64) bool {
	if r.hasFrom && v < r.from {
		return false
	}
	if r.hasTo && v >= r.to {
		return false
	}
	return true
}

func parseRanges(body map[string]interface{}) []rangeSpec {
	rawRanges, ok := body["ranges"].([]interface{})
	if !ok {
		return nil
	}
	specs := make([]rangeSpec, 0, len(rawRanges))
	for _, raw := range rawRanges {
		rm, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		var r rangeSpec
		if from, ok := rm["from"].(float64); ok {
			r.from = from
			r.hasFrom = true
		}
		if to, ok := rm["to"].(float64); ok {
			r.to = to
			r.hasTo = true
		}
		if key, ok := rm["key"].(string); ok {
			r.key = key
		} else {
			fromStr := "*"
			toStr := "*"
			if r.hasFrom {
				fromStr = strconv.FormatFloat(r.from, 'f', -1, 64)
			}
			if r.hasTo {
				toStr = strconv.FormatFloat(r.to, 'f', -1, 64)
			}
			r.key = fromStr + "-" + toStr
		}
		specs = append(specs, r)
	}
	return specs
}

// --- Auto Date Histogram Aggregation ---

func computeAutoDateHistogramAggData(hits []*diagon.Hit, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	targetBuckets := 10
	if b, ok := body["buckets"].(float64); ok {
		targetBuckets = int(b)
	} else if b, ok := body["buckets"].(int); ok {
		targetBuckets = b
	}
	if targetBuckets < 1 {
		targetBuckets = 1
	}

	var timestamps []time.Time
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case string:
			if t, err := tryParseTimeData(v); err == nil {
				timestamps = append(timestamps, t)
			}
		case float64:
			timestamps = append(timestamps, time.UnixMilli(int64(v)))
		}
	}

	if len(timestamps) == 0 {
		return &pb.AggregationResult{Type: "auto_date_histogram", Buckets: []*pb.AggregationBucket{}}
	}

	minTime, maxTime := timestamps[0], timestamps[0]
	for _, ts := range timestamps[1:] {
		if ts.Before(minTime) {
			minTime = ts
		}
		if ts.After(maxTime) {
			maxTime = ts
		}
	}

	interval := autoSelectInterval(minTime, maxTime, targetBuckets)

	type bucketData struct {
		key   time.Time
		count int64
		hits  []*diagon.Hit
	}
	bucketMap := make(map[int64]*bucketData)
	for i, ts := range timestamps {
		bucketTime := truncateToIntervalData(ts, interval)
		bk := bucketTime.UnixMilli()
		if bd, ok := bucketMap[bk]; ok {
			bd.count++
			if len(subAggsMap) > 0 {
				bd.hits = append(bd.hits, hits[i])
			}
		} else {
			bd := &bucketData{key: bucketTime, count: 1}
			if len(subAggsMap) > 0 {
				bd.hits = []*diagon.Hit{hits[i]}
			}
			bucketMap[bk] = bd
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

	return &pb.AggregationResult{Type: "auto_date_histogram", Buckets: buckets}
}

func computeAutoDateHistogramAggDocValues(docs []diagon.AggDocValues, field string, body map[string]interface{}, subAggsMap map[string]interface{}) *pb.AggregationResult {
	targetBuckets := 10
	if b, ok := body["buckets"].(float64); ok {
		targetBuckets = int(b)
	}
	if targetBuckets < 1 {
		targetBuckets = 1
	}

	var timestamps []time.Time
	for _, doc := range docs {
		fv, ok := doc.Fields[field]
		if !ok || fv.StringVal == "" {
			continue
		}
		if ts, err := tryParseTimeData(fv.StringVal); err == nil {
			timestamps = append(timestamps, ts)
		}
	}

	if len(timestamps) == 0 {
		return &pb.AggregationResult{Type: "auto_date_histogram", Buckets: []*pb.AggregationBucket{}}
	}

	minTime, maxTime := timestamps[0], timestamps[0]
	for _, ts := range timestamps[1:] {
		if ts.Before(minTime) {
			minTime = ts
		}
		if ts.After(maxTime) {
			maxTime = ts
		}
	}

	interval := autoSelectInterval(minTime, maxTime, targetBuckets)

	type bucketData struct {
		key   time.Time
		count int64
		docs  []diagon.AggDocValues
	}
	bucketMap := make(map[int64]*bucketData)
	for i, ts := range timestamps {
		bucketTime := truncateToIntervalData(ts, interval)
		bk := bucketTime.UnixMilli()
		if bd, ok := bucketMap[bk]; ok {
			bd.count++
			if len(subAggsMap) > 0 {
				bd.docs = append(bd.docs, docs[i])
			}
		} else {
			bd := &bucketData{key: bucketTime, count: 1}
			if len(subAggsMap) > 0 {
				bd.docs = []diagon.AggDocValues{docs[i]}
			}
			bucketMap[bk] = bd
		}
	}

	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	buckets := make([]*pb.AggregationBucket, len(keys))
	for idx, k := range keys {
		bd := bucketMap[k]
		bucket := &pb.AggregationBucket{
			Key:        bd.key.Format(time.RFC3339),
			NumericKey: float64(bd.key.UnixMilli()),
			DocCount:   bd.count,
		}
		if len(subAggsMap) > 0 {
			bucket.SubAggregations = computeAggregationsFromDocValues(bd.docs, subAggsMap)
		}
		buckets[idx] = bucket
	}

	return &pb.AggregationResult{Type: "auto_date_histogram", Buckets: buckets}
}

func autoSelectInterval(minTime, maxTime time.Time, targetBuckets int) time.Duration {
	dataRange := maxTime.Sub(minTime)
	if dataRange <= 0 {
		return time.Second
	}
	intervals := []time.Duration{
		time.Second, time.Minute, 5 * time.Minute, 10 * time.Minute,
		30 * time.Minute, time.Hour, 3 * time.Hour, 12 * time.Hour,
		24 * time.Hour, 7 * 24 * time.Hour, 30 * 24 * time.Hour,
		90 * 24 * time.Hour, 365 * 24 * time.Hour,
	}
	for _, iv := range intervals {
		if int(dataRange/iv)+1 <= targetBuckets {
			return iv
		}
	}
	return 365 * 24 * time.Hour
}

// --- Significant Terms Aggregation ---

func computeSignificantTermsAggData(hits []*diagon.Hit, field string, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	} else if s, ok := body["size"].(int); ok {
		size = s
	}

	fgCounts := make(map[string]int64)
	for _, hit := range hits {
		val := resolveField(hit.Source, field)
		if val == nil {
			continue
		}
		fgCounts[fmt.Sprintf("%v", val)]++
	}

	totalFg := int64(len(hits))
	type termScore struct {
		key   string
		count int64
		score float64
	}
	scored := make([]termScore, 0, len(fgCounts))
	for k, count := range fgCounts {
		fgPercent := float64(count) / float64(totalFg)
		scored = append(scored, termScore{k, count, fgPercent * float64(count)})
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })
	if len(scored) > size {
		scored = scored[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(scored))
	for i, ts := range scored {
		buckets[i] = &pb.AggregationBucket{Key: ts.key, DocCount: ts.count}
	}

	return &pb.AggregationResult{Type: "significant_terms", Buckets: buckets}
}

func computeSignificantTermsAggDocValues(docs []diagon.AggDocValues, field string, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	}

	fgCounts := make(map[string]int64)
	for _, doc := range docs {
		if fv, ok := doc.Fields[field]; ok && fv.StringVal != "" {
			fgCounts[fv.StringVal]++
		}
	}

	totalFg := int64(len(docs))
	type termScore struct {
		key   string
		count int64
		score float64
	}
	scored := make([]termScore, 0, len(fgCounts))
	for k, count := range fgCounts {
		fgPercent := float64(count) / float64(totalFg)
		scored = append(scored, termScore{k, count, fgPercent * float64(count)})
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })
	if len(scored) > size {
		scored = scored[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(scored))
	for i, ts := range scored {
		buckets[i] = &pb.AggregationBucket{Key: ts.key, DocCount: ts.count}
	}

	return &pb.AggregationResult{Type: "significant_terms", Buckets: buckets}
}

// --- Multi-Terms Aggregation ---

func computeMultiTermsAggData(hits []*diagon.Hit, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	} else if s, ok := body["size"].(int); ok {
		size = s
	}

	fields := extractMultiTermsFields(body)
	if len(fields) == 0 {
		return nil
	}

	counts := make(map[string]int64)
	for _, hit := range hits {
		key := buildMultiTermsKey(func(field string) interface{} {
			return resolveField(hit.Source, field)
		}, fields)
		if key != "" {
			counts[key]++
		}
	}

	return buildMultiTermsResult(counts, size)
}

func computeMultiTermsAggDocValues(docs []diagon.AggDocValues, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	}

	fields := extractMultiTermsFields(body)
	if len(fields) == 0 {
		return nil
	}

	counts := make(map[string]int64)
	for _, doc := range docs {
		key := buildMultiTermsKey(func(field string) interface{} {
			if fv, ok := doc.Fields[field]; ok && fv.StringVal != "" {
				return fv.StringVal
			}
			return nil
		}, fields)
		if key != "" {
			counts[key]++
		}
	}

	return buildMultiTermsResult(counts, size)
}

func extractMultiTermsFields(body map[string]interface{}) []string {
	terms, ok := body["terms"].([]interface{})
	if !ok {
		return nil
	}
	fields := make([]string, 0, len(terms))
	for _, t := range terms {
		if tm, ok := t.(map[string]interface{}); ok {
			if f, ok := tm["field"].(string); ok {
				fields = append(fields, f)
			}
		}
	}
	return fields
}

func buildMultiTermsKey(resolver func(string) interface{}, fields []string) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		val := resolver(f)
		if val == nil {
			return ""
		}
		parts = append(parts, fmt.Sprintf("%v", val))
	}
	return strings.Join(parts, "\x00")
}

func buildMultiTermsResult(counts map[string]int64, size int) *pb.AggregationResult {
	type kv struct {
		key   string
		count int64
	}
	sorted := make([]kv, 0, len(counts))
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })
	if len(sorted) > size {
		sorted = sorted[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(sorted))
	for i, kv := range sorted {
		displayKey := strings.ReplaceAll(kv.key, "\x00", "|")
		buckets[i] = &pb.AggregationBucket{Key: displayKey, DocCount: kv.count}
	}

	return &pb.AggregationResult{Type: "multi_terms", Buckets: buckets}
}

// --- Composite Aggregation ---

func computeCompositeAggData(hits []*diagon.Hit, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	} else if s, ok := body["size"].(int); ok {
		size = s
	}

	sources := parseCompositeSources(body)
	if len(sources) == 0 {
		return nil
	}

	counts := make(map[string]int64)
	for _, hit := range hits {
		key := buildCompositeKey(func(src compositeSource) string {
			val := resolveField(hit.Source, src.field)
			if val == nil {
				return ""
			}
			return compositeSourceValue(val, src)
		}, sources)
		if key != "" {
			counts[key]++
		}
	}

	return buildCompositeResult(counts, size, sources)
}

func computeCompositeAggDocValues(docs []diagon.AggDocValues, body map[string]interface{}) *pb.AggregationResult {
	size := 10
	if s, ok := body["size"].(float64); ok {
		size = int(s)
	}

	sources := parseCompositeSources(body)
	if len(sources) == 0 {
		return nil
	}

	counts := make(map[string]int64)
	for _, doc := range docs {
		key := buildCompositeKey(func(src compositeSource) string {
			fv, ok := doc.Fields[src.field]
			if !ok || fv.StringVal == "" {
				return ""
			}
			if src.sourceType == "date_histogram" {
				ts, err := tryParseTimeData(fv.StringVal)
				if err != nil {
					return ""
				}
				interval := parseIntervalData(src.interval)
				return truncateToIntervalData(ts, interval).Format(time.RFC3339)
			}
			return fv.StringVal
		}, sources)
		if key != "" {
			counts[key]++
		}
	}

	return buildCompositeResult(counts, size, sources)
}

type compositeSource struct {
	name        string
	sourceType  string // "terms", "date_histogram", "histogram"
	field       string
	interval    string  // For date_histogram
	numInterval float64 // For histogram
}

func parseCompositeSources(body map[string]interface{}) []compositeSource {
	rawSources, ok := body["sources"].([]interface{})
	if !ok {
		return nil
	}
	sources := make([]compositeSource, 0, len(rawSources))
	for _, raw := range rawSources {
		rm, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		for name, def := range rm {
			dm, ok := def.(map[string]interface{})
			if !ok {
				continue
			}
			for srcType, srcBody := range dm {
				sb, ok := srcBody.(map[string]interface{})
				if !ok {
					continue
				}
				src := compositeSource{name: name, sourceType: srcType}
				if f, ok := sb["field"].(string); ok {
					src.field = f
				}
				if v, ok := sb["calendar_interval"].(string); ok {
					src.interval = v
				} else if v, ok := sb["fixed_interval"].(string); ok {
					src.interval = v
				} else if v, ok := sb["interval"].(string); ok {
					src.interval = v
				}
				if v, ok := sb["interval"].(float64); ok {
					src.numInterval = v
				}
				sources = append(sources, src)
			}
		}
	}
	return sources
}

func compositeSourceValue(val interface{}, src compositeSource) string {
	switch src.sourceType {
	case "date_histogram":
		var ts time.Time
		switch v := val.(type) {
		case string:
			t, err := tryParseTimeData(v)
			if err != nil {
				return ""
			}
			ts = t
		case float64:
			ts = time.UnixMilli(int64(v))
		default:
			return ""
		}
		return truncateToIntervalData(ts, parseIntervalData(src.interval)).Format(time.RFC3339)
	case "histogram":
		f, ok := toFloat64Data(val)
		if !ok {
			return ""
		}
		iv := src.numInterval
		if iv <= 0 {
			iv = 1
		}
		return strconv.FormatFloat(math.Floor(f/iv)*iv, 'f', -1, 64)
	default: // "terms"
		return fmt.Sprintf("%v", val)
	}
}

func buildCompositeKey(resolver func(compositeSource) string, sources []compositeSource) string {
	parts := make([]string, 0, len(sources))
	for _, src := range sources {
		v := resolver(src)
		if v == "" {
			return ""
		}
		parts = append(parts, src.name+"="+v)
	}
	return strings.Join(parts, "\x00")
}

func buildCompositeResult(counts map[string]int64, size int, sources []compositeSource) *pb.AggregationResult {
	type kv struct {
		key   string
		count int64
	}
	sorted := make([]kv, 0, len(counts))
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].key < sorted[j].key })
	if len(sorted) > size {
		sorted = sorted[:size]
	}

	buckets := make([]*pb.AggregationBucket, len(sorted))
	for i, kv := range sorted {
		displayKey := strings.ReplaceAll(kv.key, "\x00", "|")
		buckets[i] = &pb.AggregationBucket{Key: displayKey, DocCount: kv.count}
	}

	return &pb.AggregationResult{Type: "composite", Buckets: buckets}
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
