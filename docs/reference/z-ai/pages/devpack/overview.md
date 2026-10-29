> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Overview

The GLM Coding Plan is your AI partner for tasks and code — a single subscription that unlocks a suite of agent tools for AI-powered coding and everyday office work.

## Usage

The plan can be applied to agent tools such as ZCode, Claude Code, Codex, OpenCode, AutoClaw, and OpenClaw, covering both coding and office scenarios:

### Coding

> **Natural Language Programming**
>
> Describe requirements in plain language to automatically generate plans, write code, debug issues, and ensure smooth execution.

> **Intelligent Code Completion**
>
> Get real-time, context-aware completion suggestions that reduce manual typing and significantly improve productivity.

> **Code Debugging & Repair**
>
> Input error messages or descriptions to automatically analyze your codebase, locate problems, and provide fixes.

> **Codebase Q\&A**
>
> Ask questions about your team’s codebase anytime, maintain global understanding, and receive precise answers with external data integration.

> **Automated Task Handling**
>
> Automatically fix lint issues, resolve merge conflicts, and generate release notes—allowing developers to stay focused on core logic.

### Office Work

> **Deep Search & Research**
>
> Cross-source intelligent retrieval and information synthesis, delivering structured research reports fast—so every decision is backed by evidence.

> **Business Data Analysis**
>
> Upload multi-dimensional business data to automatically generate visual charts, along with in-depth business insights and operational recommendations.

> **Creative Ideation & Brainstorming**
>
> Break through mental boundaries and spark inspiration from multiple angles, turning vague ideas into actionable plans — fast.

> **Automated Task Execution**
>
> Autonomously run complex workflows end to end, from instruction to delivery — cutting repetitive work and freeing up time for what matters.

## Advantages

* **Access to High-Intelligence Models:** Upon release, the GLM series achieved SOTA performance among open-source models in reasoning, coding, and agent capabilities, delivering outstanding results in tool use and complex task execution.
* **Works with 20+ Agent Tools:** Beyond Claude Code, it also supports ZCode, Cline, OpenCode, AutoClaw, and other [supported tools](https://docs.z.ai/devpack/tool/others.md#step-1-supported-tools), giving you flexibility across both coding and office workflows.
* **Generous Usage at a Fair Price:** Get higher call limits than standard plans. Starting at just 18 USD per month, with Pro and Max plans designed for high-frequency, complex projects.
* **Expanded Capabilities:** All plans support Vision Understanding, Web Search MCP， Web Reader MCP and Zread MCP helping you tackle a wider range of coding and office tasks.

## Benefits

### Supported Models

* All plans support **GLM-5.3**, GLM-5.3-Flash.
* Requests for GLM-5.2/GLM-5.1 will be automatically routed to GLM-5.3, requests for GLM-4.7 will automatically be routed to GLM-5.3-Flash.

### Usage Instruction

> [!WARNING]
>
> For information on Team Plan usage limits, please visit [Team Plan Benefits](https://docs.z.ai/devpack/teamplan.md).

#### Usage Credit Allowance

Each plan is subject to both a 5-hour usage limit and a weekly usage limit.

| Plan Type | 5-Hour Credits | Weekly Credits |
| :-: | :-: | :-: |
| Lite | 2,000 | 10,000 |
| Pro | 12,000 | 60,000 |
| Max | 28,000 | 140,000 |

**Credit Reset Rules**

* **5-hour credits**: Dynamically refreshed; credit quota resets 5 hours after consumption.
* **Weekly credits**: Activated upon subscription; resets every 7 days.

#### Credit Calculation

* Model credit usage = (Input tokens × Input multiplier + Cached Input tokens × Cached Input multiplier + Output tokens × Output multiplier) / 10,000
* MCP tool credit usage = Number of calls × Output multiplier

> [!TIP]
>
> You can view the number of tokens consumed under each pricing type and the number of tool calls on the [Charge Type](https://z.ai/manage-apikey/billing) page.

<table>
  <thead>
    <tr>
      <th>Product Type</th>
      <th>Product</th>
      <th>Input Multiplier</th>
      <th>Cached Input Multiplier</th>
      <th>Output Multiplier</th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td rowSpan="2">Model</td>
      <td>GLM-5.3</td>
      <td>6.9</td>
      <td>1.7</td>
      <td>24</td>
    </tr>

    <tr>
      <td>GLM-5.3-Flash\
      (Including MCP for visual understanding)</td>
      <td>2.3</td>
      <td>0.56</td>
      <td>8</td>
    </tr>

    <tr>
      <td rowSpan="3">MCP Server</td>
      <td>Web Search</td>
      <td>—</td>
      <td>—</td>
      <td>1.2</td>
    </tr>

    <tr>
      <td>Web Reader</td>
      <td>—</td>
      <td>—</td>
      <td>1.2</td>
    </tr>

    <tr>
      <td>Zread</td>
      <td>—</td>
      <td>—</td>
      <td>1.2</td>
    </tr>
  </tbody>
</table>

**During off-peak hours, model usage is charged at 50% of the standard credit rate**.

> [!NOTE]
>
> **Peak hours**: Monday to Friday, 14:00–18:00 Singapore Standard Time (UTC+8).

#### Estimated Token Allowance

Token usage varies depending on the cache hit rate, as shown below:

| Cache Hit Rate | Model | Lite (M Tokens/week) | Pro (M Tokens/week) | Max (M Tokens/week) |
| :-: | :-: | :-: | :-: | :-: |
| 95% | GLM‑5.3 | 48–97 | 290–580 | 676–1,352 |
| 95% | GLM‑5.3‑Flash | 146～292 | 877–1,755 | 2,047–4,095 |
| 96% | GLM‑5.3 | 50–99 | 297–595 | 694–1,387 |
| 96% | GLM‑5.3‑Flash | 150–300 | 900–1,801 | 2,101–4,202 |
| 98% | GLM‑5.3 | 52–104 | 313–627 | 731–1,463 |
| 98% | GLM‑5.3‑Flash | 158–317 | 950–1,900 | 2,217–4,433 |

**How the Range Is Calculated**

* Maximum token allowance: All usage occurs during off-peak hours and is charged at 0.5× the standard credit rate.
* Minimum token allowance: All usage occurs during peak hours and is charged at 1× the standard credit rate.

**By fully utilizing the off-peak discounts, you can save up to 92% compared with pay-as-you-go calls to the GLM-5.3 standard API**

### Exclusive MCP Access

> **[Vision Understanding](https://docs.z.ai/devpack/mcp/vision-mcp-server.md)**

> **[Web Search](https://docs.z.ai/devpack/mcp/search-mcp-server.md)**

> **[Web Reader](https://docs.z.ai/devpack/mcp/reader-mcp-server.md)**

> **[Zread](https://docs.z.ai/devpack/mcp/zread-mcp-server.md)**

## Next Steps

> **[Quick Start](https://docs.z.ai/devpack/quick-start.md)**
>
> Get up and running in minutes — from subscribing to the plan to using it in your agent tools.

> **[Usage Policy](https://docs.z.ai/devpack/usage-policy.md)**
>
> Learn about account usage rules, rate limits, refund policies, and other important guidelines.

> **[FAQ](https://docs.z.ai/devpack/faq.md)**
>
> Find answers to common questions about subscriptions, promotions, and using the plan.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.