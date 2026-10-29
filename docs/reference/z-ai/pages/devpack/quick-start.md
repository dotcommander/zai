> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Quick Start

> [!TIP]
>
> This guide will help you get started with [GLM Coding Plan](https://z.ai/subscribe?utm_source=zai\&utm_medium=link\&utm_term=quickstart\&utm_campaign=Platform_Ops&_channel_track_key=DRUfXN42) in minutes—from subscribing to using GLM models in officially [supported tools and products](https://docs.z.ai/devpack/tool/others.md#step-1-supported-tools).

## Getting Started

**1. Register or Login**

* Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.

**2. Subscribe to GLM Coding Plan**

After logging in, navigate to the [GLM Coding Plan](https://z.ai/subscribe?utm_source=zai\&utm_medium=link\&utm_term=quickstart\&utm_campaign=Platform_Ops&_channel_track_key=DRUfXN42) to select your preferred subscription plan.

**3. Obtain API Key**

After subscribing,

Individual Plan users (including trial card recipients) can create an API Key under [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list).

Team Plan members can create an API Key under [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan).

> [!TIP]
>
> API Keys for Individual Plans, Team Plans, and the standard model API are now managed separately. Newly created API Keys for these three product types are not interchangeable—use the one that matches your needs. For details, see the [API Key Management Guide](https://docs.z.ai/devpack/notice/api-key-management.md).

> [!WARNING]
>
> Safeguard your API Key by keeping it confidential and avoiding hard-coding it in your code.

**4. Connect a Tool**

The GLM Coding Plan is strictly limited to use within officially [supported tools and products](https://docs.z.ai/devpack/tool/others.md#step-1-supported-tools). Click a tool below to open its configuration guide:

> **[Claude Code](https://docs.z.ai/devpack/tool/claude)**

> **[Roo Code](https://docs.z.ai/devpack/tool/roo)**

> **[Kilo Code](https://docs.z.ai/devpack/tool/kilo)**

> **[Cline](https://docs.z.ai/devpack/tool/cline)**

> **[OpenCode](https://docs.z.ai/devpack/tool/opencode)**

> **[OpenClaw](https://docs.z.ai/devpack/tool/openclaw)**

> **[Crush](https://docs.z.ai/devpack/tool/crush)**

> **[Goose](https://docs.z.ai/devpack/tool/goose)**

> **[Cursor](https://docs.z.ai/devpack/tool/cursor)**

> **[Other Tools](https://docs.z.ai/devpack/tool/others.md)**

**5. Endpoint Guide**

GLM Coding Plan supports both the Anthropic and OpenAI protocols. Make sure to configure the correct `Base URL`:

| Protocol | Base URL |
| - | - |
| Anthropic Messages | `https://api.z.ai/api/anthropic` |
| OpenAI Chat Completions | `https://api.z.ai/api/coding/paas/v4` |
| OpenAI Responses | `https://api.z.ai/api/v1` |

**6. Start Coding**

Once configured, you can begin coding with the GLM model!

**Conversational Programming**

```bash theme={null}
# Using natural language commands in Claude Code
Please create a React component containing a user login form
```

**Code Debugging**

```bash theme={null}
# Describe the issue encountered
My API request returns a 404 error. Please help me check the code.
```

**Code Optimization**

```bash theme={null}
# Code optimization
This function performs poorly. Please optimize it for me.
```

## Advanced Features

> **Vision MCP Server (Coding Plan Exclusive)**
>
> All users can utilize the Vision MCP Server, which employs the flagship vision reasoning model to comprehend and analyze image content.
>
> * Analyze UI design mockups and generate corresponding code
> * Understand flowcharts and architecture diagrams
> * Extract text and information from screenshots
>
> For detailed usage instructions, refer to the [Vision MCP Server](https://docs.z.ai/devpack/mcp/vision-mcp-server.md) documentation.

> **Web Search MCP Server (Coding Plan Exclusive)**
>
> All users can utilize the Web Search MCP Server to access the latest technical information.
>
> * Search for the latest technical documentation and API changes
> * Obtain the latest information on open-source projects
> * Find solutions and best practices
>
> For detailed usage instructions, refer to the [Web Search MCP Server](https://docs.z.ai/devpack/mcp/search-mcp-server.md) documentation.

> **Web Reader MCP Server (Coding Plan Exclusive)**
>
> All users can utilize the Web Reader MCP Server to fetch full webpage content and extract structured data.
>
> * Fetch complete webpage content including text, and links
> * Extract structured data such as title, body, and metadata
> * Remote HTTP-based MCP service, no local installation required
>
> For detailed usage instructions, refer to the [Web Reader MCP Server](https://docs.z.ai/devpack/mcp/reader-mcp-server.md) documentation.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.