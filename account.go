package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray/openapi/account"
	model "github.com/dooray-go/dooray/openapi/model/account"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AccountTools(s *server.MCPServer, token *string) {
	accountGetMembers := mcp.NewTool("dooray_account_members",
		mcp.WithDescription("find dooray account members by name or userCode"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find member id)"),
			mcp.Enum("find_member_id"),
		),
		mcp.WithString("member_name",
			mcp.Required(),
			mcp.Description("member name"),
		),
		mcp.WithString("user_code",
			mcp.Required(),
			mcp.Description("user code, user code is composed of english letters and numbers, and the length is 4-20 characters"),
		),
	)

	s.AddTool(accountGetMembers, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		memberName := request.Params.Arguments["member_name"].(string)
		userCode, ok := request.Params.Arguments["user_code"].(string)
		if !ok {
			userCode = ""
		}

		var result *model.GetMembersResponse
		switch op {
		case "find_member_id":
			var err error
			result, err = account.NewDefaultAccount().GetMembers(*token, memberName, userCode)
			if err != nil {
				return nil, err
			}
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result.RawJSON)), nil
	})

	accountGetMember := mcp.NewTool("dooray_account_member",
		mcp.WithDescription("find dooray account members by id"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find member details include nickname, english name, native name by member id)"),
			mcp.Enum("find_member_details"),
		),
		mcp.WithString("member_id",
			mcp.Required(),
			mcp.Description("member id"),
		),
	)

	s.AddTool(accountGetMember, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		memberId := request.Params.Arguments["member_id"].(string)

		var result *model.GetMemberResponse
		switch op {
		case "find_member_details":
			var err error
			result, err = account.NewDefaultAccount().GetMember(*token, memberId)
			if err != nil {
				return nil, err
			}
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result.RawJSON)), nil
	})
}
