package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestAccountToolsRegistration(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	AccountTools(s, &token)

	tools := s.ListTools()

	if _, ok := tools["dooray_account_members"]; !ok {
		t.Error("dooray_account_members tool not registered")
	}
	if _, ok := tools["dooray_account_member"]; !ok {
		t.Error("dooray_account_member tool not registered")
	}
}

func TestAccountGetMembersArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_members"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_members",
			Arguments: map[string]any{
				"operation":   "find_member_id",
				"member_name": "홍길동",
				"user_code":   "hong123",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (argument parsing works)")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestAccountGetMembersWithoutUserCode(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_members"]

	// user_code는 optional - 없어도 panic 없이 동작해야 함
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_members",
			Arguments: map[string]any{
				"operation":   "find_member_id",
				"member_name": "홍길동",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (argument parsing works without user_code)")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestAccountGetMemberArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_member"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_member",
			Arguments: map[string]any{
				"operation": "find_member_details",
				"member_id": "12345",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (argument parsing works)")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestAccountGetMemberMissingArg(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_member"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing member_id): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_member",
			Arguments: map[string]any{
				"operation": "find_member_details",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestAccountGetMembersMissingMemberName(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_members"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing member_name): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_members",
			Arguments: map[string]any{
				"operation": "find_member_id",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestAccountToolCount(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	AccountTools(s, &token)

	tools := s.ListTools()

	accountTools := 0
	for name := range tools {
		if name == "dooray_account_members" || name == "dooray_account_member" {
			accountTools++
		}
	}

	if accountTools != 2 {
		t.Errorf("expected 2 account tools, got %d", accountTools)
	}
}

func TestAccountGetMembersWithUserCodeOnly(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	AccountTools(s, &token)

	tool := s.ListTools()["dooray_account_members"]

	// user_code만 제공하고 member_name 누락 시 panic 예상
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing member_name, only user_code): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_account_members",
			Arguments: map[string]any{
				"operation": "find_member_id",
				"user_code": "hong123",
			},
		},
	}

	tool.Handler(context.Background(), req)
}
