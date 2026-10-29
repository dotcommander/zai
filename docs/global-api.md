# Global API coverage

The [coverage manifest](api-coverage.json) inventories all 14 operations in the
[official global OpenAPI](https://docs.z.ai/openapi.json), retrieved on
2026-10-28. The [unmodified snapshot](openapi.json) preserves the complete source
schemas, constraints, examples, parameters, security schemes, and reference
targets. The manifest records request and response media types, original schema
references, and nested field paths, including array items, maps, and composition
branches. Its SHA-256 identifies the snapshot used for this inventory.

This is local transport coverage. It does not establish live endpoint, model,
account, region, or beta-agent availability. Arbitrary JSON fields are preserved;
that does not establish that a provider accepts undocumented fields.

## Published operations

Paths below are relative to the official server `https://api.z.ai/api`.

| Method | Published path | CLI |
|---|---|---|
| POST | `/paas/v4/chat/completions` | `chat`, `vision` |
| POST | `/paas/v4/videos/generations` | `video` |
| GET | `/paas/v4/async-result/{id}` | `task get ID` |
| POST | `/paas/v4/images/generations` | `image` |
| POST | `/paas/v4/async/images/generations` | `image --async` or `image async` |
| POST | `/paas/v4/audio/transcriptions` | `audio` |
| POST | `/paas/v4/tokenizer` | `tokenize` |
| POST | `/paas/v4/layout_parsing` | `layout` |
| POST | `/paas/v4/web_search` | `search` |
| POST | `/paas/v4/reader` | `reader` |
| POST | `/paas/v4/files` | `files upload PATH` |
| POST | `/v1/agents` | `agent run` |
| POST | `/v1/agents/async-result` | `agent result` |
| POST | `/v1/agents/conversation` | `agent conversation` |

## Requests, responses, and events

`--request-json FILE` reads one JSON object; `--request-json -` reads it from
stdin. Nested objects, arrays, unknown fields, nulls, and JSON numbers are
preserved. Supported flags override corresponding fields only when explicitly
supplied, including explicit false/zero values. Flag defaults do not replace
values from JSON. A positional prompt replaces the corresponding prompt or
message field; omit it to preserve the JSON field unchanged.

Use `--json` for the complete provider JSON response, including fields outside
the convenience response types. Ordinary output extracts standard text, URLs,
or task identifiers; structured or tool-only responses fall back to complete
JSON. Explicit `--json` also enables complete responses on existing convenience
commands, including chat and vision with file inputs. Advanced requests do not
perform convenience prompt enhancement or summarization. Keep credentials out
of request files; authentication comes from the configured API key.

`--download PATH` saves the first response artifact through the bounded atomic
downloader. An explicitly supplied existing `--output` flag also selects this
behavior. Use `--max-download-bytes` to set its size limit. Submission operations
`image_async`, `videos`, and `agent_run` require `--wait` when downloading.
`--accept-language VALUE` overrides the request's Accept-Language header.

```sh
zai chat --request-json chat.json
zai vision --request-json vision.json
zai search --request-json search.json
zai reader --request-json reader.json
zai tokenize --request-json tokenizer.json
zai layout --request-json layout.json
```

`--events` requests streaming and writes one JSON object per SSE event.
A JSON body containing `"stream": true` also selects this event output. Each
object has `event` and `data` strings and an `id` when supplied. The data string
preserves the event payload rather than projecting it into a text-only delta.
Multiline SSE data is joined with newlines. An omitted event name uses the SSE default `message`.
Provider stream errors terminate the command with an error; cancellation closes
the response body. Use events only with operations that support streaming. For audio convenience
streaming, `--stream --json` without `--events` is rejected; select `--events` for
NDJSON. Beta-agent terminal SSE conventions are not fully documented by the
vendor. The client recognizes `[DONE]` and known finish events; an unexpected
EOF is an error, rather than evidence that the agent completed.

```sh
cat chat.json | zai chat --request-json - --events
zai agent run --request-json agent.json --events
```

The finite `api OPERATION` interface accepts registry names: `chat`, `videos`,
`task_get`, `images`, `image_async`, `audio`, `tokenize`, `layout`, `search`,
`reader`, `files_upload`, `agent_run`, `agent_result`, `agent_conversation`,
`responses`, and `messages`. It does not accept arbitrary endpoint URLs. Prefer
`task get ID` for path parameters and `files upload PATH` for local uploads.

## Chat tools and context caching

Use `--tool-stream` on chat or vision commands to request streamed tool-call
arguments, text, and reasoning. It sets both `tool_stream: true` and
`stream: true`, selects the advanced path, and emits the same lossless NDJSON
events as `--events`. Supply function definitions in the request JSON:

```sh
zai chat --request-json tools.json --tool-stream
```

The JSON form, `"stream": true, "tool_stream": true`, is also supported.
An explicit `--tool-stream=false` overrides only `tool_stream`; it does not
disable streaming requested in JSON. An unchanged flag leaves the JSON value
intact. Events preserve reasoning, content, incremental tool arguments, finish
reasons, and usage when the provider sends them. The CLI returns tool calls to
the caller; it does not execute local functions. Model support follows the
[provider's tool-streaming documentation](https://docs.z.ai/guides/capabilities/stream-tool).

[Server-side MCP Calling](https://docs.z.ai/guides/capabilities/mcp-call) is a
chat feature distinct from connecting the CLI to a Coding Plan MCP server.
Supply the documented `tools` entry in a request file, for example:

```json
{
  "model": "glm-5.3",
  "messages": [{"role": "user", "content": "Read the report"}],
  "tools": [{
    "type": "mcp",
    "mcp": {
      "server_label": "documents",
      "server_url": "https://your-mcp-server.example/mcp",
      "transport_type": "streamable-http",
      "allowed_tools": ["read_document"]
    }
  }],
  "tool_choice": "auto"
}
```

Run `zai chat --request-json mcp-chat.json --json` for the complete response,
or add `--events` for streaming. The CLI preserves `server_label`, `server_url`,
`transport_type`, `allowed_tools`, custom `headers`, and additional documented
fields. The provider connects to the supplied server and invokes its tools;
choose servers and permissions accordingly. Both documented `sse` and
`streamable-http` transport values are passed through unchanged. Actual server,
model, and account availability require provider validation.

[Context caching](https://docs.z.ai/guides/capabilities/cache) is automatic for
repeated context. There is no cache-enable CLI flag or separate caching
operation. Keep repeated system prompts and history consistent when reuse is
desired. `--json` preserves `usage.prompt_tokens_details.cached_tokens`; events
preserve usage payloads when supplied. The provider determines whether tokens
were cached; local transport tests do not establish cache hits or billing.

## Multipart uploads

Audio transcription and file upload send multipart bodies. Local file contents
are streamed rather than embedded in JSON or loaded in full. JSON supplies the
remaining API fields; the local file supplies the binary `file` part.

```sh
zai audio --file recording.wav --request-json transcription.json
# A positional upload path is also accepted in advanced audio mode:
zai audio --request-json transcription.json recording.wav
zai files upload document.pdf --purpose user_data
zai files upload document.pdf --request-json upload.json --purpose agent
```

For audio JSON, use `--file` for the local upload path. For file uploads, use the
positional path. An explicitly supplied `--purpose` overrides JSON. Upload
purpose, model support, formats, and provider limits remain subject to the
snapshot and provider validation; no live upload is implied by these examples.

Audio `--hotwords 'alpha,beta'` is encoded as the documented string array in
both convenience and advanced requests. Without that explicit flag, a
`"hotwords": ["alpha", "beta"]` array in JSON is retained. An explicit empty
`--hotwords ''` overrides it with an empty array. The comma-separated flag uses
the convenience parser's limit of 100 entries; JSON requests preserve the
supplied array for provider validation.

## Asynchronous workflows

Image submission can return immediately or wait for its result. Existing video
convenience mode continues to wait and download; advanced video mode returns
submission JSON unless `--wait` is selected. Task retrieval works for both image
and video tasks. Agent result requests use their own JSON fields and route.

```sh
zai image --async --request-json image.json
zai image async --request-json image.json --wait --wait-timeout 3m --poll-interval 2s
zai video --request-json video.json --wait --wait-timeout 5m --poll-interval 2s
zai task get TASK_ID
zai task get TASK_ID --wait --wait-timeout 3m --poll-interval 2s
zai agent run AGENT_ID --request-json agent.json
zai agent result ASYNC_ID --agent-id AGENT_ID --request-json agent-result.json
zai agent conversation CONVERSATION_ID --agent-id AGENT_ID --request-json conversation.json
```

`--wait` requires an explicit positive `--wait-timeout`. `--poll-interval` defaults
to `1s` and must be positive. Waiting respects the caller deadline and terminal
failure states; a submission ID alone is not a successful generation result.
Waiting is supported for `image_async`, `videos`, `task_get`, `agent_run`, and
`agent_result`. Agent statuses include `pending`, `success`, and `failed`.
Positional agent identifiers override their matching JSON fields. Conversation
requests default to `slides_glm_agent` when no agent ID is supplied. If polling
fails after submission, the submission JSON is retained on stderr.
With `--wait --json`, output is an object containing the complete `submission`
and `result` responses, so identifiers remain available after polling completes.
Agent runs default to `stream: false`; an explicit JSON value or `--events`
selects streaming. Ordinary output retains complete JSON when tool calls are
present, including responses that also contain text.

## Routing and compatibility

For Paas operations, a configured `BaseURL` containing the full deployment
prefix is retained, and the operation suffix is appended. For example,
`https://proxy.example/custom/paas/v4` routes chat to
`https://proxy.example/custom/paas/v4/chat/completions`. An origin-only BaseURL,
or one ending at `/api`, receives `/api/paas/v4`.

The v1 and Anthropic families derive their routes independently from the
BaseURL's origin: `/api/v1/...` and `/api/anthropic/...`. They do not append to a
Paas prefix. A custom `CodingBaseURL` is respected for chat when coding mode is
enabled; it does not redirect image, video, search, agent, or other families.

`responses` sends JSON to `/api/v1/responses`; `messages` sends JSON to
`/api/anthropic/v1/messages`. These protocol compatibility routes are outside the
14-operation OpenAPI snapshot. Both preserve arbitrary JSON and SSE events;
vendor-specific fields and live protocol compatibility remain unverified.
Messages includes the API key and Anthropic version headers.

```sh
zai responses --request-json responses.json
zai messages --request-json messages.json --events
```

The existing `tts`, `embed`, and `model` convenience commands are retained.
Their absence from this snapshot does **not** mean their endpoints are
unavailable. Consult their command-specific documentation and current provider
support before drawing availability conclusions.

## Coding Plan MCP products

Use `zai mcp vision|search|reader|zread tools` to discover schemas and
`zai mcp PRODUCT call TOOL --request-json arguments.json` to invoke a tool.
See [MCP commands](mcp.md) for transports, Node requirements, pagination, and
provider verification limits. [Capability coverage](capability-coverage.json)
tracks these prose-documented features alongside the OpenAPI manifest.
