package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	model "github.com/dooray-go/dooray-sdk/openapi/model/project"
	"github.com/dooray-go/dooray-sdk/openapi/project"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const doorayAPIEndpoint = "https://api.dooray.com"

// getPost fetches a single post directly via the Dooray REST API
// because dooray-sdk (v0.4.1) does not provide a single-post lookup.
func getPost(ctx context.Context, token, projectId, postId string) (string, error) {
	url := fmt.Sprintf("%s/project/v1/projects/%s/posts/%s", doorayAPIEndpoint, projectId, postId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "dooray-api "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get post failed: status %d, body: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

// getURL performs an authenticated GET and returns the raw JSON body.
// Generalization of getPost for endpoints the SDK does not cover (e.g. logs).
func getURL(ctx context.Context, token, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "dooray-api "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("GET %s failed: status %d, body: %s", url, resp.StatusCode, string(body))
	}
	return string(body), nil
}

// putJSON performs an authenticated PUT with a JSON payload and returns the raw JSON body.
func putJSON(ctx context.Context, token, url string, payload []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "dooray-api "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("PUT %s failed: status %d, body: %s", url, resp.StatusCode, string(body))
	}
	return string(body), nil
}

// postJSON performs an authenticated POST with a JSON payload and returns the raw JSON body.
func postJSON(ctx context.Context, token, url string, payload []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "dooray-api "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("POST %s failed: status %d, body: %s", url, resp.StatusCode, string(body))
	}
	return string(body), nil
}

// recipientWorkflow is the inline workflow field observed on closed posts:
// {"type":"member","member":{...},"workflow":{"id":"..."}}.
type recipientWorkflow struct {
	ID string `json:"id"`
}

// updateRecipient mirrors model.PostRecipient but adds the inline workflow field.
// The SDK PostRecipient has no workflow field, so update_post uses this local type
// to be able to reassign a post and set the assignee's workflow in one PUT.
type updateRecipient struct {
	Type     string             `json:"type"`
	Member   *model.PostMember  `json:"member,omitempty"`
	Workflow *recipientWorkflow `json:"workflow,omitempty"`
}

type updatePostUsers struct {
	To []updateRecipient `json:"to,omitempty"`
	Cc []updateRecipient `json:"cc,omitempty"`
}

// updatePostRequest is the PUT .../posts/{postId} payload. PUT is full-replacement,
// so subject and body must always be resent or they get cleared.
type updatePostRequest struct {
	Subject  string           `json:"subject"`
	Body     model.PostBody   `json:"body"`
	Users    *updatePostUsers `json:"users,omitempty"`
	TagIDs   []string         `json:"tagIds,omitempty"`
	Priority string           `json:"priority,omitempty"`
}

func ProjectTools(s *server.MCPServer, token *string) {
	projectTools(s, token)
	postTools(s, token)
}

func projectTools(s *server.MCPServer, token *string) {
	doorayPostTool := mcp.NewTool("dooray_project",
		mcp.WithDescription("find dooray projects"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find projects)"),
			mcp.Enum("find_projects"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("project type, it can be either 'public' or 'private', default is 'public', it can not be 'all' to get all projects. "),
		),
		mcp.WithString("state",
			mcp.Required(),
			mcp.Description("project state, it can be either 'active' or 'archived', default is 'active'"),
		),
		mcp.WithString("scope",
			mcp.Required(),
			mcp.Description(
				`project state, it can be either 'private' or 'public', default is 'private',
				'private' - only the project member can see it,
				'public' - all users can see the project,
				it can not be 'all' to get all projects
`),
		),
	)

	// Add the calculator handler
	s.AddTool(doorayPostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.GetArguments()["operation"].(string)
		projectType := request.GetArguments()["type"].(string)
		scope := request.GetArguments()["scope"].(string)
		state := request.GetArguments()["state"].(string)

		var result string
		switch op {
		case "find_projects":
			var err error
			res, err := project.NewDefaultProject().GetProjects(*token, projectType, scope, state)
			if err != nil {
				return nil, err
			}
			result = res.RawJSON
		}
		return mcp.NewToolResultText(result), nil
	})
}

