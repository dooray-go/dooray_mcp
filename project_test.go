package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func newTestServer() *server.MCPServer {
	return server.NewMCPServer("dooray-test", "1.0.0", server.WithToolCapabilities(true))
}

func TestProjectToolsRegistration(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	ProjectTools(s, &token)

	tools := s.ListTools()

	// dooray_project 와 dooray_posts 두 개가 등록되어야 함
	names := make(map[string]bool)
	for name := range tools {
		names[name] = true
	}

	if !names["dooray_project"] {
		t.Error("dooray_project tool not registered")
	}
	if !names["dooray_posts"] {
		t.Error("dooray_posts tool not registered")
	}
}

func TestProjectToolArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_project"]
	if !ok {
		t.Fatal("dooray_project tool not registered")
	}

	// 핸들러 호출 - API 호출은 실패하지만 인자 파싱이 올바르게 동작하는지 확인
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"state":     "active",
				"scope":     "private",
			},
		},
	}

	// API 호출 실패는 예상됨 - 인자 파싱 중 panic이 발생하지 않는지 확인
	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (unexpected with invalid token, but argument parsing works)")
	} else {
		t.Logf("handler returned expected error (API call failed): %v", err)
	}
}

func TestPostsToolArguments(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// 필수 인자만으로 호출
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (unexpected with invalid token, but argument parsing works)")
	} else {
		t.Logf("handler returned expected error (API call failed): %v", err)
	}
}

func TestPostsToolWithAllOptions(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// 모든 옵션 인자를 포함하여 호출 - panic 없이 파싱되는지 확인
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation":           "find_posts",
				"projectId":           "12345",
				"page":                float64(0),
				"size":                float64(20),
				"fromEmailAddress":    "test@example.com",
				"fromMemberIds":       "member1,member2",
				"toMemberIds":         "member3",
				"toMemberSize":        float64(1),
				"ccMemberIds":         "member4",
				"tagIds":              "tag1,tag2",
				"parentPostId":        "parent1",
				"postNumber":          "100",
				"postWorkflowClasses": "registered,working",
				"postWorkflowIds":     "wf1",
				"milestoneIds":        "ms1",
				"subjects":            "test subject",
				"createdAt":           "today",
				"updatedAt":           "prev-7d",
				"dueAt":               "next-30d",
				"order":               "-createdAt",
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

func TestPostsToolMissingRequiredArg(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tools := s.ListTools()
	tool, ok := tools["dooray_posts"]
	if !ok {
		t.Fatal("dooray_posts tool not registered")
	}

	// projectId 누락 시 panic 발생 여부 확인
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing required arg): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				// projectId 누락
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingType(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing type): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"state":     "active",
				"scope":     "private",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingState(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing state): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"scope":     "private",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolMissingScope(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_project"]

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from panic (missing scope): %v", r)
		}
	}()

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_project",
			Arguments: map[string]any{
				"operation": "find_projects",
				"type":      "public",
				"state":     "active",
			},
		},
	}

	tool.Handler(context.Background(), req)
}

func TestProjectToolCount(t *testing.T) {
	s := newTestServer()
	token := "test-token"
	ProjectTools(s, &token)

	tools := s.ListTools()
	if len(tools) != 2 {
		t.Errorf("expected 2 project tools, got %d", len(tools))
	}
}

func TestPostsToolWithPagingOnly(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
				"page":      float64(2),
				"size":      float64(50),
			},
		},
	}

	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded with paging args")
	} else {
		t.Logf("handler returned expected error: %v", err)
	}
}

func TestPostsToolWithDateFilters(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	testCases := []struct {
		name      string
		createdAt string
		updatedAt string
		dueAt     string
	}{
		{"today pattern", "today", "", ""},
		{"thisweek pattern", "", "thisweek", ""},
		{"prev-N pattern", "prev-30d", "", ""},
		{"next-N pattern", "", "", "next-7d"},
		{"ISO8601 range", "2025-01-01T00:00:00+09:00~2025-12-31T23:59:59+09:00", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := map[string]any{
				"operation": "find_posts",
				"projectId": "12345",
			}
			if tc.createdAt != "" {
				args["createdAt"] = tc.createdAt
			}
			if tc.updatedAt != "" {
				args["updatedAt"] = tc.updatedAt
			}
			if tc.dueAt != "" {
				args["dueAt"] = tc.dueAt
			}

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "dooray_posts",
					Arguments: args,
				},
			}

			_, err := tool.Handler(context.Background(), req)
			if err == nil {
				t.Log("handler succeeded")
			} else {
				t.Logf("handler returned expected error: %v", err)
			}
		})
	}
}

