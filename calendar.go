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
			mcp.Description("The operation to perform (find events)"),
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

	calendarPostEvent := mcp.NewTool("dooray_calendar_post_event",
		mcp.WithDescription("register dooray events on a calendar"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (register dooray events on a calendar)"),
			mcp.Enum("create_event"),
		),
		mcp.WithString("calendarId",
			mcp.Description("calender id to register a event. calendar ids can be found in the find_calendars operation"),
		),
		mcp.WithString("subject",
			mcp.Required(),
			mcp.Description("event subject"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("event content"),
		),
		mcp.WithString("startedAt",
			mcp.Required(),
			mcp.Description("event start time, format must be format ISO 8601, example: 2025-04-11T00:00:00+09:00"),
		),
		mcp.WithString("endedAt",
			mcp.Required(),
			mcp.Description("event end time, format must be format ISO 8601, example: 2025-04-11T00:00:00+09:00"),
		),
	)

	s.AddTool(calendarPostEvent, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)

		var result string
		switch op {
		case "create_event":
			calendarId, ok := request.Params.Arguments["calendarId"].(string)
			if !ok {
				calendarId = ""
			}

			subject, ok := request.Params.Arguments["subject"].(string)
			if !ok {
				subject = ""
			}
			content, ok := request.Params.Arguments["content"].(string)
			if !ok {
				content = ""
			}

			startedAt, ok := request.Params.Arguments["startedAt"].(string)
			if !ok {
				startedAt = ""
			}
			endedAt, ok := request.Params.Arguments["endedAt"].(string)
			if !ok {
				endedAt = ""
			}

			startedAtTime, err := utils.ConvertISO8601ToTime(startedAt)
			if err != nil {
				return nil, err
			}

			endedAtTime, err := utils.ConvertISO8601ToTime(endedAt)
			if err != nil {
				return nil, err
			}

			event := calendar.EventRequest{
				Users:            calendar.Users{},
				Subject:          subject,
				Body:             calendar.Body{MimeType: "text/html", Content: content},
				StartedAt:        utils.JsonTime(startedAtTime),
				EndedAt:          utils.JsonTime(endedAtTime),
				WholeDayFlag:     false,
				Location:         "",
				RecurrenceRule:   nil,
				PersonalSettings: nil,
			}

			result, err = calendar.NewDefaultCalendar().CreateEvent(*token, calendarId, event)
			if err != nil {
				return nil, err
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("%s", result)), nil
	})
}
