package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

type advancedAPI interface {
	ExecuteOperation(context.Context, app.OperationRequest) (json.RawMessage, error)
	StreamOperation(context.Context, app.OperationRequest, func(app.Event) error) error
	WaitTask(context.Context, string, time.Duration) (json.RawMessage, error)
	WaitAgent(context.Context, string, string, time.Duration) (json.RawMessage, error)
}

var advancedClientFactory = func() advancedAPI { return newClientWithoutHistory() }

// registerAdvancedCommands preserves the convenience commands while exposing
// lossless request and response contracts for every registered API operation.
func registerAdvancedCommands() {
	for _, entry := range []struct {
		cmd       *cobra.Command
		operation string
	}{
		{rootCmd, "chat"}, {chatCmd, "chat"}, {visionCmd, "chat"}, {imageCmd, "images"},
		{videoCmd, "videos"}, {audioCmd, "audio"}, {searchCmd, "search"}, {readerCmd, "reader"},
	} {
		attachAdvanced(entry.cmd, entry.operation)
	}
	api := &cobra.Command{Use: "api operation", Short: "Execute a registered API operation with lossless JSON", Args: cobra.ExactArgs(1)}
	addAdvancedFlags(api)
	api.RunE = func(cmd *cobra.Command, args []string) error { return executeAdvanced(cmd, nil, args[0], nil) }
	rootCmd.AddCommand(api)
	for _, entry := range []struct{ name, op string }{{"tokenize", "tokenize"}, {"layout", "layout"}, {"responses", "responses"}, {"messages", "messages"}} {
		cmd := newAdvancedCommand(entry.name, entry.op)
		rootCmd.AddCommand(cmd)
	}
	task := &cobra.Command{Use: "task", Short: "Inspect asynchronous tasks"}
	get := newAdvancedCommand("get id", "task_get")
	get.Args = cobra.ExactArgs(1)
	get.RunE = func(cmd *cobra.Command, args []string) error {
		return executeAdvanced(cmd, nil, "task_get", map[string]string{"id": args[0]})
	}
	task.AddCommand(get)
	rootCmd.AddCommand(task)
	files := &cobra.Command{Use: "files", Short: "Manage API files"}
	upload := newAdvancedCommand("upload path", "files_upload")
	upload.Args = cobra.ExactArgs(1)
	upload.Flags().String("purpose", "", "File purpose")
	upload.RunE = func(cmd *cobra.Command, args []string) error { return executeAdvancedUpload(cmd, args[0]) }
	files.AddCommand(upload)
	rootCmd.AddCommand(files)
	agent := &cobra.Command{Use: "agent", Short: "Run beta agents (vendor support may vary)"}
	for _, entry := range []struct{ name, op string }{{"run", "agent_run"}, {"result", "agent_result"}, {"conversation", "agent_conversation"}} {
		cmd := newAdvancedCommand(entry.name+" [id]", entry.op)
		cmd.Args = cobra.MaximumNArgs(1)
		cmd.Flags().String("agent-id", "", "Beta agent identifier")
		agent.AddCommand(cmd)
	}
	rootCmd.AddCommand(agent)
	imageCmd.Flags().Bool("async", false, "Submit an asynchronous image request")
	asyncImage := newAdvancedCommand("async", "image_async")
	imageCmd.AddCommand(asyncImage)
}

func newAdvancedCommand(use, operation string) *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: "Execute " + operation + "; supply all API fields with --request-json", Args: cobra.NoArgs}
	addAdvancedFlags(cmd)
	if operation == "chat" {
		addChatAdvancedFlags(cmd)
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return executeAdvanced(cmd, args, operation, nil) }
	return cmd
}

func addChatAdvancedFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("tool-stream", false, "Stream tool-call arguments, text, and reasoning as NDJSON events")
}

func addAdvancedFlags(cmd *cobra.Command) {
	addAdvancedDownloadFlags(cmd)
	cmd.Flags().String("accept-language", "", "Provider Accept-Language request header")
	cmd.Flags().String("request-json", "", "Read an API JSON object from a file or - (stdin)")
	cmd.Flags().Bool("events", false, "Emit lossless SSE events as NDJSON")
	cmd.Flags().Bool("wait", false, "Wait for an asynchronous task (requires --wait-timeout)")
	cmd.Flags().Duration("wait-timeout", 0, "Caller deadline for asynchronous waiting")
	cmd.Flags().Duration("poll-interval", time.Second, "Interval between asynchronous status requests")
}