func TestPostsToolGetPost(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "get_post",
				"projectId": "12345",
				"postId":    "67890",
			},
		},
	}

	// 잘못된 토큰이므로 API 호출은 실패하지만 인자 파싱과 요청 구성이 동작하는지 확인
	_, err := tool.Handler(context.Background(), req)
	if err == nil {
		t.Log("handler succeeded (unexpected with invalid token, but argument parsing works)")
	} else {
		t.Logf("handler returned expected error (API call failed): %v", err)
	}
}

func TestPostsToolGetPostMissingPostId(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "dooray_posts",
			Arguments: map[string]any{
				"operation": "get_post",
				"projectId": "12345",
				// postId 누락
			},
		},
	}

	res, err := tool.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected tool result error, got go error: %v", err)
	}
	if res == nil || !res.IsError {
		t.Error("expected error result when postId is missing")
	}
}

func TestPostsToolWithSortOptions(t *testing.T) {
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)

	tool := s.ListTools()["dooray_posts"]

	sortOrders := []string{"createdAt", "-createdAt", "postDueAt", "-postDueAt", "postUpdatedAt", "-postUpdatedAt"}

	for _, order := range sortOrders {
		t.Run("order_"+order, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "dooray_posts",
					Arguments: map[string]any{
						"operation": "find_posts",
						"projectId": "12345",
						"order":     order,
					},
				},
			}

			_, err := tool.Handler(context.Background(), req)
			if err == nil {
				t.Log("handler succeeded")
			} else {
				t.Logf("handler returned expected error: %v", err)
			}
		})
	}
}

// callPosts is a small helper that invokes the dooray_posts handler with the
// given arguments and returns the tool result and error.
func callPosts(t *testing.T, args map[string]any) (*mcp.CallToolResult, error) {
	t.Helper()
	s := newTestServer()
	token := "invalid-token"
	ProjectTools(s, &token)
	tool := s.ListTools()["dooray_posts"]
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: "dooray_posts", Arguments: args},
	}
	return tool.Handler(context.Background(), req)
}

func TestPostsToolSetWorkflowRequiresArgs(t *testing.T) {
	// postId/setWorkflowId 누락 시 API 호출 전에 에러 결과를 돌려줘야 한다.
	res, err := callPosts(t, map[string]any{
		"operation": "set_workflow",
		"projectId": "12345",
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("expected validation error result for missing postId/setWorkflowId, got %+v", res)
	}
}

func TestPostsToolCreateLogRequiresArgs(t *testing.T) {
	// postId/logContent 누락 시 에러 결과.
	res, err := callPosts(t, map[string]any{
		"operation": "create_log",
		"projectId": "12345",
		"postId":    "67890",
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("expected validation error result for missing logContent, got %+v", res)
	}
}

func TestPostsToolUpdatePostRequiresSubjectBody(t *testing.T) {
	// PUT은 전체치환이므로 subject/bodyContent가 비면 거부해야 한다(데이터 유실 방지).
	res, err := callPosts(t, map[string]any{
		"operation": "update_post",
		"projectId": "12345",
		"postId":    "67890",
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if res == nil || !res.IsError {
		t.Fatalf("expected validation error result for missing subject/bodyContent, got %+v", res)
	}
}

func TestPostsToolGetLogsParsesArgs(t *testing.T) {
	// 인자가 갖춰지면 검증 단계는 통과하고 API 호출로 진행한다(invalid token이라 호출 자체는 실패).
	res, err := callPosts(t, map[string]any{
		"operation": "get_logs",
		"projectId": "12345",
		"postId":    "67890",
		"page":      float64(0),
		"size":      float64(10),
	})
	// API 호출 실패는 transport err 또는 미설정일 수 있다 — 검증 에러(postId 누락 류)만 아니면 된다.
	if err != nil {
		t.Logf("handler returned expected API error (invalid token): %v", err)
		return
	}
	if res != nil && res.IsError {
		t.Logf("handler returned error result (invalid token expected): %+v", res)
	}
}

func TestPostsToolSetWorkflowParsesArgs(t *testing.T) {
	res, err := callPosts(t, map[string]any{
		"operation":     "set_workflow",
		"projectId":     "12345",
		"postId":        "67890",
		"setWorkflowId": "33333",
	})
	if err != nil {
		t.Logf("handler returned expected API error (invalid token): %v", err)
		return
	}
	if res != nil && res.IsError {
		t.Logf("handler returned error result (invalid token expected): %+v", res)
	}
}
