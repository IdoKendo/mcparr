package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"MCP Arr",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	cfg := NewConfig()

	s.AddTools(
		RequestDownload(cfg),
		SearchMediaID(cfg),
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
