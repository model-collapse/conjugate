// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

// registerCryptoFunctions registers cryptographic and encoding functions
func (r *FunctionRegistry) registerCryptoFunctions() {
	cryptoFunctions := []*FunctionInfo{
		// Hash functions
		{
			PPLName:     "md5",
			UDFName:     "crypto_md5",
			Category:    "cryptographic",
			Description: "Calculate MD5 hash of a string",
		},
		{
			PPLName:     "sha1",
			UDFName:     "crypto_sha1",
			Category:    "cryptographic",
			Description: "Calculate SHA-1 hash of a string",
		},
		{
			PPLName:     "sha256",
			UDFName:     "crypto_sha256",
			Category:    "cryptographic",
			Description: "Calculate SHA-256 hash of a string",
			Aliases:     []string{"sha2"},
		},
		{
			PPLName:     "sha512",
			UDFName:     "crypto_sha512",
			Category:    "cryptographic",
			Description: "Calculate SHA-512 hash of a string",
		},

		// Base64 encoding
		{
			PPLName:     "base64",
			UDFName:     "crypto_base64_encode",
			Category:    "cryptographic",
			Description: "Encode string to base64",
			Aliases:     []string{"base64encode"},
		},
		{
			PPLName:     "base64decode",
			UDFName:     "crypto_base64_decode",
			Category:    "cryptographic",
			Description: "Decode base64 string",
			Aliases:     []string{"unbase64"},
		},

		// URL encoding
		{
			PPLName:     "urlencode",
			UDFName:     "crypto_url_encode",
			Category:    "cryptographic",
			Description: "URL encode string",
		},
		{
			PPLName:     "urldecode",
			UDFName:     "crypto_url_decode",
			Category:    "cryptographic",
			Description: "URL decode string",
		},

		// Hex encoding
		{
			PPLName:     "hex",
			UDFName:     "crypto_hex_encode",
			Category:    "cryptographic",
			Description: "Encode string to hexadecimal",
		},
		{
			PPLName:     "unhex",
			UDFName:     "crypto_hex_decode",
			Category:    "cryptographic",
			Description: "Decode hexadecimal string",
		},
	}

	for _, fn := range cryptoFunctions {
		r.RegisterFunction(fn)
	}
}
