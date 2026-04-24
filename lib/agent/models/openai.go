// Copyright 2025 achetronic
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package openai provides an OpenAI-compatible LLM implementation for the ADK.
// It supports both native OpenAI API and compatible providers like Ollama.
package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func OpenAI(ctx context.Context, model, baseurl string) (model.LLM, error) {
	fmt.Println("openai:", model, baseurl)
	m := New(Config{
		APIKey:    os.Getenv("OPENAI_API_KEY"),
		BaseURL:   baseurl, // For Ollama
		ModelName: model,   // Or "qwen3:8b" for Ollama
	})
	return m, nil
	// use := os.Getenv("GOOGLE_GENAI_USE_VERTEXAI")
	// if use != "" { // todo, be more truthy
	// 	// creds file
	// 	// creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	// 	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	// 	loc := os.Getenv("GOOGLE_CLOUD_LOCATION")
	// 	if loc == "" {
	// 		loc = "global"
	// 	}

	// 	fmt.Println("vertex:", proj, loc)

	// 	// default inference (same as Go SDK) (typically a service account)
	// 	return gemini.NewModel(ctx, model, &genai.ClientConfig{
	// 		Project:  proj,
	// 		Location: loc,
	// 		Backend:  genai.BackendVertexAI,
	// 	})
	// }

	// // default inference (same as Go SDK)
	// return gemini.NewModel(ctx, model, &genai.ClientConfig{})
}

var _ model.LLM = &Model{}

var (
	ErrNoChoicesInResponse = errors.New("no choices in OpenAI response")
)

// OpenAI enforces a 40-character limit on tool_call_id fields.
const maxToolCallIDLength = 40

// Model implements model.LLM using the official OpenAI Go SDK.
// Works with OpenAI API and compatible providers (Ollama, vLLM, etc.).
type Model struct {
	client    *openai.Client
	modelName string

	// toolCallIDMap stores original IDs when they exceed OpenAI's limit.
	// Keys are shortened hashes, values are original IDs.
	toolCallIDMap   map[string]string
	toolCallIDMapMu sync.RWMutex
}

// HTTPOptions holds optional HTTP-level configuration for the OpenAI client.
type HTTPOptions struct {
	Headers http.Header
}

// Config holds the configuration for creating an OpenAI Model.
type Config struct {
	// APIKey for authentication. Falls back to OPENAI_API_KEY env var if empty.
	APIKey string
	// BaseURL for the API endpoint. Use for OpenAI-compatible providers.
	// Example: "http://localhost:11434/v1" for Ollama.
	BaseURL string
	// ModelName specifies which model to use (e.g., "gpt-4o", "qwen3:8b").
	ModelName string
	// HTTPOptions holds optional HTTP-level overrides (e.g. extra headers).
	HTTPOptions HTTPOptions
}

// New creates a new OpenAI Model with the given configuration.
func New(cfg Config) *Model {
	var opts []option.RequestOption

	if cfg.APIKey != "" {
		opts = append(opts, option.WithAPIKey(cfg.APIKey))
	}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	for k, vals := range cfg.HTTPOptions.Headers {
		for _, v := range vals {
			opts = append(opts, option.WithHeaderAdd(k, v))
		}
	}

	client := openai.NewClient(opts...)

	return &Model{
		client:        &client,
		modelName:     cfg.ModelName,
		toolCallIDMap: make(map[string]string),
	}
}

// Name returns the model name.
func (m *Model) Name() string {
	return m.modelName
}

// GenerateContent sends a request to the LLM and returns responses.
// Set stream=true for streaming responses, false for a single response.
func (m *Model) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	fmt.Println("openai: generating content!", req.Model, stream)
	if stream {
		return m.generateStream(ctx, req)
	}
	return m.generate(ctx, req)
}

// generate sends a non-streaming request and yields a single response.
func (m *Model) generate(ctx context.Context, req *model.LLMRequest) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		params, err := m.buildChatCompletionParams(req)
		if err != nil {
			fmt.Println("gen.params.error:", err)
			yield(nil, err)
			return
		}

		resp, err := m.client.Chat.Completions.New(ctx, params)
		if err != nil {
			fmt.Println("gen.chat.error:", err)
			yield(nil, err)
			return
		}

		llmResp, err := m.convertResponse(resp)
		if err != nil {
			fmt.Println("gen.resp.error:", err)
			yield(nil, err)
			return
		}

		yield(llmResp, nil)
	}
}

