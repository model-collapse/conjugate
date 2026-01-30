// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

// registerIPFunctions registers IP address manipulation and validation functions
func (r *FunctionRegistry) registerIPFunctions() {
	ipFunctions := []*FunctionInfo{
		// IP validation
		{
			PPLName:     "isValidIP",
			UDFName:     "ip_is_valid",
			Category:    "ip",
			Description: "Check if string is a valid IP address (IPv4 or IPv6)",
		},
		{
			PPLName:     "isValidIPv4",
			UDFName:     "ip_is_valid_ipv4",
			Category:    "ip",
			Description: "Check if string is a valid IPv4 address",
		},
		{
			PPLName:     "isValidIPv6",
			UDFName:     "ip_is_valid_ipv6",
			Category:    "ip",
			Description: "Check if string is a valid IPv6 address",
		},

		// IP classification
		{
			PPLName:     "isPrivateIP",
			UDFName:     "ip_is_private",
			Category:    "ip",
			Description: "Check if IP address is in a private range (RFC 1918)",
		},
		{
			PPLName:     "isPublicIP",
			UDFName:     "ip_is_public",
			Category:    "ip",
			Description: "Check if IP address is in a public range",
		},
		{
			PPLName:     "isLoopbackIP",
			UDFName:     "ip_is_loopback",
			Category:    "ip",
			Description: "Check if IP address is a loopback address",
		},
		{
			PPLName:     "isMulticastIP",
			UDFName:     "ip_is_multicast",
			Category:    "ip",
			Description: "Check if IP address is a multicast address",
		},

		// CIDR operations
		{
			PPLName:     "cidr",
			UDFName:     "ip_cidr_match",
			Category:    "ip",
			Description: "Check if IP address matches CIDR notation",
		},
		{
			PPLName:     "cidrContains",
			UDFName:     "ip_cidr_contains",
			Category:    "ip",
			Description: "Check if CIDR range contains IP address",
		},

		// IP manipulation
		{
			PPLName:     "ipToInt",
			UDFName:     "ip_to_int",
			Category:    "ip",
			Description: "Convert IPv4 address to integer",
		},
		{
			PPLName:     "intToIP",
			UDFName:     "ip_from_int",
			Category:    "ip",
			Description: "Convert integer to IPv4 address",
		},

		// IP network operations
		{
			PPLName:     "ipNetwork",
			UDFName:     "ip_network",
			Category:    "ip",
			Description: "Get network address from IP and netmask",
		},
		{
			PPLName:     "ipBroadcast",
			UDFName:     "ip_broadcast",
			Category:    "ip",
			Description: "Get broadcast address from IP and netmask",
		},
		{
			PPLName:     "ipNetmask",
			UDFName:     "ip_netmask",
			Category:    "ip",
			Description: "Get netmask from CIDR notation",
		},

		// IP range operations
		{
			PPLName:     "ipRange",
			UDFName:     "ip_in_range",
			Category:    "ip",
			Description: "Check if IP is within start and end range",
		},
	}

	for _, fn := range ipFunctions {
		r.RegisterFunction(fn)
	}
}
