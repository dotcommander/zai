> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Claude Code

> Methods for Using the GLM Coding Plan in Claude Code

Claude Code is an agentic coding tool that lives in your terminal, understands your codebase, and helps you code faster by executing routine tasks, explaining complex code, and handling git workflows -- all through natural language commands.

> [!TIP]
>
> To use the latest GLM-5.3 model, see [How to Switch the Model in Use](https://docs.z.ai/devpack/latest-model.md#switching-models-in-claude-code).

## Step 1: Installing the Claude Code

**Recommended Installation Method**

Prerequisites:

* [Node.js 18 or newer](https://nodejs.org/en/download/)
* For MacOS, please use [nvm](https://nodejs.org/en/download/) to install Nodejs, if you directly install the package, maybe encounter permission issues
* For Windows, please additionally install [Git for Windows](https://git-scm.com/download/win)

```
# Install Claude Code
npm install -g @anthropic-ai/claude-code

# Navigate to your project
cd your-awesome-project

# Complete
claude
```

**Cursor Guided Installation Method**

If you are not familiar with npm but have Cursor, you can enter the command in Cursor, and Cursor will guide you through the installation of Claude Code.

```bash theme={null}
https://docs.anthropic.com/en/docs/claude-code/overview Help me install Claude Code
```

> [!NOTE]
>
> **Note**: If MacOS users encounter permission issues during installation, please use [nvm](https://nodejs.org/en/download/) to install Nodejs.

## Step 2: Config GLM Coding Plan

**1. Get API Key**

* Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
* Individual Plan users (including trial card recipients) can create an API Key under [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list); Team Plan members can create one under [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan).
* Copy your API Key for use.

**2. Configure Environment Variables**

Mapping between Claude Code internal model environment variables and GLM models, with the default configuration as follows:

* `ANTHROPIC_DEFAULT_OPUS_MODEL: GLM-5.3-Flash`
* `ANTHROPIC_DEFAULT_SONNET_MODEL: GLM-5.3-Flash`
* `ANTHROPIC_DEFAULT_HAIKU_MODEL: GLM-5.3-Flash`

Set up environment variables **using one of the following methods** in the **macOS Linux** or **Windows**:

> [!TIP]
>
> **Note**: Some commands show no output when setting environment variables — that’s normal as long as no errors appear.

**Automated Coding Tool Helper**

Coding Tool Helper is a coding-tool companion that quickly loads **GLM Coding Plan** into your favorite **Coding Tools**. Install and run it, then follow the on-screen guidance to automatically install tools, configure plan, and manage MCP servers.

```bash theme={null}
# Run Coding Tool Helper directly in the terminal
npx @z_ai/coding-helper
```

For more details, please refer to the [Coding Tool Helper](https://docs.z.ai/devpack/extension/coding-tool-helper.md) documentation.

**Automated Script**

Just run the following command in your terminal \
Attention only macOS Linux environment is supported, this method does not support Windows

```bash theme={null}
curl -O "https://cdn.bigmodel.cn/install/claude_code_zai_env.sh" && bash ./claude_code_zai_env.sh
```

The script will automatically modify `~/.claude/settings.json` to configure the following environment variables(You don't need to edit manually):

```json theme={null}
{
    "env": {
        "ANTHROPIC_AUTH_TOKEN": "your_zai_api_key",
        "ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
        "API_TIMEOUT_MS": "3000000"
    }
}
```

**Manual configuration**

If you have previously configured environment variables for Claude Code, you can manually configure them as follows. A new window is required for the changes to take effect.

```bash MacOS & Linux & Windows theme={null}
# Edit the Claude Code configuration file `~/.claude/settings.json`
# Add or modify the env fields ANTHROPIC_BASE_URL, ANTHROPIC_AUTH_TOKEN
# Note to replace `your_zai_api_key` with the API Key you obtained in the previous step

{
        "env": {
        "ANTHROPIC_AUTH_TOKEN": "your_zai_api_key",
        "ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
        "ANTHROPIC_DEFAULT_HAIKU_MODEL": "glm-5.3-flash[1m]",
        "ANTHROPIC_DEFAULT_SONNET_MODEL": "glm-5.3[1m]",
        "ANTHROPIC_DEFAULT_OPUS_MODEL": "glm-5.3[1m]",
        "CLAUDE_CODE_AUTO_COMPACT_WINDOW": "1000000",
        "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1,
        "API_TIMEOUT_MS": "3000000"
    }
}
```

## Step 3: Start with Claude Code

Once the configuration is complete, you can start using **Claude Code** in your terminal or cmd:

```
cd your-project-directory
claude
```

> If prompted with "Do you want to use this API key," select "Yes."

After launching, grant Claude Code permission to access files in your folder as shown below:

![Description](https://cdn.bigmodel.cn/markdown/1753631613096claude-2.png?attname=claude-2.png)

You can use Claude Code for development Now!

***

## FAQ

### Manual Configuration Not Work

If you manually modified the `~/.claude/settings.json` configuration file but found the changes did not take effect, refer to the following troubleshooting steps.

* Close all Claude Code windows, open a new command-line window, and run `claude` again to start.
* If the issue persists, try deleting the `~/.claude/settings.json` file and then reconfigure the environment variables; Claude Code will automatically generate a new configuration file.
* Confirm that the JSON format of the configuration file is correct, check the variable names, and ensure there are no missing or extra commas; you can use an online JSON validator tool to check.

### Recommended Claude Code Version

We recommend using the latest version of Claude Code. You can check the current version and upgrade with the following commands:

> We have verified compatibility with Claude Code 2.0.14 and other versions.

```bash theme={null}
# Check the current version
claude --version

2.0.14 (Claude Code)

# Upgrade to the latest
claude update
```


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.