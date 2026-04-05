package main

import (
	"context"
	"fmt"
	"github.com/dooray-go/dooray/openapi/calendar"
	model "github.com/dooray-go/dooray/openapi/model/calendar"
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
		op := request.GetArguments()["operation"].(string)

		var result string
		switch op {
		case "find_calendars":
			var err error
			res, err := calendar.NewDefaultCalendar().GetCalendars(*token)
			if err != nil {
				return nil, err
			}

			result = res.RawJSON
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
		op := request.GetArguments()["operation"].(string)

		var result *model.EventsResponse
		switch op {
		case "find_events":
			calendars, ok := request.GetArguments()["calendars"].(string)
			if !ok {
				calendars = ""
			}
			timeMin, ok := request.GetArguments()["timeMin"].(string)
			if !ok {
				timeMin = ""
			}
			timeMax, ok := request.GetArguments()["timeMax"].(string)
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

		return mcp.NewToolResultText(fmt.Sprintf("%s", result.RawJSON)), nil
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
		mcp.WithBoolean("wholeDayFlag",
			mcp.Description("set true for whole day event. when true, startedAt/endedAt format must be date only, example: 2025-04-11+09:00"),
		),
		mcp.WithString("recurrenceFrequency",
			mcp.Description("recurrence frequency: daily, weekly, monthly, yearly. leave empty for non-recurring event"),
			mcp.Enum("", "daily", "weekly", "monthly", "yearly"),
		),
		mcp.WithNumber("recurrenceInterval",
			mcp.Description("recurrence interval, e.g. 1 = every, 2 = every other. default is 1"),
		),
		mcp.WithString("recurrenceUntil",
			mcp.Description("recurrence end date in ISO 8601, e.g. 2025-12-31T23:59:59+09:00. leave empty for no end"),
		),
		mcp.WithString("recurrenceByday",
			mcp.Description("recurrence by day, e.g. MO,TU,WE for weekly. values: SU,MO,TU,WE,TH,FR,SA. for monthly: 1MO, 2TU, -1FR etc."),
		),
		mcp.WithString("recurrenceBymonth",
			mcp.Description("recurrence by month (1-12), for yearly recurrence"),
		),
		mcp.WithString("recurrenceBymonthday",
			mcp.Description("recurrence by day of month (1-31), for monthly/yearly recurrence"),
		),
		mcp.WithString("recurrenceTimezoneName",
			mcp.Description("timezone for recurrence rule, e.g. Asia/Seoul. defaults to Asia/Seoul"),
		),
	)

	s.AddTool(calendarPostEvent, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.GetArguments()["operation"].(string)

		var result *model.EventResponse
		switch op {
		case "create_event":
			calendarId, ok := request.GetArguments()["calendarId"].(string)
			if !ok {
				calendarId = ""
			}

			subject, ok := request.GetArguments()["subject"].(string)
			if !ok {
				subject = ""
			}
			content, ok := request.GetArguments()["content"].(string)
			if !ok {
				content = ""
			}

			startedAt, ok := request.GetArguments()["startedAt"].(string)
			if !ok {
				startedAt = ""
			}
			endedAt, ok := request.GetArguments()["endedAt"].(string)
			if !ok {
				endedAt = ""
			}

			wholeDayFlag, _ := request.GetArguments()["wholeDayFlag"].(bool)

			startedAtTime, err := utils.ConvertISO8601ToTime(startedAt)
			if err != nil {
				return nil, err
			}

			endedAtTime, err := utils.ConvertISO8601ToTime(endedAt)
			if err != nil {
				return nil, err
			}

			var startJsonTime, endJsonTime utils.JsonTime
			if wholeDayFlag {
				startJsonTime = utils.NewJsonDate(startedAtTime)
				endJsonTime = utils.NewJsonDate(endedAtTime)
			} else {
				startJsonTime = utils.NewJsonTime(startedAtTime)
				endJsonTime = utils.NewJsonTime(endedAtTime)
			}

			var recurrenceRule *model.RecurrenceRule
			recurrenceFrequency, _ := request.GetArguments()["recurrenceFrequency"].(string)
			if recurrenceFrequency != "" {
				recurrenceInterval := 1
				if v, ok := request.GetArguments()["recurrenceInterval"].(float64); ok && v > 0 {
					recurrenceInterval = int(v)
				}
				recurrenceUntil, _ := request.GetArguments()["recurrenceUntil"].(string)
				recurrenceByday, _ := request.GetArguments()["recurrenceByday"].(string)
				recurrenceBymonth, _ := request.GetArguments()["recurrenceBymonth"].(string)
				recurrenceBymonthday, _ := request.GetArguments()["recurrenceBymonthday"].(string)
				recurrenceTimezoneName, _ := request.GetArguments()["recurrenceTimezoneName"].(string)
				if recurrenceTimezoneName == "" {
					recurrenceTimezoneName = "Asia/Seoul"
				}
				recurrenceRule = &model.RecurrenceRule{
					Frequency:    recurrenceFrequency,
					Interval:     recurrenceInterval,
					Until:        recurrenceUntil,
					Byday:        recurrenceByday,
					Bymonth:      recurrenceBymonth,
					Bymonthday:   recurrenceBymonthday,
					TimezoneName: recurrenceTimezoneName,
				}
			}

			event := model.EventRequest{
				Users:            model.Users{},
				Subject:          subject,
				Body:             model.Body{MimeType: "text/html", Content: content},
				StartedAt:        startJsonTime,
				EndedAt:          endJsonTime,
				WholeDayFlag:     wholeDayFlag,
				Location:         "",
				RecurrenceRule:   recurrenceRule,
				PersonalSettings: nil,
			}

			result, err = calendar.NewDefaultCalendar().CreateEvent(*token, calendarId, event)
			if err != nil {
				return nil, err
			}
		}

		return mcp.NewToolResultText(fmt.Sprintf("%s", result.RawJSON)), nil
	})
}
