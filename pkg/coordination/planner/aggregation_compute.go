package planner

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// computeAggregations computes aggregation results from document rows.
// This runs on the coordinator when aggregation specs are not pushed down to data nodes.
func computeAggregations(rows []map[string]interface{}, aggregations []*Aggregation) map[string]*AggregationResult {
	results := make(map[string]*AggregationResult)
	for _, agg := range aggregations {
		result := computeSingleAggregation(rows, agg)
		if result != nil {
			results[agg.Name] = result
		}
	}
	return results
}

func computeSingleAggregation(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	switch agg.Type {
	case AggTypeTerms:
		return computeTermsAgg(rows, agg)
	case AggTypeDateHistogram:
		return computeDateHistogramAgg(rows, agg)
	case AggTypeAvg:
		return computeAvgAgg(rows, agg)
	case AggTypeSum:
		return computeSumAgg(rows, agg)
	case AggTypeMin:
		return computeMinAgg(rows, agg)
	case AggTypeMax:
		return computeMaxAgg(rows, agg)
	case AggTypeCount:
		return &AggregationResult{
			Type:  AggTypeCount,
			Value: float64(len(rows)),
		}
	case AggTypeStats:
		return computeStatsAgg(rows, agg)
	case AggTypeCardinality:
		return computeCardinalityAgg(rows, agg)
	default:
		return nil
	}
}

// computeTermsAgg groups documents by field value and counts per group.
func computeTermsAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	size := 10
	if s, ok := agg.Params["size"].(int); ok {
		size = s
	}

	// Count by key
	counts := make(map[string]int64)
	bucketRows := make(map[string][]map[string]interface{}) // for sub-aggs
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if val == nil {
			continue
		}
		key := fmt.Sprintf("%v", val)
		counts[key]++
		if len(agg.SubAggs) > 0 {
			bucketRows[key] = append(bucketRows[key], row)
		}
	}

	// Sort by count descending
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

	// Take top N
	if len(sorted) > size {
		sorted = sorted[:size]
	}

	buckets := make([]*Bucket, len(sorted))
	for i, kv := range sorted {
		bucket := &Bucket{
			Key:      kv.key,
			DocCount: kv.count,
			SubAggs:  make(map[string]*AggregationResult),
		}
		// Compute sub-aggregations within this bucket
		if len(agg.SubAggs) > 0 {
			bucket.SubAggs = computeAggregations(bucketRows[kv.key], agg.SubAggs)
		}
		buckets[i] = bucket
	}

	return &AggregationResult{
		Type:    AggTypeTerms,
		Buckets: buckets,
	}
}

// computeDateHistogramAgg buckets documents by time intervals.
func computeDateHistogramAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	// Determine interval
	interval := "1h" // default
	if v, ok := agg.Params["calendar_interval"].(string); ok {
		interval = v
	} else if v, ok := agg.Params["fixed_interval"].(string); ok {
		interval = v
	} else if v, ok := agg.Params["interval"].(string); ok {
		interval = v
	}

	duration := parseInterval(interval)
	if duration == 0 {
		duration = time.Hour // fallback
	}

	// Collect timestamps and bucket
	type bucketData struct {
		key   time.Time
		count int64
		rows  []map[string]interface{}
	}
	bucketMap := make(map[int64]*bucketData)

	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if val == nil {
			continue
		}

		var ts time.Time
		switch v := val.(type) {
		case string:
			t, err := tryParseTime(v)
			if err != nil {
				continue
			}
			ts = t
		case float64:
			// Epoch millis (from double-indexed date fields)
			ts = time.UnixMilli(int64(v))
		default:
			continue
		}

		// Truncate to interval bucket
		bucketTime := truncateToInterval(ts, duration)
		bucketKey := bucketTime.UnixMilli()

		if bd, ok := bucketMap[bucketKey]; ok {
			bd.count++
			if len(agg.SubAggs) > 0 {
				bd.rows = append(bd.rows, row)
			}
		} else {
			bd := &bucketData{key: bucketTime, count: 1}
			if len(agg.SubAggs) > 0 {
				bd.rows = []map[string]interface{}{row}
			}
			bucketMap[bucketKey] = bd
		}
	}

	// Sort buckets by time
	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	buckets := make([]*Bucket, len(keys))
	for i, k := range keys {
		bd := bucketMap[k]
		bucket := &Bucket{
			Key:      bd.key.Format(time.RFC3339),
			DocCount: bd.count,
			SubAggs:  make(map[string]*AggregationResult),
		}
		if len(agg.SubAggs) > 0 {
			bucket.SubAggs = computeAggregations(bd.rows, agg.SubAggs)
		}
		buckets[i] = bucket
	}

	return &AggregationResult{
		Type:    AggTypeDateHistogram,
		Buckets: buckets,
	}
}

func computeAvgAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	sum := 0.0
	count := 0
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if val == nil {
			continue
		}
		if f, ok := toFloat64(val); ok {
			sum += f
			count++
		}
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	}
	return &AggregationResult{
		Type:  AggTypeAvg,
		Value: avg,
	}
}

func computeSumAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	sum := 0.0
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if f, ok := toFloat64(val); ok {
			sum += f
		}
	}
	return &AggregationResult{
		Type:  AggTypeSum,
		Value: sum,
	}
}

func computeMinAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	min := math.Inf(1)
	found := false
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if f, ok := toFloat64(val); ok {
			if f < min {
				min = f
			}
			found = true
		}
	}
	if !found {
		min = 0
	}
	return &AggregationResult{
		Type:  AggTypeMin,
		Value: min,
	}
}

func computeMaxAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	max := math.Inf(-1)
	found := false
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if f, ok := toFloat64(val); ok {
			if f > max {
				max = f
			}
			found = true
		}
	}
	if !found {
		max = 0
	}
	return &AggregationResult{
		Type:  AggTypeMax,
		Value: max,
	}
}

func computeStatsAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	sum := 0.0
	min := math.Inf(1)
	max := math.Inf(-1)
	count := int64(0)
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if f, ok := toFloat64(val); ok {
			sum += f
			if f < min {
				min = f
			}
			if f > max {
				max = f
			}
			count++
		}
	}
	avg := 0.0
	if count > 0 {
		avg = sum / float64(count)
	} else {
		min = 0
		max = 0
	}
	return &AggregationResult{
		Type: AggTypeStats,
		Stats: &Stats{
			Count: count,
			Min:   min,
			Max:   max,
			Avg:   avg,
			Sum:   sum,
		},
	}
}

func computeCardinalityAgg(rows []map[string]interface{}, agg *Aggregation) *AggregationResult {
	unique := make(map[string]struct{})
	for _, row := range rows {
		val := resolveNestedField(row, agg.Field)
		if val != nil {
			unique[fmt.Sprintf("%v", val)] = struct{}{}
		}
	}
	return &AggregationResult{
		Type:  AggTypeCardinality,
		Value: float64(len(unique)),
	}
}

// resolveNestedField resolves dotted field paths like "cloud.region" in nested maps.
func resolveNestedField(doc map[string]interface{}, field string) interface{} {
	// Try direct lookup first (fastest path)
	if val, ok := doc[field]; ok {
		return val
	}

	// Try dotted path resolution
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

// parseInterval parses a date histogram interval string to a Duration.
func parseInterval(interval string) time.Duration {
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
		return 30 * 24 * time.Hour // approximate
	case "1q", "quarter":
		return 90 * 24 * time.Hour // approximate
	case "1y", "year":
		return 365 * 24 * time.Hour // approximate
	default:
		return time.Hour
	}
}

// truncateToInterval truncates a time to the nearest interval bucket.
func truncateToInterval(t time.Time, interval time.Duration) time.Time {
	if interval >= 24*time.Hour {
		// For day+ intervals, truncate to midnight UTC
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

// tryParseTime attempts to parse a time string in common formats.
func tryParseTime(s string) (time.Time, error) {
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
