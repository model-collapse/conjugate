// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

import (
	"strings"
	"sync"
)

// FunctionInfo describes a PPL function and its WASM mapping
type FunctionInfo struct {
	PPLName     string   // PPL function name (e.g., "abs")
	UDFName     string   // WASM UDF name (e.g., "math_abs")
	Category    string   // Category (math, string, date, etc.)
	Description string   // Description
	Aliases     []string // Alternative names
}

// FunctionRegistry maps PPL functions to WASM UDFs
type FunctionRegistry struct {
	functions map[string]*FunctionInfo // lowercase name → info
	mu        sync.RWMutex
}

// Global registry instance
var globalRegistry *FunctionRegistry
var once sync.Once

// GetBuiltinRegistry returns the global built-in function registry
func GetBuiltinRegistry() *FunctionRegistry {
	once.Do(func() {
		globalRegistry = newBuiltinRegistry()
	})
	return globalRegistry
}

// newBuiltinRegistry creates a new registry with built-in functions
func newBuiltinRegistry() *FunctionRegistry {
	registry := &FunctionRegistry{
		functions: make(map[string]*FunctionInfo),
	}

	// Register all built-in functions
	registry.registerMathFunctions()
	registry.registerStringFunctions()
	registry.registerDateFunctions()
	registry.registerTypeFunctions()
	registry.registerConditionalFunctions()
	registry.registerRelevanceFunctions()
	registry.registerAggregationFunctions()

	return registry
}

// HasFunction checks if a function is in the registry
func (r *FunctionRegistry) HasFunction(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.functions[strings.ToLower(name)]
	return ok
}

// GetFunction retrieves function info
func (r *FunctionRegistry) GetFunction(name string) *FunctionInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.functions[strings.ToLower(name)]
}

// RegisterFunction registers a new function mapping
func (r *FunctionRegistry) RegisterFunction(info *FunctionInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := strings.ToLower(info.PPLName)
	r.functions[name] = info

	// Register aliases
	for _, alias := range info.Aliases {
		r.functions[strings.ToLower(alias)] = info
	}
}

