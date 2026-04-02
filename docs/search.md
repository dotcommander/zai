Zai's search command queries Z.AI's web search API and displays results directly in your terminal. You get titles, domains, URLs, and optional content previews -- formatted for quick scanning or structured as JSON for scripts.

- [Running a Search](#running-a-search)
- [Filtering Results](#filtering-results)
    - [Result Count](#result-count)
    - [Recency](#recency)
    - [Domain Restriction](#domain-restriction)
    - [Language](#language)
- [Search Engines](#search-engines)
- [Content Size](#content-size)
- [Source Details](#source-details)
- [Output Formats](#output-formats)
- [Searching in Chat](#searching-in-chat)
- [Search-Augmented Chat](#search-augmented-chat)

<a name="running-a-search"></a>
## Running a Search

```bash
zai search "golang error handling best practices"
```

Zai displays a table of results with titles, domains, and URLs. You may also pipe the query from stdin:

```bash
echo "kubernetes networking" | zai search
```

<a name="filtering-results"></a>
## Filtering Results

<a name="result-count"></a>
### Result Count

Limit the number of results with `-c`:

```bash
zai search "AI news" -c 5
```

<a name="recency"></a>
### Recency

Filter by time window with `-r`:

```bash
zai search "golang release notes" -r oneMonth
```

| Value | Time Window |
|-------|-------------|
| `oneDay` | Last 24 hours |
| `oneWeek` | Last 7 days |
| `oneMonth` | Last 30 days |
| `oneYear` | Last 365 days |
| `noLimit` | No restriction (default) |

<a name="domain-restriction"></a>
### Domain Restriction

Restrict results to a specific domain with `-d`:

```bash
zai search "kubernetes tutorials" -d kubernetes.io
```

<a name="language"></a>
### Language

Set the result language with `-l`:

```bash
zai search "machine learning" -l zh
```

<a name="search-engines"></a>
## Search Engines

Z.AI offers multiple search engines. Choose one with `--engine`:

| Engine | Description |
|--------|-------------|
| `search_std` | Standard search (default) |
| `search_pro` | Enhanced search with richer results |
| `search_pro_sogou` | Pro search via Sogou |
| `search_pro_quark` | Pro search via Quark |

```bash
zai search "quantum computing" --engine search_pro
```

<a name="content-size"></a>
## Content Size

Control how much content each result includes with `--content-size`:

| Size | Detail Level |
|------|--------------|
| `medium` | 400-600 characters per result (default) |
| `high` | Up to 2,500 characters per result |

```bash
zai search "distributed systems" --content-size high
```

<a name="source-details"></a>
## Source Details

Include detailed source information with `--sources`:

```bash
zai search "WebAssembly" --sources
```

Zai appends additional metadata about each result's origin when this flag is active.

<a name="output-formats"></a>
## Output Formats

Control how results are displayed with `-o`:

```bash
zai search "docker compose" -o table
zai search "docker compose" -o detailed
zai search "docker compose" -o json
```

The `table` format (default) gives a compact overview. The `detailed` format includes content previews from each result. The `json` format outputs a structured object with query, duration, count, and a results array -- suitable for scripting.

The `--json` global flag also selects JSON output:

```bash
zai search "docker compose" --json
```

<a name="searching-in-chat"></a>
## Searching in Chat

Inside the REPL, use the `search` command to run a search and add results to conversation context:

```bash
zai chat
```

Then at the prompt, type `search "latest golang release" -c 5 -r oneWeek`. The model can reference those results in subsequent answers. For more on the REPL, see the [chat documentation](/docs/chat#interactive-repl).

<a name="search-augmented-chat"></a>
## Search-Augmented Chat

The `--search` flag on chat prompts is a different mechanism from the `search` command. It silently searches the web and injects results as context for the model -- you never see the raw results.

```bash
zai --search "What happened at GopherCon 2026?"
```

Use `--search` when you want the model to answer using current information. Use the `search` command when you want to browse results yourself. For a detailed comparison, see the [chat documentation](/docs/chat#search-flag-vs-search-command).

Search results are cached locally by default. Identical queries within the TTL window are served from cache without an API call. For cache configuration, see the [configuration documentation](/docs/configuration#web-search-settings).
