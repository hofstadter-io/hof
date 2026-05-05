package mcp

import (
	"context"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func QuickbooksMCPToolset(ctx context.Context) (tool.Toolset, error) {
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	// )
	return mcptoolset.New(mcptoolset.Config{
		Transport: &gomcp.StreamableClientTransport{
			Endpoint: "http://localhost:8000/callback",
			// HTTPClient: oauth2.NewClient(ctx, ts),
		},
	})
}