// registerMathFunctions registers mathematical functions
func (r *FunctionRegistry) registerMathFunctions() {
	mathFunctions := []*FunctionInfo{
		// Basic math
		{
			PPLName:     "abs",
			UDFName:     "math_abs",
			Category:    "math",
			Description: "Absolute value",
		},
		{
			PPLName:     "ceil",
			UDFName:     "math_ceil",
			Category:    "math",
			Description: "Ceiling function (round up)",
			Aliases:     []string{"ceiling"},
		},
		{
			PPLName:     "floor",
			UDFName:     "math_floor",
			Category:    "math",
			Description: "Floor function (round down)",
		},
		{
			PPLName:     "round",
			UDFName:     "math_round",
			Category:    "math",
			Description: "Round to nearest integer",
		},
		{
			PPLName:     "sqrt",
			UDFName:     "math_sqrt",
			Category:    "math",
			Description: "Square root",
		},
		{
			PPLName:     "cbrt",
			UDFName:     "math_cbrt",
			Category:    "math",
			Description: "Cube root",
		},
		{
			PPLName:     "pow",
			UDFName:     "math_pow",
			Category:    "math",
			Description: "Power function (x^y)",
			Aliases:     []string{"power"},
		},
		{
			PPLName:     "mod",
			UDFName:     "math_mod",
			Category:    "math",
			Description: "Modulo (remainder)",
		},
		{
			PPLName:     "sign",
			UDFName:     "math_sign",
			Category:    "math",
			Description: "Sign of a number (-1, 0, 1)",
		},
		{
			PPLName:     "truncate",
			UDFName:     "math_truncate",
			Category:    "math",
			Description: "Truncate to specified decimal places",
			Aliases:     []string{"trunc"},
		},
		{
			PPLName:     "rand",
			UDFName:     "math_rand",
			Category:    "math",
			Description: "Random number between 0 and 1",
			Aliases:     []string{"random"},
		},

		// Logarithms
		{
			PPLName:     "log",
			UDFName:     "math_log",
			Category:    "math",
			Description: "Natural logarithm",
			Aliases:     []string{"ln"},
		},
		{
			PPLName:     "log10",
			UDFName:     "math_log10",
			Category:    "math",
			Description: "Base-10 logarithm",
		},
		{
			PPLName:     "log2",
			UDFName:     "math_log2",
			Category:    "math",
			Description: "Base-2 logarithm",
		},
		{
			PPLName:     "exp",
			UDFName:     "math_exp",
			Category:    "math",
			Description: "Exponential function (e^x)",
		},

		// Trigonometric functions
		{
			PPLName:     "sin",
			UDFName:     "math_sin",
			Category:    "math",
			Description: "Sine function",
		},
		{
			PPLName:     "cos",
			UDFName:     "math_cos",
			Category:    "math",
			Description: "Cosine function",
		},
		{
			PPLName:     "tan",
			UDFName:     "math_tan",
			Category:    "math",
			Description: "Tangent function",
		},
		{
			PPLName:     "asin",
			UDFName:     "math_asin",
			Category:    "math",
			Description: "Arc sine (inverse sine)",
		},
		{
			PPLName:     "acos",
			UDFName:     "math_acos",
			Category:    "math",
			Description: "Arc cosine (inverse cosine)",
		},
		{
			PPLName:     "atan",
			UDFName:     "math_atan",
			Category:    "math",
			Description: "Arc tangent (inverse tangent)",
		},
		{
			PPLName:     "atan2",
			UDFName:     "math_atan2",
			Category:    "math",
			Description: "Arc tangent of y/x with quadrant handling",
		},
		{
			PPLName:     "cot",
			UDFName:     "math_cot",
			Category:    "math",
			Description: "Cotangent function",
		},
		{
			PPLName:     "degrees",
			UDFName:     "math_degrees",
			Category:    "math",
			Description: "Convert radians to degrees",
		},
		{
			PPLName:     "radians",
			UDFName:     "math_radians",
			Category:    "math",
			Description: "Convert degrees to radians",
		},

		// Constants
		{
			PPLName:     "e",
			UDFName:     "math_e",
			Category:    "math",
			Description: "Euler's number (e ≈ 2.71828)",
		},
		{
			PPLName:     "pi",
			UDFName:     "math_pi",
			Category:    "math",
			Description: "Pi (π ≈ 3.14159)",
		},

		// Bitwise functions
		{
			PPLName:     "bit_and",
			UDFName:     "math_bit_and",
			Category:    "math",
			Description: "Bitwise AND",
			Aliases:     []string{"bitwise_and"},
		},
		{
			PPLName:     "bit_or",
			UDFName:     "math_bit_or",
			Category:    "math",
			Description: "Bitwise OR",
			Aliases:     []string{"bitwise_or"},
		},
		{
			PPLName:     "bit_xor",
			UDFName:     "math_bit_xor",
			Category:    "math",
			Description: "Bitwise XOR",
			Aliases:     []string{"bitwise_xor"},
		},
		{
			PPLName:     "bit_not",
			UDFName:     "math_bit_not",
			Category:    "math",
			Description: "Bitwise NOT",
			Aliases:     []string{"bitwise_not"},
		},
	}

	for _, fn := range mathFunctions {
		r.RegisterFunction(fn)
	}
}

