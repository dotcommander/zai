package cmd

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
)

// Advanced requests retain the complete JSON under --json. The ordinary view
// extracts only standard display fields and falls back to JSON for structured
// responses such as tool calls.
func advancedResultText(cmd *cobra.Command, operation string, result json.RawMessage) string {
	if flag := cmd.Flag("json"); flag != nil && flag.Value.String() == "true" {
		return string(result)
	}
	var raw any
	if json.Unmarshal(result, &raw) == nil && containsToolCall(raw) {
		return string(result)
	}
	var response struct {
		Text       string `json:"text"`
		OutputText string `json:"output_text"`
		ID         string `json:"id"`
		AsyncID    string `json:"async_id"`
		Choices    []struct {
			Message struct {
				Content json.RawMessage `json:"content"`
			}
		} `json:"choices"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
		VideoResult []struct {
			URL string `json:"url"`
		} `json:"video_result"`
	}
	if json.Unmarshal(result, &response) != nil {
		return string(result)
	}
	if response.Text != "" {
		return response.Text
	}
	if response.OutputText != "" {
		return response.OutputText
	}
	var parts []string
	for _, choice := range response.Choices {
		var text string
		if json.Unmarshal(choice.Message.Content, &text) == nil && text != "" {
			parts = append(parts, text)
		}
	}
	for _, block := range response.Content {
		if block.Text != "" {
			parts = append(parts, block.Text)
		}
	}
	for _, output := range response.Output {
		for _, block := range output.Content {
			if block.Text != "" {
				parts = append(parts, block.Text)
			}
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "\n")
	}
	for _, data := range response.Data {
		if data.URL != "" {
			parts = append(parts, data.URL)
		}
	}
	for _, video := range response.VideoResult {
		if video.URL != "" {
			parts = append(parts, video.URL)
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "\n")
	}
	switch operation {
	case "videos", "image_async", "task_get", "agent_run", "agent_result", "agent_conversation":
		if response.AsyncID != "" {
			return response.AsyncID
		}
		if response.ID != "" {
			return response.ID
		}
	}
	return string(result)
}

func containsToolCall(value any) bool {
	switch value := value.(type) {
	case []any:
		for _, item := range value {
			if containsToolCall(item) {
				return true
			}
		}
	case map[string]any:
		if calls, ok := value["tool_calls"].([]any); ok && len(calls) > 0 {
			return true
		}
		if kind, _ := value["type"].(string); kind == "tool_use" || kind == "function_call" {
			return true
		}
		for _, child := range value {
			if containsToolCall(child) {
				return true
			}
		}
	}
	return false
}
