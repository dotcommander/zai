> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Tool Integration

## Supported Tools

> [!TIP]
>
> The GLM Coding Plan is limited to use within the following officially supported tools and product environments; users may not use their subscription benefits for tools or scenarios outside of this scope.

### 1. Coding Agent Tool

Click on the tool documentation below that you wish to use, and follow the corresponding integration guide to set it up.

> **[ZCode (1.5× Usage)](https://docs.z.ai/devpack/tool/zcode)**
>
> ZCode integrates AI agents into your existing toolchain.

> **[Claude Code](https://docs.z.ai/devpack/tool/claude)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Claude for IDE](https://docs.z.ai/devpack/tool/claude-for-ide)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Codex](https://docs.z.ai/devpack/tool/codex)**
>
> OpenAI's AI coding agent that helps you write, review, and debug code.

> **[OpenCode](https://docs.z.ai/devpack/tool/opencode)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Pi](https://docs.z.ai/devpack/tool/pi)**
>
> A minimalist terminal coding agent with a small core, extensible via extensions, Skills, and Pi Packages.

> **[Cursor](https://docs.z.ai/devpack/tool/cursor)**
>
> An AI-first code editor that supports custom model configurations.

> **[Cline](https://docs.z.ai/devpack/tool/cline)**
>
> An AI programming extension for VS Code that supports code generation and file operations.

> **[TRAE](https://docs.z.ai/devpack/tool/trae)**
>
> An AI editor capable of independently completing various development tasks.

> **[Qoder](https://docs.z.ai/devpack/tool/qoder)**
>
> An agentic coding platform designed for real software development.

> **[Droid](https://docs.z.ai/devpack/tool/droid)**
>
> Enterprise-grade AI coding agent that runs in the terminal to handle end-to-end workflows.

> **[Kilo Code](https://docs.z.ai/devpack/tool/kilo)**
>
> A powerful VS Code extension for code generation and project management.

> **[Roo Code](https://docs.z.ai/devpack/tool/roo)**
>
> A smart VS Code extension for code writing and refactoring.

> **[Crush](https://docs.z.ai/devpack/tool/crush)**
>
> A terminal-based AI programming tool that supports both CLI and TUI interfaces.

> **[Goose](https://docs.z.ai/devpack/tool/goose)**
>
> AI Agent tool, supporting local execution and automated engineering tasks.

> **[Eigent](https://docs.z.ai/devpack/tool/eigent)**
>
> A desktop AI agent built on a multi-agent architecture, capable of automating browser, terminal, and MCP-powered workflows.

### 2. General-purpose Agent Tool

> [!TIP]
>
> The general-purpose agent tools listed below are also supported and will continue to be served on a best-effort basis. Under high inference load, some requests may face temporary rate limits.

> **[AutoClaw (1.5× Usage)](https://docs.z.ai/devpack/tool/autoclaw)**
>
> Deeply optimized for GLM models, integrating specialized domain knowledge bases and workflows.

> **[OpenClaw](https://docs.z.ai/devpack/tool/openclaw)**
>
> An open-source AI assistant that runs on local devices, supports multi-platform use.

> **[Hermes Agent](#)**
>
> An open-source evolving AI agent with persistent memory,getting smarter with use.

> **[SillyTavern](#)**
>
> A highly customizable AI chat frontend for immersive roleplay with multi-model & media support.

## Coding Endpoint

GLM Coding Plan supports both the Anthropic and OpenAI protocols. Make sure to configure the correct `Base URL`:

| Protocol | Base URL |
| - | - |
| Anthropic Messages | `https://api.z.ai/api/anthropic` |
| OpenAI Chat Completions | `https://api.z.ai/api/coding/paas/v4` |
| OpenAI Responses | `https://api.z.ai/api/v1` |

**Core Steps**

1. Choose the protocol that fits your tool (Anthropic Messages or OpenAI Chat Completions).
2. Configure the correct Base URL.
3. Enter your API Key and select a GLM model.

> [!WARNING]
>
> Please select the correct endpoint address based on the tool you are using. Incorrect endpoint configuration will result in inability to use GLM Coding Plan subscription quota.

## Config Example

Using **Cline** as an example, the following steps demonstrate how to integrate GLM models via the OpenAI Compatible protocol. Similarly, other tools supporting the OpenAI-compatible protocol can adopt the same configuration approach.

### 1. Install the Cline Plugin

1. Open VS Code and click the Extensions Marketplace icon on the left.
2. Enter `cline` in the search box and locate the `Cline` extension.
3. Click `Install` to install it, then choose to trust the developer.

### 2. Configure API Endpoint

In Cline, select `Use your own API Key`, then fill in the following configuration:

* **API Provider**: Select `OpenAI Compatible`
* **Base URL**: Enter `https://api.z.ai/api/coding/paas/v4`
* **API Key**: Enter your Z.AI API Key
* **Model**: Select "Use custom" and enter the GLM model code you want to use (e.g., `glm-5.2`)
* **Other Configurations**:
  * Uncheck **Support Images**
  * Adjust **Context Window Size** based on your model (`glm-5.2` is `1000000`; other models `200000`)

### 3. Get Started

Once configured, you can enter your requirements in the input box to let the model assist you with code generation, file editing, refactoring, explaining code logic, debugging, and more.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.