// registerStringFunctions registers string manipulation functions
func (r *FunctionRegistry) registerStringFunctions() {
	stringFunctions := []*FunctionInfo{
		// Case conversion
		{
			PPLName:     "upper",
			UDFName:     "string_upper",
			Category:    "string",
			Description: "Convert to uppercase",
			Aliases:     []string{"ucase"},
		},
		{
			PPLName:     "lower",
			UDFName:     "string_lower",
			Category:    "string",
			Description: "Convert to lowercase",
			Aliases:     []string{"lcase"},
		},

		// Trimming
		{
			PPLName:     "trim",
			UDFName:     "string_trim",
			Category:    "string",
			Description: "Remove leading/trailing whitespace",
		},
		{
			PPLName:     "ltrim",
			UDFName:     "string_ltrim",
			Category:    "string",
			Description: "Remove leading whitespace",
		},
		{
			PPLName:     "rtrim",
			UDFName:     "string_rtrim",
			Category:    "string",
			Description: "Remove trailing whitespace",
		},

		// Length and substring
		{
			PPLName:     "length",
			UDFName:     "string_length",
			Category:    "string",
			Description: "String length",
			Aliases:     []string{"len", "char_length", "character_length"},
		},
		{
			PPLName:     "substring",
			UDFName:     "string_substring",
			Category:    "string",
			Description: "Extract substring",
			Aliases:     []string{"substr", "mid"},
		},
		{
			PPLName:     "left",
			UDFName:     "string_left",
			Category:    "string",
			Description: "Get leftmost characters",
		},
		{
			PPLName:     "right",
			UDFName:     "string_right",
			Category:    "string",
			Description: "Get rightmost characters",
		},

		// Concatenation and manipulation
		{
			PPLName:     "concat",
			UDFName:     "string_concat",
			Category:    "string",
			Description: "Concatenate strings",
		},
		{
			PPLName:     "concat_ws",
			UDFName:     "string_concat_ws",
			Category:    "string",
			Description: "Concatenate with separator",
		},
		{
			PPLName:     "replace",
			UDFName:     "string_replace",
			Category:    "string",
			Description: "Replace substring",
		},
		{
			PPLName:     "reverse",
			UDFName:     "string_reverse",
			Category:    "string",
			Description: "Reverse string",
		},
		{
			PPLName:     "repeat",
			UDFName:     "string_repeat",
			Category:    "string",
			Description: "Repeat string n times",
		},

		// Search functions
		{
			PPLName:     "locate",
			UDFName:     "string_locate",
			Category:    "string",
			Description: "Find position of substring",
			Aliases:     []string{"position", "instr"},
		},
		{
			PPLName:     "split",
			UDFName:     "string_split",
			Category:    "string",
			Description: "Split string by delimiter",
		},

		// Pattern matching
		{
			PPLName:     "regexp",
			UDFName:     "string_regexp",
			Category:    "string",
			Description: "Regular expression match",
			Aliases:     []string{"regex", "regexp_like"},
		},
		{
			PPLName:     "regexp_replace",
			UDFName:     "string_regexp_replace",
			Category:    "string",
			Description: "Regular expression replace",
		},
		{
			PPLName:     "regexp_extract",
			UDFName:     "string_regexp_extract",
			Category:    "string",
			Description: "Extract regex match groups",
		},
		{
			PPLName:     "like",
			UDFName:     "string_like",
			Category:    "string",
			Description: "SQL LIKE pattern matching",
		},

		// Padding
		{
			PPLName:     "lpad",
			UDFName:     "string_lpad",
			Category:    "string",
			Description: "Left pad string",
		},
		{
			PPLName:     "rpad",
			UDFName:     "string_rpad",
			Category:    "string",
			Description: "Right pad string",
		},

		// ASCII functions
		{
			PPLName:     "ascii",
			UDFName:     "string_ascii",
			Category:    "string",
			Description: "ASCII code of first character",
		},
		{
			PPLName:     "chr",
			UDFName:     "string_chr",
			Category:    "string",
			Description: "Character from ASCII code",
			Aliases:     []string{"char"},
		},
	}

	for _, fn := range stringFunctions {
		r.RegisterFunction(fn)
	}
}

