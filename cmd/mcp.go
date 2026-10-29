package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
	"github.com/dotcommander/zai/internal/mcp"
)

type mcpSession interface {
	Call(context.Context, string, json.RawMessage) (json.RawMessage, error)
	Close() error
}

var mcpFactory = func(ctx context.Context, opts mcp.Options) (mcpSession, error) { return mcp.New(ctx, opts) }

var mcpEndpoints = map[string]string{
	"search": "https://api.z.ai/api/mcp/web_search_prime/mcp",
	"reader": "https://api.z.ai/api/mcp/web_reader/mcp",
	"zread":  "https://api.z.ai/api/mcp/zread/mcp",
}

func newMCPCommand() *cobra.Command {
	group := &cobra.Command{Use: "mcp", Short: "Discover and call Coding Plan MCP tools (Vision, Search, Reader, Zread)"}
	for _, product := range []string{"vision", "search", "reader", "zread"} {
		product := product
		server := &cobra.Command{Use: product, Short: "Use the " + product + " MCP server"}
		server.PersistentFlags().Duration("timeout", 2*time.Minute, "Deadline for server startup and the complete request")
		server.PersistentFlags().String("server-url", "", "Override the remote MCP endpoint (Search, Reader, Zread)")
		server.PersistentFlags().String("command", "", "Override Vision's local stdio server executable")
		server.PersistentFlags().StringArray("command-arg", nil, "Argument to the Vision executable; repeat for each argument")
		tools := &cobra.Command{Use: "tools", Short: "List available tools and their input schemas", Args: cobra.NoArgs}
		tools.Flags().String("cursor", "", "Opaque nextCursor from a previous tools response")
		tools.RunE = func(cmd *cobra.Command, args []string) error { return runMCP(cmd, product, "tools/list", "") }
		call := &cobra.Command{Use: "call tool", Short: "Call a named tool with JSON arguments", Args: cobra.ExactArgs(1)}
		call.Flags().String("request-json", "", "Tool argument JSON object from a file or - (stdin); defaults to {}")
		call.RunE = func(cmd *cobra.Command, args []string) error { return runMCP(cmd, product, "tools/call", args[0]) }
		server.AddCommand(tools, call)
		group.AddCommand(server)
	}
	return group
}

func mcpOptions(cmd *cobra.Command, product string) (mcp.Options, time.Duration, error) {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	if timeout <= 0 {
		return mcp.Options{}, 0, fmt.Errorf("--timeout must be positive")
	}
	endpoint, _ := cmd.Flags().GetString("server-url")
	command, _ := cmd.Flags().GetString("command")
	args, _ := cmd.Flags().GetStringArray("command-arg")
	key := viper.GetString("api.key")
	if key == "" {
		return mcp.Options{}, 0, fmt.Errorf("API key is required; configure api.key or ZAI_API_KEY")
	}
	opts := mcp.Options{APIKey: key}
	if product == "vision" {
		if endpoint != "" {
			return opts, 0, fmt.Errorf("vision uses stdio; --server-url is unavailable")
		}
		if command == "" {
			if cmd.Flags().Changed("command-arg") {
				return opts, 0, fmt.Errorf("--command-arg requires --command")
			}
			command, args = "npx", []string{"-y", "@z_ai/mcp-server"}
		}
		opts.Command, opts.Args = command, args
		for _, variable := range os.Environ() {
			if !strings.HasPrefix(variable, "Z_AI_API_KEY=") && !strings.HasPrefix(variable, "Z_AI_MODE=") {
				opts.Env = append(opts.Env, variable)
			}
		}
		opts.Env = append(opts.Env, "Z_AI_API_KEY="+key, "Z_AI_MODE=ZAI")
	} else {
		if command != "" || cmd.Flags().Changed("command-arg") {
			return opts, 0, fmt.Errorf("--command and --command-arg are only available for Vision")
		}
		if endpoint == "" {
			endpoint = mcpEndpoints[product]
		}
		if endpoint == "" {
			return opts, 0, fmt.Errorf("unknown MCP product %q", product)
		}
		opts.URL = endpoint
		// Reuse the app's streaming transport and shared limiter for initialize,
		// notifications, tool listing, and calls. Tool calls are never retried here.
		transport := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		cfg := buildClientConfig()
		client := app.NewClientWithDeps(cfg, app.NewLogger(cfg.Verbose), nil, &app.ClientDeps{HTTPClient: transport, StreamHTTPClient: transport})
		opts.HTTPClient = client.StreamHTTPClient()
	}
	return opts, timeout, nil
}

func runMCP(cmd *cobra.Command, product, method, tool string) (resultErr error) {
	var params json.RawMessage
	if method == "tools/call" {
		source, _ := cmd.Flags().GetString("request-json")
		arguments, err := readRequestObject(source, cmd.InOrStdin())
		if err != nil {
			return err
		}
		if strings.TrimSpace(tool) == "" {
			return fmt.Errorf("tool name is required")
		}
		params, err = json.Marshal(struct {
			Name      string                     `json:"name"`
			Arguments map[string]json.RawMessage `json:"arguments"`
		}{tool, arguments})
		if err != nil {
			return err
		}
	} else {
		cursor, _ := cmd.Flags().GetString("cursor")
		values := map[string]string{}
		if cursor != "" {
			values["cursor"] = cursor
		}
		params, _ = json.Marshal(values)
	}
	opts, timeout, err := mcpOptions(cmd, product)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer cancel()
	session, err := mcpFactory(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := session.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close MCP session: %w", err))
		}
	}()
	result, err := session.Call(ctx, method, params)
	if err != nil {
		return err
	}
	if _, err = app.JSONObject(result); err != nil {
		return fmt.Errorf("MCP result must be an object: %w", err)
	}
	// Tool schemas, structuredContent, images, resource links, pagination and
	// metadata are all useful: return the full JSON result in every output mode.
	if _, err = fmt.Fprintln(cmd.OutOrStdout(), string(result)); err != nil {
		return err
	}
	var status struct {
		IsError bool `json:"isError"`
	}
	_ = json.Unmarshal(result, &status)
	if status.IsError {
		return fmt.Errorf("MCP tool %q reported an error (result retained above)", tool)
	}
	return nil
}