// generateStream sends a streaming request and yields partial responses
// as they arrive, followed by a final aggregated response.
func (m *Model) generateStream(ctx context.Context, req *model.LLMRequest) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		params, err := m.buildChatCompletionParams(req)
		if err != nil {
			fmt.Println("gen.params.error:", err)
			yield(nil, err)
			return
		}

		stream := m.client.Chat.Completions.NewStreaming(ctx, params)
		acc := openai.ChatCompletionAccumulator{}

		// Yield partial responses as chunks arrive
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				llmResp := &model.LLMResponse{
					Content: &genai.Content{
						Role:  genai.RoleModel,
						Parts: []*genai.Part{{Text: chunk.Choices[0].Delta.Content}},
					},
					Partial:      true,
					TurnComplete: false,
				}
				if !yield(llmResp, nil) {
					return
				}
			}
		}

		if err := stream.Err(); err != nil {
			fmt.Println("gen.stream.error:", err)
			yield(nil, err)
			return
		}

		// Build and yield final aggregated response
		yield(m.buildStreamFinalResponse(&acc), nil)
	}
}

// buildStreamFinalResponse creates the final LLMResponse from accumulated stream data.
func (m *Model) buildStreamFinalResponse(acc *openai.ChatCompletionAccumulator) *model.LLMResponse {
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	if len(acc.Choices) > 0 {
		choice := acc.Choices[0]

		if choice.Message.Content != "" {
			content.Parts = append(content.Parts, &genai.Part{Text: choice.Message.Content})
		}

		for _, tc := range choice.Message.ToolCalls {
			content.Parts = append(content.Parts, &genai.Part{
				FunctionCall: &genai.FunctionCall{
					ID:   tc.ID,
					Name: tc.Function.Name,
					Args: parseJSONArgs(tc.Function.Arguments),
				},
			})
		}
	}

	var finishReason genai.FinishReason
	if len(acc.Choices) > 0 {
		finishReason = convertFinishReason(string(acc.Choices[0].FinishReason))
	}

	return &model.LLMResponse{
		Content:       content,
		UsageMetadata: convertUsageMetadata(acc.Usage),
		FinishReason:  finishReason,
		Partial:       false,
		TurnComplete:  true,
	}
}

// buildChatCompletionParams converts an LLMRequest into OpenAI API parameters.
func (m *Model) buildChatCompletionParams(req *model.LLMRequest) (openai.ChatCompletionNewParams, error) {
	var messages []openai.ChatCompletionMessageParamUnion

	// Add system instruction
	if req.Config != nil && req.Config.SystemInstruction != nil {
		if text := extractText(req.Config.SystemInstruction); text != "" {
			messages = append(messages, openai.SystemMessage(text))
		}
	}

	// Convert conversation messages
	for _, content := range req.Contents {
		msgs, err := m.convertContentToMessages(content)
		if err != nil {
			return openai.ChatCompletionNewParams{}, err
		}
		messages = append(messages, msgs...)
	}

	params := openai.ChatCompletionNewParams{
		Model:           openai.ChatModel(m.modelName),
		Messages:        messages,
		ReasoningEffort: "low",
	}

	// Apply optional configuration
	if req.Config != nil {
		m.applyGenerationConfig(&params, req.Config)
	}

	return params, nil
}

// applyGenerationConfig applies optional generation settings to the request params.
func (m *Model) applyGenerationConfig(params *openai.ChatCompletionNewParams, cfg *genai.GenerateContentConfig) {
	if cfg.Temperature != nil {
		params.Temperature = openai.Float(float64(*cfg.Temperature))
	}
	if cfg.MaxOutputTokens > 0 {
		params.MaxTokens = openai.Int(int64(cfg.MaxOutputTokens))
	}
	if cfg.TopP != nil {
		params.TopP = openai.Float(float64(*cfg.TopP))
	}

	// Stop sequences
	if len(cfg.StopSequences) == 1 {
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfString: openai.String(cfg.StopSequences[0]),
		}
	} else if len(cfg.StopSequences) > 1 {
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: cfg.StopSequences,
		}
	}

	// Reasoning effort (for o-series models)
	if cfg.ThinkingConfig != nil {
		params.ReasoningEffort = convertThinkingLevel(cfg.ThinkingConfig.ThinkingLevel)
	}

	// JSON mode
	if cfg.ResponseMIMEType == "application/json" {
		params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		}
	}

	// Structured output with schema
	if cfg.ResponseSchema != nil {
		if schemaMap, err := convertSchema(cfg.ResponseSchema); err == nil {
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "response",
						Description: openai.String(cfg.ResponseSchema.Description),
						Schema:      schemaMap,
						Strict:      openai.Bool(true),
					},
				},
			}
		}
	}

	// Tools
	if len(cfg.Tools) > 0 {
		if tools, err := m.convertTools(cfg.Tools); err == nil {
			params.Tools = tools
		}
	}
}