// registerDateFunctions registers date/time functions
func (r *FunctionRegistry) registerDateFunctions() {
	dateFunctions := []*FunctionInfo{
		// Extraction functions
		{
			PPLName:     "year",
			UDFName:     "date_year",
			Category:    "date",
			Description: "Extract year from date",
		},
		{
			PPLName:     "month",
			UDFName:     "date_month",
			Category:    "date",
			Description: "Extract month from date",
		},
		{
			PPLName:     "day",
			UDFName:     "date_day",
			Category:    "date",
			Description: "Extract day from date",
			Aliases:     []string{"dayofmonth"},
		},
		{
			PPLName:     "hour",
			UDFName:     "date_hour",
			Category:    "date",
			Description: "Extract hour from datetime",
		},
		{
			PPLName:     "minute",
			UDFName:     "date_minute",
			Category:    "date",
			Description: "Extract minute from datetime",
		},
		{
			PPLName:     "second",
			UDFName:     "date_second",
			Category:    "date",
			Description: "Extract second from datetime",
		},
		{
			PPLName:     "microsecond",
			UDFName:     "date_microsecond",
			Category:    "date",
			Description: "Extract microsecond from datetime",
		},
		{
			PPLName:     "dayofweek",
			UDFName:     "date_dayofweek",
			Category:    "date",
			Description: "Day of week (1=Sunday, 7=Saturday)",
			Aliases:     []string{"dow"},
		},
		{
			PPLName:     "dayofyear",
			UDFName:     "date_dayofyear",
			Category:    "date",
			Description: "Day of year (1-366)",
			Aliases:     []string{"doy"},
		},
		{
			PPLName:     "weekofyear",
			UDFName:     "date_weekofyear",
			Category:    "date",
			Description: "Week of year (1-53)",
			Aliases:     []string{"week"},
		},
		{
			PPLName:     "quarter",
			UDFName:     "date_quarter",
			Category:    "date",
			Description: "Quarter of year (1-4)",
		},
		{
			PPLName:     "dayname",
			UDFName:     "date_dayname",
			Category:    "date",
			Description: "Name of the day (Monday, Tuesday, etc.)",
		},
		{
			PPLName:     "monthname",
			UDFName:     "date_monthname",
			Category:    "date",
			Description: "Name of the month (January, February, etc.)",
		},

		// Current date/time functions
		{
			PPLName:     "now",
			UDFName:     "date_now",
			Category:    "date",
			Description: "Current datetime",
			Aliases:     []string{"current_timestamp"},
		},
		{
			PPLName:     "curdate",
			UDFName:     "date_curdate",
			Category:    "date",
			Description: "Current date",
			Aliases:     []string{"current_date"},
		},
		{
			PPLName:     "curtime",
			UDFName:     "date_curtime",
			Category:    "date",
			Description: "Current time",
			Aliases:     []string{"current_time"},
		},
		{
			PPLName:     "sysdate",
			UDFName:     "date_sysdate",
			Category:    "date",
			Description: "System datetime",
		},
		{
			PPLName:     "utc_date",
			UDFName:     "date_utc_date",
			Category:    "date",
			Description: "Current UTC date",
		},
		{
			PPLName:     "utc_time",
			UDFName:     "date_utc_time",
			Category:    "date",
			Description: "Current UTC time",
		},
		{
			PPLName:     "utc_timestamp",
			UDFName:     "date_utc_timestamp",
			Category:    "date",
			Description: "Current UTC datetime",
		},

		// Date construction
		{
			PPLName:     "date",
			UDFName:     "date_date",
			Category:    "date",
			Description: "Extract date part from datetime",
		},
		{
			PPLName:     "time",
			UDFName:     "date_time",
			Category:    "date",
			Description: "Extract time part from datetime",
		},
		{
			PPLName:     "makedate",
			UDFName:     "date_makedate",
			Category:    "date",
			Description: "Create date from year and day of year",
		},
		{
			PPLName:     "maketime",
			UDFName:     "date_maketime",
			Category:    "date",
			Description: "Create time from hour, minute, second",
		},

		// Date arithmetic
		{
			PPLName:     "date_add",
			UDFName:     "date_add",
			Category:    "date",
			Description: "Add interval to date",
			Aliases:     []string{"adddate"},
		},
		{
			PPLName:     "date_sub",
			UDFName:     "date_sub",
			Category:    "date",
			Description: "Subtract interval from date",
			Aliases:     []string{"subdate"},
		},
		{
			PPLName:     "addtime",
			UDFName:     "date_addtime",
			Category:    "date",
			Description: "Add time to datetime",
		},
		{
			PPLName:     "subtime",
			UDFName:     "date_subtime",
			Category:    "date",
			Description: "Subtract time from datetime",
		},
		{
			PPLName:     "datediff",
			UDFName:     "date_datediff",
			Category:    "date",
			Description: "Difference in days between dates",
		},
		{
			PPLName:     "timediff",
			UDFName:     "date_timediff",
			Category:    "date",
			Description: "Difference between times",
		},
		{
			PPLName:     "timestampdiff",
			UDFName:     "date_timestampdiff",
			Category:    "date",
			Description: "Difference in specified units",
		},
		{
			PPLName:     "period_add",
			UDFName:     "date_period_add",
			Category:    "date",
			Description: "Add months to period (YYYYMM)",
		},
		{
			PPLName:     "period_diff",
			UDFName:     "date_period_diff",
			Category:    "date",
			Description: "Difference between periods (YYYYMM)",
		},

		// Date conversion
		{
			PPLName:     "from_days",
			UDFName:     "date_from_days",
			Category:    "date",
			Description: "Date from day number",
		},
		{
			PPLName:     "to_days",
			UDFName:     "date_to_days",
			Category:    "date",
			Description: "Day number from date",
		},
		{
			PPLName:     "to_seconds",
			UDFName:     "date_to_seconds",
			Category:    "date",
			Description: "Seconds since year 0",
		},
		{
			PPLName:     "from_unixtime",
			UDFName:     "date_from_unixtime",
			Category:    "date",
			Description: "Datetime from Unix timestamp",
		},
		{
			PPLName:     "unix_timestamp",
			UDFName:     "date_unix_timestamp",
			Category:    "date",
			Description: "Unix timestamp from datetime",
		},

		// Date utilities
		{
			PPLName:     "last_day",
			UDFName:     "date_last_day",
			Category:    "date",
			Description: "Last day of the month",
		},
		{
			PPLName:     "convert_tz",
			UDFName:     "date_convert_tz",
			Category:    "date",
			Description: "Convert between time zones",
		},

		// Formatting
		{
			PPLName:     "date_format",
			UDFName:     "date_format",
			Category:    "date",
			Description: "Format date according to pattern",
		},
		{
			PPLName:     "str_to_date",
			UDFName:     "date_str_to_date",
			Category:    "date",
			Description: "Parse string to date",
		},
		{
			PPLName:     "time_format",
			UDFName:     "date_time_format",
			Category:    "date",
			Description: "Format time according to pattern",
		},
	}

	for _, fn := range dateFunctions {
		r.RegisterFunction(fn)
	}
}

