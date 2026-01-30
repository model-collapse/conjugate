// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

// registerSystemFunctions registers system information functions
func (r *FunctionRegistry) registerSystemFunctions() {
	systemFunctions := []*FunctionInfo{
		// System information
		{
			PPLName:     "version",
			UDFName:     "system_version",
			Category:    "system",
			Description: "Return CONJUGATE version string",
		},
		{
			PPLName:     "database",
			UDFName:     "system_database",
			Category:    "system",
			Description: "Return current database/index name",
		},
		{
			PPLName:     "user",
			UDFName:     "system_user",
			Category:    "system",
			Description: "Return current user name",
			Aliases:     []string{"current_user"},
		},

		// Session information
		{
			PPLName:     "connection_id",
			UDFName:     "system_connection_id",
			Category:    "system",
			Description: "Return current connection ID",
		},
		{
			PPLName:     "session_user",
			UDFName:     "system_session_user",
			Category:    "system",
			Description: "Return session user name",
		},

		// System constants
		{
			PPLName:     "null",
			UDFName:     "system_null",
			Category:    "system",
			Description: "Return NULL value",
		},

		// Environment
		{
			PPLName:     "current_catalog",
			UDFName:     "system_current_catalog",
			Category:    "system",
			Description: "Return current catalog name",
		},
		{
			PPLName:     "current_schema",
			UDFName:     "system_current_schema",
			Category:    "system",
			Description: "Return current schema name",
		},
	}

	for _, fn := range systemFunctions {
		r.RegisterFunction(fn)
	}
}
