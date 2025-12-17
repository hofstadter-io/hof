# MCP Tools

This directory contains integrations with the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) to provide toolsets for agents.

## Toolsets

### Github
`github.go`
Provides access to GitHub Copilot API via MCP.
- **Config**: Uses `GITHUB_PAT` environment variable.
- **Transport**: `StreamableClientTransport` to `api.githubcopilot.com`.

### Tavily
`tavily.go`
Provides web search capabilities using Tavily.
- **Config**: Hardcoded API key (DEV ONLY).
- **Transport**: `StreamableClientTransport` to `mcp.tavily.com`.

### Local
`local.go`
An in-memory example of an MCP server providing a weather tool.
- **Features**: `get_weather` tool.
- **Implementation**: Uses `gomcp.NewInMemoryTransports`.

#### Types

```go
type Input struct {
	City string `json:"city" jsonschema:"city name"`
}

type Output struct {
	WeatherSummary string `json:"weather_summary" jsonschema:"weather summary in the given city"`
}
```

### Quickbooks
`quickbooks.go`
Integration for Quickbooks (WIP).

## Usage

Toolsets are initialized using `mcptoolset.New`.

```go
// Example: Github
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
```