// registerTypeFunctions registers type conversion functions
func (r *FunctionRegistry) registerTypeFunctions() {
	typeFunctions := []*FunctionInfo{
		{
			PPLName:     "int",
			UDFName:     "type_int",
			Category:    "type",
			Description: "Convert to integer",
			Aliases:     []string{"toint"},
		},
		{
			PPLName:     "long",
			UDFName:     "type_long",
			Category:    "type",
			Description: "Convert to long integer",
			Aliases:     []string{"tolong"},
		},
		{
			PPLName:     "float",
			UDFName:     "type_float",
			Category:    "type",
			Description: "Convert to float",
			Aliases:     []string{"tofloat"},
		},
		{
			PPLName:     "double",
			UDFName:     "type_double",
			Category:    "type",
			Description: "Convert to double",
			Aliases:     []string{"todouble"},
		},
		{
			PPLName:     "string",
			UDFName:     "type_string",
			Category:    "type",
			Description: "Convert to string",
			Aliases:     []string{"tostring"},
		},
		{
			PPLName:     "bool",
			UDFName:     "type_bool",
			Category:    "type",
			Description: "Convert to boolean",
			Aliases:     []string{"tobool"},
		},
		{
			PPLName:     "cast",
			UDFName:     "type_cast",
			Category:    "type",
			Description: "Cast to specified type",
		},
		{
			PPLName:     "convert",
			UDFName:     "type_convert",
			Category:    "type",
			Description: "Convert to specified type",
		},
		{
			PPLName:     "try_cast",
			UDFName:     "type_try_cast",
			Category:    "type",
			Description: "Cast to type, returning null on failure",
		},
		{
			PPLName:     "typeof",
			UDFName:     "type_typeof",
			Category:    "type",
			Description: "Get type name of value",
		},
	}

	for _, fn := range typeFunctions {
		r.RegisterFunction(fn)
	}
}

