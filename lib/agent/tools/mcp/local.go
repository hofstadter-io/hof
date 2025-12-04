package mcp

import (
	"context"
	"fmt"
	"log"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	City string `json:"city" jsonschema:"city name"`
}

type Output struct {
	WeatherSummary string `json:"weather_summary" jsonschema:"weather summary in the given city"`
}

func GetWeather(ctx context.Context, req *gomcp.CallToolRequest, input Input) (*gomcp.CallToolResult, Output, error) {
	return nil, Output{
		WeatherSummary: fmt.Sprintf("Today in %q is sunny\n", input.City),
	}, nil
}

func localMCPTransport(ctx context.Context) gomcp.Transport {
	clientTransport, serverTransport := gomcp.NewInMemoryTransports()

	// Run in-memory MCP server.
	server := gomcp.NewServer(&gomcp.Implementation{Name: "weather_server", Version: "v1.0.0"}, nil)
	gomcp.AddTool(server, &gomcp.Tool{Name: "get_weather", Description: "returns weather in the given city"}, GetWeather)
	_, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	return clientTransport
}