func attachAdvanced(cmd *cobra.Command, operation string) {
	addAdvancedFlags(cmd)
	if operation == "chat" {
		addChatAdvancedFlags(cmd)
	}
	run, argsCheck := cmd.RunE, cmd.Args
	advanced := func(cmd *cobra.Command) bool {
		events, _ := cmd.Flags().GetBool("events")
		async, _ := cmd.Flags().GetBool("async")
		rawJSON := false
		if flag := cmd.Flag("json"); flag != nil && flag.Changed && flag.Value.String() == "true" {
			switch operation {
			case "images", "videos", "search", "reader", "chat":
				rawJSON = true
			case "audio":
				stream, _ := cmd.Flags().GetBool("stream")
				rawJSON = !stream
			}
		}
		wait, _ := cmd.Flags().GetBool("wait")
		return cmd.Flags().Changed("request-json") || events || async || rawJSON || wait || cmd.Flags().Changed("tool-stream")
	}
	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if advanced(cmd) {
			return nil
		}
		if argsCheck != nil {
			return argsCheck(cmd, args)
		}
		return nil
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if advanced(cmd) {
			async, _ := cmd.Flags().GetBool("async")
			if async && operation == "images" {
				return executeAdvanced(cmd, args, "image_async", nil)
			}
			return executeAdvanced(cmd, args, operation, nil)
		}
		return run(cmd, args)
	}
}

func readRequestObject(source string, input io.Reader) (map[string]json.RawMessage, error) {
	object := make(map[string]json.RawMessage)
	if source == "" {
		return object, nil
	}
	reader := input
	if source != "-" {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		// The file is read-only; closing it is best-effort cleanup.
		defer func() { _ = file.Close() }()
		reader = file
	}
	data, err := io.ReadAll(io.LimitReader(reader, MaxFileSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > MaxFileSize {
		return nil, fmt.Errorf("request JSON exceeds %d bytes", MaxFileSize)
	}
	if err = json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("request JSON must be an object: %w", err)
	}
	if object == nil || !strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		return nil, fmt.Errorf("request JSON must be an object")
	}
	return object, nil
}

func advancedBody(cmd *cobra.Command, args []string, operation string) (json.RawMessage, error) {
	source, _ := cmd.Flags().GetString("request-json")
	args, err := advancedInputArgs(cmd, args, operation, source)
	if err != nil {
		return nil, err
	}
	body, err := readRequestObject(source, cmd.InOrStdin())
	if err != nil {
		return nil, err
	}
	if err := mergeAdvancedFlags(cmd, operation, body); err != nil {
		return nil, err
	}
	if err := mergeAdvancedArguments(cmd, args, operation, body); err != nil {
		return nil, err
	}
	if err := applyAdvancedDefaults(cmd, operation, body); err != nil {
		return nil, err
	}
	if operation == "chat" {
		if err := advancedChatInputs(cmd, args, body); err != nil {
			return nil, err
		}
	}
	return json.Marshal(body)
}

func advancedInputArgs(cmd *cobra.Command, args []string, operation, source string) ([]string, error) {
	if source == "-" {
		for _, name := range []string{"system", "file"} {
			f := cmd.Flag(name)
			if f != nil && (f.Value.String() == "-" || f.Value.String() == "/dev/stdin") {
				return nil, fmt.Errorf("--request-json - conflicts with --%s stdin input", name)
			}
		}
	}
	if source == "" && len(args) == 0 && hasStdinData() {
		switch operation {
		case "search", "reader", "images", "image_async", "videos":
			data, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), MaxStdinSize+1))
			if err != nil {
				return nil, err
			}
			if len(data) > MaxStdinSize {
				return nil, fmt.Errorf("stdin exceeds %d bytes", MaxStdinSize)
			}
			if text := strings.TrimSpace(string(data)); text != "" {
				args = []string{text}
			}
		}
	}
	return args, nil
}

func mergeAdvancedFlags(cmd *cobra.Command, operation string, body map[string]json.RawMessage) error {
	mappings := advancedMappings(operation)
	var mergeErr, err error
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		field, ok := mappings[flag.Name]
		if !ok || !flag.Changed || mergeErr != nil {
			return
		}
		var value any
		switch flag.Value.Type() {
		case "bool":
			value, mergeErr = cmd.Flags().GetBool(flag.Name)
		case "int":
			value, mergeErr = cmd.Flags().GetInt(flag.Name)
		case "float64":
			value, mergeErr = cmd.Flags().GetFloat64(flag.Name)
		case "stringArray":
			value, mergeErr = cmd.Flags().GetStringArray(flag.Name)
		default:
			value = flag.Value.String()
		}
		if flag.Name == "no-retain-images" {
			value = !value.(bool)
		}
		if operation == "audio" && flag.Name == "hotwords" {
			hotwords := parseHotwords(flag.Value.String())
			if hotwords == nil {
				hotwords = []string{}
			}
			value = hotwords
		}
		if mergeErr == nil {
			body[field], mergeErr = json.Marshal(value)
		}
	})
	// Inherited flags are not part of Flags until Cobra merges them.
	if f := cmd.Flag("model"); f != nil && f.Changed {
		body["model"], err = json.Marshal(f.Value.String())
		if err != nil {
			return err
		}
	}
	if mergeErr != nil {
		return mergeErr
	}
	return nil
}

