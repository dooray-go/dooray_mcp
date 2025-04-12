package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dooray-go/dooray/openapi/account"
	"github.com/dooray-go/dooray/openapi/calendar"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/dooray-go/dooray/openapi/messenger"
)

func main() {
	token := flag.String("token", "", "개인설정 > API > 개인 인증 토큰 메뉴에서 생성할 수 있습니다.")
	flag.Parse()

	if *token == "" {
		fmt.Printf("token must be set!!")
		return
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"dooray",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

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
			err := messenger.NewDefaultMessenger().DirectSend(*token,
				&messenger.DirectSendRequest{
					Text:                 message,
					OrganizationMemberId: to,
				})

			if err != nil {
				return nil, err
			}
			result = "Message sent successfully"
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})

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

		var result string
		switch op {
		case "find_member_id":
			var err error
			result, err = account.NewDefaultAccount().GetMembers(*token, memberName, userCode)
			if err != nil {
				return nil, err
			}
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
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

		var result string
		switch op {
		case "find_member_details":
			var err error
			result, err = account.NewDefaultAccount().GetMember(*token, memberId)
			if err != nil {
				return nil, err
			}
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})

	calendarGetCalendars := mcp.NewTool("dooray_calendar_calendars",
		mcp.WithDescription("find dooray calendars"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find calendars)"),
			mcp.Enum("find_calendars"),
		),
	)

	s.AddTool(calendarGetCalendars, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)

		var result string
		switch op {
		case "find_calendars":
			var err error
			result, err = calendar.NewDefaultCalendar().GetCalendars(*token)
			if err != nil {
				return nil, err
			}

		}

		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})

	calendarGetEvents := mcp.NewTool("dooray_calendar_events",
		mcp.WithDescription("find dooray events of calendars"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (find calendars)"),
			mcp.Enum("find_events"),
		),
		mcp.WithString("calendars",
			mcp.Description("calender ids to find events. calender ids are separated by commas, calendar ids can be found in the find_calendars operation"),
		),
		mcp.WithString("timeMin",
			mcp.Required(),
			mcp.Description("inclusive start time to find events, format must be format ISO 8601, example: 2025-04-11T00:00:00+09:00"),
		),
		mcp.WithString("timeMax",
			mcp.Required(),
			mcp.Description("exclusive end time to find events, format must be format ISO 8601, example: 2025-04-12T00:00:00+09:00"),
		),
	)

	s.AddTool(calendarGetEvents, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)

		var result string
		switch op {
		case "find_events":
			calendars, ok := request.Params.Arguments["calendars"].(string)
			if !ok {
				calendars = ""
			}
			timeMin, ok := request.Params.Arguments["timeMin"].(string)
			if !ok {
				timeMin = ""
			}
			timeMax, ok := request.Params.Arguments["timeMax"].(string)
			if !ok {
				timeMax = ""
			}

			fromTime, err := calendar.ConvertISO8601ToTime(timeMin)
			if err != nil {
				return nil, err
			}

			toTime, err := calendar.ConvertISO8601ToTime(timeMax)
			if err != nil {
				return nil, err
			}

			result, err = calendar.NewDefaultCalendar().GetEvents(*token, calendars, fromTime, toTime)
			if err != nil {
				return nil, err
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