// convertContentToMessages converts a genai.Content into OpenAI message format.
// Handles text, images, audio, files, function calls, and function responses.
func (m *Model) convertContentToMessages(content *genai.Content) ([]openai.ChatCompletionMessageParamUnion, error) {
	var messages []openai.ChatCompletionMessageParamUnion
	var textParts []string
	var toolCalls []openai.ChatCompletionMessageToolCallUnionParam
	var mediaParts []openai.ChatCompletionContentPartUnionParam

	for _, part := range content.Parts {
		switch {
		case part.FunctionResponse != nil:
			responseJSON, err := json.Marshal(part.FunctionResponse.Response)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal function response: %w", err)
			}
			normalizedID := m.normalizeToolCallID(part.FunctionResponse.ID)
			messages = append(messages, openai.ToolMessage(string(responseJSON), normalizedID))

		case part.FunctionCall != nil:
			argsJSON, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal function args: %w", err)
			}
			normalizedID := m.normalizeToolCallID(part.FunctionCall.ID)
			toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallUnionParam{
				OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
					ID: normalizedID,
					Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
						Name:      part.FunctionCall.Name,
						Arguments: string(argsJSON),
					},
				},
			})

		case part.Text != "":
			textParts = append(textParts, part.Text)

		// case part.Thought != "":
		// 	thoughtParts = append(textParts, part.Thought)

		case part.InlineData != nil:
			p, err := convertInlineDataToPart(part.InlineData)
			if err != nil {
				return nil, err
			}
			mediaParts = append(mediaParts, *p)
		}
	}

	if len(textParts) > 0 || len(mediaParts) > 0 || len(toolCalls) > 0 {
		msg := m.buildRoleMessage(content.Role, textParts, mediaParts, toolCalls)
		if msg != nil {
			messages = append(messages, *msg)
		}
	}

	return messages, nil
}

// buildRoleMessage creates the appropriate message type based on role.
func (m *Model) buildRoleMessage(role string, texts []string, media []openai.ChatCompletionContentPartUnionParam, toolCalls []openai.ChatCompletionMessageToolCallUnionParam) *openai.ChatCompletionMessageParamUnion {
	switch convertRole(role) {
	case "user":
		return buildUserMessage(texts, media)
	case "assistant":
		return buildAssistantMessage(texts, toolCalls)
	case "system":
		msg := openai.SystemMessage(joinTexts(texts))
		return &msg
	}
	return nil
}

// buildUserMessage creates a user message, with multi-part support for media.
func buildUserMessage(texts []string, media []openai.ChatCompletionContentPartUnionParam) *openai.ChatCompletionMessageParamUnion {
	if len(media) == 0 {
		msg := openai.UserMessage(joinTexts(texts))
		return &msg
	}

	var parts []openai.ChatCompletionContentPartUnionParam
	for _, text := range texts {
		parts = append(parts, openai.ChatCompletionContentPartUnionParam{
			OfText: &openai.ChatCompletionContentPartTextParam{Text: text},
		})
	}
	parts = append(parts, media...)

	return &openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfArrayOfContentParts: parts,
			},
		},
	}
}

// buildAssistantMessage creates an assistant message with optional tool calls.
func buildAssistantMessage(texts []string, toolCalls []openai.ChatCompletionMessageToolCallUnionParam) *openai.ChatCompletionMessageParamUnion {
	msg := openai.ChatCompletionAssistantMessageParam{}

	if len(texts) > 0 {
		msg.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
			OfString: openai.String(joinTexts(texts)),
		}
	}
	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}

	return &openai.ChatCompletionMessageParamUnion{OfAssistant: &msg}
}

// convertResponse transforms an OpenAI response into an LLMResponse.
func (m *Model) convertResponse(resp *openai.ChatCompletion) (*model.LLMResponse, error) {
	if len(resp.Choices) == 0 {
		return nil, ErrNoChoicesInResponse
	}

	choice := resp.Choices[0]
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	if choice.Message.Content != "" {
		content.Parts = append(content.Parts, &genai.Part{Text: choice.Message.Content})
	}

	for _, tc := range choice.Message.ToolCalls {
		content.Parts = append(content.Parts, &genai.Part{
			FunctionCall: &genai.FunctionCall{
				ID:   tc.ID,
				Name: tc.Function.Name,
				Args: parseJSONArgs(tc.Function.Arguments),
			},
		})
	}

	return &model.LLMResponse{
		Content:       content,
		UsageMetadata: convertUsageMetadata(resp.Usage),
		FinishReason:  convertFinishReason(string(choice.FinishReason)),
		TurnComplete:  true,
	}, nil
}

