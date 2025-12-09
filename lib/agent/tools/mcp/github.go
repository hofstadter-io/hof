package mcp

import (
	"context"
	"os"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func GithubMCPToolset(ctx context.Context) (tool.Toolset, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	return mcptoolset.New(mcptoolset.Config{
		Transport: &gomcp.StreamableClientTransport{
			Endpoint:   "https://api.githubcopilot.com/mcp/",
			HTTPClient: oauth2.NewClient(ctx, ts),
		},
	})
}
