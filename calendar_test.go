package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCalendarToolsRegistration(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	CalendarTools(s, &token)

	tools := s.ListTools()

	expected := []string{"dooray_calendar_calendars", "dooray_calendar_events", "dooray_calendar_post_event"}
	for _, name := range expected {
		if _, ok := tools[name]; !ok {
			t.Errorf("%s tool not registered", name)
		}
	}
}

func TestCalendarGetCalendarsArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_calendars"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_calendars",
			Arguments: map[string]any{
				"operation": "find_calendars",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestCalendarGetEventsArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_events"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_events",
			Arguments: map[string]any{
				"operation": "find_events",
				"calendars": "cal-123,cal-456",
				"timeMin":   "2025-04-11T00:00:00+09:00",
				"timeMax":   "2025-04-12T00:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestCalendarGetEventsWithoutOptionalCalendars(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_events"]

	// calendars는 optional
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_events",
			Arguments: map[string]any{
				"operation": "find_events",
				"timeMin":   "2025-04-11T00:00:00+09:00",
				"timeMax":   "2025-04-12T00:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded without optional calendars")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestCalendarGetEventsInvalidTimeFormat(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_events"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_events",
			Arguments: map[string]any{
				"operation": "find_events",
				"timeMin":   "not-a-date",
				"timeMax":   "2025-04-12T00:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Error("expected error for invalid time format, got nil")
	} else {
		t.Logf("got expected error for invalid time: %v", err)
	}
}

func TestCalendarPostEventArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_post_event"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_post_event",
			Arguments: map[string]any{
				"operation":  "create_event",
				"calendarId": "cal-123",
				"subject":    "회의",
				"content":    "주간 회의",
				"startedAt":  "2025-04-11T10:00:00+09:00",
				"endedAt":    "2025-04-11T11:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestCalendarPostEventWithoutOptionalFields(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_post_event"]

	// calendarId, subject, content는 optional로 처리됨 (ok 체크)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_post_event",
			Arguments: map[string]any{
				"operation": "create_event",
				"startedAt": "2025-04-11T10:00:00+09:00",
				"endedAt":   "2025-04-11T11:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded without optional fields")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestCalendarToolCount(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	CalendarTools(s, &token)

	tools := s.ListTools()
	expected := []string{"dooray_calendar_calendars", "dooray_calendar_events", "dooray_calendar_post_event"}
	if len(tools) != len(expected) {
		t.Errorf("expected %d calendar tools, got %d", len(expected), len(tools))
	}
}

func TestCalendarGetEventsInvalidTimeMaxFormat(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_events"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_events",
			Arguments: map[string]any{
				"operation": "find_events",
				"timeMin":   "2025-04-11T00:00:00+09:00",
				"timeMax":   "invalid-date",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Error("expected error for invalid timeMax format, got nil")
	} else {
		t.Logf("got expected error for invalid timeMax: %v", err)
	}
}

func TestCalendarPostEventInvalidStartedAt(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_post_event"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_post_event",
			Arguments: map[string]any{
				"operation":  "create_event",
				"calendarId": "cal-123",
				"subject":    "테스트",
				"content":    "내용",
				"startedAt":  "bad-date",
				"endedAt":    "2025-04-11T11:00:00+09:00",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Error("expected error for invalid startedAt format, got nil")
	} else {
		t.Logf("got expected error for invalid startedAt: %v", err)
	}
}

func TestCalendarPostEventInvalidEndedAt(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_post_event"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_post_event",
			Arguments: map[string]any{
				"operation":  "create_event",
				"calendarId": "cal-123",
				"subject":    "테스트",
				"content":    "내용",
				"startedAt":  "2025-04-11T10:00:00+09:00",
				"endedAt":    "bad-date",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Error("expected error for invalid endedAt format, got nil")
	} else {
		t.Logf("got expected error for invalid endedAt: %v", err)
	}
}

func TestCalendarGetEventsBothTimesInvalid(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	CalendarTools(s, &token)

	tool := s.ListTools()["dooray_calendar_events"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_calendar_events",
			Arguments: map[string]any{
				"operation": "find_events",
				"timeMin":   "bad-date",
				"timeMax":   "bad-date",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Error("expected error for both invalid time formats, got nil")
	} else {
		t.Logf("got expected error: %v", err)
	}
}
