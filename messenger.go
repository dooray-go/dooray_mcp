package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray/openapi/messenger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func MessengerTools(s *server.MCPServer, token *string) {
	doorayMessengerTool := mcp.NewTool("dooray_messenger",
		mcp.WithDescription("send message to dooray messenger"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (send)"),
			mcp.Enum("send"),
		),
		mcp.WithString("to",
			mcp.Required(),
			mcp.Description("recipient"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("message to send"),
		),
	)

	// Add the calculator handler
	s.AddTool(doorayMessengerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		to := request.Params.Arguments["to"].(string)
		message := request.Params.Arguments["message"].(string)

		var result string
		switch op {
		case "send":
			res, err := messenger.NewDefaultMessenger().DirectSend(*token,
				&messenger.DirectSendRequest{
					Text:                 message,
					OrganizationMemberId: to,
				})

			if err != nil {
				return nil, err
			}
			result = res.RawJSON
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})
}
