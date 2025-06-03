package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Define a struct for curl command arguments
type CurlArgs struct {
	URL string `json:"url"`
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"MCP Browser",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Add curl tool
	curlTool := mcp.NewTool("curl",
		mcp.WithDescription("Fetch content from a URL using curl"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to fetch"),
		),
	)

	// Add tool handler
	s.AddTool(curlTool, mcp.NewTypedToolHandler(curlHandler))

	// Start the stdio server
	fmt.Println("Starting MCP Browser server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Curl handler function
func curlHandler(ctx context.Context, request mcp.CallToolRequest, args CurlArgs) (*mcp.CallToolResult, error) {
	if args.URL == "" {
		return mcp.NewToolResultError("URL is required"), nil
	}

	// Create and execute the curl command
	cmd := exec.Command("curl", "-s", args.URL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error executing curl: %v", err)), nil
	}

	// Return the output as text
	return mcp.NewToolResultText(string(output)), nil
}