// convertTools transforms genai tools into OpenAI function tool format.
func (m *Model) convertTools(genaiTools []*genai.Tool) ([]openai.ChatCompletionToolUnionParam, error) {
	var tools []openai.ChatCompletionToolUnionParam

	for _, genaiTool := range genaiTools {
		if genaiTool == nil {
			continue
		}

		for _, funcDecl := range genaiTool.FunctionDeclarations {
			params := funcDecl.ParametersJsonSchema
			if params == nil {
				params = funcDecl.Parameters
			}

			tools = append(tools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
				Name:        funcDecl.Name,
				Description: openai.String(funcDecl.Description),
				Parameters:  convertToFunctionParams(params),
			}))
		}
	}

	return tools, nil
}

// convertToFunctionParams converts various parameter types to OpenAI format.
// OpenAI requires object schemas to have a "properties" field, even if empty.
func convertToFunctionParams(params any) shared.FunctionParameters {
	if params == nil {
		return nil
	}

	var m map[string]any

	// Direct map
	if dm, ok := params.(map[string]any); ok {
		m = dm
	} else {
		// Convert via JSON for other types (e.g., *jsonschema.Schema)
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil
		}
		if json.Unmarshal(jsonBytes, &m) != nil {
			return nil
		}
	}

	// OpenAI requires "properties" for object types
	ensureObjectProperties(m)

	return shared.FunctionParameters(m)
}

// ensureObjectProperties recursively ensures all object schemas have a properties field.
func ensureObjectProperties(schema map[string]any) {
	if schema == nil {
		return
	}

	// If type is "object" and no properties, add empty properties
	if t, ok := schema["type"].(string); ok && t == "object" {
		if _, hasProps := schema["properties"]; !hasProps {
			schema["properties"] = map[string]any{}
		}
	}

	// Recursively process nested properties
	if props, ok := schema["properties"].(map[string]any); ok {
		for _, prop := range props {
			if propMap, ok := prop.(map[string]any); ok {
				ensureObjectProperties(propMap)
			}
		}
	}

	// Process array items
	if items, ok := schema["items"].(map[string]any); ok {
		ensureObjectProperties(items)
	}
}

// convertSchema recursively converts a genai.Schema to OpenAI JSON schema format.
func convertSchema(schema *genai.Schema) (map[string]any, error) {
	if schema == nil {
		return map[string]any{"type": "object", "properties": map[string]any{}}, nil
	}

	result := make(map[string]any)

	if schema.Type != genai.TypeUnspecified {
		result["type"] = schemaTypeToString(schema.Type)
	}
	if schema.Description != "" {
		result["description"] = schema.Description
	}
	if len(schema.Required) > 0 {
		result["required"] = schema.Required
	}
	if len(schema.Enum) > 0 {
		result["enum"] = schema.Enum
	}

	if len(schema.Properties) > 0 {
		props := make(map[string]any)
		for name, propSchema := range schema.Properties {
			converted, err := convertSchema(propSchema)
			if err != nil {
				return nil, err
			}
			props[name] = converted
		}
		result["properties"] = props
	}

	if schema.Items != nil {
		items, err := convertSchema(schema.Items)
		if err != nil {
			return nil, err
		}
		result["items"] = items
	}

	return result, nil
}

// normalizeToolCallID shortens IDs exceeding OpenAI's 40-char limit using a hash.
// The mapping is stored to allow reverse lookup if needed.
func (m *Model) normalizeToolCallID(id string) string {
	if len(id) <= maxToolCallIDLength {
		return id
	}

	hash := sha256.Sum256([]byte(id))
	shortID := "tc_" + hex.EncodeToString(hash[:])[:maxToolCallIDLength-3]

	m.toolCallIDMapMu.Lock()
	m.toolCallIDMap[shortID] = id
	m.toolCallIDMapMu.Unlock()

	return shortID
}

// denormalizeToolCallID restores the original ID from a shortened one.
func (m *Model) denormalizeToolCallID(shortID string) string {
	m.toolCallIDMapMu.RLock()
	defer m.toolCallIDMapMu.RUnlock()

	if original, exists := m.toolCallIDMap[shortID]; exists {
		return original
	}
	return shortID
}

// --- Helper functions ---

