package main

import (
	"flag"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
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

	MessengerTools(s, token)

	AccountTools(s, token)

	CalendarTools(s, token)

	ProjectTools(s, token)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