// registerConditionalFunctions registers conditional and null-handling functions
func (r *FunctionRegistry) registerConditionalFunctions() {
	conditionalFunctions := []*FunctionInfo{
		// Null handling
		{
			PPLName:     "isnull",
			UDFName:     "cond_isnull",
			Category:    "conditional",
			Description: "Check if value is null",
		},
		{
			PPLName:     "isnotnull",
			UDFName:     "cond_isnotnull",
			Category:    "conditional",
			Description: "Check if value is not null",
		},
		{
			PPLName:     "ifnull",
			UDFName:     "cond_ifnull",
			Category:    "conditional",
			Description: "Return alternative if null",
			Aliases:     []string{"nvl"},
		},
		{
			PPLName:     "nvl2",
			UDFName:     "cond_nvl2",
			Category:    "conditional",
			Description: "Return one of two values based on null",
		},
		{
			PPLName:     "nullif",
			UDFName:     "cond_nullif",
			Category:    "conditional",
			Description: "Return null if values are equal",
		},
		{
			PPLName:     "coalesce",
			UDFName:     "cond_coalesce",
			Category:    "conditional",
			Description: "Return first non-null value",
		},

		// Conditional logic
		{
			PPLName:     "if",
			UDFName:     "cond_if",
			Category:    "conditional",
			Description: "If-then-else expression",
		},
		{
			PPLName:     "case",
			UDFName:     "cond_case",
			Category:    "conditional",
			Description: "Multi-way conditional",
		},
		{
			PPLName:     "greatest",
			UDFName:     "cond_greatest",
			Category:    "conditional",
			Description: "Return maximum value",
		},
		{
			PPLName:     "least",
			UDFName:     "cond_least",
			Category:    "conditional",
			Description: "Return minimum value",
		},

		// Comparison utilities
		{
			PPLName:     "in",
			UDFName:     "cond_in",
			Category:    "conditional",
			Description: "Check if value is in list",
		},
		{
			PPLName:     "between",
			UDFName:     "cond_between",
			Category:    "conditional",
			Description: "Check if value is between bounds",
		},
	}

	for _, fn := range conditionalFunctions {
		r.RegisterFunction(fn)
	}
}

// registerRelevanceFunctions registers OpenSearch relevance functions
func (r *FunctionRegistry) registerRelevanceFunctions() {
	relevanceFunctions := []*FunctionInfo{
		{
			PPLName:     "match",
			UDFName:     "rel_match",
			Category:    "relevance",
			Description: "Full-text match query",
		},
		{
			PPLName:     "match_phrase",
			UDFName:     "rel_match_phrase",
			Category:    "relevance",
			Description: "Match phrase query",
		},
		{
			PPLName:     "match_phrase_prefix",
			UDFName:     "rel_match_phrase_prefix",
			Category:    "relevance",
			Description: "Match phrase with prefix",
		},
		{
			PPLName:     "match_bool_prefix",
			UDFName:     "rel_match_bool_prefix",
			Category:    "relevance",
			Description: "Match with boolean prefix",
		},
		{
			PPLName:     "multi_match",
			UDFName:     "rel_multi_match",
			Category:    "relevance",
			Description: "Match across multiple fields",
		},
		{
			PPLName:     "query_string",
			UDFName:     "rel_query_string",
			Category:    "relevance",
			Description: "Lucene query string syntax",
		},
		{
			PPLName:     "simple_query_string",
			UDFName:     "rel_simple_query_string",
			Category:    "relevance",
			Description: "Simplified query string syntax",
		},
	}

	for _, fn := range relevanceFunctions {
		r.RegisterFunction(fn)
	}
}

