package main

import (
	"context"

	"github.com/dooray-go/dooray/openapi/project"
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
		op := request.Params.Arguments["operation"].(string)
		projectType := request.Params.Arguments["type"].(string)
		scope := request.Params.Arguments["scope"].(string)
		state := request.Params.Arguments["state"].(string)

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
			mcp.Description("The operation to perform (find dooray posts in projects)"),
			mcp.Enum("find_posts"),
		),
		mcp.WithString("projectId",
			mcp.Required(),
			mcp.Description("project id, it can be a single id or a comma separated list of projectIds. it can be obtained from the find_projects tool"),
		),
		mcp.WithString("toMemberIds",
			mcp.Required(),
			mcp.Description("member ids in charge, it can be a single organizationMemberId or a comma separated list of organizationMemberIds. "),
		),
		mcp.WithString("postWorkflowClasses",
			mcp.Required(),
			mcp.Description("It can be 'backlog', 'registered', 'working', 'closed'"),
		),
	)

	// Add the calculator handler
	s.AddTool(doorayPostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		projectId := request.Params.Arguments["projectId"].(string)
		toMemberIds := request.Params.Arguments["toMemberIds"].(string)
		postWorkflowClasses := request.Params.Arguments["postWorkflowClasses"].(string)

		var result string
		switch op {
		case "find_posts":
			var err error
			res, err := project.NewDefaultProject().GetPosts(*token, projectId, toMemberIds, postWorkflowClasses)
			if err != nil {
				return nil, err
			}
			result = res.RawJSON
		}
		return mcp.NewToolResultText(result), nil
	})
}
