package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app/fileutil"
)

// advancedChatInputs adds convenience inputs without decoding multimodal/tool
// content into the legacy string-only Message type.
func advancedChatInputs(cmd *cobra.Command, args []string, body map[string]json.RawMessage) error {
	source, _ := cmd.Flags().GetString("request-json")
	text := strings.Join(args, " ")
	if cmd == visionCmd {
		if f := cmd.Flag("prompt"); f != nil && f.Changed {
			text = f.Value.String()
		}
	}
	file := ""
	if f := cmd.Flag("file"); f != nil && f.Changed {
		file = f.Value.String()
	}
	systemValue := ""
	if f := cmd.Flag("system"); f != nil && f.Changed {
		systemValue = f.Value.String()
	} else if source == "" {
		systemValue = viper.GetString("system")
	}
	systemStdin := systemValue == "-" || systemValue == "/dev/stdin"
	fileStdin := file == "-" || file == "/dev/stdin"
	stdin, err := readAdvancedInputStdin(cmd, source, systemStdin, fileStdin)
	if err != nil {
		return err
	}
	if systemStdin {
		systemValue = stdin
	} else if systemValue != "" {
		resolved, _, err := resolveSystemPrompt(systemValue)
		if err != nil {
			return err
		}
		systemValue = resolved
	}
	var messages []json.RawMessage
	if raw, ok := body["messages"]; ok {
		if err := json.Unmarshal(raw, &messages); err != nil {
			return fmt.Errorf("messages must be an array: %w", err)
		}
	}
	if systemValue != "" {
		messages = replaceAdvancedSystemMessages(messages, systemValue)
	}
	var message json.RawMessage
	if cmd == visionCmd && file != "" {
		message, err = advancedVisionMessage(file, text, fileStdin)
	} else {
		message, err = advancedTextMessage(text, file, stdin, systemStdin, fileStdin)
	}
	if err != nil {
		return err
	}
	if message != nil {
		messages = append(messages, message)
	}
	if cmd == visionCmd && len(messages) == 0 {
		return fmt.Errorf("image required: use -f or supply multimodal messages with --request-json")
	}
	if len(messages) > 0 {
		raw, err := json.Marshal(messages)
		if err != nil {
			return err
		}
		body["messages"] = raw
	}
	return applyAdvancedThinking(cmd, body)
}

func replaceAdvancedSystemMessages(messages []json.RawMessage, systemValue string) []json.RawMessage {
	// Explicit system input replaces previous system messages; other message
	// fields (including reasoning and tool calls) remain byte-preserved.
	kept := make([]json.RawMessage, 0, len(messages)+1)
	systemMessage, _ := json.Marshal(map[string]string{"role": "system", "content": systemValue})
	kept = append(kept, systemMessage)
	for _, raw := range messages {
		var role struct {
			Role string `json:"role"`
		}
		if json.Unmarshal(raw, &role) == nil && role.Role == "system" {
			continue
		}
		kept = append(kept, raw)
	}
	return kept
}

func advancedVisionMessage(file, text string, fileStdin bool) (json.RawMessage, error) {
	if fileStdin {
		return nil, fmt.Errorf("vision binary input requires a local file path or HTTP(S) URL")
	}
	image := file
	if !strings.HasPrefix(image, "http://") && !strings.HasPrefix(image, "https://") {
		var err error
		image, err = encodeLocalImage(file, fileutil.OSFileReader{})
		if err != nil {
			return nil, err
		}
	}
	if text == "" {
		text = "What do you see in this image? Please describe it in detail."
	}
	message, _ := json.Marshal(map[string]any{"role": "user", "content": []any{map[string]string{"type": "text", "text": text}, map[string]any{"type": "image_url", "image_url": map[string]string{"url": image}}}})
	return message, nil
}

func readAdvancedContextFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("open context file: %w", err)
	}
	raw, readErr := io.ReadAll(io.LimitReader(f, MaxFileSize+1))
	closeErr := f.Close()
	if readErr != nil {
		return "", readErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if len(raw) > MaxFileSize {
		return "", fmt.Errorf("context file exceeds %d bytes", MaxFileSize)
	}
	return string(raw), nil
}

func advancedTextMessage(text, file, stdin string, systemStdin, fileStdin bool) (json.RawMessage, error) {
	if file != "" {
		data := stdin
		if !fileStdin {
			var err error
			data, err = readAdvancedContextFile(file)
			if err != nil {
				return nil, err
			}
		}
		text = "<file>\n" + data + "\n</file>\n\n" + text
	}
	if stdin != "" && !systemStdin && !fileStdin {
		text = "<stdin>\n" + stdin + "\n</stdin>\n\n" + text
	}
	if text == "" {
		return nil, nil
	}
	message, _ := json.Marshal(map[string]string{"role": "user", "content": text})
	return message, nil
}

func applyAdvancedThinking(cmd *cobra.Command, body map[string]json.RawMessage) error {
	if f := cmd.Flag("think"); f != nil && f.Changed {
		enabled, _ := cmd.Flags().GetBool("think")
		kind := "disabled"
		if enabled {
			kind = "enabled"
		}
		thinking := make(map[string]json.RawMessage)
		if raw := body["thinking"]; len(raw) > 0 {
			if err := json.Unmarshal(raw, &thinking); err != nil || thinking == nil {
				return fmt.Errorf("thinking must be an object")
			}
		}
		thinking["type"], _ = json.Marshal(kind)
		raw, _ := json.Marshal(thinking)
		body["thinking"] = raw
	}
	return nil
}

func readAdvancedInputStdin(cmd *cobra.Command, source string, systemStdin, fileStdin bool) (string, error) {
	if (source == "-" && (systemStdin || fileStdin)) || (systemStdin && fileStdin) {
		return "", fmt.Errorf("request JSON, system prompt, and file cannot share stdin")
	}
	var stdin string
	if systemStdin || fileStdin || (source == "" && cmd == rootCmd && hasStdinData()) {
		data, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), MaxStdinSize+1))
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		if len(data) > MaxStdinSize {
			return "", fmt.Errorf("stdin exceeds %d bytes", MaxStdinSize)
		}
		stdin = string(data)
	}
	return stdin, nil
}
