# Residual audit revalidation

Revalidated against the current source on 2026-10-28. The earlier critical/high
repair plan is retained as historical evidence; this pass implements the supported
remaining findings under `.work/plans/global-api-coverage.md`. Existing uncommitted
critical/high repairs are preserved.

## Supported findings repaired

| Audit finding | Current outcome and evidence |
| --- | --- |
| C1-02 version initialization | `Execute` refreshes Cobra's version after `main` initializes build metadata (`cmd/root.go`). |
| C1-03 auth-free help | Root without arguments and the `model` group load configuration without demanding a key. Actual requests still validate credentials. |
| C1-04 system files | Missing values that look like paths return errors; ordinary prose remains a literal prompt. Interactive chat propagates resolution errors. |
| C1-05 stdin boundary | Read limit is maximum plus one; exactly 10 MiB succeeds and excess input fails (`TestResidualStdinBoundary`). |
| Model selection/listing | Root exposes `--model`; chat respects its session override and configured default. Configured model listings include chat; API listings avoid claiming every returned model is a chat model. |
| C2-03 stale events | Streaming tokens, completion and errors carry turn identities; canceled turns cannot mutate newer buffers. Search/web results also carry identities (`TestResidualStaleStream`). |
| C2-04 SIGINT recovery | Raw mode creates and stops a signal context for each operation, so cancellation cannot poison all later turns. |
| C2-05 raw command parity | Raw mode handles help, clear, history, context, exit, search and web commands locally. Successful search/web results enter conversation context; incomplete streams do not (`TestResidualRawCommandsAndFileContext`). |
| C2-06 retained file | Both raw and TTY chat retain the file on subsequent turns; raw two-turn HTTP fixture checks payloads. |
| C2-07 / C4-05 UTF-8 | Chat web previews and search output truncate by runes; tiny truncation limits are safe (`TestResidualMissingSystemAndUTF8`). |
| C2-09 terminal detection | Interactive rendering requires both stdin and stdout to be terminals. |
| C2-10 active model history | Completed streams carry the selected request model; stream and ordinary chat history use the active override (`TestResidualActiveModelHistory`). |
| C3-04 chat breaker | Chat and streaming request headers pass through the configured chat breaker (`TestResidualChatBreaker`). |
| C3-07 typed retry | HTTP retries use `APIError` status codes, network errors use typed errors, and caller cancellation is excluded. Error-body substrings cannot trigger retries (`TestResidualRetryAndBackoff`). |
| C3-08 overflow | Backoff doubles only while below a checked cap; large, zero and negative attempt inputs remain bounded. |
| C3-09 opener | Select by OS, use Windows `cmd /c start`, and wait/reap the helper process (`internal/app/platform.go`). |
| C3-11 cache identity | Hash structured query/options JSON, including engine, content size and all other options; delimiter collisions disappear (`TestResidualCacheOptionIdentity`). |
| C3-12 cache I/O | Writes use temporary files and atomic rename. Cleanup scans without the cache lock and checks unchanged contents under the lock before deletion (`TestResidualCacheConcurrentCleanup`). |
| C4-07 image save errors | `ProcessImageResult` returns the save failure, producing command failure instead of success. |
| C4-08 audio chunk options | Chunk workers reuse `buildTranscriptionOptions` and disable only streaming; model, language prompt, user ID and hotwords propagate. |
| R-02 cancelable web/search | TTY search and reader operations share the active operation cancellation and stale-result checks. |

## Findings already repaired before this pass

C4-01 audio ownership; C2-01 terminal stream delivery; C2-02 oversized history
records; C3-01 message ordering and slice ownership; C3-02/03 concurrent breaker
admission and error classification; C3-05 separate stream HTTP policy; C3-06
zero-burst normalization; C3-10 bounded atomic downloads; C4-02 URL-only video input;
C1-01/C4-03 embedding defaults; C4-04 SSE audio; C4-06 typed audio errors.
Their source and regression fixtures remain in place. See
`.work/plans/audit-critical-high.md` for their original acceptance receipts.

## Unsupported claims and residual limits

- R-01 is an architectural limitation, not a demonstrated correctness failure:
  context remains session-local with a fixed message window, without token-budget
  estimation or resume persistence. No such capabilities are claimed.
- R-03's claim that ignoring SSE comments is defective is unsupported: comments
  are not content. Convenience chat still consumes its established single-line
  JSON delta format. Multi-line/event-preserving SSE belongs to the advanced API
  interface; this pass does not claim lossless convenience-chat SSE coverage.
  Callers close streams, including error/cancel paths.
- R-04's file-permission warning is a policy suggestion. This pass does not change
  existing configuration permissions or introduce unsolicited warnings.
- R-05 remains a separate size-policy limitation: convenience vision input still
  reads local images into memory. No safe maximum was established by the finding;
  no arbitrary vision limit was introduced.
- Windows/Linux opener behavior is source-validated only on this macOS host.
  Raw SIGINT recovery and TTY behavior are source/fixture evidence, not a live
  interactive terminal acceptance claim.

## Checks

Passed focused tests: `GOWORK=off go test ./cmd ./internal/app -run
'TestResidual|TestRepairChat|TestIsRetryableError|TestCalculateBackoff|TestFileSearchCache|TestGenerateCacheKey' -count=1`.

Passed affected race tests: `GOWORK=off go test -race ./cmd ./internal/app
-run 'TestResidual|TestRepairChat|TestCircuitBreaker|TestStreamingHTTP|TestZeroBurst'
-count=1`. `git diff --check` passed. No full suite, installed binary, live provider,
paid API request, commit or publication is claimed here.
