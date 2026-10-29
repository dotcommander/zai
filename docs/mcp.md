# Coding Plan MCP tools

`zai mcp` connects to the four MCP products documented by Z.AI. This is separate
from server-side MCP Calling inside a chat request: `zai mcp … call` invokes a
tool explicitly and never automatically executes a model's tool calls.

```sh
zai mcp search tools
zai mcp reader tools --json
zai mcp zread tools
zai mcp vision tools
zai mcp search call TOOL_NAME --request-json arguments.json
printf '%s' '{"url":"https://example.com"}' | zai mcp reader call TOOL_NAME --request-json -
```

Use the exact tool name and `inputSchema` returned by `tools`; the deployed tool
set can change. `--request-json` contains the tool's **arguments object**, not a
JSON-RPC envelope. Omission supplies `{}`. Unknown fields and large JSON numbers
are retained. Output is always the complete JSON result, with input schemas,
structured content, image/resource blocks, pagination, and metadata intact.
A tool result with `isError: true` is printed and exits unsuccessfully.

`tools --cursor CURSOR` sends the opaque `nextCursor` from the previous result.
`--timeout 2m` controls the whole initialization and call; it must be positive.
Requests are never automatically retried, including calls whose outcome is
unknown after a disconnect. A server that needs sampling or other client-side
capabilities is not given those capabilities, and no local functions run.

## Authentication and remote servers

The commands reuse the configured `api.key` / `ZAI_API_KEY`. Search, Reader, and
Zread use bearer authentication and MCP Streamable HTTP at the documented
endpoints:

| Product | Endpoint |
|---|---|
| Search | `https://api.z.ai/api/mcp/web_search_prime/mcp` |
| Reader | `https://api.z.ai/api/mcp/web_reader/mcp` |
| Zread | `https://api.z.ai/api/mcp/zread/mcp` |

Use `--server-url URL` to select a proxy or test endpoint. API chat/Coding Plan
base URL settings do not change these product endpoints. HTTP JSON and
request-scoped SSE responses are supported, with MCP initialization and session
headers. This adapter implements the initialized-session protocol used by the
2024/2025 MCP revisions; it is not a general implementation of every later MCP
revision. Unsupported negotiated revisions fail explicitly. Legacy endpoints
that put the API key in a URL query are not used.

## Local Vision server

Vision uses the documented stdio server command:

```sh
npx -y @z_ai/mcp-server
```

It needs Node.js 22 or newer and `npx` on PATH. Invoking a Vision command may
cause `npx` to download the upstream package; adding this CLI support does not
install or start it. The CLI supplies `Z_AI_API_KEY` from its configured key and
sets `Z_AI_MODE=ZAI`. Server stdout is reserved for newline-delimited JSON-RPC.
The owned subprocess and pipes are closed on completion or cancellation.
On Unix, cleanup also terminates the owned process group; other platforms
terminate the direct child only.

For a preinstalled/custom executable, use an argument vector rather than a shell
command:

```sh
zai mcp vision tools --command /path/to/server --command-arg=--some-option
```

`--command-arg` can be repeated and requires `--command`. It replaces the default
`npx` arguments. These flags apply only to Vision; `--server-url` applies only to
the remote products. No npm package was installed during local verification.

## Evidence and provider limits

Local HTTP/SSE and subprocess fixtures verify protocol framing, authentication,
raw results, errors, deadlines, and cleanup. They do not prove account entitlement,
quota, deployed tool schemas, npm package behavior, or live-provider availability.
The archived provider docs do not identify every deployed protocol revision.

Official setup guides: [Vision](https://docs.z.ai/devpack/mcp/vision-mcp-server),
[Search](https://docs.z.ai/devpack/mcp/search-mcp-server),
[Reader](https://docs.z.ai/devpack/mcp/reader-mcp-server),
[Zread](https://docs.z.ai/devpack/mcp/zread-mcp-server).
