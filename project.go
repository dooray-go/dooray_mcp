package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray/openapi/project"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ProjectTools(s *server.MCPServer, token *string) {
	doorayMessengerTool := mcp.NewTool("dooray_project",
		mcp.WithDescription("find dooray projects"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find projects)"),
			mcp.Enum("find_projects"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("project type, it can be either 'public' or 'private', default is 'public'"),
		),
		mcp.WithString("state",
			mcp.Required(),
			mcp.Description("project state, it can be either 'active' or 'archived', default is 'active'"),
		),
		mcp.WithString("scope",
			mcp.Required(),
			mcp.Description(
				`project state, it can be either 'private' or 'public', default is 'active',
				'private' - only the project member can see it,
				'public' - all users can see the project
`),
		),
	)

	// Add the calculator handler
	s.AddTool(doorayMessengerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})
}