// convertInlineDataToPart converts inline data to the appropriate OpenAI content part.
// Supports images (as data URI), audio (wav, mp3), and generic files (PDF, etc.).
// Returns an error for unsupported MIME types, matching Gemini's behavior of letting
// the request fail rather than silently dropping content.
func convertInlineDataToPart(data *genai.Blob) (*openai.ChatCompletionContentPartUnionParam, error) {
	if data == nil {
		return nil, fmt.Errorf("inline data is nil")
	}

	mediaType := data.MIMEType
	base64Data := base64.StdEncoding.EncodeToString(data.Data)

	switch {
	case mediaType == "image/jpeg" || mediaType == "image/jpg" || mediaType == "image/png" ||
		mediaType == "image/gif" || mediaType == "image/webp":
		return &openai.ChatCompletionContentPartUnionParam{
			OfImageURL: &openai.ChatCompletionContentPartImageParam{
				ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
					URL:    fmt.Sprintf("data:%s;base64,%s", mediaType, base64Data),
					Detail: "auto",
				},
			},
		}, nil

	case mediaType == "audio/wav" || mediaType == "audio/mp3" ||
		mediaType == "audio/mpeg" || mediaType == "audio/webm":
		format := "wav"
		if mediaType == "audio/mp3" || mediaType == "audio/mpeg" {
			format = "mp3"
		}
		return &openai.ChatCompletionContentPartUnionParam{
			OfInputAudio: &openai.ChatCompletionContentPartInputAudioParam{
				InputAudio: openai.ChatCompletionContentPartInputAudioInputAudioParam{
					Data:   base64Data,
					Format: format,
				},
			},
		}, nil

	case mediaType == "application/pdf" || strings.HasPrefix(mediaType, "text/"):
		return &openai.ChatCompletionContentPartUnionParam{
			OfFile: &openai.ChatCompletionContentPartFileParam{
				File: openai.ChatCompletionContentPartFileFileParam{
					FileData: openai.String(fmt.Sprintf("data:%s;base64,%s", mediaType, base64Data)),
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported inline data MIME type for OpenAI: %s", mediaType)
	}
}

// convertUsageMetadata converts OpenAI usage stats to genai format.
func convertUsageMetadata(usage openai.CompletionUsage) *genai.GenerateContentResponseUsageMetadata {
	if usage.TotalTokens == 0 {
		return nil
	}
	return &genai.GenerateContentResponseUsageMetadata{
		PromptTokenCount:     int32(usage.PromptTokens),
		CandidatesTokenCount: int32(usage.CompletionTokens),
		TotalTokenCount:      int32(usage.TotalTokens),
	}
}

// convertRole maps genai roles to OpenAI roles.
func convertRole(role string) string {
	if role == "model" {
		return "assistant"
	}
	return role // "user" and "system" are the same
}

// convertFinishReason maps OpenAI finish reasons to genai format.
func convertFinishReason(reason string) genai.FinishReason {
	switch reason {
	case "stop", "tool_calls", "function_call":
		return genai.FinishReasonStop
	case "length":
		return genai.FinishReasonMaxTokens
	case "content_filter":
		return genai.FinishReasonSafety
	default:
		return genai.FinishReasonUnspecified
	}
}

// convertThinkingLevel maps genai thinking levels to OpenAI reasoning effort.
func convertThinkingLevel(level genai.ThinkingLevel) shared.ReasoningEffort {
	switch level {
	case genai.ThinkingLevelLow:
		return shared.ReasoningEffortLow
	case genai.ThinkingLevelHigh:
		return shared.ReasoningEffortHigh
	default:
		return shared.ReasoningEffortMedium
	}
}

// schemaTypeToString converts genai.Type to JSON schema type string.
func schemaTypeToString(t genai.Type) string {
	types := map[genai.Type]string{
		genai.TypeString:  "string",
		genai.TypeNumber:  "number",
		genai.TypeInteger: "integer",
		genai.TypeBoolean: "boolean",
		genai.TypeArray:   "array",
		genai.TypeObject:  "object",
	}
	if s, ok := types[t]; ok {
		return s
	}
	return "string"
}

// extractText extracts all text parts from a Content and joins them.
func extractText(content *genai.Content) string {
	if content == nil {
		return ""
	}
	var texts []string
	for _, part := range content.Parts {
		if part.Text != "" {
			texts = append(texts, part.Text)
		}
	}
	return joinTexts(texts)
}

// joinTexts joins multiple text strings with newlines.
func joinTexts(texts []string) string {
	return strings.Join(texts, "\n")
}

// parseJSONArgs parses a JSON string into a map. Returns empty map on error.
func parseJSONArgs(argsJSON string) map[string]any {
	if argsJSON == "" {
		return make(map[string]any)
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return make(map[string]any)
	}
	return args
}