func mergeAdvancedArguments(cmd *cobra.Command, args []string, operation string, body map[string]json.RawMessage) error {
	var err error
	if len(args) > 0 {
		text := strings.Join(args, " ")
		switch operation {
		case "audio":
			if len(args) != 1 {
				return fmt.Errorf("audio accepts one upload path")
			}
		case "agent_run":
			body["agent_id"], err = json.Marshal(text)
		case "agent_result":
			body["async_id"], err = json.Marshal(text)
		case "agent_conversation":
			body["conversation_id"], err = json.Marshal(text)
		case "chat":
			// The chat input helper merges positional and multimodal inputs without dropping tools.
		default:
			field := "prompt"
			if operation == "search" {
				field = "search_query"
			}
			if operation == "reader" {
				field = "url"
			}
			body[field], err = json.Marshal(text)
		}
	}
	if f := cmd.Flag("agent-id"); f != nil && f.Changed {
		body["agent_id"], err = json.Marshal(f.Value.String())
	}
	if operation == "agent_conversation" {
		if _, ok := body["agent_id"]; !ok {
			body["agent_id"] = json.RawMessage(`"slides_glm_agent"`)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func applyAdvancedDefaults(cmd *cobra.Command, operation string, body map[string]json.RawMessage) error {
	var err error
	if _, exists := body["model"]; !exists {
		key := map[string]string{"chat": "api.model", "images": "api.image_model", "image_async": "api.image_model", "videos": "api.video_model", "audio": "api.audio_model", "tokenize": "api.model"}[operation]
		if cmd == visionCmd {
			key = "api.vision_model"
		}
		if key != "" && viper.GetString(key) != "" {
			body["model"], err = json.Marshal(viper.GetString(key))
			if err != nil {
				return err
			}
		}
	}
	events, _ := cmd.Flags().GetBool("events")
	if operation == "chat" {
		if flag := cmd.Flag("tool-stream"); flag != nil && flag.Changed && flag.Value.String() == "true" {
			body["stream"] = json.RawMessage("true")
		}
	}
	// Some agent schemas default to streaming; make the CLI's ordinary JSON
	// response mode explicit while preserving the caller's chosen mode.
	if operation == "agent_run" {
		if _, exists := body["stream"]; !exists {
			body["stream"] = json.RawMessage("false")
		}
	}
	if events {
		body["stream"] = json.RawMessage("true")
	}
	return nil
}

func advancedMappings(operation string) map[string]string {
	fields := map[string]string{"model": "model", "user-id": "user_id", "request-id": "request_id"}
	add := func(pairs ...string) {
		for i := 0; i < len(pairs); i += 2 {
			fields[pairs[i]] = pairs[i+1]
		}
	}
	switch operation {
	case "images", "image_async":
		add("quality", "quality", "size", "size")
	case "videos":
		add("quality", "quality", "size", "size", "fps", "fps", "duration", "duration", "with-audio", "with_audio", "file", "image_url")
	case "chat":
		add("temperature", "temperature", "top-p", "top_p", "max-tokens", "max_tokens", "tool-stream", "tool_stream")
	case "audio":
		add("prompt", "prompt", "language", "language", "hotwords", "hotwords", "stream", "stream")
	case "search":
		add("count", "count", "recency", "search_recency_filter", "domain", "search_domain_filter", "engine", "search_engine", "content-size", "content_size", "sources", "include_sources", "require-search", "require_search", "search-prompt", "search_prompt", "result-sequence", "result_sequence", "lang", "language")
	case "reader":
		add("format", "return_format", "timeout", "timeout", "no-cache", "no_cache", "no-gfm", "no_gfm", "keep-img-data-url", "keep_img_data_url", "with-images-summary", "with_images_summary", "with-links-summary", "with_links_summary", "no-retain-images", "retain_images")
	case "files_upload":
		add("purpose", "purpose")
	}
	return fields
}

func executeAdvancedUpload(cmd *cobra.Command, path string) error {
	if path == "-" {
		return fmt.Errorf("upload requires a file path; stdin is reserved for --request-json")
	}
	return runAdvanced(cmd, nil, "files_upload", nil, path)
}
func executeAdvanced(cmd *cobra.Command, args []string, operation string, path map[string]string) error {
	return runAdvanced(cmd, args, operation, path, "")
}

func validateAdvancedOptions(cmd *cobra.Command, operation string) error {
	if operation == "audio" {
		stream, _ := cmd.Flags().GetBool("stream")
		events, _ := cmd.Flags().GetBool("events")
		jsonFlag := cmd.Flag("json")
		if stream && !events && jsonFlag != nil && jsonFlag.Value.String() == "true" {
			return fmt.Errorf("--stream and --json require --events for lossless event output")
		}
	}
	waitRequested, _ := cmd.Flags().GetBool("wait")
	if waitRequested && operation != "image_async" && operation != "videos" && operation != "task_get" && operation != "agent_run" && operation != "agent_result" {
		return fmt.Errorf("--wait is unavailable for %s", operation)
	}
	return nil
}

func runAdvanced(cmd *cobra.Command, args []string, operation string, path map[string]string, upload string) error {
	if err := validateAdvancedOptions(cmd, operation); err != nil {
		return err
	}
	body, err := advancedBody(cmd, args, operation)
	if err != nil {
		return err
	}
	events, _ := cmd.Flags().GetBool("events")
	var streamRequest struct {
		Stream bool `json:"stream"`
	}
	_ = json.Unmarshal(body, &streamRequest)
	events = events || streamRequest.Stream
	wait, _ := cmd.Flags().GetBool("wait")
	timeout, _ := cmd.Flags().GetDuration("wait-timeout")
	interval, _ := cmd.Flags().GetDuration("poll-interval")
	if wait && (timeout <= 0 || interval <= 0) {
		return fmt.Errorf("--wait requires positive --wait-timeout and --poll-interval")
	}
	if wait && events {
		return fmt.Errorf("--wait and --events cannot be combined")
	}
	if operation == "audio" {
		if len(args) > 0 {
			upload = args[0]
		}
		if f := cmd.Flag("file"); f != nil && f.Changed {
			upload = f.Value.String()
		}
	}
	if err := validateAdvancedDownload(cmd, operation, wait, events); err != nil {
		return err
	}
	req := app.OperationRequest{Operation: operation, PathParams: path, JSON: body, UploadPath: upload}
	if language, _ := cmd.Flags().GetString("accept-language"); language != "" {
		req.AcceptLanguage = language
	}
	ctx := cmd.Context()
	if wait {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	client := advancedClientFactory()
	if events {
		return client.StreamOperation(ctx, req, func(event app.Event) error { return json.NewEncoder(cmd.OutOrStdout()).Encode(event) })
	}
	result, err := client.ExecuteOperation(ctx, req)
	if err != nil {
		return err
	}
	var submission json.RawMessage
	if wait {
		submission = append(json.RawMessage(nil), result...)
		result, err = waitAdvancedResult(ctx, cmd, client, operation, path, body, submission, interval)
		if err != nil {
			return err
		}
	}
	if !json.Valid(result) {
		return fmt.Errorf("API response is not valid JSON")
	}
	if err := downloadAdvancedResult(ctx, cmd, client, result); err != nil {
		return err
	}
	if wait {
		if flag := cmd.Flag("json"); flag != nil && flag.Value.String() == "true" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(struct {
				Submission json.RawMessage `json:"submission"`
				Result     json.RawMessage `json:"result"`
			}{submission, result})
		}
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), advancedResultText(cmd, operation, result))
	return err
}

func waitAdvancedResult(ctx context.Context, cmd *cobra.Command, client advancedAPI, operation string, path map[string]string, body, submission json.RawMessage, interval time.Duration) (json.RawMessage, error) {
	result := submission
	var err error
	var task struct {
		ID      string `json:"id"`
		AsyncID string `json:"async_id"`
		AgentID string `json:"agent_id"`
	}
	if err = json.Unmarshal(result, &task); err != nil {
		return nil, fmt.Errorf("decode asynchronous task: %w", err)
	}
	if operation == "task_get" {
		task.ID = path["id"]
	}
	if operation == "agent_run" || operation == "agent_result" {
		var request struct {
			AgentID string `json:"agent_id"`
			AsyncID string `json:"async_id"`
		}
		_ = json.Unmarshal(body, &request)
		if task.AgentID == "" {
			task.AgentID = request.AgentID
		}
		if task.AsyncID == "" {
			task.AsyncID = request.AsyncID
		}
		if task.AgentID == "" || task.AsyncID == "" {
			return nil, fmt.Errorf("agent waiting requires agent_id and async_id")
		}
		result, err = client.WaitAgent(ctx, task.AgentID, task.AsyncID, interval)
	} else {
		if task.ID == "" {
			return nil, fmt.Errorf("asynchronous response has no task id")
		}
		result, err = client.WaitTask(ctx, task.ID, interval)
	}
	if err != nil {
		if _, outputErr := fmt.Fprintln(cmd.ErrOrStderr(), string(submission)); outputErr != nil {
			return nil, errors.Join(err, fmt.Errorf("write submission: %w", outputErr))
		}
		return nil, err
	}
	return result, nil
}
