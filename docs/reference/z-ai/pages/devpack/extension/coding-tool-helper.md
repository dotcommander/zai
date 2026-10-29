> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Coding Tool Helper

> [!TIP]
>
> **One-Click Install Assistant** \\
>
> A command-line assistant that helps GLM Coding Plan users centrally manage and configure CLI tools such as Claude Code.

**NPM Package**: [@z\_ai/coding-helper](https://www.npmjs.com/package/@z_ai/coding-helper) \
**Prerequisite**: [Node.js >= v18.0.0](https://nodejs.org/en/download/)

## Tool Overview

Coding Tool Helper is a coding-tool companion that quickly loads **GLM Coding Plan** into your favorite **Coding Tools**. Install and run it, then follow the on-screen guidance to automatically install tools, configure plan, and manage MCP servers.

The current coding tools supported are:

* **Claude Code**
* **Codex**
* **OpenCode**
* **Crush**
* **Factory Droid**

## Key Features

> **Interactive Wizard**
>
> Friendly setup guidance

> **Plan Integration**
>
> Connect GLM Plan to your preferred coding tools

> **Tool Management**
>
> Automatically detect, install, and configure coding tools

> **MCP Configuration**
>
> Easily manage MCP services

> **Local Storage**
>
> Secure local storage configuration

> **I18n Support**
>
> Interfaces support multiple languages

## Quick Start

**1. Get Your API Key**

Get your API Key from [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list) (Individual Plan) or [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) (Team Plan).

**2. Install & Launch**

Prerequisite: You need [Node.js 18+ or newer](https://nodejs.org/en/download/) \
Choose either installation method below.

**Method 1 (Recommended: npx on demand)**

Best for occasional users—no global install required. Run via npx to start instantly.

```shell theme={null}
## Run Coding Tool Helper directly in the terminal
npx @z_ai/coding-helper
```

**Method 2 (Global install for frequent use)**

Ideal for heavy users. Install globally, then launch with `coding-helper` or `chelper`.

> [!TIP]
>
> If `npm install` fails with `permission denied`, add sudo (macOS/Linux) or run the terminal as administrator (Windows).\
> Example: `sudo npm install -g @z_ai/coding-helper` \
> Alternatively, simply use npx: `npx @z_ai/coding-helper`

```shell theme={null}
## Install @z_ai/coding-helper globally
npm install -g @z_ai/coding-helper
## Then run coding-helper or chelper
coding-helper
```

**3. Complete the wizard**

Inside the wizard, use the arrow keys to choose options and Enter to confirm. Follow the guide to: \
Select UI language --> Choose a coding plan --> Enter API key --> Pick tools to manage \
\--> Auto-install tools (if needed) --> Open the tool management menu --> Load plan into tools \
\--> Manage MCP services (optional) --> Finish setup and launch your coding tools

## Additional Information

### Command List

> Beyond the interactive wizard, Coding Tool Helper also supports running specific commands via `coding-helper` or `chelper` with arguments:

```bash theme={null}
# Launch the initialization wizard
coding-helper init

# Language management
coding-helper lang show              # Display the current language
coding-helper lang set en_US         # Switch to English
coding-helper lang --help            # View language command help

# API key management
coding-helper auth                   # Configure the key interactively
coding-helper auth glm_coding_plan_global <token>     # Select the Global plan and set the key directly
coding-helper auth revoke            # Remove the stored key
coding-helper auth reload claude     # Load the latest plan into Claude Code
coding-helper auth --help            # View auth command help

coding-helper doctor                 # Check system configuration and tool status
coding-helper --help                 # Show help information
coding-helper --version              # Show version
```

### Troubleshooting

If issues arise, run `coding-helper doctor` first for a health check.

> **Network error, please check your connection**
>
> **Issue:** When saving or validating the API KEY or performing other network operations, you may see network errors such as `Network Error`.
>
> **Solution:**
>
> 1. Check your network connection or configure a proxy.
> 2. Note: If you must use a proxy to access external networks, Node.js does not automatically use the system proxy settings. Set the environment variables `HTTP_PROXY` and `HTTPS_PROXY` to make Node.js use the proxy.
>
> ```shell theme={null}
> # Example:
> export HTTP_PROXY=http://your.proxy.server:port
> export HTTPS_PROXY=http://your.proxy.server:port
> ```

> **Network timeout**
>
> **Issue:** When running or installing the coding tools, timeout or other network timeout errors appear.
>
> **Solution:**
>
> 1. Check the network connection or configure a proxy.

> **Insufficient permissions EACCES: permission denied**
>
> **Issue:** `npm install -g` throws EACCES: permission denied.
>
> **Solution:**
>
> 1. Retry with sudo (macOS / Linux).
> 2. Run the terminal as administrator (Windows).
> 3. Start directly via `npx @z_ai/coding-helper`.
> 4. Use nvm to manage Node.js versions and avoid global permission issues.

> **Incorrect plugin status in the Claude Code Marketplace**
>
> **Issue:** While using the Claude Code Marketplace, the plugin status is incorrect (e.g., it shows “not installed” even though it is already installed).
>
> **Solution:**
>
> 1. Run `claude update` to upgrade Claude Code to version 2.0.70 or later.

> **API Key invalid**
>
> **Issue:** API Key reported as invalid.
>
> **Solution:**
>
> 1. Confirm the API Key was copied correctly.
> 2. Check that the associated account has sufficient balance.

> **Connection timeout**
>
> **Issue:** Service connection timed out.
>
> **Solution:**
>
> 1. Check network connectivity.
> 2. Verify firewall settings.
> 3. Ensure Node.js and the network environment are ready.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.