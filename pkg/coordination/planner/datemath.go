package planner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DateMathParser parses and evaluates OpenSearch date math expressions
// Examples: "now", "now-1d", "now+3h", "now/d", "2024-01-01||+1M/M"
type DateMathParser struct {
	nowFunc func() time.Time // Allows testing with fixed time
}

// NewDateMathParser creates a new date math parser
func NewDateMathParser() *DateMathParser {
	return &DateMathParser{
		nowFunc: time.Now,
	}
}

// Parse converts a date math expression to an absolute timestamp
// Returns the original value if it's not a date math expression
func (p *DateMathParser) Parse(expr interface{}) interface{} {
	// Only handle string expressions
	str, ok := expr.(string)
	if !ok {
		return expr
	}

	// Check if it's a date math expression
	if !p.isDateMath(str) {
		return expr
	}

	// Parse and evaluate
	result, err := p.evaluate(str)
	if err != nil {
		// If parsing fails, return original (let Diagon handle it)
		return expr
	}

	// Return ISO8601 format timestamp
	return result.UTC().Format("2006-01-02T15:04:05.000Z")
}

// isDateMath checks if a string looks like a date math expression
func (p *DateMathParser) isDateMath(str string) bool {
	// Fast checks for common patterns
	if strings.HasPrefix(str, "now") {
		return true
	}
	if strings.Contains(str, "||") {
		return true
	}
	return false
}

// evaluate parses and evaluates a date math expression
func (p *DateMathParser) evaluate(expr string) (time.Time, error) {
	// Handle "now" based expressions
	if strings.HasPrefix(expr, "now") {
		return p.evaluateNowExpression(expr)
	}

	// Handle anchor||expression format (e.g., "2024-01-01||+1d")
	if strings.Contains(expr, "||") {
		return p.evaluateAnchoredExpression(expr)
	}

	// Try parsing as ISO8601
	t, err := time.Parse(time.RFC3339, expr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unsupported date math expression: %s", expr)
}

// evaluateNowExpression handles "now", "now-1d", "now+3h/d" etc.
func (p *DateMathParser) evaluateNowExpression(expr string) (time.Time, error) {
	t := p.nowFunc().UTC()

	// Just "now"
	if expr == "now" {
		return t, nil
	}

	// Remove "now" prefix
	remainder := strings.TrimPrefix(expr, "now")

	// Apply operations
	return p.applyOperations(t, remainder)
}

// evaluateAnchoredExpression handles "2024-01-01||+1d" format
func (p *DateMathParser) evaluateAnchoredExpression(expr string) (time.Time, error) {
	parts := strings.Split(expr, "||")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid anchored expression: %s", expr)
	}

	// Parse anchor time
	t, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		// Try other common formats
		t, err = time.Parse("2006-01-02", parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid anchor date: %s", parts[0])
		}
	}

	// Apply operations
	return p.applyOperations(t, parts[1])
}

// applyOperations applies +/- operations and rounding
// Examples: "-1d", "+3h/d", "-1M/M+1d"
func (p *DateMathParser) applyOperations(t time.Time, ops string) (time.Time, error) {
	if ops == "" {
		return t, nil
	}

	// Regex to match operations: (+/-)?(number)(unit)(rounding)?
	// Examples: "-1d", "+3h", "/d", "+1M/M"
	re := regexp.MustCompile(`([+\-/])?(\d*)([yMwdhHms])`)
	matches := re.FindAllStringSubmatch(ops, -1)

	for _, match := range matches {
		op := match[1]     // "+", "-", "/" or ""
		numStr := match[2] // Number (may be empty for rounding)
		unit := match[3]   // Unit: y, M, w, d, h, H, m, s

		if op == "/" || op == "" && numStr == "" {
			// Rounding operation (e.g., "/d")
			t = p.roundDown(t, unit)
		} else {
			// Addition/subtraction
			num := 1
			if numStr != "" {
				num, _ = strconv.Atoi(numStr)
			}

			if op == "-" {
				num = -num
			}

			t = p.addDuration(t, num, unit)
		}
	}

	return t, nil
}

// roundDown rounds time down to the specified unit
func (p *DateMathParser) roundDown(t time.Time, unit string) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	switch unit {
	case "y": // Round to start of year
		return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
	case "M": // Round to start of month
		return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	case "w": // Round to start of week (Monday)
		weekday := t.Weekday()
		daysBack := int(weekday - time.Monday)
		if daysBack < 0 {
			daysBack += 7
		}
		return time.Date(year, month, day-daysBack, 0, 0, 0, 0, t.Location())
	case "d": // Round to start of day
		return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	case "h", "H": // Round to start of hour
		return time.Date(year, month, day, hour, 0, 0, 0, t.Location())
	case "m": // Round to start of minute
		return time.Date(year, month, day, hour, min, 0, 0, t.Location())
	case "s": // Round to start of second
		return time.Date(year, month, day, hour, min, sec, 0, t.Location())
	default:
		return t
	}
}

// addDuration adds a duration to time
func (p *DateMathParser) addDuration(t time.Time, num int, unit string) time.Time {
	switch unit {
	case "y": // Years
		return t.AddDate(num, 0, 0)
	case "M": // Months
		return t.AddDate(0, num, 0)
	case "w": // Weeks
		return t.AddDate(0, 0, num*7)
	case "d": // Days
		return t.AddDate(0, 0, num)
	case "h", "H": // Hours
		return t.Add(time.Duration(num) * time.Hour)
	case "m": // Minutes
		return t.Add(time.Duration(num) * time.Minute)
	case "s": // Seconds
		return t.Add(time.Duration(num) * time.Second)
	default:
		return t
	}
}

// ParseRangeParams processes range query parameters and converts date math
func (p *DateMathParser) ParseRangeParams(params map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range params {
		result[key] = p.Parse(value)
	}
	return result
}
