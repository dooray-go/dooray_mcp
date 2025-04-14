package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray/openapi/calendar"
	"github.com/dooray-go/dooray/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CalendarTools(s *server.MCPServer, token *string) {

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

			fromTime, err := utils.ConvertISO8601ToTime(timeMin)
			if err != nil {
				return nil, err
			}

			toTime, err := utils.ConvertISO8601ToTime(timeMax)
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
}
