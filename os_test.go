package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestOsToolsRegistration(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tools := s.ListTools()
	if _, ok := tools["os"]; !ok {
		t.Fatal("os tool not registered")
	}
}

func TestOsGetDateTime(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tool := s.ListTools()["os"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "os",
			Arguments: map[string]any{
				"operation": "get_date_time",
			},
		},
	}

	before := time.Now()
	result, err := tool.Handler(context.Background(), req)
	after := time.Now()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}

	// 결과에서 시간 문자열 추출
	text := result.Content[0].(mcp.TextContent).Text
	var parsed struct {
		Time string `json:"time"`
	}
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("failed to parse result JSON: %v", err)
	}

	resultTime, err := time.ParseInLocation("2006-01-02 15:04:05", parsed.Time, time.Now().Location())
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	// 반환된 시간이 호출 전후 사이에 있는지 확인
	if resultTime.Before(before.Truncate(time.Second)) || resultTime.After(after.Add(time.Second)) {
		t.Errorf("returned time %v is not between %v and %v", resultTime, before, after)
	}
}

func TestOsMissingOperation(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tool := s.ListTools()["os"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing operation): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "os",
			Arguments: map[string]any{},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestOsGetDateTimeResultFormat(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tool := s.ListTools()["os"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "os",
			Arguments: map[string]any{
				"operation": "get_date_time",
			},
		},
	}

	result, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text

	// JSON 형식 확인
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	// "time" 키가 존재하는지 확인
	timeVal, ok := parsed["time"]
	if !ok {
		t.Fatal("result JSON does not contain 'time' key")
	}

	// time 값이 string인지 확인
	timeStr, ok := timeVal.(string)
	if !ok {
		t.Fatal("'time' value is not a string")
	}

	// 포맷이 올바른지 확인
	_, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Now().Location())
	if err != nil {
		t.Errorf("time format is invalid: %v", err)
	}
}

func TestOsToolCount(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tools := s.ListTools()
	if _, ok := tools["os"]; !ok {
		t.Fatal("os tool not registered")
	}
}

func TestOsGetDateTimeNilArguments(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	OsTools(s, &token)

	tool := s.ListTools()["os"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (nil arguments): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "os",
			Arguments: nil,
		},
	}

	tool.Handler(context.Background(), req)
}
