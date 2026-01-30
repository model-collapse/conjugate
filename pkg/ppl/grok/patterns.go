// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package grok

// Base patterns - building blocks for more complex patterns
var basePatterns = map[string]string{
	// Numbers
	"INT":        `(?:[+-]?(?:[0-9]+))`,
	"BASE10NUM":  `(?:[+-]?(?:[0-9]+(?:\.[0-9]+)?))`,
	"NUMBER":     `(?:%{BASE10NUM})`,
	"BASE16NUM":  `(?:0[xX]?[0-9a-fA-F]+)`,
	"BASE16FLOAT": `\b(?:0[xX]?[0-9a-fA-F]+(?:\.[0-9a-fA-F]*)?(?:[pP][+-]?[0-9]+)?)\b`,
	"POSINT":     `\b(?:[1-9][0-9]*)\b`,
	"NONNEGINT":  `\b(?:[0-9]+)\b`,

	// Words and identifiers
	"WORD":       `\b\w+\b`,
	"NOTSPACE":   `\S+`,
	"SPACE":      `\s*`,
	"DATA":       `.*?`,
	"GREEDYDATA": `.*`,
	"QUOTEDSTRING": `"(?:[^"\\]*(?:\\.[^"\\]*)*)"|'(?:[^'\\]*(?:\\.[^'\\]*)*)'`,

	// Network - IP addresses
	"IPV6":        `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?`,
	"IPV4":        `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`,
	"IP":          `(?:%{IPV6}|%{IPV4})`,
	"HOSTNAME":    `\b(?:[0-9A-Za-z][0-9A-Za-z-]{0,62})(?:\.(?:[0-9A-Za-z][0-9A-Za-z-]{0,62}))*(?:\.?|\b)`,
	"IPORHOST":    `(?:%{IP}|%{HOSTNAME})`,
	"HOSTPORT":    `%{IPORHOST}:%{POSINT}`,

	// Paths and URIs
	"PATH":        `(?:%{UNIXPATH}|%{WINPATH})`,
	"UNIXPATH":    `(?:/[A-Za-z0-9_./-]+)+`,
	"TTY":         `(?:/dev/(?:pts|tty(?:[pq])?)(?:\w+)?/?(?:[0-9]+))`,
	"WINPATH":     `(?:[A-Za-z]:|\\)(?:\\[^\\?*]*)+`,
	"URIPROTO":    `[A-Za-z][A-Za-z0-9+\-.]+`,
	"URIHOST":     `%{IPORHOST}(?::%{POSINT})?`,
	"URIPATH":     `(?:/[A-Za-z0-9$.+!*'(){},~:;=@#%&_\-]*)+`,
	"URIPARAM":    `\?[A-Za-z0-9$.+!*'|(){},~@#%&/=:;_?\-\[\]<>]*`,
	"URIPATHPARAM": `%{URIPATH}(?:%{URIPARAM})?`,
	"URI":         `%{URIPROTO}://(?:%{USER}(?::[^@]*)?@)?(?:%{URIHOST})?(?:%{URIPATHPARAM})?`,

	// Users and auth
	"USER":        `[a-zA-Z0-9._-]+`,
	"EMAILLOCALPART": `[a-zA-Z0-9!#$%&'*+/=?^_\x60{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+/=?^_\x60{|}~-]+)*`,
	"EMAILADDRESS": `%{EMAILLOCALPART}@%{HOSTNAME}`,

	// MAC address
	"MAC":         `(?:(?:[A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2})`,
	"CISCOMAC":    `(?:(?:[A-Fa-f0-9]{4}\.){2}[A-Fa-f0-9]{4})`,
	"WINDOWSMAC":  `(?:(?:[A-Fa-f0-9]{2}-){5}[A-Fa-f0-9]{2})`,
	"COMMONMAC":   `(?:%{MAC}|%{CISCOMAC}|%{WINDOWSMAC})`,

	// Date/Time
	"MONTHNUM":    `(?:0?[1-9]|1[0-2])`,
	"MONTHNUM2":   `(?:0[1-9]|1[0-2])`,
	"MONTHDAY":    `(?:(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9])`,
	"DAY":         `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)`,
	"MONTH":       `(?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)`,
	"YEAR":        `(?:\d\d){1,2}`,
	"HOUR":        `(?:2[0123]|[01]?[0-9])`,
	"MINUTE":      `(?:[0-5][0-9])`,
	"SECOND":      `(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)`,
	"TIME":        `%{HOUR}:%{MINUTE}(?::%{SECOND})?`,
	"DATE_US":     `%{MONTHNUM}[/-]%{MONTHDAY}[/-]%{YEAR}`,
	"DATE_EU":     `%{MONTHDAY}[./-]%{MONTHNUM}[./-]%{YEAR}`,
	"ISO8601_TIMEZONE": `(?:Z|[+-]%{HOUR}(?::?%{MINUTE}))`,
	"ISO8601_SECOND":   `(?:%{SECOND}|60)`,
	"TIMESTAMP_ISO8601": `%{YEAR}-%{MONTHNUM}-%{MONTHDAY}[T ]%{HOUR}:?%{MINUTE}(?::?%{SECOND})?%{ISO8601_TIMEZONE}?`,
	"DATE":             `%{DATE_US}|%{DATE_EU}`,
	"DATESTAMP":        `%{DATE}[- ]%{TIME}`,
	"TZ":               `(?:[APMCE][SD]T|UTC)`,
	"DATESTAMP_RFC822":  `%{DAY} %{MONTH} %{MONTHDAY} %{YEAR} %{TIME} %{TZ}`,
	"DATESTAMP_RFC2822": `%{DAY}, %{MONTHDAY} %{MONTH} %{YEAR} %{TIME} %{ISO8601_TIMEZONE}`,
	"DATESTAMP_OTHER":   `%{DAY} %{MONTH} %{MONTHDAY} %{TIME} %{TZ} %{YEAR}`,
	"DATESTAMP_EVENTLOG": `%{YEAR}%{MONTHNUM2}%{MONTHDAY}%{HOUR}%{MINUTE}%{SECOND}`,

	// Syslog
	"SYSLOGTIMESTAMP":  `%{MONTH} +%{MONTHDAY} %{TIME}`,
	"PROG":             `[\x21-\x5a\x5c\x5e-\x7e]+`,
	"SYSLOGPROG":       `%{PROG:program}(?:\[%{POSINT:pid}\])?`,
	"SYSLOGHOST":       `%{IPORHOST}`,
	"SYSLOGFACILITY":   `<%{NONNEGINT:facility}.%{NONNEGINT:priority}>`,
	"HTTPDATE":         `%{MONTHDAY}/%{MONTH}/%{YEAR}:%{TIME} %{INT}`,

	// Log levels
	"LOGLEVEL":    `(?:[Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)`,

	// Common identifiers
	"UUID":        `[A-Fa-f0-9]{8}-(?:[A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}`,
	"URN":         `urn:[0-9A-Za-z][0-9A-Za-z-]{0,31}:(?:%[0-9a-fA-F]{2}|[0-9A-Za-z()+,.:=@;$_!*'/?#-])+`,

	// HTTP
	"HTTPVERSION": `HTTP/%{NUMBER}`,
}