// registerAggregationFunctions registers aggregation functions
func (r *FunctionRegistry) registerAggregationFunctions() {
	aggFunctions := []*FunctionInfo{
		// Basic aggregations
		{
			PPLName:     "count",
			UDFName:     "agg_count",
			Category:    "aggregation",
			Description: "Count values",
		},
		{
			PPLName:     "sum",
			UDFName:     "agg_sum",
			Category:    "aggregation",
			Description: "Sum of values",
		},
		{
			PPLName:     "avg",
			UDFName:     "agg_avg",
			Category:    "aggregation",
			Description: "Average of values",
			Aliases:     []string{"mean"},
		},
		{
			PPLName:     "min",
			UDFName:     "agg_min",
			Category:    "aggregation",
			Description: "Minimum value",
		},
		{
			PPLName:     "max",
			UDFName:     "agg_max",
			Category:    "aggregation",
			Description: "Maximum value",
		},

		// Statistical aggregations
		{
			PPLName:     "stddev",
			UDFName:     "agg_stddev",
			Category:    "aggregation",
			Description: "Standard deviation",
			Aliases:     []string{"stdev", "stddev_samp"},
		},
		{
			PPLName:     "stddev_pop",
			UDFName:     "agg_stddev_pop",
			Category:    "aggregation",
			Description: "Population standard deviation",
		},
		{
			PPLName:     "variance",
			UDFName:     "agg_variance",
			Category:    "aggregation",
			Description: "Variance",
			Aliases:     []string{"var_samp"},
		},
		{
			PPLName:     "var_pop",
			UDFName:     "agg_var_pop",
			Category:    "aggregation",
			Description: "Population variance",
		},

		// Distinct/cardinality
		{
			PPLName:     "distinct_count",
			UDFName:     "agg_distinct_count",
			Category:    "aggregation",
			Description: "Count of distinct values",
			Aliases:     []string{"dc", "cardinality"},
		},
		{
			PPLName:     "approx_count_distinct",
			UDFName:     "agg_approx_count_distinct",
			Category:    "aggregation",
			Description: "Approximate distinct count (HyperLogLog)",
		},

		// Percentile aggregations
		{
			PPLName:     "percentile",
			UDFName:     "agg_percentile",
			Category:    "aggregation",
			Description: "Percentile value",
		},
		{
			PPLName:     "percentile_approx",
			UDFName:     "agg_percentile_approx",
			Category:    "aggregation",
			Description: "Approximate percentile (t-digest)",
		},
		{
			PPLName:     "median",
			UDFName:     "agg_median",
			Category:    "aggregation",
			Description: "Median value (50th percentile)",
		},

		// Collection aggregations
		{
			PPLName:     "values",
			UDFName:     "agg_values",
			Category:    "aggregation",
			Description: "List of values",
		},
		{
			PPLName:     "list",
			UDFName:     "agg_list",
			Category:    "aggregation",
			Description: "Collect values into list",
		},
		{
			PPLName:     "first",
			UDFName:     "agg_first",
			Category:    "aggregation",
			Description: "First value",
		},
		{
			PPLName:     "last",
			UDFName:     "agg_last",
			Category:    "aggregation",
			Description: "Last value",
		},

		// Earliest/latest for time series
		{
			PPLName:     "earliest",
			UDFName:     "agg_earliest",
			Category:    "aggregation",
			Description: "Earliest value by time",
		},
		{
			PPLName:     "latest",
			UDFName:     "agg_latest",
			Category:    "aggregation",
			Description: "Latest value by time",
		},
	}

	for _, fn := range aggFunctions {
		r.RegisterFunction(fn)
	}
}

// ListFunctions returns all registered functions
func (r *FunctionRegistry) ListFunctions() []*FunctionInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Deduplicate (aliases point to same info)
	seen := make(map[*FunctionInfo]bool)
	result := make([]*FunctionInfo, 0, len(r.functions))

	for _, info := range r.functions {
		if !seen[info] {
			seen[info] = true
			result = append(result, info)
		}
	}

	return result
}
