# Project Ledger
<!-- agent-ledger:v1 -->

Next-ID: L0006

## Current Focus

## Backlog

## Completed
- [x] L0005 Style human quota output like the main ZAI CLI | done: matching terminal design, plain pipe/JSON compatibility, focused tests, build and canonical install | source: Gary, 2026-10-28, “give this output a fresh new style design like the main 'zai' output!” | completed: 2026-10-28 | evidence: TestQuotaDisplayDefaultsMerge passed; all focused TestUsage tests passed after fixing low-percentage bars and preview spacing; project build passed; canonical symlink and installed usage --help verified. Offline fixture preview only; no new provider requests.
- [x] L0004 Verify live account quota access once | done: installed usage JSON command yields real quota windows; no retries, generation or account/config changes | source: Gary, 2026-10-28, “proceed all remaining” after residual live-access limitation | completed: 2026-10-28 | evidence: /Users/vampire/code/go/bin/zai usage --json succeeded; TIME_LIMIT 4%, 40/1000 used, 960 remaining; TOKENS_LIMIT 1%, 300-minute duration; both returned reset times. Live entitlement demonstrated for the configured account at this observation, not guaranteed indefinitely.
- [x] L0003 Verify and install quota-enabled ZAI | done: focused quota/usage tests, project build, task-owned diff review and canonical symlink/offline smoke | source: Gary, 2026-10-28, “write specs and implement”; .work/spec-account-quota-r1.md | completed: 2026-10-28 | evidence: focused app/config/cmd tests and build passed; review findings fixed and affected cmd tests/rebuild passed; installed symlink and offline help/JSON/bootstrap checks passed; FUI pst_01M62N0XJJGQ5T5QQ9EGSS0EZC; live entitlement unverified
- [x] L0002 Implement account quota visibility, tests, configuration and docs | done: complete task-owned slice per .work/spec-account-quota-r1.md | source: Gary, 2026-10-28, “write specs and implement” after quota comparison | completed: 2026-10-28 | evidence: new quota/usage source, tests and YAML asset; focused tests and build passed; review findings fixed and installation completed under L0003
- [x] L0001 Specify account quota visibility | done: decision-complete r1 scope and acceptance in .work/spec-account-quota-r1.md | source: same request as L0002 | completed: 2026-10-28 | evidence: .work/spec-account-quota-r1.md
