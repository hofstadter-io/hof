package mcp

import (
	"context"
	"os"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
)

func GithubMCPTransport(ctx context.Context) gomcp.Transport {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	return &gomcp.StreamableClientTransport{
		Endpoint:   "https://api.githubcopilot.com/mcp/",
		HTTPClient: oauth2.NewClient(ctx, ts),
	}
}
