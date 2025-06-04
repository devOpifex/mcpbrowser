package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Define a struct for browser arguments
type BrowserArgs struct {
	URL string `json:"url"`
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"MCP Browser",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Add browser tool
	browserTool := mcp.NewTool("browser",
		mcp.WithDescription("Fetches and extracts content from a web URL. For HTML pages, returns the text content of the body. For non-HTML resources, returns the raw content. Useful for retrieving information from websites, APIs, or documents available online."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to fetch"),
		),
	)

	// Add tool handler
	s.AddTool(browserTool, mcp.NewTypedToolHandler(browserHandler))

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Browser handler function using colly
func browserHandler(ctx context.Context, request mcp.CallToolRequest, args BrowserArgs) (*mcp.CallToolResult, error) {
	if args.URL == "" {
		return mcp.NewToolResultError("URL is required"), nil
	}

	// Create a new collector
	c := colly.NewCollector()

	var content string
	var isHTML bool

	// Check if the response is HTML by inspecting Content-Type header
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		isHTML = strings.Contains(strings.ToLower(contentType), "text/html")

		if !isHTML {
			// If not HTML, return the entire content
			content = string(r.Body)
		}
	})

	// Extract only the body content if it's HTML
	c.OnHTML("body", func(e *colly.HTMLElement) {
		content = e.Text
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		content = fmt.Sprintf("Error fetching URL: %v", err)
	})

	// Visit the URL
	err := c.Visit(args.URL)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error visiting URL: %v", err)), nil
	}

	// Return the content
	return mcp.NewToolResultText(content), nil
}
