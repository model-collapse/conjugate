// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGrokOperator_SimplePattern(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "192.168.1.1 - admin"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{IP:client_ip} - %{USER:username}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should have extracted fields
	clientIP, exists := row.Get("client_ip")
	assert.True(t, exists)
	assert.Equal(t, "192.168.1.1", clientIP)

	username, exists := row.Get("username")
	assert.True(t, exists)
	assert.Equal(t, "admin", username)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_CommonApacheLog(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{COMMONAPACHELOG}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Verify extracted fields
	clientIP, _ := row.Get("clientip")
	assert.Equal(t, "127.0.0.1", clientIP)

	auth, _ := row.Get("auth")
	assert.Equal(t, "frank", auth)

	verb, _ := row.Get("verb")
	assert.Equal(t, "GET", verb)

	request, _ := row.Get("request")
	assert.Equal(t, "/apache_pb.gif", request)

	response, _ := row.Get("response")
	assert.Equal(t, "200", response)

	bytes, _ := row.Get("bytes")
	assert.Equal(t, "2326", bytes)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_TypeConversion(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "User 12345 transferred 9876.54 bytes"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "User %{INT:user_id:int} transferred %{NUMBER:bytes:float} bytes",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Verify type conversions
	userID, exists := row.Get("user_id")
	assert.True(t, exists)
	assert.IsType(t, int64(0), userID)
	assert.Equal(t, int64(12345), userID)

	bytes, exists := row.Get("bytes")
	assert.True(t, exists)
	assert.IsType(t, float64(0), bytes)
	assert.Equal(t, float64(9876.54), bytes)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_NoMatch(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "This doesn't match the pattern"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
			"id":   1,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{IP:client_ip}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Original row should be returned unchanged
	_, exists := row.Get("client_ip")
	assert.False(t, exists)

	id, exists := row.Get("id")
	assert.True(t, exists)
	assert.Equal(t, 1, id)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_MissingInputField(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"id": 1,
			// No _raw field
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{IP:client_ip}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Row should be returned as-is
	_, exists := row.Get("client_ip")
	assert.False(t, exists)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_CustomInputField(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "ERROR: Connection failed"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"message": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern:    "%{LOGLEVEL:level}: %{GREEDYDATA:error_msg}",
		InputField: "message",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	level, exists := row.Get("level")
	assert.True(t, exists)
	assert.Equal(t, "ERROR", level)

	errorMsg, exists := row.Get("error_msg")
	assert.True(t, exists)
	assert.Equal(t, "Connection failed", errorMsg)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_KeepOriginal(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "192.168.1.1 GET /api"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern:      "%{IP:ip} %{WORD:method} %{URIPATH:path}",
		KeepOriginal: true,
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Extracted fields should exist
	ip, _ := row.Get("ip")
	assert.Equal(t, "192.168.1.1", ip)

	// Original field should still exist
	raw, exists := row.Get("_raw")
	assert.True(t, exists)
	assert.Equal(t, logLine, raw)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_RemoveOriginal(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "192.168.1.1 GET /api"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern:      "%{IP:ip} %{WORD:method} %{URIPATH:path}",
		KeepOriginal: false, // Default
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Extracted fields should exist
	ip, _ := row.Get("ip")
	assert.Equal(t, "192.168.1.1", ip)

	// Original field should be removed
	_, exists := row.Get("_raw")
	assert.False(t, exists)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_MultipleRows(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": "INFO: Process started",
		}),
		NewRow(map[string]interface{}{
			"_raw": "WARN: Low memory",
		}),
		NewRow(map[string]interface{}{
			"_raw": "ERROR: Connection failed",
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{LOGLEVEL:level}: %{GREEDYDATA:message}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	// First row
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	level1, _ := row1.Get("level")
	assert.Equal(t, "INFO", level1)

	// Second row
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	level2, _ := row2.Get("level")
	assert.Equal(t, "WARN", level2)

	// Third row
	row3, err := op.Next(ctx)
	require.NoError(t, err)
	level3, _ := row3.Get("level")
	assert.Equal(t, "ERROR", level3)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_ComplexNginxLog(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Nginx access log format
	logLine := `10.0.0.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users?id=123 HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: `%{IP:remote_addr} - - \[%{HTTPDATE:time_local}\] "%{WORD:method} %{URIPATHPARAM:uri} HTTP/%{NUMBER:http_version}" %{INT:status:int} %{INT:body_bytes_sent:int}`,
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Verify fields
	remoteAddr, _ := row.Get("remote_addr")
	assert.Equal(t, "10.0.0.1", remoteAddr)

	method, _ := row.Get("method")
	assert.Equal(t, "GET", method)

	uri, _ := row.Get("uri")
	assert.Equal(t, "/api/users?id=123", uri)

	status, _ := row.Get("status")
	assert.Equal(t, int64(200), status)
	assert.IsType(t, int64(0), status)

	bodyBytes, _ := row.Get("body_bytes_sent")
	assert.Equal(t, int64(1234), bodyBytes)
	assert.IsType(t, int64(0), bodyBytes)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_SyslogFormat(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := `Jan 30 14:30:00 hostname sshd[12345]: Failed password for invalid user admin from 192.168.1.100 port 22 ssh2`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: `%{SYSLOGTIMESTAMP:timestamp} %{SYSLOGHOST:hostname} %{DATA:program}\[%{POSINT:pid:int}\]: %{GREEDYDATA:message}`,
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	timestamp, _ := row.Get("timestamp")
	assert.Equal(t, "Jan 30 14:30:00", timestamp)

	hostname, _ := row.Get("hostname")
	assert.Equal(t, "hostname", hostname)

	program, _ := row.Get("program")
	assert.Equal(t, "sshd", program)

	pid, _ := row.Get("pid")
	assert.Equal(t, int64(12345), pid)
	assert.IsType(t, int64(0), pid)

	message, _ := row.Get("message")
	assert.Contains(t, message, "Failed password")

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_CustomPattern(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "Transaction TXN12345 completed in 1234ms"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	// Define custom pattern
	customPatterns := map[string]string{
		"TXNID": `TXN[0-9]+`,
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern:        "Transaction %{TXNID:txn_id} completed in %{INT:duration:int}ms",
		CustomPatterns: customPatterns,
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	txnID, exists := row.Get("txn_id")
	assert.True(t, exists)
	assert.Equal(t, "TXN12345", txnID)

	duration, exists := row.Get("duration")
	assert.True(t, exists)
	assert.Equal(t, int64(1234), duration)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_EmailPattern(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "User alice@example.com logged in from 10.0.0.1"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "User %{EMAILADDRESS:email} logged in from %{IP:ip}",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	email, _ := row.Get("email")
	assert.Equal(t, "alice@example.com", email)

	ip, _ := row.Get("ip")
	assert.Equal(t, "10.0.0.1", ip)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_UUIDPattern(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	logLine := "Request 550e8400-e29b-41d4-a716-446655440000 processed"

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": logLine,
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "Request %{UUID:request_id} processed",
	}

	op, err := NewGrokOperator(input, config, logger)
	require.NoError(t, err)

	err = op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	requestID, exists := row.Get("request_id")
	assert.True(t, exists)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", requestID)

	err = op.Close()
	require.NoError(t, err)
}

func TestGrokOperator_InvalidPattern(t *testing.T) {
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": "test",
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "%{NONEXISTENTPATTERN:field}",
	}

	_, err := NewGrokOperator(input, config, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown pattern")
}

func TestGrokOperator_EmptyPattern(t *testing.T) {
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": "test",
		}),
	}

	input := NewSliceIterator(rows)
	config := GrokConfig{
		Pattern: "",
	}

	_, err := NewGrokOperator(input, config, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pattern is required")
}
