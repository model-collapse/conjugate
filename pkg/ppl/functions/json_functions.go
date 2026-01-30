// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

// registerJSONFunctions registers JSON manipulation functions
func (r *FunctionRegistry) registerJSONFunctions() {
	jsonFunctions := []*FunctionInfo{
		// JSON extraction and parsing
		{
			PPLName:     "json_extract",
			UDFName:     "json_extract",
			Category:    "json",
			Description: "Extract value from JSON using JSONPath",
		},
		{
			PPLName:     "json_extract_scalar",
			UDFName:     "json_extract_scalar",
			Category:    "json",
			Description: "Extract scalar value from JSON",
		},

		// JSON validation
		{
			PPLName:     "json_valid",
			UDFName:     "json_valid",
			Category:    "json",
			Description: "Check if string is valid JSON",
			Aliases:     []string{"is_json"},
		},

		// JSON structure operations
		{
			PPLName:     "json_keys",
			UDFName:     "json_keys",
			Category:    "json",
			Description: "Get array of keys from JSON object",
		},
		{
			PPLName:     "json_values",
			UDFName:     "json_values",
			Category:    "json",
			Description: "Get array of values from JSON object",
		},
		{
			PPLName:     "json_length",
			UDFName:     "json_length",
			Category:    "json",
			Description: "Get number of elements in JSON array or object",
			Aliases:     []string{"json_array_length"},
		},

		// JSON construction
		{
			PPLName:     "json_array",
			UDFName:     "json_array",
			Category:    "json",
			Description: "Create JSON array from values",
		},
		{
			PPLName:     "json_object",
			UDFName:     "json_object",
			Category:    "json",
			Description: "Create JSON object from key-value pairs",
		},

		// JSON type checking
		{
			PPLName:     "json_type",
			UDFName:     "json_type",
			Category:    "json",
			Description: "Get type of JSON value (object, array, string, number, boolean, null)",
		},

		// JSON set operations
		{
			PPLName:     "json_set",
			UDFName:     "json_set",
			Category:    "json",
			Description: "Set value in JSON at specified path",
		},
		{
			PPLName:     "json_delete",
			UDFName:     "json_delete",
			Category:    "json",
			Description: "Delete value from JSON at specified path",
			Aliases:     []string{"json_remove"},
		},

		// JSON array operations
		{
			PPLName:     "json_array_contains",
			UDFName:     "json_array_contains",
			Category:    "json",
			Description: "Check if JSON array contains a value",
		},
		{
			PPLName:     "json_array_append",
			UDFName:     "json_array_append",
			Category:    "json",
			Description: "Append value to JSON array",
		},

		// JSON formatting
		{
			PPLName:     "json_format",
			UDFName:     "json_format",
			Category:    "json",
			Description: "Format JSON with indentation",
			Aliases:     []string{"json_pretty"},
		},
		{
			PPLName:     "json_compact",
			UDFName:     "json_compact",
			Category:    "json",
			Description: "Remove whitespace from JSON",
			Aliases:     []string{"json_minify"},
		},
	}

	for _, fn := range jsonFunctions {
		r.RegisterFunction(fn)
	}
}
