package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func OsTools(s *server.MCPServer, token *string) {
	osTool := mcp.NewTool("os",
		mcp.WithDescription("get os time date"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to get date time"),
			mcp.Enum("get_date_time"),
		),
	)

	// Add the calculator handler
	s.AddTool(osTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)

		var result string
		switch op {
		case "get_date_time":
			result = fmt.Sprintf(`{"time":"%s"}`, time.Now().Format("2006-01-02 15:04:05"))
		}
		return mcp.NewToolResultText(result), nil
	})
}