// Web/HTTP specific patterns
var webPatterns = map[string]string{
	"COMMONAPACHELOG": `%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes}|-)`,
	"COMBINEDAPACHELOG": `%{COMMONAPACHELOG} %{QS:referrer} %{QS:agent}`,
	"QS":                `%{QUOTEDSTRING}`,
}

// Java/Application patterns
var javaPatterns = map[string]string{
	"JAVACLASS":       `(?:[a-zA-Z$_][a-zA-Z$_0-9]*\.)*[a-zA-Z$_][a-zA-Z$_0-9]*`,
	"JAVAFILE":        `(?:[a-zA-Z$_0-9. -]+)`,
	"JAVASTACKTRACEPART": `at %{JAVACLASS:class}\.%{WORD:method}\(%{JAVAFILE:file}(?::%{NUMBER:line})?\)`,
	"JAVATHREAD":      `(?:[A-Z]{2}-Processor[\d]+)`,
	"CATALINA_DATESTAMP": `%{YEAR}-%{MONTHNUM}-%{MONTHDAY} %{HOUR}:%{MINUTE}:%{SECOND}`,
}

// AllPatterns contains all predefined patterns
var AllPatterns map[string]string

func init() {
	// Merge all pattern sets
	AllPatterns = make(map[string]string)

	// Add base patterns
	for k, v := range basePatterns {
		AllPatterns[k] = v
	}

	// Add web patterns
	for k, v := range webPatterns {
		AllPatterns[k] = v
	}

	// Add Java patterns
	for k, v := range javaPatterns {
		AllPatterns[k] = v
	}
}

// GetPattern returns a pattern by name
func GetPattern(name string) (string, bool) {
	pattern, ok := AllPatterns[name]
	return pattern, ok
}

// AddPattern adds or updates a custom pattern
func AddPattern(name, pattern string) {
	AllPatterns[name] = pattern
}
