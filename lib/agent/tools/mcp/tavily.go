package mcp

import (
	"context"
	"fmt"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

const TAVILY_APIKEY = "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"

func TavilyMCPToolset(ctx context.Context) (tool.Toolset, error) {
	return mcptoolset.New(mcptoolset.Config{
		Transport: &gomcp.StreamableClientTransport{
			Endpoint: fmt.Sprintf("https://mcp.tavily.com/mcp/?tavilyApiKey=%s", TAVILY_APIKEY),
		},
	})
}
