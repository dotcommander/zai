> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Codex

> Methods for Using the GLM Coding Plan in Codex CLI

Codex is an AI coding agent from OpenAI that helps you write, review, and debug code. It is available both as a terminal tool and a desktop application.

> [!WARNING]
>
> Codex requires a dedicated OpenAI Responses protocol endpoint: `https://api.z.ai/api/v1`.

## Step 1: Install Codex

Codex is available in two forms: the **Codex CLI** that runs in your terminal, and the desktop **Codex app**. Choose whichever fits your workflow.

**Install the Codex App**

1. Visit the [ChatGPT download page](https://www.chatgpt.com/download) and download the installer for your operating system.
2. Follow the system prompts to complete the installation and launch the app.
3. If you already have Codex CLI installed, you can also launch the desktop app directly from the terminal:
   ```
   codex app
   ```

**Install Codex CLI**

Prerequisites:

* [Node.js 18 or newer](https://nodejs.org/en/download/)
* macOS users are recommended to install Node.js via [nvm](https://nodejs.org/en/download/). Direct package installation is not recommended (it may cause permission issues later).

Open a terminal and install Codex CLI:

```
npm install -g @openai/codex
```

Run the following command to verify the installation. A version number indicates success:

```
codex --version
```

## Step 2: Configure GLM Coding Plan

**1. Register an Account**

Visit the [Z.AI Open Platform](https://z.ai/model-api), click the "Register/Login" button in the top right, and follow the prompts to complete the registration process.

**2. Get Your API Key**

* Individual Plan users (including trial card recipients) can create an API Key under [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list).
* Team Plan members can create an API Key under [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan).
* API Keys for Individual Plans, Team Plans, and the standard model API are managed separately and are not interchangeable. For details, see the [API Key Management Guide](https://docs.z.ai/devpack/notice/api-key-management.md).

> [!WARNING]
>
> Keep your API Key safe. Do not share it with others or hard-code it directly in your code.

**3. Configure the Model and API Key**

Codex connects via the OpenAI Responses protocol, using the endpoint `https://api.z.ai/api/v1`.

> [!TIP]
>
> Choose whichever configuration method below suits you.

**Manual Configuration**

Supported on macOS & Linux & Windows. Note that config file paths differ by system. Make sure the file format is correct after editing.

#### 1. Create the Model Catalog File

Create `~/.codex/models.json` to declare GLM model metadata to Codex. If the file or directory does not exist, create it first.

Copy the following content entirely into `~/.codex/models.json`:

```json theme={null}
{
    "models": [
        {
            "slug": "glm-5.3",
            "display_name": "glm-5.3",
            "description": "Z.ai's latest flagship model",
            "default_reasoning_level": "max",
            "supported_reasoning_levels": [
                {
                    "effort": "low",
                    "description": "Light reasoning"
                },
                {
                    "effort": "high",
                    "description": "Enhanced reasoning"
                },
                {
                    "effort": "max",
                    "description": "Deep reasoning"
                }
            ],
            "shell_type": "shell_command",
            "visibility": "list",
            "supported_in_api": true,
            "priority": 0,
            "base_instructions": "",
            "supports_reasoning_summaries": true,
            "default_reasoning_summary": "none",
            "support_verbosity": false,
            "apply_patch_tool_type": "freeform",
            "truncation_policy": {
                "mode": "bytes",
                "limit": 10000
            },
            "context_window": 1048576,
            "max_context_window": 1048576,
            "effective_context_window_percent": 95,
            "supports_parallel_tool_calls": true,
            "experimental_supported_tools": [],
            "input_modalities": [
                "text"
            ]
        }
    ]
}
```

#### 2. Configure the API Key and Model

Locate the Codex configuration file at the following path. Create it if it does not exist:

* macOS/Linux: `~/.codex/config.toml`
* Windows: `C:\Users\<YourUsername>\.codex\config.toml`

Copy the following content entirely into `~/.codex/config.toml`, replacing `<Your API Key>` with the API Key you obtained in the previous step:

```toml theme={null}
model_provider = "ZAI"
model = "glm-5.3"
model_reasoning_effort = "max"
model_catalog_json = "~/.codex/models.json"

[model_providers.ZAI]
name = "ZAI"
base_url = "https://api.z.ai/api/v1"
experimental_bearer_token = "<Your API Key>"
wire_api = "responses"
```

> [!NOTE]
>
> * `wire_api` must be set to `responses`.
> * `model_catalog_json` points to the model metadata file created in the previous step.
> * `model` can be switched to `glm-5.3` or any other model declared in `models.json` as needed.

**Automated Coding Tool Helper**

Coding Tool Helper is a coding-tool companion that quickly loads **GLM Coding Plan** into your Codex. Install and run it, then follow the on-screen guidance to automatically install tools, configure plan.

```bash theme={null}
# Run Coding Tool Helper directly in the terminal
npx @z_ai/coding-helper
```

For more details, please refer to the [Coding Tool Helper](https://docs.z.ai/devpack/extension/coding-tool-helper.md) documentation.

> [!NOTE]
>
> After a successful configuration, restart the app for the changes to take effect.

## Step 3: Start Using Codex

Once configured, restart your `Codex app` or open a new terminal and run the `codex` command to get started.

## Step 4: FAQ

### Manual Configuration Not Taking Effect

* Close all Codex windows, open a new terminal window, and run `codex` again to start.
* Confirm that `~/.codex/models.json` exists and is valid JSON. You can use an online JSON validator to check.
* Confirm that the `model_catalog_json` path in `~/.codex/config.toml` is correct.
* Confirm that `experimental_bearer_token` has been replaced with your actual API Key.
* Confirm that the TOML format of the configuration file is correct—check field names, quotes, and brackets for completeness.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.