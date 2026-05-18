package main

import (
	"context"
	"strings"

	model "github.com/dooray-go/dooray-sdk/openapi/model/project"
	"github.com/dooray-go/dooray-sdk/openapi/project"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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
			mcp.Description("The operation to perform. 'find_posts': list posts with filters. 'create_post': create a new post (requires subject, bodyContent)"),
			mcp.Enum("find_posts", "create_post"),
		),
		mcp.WithString("projectId",
			mcp.Required(),
			mcp.Description("project id, it can be a single id or a comma separated list of projectIds. it can be obtained from the find_projects tool"),
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
		}
		return mcp.NewToolResultText(result), nil
	})
}