func postTools(s *server.MCPServer, token *string) {
	doorayPostTool := mcp.NewTool("dooray_posts",
		mcp.WithDescription("find dooray posts in projects"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform. 'find_posts': list posts with filters. 'get_post': get a single post with full body (requires postId). 'create_post': create a new post (requires subject, bodyContent). 'update_post': FULL-REPLACEMENT edit of a post (requires postId, subject, bodyContent — any field not resent is cleared; fetch current subject/body via get_post first). 'set_workflow': change a post's status/workflow (requires postId, setWorkflowId). 'create_log': add a comment/log to a post (requires postId, logContent). 'get_logs': list a post's comments/logs (requires postId)."),
			mcp.Enum("find_posts", "get_post", "create_post", "update_post", "set_workflow", "create_log", "get_logs"),
		),
		mcp.WithString("projectId",
			mcp.Required(),
			mcp.Description("project id, it can be a single id or a comma separated list of projectIds. it can be obtained from the find_projects tool. for get_post it must be a single id"),
		),
		// get_post fields
		mcp.WithString("postId",
			mcp.Description("post id (required for get_post). it can be obtained from the find_posts tool"),
		),
		// create_post fields
		mcp.WithString("subject",
			mcp.Description("post subject/title (required for create_post)"),
		),
		mcp.WithString("bodyContent",
			mcp.Description("post body markdown/html content (required for create_post)"),
		),
		mcp.WithString("bodyMimeType",
			mcp.Description("post body mime type for create_post: 'text/x-markdown' (default) or 'text/html'"),
		),
		mcp.WithString("toMemberIdsCreate",
			mcp.Description("assignee organizationMemberIds for create_post, comma separated. type=member"),
		),
		mcp.WithString("ccMemberIdsCreate",
			mcp.Description("cc organizationMemberIds for create_post, comma separated. type=member"),
		),
		mcp.WithString("tagIdsCreate",
			mcp.Description("tag ids for create_post, comma separated"),
		),
		mcp.WithString("priority",
			mcp.Description("priority for create_post: urgent | high | normal | low"),
		),
		mcp.WithString("parentPostIdCreate",
			mcp.Description("parent post id for create_post (sub-task)"),
		),
		mcp.WithString("milestoneIdCreate",
			mcp.Description("milestone id for create_post"),
		),
		mcp.WithString("workflowId",
			mcp.Description("workflow id for create_post"),
		),
		// set_workflow
		mcp.WithString("setWorkflowId",
			mcp.Description("target workflow id for set_workflow (the post's new status). Required for set_workflow."),
		),
		// create_log / get_logs
		mcp.WithString("logContent",
			mcp.Description("comment/log body content for create_log. Required for create_log."),
		),
		mcp.WithString("logMimeType",
			mcp.Description("mime type for create_log body: 'text/x-markdown' (default) or 'text/html'"),
		),
		// update_post
		mcp.WithString("toMemberWorkflowId",
			mcp.Description("update_post only: attaches an inline workflow.id to each assignee in toMemberIdsCreate. NOTE: Dooray ignores this inline workflow on the update PUT (verified) — use the separate set_workflow operation to change status. Kept for completeness; normally leave empty."),
		),
		// Paging
		mcp.WithNumber("page",
			mcp.Description("page number, default is 0"),
		),
		mcp.WithNumber("size",
			mcp.Description("number of posts per page, default is 20, max is 100"),
		),
		// Filters
		mcp.WithString("fromEmailAddress",
			mcp.Description("filter posts by sender email address"),
		),
		mcp.WithString("fromMemberIds",
			mcp.Description("filter posts created by specific members, comma separated organizationMemberIds"),
		),
		mcp.WithString("toMemberIds",
			mcp.Description("filter posts assigned to specific members, comma separated organizationMemberIds"),
		),
		mcp.WithNumber("toMemberSize",
			mcp.Description("filter by number of assignees (0: no assignees, 1: single assignee matching toMemberIds[0])"),
		),
		mcp.WithString("ccMemberIds",
			mcp.Description("filter posts where specific members are CC'd, comma separated organizationMemberIds"),
		),
		mcp.WithString("tagIds",
			mcp.Description("filter posts by tag ids, comma separated"),
		),
		mcp.WithString("parentPostId",
			mcp.Description("filter sub-tasks of a specific parent post"),
		),
		mcp.WithString("postNumber",
			mcp.Description("filter by post number"),
		),
		mcp.WithString("postWorkflowClasses",
			mcp.Description("filter by workflow class: 'backlog', 'registered', 'working', 'closed' (comma separated)"),
		),
		mcp.WithString("postWorkflowIds",
			mcp.Description("filter by workflow ids, comma separated"),
		),
		mcp.WithString("milestoneIds",
			mcp.Description("filter by milestone ids, comma separated"),
		),
		mcp.WithString("subjects",
			mcp.Description("filter by post subject keyword"),
		),
		// Date filters
		mcp.WithString("createdAt",
			mcp.Description("filter by creation date. supported patterns: 'today', 'thisweek', 'prev-{N}d' (e.g. prev-30d), 'next-{N}d' (e.g. next-7d), or ISO8601 range with ~ separator (e.g. 2025-01-01T00:00:00+09:00~2025-12-31T23:59:59+09:00)"),
		),
		mcp.WithString("updatedAt",
			mcp.Description("filter by update date. supported patterns: 'today', 'thisweek', 'prev-{N}d' (e.g. prev-30d), 'next-{N}d' (e.g. next-7d), or ISO8601 range with ~ separator (e.g. 2025-01-01T00:00:00+09:00~2025-12-31T23:59:59+09:00)"),
		),
		mcp.WithString("dueAt",
			mcp.Description("filter by due date. supported patterns: 'today', 'thisweek', 'prev-{N}d' (e.g. prev-30d), 'next-{N}d' (e.g. next-7d), or ISO8601 range with ~ separator (e.g. 2025-01-01T00:00:00+09:00~2025-12-31T23:59:59+09:00)"),
		),
		// Sort
		mcp.WithString("order",
			mcp.Description("sort order: postDueAt, postUpdatedAt, createdAt (prefix with - for descending, e.g. -createdAt)"),
		),
	)

	s.AddTool(doorayPostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.GetArguments()["operation"].(string)
		projectId := request.GetArguments()["projectId"].(string)

		opts := project.GetPostsOptions{}

		// Paging
		if v, ok := request.GetArguments()["page"]; ok {
			page := int(v.(float64))
			opts.Page = &page
		}
		if v, ok := request.GetArguments()["size"]; ok {
			size := int(v.(float64))
			opts.Size = &size
		}

		// Filters
		if v, ok := request.GetArguments()["fromEmailAddress"].(string); ok {
			opts.FromEmailAddress = v
		}
		if v, ok := request.GetArguments()["fromMemberIds"].(string); ok {
			opts.FromMemberIds = v
		}
		if v, ok := request.GetArguments()["toMemberIds"].(string); ok {
			opts.ToMemberIds = v
		}
		if v, ok := request.GetArguments()["toMemberSize"]; ok {
			size := int(v.(float64))
			opts.ToMemberSize = &size
		}
		if v, ok := request.GetArguments()["ccMemberIds"].(string); ok {
			opts.CcMemberIds = v
		}
		if v, ok := request.GetArguments()["tagIds"].(string); ok {
			opts.TagIds = v
		}
		if v, ok := request.GetArguments()["parentPostId"].(string); ok {
			opts.ParentPostId = v
		}
		if v, ok := request.GetArguments()["postNumber"].(string); ok {
			opts.PostNumber = v
		}
		if v, ok := request.GetArguments()["postWorkflowClasses"].(string); ok {
			opts.PostWorkflowClasses = v
		}
		if v, ok := request.GetArguments()["postWorkflowIds"].(string); ok {
			opts.PostWorkflowIds = v
		}
		if v, ok := request.GetArguments()["milestoneIds"].(string); ok {
			opts.MilestoneIds = v
		}
		if v, ok := request.GetArguments()["subjects"].(string); ok {
			opts.Subjects = v
		}

		// Date filters
		if v, ok := request.GetArguments()["createdAt"].(string); ok {
			opts.CreatedAt = v
		}
		if v, ok := request.GetArguments()["updatedAt"].(string); ok {
			opts.UpdatedAt = v
		}
		if v, ok := request.GetArguments()["dueAt"].(string); ok {
			opts.DueAt = v
		}

		// Sort
		if v, ok := request.GetArguments()["order"].(string); ok {
			opts.Order = v
		}

		var result string
		switch op {
		case "find_posts":
			res, err := project.NewDefaultProject().GetPostsWithOptions(*token, projectId, opts)
			if err != nil {
				return nil, err
			}
			result = res.RawJSON
		case "get_post":
			postId, _ := request.GetArguments()["postId"].(string)
			if postId == "" {
				return mcp.NewToolResultError("postId is required for get_post"), nil
			}
			res, err := getPost(ctx, *token, projectId, postId)
			if err != nil {
				return nil, err
			}
			result = res
		case "create_post":
			subject, _ := request.GetArguments()["subject"].(string)
			bodyContent, _ := request.GetArguments()["bodyContent"].(string)
			if subject == "" || bodyContent == "" {
				return mcp.NewToolResultError("subject and bodyContent are required for create_post"), nil
			}
			bodyMimeType, _ := request.GetArguments()["bodyMimeType"].(string)
			if bodyMimeType == "" {
				bodyMimeType = "text/x-markdown"
			}

			post := model.PostRequest{
				Subject: subject,
				Body: model.PostBody{
					MimeType: bodyMimeType,
					Content:  bodyContent,
				},
			}

			if v, _ := request.GetArguments()["priority"].(string); v != "" {
				post.Priority = v
			}
			if v, _ := request.GetArguments()["parentPostIdCreate"].(string); v != "" {
				post.ParentPostID = v
			}
			if v, _ := request.GetArguments()["milestoneIdCreate"].(string); v != "" {
				post.MilestoneID = v
			}
			if v, _ := request.GetArguments()["workflowId"].(string); v != "" {
				post.WorkflowID = v
			}
			if v, _ := request.GetArguments()["tagIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					id = strings.TrimSpace(id)
					if id != "" {
						post.TagIDs = append(post.TagIDs, id)
					}
				}
			}

			users := &model.PostUsers{}
			hasUsers := false
			if v, _ := request.GetArguments()["toMemberIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					id = strings.TrimSpace(id)
					if id != "" {
						users.To = append(users.To, model.PostRecipient{
							Type:   "member",
							Member: &model.PostMember{OrganizationMemberID: id},
						})
						hasUsers = true
					}
				}
			}
			if v, _ := request.GetArguments()["ccMemberIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					id = strings.TrimSpace(id)
					if id != "" {
						users.Cc = append(users.Cc, model.PostRecipient{
							Type:   "member",
							Member: &model.PostMember{OrganizationMemberID: id},
						})
						hasUsers = true
					}
				}
			}
			if hasUsers {
				post.Users = users
			}

			res, err := project.NewDefaultProject().CreatePost(*token, projectId, post)
			if err != nil {
				return nil, err
			}
			result = res.RawJSON
		case "set_workflow":
			postId, _ := request.GetArguments()["postId"].(string)
			setWorkflowId, _ := request.GetArguments()["setWorkflowId"].(string)
			if postId == "" || setWorkflowId == "" {
				return mcp.NewToolResultError("postId and setWorkflowId are required for set_workflow"), nil
			}
			payload, err := json.Marshal(map[string]string{"workflowId": setWorkflowId})
			if err != nil {
				return nil, err
			}
			// Dooray's set-workflow endpoint is POST (not PUT) — PUT returns 404.
			url := fmt.Sprintf("%s/project/v1/projects/%s/posts/%s/set-workflow", doorayAPIEndpoint, projectId, postId)
			res, err := postJSON(ctx, *token, url, payload)
			if err != nil {
				return nil, err
			}
			result = res
		case "create_log":
			postId, _ := request.GetArguments()["postId"].(string)
			logContent, _ := request.GetArguments()["logContent"].(string)
			if postId == "" || logContent == "" {
				return mcp.NewToolResultError("postId and logContent are required for create_log"), nil
			}
			logMimeType, _ := request.GetArguments()["logMimeType"].(string)
			if logMimeType == "" {
				logMimeType = "text/x-markdown"
			}
			payload, err := json.Marshal(map[string]any{
				"body": model.PostBody{MimeType: logMimeType, Content: logContent},
			})
			if err != nil {
				return nil, err
			}
			url := fmt.Sprintf("%s/project/v1/projects/%s/posts/%s/logs", doorayAPIEndpoint, projectId, postId)
			res, err := postJSON(ctx, *token, url, payload)
			if err != nil {
				return nil, err
			}
			result = res
		case "get_logs":
			postId, _ := request.GetArguments()["postId"].(string)
			if postId == "" {
				return mcp.NewToolResultError("postId is required for get_logs"), nil
			}
			page := 0
			size := 20
			if v, ok := request.GetArguments()["page"]; ok {
				page = int(v.(float64))
			}
			if v, ok := request.GetArguments()["size"]; ok {
				size = int(v.(float64))
			}
			url := fmt.Sprintf("%s/project/v1/projects/%s/posts/%s/logs?page=%d&size=%d", doorayAPIEndpoint, projectId, postId, page, size)
			res, err := getURL(ctx, *token, url)
			if err != nil {
				return nil, err
			}
			result = res
		case "update_post":
			postId, _ := request.GetArguments()["postId"].(string)
			subject, _ := request.GetArguments()["subject"].(string)
			bodyContent, _ := request.GetArguments()["bodyContent"].(string)
			if postId == "" || subject == "" || bodyContent == "" {
				return mcp.NewToolResultError("postId, subject and bodyContent are required for update_post (PUT is full-replacement; resend subject+body fetched via get_post or they are cleared)"), nil
			}
			bodyMimeType, _ := request.GetArguments()["bodyMimeType"].(string)
			if bodyMimeType == "" {
				bodyMimeType = "text/x-markdown"
			}

			upd := updatePostRequest{
				Subject: subject,
				Body:    model.PostBody{MimeType: bodyMimeType, Content: bodyContent},
			}
			if v, _ := request.GetArguments()["priority"].(string); v != "" {
				upd.Priority = v
			}
			if v, _ := request.GetArguments()["tagIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					if id = strings.TrimSpace(id); id != "" {
						upd.TagIDs = append(upd.TagIDs, id)
					}
				}
			}

			toWorkflowID, _ := request.GetArguments()["toMemberWorkflowId"].(string)
			users := &updatePostUsers{}
			hasUsers := false
			if v, _ := request.GetArguments()["toMemberIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					if id = strings.TrimSpace(id); id != "" {
						r := updateRecipient{Type: "member", Member: &model.PostMember{OrganizationMemberID: id}}
						if toWorkflowID != "" {
							r.Workflow = &recipientWorkflow{ID: toWorkflowID}
						}
						users.To = append(users.To, r)
						hasUsers = true
					}
				}
			}
			if v, _ := request.GetArguments()["ccMemberIdsCreate"].(string); v != "" {
				for _, id := range strings.Split(v, ",") {
					if id = strings.TrimSpace(id); id != "" {
						users.Cc = append(users.Cc, updateRecipient{Type: "member", Member: &model.PostMember{OrganizationMemberID: id}})
						hasUsers = true
					}
				}
			}
			if hasUsers {
				upd.Users = users
			}

			payload, err := json.Marshal(upd)
			if err != nil {
				return nil, err
			}
			url := fmt.Sprintf("%s/project/v1/projects/%s/posts/%s", doorayAPIEndpoint, projectId, postId)
			res, err := putJSON(ctx, *token, url, payload)
			if err != nil {
				return nil, err
			}
			result = res
		}
		return mcp.NewToolResultText(result), nil
	})
}
