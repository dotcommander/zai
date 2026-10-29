# Agent Chat
Source: https://docs.z.ai/api-reference/agents/agent

POST /v1/agents
General Translation: General Translation API provides large model-based multilingual translation services, including general translation, paraphrase translation, two-step translation, and three-pass translation strategies. It supports automatic language detection, glossary customization, translation suggestions, and streaming output. Users only need to call the Translation API, input the text to be processed, specify the source language (auto-detection supported) and target language to receive high-quality translation results. 

 Popular Special Effects Videos: Popular special effects videos are intelligent templates launched based on trending features from pan-entertainment platforms, designed to precisely adapt to short video creative production needs. Currently, three effect templates are available: `French Kiss`, `Body Shake Dance`, and `Sexy Me`. After selecting a template, users only need to upload an image and enter corresponding prompts to generate a special effects video. 

 GLM Slide/Poster Agent: An intelligent creation agent built for working people and creators. It goes beyond traditional engineering-style assembly tools—supporting one-click generation of slides or posters from natural language instructions. By natively integrating content generation with layout aesthetics and design conventions, it helps you quickly produce polished, professional-grade materials while lowering design barriers and boosting creative efficiency.



# Conversation History
Source: https://docs.z.ai/api-reference/agents/agent-conversation

POST /v1/agents/conversation
This endpoint is used to query the agent conversation history.Only support slides_glm_agent



# File Upload
Source: https://docs.z.ai/api-reference/agents/file-upload

POST /paas/v4/files
This API is designed for uploading auxiliary files (such as glossaries, terminology lists) to support the translation service. It allows users to upload reference materials that can enhance translation accuracy and consistency.

## File Limitations

* Maximum 100 files
* Maximum 100MB per file
* Files retained for 180 days


# Retrieve Result
Source: https://docs.z.ai/api-reference/agents/get-async-result

POST /v1/agents/async-result
This endpoint is used to query the result of an asynchronous request.



# Errors
Source: https://docs.z.ai/api-reference/api-code



When calling the API, the response code consists of two parts: the outer layer is the HTTP status code, and the inner layer is the business error code defined by Z.AI in the response body, which provides a more detailed error description.

| Error Codes | ​HTTP Status Code | Error messege |
| :- | :- | :- |
| - | 500 | Internal Error |
| 1000 | 401 | Authentication Failed |
| 1001 | 401 | Authentication parameter not received in Header, unable to authenticate |
| 1003 | 401 | Authentication Token expired, please regenerate/obtain |
| 1005 | 401 | Need Two-Factor Authentication |
| 1113 | 429 | Insufficient balance or no resource package. Please recharge. |
| 1200 | 500 | API Call Error |
| 1210 | 400 | Invalid API parameter, please check the documentation. |
| 1211 | 400 | Unknown Model, please check the model code. |
| 1212 | 400 | Current model does not support `${method}` call method |
| 1213 | 400 | Parameter `${field}` was not received. |
| 1214 | 400 | Parameter `${field}` is invalid. Please check the documentation. |
| 1215 | 400 | `${field1}` and `${field2}` cannot be set simultaneously, please check the documentation |
| 1220 | 403 | You do not have permission to access \$ |
| 1221 | 400 | API `${API_name}` has been taken offline |
| 1222 | 400 | API `${API_name}` does not exist |
| 1230 | 500 | API call process error |
| 1234 | 500 | Network error, error id: `${error_id}`, please try again later |
| 1261 | 400 | Prompt too long |
| 1301 | 400 | System detected potentially unsafe or sensitive content in input or generation. Please avoid using prompts that may generate sensitive content. Thank you for your cooperation. |
| 1302 | 429 | Rate limit reached for requests |
| 1305 | 429 | The service may be temporarily overloaded, please try again later |
| 1308 | 429 | Usage limit reached for `{number}` `{unit}`. Your limit will reset at `{next_flush_time}` |
| 1309 | 429 | Your GLM Coding Plan package has expired and is temporarily unavailable. You can resume using it after renewing the subscription on the official website. [https://z.ai/subscribe。](https://z.ai/subscribe。) |
| 1310 | 429 | Weekly/Monthly Limit Exhausted. Your limit will reset at `{next_flush_time}` |
| 1311 | 429 | Your current subscription plan does not yet include access to `${model_name}` |
| 1313 | 429 | Your account's current usage pattern does not comply with the Fair Usage Policy, and your request frequency has been limited. For details, please refer to the Subscription Service Agreement. To restore access, please submit a request. |
| 1314 | 429 | Your enterprise package has expired. Please contact your enterprise administrator. |
| 1315 | 429 | This API Key is limited to enterprise coding package scenarios. Please go to the official website to replace the API Key of the corresponding product type. |
| 1316 | 429 | Usage limit reached for the past 5 hours. Insufficient balance for extra usage. Resets at `{next_flush_time}`. |
| 1317 | 429 | Usage limit reached for the past 7 days. Insufficient balance for extra usage. Resets at `{next_flush_time}`. |
| 1318 | 429 | Usage limit reached for the past 5 hours. Extra usage is not available due to monthly spend limit. Resets at `{next_flush_time}`. |
| 1319 | 429 | Usage limit reached for the past 7 days. Extra usage is not available due to monthly spend limit. Resets at `{next_flush_time}`. |
| 1320 | 429 | Usage limit reached for the past 5 hours. Extra usage is not available due to monthly spend limit. Resets at `{next_flush_time}`. |
| 1321 | 429 | Usage limit reached for the past 7 days. Extra usage is not available due to monthly spend limit. Resets at `{next_flush_time}`. |

## Error Shapes

Errors are always returned as JSON, with a top-level error object that includes a `code` and `message` value.

```json theme={null}
{
  "error": {
    "code": "1214",
    "message": "Parameter `${field}` is invalid. Please check the documentation."
  }
}
```

## Error Example

The following is the response message of a curl request, where 401 is the HTTP status code and 1001 is the business error code.

```
* We are completely uploaded and fine
< HTTP/2 401
< date: Wed, 20 Mar 2024 03:06:05 GMT
< content-type: application/json
< set-cookie: acw_tc=76b20****a0e42;path=/;HttpOnly;Max-Age=1800
< server: nginx/1.21.6
< vary: Origin
< vary: Access-Control-Request-Method
< vary: Access-Control-Request-Headers
<
* Connection #0 to host open.z.ai left intact
{"error":{"code":"1001","message":"Authentication parameter not received in Header, unable to authenticate"}}
```

> **Note**: When using streaming (SSE) calls, if the API terminates abnormally during inference, the above error codes will not be returned. Instead, the reason for the exception will be provided in the `finish_reason` parameter of the response body. For details, please refer to the description of the `finish_reason` parameter.


# Audio Transcriptions
Source: https://docs.z.ai/api-reference/audio/audio-transcriptions

POST /paas/v4/audio/transcriptions
Use the [GLM-ASR-2512](/guides/audio/glm-asr-2512) model to transcribe audio files into text, supporting multiple languages and real-time streaming transcription.



# Generate Image
Source: https://docs.z.ai/api-reference/image/generate-image

POST /paas/v4/images/generations
Use [GLM-Image](/guides/image/glm-image) series models to generate high-quality images from text prompts. Through quick and accurate understanding of user text descriptions, `AI` image expression becomes more precise and personalized.



# Generate Image(Async)
Source: https://docs.z.ai/api-reference/image/generate-image-async

POST /paas/v4/async/images/generations
Use the [GLM-Image](/guides/image/glm-image) series models to generate high-quality images from text prompts. Through quick and accurate understanding of user text descriptions, `AI` image expression becomes more precise and personalized. Only supports `GLM-Image` model.



# Retrieve Result
Source: https://docs.z.ai/api-reference/image/get-image-status

GET /paas/v4/async-result/{id}
This endpoint is used to query the result of an asynchronous request.



# Introduction
Source: https://docs.z.ai/api-reference/introduction



> [!NOTE]
>
> The API reference describes the RESTful APIs you can use to interact with the Z.AI platform.

Z.AI provides standard HTTP API interfaces that support multiple programming languages and development environments, with [SDKs](/guides/develop/python/introduction) also available.

## API Endpoint

Z.ai Platform's general API endpoint is as follows:

```
https://api.z.ai/api/paas/v4
```

> [!WARNING]
>
> When using GLM Coding Plan, please follow the [tutorial](/devpack/quick-start) to configure your dedicated endpoint.

## Authentication

The Z.AI API uses the standard **HTTP Bearer** for authentication.
A Standard API Key is required, which you can create or manage on the [API Keys Page](https://z.ai/manage-apikey/apikey-list).

API keys should be provided via HTTP Bearer Authentication in HTTP Request Headers.

```
Authorization: Bearer ZAI_API_KEY
```

## Playground

The API Playground allows developers to quickly try out API calls. Simply click **Try it** on the API details page to get started.

* On the API details page, there are many interactive options, such as **switching input types**, **switching tabs**, and **adding new content**.
* You can click **Add an item** or **Add new property** to add more properties the API need.
* **Note** that when switching the tabs, the previous properties value you need re-input or re-switch.

## API Usage Bundle Billing Rules

If you purchase an [API usage bundle](https://z.ai/usage-bundles) for a model, we will prioritize consuming the included usage quota from that bundle. Once the bundle is exhausted, your API access remains fully unaffected — you can continue making calls as usual with no service interruption whatsoever. The only change is that, going forward, each API call will be billed at the standard pay-as-you-go rate, with the corresponding charges deducted directly from your account.

## Call Examples

**cURL**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Accept-Language: en-US,en" \
-H "Authorization: Bearer YOUR_API_KEY" \
-d '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful AI assistant."
        },
        {
            "role": "user",
            "content": "Hello, please introduce yourself."
        }
    ],
    "temperature": 1.0,
    "stream": true
}'
```

**Official Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Usage Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key="YOUR_API_KEY")

# Create chat completion request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a helpful AI assistant."
        },
        {
            "role": "user",
            "content": "Hello, please introduce yourself."
        }
    ]
)

# Get response
print(response.choices[0].message.content)
```

**Official Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

**Usage Example**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import java.util.Arrays;

public class QuickStart {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();

        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Hello, who are you?")
                    .build()
            ))
            .stream(false)
            .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        // Get response
        System.out.println(response.getData().getChoices().get(0).getMessage().getContent());
    }
}
```

**OpenAI Python SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
pip install --upgrade 'openai>=1.0'
```

**Verify Installation**

```python theme={null}
python -c "import openai; print(openai.__version__)"
```

**Usage Example**

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)

completion = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "system", "content": "You are a smart and creative novelist"},
        {"role": "user", "content": "Please write a short fairy tale story as a fairy tale master"}
    ]
)

print(completion.choices[0].message.content)
```

**OpenAI NodeJs SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
npm install openai

# Or using yarn
yarn add openai
```

**Usage Example**

```javascript theme={null}
import OpenAI from "openai";

const client = new OpenAI({
    apiKey: "your-Z.AI-api-key",
    baseURL: "https://api.z.ai/api/paas/v4/"
});

async function main() {
    const completion = await client.chat.completions.create({
        model: "glm-5.3",
        messages: [
            { role: "system", content: "You are a helpful AI assistant." },
            { role: "user", content: "Hello, please introduce yourself." }
        ]
    });

    console.log(completion.choices[0].message.content);
}

main();
```

**OpenAI Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>com.openai</groupId>
    <artifactId>openai-java</artifactId>
    <version>2.20.1</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'com.openai:openai-java:2.20.1'
```

**Usage Example**

```java theme={null}
import com.openai.client.OpenAIClient;
import com.openai.client.okhttp.OpenAIOkHttpClient;
import com.openai.models.chat.completions.ChatCompletion;
import com.openai.models.chat.completions.ChatCompletionCreateParams;

public class QuickStart {
    public static void main(String[] args) {
        // Initialize client
        OpenAIClient client = OpenAIOkHttpClient.builder()
            .apiKey("your-Z.AI-api-key")
            .baseUrl("https://api.z.ai/api/paas/v4/")
            .build();

        // Create chat completion request
        ChatCompletionCreateParams params = ChatCompletionCreateParams.builder()
            .addSystemMessage("You are a helpful AI assistant.")
            .addUserMessage("Hello, please introduce yourself.")
            .model("glm-5.3")
            .build();

        // Send request and get response
        ChatCompletion chatCompletion = client.chat().completions().create(params);
        Object response = chatCompletion.choices().get(0).message().content();

        System.out.println(response);
    }
}
```


# Chat Completion
Source: https://docs.z.ai/api-reference/llm/chat-completion

POST /paas/v4/chat/completions
Create a chat completion model that generates AI replies for given conversation messages. It supports multimodal inputs (text, images, audio, video, file), offers configurable parameters (like temperature, max tokens, tool use), and supports both streaming and non-streaming output modes.



# Rate Limits
Source: https://docs.z.ai/api-reference/rate-limit





# Layout Parsing
Source: https://docs.z.ai/api-reference/tools/layout-parsing

POST /paas/v4/layout_parsing
Use the [GLM-OCR](/guides/vlm/glm-ocr) model to parse the layout of documents and images and extract text content. Support OCR recognition of images and PDF documents, returning detailed layout information and visualization results.



# Tokenizer
Source: https://docs.z.ai/api-reference/tools/tokenizer

POST /paas/v4/tokenizer
`Tokenizer` is used to split text into `tokens` recognizable by the model and calculate the count. It receives user input text, processes it through the model for tokenization, and finally returns the corresponding `token` count. It is suitable for text length evaluation, model input estimation, dialogue context truncation, cost calculation, etc.



# Web Reader
Source: https://docs.z.ai/api-reference/tools/web-reader

POST /paas/v4/reader
Reads and parses the content of the specified URL. Supports selectable return formats, cache control, image retention, and summary options.



# Web Search
Source: https://docs.z.ai/api-reference/tools/web-search

POST /paas/v4/web_search
The [Web Search](/guides/tools/web-search) is a specialized search engine for large language models. Building upon traditional search engine capabilities like web crawling and ranking, it enhances intent recognition to return results better suited for LLM processing (including webpage titles, URLs, summaries, site names, favicons etc.).



# Generate Video(Async)
Source: https://docs.z.ai/api-reference/video/generate-video

POST /paas/v4/videos/generations
CogVideoX is a video generation large model developed by Z.AI, equipped with powerful video generation capabilities. Simply inputting text or images allows for effortless video creation.



# Retrieve Result
Source: https://docs.z.ai/api-reference/video/get-video-status

GET /paas/v4/async-result/{id}
This endpoint is used to query the result of an asynchronous request.



# Invite Friends, Get Credits
Source: https://docs.z.ai/devpack/credit-campaign-rules



> [!TIP]
>
> Please read and fully understand all the rules below before participating in the campaign. Participation constitutes your acceptance and agreement to be bound by these rules.

Last Update: March 15, 2026

## 1. General Principles

* **Campaign Organizer**: This campaign is initiated and operated by the service provider of Z.ai Platform (hereinafter referred to as "we," "us," or the "Platform").
* **Campaign Purpose**: To reward users and enhance product experience, we are launching the "Invite Friends, Get Credits" campaign.
* **Campaign Nature**: This campaign offers incentive Credits. All Credits are preferential benefits provided by the Platform, **not cash, currency, or financial instruments**, and do not possess monetary attributes.
* **User Agreement**: By participating, you represent that you have **read, understood, and agreed** to be bound by these rules. We reserve the right to modify or adjust these rules based on laws, regulations, and operational needs. Changes will be announced on the campaign page in advance.

Event link: >>[Click here to join](https://z.ai/subscribe?invitedialog=true)

## 2. Benefits for Invited Users

* **Invitation Validity Period**:
  * Friends enter via referral link, log in on the event page to bind. Inviter gets bonus if friend completes first subscription payment within 72hrs. If friend clicks multiple links, reward goes to the last valid touchpoint inviter.
* **Instant Discount**:
  * Invited users who complete registration via a unique invitation link or code and **successfully pay for their first subscription to 【GLM Coding】 service exclusively** will enjoy a **10% instant discount** on the order amount (the discount amount will be directly deducted from the total amount payable).
  * **Important Payment Note**: Due to requirements of our payment processor (Stripe), the final amount payable **after applying all discounts and Credits must be at least 0.50 USD**. If the discount reduces the payable amount below 0.50 USD, the discount applied may be adjusted to ensure the minimum payment threshold is met.
* **Eligible Users**:
  * Only **newly registered users** of the Platform, or **existing registered users who have never had any paid subscription history**.
  * Each user (determined by mobile number or email address) can only enjoy the first-order discount **once as an invited user**.
* **Exclusions**:
  * This discount **applies only to the current initial subscription order**.
  * **Subsequent renewals, upgrades, downgrades, or any other follow-up subscription orders are NOT eligible** for this discount.
  * This discount cannot be stacked with other similar first-order discount campaigns on the Platform.

## 3. Benefits for Inviters

Inviting users (hereinafter "Inviter") can earn the following Credits rewards after successfully inviting friends (i.e., "Invited Users") who complete registration and their first valid payment:

* **Bonus Eligibility Threshold (Payout begins once 3 people are reached)**:
  * The referrer must **successfully invite a total of 3 eligible friends** who have completed their first package purchase before receiving the bonus reward.
  * **Once the 3-person threshold is met, the referral bonuses for the first 3 valid friends will be disbursed in a single lump-sum payment**; thereafter, the corresponding bonus for each new valid friend will be disbursed immediately.
  * While the 3-person threshold has not yet been met, any accrued referral bonuses will remain in a **pending** status. They can be viewed in the referrer’s account but cannot be used until the 3-person threshold is met, at which point they will be automatically unlocked and disbursed.
* **Base Reward**:
  * For each successfully invited valid friend, the Inviter receives a Credits equivalent to **10% of the friend's first order Actual Payment Amount** (defined below).
  * There is **no upper limit** on the number of invites; rewards accumulate continuously.
* **Tiered Extra Reward**:
  * For **every cumulative 30** successfully invited valid friends, the Inviter receives an additional one-time Credits equivalent to **10% of the total Actual Payment Amount** from those 30 friends.
  * Extra rewards are calculated per "**round**" (each set of 30 invites). Rounds can repeat and accumulate without an upper limit.
* **Definition of "Actual Payment Amount"**:
  * Refers to **the final cash amount actually paid** by the Invited User for their first subscription order. This amount specifically includes payments made via cash balance in the account or third-party payment methods.
  * This amount **excludes** any reductions from coupons, discount codes, credits, points, cash rebates, or any other promotional benefits, **regardless of** any subsequent refunds.
* **Credits Nature & Usage**:
  * **Nature**: Credits are **non-cash** preferential benefits provided by the Platform.
  * **Usage**: Credits can only be used to offset the following consumption on Z.ai Platform:
  * Purchasing or renewing GLM Coding and other subscription services.
  * Purchasing virtual products like resource packs or feature expansion packs within the platform.
  * Paying for platform API call fees.
  * **Restrictions**: Credits are **non-refundable, non-transferable, non-giftable**, and **cannot be used to request refunds for subscription services**.
  * **Validity Period**: Credits are **long-term valid** unless the user account is deactivated or the Platform ceases operation. We reserve the right to adjust the Credits validity period based on business needs in the future, with 30 days' advance notice provided.

## 4. Important Notes & Risk Warnings

* **Validity Criteria & Single Account Principle**:
  * The **sole criterion** for a valid invitation is: The friend registers via the Inviter's **unique invitation link or code** and completes the first subscription payment without refund within the campaign period.
  * **Single Account Principle**: Each natural person may only use **one platform account** to participate. The Platform determines user identity based on multiple factors including **phone number, email address, device ID (e.g., IMEI/OAID/IMSI), IP address, etc. The same phone number, device, payment account, or network environment may be deemed the same user**.
  * The Platform's backend system records are the **sole and final basis** for determining invitation relationships, validity, and user identity. In case of discrepancy with user records, the system records prevail.
  * If you have objections regarding invitation records or Credits issuance, please contact customer service within **15 days** of the objection arising for investigation. Objections raised after this period will be considered waived.
* **User Information Authorization & Privacy**:
  * To facilitate invitation tracking, Credits distribution, and anti-fraud checks, you understand and authorize the Platform to process necessary account and information (e.g., account ID, registration time, device info, payment info) for you and your friends within the scope of this campaign.
  * The Platform will strictly adhere to its privacy policy to protect all users' personal information. Invitation relationship information will not be disclosed to third parties without user consent.
* **Violations & Penalties**:
  * Users must comply with the Platform's User Agreement, Terms of Service, Subscription Agreement, Top-Up Agreement, campaign rules, etc. **Any form of cheating, fraud, or disruption is strictly prohibited**, including but not limited to:
  * **Technical Cheating**: Using plugins, scripts, emulators, virtual location software, hacker tools, automated scripts, bulk registration tools, etc., to fake invitations, registrations, payments, or tamper with device/network data.
  * **False Info & Account Trading**: Using false information for registration/verification, infringing on third-party rights, buying/selling platform accounts.
  * **Artificial Engagement**: Paying for services, joining booster groups, or exchanging non-genuine social engagements to complete tasks.
  * **Disruption**: Launching network attacks, exploiting loopholes, bypassing security measures, or any activity disrupting normal platform operations.
  * **Account Association**: If any account is found involved in fraud/cheating, we reserve the right to **disqualify all associated accounts** from the campaign and **reclaim all rewards** earned through those accounts based on reasonable judgment.
  * **Penalties**: If users are found violating these rules per our risk control systems or reasonable judgment, we may take one or more actions **without prior notice**:
  * **Deduct** all issued but unused Credits; **Reclaim the cash value** of any consumed Credits.
  * **Restrict, freeze, or permanently ban** violating accounts.
  * **Permanently prohibit** violating and associated accounts from all current and future Platform campaigns and benefits.
  * **Reserve the right to claim compensation** for all direct/indirect losses caused (e.g., extra operational costs, reputational damage, recovery costs) and **report to authorities**.
  * **Pursue legal liability** for severe cases involving significant amounts or illegal activities.
* **Campaign Changes, Suspension, & Termination**:
  * We reserve the right to **adjust, pause, or terminate** this campaign early based on actual operations, laws, regulations, or force majeure, announced via our website, in-app messages, or other reasonable means. Users may not hold the Platform liable for such actions.
* **Final Interpretation Right**:
  * **To the extent permitted by law**, Z.ai Platform reserves the **final right to interpret** all terms, details, and reward distributions related to this campaign. We have the right to make final judgments on unforeseen circumstances or ambiguities in a **fair, reasonable manner consistent with the campaign's purpose**.

***

**If you have any questions about these rules, please contact our customer service. Thank you for your participation!**


# Coding Tool Helper
Source: https://docs.z.ai/devpack/extension/coding-tool-helper



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


# FAQ
Source: https://docs.z.ai/devpack/faq



## GLM Coding Plan Details

**Q: Which models are supported by GLM Coding Plan? How can I switch between them?**

**A:** All plans support GLM-5.3, GLM-5.3-Flash.

> [!NOTE]
>
> For instructions on using the latest models, see [How to Switch Models](/devpack/latest-model).

***

**Q: How much usage quota does the plan provide?**

**A:** Subscribe once, unlock unmatched usage and unbeatable value.

For usage details, see:

* [Individual Plan Usage Quotas](https://docs.z.ai/devpack/overview#usage-instruction)
* [Team Plan Usage Quotas](https://docs.z.ai/devpack/teamplan#usage-details)

***

**Q: Apart from the 5-hour usage limit, are there any other restrictions on usage?**

**A:** To manage resources and ensure fair access for all users, we have introduced a weekly usage quota limit. This quota starts counting from the time you place your order, and the quota is refreshed and reset on a 7-day cycle. You can view the progress of your quota consumption in the [subscription](https://z.ai/manage-apikey/subscription).

***

**Q: What is a Quota Reset Card?**

**A:** A **Quota Reset Card** is a benefit occasionally issued to GLM Coding Plan users during special campaigns. Once used, it restores your consumed plan quota to the full amount, so you do not need to wait for the next scheduled reset.

The rules are as follows:

1. **Reset Types**: Quota Reset Cards are available in two types: **5-hour quota reset** and **weekly quota reset**. When you use a weekly quota reset card, your 5-hour quota will also be restored to full without consuming an additional 5-hour reset.

2. **Reset Cycle**: After a reset card is used, the corresponding quota reset cycle will also restart. The next reset time will be recalculated from the moment the card is used.

3. **Validity Period**: Each reset card has its own expiration date. Unused cards will automatically expire after the validity period ends.

***

**Q: Which tools are supported?**

**A:** The GLM Coding Plan is strictly limited to use within officially [supported tools and products](https://docs.z.ai/devpack/tool/others#step-1-supported-tools). The subscriber shall not use the subscription benefits in any unsupported tools or scenarios.\
Please refer to our tool guide for step-by-step setup. All supported tools share the same usage quota under your subscription.

***

**Q: What happens when my plan quota runs out? Will the system consume my account balance?**

**A:** No. GLM calls made in supported tools will **only** use your Coding Plan quota.

* Once the quota is used up, you’ll need to wait until the next **5-hour cycle** for it to refresh. The system will not deduct from your account balance.
* Users subscribed to the Coding Plan can only make calls via the plan’s quota in supported tools. API calls outside the plan are not available.

***

## MCP Call

**Q：Which plan supports Vision Understanding, Web Search and Web Reader MCP tools?**

**A：** All plans support these MCP tools.

***

**Q：What is the quota of Vision Understanding, Web Search and Web Reader MCP in the plan?**

**A：** Usage for both models and MCP tools is calculated and deducted in credits. For details on your usage allowance, please see [Usage Instruction](https://docs.z.ai/devpack/overview#usage-instruction)

***

**Q：Besides using the GLM Coding Plan resource package, can I use other methods to call these three MCP tools for Vision Understanding, Web Search and Web Reader?**

**A：** Other than the resource package, we currently do not provide any other access solutions for calling these three MCP tools. If you use other similar MCP tools and incur billing issues during their use, such issues are not within the scope of this package.

***

## Subscription Management

**Q: How are subscription fees deducted?**

**A:** The system deducts fees in the following order:

1. Credits balance will be used first.
2. If insufficient, cash balance will be used;.
3. If still insufficient, the remaining amount will be charged from your linked payment method (e.g., bank card or PayPal).

> [!TIP]
>
> Please note that a small minimum applies when charging your credit card. If the remaining amount is less, we will round up the deduction to meet this minimum.

***

**Q: How can I cancel my subscription?**

**A:** You can cancel your subscription on the subscription management page. Please make sure to cancel at least 24 hours before your next billing date to avoid auto-renewal. After cancellation, your current plan remains valid until it expires.

***

**Q: Can I get a refund?**

**A:** Please note that subscriptions are non-refundable once purchased. Even if you haven’t used your full plan, the fees cannot be returned. We recommend choosing a subscription plan and billing cycle that best fits your usage needs.

***

## Upgrade your Plan

**Q: How can I upgrade my coding plan?**

**A:** For **same-tier plan upgrades** (e.g., monthly plan to annual plan), the newly purchased plan will **not take effect immediately**. It will only become active after the expiration of your current plan, meaning the validity periods of the plans are cumulative.

Example: Upgrading from a **Lite Monthly Plan** to a **Lite Annual Plan** will grant you a total of 13 months of service validity.

For **cross-tier plan upgrades**, the new plan will take effect immediately. The remaining value of your original plan will be converted into account balance on a pro-rated basis, which can be used to offset the price difference of the upgrade.

***

**Q: How can I change my billing cycle?**

**A:** On the subscription page, select your desired plan. The change will take effect after the current billing cycle ends. You can also cancel your current subscription and re-subscribe to the desired plan after the billing cycle ends. Steps are as follows:

1. Go to your [subscription page](https://z.ai/subscribe?utm_source=zai\&utm_medium=index\&utm_term=glm-coding-plan\&utm_campaign=Platform_Ops&_channel_track_key=6lShUDnv).
2. Click “Subscribe” and select your desired billing cycle.
3. Confirm the change.
4. The new cycle will take effect after the current billing period ends.

***

## 'Invite Friends, Get Credits'

**Q: Where can I view the credits I earned?**

**A:** You can view your credit rewards in the \[Invite Friends, Earn Credits pop-up - Your Credit Rewards] or check the details under \[Billing - Transaction History].

***

**Q: Where can I use my credits? Can I withdraw them?**

**A:** Credits can be used on Z.ai platform to offset various purchases, including subscriptions like GLM Coding Plan, resource packs, feature extensions, and API calls.

> [!TIP]
>
> Credits are promotional benefits and cannot be withdrawn, transferred, or refunded.

***

**Q: How is a successful invitation defined? What does my friend need to do?**

**A:** A successful invitation must meet all of the following:

1. The referred friend shall access the campaign page through the inviter's unique referral link or code.
2. Your friend is a new user (has never paid for a subscription).
3. They complete their first GLM Coding Plan subscription payment within 72 hours.
4. The order is not refunded within 24 hours.

***

**Q: When will I receive the credits after inviting a friend?**

**A:** Credits are usually issued within 24-48 hours after your friend’s payment is confirmed and the order passes review. Check \[Finance - Transaction History]. If not received after 72 hours, contact support.

***

**Q: What does "friend’s actual payment amount" mean? How are credits calculated if a coupon is used?**

**A:** "Actual payment amount" refers to the final cash amount your friend paid for their first order.

For example: If a plan costs 100 USD, and your friend gets a 50% first-order discount, an extra 10% off from your referral, and uses 5 USD in credits, the final actual payment is 40 USD. Your credits will be calculated based on 40 USD (e.g., 10% of 40 USD = 4 USD). Credit calculation does not include any coupons, discounts, point deductions, or refunded amounts.

***

**Q: Why didn’t I receive credits after referring a friend?**

**A:** Common reasons include:

1. Your friend didn’t register using your exclusive link/code.
2. Your friend completed payment after 72 hours.
3. Your friend was not a new subscribing user.
4. The order was refunded.
5. The system detected unusual activity (e.g., bulk registrations, fraud).

**If none of these apply, please contact customer support with your account and your friend's registration information.**

***

## Technical Support

***

**Q：Why does it still report error "1113 Insufficient Balance" after purchasing the coding package? Why is the account balance still deducted after purchasing the coding package?**

**A：** The situation of reporting insufficient balance or deducting account balance may be due to not meeting the usage conditions of the GLM Coding Plan coding package:

1. The GLM Coding Plan is strictly limited to use within officially [supported tools and products](https://docs.z.ai/devpack/tool/others#step-1-supported-tools).
2. A specific baseurl address must be configured to use it:

* API endpoint for Claude Code and Goose is：`https://api.z.ai/api/anthropic`
* API endpoint for other tools is: `https://api.z.ai/api/coding/paas/v4`

3. Only the following two models can be called: GLM-5.3, GLM-5.3-Flash. For detailed instructions, please refer to:
   [https://docs.z.ai/devpack/overview](https://docs.z.ai/devpack/overview)


# How to Switch Models
Source: https://docs.z.ai/devpack/latest-model



> [!TIP]
>
> **GLM-5.3 & GLM-5.3-Flash Now Live** \\
>
> The GLM Coding Plan now supports the latest GLM-5.3 and GLM-5.3-Flash models for all users (Max, Pro, and Lite), and you can switch between models within your preferred Coding Agent.

## Before You Start

This is a model-switching guide for existing GLM Coding Plan users—not an onboarding tutorial. Make sure you have already:

1. **Active GLM Coding Plan subscription** and a valid **Z.AI API key**.
2. **Configured the correct endpoint for your tool:**
   * **Claude Code / Goose** (Anthropic-compatible): `https://api.z.ai/api/anthropic`
   * **Codex**: `https://api.z.ai/api/v1`
   * **Other OpenAI-compatible tools**: `https://api.z.ai/api/coding/paas/v4`
3. **Verified that your tool can successfully call an existing GLM model** (for example, `glm-5.3` or `glm-5.3-flash`). If a basic call fails, please resolve that first.

## Switching Models in Claude Code

### Step 1: Update the default configuration

**macOS**

* Method 1: In the terminal, run `vim ~/.claude/settings.json` to open and edit the file. When you're done, press Esc, type `:wq`, and save your changes.
* Method 2: In Finder, choose Go → Go to Folder, then enter `~/.claude/settings.json` to locate and edit the configuration file.

**Windows**

The user-level Claude Code settings file is located at:

`%USERPROFILE%\.claude\settings.json`

To create or open it in PowerShell:

```powershell theme={null}
New-Item -ItemType Directory -Force "$env:USERPROFILE\.claude"
notepad "$env:USERPROFILE\.claude\settings.json"
```

> [!TIP]
>
> If you use Git Bash or WSL, the path resolution may differ. Make sure the file is read by the same Claude Code installation you launch.

**Linux / WSL**

The user-level Claude Code settings file is located at:

`~/.claude/settings.json`

To create or open it:

```bash theme={null}
mkdir -p ~/.claude
${EDITOR:-nano} ~/.claude/settings.json
```

> [!TIP]
>
> If you are using WSL alongside native Windows Claude Code, the configuration file locations may differ. Make sure you edit the file used by the Claude Code installation you actually launch.

To use GLM-5.3-FLASH, add or replace the following environment variables in `settings.json`. **Do not replace the entire file** if it already contains other fields —only add or update the keys shown below.

```json theme={null}
{
  "env": {
    "CLAUDE_CODE_AUTO_COMPACT_WINDOW": "1000000",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "glm-5.3-flash[1m]",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "glm-5.3-flash[1m]",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "glm-5.3-flash[1m]"
  }
}
```

#### Enable 1M Context

To enable GLM 1M context, add the `[1m]` suffix to the model name (e.g., `glm-5.3-flash[1m]`) and configure the compression window size parameter `"CLAUDE_CODE_AUTO_COMPACT_WINDOW": "1000000"`.

If Claude Code reports that the model with the `[1m]` suffix does not exist, please upgrade Claude Code to the latest version and try again.

#### Switch Effort (Thinking Intensity)

In a Claude Code session, type the `/effort` command to switch thinking intensity. The default level is `max`.

| Tool input value | Actual level | Handling |
| - | - | - |
| `thinking.type` not set, `true`, `enabled`, `adaptive` | max | Uses the default level |
| `thinking.type` is `false`, `disabled`, `none`, `off` | low | Continues the request; still performs lightweight thinking |
| `reasoning_effort` is `minimal`, `light`, `low` | low | Auto-converted |
| `reasoning_effort` is `medium`, `high` | high | Auto-converted |
| `reasoning_effort` is `xhigh`, `max`, `ultra` | max | Auto-converted |
| `reasoning_effort` is any other unknown string | max | Falls back to the default level and logs a hint |

**Processing priority**: Explicit Effort > thinking toggle > default `max`.

> [!NOTE]
>
> Claude Code uses `thinking.type` and `output_config.effort`; Codex uses `reasoning.effort`. Disabling the thinking configuration is converted to `low` and does not switch to another model.

> [!TIP]
>
> For coding tasks, we recommend switching to `max` effort for deeper reasoning and more stable complex task performance.

### Step 2: Confirm that the model has been switched

Open a terminal and run `claude`, then type `/status`. Verify:

1. **Settings source** shows your `~/.claude/settings.json`.
2. **Model** shows `glm-5.3-flash` or `glm-5.3-flash[1m]`.

![Description](https://cdn.bigmodel.cn/markdown/1785421261069image.png?attname=image.png)

## Switching Models in Other Tools

This currently only works with coding agents that support custom model configuration. If your agent tool does not allow custom model settings, you will need to wait for official support in a future release.

### Using Cline as an example:

Please use the following settings:

* API Provider：Select OpenAI Compatible
* Base URL：Enter [https://api.z.ai/api/coding/paas/v4](https://api.z.ai/api/coding/paas/v4)
* API Key：Enter your Z.AI API key
* Model：Choose Custom Model and enter the model name, such as `glm-5.3` or `glm-5.3-flash`
* Other settings:
  * GLM-5.3 is a text-only model, so uncheck **Support Images**; GLM-5.3-FLASH is a multimodal model, so **Support Images** can be checked
  * Set Context Window Size to 1000000
  * Adjust temperature and any other parameters based on your task requirements

![Description](https://cdn.bigmodel.cn/markdown/1785421409503image.png?attname=image.png)


# Web Reader MCP Server
Source: https://docs.z.ai/devpack/mcp/reader-mcp-server



The Web Reader MCP Server is an exclusive Remote MCP Server developed by Z.AI for GLM Coding Plan users. Built on the Model Context Protocol (MCP), it connects to webpage content extraction capabilities to provide webpage content extraction, detailed page reading, and structured data retrieval for MCP-compatible clients, including Claude Code and Cline.

## Features

> **Web Content Reading**
>
> Fetch the complete content of any webpage, including text, and links

> **Structured Data**
>
> Extract structured data such as title, main body, and metadata

> **Remote Service**
>
> HTTP-based remote MCP service, no local installation required

## Tools

This server implements the Model Context Protocol and works with any MCP-compatible client. Currently, it provides the following tool:

* **`webReader`** — Fetch webpage content for a specified URL. Returns the page title, main content, metadata, list of links, and more.

## Example Scenarios

> **API Documentation Reading and Summarization**
>
> Automatically fetch and parse titles, body content, examples, and release notes from official documentation pages, distilling key takeaways to accelerate integration and implementation.

> **Open Source Project Page Parsing**
>
> Parse project websites or repository pages (such as README, release notes, and usage guides) to extract core information and link lists, assisting evaluation and integration.

> **Technical Article Understanding and Knowledge Extraction**
>
> Extract steps, commands, and caveats from blogs, tutorials, and guide pages, organizing unstructured content into actionable developer notes and task lists.

> **Bug Resolution Using Reference Documentation**
>
> For issue remediation, read the publicly available steps on the specified web page and use them as references to resolve the problem.

> **Knowledge Base Construction and Synchronization**
>
> Convert content from designated web pages into structured data and leverage in-page links for incremental synchronization to build a team technical knowledge base.

## Installation and Usage

### Quick Start

**1. Get API Key**

Get your API Key from [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list) (Individual Plan) or [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) (Team Plan)

**2. Configure MCP Server**

According to the client you’re using, **choose the corresponding installation method from the options below**.

### Supported Clients

**Claude Code**

**One-click install command**

Replace `your_api_key` with the API key you obtained in the previous step

```bash theme={null}
claude mcp add -s user -t http web-reader https://api.z.ai/api/mcp/web_reader/mcp --header "Authorization: Bearer your_api_key"
```

**Manual configuration**

Edit the Claude Code configuration file under your home directory, the MCP section of `.claude.json`:

```json theme={null}
{
  "mcpServers": {
    "web-reader": {
      "type": "http",
      "url": "https://api.z.ai/api/mcp/web_reader/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

**Cline (VS Code)**

Add the MCP server configuration in the Cline extension settings:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
  "mcpServers": {
    "web-reader": {
      "type": "streamableHttp",
      "url": "https://api.z.ai/api/mcp/web_reader/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

If Cline older version does not support StreamableHttp type MCP server, you can use SSE type configuration:

```json theme={null}
{
  "mcpServers": {
    "web-reader": {
      "type": "sse",
      "url": "https://api.z.ai/api/mcp/web_reader/sse?Authorization=your_api_key"
    }
  }
}
```

**OpenCode**

Add the MCP server configuration in OpenCode settings:

See the [OpenCode MCP documentation](https://opencode.ai/docs/mcp-servers)

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
    "$schema": "https://opencode.ai/config.json",
    "mcp": {
        "web-reader": {
            "type": "remote",
            "url": "https://api.z.ai/api/mcp/web_reader/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Crush**

Add the MCP server configuration in Crush settings:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
    "$schema": "https://charm.land/crush.json",
    "mcp": {
        "web-reader": {
            "type": "http",
            "url": "https://api.z.ai/api/mcp/web_reader/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Goose**

Not support Goose now，refer [Issue](https://github.com/block/goose/issues/6576)

Add the MCP server in Goose:

Go to `Extensions` -> `Add custom extension`

Set Extension Name to `web-reader`, Type to `HTTP`, and use the following endpoint:

```
https://api.z.ai/api/mcp/web_reader/mcp
```

Set Request Headers Add `Authorization` : `your_api_key`

Click `Add Extension` at the bottom. Remember to replace `your_api_key` with the API key you obtained in the previous step.

**Roo Code, Kilo Code, Others**

For Roo Code, Kilo Code, and other MCP-compatible clients, use the following general configuration:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
  "mcpServers": {
    "web-reader": {
      "type": "streamable-http",
      "url": "https://api.z.ai/api/mcp/web_reader/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

## Troubleshooting

> **Invalid access token**
>
> **Issue:** Received an invalid access token error
>
> **Solutions:**
>
> 1. Verify the token was copied correctly
> 2. Check that the token is activated
> 3. Ensure the token has sufficient balance
> 4. Confirm the Authorization header format is correct

> **Connection timeout**
>
> **Issue:** Connection to the MCP server timed out
>
> **Solutions:**
>
> 1. Check your network connection
> 2. Verify firewall settings
> 3. Ensure the server URL is correct
> 4. Increase client timeout settings

> **Webpage fetch failed**
>
> **Issue:** Web content reading returned empty result or error
>
> **Solutions:**
>
> 1. Confirm the target URL is accessible
> 2. Check if the page has anti-scraping mechanisms
> 3. Try different URLs
> 4. Ensure network connectivity is normal
> 5. Contact technical support for assistance

## Resources

* [Model Context Protocol (MCP) Documentation](https://modelcontextprotocol.io/)
* [Claude Code MCP Configuration Guide](https://docs.anthropic.com/en/docs/claude-code/mcp)
* [MCP Usage Limits](/devpack/overview#usage-instruction)
* [GLM Coding Plan Overview](/devpack/overview)


# Web Search MCP Server
Source: https://docs.z.ai/devpack/mcp/search-mcp-server



The Web Search MCP Server is an exclusive Remote MCP Server developed by Z.AI for GLM Coding Plan users. Built on the Model Context Protocol (MCP), it connects to search capabilities to provide web search and real-time information retrieval for MCP-compatible clients, including Claude Code and Cline.

## Features

> **Web Search**
>
> Supports comprehensive web search to retrieve the latest web information and resources

> **Real-time Information**
>
> Retrieves real-time updated information including news, stock prices, weather, and more

> **Remote Service**
>
> HTTP protocol-based remote MCP service, no local installation required

## Supported Tools

This server implements the Model Context Protocol and can be used with any MCP-compatible client. Currently provides the following tools:

* **`webSearchPrime`** - Search web information, returning results including page titles, URLs, summaries, site names, site icons, and more.

## Installation and Usage

### Quick Start

**1. Get API Key**

Get your API Key from [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list) (Individual Plan) or [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) (Team Plan)

**2. Configure MCP Server**

According to the client you’re using, **choose the corresponding installation method from the options below**.

### Supported Clients

**Claude Code**

**One-click Installation Command**

Be sure to replace `your_api_key` with the API Key you obtained.

```bash theme={null}
claude mcp add -s user -t http web-search-prime https://api.z.ai/api/mcp/web_search_prime/mcp --header "Authorization: Bearer your_api_key"
```

**Manual Configuration**

Edit Claude Code's configuration file `.claude.json` in the user directory, MCP section:

```json theme={null}
{
  "mcpServers": {
    "web-search-prime": {
      "type": "http",
      "url": "https://api.z.ai/api/mcp/web_search_prime/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

**Cline (VS Code)**

Add MCP server configuration in Cline extension settings:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
  "mcpServers": {
    "web-search-prime": {
      "type": "streamableHttp",
      "url": "https://api.z.ai/api/mcp/web_search_prime/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

If Cline older version does not support StreamableHttp type MCP server, you can use SSE type configuration:

```json theme={null}
{
  "mcpServers": {
    "web-search-prime": {
      "type": "sse",
      "url": "https://api.z.ai/api/mcp/web_search_prime/sse?Authorization=your_api_key"
    }
  }
}
```

**OpenCode**

Add MCP server configuration in OpenCode settings:

Refer [OpenCode MCP Doc](https://opencode.ai/docs/mcp-servers)

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
    "$schema": "https://opencode.ai/config.json",
    "mcp": {
        "web-search-prime": {
            "type": "remote",
            "url": "https://api.z.ai/api/mcp/web_search_prime/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Crush**

Add MCP server configuration in Crush settings:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
    "$schema": "https://charm.land/crush.json",
    "mcp": {
        "web-search-prime": {
            "type": "http",
            "url": "https://api.z.ai/api/mcp/web_search_prime/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Goose**

Not support Goose now，refer [Issue](https://github.com/block/goose/issues/6576)

Add MCP server configuration in Goose settings:

Click `Extensions` -> `Add custom extension`

Set `Extension Name` is `web-search-prime`，`Type` switch `HTTP`，`Endpoint` as follow：

```
https://api.z.ai/api/mcp/web_search_prime/mcp
```

Set Request Headers Add `Authorization` : `your_api_key`

Finally, click `Add Extension` at the bottom. Remember to replace your\_api\_key with the API Key you obtained in the previous step.

**Roo Code, Kilo Code and Other MCP Clients**

For Roo Code, Kilo Code and other clients that support MCP protocol, use the following general configuration:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
  "mcpServers": {
    "web-search-prime": {
      "type": "streamable-http",
      "url": "https://api.z.ai/api/mcp/web_search_prime/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

## Usage Example

Through the previous step of installing the Search MCP server to the client, you can directly use MCP in your Coding client.\
You can directly use search functionality in conversations:

* "Help me search for the latest AI technology developments"
* "Find best practices for Python asynchronous programming"

## Troubleshooting

> **Invalid API Key**
>
> **Issue:** Receiving invalid api key error
>
> **Solutions:**
>
> 1. Confirm the api key is correctly copied
> 2. Check if the api key is activated
> 3. Confirm the api key has sufficient balance
> 4. Check if the Authorization header format is correct

> **Connection Timeout**
>
> **Issue:** MCP server connection timeout
>
> **Solutions:**
>
> 1. Check network connection
> 2. Confirm firewall settings
> 3. Verify the server URL is correct
> 4. Increase timeout settings

> **Empty Search Results**
>
> **Issue:** Search returns empty results
>
> **Solutions:**
>
> 1. Try using different search keywords
> 2. Check if the search query is too specific
> 3. Confirm network connection is normal
> 4. Contact technical support for assistance

## Related Resources

* [Model Context Protocol (MCP) Official Documentation](https://modelcontextprotocol.io/)
* [Claude Code MCP Configuration Guide](https://docs.anthropic.com/en/docs/claude-code/mcp)
* [MCP Usage Limits](/devpack/overview#usage-instruction)
* [GLM Coding Plan Overview](/devpack/overview)


# Vision MCP Server
Source: https://docs.z.ai/devpack/mcp/vision-mcp-server



The Visual Understanding MCP Server is an exclusive Local MCP Server developed by Z.AI for GLM Coding Plan users. Built on the Model Context Protocol (MCP), it provides visual capabilities such as image analysis and video understanding for MCP-compatible clients, including Claude Code and Cline.

> [!TIP]
>
> Please install the latest version(>= 0.1.2) of the vision mcp server to experience the GLM-5.3-Flash capability. \
> Existing users might still be using a cached older version. Please clear the npx cache, or append the `@latest` tag to `z_ai/mcp-server` to force-install the newest version (i.e., `z_ai/mcp-server@latest`).

> [!NOTE]
>
> Except in Claude Code, pasting an image directly into the client cannot call this MCP Server, as the client will by default transcode the image and call the model interface directly. \
> The best practice is to place the image in a local directory and invoke the MCP Server by specifying the image name or path in the conversation. \
> For example: `What does demo.png describe?`

## Features

> **Image Analysis**
>
> Supports intelligent analysis and content understanding of multiple image formats, giving your AI Agent visual capabilities

> **Video Understanding**
>
> Supports visual understanding of both local and remote videos

> **Easy Integration**
>
> One-click installation, quick integration with Claude Code and other MCP-compatible clients

## Supported Tools

This server implements the Model Context Protocol and can be used with any MCP-compatible client. Currently provides the following tools:

* **`ui_to_artifact`** - Turn UI screenshots into code, prompts, specs, or descriptions.
* **`extract_text_from_screenshot`** - OCR screenshots for code, terminals, docs, and general text.
* **`diagnose_error_screenshot`** - Analyze error snapshots and propose actionable fixes.
* **`understand_technical_diagram`** - Interpret architecture, flow, UML, ER, and system diagrams.
* **`analyze_data_visualization`** - Read charts and dashboards to surface insights and trends.
* **`ui_diff_check`** - Compare two UI shots to flag visual or implementation drift.
* **`image_analysis`** - General-purpose image understanding when other tools don’t fit.
* **`video_analysis`** - Inspect videos (local/remote ≤8 MB; MP4/MOV/M4V) to describe scenes, moments, and entities.

## Environment Variable Configuration

### Detailed Configuration

| Environment Variable | Description | Default Value | Optional Values |
| :- | :- | :- | :- |
| `Z_AI_API_KEY` | Z.AI API KEY | Required | Your API key |
| `Z_AI_MODE` | Service platform selection | Required | `ZAI` |

## Installation and Usage

### Quick Start

**1. Get API Key**

Get your API Key from [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list) (Individual Plan) or [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) (Team Plan)

**2. Install MCP Server**

Prerequisites: [Node.js >= v22.0.0](https://nodejs.org/en/download/) \
According to the client you’re using, **choose the corresponding installation method from the options below**.

### Supported Clients

**Claude Desktop**

**Method A: One-click Installation Command**

Be sure to replace `your_api_key` with the API Key you obtained.

```bash theme={null}
claude mcp add -s user zai-mcp-server --env Z_AI_API_KEY=your_api_key Z_AI_MODE=ZAI -- npx -y "@z_ai/mcp-server"
```

If you forgot to replace the API Key, you need to uninstall the old MCP Server before re-executing the installation command:

```bash theme={null}
claude mcp list
claude mcp remove zai-mcp-server
```

**Method B: Manual Configuration**

Edit Claude Desktop's configuration file `.claude.json` `mcpServers` content: \
Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
  "mcpServers": {
    "zai-mcp-server": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@z_ai/mcp-server"
      ],
      "env": {
        "Z_AI_API_KEY": "your_api_key",
        "Z_AI_MODE": "ZAI"
      }
    }
  }
}
```

**Cline (VS Code)**

Add MCP server configuration in Cline extension settings:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
  "mcpServers": {
    "zai-mcp-server": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@z_ai/mcp-server"
      ],
      "env": {
        "Z_AI_API_KEY": "your_api_key",
        "Z_AI_MODE": "ZAI"
      }
    }
  }
}
```

**OpenCode**

Add MCP server configuration in OpenCode settings:

Refer [OpenCode MCP Doc](https://opencode.ai/docs/mcp-servers)

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
    "$schema": "https://opencode.ai/config.json",
    "mcp": {
        "zai-mcp-server": {
            "type": "local",
            "command": ["npx","-y","@z_ai/mcp-server"],
            "environment": {
                "Z_AI_API_KEY": "your_api_key",
                "Z_AI_MODE": "ZAI"
            }
        }
    }
}
```

**Crush**

Add MCP server configuration in Crush settings:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
    "$schema": "https://charm.land/crush.json",
    "mcp": {
        "zai-mcp-server": {
            "type": "stdio",
            "command": "npx",
            "args": [
                "-y",
                "@z_ai/mcp-server"
            ],
            "env": {
                "Z_AI_API_KEY": "your_api_key",
                "Z_AI_MODE": "ZAI"
            }
        }
    }
}
```

**Roo Code, Kilo Code and Other MCP Clients**

For Roo Code, Kilo Code and other clients that support MCP protocol, use the following general configuration:

Be sure to replace `your_api_key` with the API Key you obtained.

```json theme={null}
{
  "mcpServers": {
    "zai-mcp-server": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@z_ai/mcp-server"
      ],
      "env": {
        "Z_AI_API_KEY": "your_api_key",
        "Z_AI_MODE": "ZAI"
      }
    }
  }
}
```

## Usage Example

Through the previous step of installing the Vision MCP server to the client, you can directly use MCP in your Coding client. \
For example, in Claude Code, inputting `hi describe this xx.png` in the conversation, the MCP Server will process the image and return the description result. (The prerequisite is that you have the image in your current directory)

![Description](https://cdn.bigmodel.cn/markdown/1760501186683image.png?attname=image.png)
![Description](https://cdn.bigmodel.cn/markdown/1782359174332img_v3_02130_be0857a4-7405-40cc-8566-c8f45bbf1f1g.jpg?attname=img_v3_02130_be0857a4-7405-40cc-8566-c8f45bbf1f1g.jpg)

## Troubleshooting

Run the following command in your local terminal to verify if it can be installed locally, to troubleshoot environment, permission, and other issues:

```bash Linux/macOS theme={null}
Z_AI_API_KEY=your_api_key npx -y @z_ai/mcp-server
```

```cmd Windows Cmd theme={null}
set Z_AI_API_KEY=your_api_key && npx -y @z_ai/mcp-server
```

```powershell Windows PowerShell theme={null}
$env:Z_AI_API_KEY="your_api_key"; npx -y @z_ai/mcp-server
```

* If installed successfully, it indicates that the environment is correct, and the issue may be with the client configuration. Please check the client's MCP configuration.
* If installation fails, please troubleshoot based on the error message. It is recommended to paste the error message to a large model for analysis and resolution.

Other common issues:

> **Connection Closed**
>
> **Issue：** Mcp server connection closed
>
> **Solutions：**
>
> 1. Check whether Node.js 22 or a newer version is installed locally.
> 2. Run `node -v` and `npx -v` to verify that the execution environment is available.
> 3. Check the environment variable `Z_AI_API_KEY` is configured correctly.

> **Invalid API Key**
>
> **Issue:** Receiving invalid API Key error
>
> **Solutions:**
>
> 1. Confirm the API Key is correctly copied
> 2. Check if the API Key is activated
> 3. Confirm the selected platform (`Z_AI_MODE`) matches the API Key
> 4. Check if the API Key has sufficient balance

> **Connection Timeout**
>
> **Issue:** MCP server connection timeout
>
> **Solutions:**
>
> 1. Check network connection
> 2. Confirm firewall settings
> 3. Increase timeout settings

## Related Resources

* [Model Context Protocol (MCP) Official Documentation](https://modelcontextprotocol.io/)
* [Claude Desktop MCP Configuration Guide](https://docs.anthropic.com/en/docs/claude-code/mcp)
* [MCP Usage Limits](/devpack/overview#usage-instruction)


# Zread MCP Server
Source: https://docs.z.ai/devpack/mcp/zread-mcp-server



The ZRead MCP Server is an exclusive Remote MCP Server developed by Z.AI for GLM Coding Plan users. Built on the Model Context Protocol (MCP) and powered by zread.ai, it provides access to open-source repository documentation, code structures, and file contents for MCP-compatible clients, including Claude Code and Cline.

## Features

> **Documentation Search**
>
> Search documentation, code, and comments in Github repositories

> **Repository Structure**
>
> Get the directory structure and file list of GitHub repositories to quickly master project layout

> **Code Reading**
>
> Read the complete code content of specified files in GitHub repositories to deeply analyze implementation details

## Tools

This server implements the Model Context Protocol and works with any MCP-compatible client. Currently, it provides the following tools:

* **`search_doc`** — Search for knowledge documentation corresponding to the GitHub repository, quickly understanding repository knowledge, news, recent issues, PRs, and contributors.
* **`get_repo_structure`** — Get the directory structure and file list of the GitHub repository to understand project module splitting and directory organization.
* **`read_file`** — Read the complete code content of specified files in the GitHub repository to deeply analyze the implementation details of the file code.

## Example Scenarios

> **Quick Start with Open Source Libraries**
>
> Quickly understand the core concepts, installation steps, and code organization of open source libraries by searching documentation and obtaining repository structures, accelerating the learning curve.

> **Issue Troubleshooting and History**
>
> When encountering problems, search the repository's Issue and Commit history to find solutions or fix records for similar problems.

> **Deep Source Code Analysis**
>
> Directly read the code content of core files, analyze implementation logic, and assist in secondary development or Debugging.

> **Dependency Library Research**
>
> Before introducing a new dependency library, evaluate its activity, code quality, and maintenance status by viewing its repository structure and documentation.

## Installation and Usage

### Quick Start

**1. Get API Key**

Get your API Key from [Personal Coding Plan > Plan Overview](https://z.ai/manage-apikey/apikey-list) (Individual Plan) or [Team Coding Plan > My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) (Team Plan)

**2. Configure MCP Server**

According to the client you’re using, **choose the corresponding installation method from the options below**.

### Supported Clients

**Claude Code**

**One-click install command**

Replace `your_api_key` with the API key you obtained in the previous step

```bash theme={null}
claude mcp add -s user -t http zread https://api.z.ai/api/mcp/zread/mcp --header "Authorization: Bearer your_api_key"
```

**Manual configuration**

Edit the Claude Code configuration file under your home directory, the MCP section of `.claude.json`:

```json theme={null}
{
  "mcpServers": {
    "zread": {
      "type": "http",
      "url": "https://api.z.ai/api/mcp/zread/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

**Cline (VS Code)**

Add the MCP server configuration in the Cline extension settings:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
  "mcpServers": {
    "zread": {
      "type": "streamableHttp",
      "url": "https://api.z.ai/api/mcp/zread/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

If Cline older version does not support StreamableHttp type MCP server, you can use SSE type configuration:

```json theme={null}
{
  "mcpServers": {
    "zread": {
      "type": "sse",
      "url": "https://api.z.ai/api/mcp/zread/sse?Authorization=your_api_key"
    }
  }
}
```

**OpenCode**

Add the MCP server configuration in OpenCode settings:

See the [OpenCode MCP documentation](https://opencode.ai/docs/mcp-servers)

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
    "$schema": "https://opencode.ai/config.json",
    "mcp": {
        "zread": {
            "type": "remote",
            "url": "https://api.z.ai/api/mcp/zread/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Crush**

Add the MCP server configuration in Crush settings:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
    "$schema": "https://charm.land/crush.json",
    "mcp": {
        "zread": {
            "type": "http",
            "url": "https://api.z.ai/api/mcp/zread/mcp",
            "headers": {
                "Authorization": "Bearer your_api_key"
            }
        }
    }
}
```

**Goose**

Not support Goose now，refer [Issue](https://github.com/block/goose/issues/6576)

Add the MCP server in Goose:

Go to `Extensions` -> `Add custom extension`

Set Extension Name to `zread`, Type to `HTTP`, and use the following endpoint:

```
https://api.z.ai/api/mcp/zread/mcp
```

Set Request Headers Add `Authorization` : `your_api_key`

Click `Add Extension` at the bottom. Remember to replace `your_api_key` with the API key you obtained in the previous step.

**Roo Code, Kilo Code, Others**

For Roo Code, Kilo Code, and other MCP-compatible clients, use the following general configuration:

Replace `your_api_key` with the API key you obtained in the previous step

```json theme={null}
{
  "mcpServers": {
    "zread": {
      "type": "streamable-http",
      "url": "https://api.z.ai/api/mcp/zread/mcp",
      "headers": {
        "Authorization": "Bearer your_api_key"
      }
    }
  }
}
```

## Troubleshooting

> **Invalid access token**
>
> **Issue:** Received an invalid access token error
>
> **Solutions:**
>
> 1. Verify the token was copied correctly
> 2. Check that the token is activated
> 3. Ensure the token has sufficient balance
> 4. Confirm the Authorization header format is correct

> **Connection timeout**
>
> **Issue:** Connection to the MCP server timed out
>
> **Solutions:**
>
> 1. Check your network connection
> 2. Verify firewall settings
> 3. Ensure the server URL is correct
> 4. Increase client timeout settings

> **Repository access failed**
>
> **Issue:** Unable to search or read specified repository content
>
> **Solutions:**
>
> 1. Confirm the repository exists and is open source (public)
> 2. Check if the repository name is spelled correctly (owner/repo)
> 3. Visit zread.ai to search if this open source repository is supported

## Resources

* [Model Context Protocol (MCP) Documentation](https://modelcontextprotocol.io/)
* [Claude Code MCP Configuration Guide](https://docs.anthropic.com/en/docs/claude-code/mcp)
* [MCP Usage Limits](/devpack/overview#usage-instruction)
* [GLM Coding Plan Overview](/devpack/overview)


# API Key Management Guide
Source: https://docs.z.ai/devpack/notice/api-key-management

Starting October 16, 2026, API keys for the standard model API and GLM Coding Plan are managed separately. Learn which key to use, how existing keys are classified, and how billing works for each key type.

* Publication date: October 16, 2026

Starting October 16, 2026, API keys for the standard model API (pay-as-you-go or usage bundles) and GLM Coding Plan are managed separately. When creating a key, choose the entry point for the type of usage you need. Existing keys will be classified automatically; you do not need to replace them as part of this update.

## Which API key should I use?

| Your needs | Key type | Where to get it and how it is billed |
| :- | :- | :- |
| Pay-as-you-go usage or usage bundles | Standard API Key | Open [API Keys](https://z.ai/manage-apikey/apikey-list) and create a key under Standard API Key. Calls use your account balance or applicable usage bundles, not your Coding Plan quota. |
| Individual Plan usage | API Key | Create an API Key under [API Keys → Personal Coding Plan](https://z.ai/manage-apikey/apikey-list). Calls use your Individual Plan quota. |
| Team Plan usage | Team Plan Key | Get your key from [Team Coding Plan → My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan). Usage follows your Team Plan and seat rules. |

**For Standard and Coding Plan keys created from October 16 onward, the key type determines how calls are billed.** Across supported endpoints, Standard keys use standard API billing, while Coding Plan keys use the corresponding plan. Changing the endpoint does not change the key's billing type.

> For example, using a Standard API Key in a coding tool does not consume your Coding Plan quota, even if you have a subscription. It may incur pay-as-you-go charges against your account balance.

## What happens to existing API keys?

Keys created before October 16, 2026 are classified as follows. **Your key values remain unchanged, and this classification update does not require you to reconfigure existing keys.**

| Existing key belongs to | Classification | How it works |
| :- | :- | :- |
| A Team Plan project | Team Plan Key | Follows Team Plan and seat rules. |
| An account that has ever purchased or had an active Coding Plan, outside a Team Plan project | Compatible key | Retains its existing usage and billing rules. See "How Compatible keys work and are billed" below. |
| An account that has never had a Coding Plan | Standard key | Uses standard API billing. |

Some existing keys for dedicated enterprise services or partner channels retain their previous compatibility behavior. Refer to the key's classification in the console and your existing service agreement.

### How Compatible keys work and are billed

A **Compatible** label means your key retains its existing calling behavior and billing rules. You can continue using it as before. **Billing for an existing Compatible key depends on both the endpoint and whether the account has ever had a Coding Plan, regardless of which list the key appears in.**

Common endpoints and their billing behavior are listed below:

<table>
  <thead>
    <tr>
      <th>API / endpoint</th>
      <th>Base URL configured in your tool</th>
      <th>Billing for a Compatible key</th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>OpenAI Chat Completions — standard API</td>
      <td>[https://api.z.ai/api/paas/v4](https://api.z.ai/api/paas/v4)</td>
      <td>Uses your account balance or applicable usage bundles. Does not consume Coding Plan quota.</td>
    </tr>

    <tr>
      <td>OpenAI Chat Completions — Coding Plan endpoint</td>
      <td>[https://api.z.ai/api/coding/paas/v4](https://api.z.ai/api/coding/paas/v4)</td>

      <td>
        <ul>
          <li>If the account has ever had a Coding Plan, calls use plan quota. Calls stop when the plan expires.</li>
          <li>If the account has never had a Coding Plan, calls use the account balance or applicable usage bundles.</li>
        </ul>
      </td>
    </tr>

    <tr>
      <td>Anthropic Messages — general endpoint</td>
      <td>[https://api.z.ai/api/anthropic](https://api.z.ai/api/anthropic)</td>
    </tr>

    <tr>
      <td>OpenAI Responses — general endpoint</td>
      <td>[https://api.z.ai/api/v1](https://api.z.ai/api/v1)</td>
    </tr>
  </tbody>
</table>

> For example, if your account currently has an active Coding Plan with available quota, the same existing Compatible key uses plan quota at [https://api.z.ai/api/coding/paas/v4](https://api.z.ai/api/coding/paas/v4). If you change the Base URL to [https://api.z.ai/api/paas/v4](https://api.z.ai/api/paas/v4), calls use standard API billing and **may deduct from your account balance even while plan quota remains**.

## FAQ

### How does this update affect me? Do I need to reconfigure anything?

* **You can keep using your existing keys.** This update does not require you to replace keys or reconfigure your tools.
* **Choose the right entry point when creating a new key.** To use Coding Plan quota, get a key from [API Keys → Personal Coding Plan](https://z.ai/manage-apikey/apikey-list) or [Team Coding Plan → My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan). For pay-as-you-go usage or usage bundles, create a Standard API Key on the [API Keys page](https://z.ai/manage-apikey/apikey-list).

### Why is my account balance being charged when I have a Coding Plan?

Check the key configured in your tool. **Standard keys do not use Coding Plan quota. An existing personal Compatible key also uses standard API billing when calling the standard API endpoint.** To use plan quota, get a key from the corresponding Coding Plan entry point and replace the key in your tool. For setup instructions, see the [Quick Start](https://docs.z.ai/devpack/quick-start).

### How do I switch an existing Compatible key to pay-as-you-go after my plan expires?

After your plan expires, calls made with an existing Compatible key to plan endpoints do not automatically switch to pay-as-you-go billing. To continue:

1. Open the [API Keys page](https://z.ai/manage-apikey/apikey-list) and create a Standard API Key.
2. Replace the Compatible key in your tool or application with the new Standard key. **You can keep the same endpoint.**

Make sure your account has an available balance or applicable usage bundles. Subsequent calls use standard API billing instead of Coding Plan quota.

If on-demand usage overage is enabled for your Team Plan, it continues to follow the existing settings and billing rules.

### Why can't I find the API Key entry for my Individual Plan?

The **Personal Coding Plan** tab on the [API Keys page](https://z.ai/manage-apikey/apikey-list) appears for accounts that have subscribed to an Individual Plan. If it is missing, make sure you are signed in to the account used for that subscription. Team Plan members should get their keys from [Team Coding Plan → My Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan).

***

If your key's classification or API responses do not match the expected behavior, contact [Technical Support](https://z.ai/contact) with the error details and request ID. Keep your API keys secure and never share a complete key.


# Plan Update Announcement
Source: https://docs.z.ai/devpack/notice/usage-revision

This notice explains the plan update, how your current benefits are protected, and what happens to future subscriptions.

* Publication date: July 30, 2026
* All dates and times mentioned in this notice are in Singapore Standard Time (UTC+8).

To make usage clearer and easier to predict, GLM Coding Plan is moving to a credits-based system. See the [Subscription](https://z.ai/subscribe) page and the usage guides for [Individual Plans](https://docs.z.ai/devpack/overview#usage-instruction) and [Team Plans](https://docs.z.ai/devpack/teamplan#usage-details) for new plan tiers, prices, credit amounts, and deduction rules.

> [!TIP]
>
> **Key points**
>
> * The new credits-based plan is now available. Previous plans are no longer sold to new users.
> * Existing active subscriptions are not affected. In addition, usage on weekends will be deducted at off-peak rates all day.
> * Users on plans discontinued on July 30 can still renew or upgrade them.
> * Discontinued and new plans differ only in how usage is calculated. Other benefits will remain aligned.

## Identify Your Current Plan

Sign in and use the table below to identify your current plan and see how this update affects you.

| Plan type used in this notice | How to identify it |
| :- | :- |
| Legacy Plan V1 (individual plans only) | [My plan](https://z.ai/manage-apikey/coding-plan/personal/my-plan) shows "Legacy Plan V1" |
| Legacy Plan V2 (individual plans only) | [My plan](https://z.ai/manage-apikey/coding-plan/personal/my-plan) shows "Legacy Plan V2" |
| Team Plan | [My Plan](https://z.ai/manage-apikey/coding-plan/team/plans) shows “Team Edition” |
| No active plan | Your plan has expired, or you have not subscribed |

## How This Update Affects You

### 1. If you have a legacy plan V1

* Your current price, benefits, usage limits, and calculation method stay unchanged until the end of the current billing cycle.
* You cannot switch to the new credits-based plan before your current plan expires. After it expires, you can subscribe to a currently available plan.

### 2. If you have a legacy plan V2

* Your current price, benefits, usage limits, and calculation method stay unchanged until the end of the current billing cycle.
* While your plan is active, you can keep auto-renew enabled and upgrade your tier or billing term.
* You may switch immediately only when upgrading to a higher tier of the new credits-based plan. For a same-tier switch or downgrade, wait until your current plan expires. If you switch, your current plan ends immediately and its unused value is applied toward the new plan price.

### 3. If you have a team plan

* Your current price, benefits, usage limits, and calculation method stay unchanged until the end of the current billing cycle.
* You cannot switch to the new credits-based plan before your current plan expires. After it expires, you can subscribe to a currently available plan.

### 4. If you have no active plan or are subscribing for the first time

* Starting July 30, 2026, you can subscribe directly to the new credits-based plan. See the subscription page for current prices, credit amounts, and plan rules.

## FAQ

**Q: Why is the usage calculation changing?**

**A:** We want plan usage to be fully transparent. The new credits-based plan is based on token usage, so credit deductions better reflect actual model usage. You can also see how input, cached, and output tokens contribute to total usage and plan your usage more easily.

**Q: Will this update affect my current plan?**

**A:** No. If you subscribed before July 30, 2026, your active plan keeps its current price, benefits, usage limits, and calculation method until the end of the current billing cycle.

**Q: Can I renew or upgrade a legacy plan V2 discontinued on July 30?**

**A:** Yes. If you currently have a legacy plan V2, you can enable auto-renew or upgrade your plan.

**Q: When can I switch from a discontinued plan to a new plan?**

**A:** It depends on your current plan:

* **Legacy Plan V1**: You can subscribe to the new credits-based plan after your current plan expires.
* **Legacy Plan V2**: You can switch immediately when upgrading to a higher-tier new credits-based plan. For a same-tier switch or downgrade, wait until your current plan expires.

Plans are not switched automatically. You must subscribe or switch manually.

**Q: Can I still use the migration discount I received on April 30?**

**A:** Yes. If you met the eligibility requirements in the [Legacy Plan Migration Notice](/devpack/transition), your **50% Migration Support discount** remains valid through its original validity period and can be used for the new credits-based plan.

## Usage Reference for Legacy Plans

### Legacy Plan V2

To manage resources and ensure fair access for all users, we apply usage limits on a 5-hour and weekly basis. You can check your quota consumption progress in [Usage Statistics](https://z.ai/manage-apikey/subscription).

One prompt refers to one query. Each prompt is estimated to invoke the model 15–20 times.

**The monthly available quota is converted based on API pricing, equivalent to approximately 15–30× the monthly subscription fee (weekly caps already factored in).**

| Plan Type | 5-Hour Limit (Dynamically refreshed; quota resets 5 hours after consumption) | Weekly Limit (Activated upon subscription; resets every 7 days) |
| - | - | - |
| Lite Plan | Up to approx. 80 prompts | Up to approx. 400 prompts |
| Pro Plan | Up to approx. 400 prompts | Up to approx. 2,000 prompts |
| Max Plan | Up to approx. 1,600 prompts | Up to approx. 8,000 prompts |

> [!TIP]
>
> The above figures are estimates. Actual available usage may vary depending on project complexity, repository size, and whether auto-accept is enabled.
>
> <li>
>   **GLM-5.3**: As the flagship model, API calls consume quota at a rate of "1× during off-peak hours and 3× during peak hours".
> </li>
>
> <li>
>   **GLM-5.3-Flash**: API calls consume quota at a rate of 0.4× during off-peak hours and 1.2× during peak hours.
> </li>

> [!NOTE]
>
> Peak hours: Monday to Friday, 14:00–18:00 Singapore Standard Time (UTC+8).

### Team Plan

<table>
  <thead>
    <tr>
      <th>
        Plan Type
      </th>

      <th>
        Quota per 5 hours / seat

        (Dynamically refreshed; quota is restored 5 hours after request consumption)
      </th>

      <th>
        Weekly quota / seat

        (Starts at purchase; quota refreshes every 7-day cycle)
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        Standard Seat
      </td>

      <td>
        GLM-5.3: Up to 60 M tokens

        GLM-5.3-Flash: Up to 150 M tokens
      </td>

      <td>
        GLM-5.3: Up to 300 M tokens

        GLM-5.3-Flash: Up to 750 M tokens
      </td>
    </tr>

    <tr>
      <td>
        Premium Seat
      </td>

      <td>
        GLM-5.3: Up to 160 M tokens

        GLM-5.3-Flash: Up to 400 M tokens
      </td>

      <td>
        GLM-5.3: Up to 800 M tokens

        GLM-5.3-Flash: Up to 2,000 M tokens
      </td>
    </tr>

    <tr>
      <td>
        Notes
      </td>

      <td>
        “Maximum” refers to the total number of Tokens that can actually be consumed when the **off-peak consumption multiplier** applies.

        The current quota consumption rules for each model are as follows:

        <ul>
          <li>
            **GLM-5.3**: As the flagship model, API calls consume quota at a rate of "1× during off-peak hours and 3× during peak hours".
          </li>

          <li>
            **GLM-5.3-Flash**: API calls consume quota at a rate of 0.4× during off-peak hours and 1.2× during peak hours.
          </li>
        </ul>

        \*Peak hours: Monday to Friday, 14:00–18:00 Singapore Standard Time (UTC+8).
      </td>
    </tr>
  </tbody>
</table>


# Overview
Source: https://docs.z.ai/devpack/overview



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
* **Works with 20+ Agent Tools:** Beyond Claude Code, it also supports ZCode, Cline, OpenCode, AutoClaw, and other [supported tools](https://docs.z.ai/devpack/tool/others#step-1-supported-tools), giving you flexibility across both coding and office workflows.
* **Generous Usage at a Fair Price:** Get higher call limits than standard plans. Starting at just 18 USD per month, with Pro and Max plans designed for high-frequency, complex projects.
* **Expanded Capabilities:** All plans support Vision Understanding, Web Search MCP， Web Reader MCP and Zread MCP helping you tackle a wider range of coding and office tasks.

## Benefits

### Supported Models

* All plans support **GLM-5.3**, GLM-5.3-Flash.
* Requests for GLM-5.2/GLM-5.1 will be automatically routed to GLM-5.3, requests for GLM-4.7 will automatically be routed to GLM-5.3-Flash.

### Usage Instruction

> [!WARNING]
>
> For information on Team Plan usage limits, please visit [Team Plan Benefits](/devpack/teamplan).

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
      <td>Model</td>
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
      <td>MCP Server</td>
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

> **[Vision Understanding](/devpack/mcp/vision-mcp-server)**

> **[Web Search](/devpack/mcp/search-mcp-server)**

> **[Web Reader](/devpack/mcp/reader-mcp-server)**

> **[Zread](/devpack/mcp/zread-mcp-server)**

## Next Steps

> **[Quick Start](/devpack/quick-start)**
>
> Get up and running in minutes — from subscribing to the plan to using it in your agent tools.

> **[Usage Policy](/devpack/usage-policy)**
>
> Learn about account usage rules, rate limits, refund policies, and other important guidelines.

> **[FAQ](/devpack/faq)**
>
> Find answers to common questions about subscriptions, promotions, and using the plan.


# Quick Start
Source: https://docs.z.ai/devpack/quick-start



> [!TIP]
>
> This guide will help you get started with [GLM Coding Plan](https://z.ai/subscribe?utm_source=zai\&utm_medium=link\&utm_term=quickstart\&utm_campaign=Platform_Ops&_channel_track_key=DRUfXN42) in minutes—from subscribing to using GLM models in officially [supported tools and products](https://docs.z.ai/devpack/tool/others#step-1-supported-tools).

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
> API Keys for Individual Plans, Team Plans, and the standard model API are now managed separately. Newly created API Keys for these three product types are not interchangeable—use the one that matches your needs. For details, see the [API Key Management Guide](/devpack/notice/api-key-management).

> [!WARNING]
>
> Safeguard your API Key by keeping it confidential and avoiding hard-coding it in your code.

**4. Connect a Tool**

The GLM Coding Plan is strictly limited to use within officially [supported tools and products](https://docs.z.ai/devpack/tool/others#step-1-supported-tools). Click a tool below to open its configuration guide:

> **[Claude Code](https://docs.z.ai/devpack/tool/claude)**

> **[Roo Code](https://docs.z.ai/devpack/tool/roo)**

> **[Kilo Code](https://docs.z.ai/devpack/tool/kilo)**

> **[Cline](https://docs.z.ai/devpack/tool/cline)**

> **[OpenCode](https://docs.z.ai/devpack/tool/opencode)**

> **[OpenClaw](https://docs.z.ai/devpack/tool/openclaw)**

> **[Crush](https://docs.z.ai/devpack/tool/crush)**

> **[Goose](https://docs.z.ai/devpack/tool/goose)**

> **[Cursor](https://docs.z.ai/devpack/tool/cursor)**

> **[Other Tools](https://docs.z.ai/devpack/tool/others)**

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
> For detailed usage instructions, refer to the [Vision MCP Server](/devpack/mcp/vision-mcp-server) documentation.

> **Web Search MCP Server (Coding Plan Exclusive)**
>
> All users can utilize the Web Search MCP Server to access the latest technical information.
>
> * Search for the latest technical documentation and API changes
> * Obtain the latest information on open-source projects
> * Find solutions and best practices
>
> For detailed usage instructions, refer to the [Web Search MCP Server](/devpack/mcp/search-mcp-server) documentation.

> **Web Reader MCP Server (Coding Plan Exclusive)**
>
> All users can utilize the Web Reader MCP Server to fetch full webpage content and extract structured data.
>
> * Fetch complete webpage content including text, and links
> * Extract structured data such as title, body, and metadata
> * Remote HTTP-based MCP service, no local installation required
>
> For detailed usage instructions, refer to the [Web Reader MCP Server](/devpack/mcp/reader-mcp-server) documentation.


# Best Practice
Source: https://docs.z.ai/devpack/resources/best-practice

Best Practices for Coding Agents: Managing Prompts, Plans, Skills, and Workflows

As foundation models continue to improve, AI coding tools are evolving from simple code-completion assistants into **coding agents** that can participate in the full software development lifecycle. Unlike traditional copilot-style tools, coding agents do more than generate code from prompts. They can read and navigate codebases, modify files, run commands, invoke external tools, and complete complex tasks through multi-step interaction.

As this shift continues, developers need more than prompt-writing techniques. They need a reliable way to work with coding agents in practice. Across leading tools, a shared usage pattern is starting to emerge: provide clear task context, plan execution steps, capture project-level guidance, connect external tools and systems, and automate repetitive workflows so agents can collaborate effectively in real development environments over time.

Drawing on official guidance from these tools, this article outlines a more general **best practices framework for coding agents**.

## **1. Treat Coding Agents as a Collaborator, Not One-Off Assistants**

A common mistake when using a coding agent is to treat it like a one-off question-and-answer tool:

> Ask a question, get a piece of code, and end the interaction.

In practice, that approach does not make full use of what a coding agent can do.

A coding agent is better understood as a configurable collaborator that can be refined over time. Through project-level guidance files, tool integrations, and reusable skills, developers can continuously shape the agent’s behavior so it becomes better aligned with the team’s development workflow.

> [!TIP]
>
> In other words, the value of a coding agent does not come from model capability alone. It comes from the combination of model capability and the development workflow around it.

## **2. Structure Task Inputs: Context Matters More Than Prompting Tricks**

When working with a coding agent, many developers focus too much on prompt-writing techniques and not enough on what matters more: **task context**.

In a complex codebase, an effective task description typically includes four elements:

> **Goal**
>
> Clearly describe what needs to be built or changed, such as fixing a bug, implementing an endpoint, or refactoring a module.

> **Context**
>
> Provide the relevant files, error messages, documentation, or examples. For example, specify which files, functions, or modules are involved.

> **Constraints**
>
> List the engineering requirements the agent should follow, such as coding standards, architectural rules, security requirements, or dependency limitations.

> **Done when**
>
> Define how completion should be evaluated, such as tests passing, behavior changing as expected, or the bug no longer reproducing.

This kind of structured task input reduces unnecessary guesswork and makes the agent’s changes more consistent and easier to review.

In most coding agents, you can provide this context by referencing files, supplying code snippets, or stating the relevant details explicitly in the prompt. Once that context is in place, the next step for complex work is to plan before making changes.

## **3. Plan Before Execution for Complex Tasks**

Once a task has clear context, the next challenge is execution. For complex requests, coding agents are most effective when they plan before they act.

For more complex requests, asking the agent to start writing code immediately often leads to logic errors, unnecessary rework, or repeated revisions. A more effective approach is to have the agent **produce a plan first, then move into implementation**.

This planning phase typically includes:

![Complex Task Planning](https://cdn.bigmodel.cn/markdown/1774428526948mermaid-diagram%20%282%29.png?attname=mermaid-diagram+%282%29.png)

Claude Code, for example, encourages an analysis-and-planning step for complex tasks, such as exploring the codebase, identifying the scope of changes, and confirming the implementation approach before making edits. Some coding agents also provide a dedicated Plan mode that generates a complete execution plan before implementation begins.

This shifts the agent from simply generating code on demand to completing work step by step against an explicit plan.

## **4. Capture Repeated Rules in Project-Level Configuration Files**

In practice, many prompts end up repeating the same project rules, such as:

```
- Project directory structure
- Build commands
- Test workflow
- Coding standards
- PR submission process
```

If these rules are restated in every prompt, the workflow becomes inefficient and the instructions can easily drift over time.

For this reason, most coding agents provide a way to store **long-lived project guidance** in project-level configuration files, so the agent can automatically load the relevant context when carrying out tasks.

In some tools, these configuration files take the form of agent-facing guidance documents that describe the repository structure, how to run the project, and the development conventions to follow.

Other systems capture the same information through configuration files, scripts, or related project settings, allowing the agent to apply consistent project constraints across sessions.

Regardless of the implementation, the goal is the same: move information that would otherwise need to be repeated in conversation into **stable project context**.

> [!TIP]
>
> From a practical standpoint, this can be reduced to one simple rule: **put temporary instructions in the prompt, and put long-lived rules in project-level configuration files.**

## **5. The Execution Environment Defines What the Agent Can Do**

When working with coding agents, developers often attribute inconsistent results to model capability. In practice, many of these issues are caused by an incomplete or poorly configured **execution environment**.

Unlike traditional code-completion tools, coding agents are typically expected to operate in a real development environment and carry out tasks such as:

```
- Reading and modifying source files
- Running build or test commands
- Calling external tools or APIs
- Interacting with version control systems
```

As a result, agent behavior depends not only on model capability, but also on whether its **execution environment is complete, stable, and accessible**. When the environment is not configured properly, the agent can easily run into problems such as:

> [!WARNING]
>
> * Being unable to locate the correct project directory
> * Lacking permission to read or modify critical files
> * Failing to run build or test commands
> * Being unable to access external tools or services

These issues often appear as model misunderstanding or poor code quality, but the underlying problem is usually that the agent does not have enough execution capability or access to the right context.

Most leading coding agents now provide some form of **environment configuration** to define the boundaries of agent behavior within a project. For example:

```
- Setting the default model or reasoning level
- Controlling file permissions and sandbox policies
- Defining which commands the agent is allowed to run
- Configuring connections to external tools or services
```

While the implementation differs across tools, the objective is largely the same: to provide the agent with a **stable, controlled, and repeatable execution environment**.

In practice, a reliable execution environment usually depends on a few categories of configuration:

```
- The project working directory and code access permissions
- The allowed command surface, such as build, test, and lint
- Connections to external tools or data sources
- Shared team defaults for agent behavior
```

Once the execution environment is configured correctly, a coding agent is much more likely to behave consistently across sessions and complete multi-step tasks reliably.

> [!TIP]
>
> At a higher level, a coding agent depends on three types of context:
>
> * **Task Context**: the prompt and input for the current task
> * **Project Context**: the repository structure and engineering rules
> * **Environment Context**: the tools, permissions, and execution environment
>
> Of these, Environment Context determines **what the agent can do, and how far it can go**.

## **6. Involve Coding Agents in the Full Development Loop**

Once a coding agent has the right execution environment, the next step is to involve it in the full development loop rather than using it for code generation alone. In real software development, a code change is rarely judged on generation alone. It also needs to pass tests, comply with engineering standards, and go through appropriate review.

A more effective approach is to have the coding agent participate in the **full development loop**, rather than treating it as a tool for code generation only.

A typical agent-driven development loop includes the following steps:

![Description](https://cdn.bigmodel.cn/markdown/1774428645882mermaid-diagram%20%283%29.png?attname=mermaid-diagram+%283%29.png)

1. **Implement code changes**
   Modify existing code or add new code based on the task requirements.
2. **Write or update tests**
   Add test coverage for new functionality or for the bug being fixed.
3. **Run the test suite**
   Execute unit or integration tests to verify that the changes behave as expected.
4. **Run code checks**
   Run linting, formatting, or type-checking tools to ensure the changes meet engineering standards.
5. **Review the code changes**
   Inspect the diff to identify potential issues, regression risks, or unintended modifications.

In this workflow, the coding agent is no longer just a code generator. It becomes an active participant in **implementation, validation, and review**.

Most leading coding agents now support this way of working. Some can automatically run tests or build commands after making changes. Others can execute validation scripts or trigger review steps once code has been modified. Through these mechanisms, agents can participate continuously in code quality workflows, whether in a local development environment or in CI pipelines.

> [!TIP]
>
> From a workflow perspective, this shifts the coding agent from a traditional **code generator** to an **execution node within the development loop**.

Once an agent is able to participate in the full loop, its role extends beyond writing code. It can also reduce repetitive manual work in testing, validation, and review, improving overall development efficiency.

## **7. Extend Agent Context with MCP**

In real development workflows, the information a coding agent needs does not always live inside the code repository. Much of the context that shapes implementation decisions is often spread across external systems, such as:

```
- Issue tracking and requirements systems
- CI/CD status and execution results
- Database schemas or production data
- API documentation and external service references
```

If this information has to be copied and pasted manually each time, the workflow becomes inefficient and the context passed to the agent is often fragmented and unreliable. For coding agents that are expected to carry out multi-step tasks over time, this approach does not scale well to more complex workflows.

For this reason, many coding agent tools support the **Model Context Protocol (MCP)**, which provides a standard way to connect external tools and systems so the agent can access the real-time information it needs beyond the repository itself.

Through MCP, a coding agent can typically access resources such as:

```
-   Code hosting and collaboration platforms
-   Databases and query interfaces
-   API services and technical documentation
-   Internal tools and automation systems
```

This changes how coding agents obtain information. Instead of relying entirely on what a developer describes in the prompt, the agent can retrieve the relevant context directly from connected tools.

> [!TIP]
>
> From a workflow perspective, this is a meaningful shift. When an agent can only use the information provided in the prompt, it is usually limited to localized tasks. Once it can connect to external systems, it becomes capable of participating in more complete development workflows, such as reading issue context, investigating failed CI runs, checking API definitions, or analyzing problems against database schemas.

The value of external tool integration is therefore not just that it gives the agent more tools. It expands the agent’s context boundary, allowing it to evolve from a **repository-level executor** into **a collaborative node that operates within a real engineering environment**.

## **8. Capture Repeated Workflows as Skills**

Over time, teams often find that certain tasks come up again and again when working with coding agents. Common examples include:

```
- PR review
- Log analysis
- Release note generation
- Standard debugging workflows
```

If these tasks are described manually in the prompt each time, the result is unnecessary repetition and less consistent outcomes.

For this reason, many coding agent systems provide a **Skill** mechanism for packaging common workflows into reusable workflow templates.

At a high level, a Skill can be understood as a **structured workflow template**. It abstracts execution logic that would otherwise be scattered across prompts, allowing the agent to apply the same workflow consistently when handling similar tasks.

> [!TIP]
>
> For a deeper discussion of what Skills are and how to write them, see [How to create Skills: Key steps, limitations, and examples](https://claude.com/blog/how-to-create-skills-key-steps-limitations-and-examples).

Different coding agent tools manage Skills in different ways. Some define them through dedicated skill files, while others register reusable workflows through configuration or scripting mechanisms.

Although the implementations vary, the goal is the same: **turn ad hoc prompting into reusable workflows**.

In practice, one simple rule works well:

> **If a prompt pattern or task flow is used repeatedly, it should probably be captured as a Skill.**

This shifts the use of coding agents from one-off, conversation-driven interaction toward more workflow-oriented task execution. As the skill library grows, agent behavior also becomes more consistent and predictable.

## **9. Automate Stable Workflows**

Once a Skill can be executed reliably, the next step is often to automate it.

In long-running development workflows, many tasks are inherently repetitive or time-based. For example:

```
- Generate commit summaries on a regular basis
- Automatically investigate failed CI runs
- Scan for potential bugs or abnormal logs
- Produce daily or weekly engineering reports
```

Even when these tasks have already been encapsulated as Skills, they still create unnecessary manual work if developers need to trigger them every time.

For this reason, many coding agent systems provide an **Automation layer** that allows workflows to run automatically based on a schedule or a defined trigger condition.

> [!NOTE]
>
> Automation can be understood as the next layer above Skills. A Skill defines how a workflow is executed, while Automation determines when **that workflow runs and how it continues to operate over time**.

For example, a Skill for generating release notes might be configured to run:

```
-   whenever a new release is published
-   once a week to produce a release summary
-   automatically after CI completes
```

This allows the coding agent to keep performing tasks in the background without requiring developers to start each run manually. The significance of this model is that it shifts the coding agent from an **interactive tool** to a **continuous development assistant**.

Once a workflow has been abstracted into a stable Skill, adding automation can significantly reduce manual effort and allow the agent to provide ongoing support within the development environment.

## **10. Manage Agent Sessions Deliberately**

When working with coding agents, a session is more than just a chat history. In practice, it functions as a **working context** that accumulates context, intermediate reasoning, and execution results over time.

As a task progresses, the agent gradually builds up information within the same session, including:

```
-   the task objective
-   relevant code context
-   changes that have already been made
-   intermediate reasoning and decisions
```

Together, these form the agent’s **active working context** for the task.

If sessions are not managed carefully, unrelated tasks can accumulate in the same session and make the context unnecessarily complex. This often reduces the quality of the agent’s reasoning and execution. For that reason, deliberate session management is an important part of working effectively with coding agents. Common practices include:

> **Use a separate session for each task**
>
> Avoid mixing unrelated tasks in the same session so the working context stays clear.

> **Avoid overly long sessions**
>
> When a session accumulates too much history, use summaries or compression to reduce context overhead.

> **Start a new session for branch explorations**
>
> If the task opens up a new line of investigation, continue it in a separate session instead of piling more changes into the original one.

> **Periodically compress historical context**
>
> Summarize older parts of the conversation to reduce pressure on the context window.

In more complex development scenarios, teams may also adopt a **multi-agent collaboration model**. For example, subtasks such as exploring the codebase, running tests, or investigating failures can be delegated to separate agents, while a primary agent coordinates the overall task. This helps preserve clarity in the main session while improving execution efficiency for more complex work.

At a higher level, session management is a form of **context management**. A well-structured session strategy helps coding agents maintain a clear reasoning path across multi-step tasks, which in turn improves consistency and execution quality.

## Conclusion

The effectiveness of a coding agent does not come from the model alone. It also depends on how developers structure the workflow around it.

In practice, a mature coding agent workflow typically includes the following stages:

![Description](https://cdn.bigmodel.cn/markdown/1774433767376image.png?attname=image.png)

Through this workflow, a coding agent can gradually evolve from a simple code generation tool into a collaborative system that participates across the full software development lifecycle.


# Memory-mechanism
Source: https://docs.z.ai/devpack/resources/memory-mechanism



Memory enables a coding agent to retain context across tasks and sessions, reducing repeated input and improving execution efficiency. With a well-designed memory system, an agent can continuously understand the project structure, engineering conventions, and user preferences, and automatically reuse that information in future work.

In coding agent systems, memory is typically organized into multiple layers, such as **automatic memory, project memory**, and **session memory**.

## **Why do coding agents need memory?**

Traditional large language models do not preserve state between calls. As a result, they cannot remember project context across sessions, accumulate problem-solving experience over time, or consistently adapt to user preferences.

Agent systems address this limitation through **external memory**.

A typical architecture looks like this:

```
User input
   ↓
Memory retrieval
   ↓
Context assembly
   ↓
LLM reasoning
   ↓
Action / tool call
   ↓
Memory update
```

In other words, the agent retrieves relevant memory before starting a task and updates memory after the task is completed.

This architecture is a common pattern in modern agent systems such as LangGraph, AutoGPT, and Devin.

## **A complete memory architecture for modern coding agents**

At a high level, a complete agent memory architecture typically looks like this:

```
Short-term memory
    ↓
Session context

Long-term memory
    ├ semantic memory
    ├ episodic memory
    └ procedural memory
```

## **Core memory types in coding agents**

**Session Memory**

Session memory is the contextual information associated with the current task. It includes the current conversation history, recent tool outputs, the current execution plan, and the contents of the files currently in scope. This information typically lives in the model’s context window.

For example:

```
User: Fix this Python bug
Agent: Analyze the error
Agent: Modify the code
Agent: Run tests
```

These execution steps all fall under session memory.

**Project Memory**

Project memory stores **long-lived information about the entire codebase**, such as the project architecture, coding standards, build workflows, and frequently used commands. This kind of memory is typically written into .md files and loaded at the start of a session.

For example:

```
your-project/
├── .claude/
│   ├── CLAUDE.md           # Main project instructions
│   └── rules/
│       ├── code-style.md   # Code style guide
│       ├── testing.md      # Testing conventions
│       └── security.md     # Security requirements
```

This way, the Agent will automatically follow these rules when modifying the code.

**Semantic Memory**

Semantic memory stores factual knowledge and reference information. Examples include API documentation, programming language rules, and project knowledge bases. In practice, this is often implemented through RAG (retrieval-augmented generation).

A typical flow looks like this:

```
query
↓
embedding
↓
vector search
↓
retrieve documents
↓
LLM reasoning
```

This is also one of the most commonly used memorization methods at Coding Agent.

**Episodic Memory**

Episodic memory records an agent’s past experiences, such as the steps used to fix a previous bug, the root cause of a previous build failure, or a debugging strategy that worked well before. This type of memory helps the agent learn from prior experience.

For example:

```
Episode:
CI failure caused by missing dependency
Solution: upgrade pip package
```

**Procedural Memory**

Procedural memory stores strategies or step-by-step workflows for completing tasks.

For example:

```
Debug_Workflow.md
1. read error log
2. locate file
3. write patch
4. run tests
```

These memories are typically used in system prompt engineering, workflow templates, and agent policies.

## **The standard memory pattern used by coding agents**

In real-world systems, agents typically follow a consistent memory workflow.

**1. Memory retrieval**

Before starting a task, the agent retrieves relevant project memory, knowledge base entries, and prior experience, then injects them into the working context.

**2. Context construction**

The retrieved memories are assembled into a complete context and passed to the model.

**3. Memory update**

After the task is completed, the agent decides whether to write new memories—for example, newly discovered project rules, debugging experience, or user preferences.

## **How to use memory correctly in coding agents**

In mainstream agent systems, memory is generally designed to be **layered, controllable, retrievable, and updatable**.

> [!TIP]
>
> In most cases, memory is divided into **short-term memory** and **long-term memory**. Short-term memory is mainly used to preserve state within the current thread or session, while long-term memory is maintained through explicit files, rule configurations, vector retrieval, or other persistent storage mechanisms.

> [!NOTE]
>
> Take **Claude Code** as an example. Its official documentation explicitly states that each session begins with a fresh context window. Knowledge is carried across sessions primarily through persistent instruction files such as CLAUDE.md and through **auto memory**. Similarly, in **LangChain / LangGraph**, memory is also divided into **thread-scoped short-term memory** and **long-term memory** that persists across sessions.

In practice, the most effective approach is not to rely on the model to “remember everything” automatically, but to establish a clear memory management pattern. For example: what should be written into project-level memory files, what should be retrieved from a knowledge base or vector store, what only needs to remain in the current session, and what is worth promoting into long-term memory after a task is completed.

### \* Separate instruction memory from learning memory

One of the most practical principles for general-purpose coding agents is to distinguish between two fundamentally different kinds of memory:

* **Instruction memory**: written by humans to tell the agent how it should work. This usually includes coding standards, directory conventions, build commands, test procedures, naming conventions, commit requirements, and team-level safety rules. In Claude Code, this maps to persistent instruction files such as `CLAUDE.md`.
* **Learning memory**: not predefined in advance, but accumulated by the agent over time from your corrections, preferences, failed attempts, common commands, and project habits. Claude Code refers to this capability as auto memory, and its documentation states that it is loaded at the start of every conversation together with instruction files. For subagents, Claude Code can also maintain a separate persistent memory directory, and the first 200 lines of `MEMORY.md` are included automatically.

> [!TIP]
>
> If these two types of memory are mixed together, system behavior often drifts over time. A better approach is to assign them clearly separated roles:
>
> * Write **rules, policies, and behavioral constraints** into **instruction memory**, so the agent’s behavior stays stable and predictable.
> * Write **experience, user preferences, temporary discoveries, and retrospective takeaways** into **learning memory**, so decisions can improve over time in future tasks.
>
> This separation helps prevent experience-driven notes from gradually polluting the system’s core rules, which in turn keeps agent behavior more stable and controllable.

### \* Layered memory management

> **Organization-level memory**
>
> This layer contains rules defined and distributed at the team or company level, and applies across all developers and all relevant projects. Typical examples include:
>
> * security and compliance requirements
> * baseline code review standards
> * restricted directories that must not be read from or written to
> * dependency and license constraints
> * organization-wide engineering standards
>
> At the organizational level, a shared `sysytem.md` can be deployed to a system-level path and should not be easily excluded by individual users. In practice, this can also be distributed through centralized management tools such as hosted configuration, MDM, Group Policy, or Ansible. In a more general agent architecture, this means **organization-level memory should be treated as the highest-priority governance layer and should not be casually bypassed**.

> **Project-level memory**
>
> This is the team-shared project context, and it is the most important memory layer for a coding agent. It should be version-controlled and shared across all collaborators. Typical examples include:
>
> * project architecture documentation
> * directory structure conventions
> * build and test commands
> * where APIs should live
> * naming conventions
> * common development workflows
>
> Claude Code recommends storing this kind of information in a project-level `project.md`, and its `/init ` command can generate an initial draft automatically. That draft can then be refined with rules the model is unlikely to infer on its own. The key property of this layer is that it is **shared across the project, tracked in version control, and stable over time**.

> **User-level memory**
>
> This layer captures a developer’s personal preferences that apply across projects. It is best stored under the user’s home directory and treated as reusable personal context for all workspaces. In Claude Code, user instructions are stored separately from project instructions, and both are loaded at session start. This layer is a good place for:
>
> * your preferred coding style
> * your usual debugging sequence
> * your preferred output format
> * your personal workflow shortcuts
>
> It should complement project conventions, not override them.

> **Local memory**
>
> This layer is specific to your local copy of a project, but should not be committed to Git. A file such as local.md is a good place to store project-specific preferences that should remain private or machine-specific, such as:
>
> * personal test accounts
> * local development ports
> * temporary mock service endpoints
> * machine-specific runtime notes
> * experimental workflows that are not ready to share
>
> The value of this layer is that it **allows individuals to work efficiently without polluting team-shared memory**.

> **Subagent / role-specific memory**
>
> Another pattern worth generalizing is role-specific memory for subagents. Different subagents can maintain their own memory scopes rather than sharing a single global memory. This is especially important in multi-agent systems, where one of the most common failure modes is context pollution across roles.
> A better pattern is to let each subagent retain only the memory relevant to its role:
>
> * let the **testing agent** remember test commands, CI behavior, and assertion style
> * let the **refactoring agent** remember module boundaries, restricted dependencies, and migration strategies
> * let the **documentation agent** remember glossary terms, documentation templates, and audience-specific style
>
> This keeps memory shorter, more precise, and more stable.

### \* Loading `.md` files by path

Claude Code’s official documentation offers a very useful pattern for organizing memory in large codebases. For larger repositories, it recommends splitting instructions into multiple Markdown files under `.claude/rules/`, with each file focused on a single topic such as `testing.md`, `api-design.md`, or `security.md`.

Claude Code also supports **scoping rules to specific subdirectories or file types**, and these rules are loaded only when Claude is working with matching files. That reduces irrelevant noise and helps conserve context window space.

As a general design pattern for coding agents, this can be summarized in three principles:

* **Keep the main memory file limited to global shared context**, such as project background, high-level architecture, and cross-project conventions.
* **Keep specialized rules modular**, with one rule file per topic.
* If a rule can be loaded by path, **do not load it globally; bring it into context only when needed**.

Based on this approach, a project’s memory structure could look like this:

```
agent-memory/
├── project.md            # Project overview
├── rules/
│   ├── code-style.md     # Code style
│   ├── testing.md        # Testing conventions
│   ├── api-design.md     # API design guidelines
│   ├── security.md       # Security requirements
│   └── frontend/
│       └── react.md      # Frontend-specific rules
└── local/
    └── developer.local.md
```

> [!TIP]
>
> This structure offers three advantages:
>
> 1. Easier to maintain. Each rule file focuses on a single topic, so the rule set is less likely to become bloated or disorganized. Claude Code explicitly recommends topic-specific files with descriptive names.
> 2. Easier to load on demand. When the agent is working on tests, it does not need to load frontend conventions or database-specific rules into the context window.
> 3. Better for team collaboration. Different teams or subteams can maintain their own rule directories instead of competing to edit a single monolithic instruction file.

### \* Write memory rules as concrete instructions

When writing agent memory, use **specific, verifiable rules** whenever possible rather than abstract principles. The clearer the instructions are, the more stable the agent’s behavior will be.

In general, it is recommended to:

* keep instructions **concise and explicit**
* keep rules **consistent** with one another
* keep the main memory file **under 200 lines** where possible
* use **Markdown headings and lists** to improve readability
* phrase requirements as rules that can be **checked and executed**

> [!WARNING]
>
> For example, avoid writing:
>
> * Keep the code clean
> * Write good tests
> * Be mindful of API design
> * Split modules when appropriate

> [!TIP]
>
> Instead, prefer rules like:
>
> * Use **2-space indentation** in all new TypeScript files
> * **Run `pnpm test`** after modifying business logic
> * Place **all API handlers under `src/api/handlers/`**
> * Keep React page components **under 300 lines**; split larger ones into hooks or child components

Concrete rules significantly reduce the agent’s room for interpretation, which improves behavioral consistency.

### \* Separate shared rules from personal preferences

When designing an agent memory structure, it is important to clearly define **the scope of each rule and who is responsible for it**. A common approach is to organize memory by scope:

* **Project**: shared by all team members and maintained through version control
* **Organization**: defined centrally by IT or DevOps, such as security standards or development processes
* **User**: applies only to an individual, such as personal coding habits
* **Local**: applies only to the current machine or working environment and should not be committed to Git
* **Role / Agent-specific**: used only by a specific specialized agent

The core principle of this hierarchy is:

> who owns it, who shares it, and who it applies to.

For example:

* team-wide conventions → project level
* company security policies → organization level
* personal coding habits → user level
* machine-specific configuration → local level
* rules for a specialized agent → role level

Defining these boundaries during the memory design phase helps avoid rule sprawl and duplicate definitions.

### \* Reuse memory through imports and rule packages

In real projects, many rules are **shared engineering conventions across repositories**. Rewriting them in every repo increases maintenance overhead and makes inconsistency more likely.

Using **Claude Code** as an example, its documentation explains that:

* `CLAUDE.md` can import other rule files using `@path/to/import`
* `.claude/rules/` can share rules through **symbolic links** (symlinks)
* imported content can be **expanded recursively**, and symlinks are resolved normally

This makes it possible for teams to build **reusable rule packages**, such as:

* `company-security-rules`
* `frontend-react-rules`
* `backend-api-rules`
* `python-testing-rules`

Each project only needs to reference the rule modules it needs, rather than maintaining a full copy of the entire rule set.

This approach brings two direct benefits:

1. **rules can be maintained centrally and updated consistently**
2. **different projects can share the same engineering language**, making agent behavior more consistent across repositories

## **Memory troubleshooting**

> **The coding agent is not following my `.md` memory files**
>
> `.md` memory files are typically provided to the agent as contextual instructions, not as enforced configuration.
>
> The agent will read them and try to follow them, but it cannot guarantee strict compliance when the rules are vague, unclear, or conflicting.
>
> If the agent is not following the rules, you can check the following:
>
> * Run `/memory` (or the equivalent command) to confirm that the .md memory files have been loaded.
> * Check whether the `.md` files are located in a path or scope that is allowed to load in the current session.
> * Check whether there are conflicting rules across multiple `.md` files. If different files give different instructions for the same behavior, the agent may choose one arbitrarily.

> **I don’t know what auto memory has saved**
>
> Most coding agents maintain auto memory in the background to capture project context, user preferences, or common actions.
>
> You can inspect it in the following ways:
>
> * Run `/memory` (or a similar command) to view the current auto memory directory.
> * Auto memory is typically stored as Markdown files that you can read, edit, or delete directly.

> **My memory files are too large**
>
> Oversized memory files consume more of the context window, reduce the agent’s adherence to instructions, and increase the likelihood of conflicts.
>
> It is recommended to split detailed content into multiple Markdown files and use file references or imports (such as `@path/to/file`), or move rules into a dedicated rules directory such as `rules/`.

> **Instructions disappear after context compression**
>
> Many coding agents **compress or summarize context** during long conversations in order to reduce context length.
>
> In most cases, memory files are **reloaded from disk** after compression, so only content that has been written into memory files will persist. If certain rules disappear after compression, that means those rules **existed only in the conversation** and were never written into a memory file.
>
> To fix this:
>
> * write long-term instructions into `.md` memory files
> * do not rely on the conversation alone to preserve rules


# Team Plan Benefits
Source: https://docs.z.ai/devpack/teamplan

Learn about GLM Coding Team Plan usage quotas, benefits, and usage rules

GLM Coding Team Plan is a self-service subscription for enterprises and teams. Building on the individual plan’s generous access to Z.AI’s top-tier models and broad agent tool compatibility, it adds flexible organization management, enterprise-grade data security, and centralized billing and invoicing—helping teams scale AI-powered work efficiently with predictable costs.

## Exclusive Capabilities

* **Centralized seat and access management**: Manage seats, roles, and permissions in one place, with clear visibility into personnel changes, access updates, and resource usage.
* **Team usage and productivity insights**: Track usage and consumption trends by member and time period to better understand AI adoption and productivity gains.
* **on-demand usage overage and budget control**: Keep services running after the included quota is used by enabling on-demand usage overage. Set per-member spending limits to protect key projects and prevent unexpected costs from high-frequency usage.
  \
  (*Limited-time offer: Overage usage is billed at a 10% discount from the model API list price.*)
* **Centralized billing and invoicing**: Consolidate billing, invoicing, and reconciliation across the organization. Verified enterprises can request special VAT invoices, reducing finance and reimbursement overhead.
* **Data is not used for model training by default**: Code, prompts, conversations, and related content are excluded from model training by default, helping protect your organization’s core R\&D assets.
* **Early access to new flagship models and features (*Premium Seat only*)**: Get priority access to the latest models to help teams continuously improve their AI coding experience and development efficiency.
* **Priority access during peak hours (*Premium Seat only*)**: Benefit from more reliable resource allocation and response times during periods of high demand, with fewer delays, rate limits, and productivity disruptions.

## Usage Details

### Usage Credit Allowance

Each plan is subject to both a 5-hour usage limit and a weekly usage limit. You can check your quota consumption progress in [Usage Statistics](https://z.ai/manage-apikey/coding-plan/team/usage-stats).

| Plan Type | 5-Hour Credits | Weekly Credits |
| :-: | :-: | :-: |
| Standard Seat | 15,000 | 66,000 |
| Premium Seat | 35,000 | 155,000 |

**Credit Reset Rules**

* **5-hour credits**: Dynamically refreshed; credit quota resets 5 hours after consumption.
* **Weekly credits**: Activated upon subscription; resets every 7 days.

### Credit Calculation

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
      <td>Model</td>
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
      <td>MCP Server</td>
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

**During off-peak hours, model usage is charged at 50% of the standard credit rate.**

> [!NOTE]
>
> **Peak hours**: Monday to Friday, 14:00–18:00 Singapore Standard Time (UTC+8).

### Estimated Token Allowance

The plan’s token usage will vary depending on the cache hit rate, as shown below:

| Cache hit rate | Model | Team Standard (million tokens/week) | Team Advanced (million tokens/week) |
| :- | :- | :- | :- |
| 95% | GLM‑5.3 | 319–638 | 749–1,497 |
| 95% | GLM‑5.3‑Flash | 965–1,930 | 2,267–4,533 |
| 96% | GLM‑5.3 | 327–654 | 768–1,536 |
| 96% | GLM‑5.3‑Flash | 990–1,981 | 2,326–4,652 |
| 98% | GLM‑5.3 | 345–689 | 810–1,619 |
| 98% | GLM‑5.3‑Flash | 1,045–2,090 | 2,454–4,908 |

**How the Range Is Calculated**

* Maximum token allowance: All usage occurs during off-peak hours and is charged at 0.5× the standard credit rate.
* Minimum token allowance: All usage occurs during peak hours and is charged at 1× the standard credit rate.

**By fully leveraging off-peak discounts, users can save up to 92% compared with using the GLM-5.3 Standard API on a pay-as-you-go basis.**

## Team Plan Key

The Team Plan Key is the dedicated access credential for the Team Plan. After each team member receives a seat assignment invitation, joins the team, and enters the [Team Plan](https://z.ai/manage-apikey/coding-plan/team/my-plan) page in the console, they can obtain their own Key.

> [!WARNING]
>
> Please note that the **Team Plan Key is independent from other platform API Keys**. To use your Team Plan quota, make sure to use the Team Plan Key in the relevant scenarios.

For details on how API Keys are classified and billed, see the [API Key Management Guide](/devpack/notice/api-key-management).

## Seat Rules

Team plans are subscribed to and assigned by seat:

1. A minimum of 2 seats is required, with no upper limit on the number of seats
2. Members and seats follow a 1:1 relationship; multiple members cannot share the same seat
3. Mixed purchases of Standard Seat and Premium Seat are not currently supported
4. Administrators can reassign seats during the validity period of the plan benefits
5. The seat validity period is the same as the validity period of the plan benefits. After the plan expires, all seat benefits will also expire

## Subscription Changes

**Plan changes:**

1. Continuous subscription users can cancel automatic renewal at any time and re-enable automatic renewal at any time
2. Monthly or annual purchases can be extended by purchasing the plan again
3. Upgrading from the Standard Seat to the Premium Seat is not currently supported

**Seat quantity changes:**

1. Additional seats can be added during the subscription period, with fees calculated based on the remaining time in the current billing cycle
2. Directly reducing the number of seats is not supported. To reduce seats, please purchase again after the current subscription cycle ends

## Account Usage Rules

To protect subscriber rights, ensure system fairness, and maintain service stability, GLM Coding Plan must be used in [officially supported tools and products](https://docs.z.ai/devpack/tool/others#1-coding-agent-tool), and must comply with the [Subscriptions, Fees, and Payment](https://docs.z.ai/legal-agreement/subscription-terms) and related usage rules.

Improper behavior such as multiple people sharing the same seat, use in unsupported tools, or abnormally high-frequency calls may trigger platform risk control rules, resulting in corresponding restrictions on subscription benefits. In serious cases, it may affect normal account usage.

## FAQ

**Q: Can the Team Plan and Individual Plan be active and used at the same time?**

**A:** Yes. Each user can have both an Individual Plan and a Team Plan at the same time, and can also be invited to join different teams and use the plan benefits assigned by those teams. However, within the same team, each member can only have one active Team Plan seat at a time.

**Q: What happens after a seat exceeds its plan quota?**

**A:** Plan usage is limited separately by seat. If a seat exceeds its quota, the model cannot be used during the limit period until the next reset cycle begins. The team administrator can enable on-demand usage overage in advance. After a seat exceeds its usage quota, the service can continue to be used and will be billed based on the actual overage. When the next reset cycle begins, the seat will resume using the quota included in the plan, helping avoid business interruptions.

**Q: How is concurrency limited for each Team Plan seat?**

**A:** Rate limits and concurrency limits are related to your plan tier, and the platform dynamically adjusts them based on available resources. Each development project can use methods such as Subagents to make concurrent model calls. Our recommended number of projects is as follows:

* Standard Seat: recommended for 1–2 concurrent development projects
* Premium Seat: recommended for 2+ concurrent development projects

During off-peak hours, plan users will enjoy higher concurrency benefits through dynamic upgrades, supporting a larger number of development projects.

**Q: Does the primary administrator, meaning the account that purchases the Team Plan, occupy a seat?**

**A:** No. By default, the primary administrator account does not occupy a team seat. If the primary administrator needs to use the quota associated with a seat, they can assign a seat to their own account. After joining the seat, they will receive the corresponding quota.

## Next Steps

* [Quick Start](https://docs.z.ai/devpack/quick-start): Complete the basic integration process in just a few minutes and get started quickly
* [Tool Integration](https://docs.z.ai/devpack/tool/others): View the agent tools supported by the plan and their configuration methods, and choose the working environment that best suits your needs
* [How to Switch Models](https://docs.z.ai/devpack/latest-model): Make sure your coding tool is using your target model version


# Tool Integration
Source: https://docs.z.ai/devpack/tool/others



## Supported Tools

> [!TIP]
>
> The GLM Coding Plan is limited to use within the following officially supported tools and product environments; users may not use their subscription benefits for tools or scenarios outside of this scope.

### 1. Coding Agent Tool

Click on the tool documentation below that you wish to use, and follow the corresponding integration guide to set it up.

> **[ZCode (1.5× Usage)](/devpack/tool/zcode)**
>
> ZCode integrates AI agents into your existing toolchain.

> **[Claude Code](/devpack/tool/claude)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Claude for IDE](/devpack/tool/claude-for-ide)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Codex](/devpack/tool/codex)**
>
> OpenAI's AI coding agent that helps you write, review, and debug code.

> **[OpenCode](/devpack/tool/opencode)**
>
> The Claude Code IDE plugin supports VSCode and JetBrains.

> **[Pi](/devpack/tool/pi)**
>
> A minimalist terminal coding agent with a small core, extensible via extensions, Skills, and Pi Packages.

> **[Cursor](/devpack/tool/cursor)**
>
> An AI-first code editor that supports custom model configurations.

> **[Cline](/devpack/tool/cline)**
>
> An AI programming extension for VS Code that supports code generation and file operations.

> **[TRAE](/devpack/tool/trae)**
>
> An AI editor capable of independently completing various development tasks.

> **[Qoder](/devpack/tool/qoder)**
>
> An agentic coding platform designed for real software development.

> **[Droid](/devpack/tool/droid)**
>
> Enterprise-grade AI coding agent that runs in the terminal to handle end-to-end workflows.

> **[Kilo Code](/devpack/tool/kilo)**
>
> A powerful VS Code extension for code generation and project management.

> **[Roo Code](/devpack/tool/roo)**
>
> A smart VS Code extension for code writing and refactoring.

> **[Crush](/devpack/tool/crush)**
>
> A terminal-based AI programming tool that supports both CLI and TUI interfaces.

> **[Goose](/devpack/tool/goose)**
>
> AI Agent tool, supporting local execution and automated engineering tasks.

> **[Eigent](/devpack/tool/eigent)**
>
> A desktop AI agent built on a multi-agent architecture, capable of automating browser, terminal, and MCP-powered workflows.

### 2. General-purpose Agent Tool

> [!TIP]
>
> The general-purpose agent tools listed below are also supported and will continue to be served on a best-effort basis. Under high inference load, some requests may face temporary rate limits.

> **[AutoClaw (1.5× Usage)](/devpack/tool/autoclaw)**
>
> Deeply optimized for GLM models, integrating specialized domain knowledge bases and workflows.

> **[OpenClaw](/devpack/tool/openclaw)**
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


# Legacy Plan Migration Notice
Source: https://docs.z.ai/devpack/transition



* Publication date: April 21, 2026

Thank you for your continued support of GLM Coding Plan.
To provide a more consistent subscription structure and better support future product upgrades, GLM Coding Plan will gradually migrate remaining legacy subscribers to the current plan system. The details of this update are as follows.

## 1. Changes and Applicability

Legacy plans (without weekly usage limits) will be phased out and will no longer remain as a supported subscription option going forward.
**This update applies to users who are on a legacy plan with auto-renew enabled as of April 30, 2026 (Singapore Time). It does not affect users who are already on the current plans.**

## 2. Details of the Changes

### 2.1 Auto-renew for legacy plans will be turned off

Starting on **April 30, 2026 (Singapore Time)**, auto-renew for eligible legacy plans will be automatically canceled.
This change **will not affect your current paid subscription period**. Your existing plan will remain active until the end of its current billing cycle.

### 2.2 Migration support

To make the migration as smooth as possible, we will provide the following migration support to eligible users.

#### 2.2.1 2 complimentary months of the equivalent online version plan

We will automatically provide 2 complimentary months of the equivalent online version plan tier. This will be issued on April 30, 2026 (Singapore Time) and will be added after your current subscription ends, with no manual action required.

| Current Legacy Plan | Transition Support |
| :- | :- |
| Legacy Lite | 2 months of the online version Lite plan |
| Legacy Pro | 2 months of the online version Pro plan |
| Legacy Max | 2 months of the online version Max plan |

#### 2.2.2 Extra benefits when migrating to the online version plan

We understand that the current plan pricing has already gone through a recent price increase.
As part of this policy, eligible users will be able to purchase the latest plan at 50% off based on the latest discounted price. This benefit will remain available from April 30, 2026 until 3 months after the end of the user’s current legacy plan period.
During this period:

* The discounted price will dynamically reflect **50%** of the latest discounted price available at the time of purchase.

![Description](https://cdn.bigmodel.cn/markdown/1777461936885image.png?attname=image.png)

* Users may purchase or renew multiple times within the eligibility period, and renewals will continue to be charged at the discounted rate.
  After the eligibility period ends, this migration support will no longer apply.
  The available renewal prices are illustrated below (based on current pricing, for reference only):

Monthly Plan

| Plan | Current Standard Price | 50% off (example based on current pricing) |
| :- | :- | :- |
| Lite | \$18 | \$9 |
| Pro | \$72 | \$36 |
| Max | \$160 | \$80 |

Quarterly Plan

| Plan | Current Discounted Price | 50% off (example based on current pricing) |
| :- | :- | :- |
| Lite | \$48.60 | \$24.30 |
| Pro | \$194.40 | \$97.20 |
| Max | \$432 | \$216 |

Annual Plan

| Plan | Current Discounted Price | 50% off (example based on current pricing) |
| :- | :- | :- |
| Lite | \$172.80 | \$96 |
| Pro | \$691.20 | \$288 |
| Max | \$1,536 | \$768 |

### 3. What this means for you

* Your current legacy subscription will remain active until the end of the current billing cycle.
* Auto-renew for the legacy plan will be canceled on **April 30, 2026 (Singapore Time)**.
* You will automatically receive 2 complimentary months of the equivalent online version plan tier, which will begin after your current subscription ends.
* As part of the migration support, you will remain eligible for an additional 50% discount from April 30, 2026 until 3 months after the end of the user’s current legacy plan period.

### 4. FAQ

**Q: Will my current legacy plan stop immediately once auto-renew is canceled?**

No. Your current subscription will remain available through the end of the existing billing period.

**Q: Do I need to manually claim the complimentary 2 months?**

No. The migration support will be granted automatically on April 30, 2026 (Singapore Time) and added after your current plan expires.

**Q: Will the complimentary online version plan be the same tier as my current legacy plan?**

Yes. We will grant 2 months of the corresponding online version plan based on the tier of your current legacy plan.

**Q: Can the value of the 2 complimentary months be used as credit when subscribing to a new plan?**

Under the current subscription mechanism, plan upgrades are only available when there is a single active personal plan on the account.

While your account holds both your current legacy plan and the complimentary plan, you will not be able to perform an upgrade. Once the complimentary plan has taken effect and becomes the only active plan on your account, upgrades will become available again, and the offset will be handled according to the applicable rules.


# Usage Policy
Source: https://docs.z.ai/devpack/usage-policy



The GLM Coding Plan is your AI partner for tasks and code — a single subscription that unlocks a suite of agent tools for AI-powered coding and everyday office work.

## Rate Limits

Rate (concurrency) limits are tied to your plan tier. The platform dynamically adjusts these limits based on resource availability, with the general principle being **Max > Pro > Lite**.

For each project, you can use methods like Subagent to make concurrent model calls. Our recommended number of concurrent projects is as follows:

* **Lite**: Recommended for developing a single project at a time
* **Pro**: Recommended for developing 1–2 projects simultaneously
* **Max**: Recommended for developing 2 or more projects simultaneously

> [!TIP]
>
> Plan users will enjoy higher concurrency limits during off-peak hours (dynamically increased), enabling support for a greater number of concurrent projects.

## Account Usage Policy

* **Subscription benefits are exclusive to the subscriber**: Account sharing or multi-user access is prohibited. Violations may result in restrictions on subscription benefits and, in serious cases, may affect normal account usage.
* **Use limited to supported tools**: GLM Coding Plan may only be used within [officially supported tools and products](https://docs.z.ai/devpack/tool/others#step-1-supported-tools). Use in unsupported tools may result in restricted benefits.
* **Violation handling**: Violations of the [Usage Rules](https://docs.z.ai/legal-agreement/subscription-terms#4-usage-rules) may trigger risk control measures, including rate limiting, account freezing, or other restrictions. Accounts with more than three violations may be banned. If your account is flagged by the risk control system, you can view the risk notice on the [Plan Overview](https://z.ai/manage-apikey/coding-plan/personal/my-plan) page in the console and submit an appeal.

## Manage Your Subscription

You can manage your subscription, view billing details, and cancel the subscription as follows:

1. Log in to the Z.ai [API Platform](https://z.ai/subscribe?utm_source=zai\&utm_medium=index\&utm_term=glm-coding-plan\&utm_campaign=Platform_Ops&_channel_track_key=6lShUDnv).
2. Click your profile icon in the top-right corner → [Payment Method](https://z.ai/manage-apikey/billing).
3. In the left menu, select [Subscription](https://z.ai/manage-apikey/subscription).
4. To view billing history, go to Billing → [Billing History](https://z.ai/manage-apikey/billing).

## Subscription Period & Renewal

* Your subscription will **automatically renew at the end of each billing cycle**, and the fee will be charged to your linked payment method.

* The system will charge you in the following order:
  1. Priority is given to the bonus balance in your platform account.
  2. If the bonus balance is insufficient, the cash balance in your platform account will be used.
  3. If the above balances are insufficient, the payment will be deducted from your linked third-party payment method (e.g., bank card or PayPal).

* If you need to cancel your subscription, you can do so on the [Subscription](https://z.ai/manage-apikey/subscription) page. Please be sure to cancel **at least 3 days** before the next billing date to avoid automatic renewal. After cancellation, the current subscription period remains valid and will not be renewed upon expiration.

## Refund Policy

Please note that once a subscription service is purchased, it is considered confirmed, and refunds are not supported. Even if you have not used up your plan, the fees cannot be refunded. We recommend that you choose the appropriate subscription plan and term based on your usage needs.


# GLM Slide/Poster Agent(beta)
Source: https://docs.z.ai/guides/agents/slide

A Slide & Poster Agent powered by the native capabilities of the GLM model — integrating information retrieval, content structuring, and visual layout design, enabling you to effortlessly create professional-grade slides and posters.

## Overview

GLM Slide/Poster Agent is an intelligent creation agent built for working people and creators. Powered by the deeply optimized GLM LLM, it goes beyond traditional engineering-style assembly tools—supporting one-click generation of slides or posters from natural language instructions. By natively integrating content generation with layout aesthetics and design conventions, it helps you quickly produce polished, professional-grade materials while lowering design barriers and boosting creative efficiency.

## Resources

* [API Documentation](https://docs.z.ai/api-reference/agents/agent): Learn how to call the API.

## Examples

**Market Analysis Report**

> **Prompt**
>
> Help an AI startup team create a big model market analysis PPT

> **Display**
>
> ![slide](https://mintcdn.com/zhipu-32152247/fQm1SxNtD2jBDQ3i/images/slide-1.png?fit=max\&auto=format\&n=fQm1SxNtD2jBDQ3i\&q=85\&s=cbb392e136bea0d648427cdb0a316867)

**Product Launch**

> **Prompt**
>
> Make a poster for the release of my new visual understanding model

> **Display**
>
> ![slide](https://mintcdn.com/zhipu-32152247/fQm1SxNtD2jBDQ3i/images/slide-2.png?fit=max\&auto=format\&n=fQm1SxNtD2jBDQ3i\&q=85\&s=0c1398a32c5d63840a57b39a1fd30436)

## Core Features

* **Smart Information Search:** Automatically retrieve and organize relevant materials, aggregate multi-source content for richness and accuracy, and support real-time web data access to ensure up-to-date information.
* **Elegant Visual Design:** Built-in professional visual standards with intelligent layout algorithms, ensuring clean, balanced, and visually impressive slides or posters.
* **Rich Content Display:** Seamless integration of text, charts, and multimedia elements—delivering data-rich, visually engaging content that communicates at a glance.
* **Efficient Creation Workflow:** Generates polished slides or posters , supporting custom page count to meet your exact content needs.

## Usage

| Scenarios | Typical Needs |
| :- | :- |
| Workplace Reporting | Weekly Reports, Project Roadshows, Business Presentations, Internal Briefings, etc. |
| Education & Training | Teaching Courseware, Academic Thesis Presentations, Research Posters, Workshop Materials, Training Program Slides, etc. |
| Personal Use | Summary Reporting, Activity Planning Programs, Personal Portfolios, Knowledge Sharing Presentations, Event Invitations, etc. |

## Use GLM Slide/Poster Agent

Describe your slide or poster needs in one sentence, and GLM Slide/Poster Agent will generate polished, professional-level content for you!

1. **Describe Your Needs:** Enter the topic or purpose of your slide/poster in natural language
2. **Smart Information Gathering:** The intelligent agent automatically searches and organizes relevant content
3. **Slide/Poster Generation:** Instantly create visually engaging slides or posters based on curated information
4. **Refinement & Adjustment:** Modify and iterate based on your feedback (*in-browser editing coming soon*)
5. **Export and Use:** Export as PDF (*PPTX and other formats coming soon*)

## Price

* **Pay-as-you-go based on Token consumption, \$0.7/M Tokens**
* Measurement Scope: Total Tokens generated by all task flow nodes of the intelligent agent


# Translation Agent
Source: https://docs.z.ai/guides/agents/translation



## Overview

An AI-powered translation agent designed for multilingual communication and content localization, supporting automatic language detection and translation across 40+ languages. It offers six flexible translation strategies—General, Paraphrasing, Two-Step, Three-Stage, Reflection, and COT—to handle diverse scenarios such as social media, literary works, and academic papers. Users can further enhance translation accuracy through terminology customization and expert guidance. This translation agent also provides a corresponding API service, making it a powerful tool for individuals, enterprises, developers, and professional translators.

## Core Features

* **Multilingual Translation Support:** Covers translation and language detection for 40 languages, including niche languages such as Classical Chinese and Cantonese.
* **Six Professional Translation Strategies:** Includes General, Paraphrasing, Two-Step, Three-Stage, Reflective, and COT translation modes to meet diverse text translation needs.
* **Terminology Customization:** Enterprises or users can upload terminology glossaries for term alignment, enhancing domain-specific accuracy.
* **Expert Guidance and Style Adaptation:** Allows specification of stylistic preferences and cultural context to ensure translations better fit the intended usage scenarios.
* **Pinyin and IPA Support:** Converts text into Pinyin or International Phonetic Alphabet, suitable for language learning or voice-assistive content creation.
* **Intra-Language “Translation” Support:** When the source and target languages are the same, applying strategies like Two-Step, Reflective, or COT will paraphrase, rewrite, or polish the text—helpful for eliminating ambiguity in secondary translations.
* **Translation Explainability:** COT and Reflective strategies can output the reasoning behind the translation process.

| **Feature** | **General-Purpose Translation Agent** | **Traditional Translation API** | **General LLM** |
| :- | :- | :- | :- |
| Multilingual Support | ✅ 40 languages | ✅ Typically supports many | ✅ Typically supports many |
| Professional Strategies | ✅ 6 specialized strategies | ❌ Limited strategies | ❌ No dedicated translation strategies |
| Terminology Support | ✅ Full support | ⚠️ Partial support | ❌ Not supported |
| Translation Suggestions | ✅ Detailed suggestions supported | ❌ Not supported | ❌ Requires manual input |
| Translation Explainability | ✅ Provides reasoning | ❌ Not available | ⚠️ Requires special prompting |
| Auto Language Detection | ✅ High accuracy | ✅ Supported | ✅ Supported |
| Special Language Support | ✅ Classical Chinese, Cantonese, etc. | ❌ Typically unsupported | ⚠️ Limited support |

## Usage

| **Target Users** | **Application Scenarios** |
| :- | :- |
| Individual Developers & Language Learners | Provides translation assistance and phonetic annotation to enhance learning efficiency |
| Content Creators / New Media Operators | Supports multilingual content publishing and social media localization |
| Enterprise Clients / Translation Providers | Handles batch translation tasks, manages custom terminology, and ensures compliant output |
| Educational / Research Institutions | Supports translation of papers, course materials, and academic communication |

| **Strategy** | **Description** | **Applicable Scenarios** | **Key Features** |
| :- | :- | :- | :- |
| General Translation | A basic translation method that retains the original format, accounts for the target language’s cultural and linguistic context, balances accuracy and fluency, and resolves over 95% of prompt injection issues. | General-purpose content translation; suitable for most daily use cases and handling prompt injection. | - Preserves original structure - Follows translation rules  - Supports terminology glossaries - Adapts to regional language variants |
| Paraphrased Translation | Respects the original meaning while rewriting content in the target language’s natural expression. | When natural adaptation to target culture is preferred; non-literal translation scenarios. | -  Focuses on conveying meaning over form - Adapts more naturally to the target language - Allows greater flexibility in restructuring |
| Two-Step Translation | Involves literal translation first, followed by freer expression; combines direct and adaptive approaches. | Literary works or content needing both accuracy and expressiveness. | - Balances literal and intended meanings - Enhances comprehension - Stepwise approach improves both accuracy and fluency |
| Three-Stage Translation | Based on the Chinese traditional theory of "faithfulness, expressiveness, elegance"; ensures content accuracy, fluent delivery, and stylistic beauty—currently supports classical/literary styles only. | High-quality literary, poetic, or culturally expressive texts. | - Follows "faithfulness, expressiveness, elegance" principle - Focuses on cultural and philosophical nuance - Produces refined, elegant output |
| Reflective Translation | Starts with a literal translation, then prompts the model to act as an expert to reflect on the result by specific dimensions, and finally optimizes the output accordingly. | Professional publications, formal documents, or content requiring high accuracy and quality. | - Improves output through iterative feedback - Corrects errors and enhances style - Ideal for complex or domain-specific texts |
| COT Translation | Applies a Chain of Thought reasoning process to analyze the source text before translating into the target language. | Conceptually complex or professional domain content requiring deeper understanding. | - Performs explicit reasoning before translation - Enhances comprehension of complex input - Increases transparency - Useful in expert domains |

## Use Translation Agent

1. **Select Source Language:** Specify the source language of the text to be translated (can be left as auto-detect by default).
2. **Select Target Language:** Choose the language you want the text to be translated into.
3. **Choose Translation Strategy:** Select a suitable translation strategy based on your needs (e.g., General, Paraphrasing, Two-Step, etc.).
4. **Enter Source Text:** Paste the text you want to translate into the input box.
5. **Get Translation Result:** The agent will generate a high-quality translation based on the selected strategy.

## Examples

> Original English Text：
>
> "Two roads diverged in a wood, and I took the one less traveled by, and that has made all the difference."
>
> — Robert Frost

1. General Translation：

> “林中有两条路分叉，而我选择了那条少有人走的路，而这让一切变得不同。”
>
> ——罗伯特·弗罗斯特

Uses a basic literal translation approach with terminology support, preserving the original format while accommodating the cultural context and linguistic norms of the target language—**suitable for over 95% of general translation needs.**

2. Paraphrased Translation：

> “树林中分出两条路，我选择了人迹罕至的一条，而这决定了我一生的道路。”
>
> ——罗伯特·弗罗斯特

While staying true to the original meaning, this approach rewrites sentence structure and expressions to present the content in a more natural and native way in the target language. **It is suitable for translation scenarios that do not require word-for-word correspondence and prioritize cultural adaptation and linguistic conventions.**

3. Two-Step Translation：

> "林中有两条路分岔，而我选择了那条少有人走的路，这一切因此而截然不同。"
>
> ——罗伯特·弗罗斯特

This approach first performs a literal word-for-word translation, then refines the result through free translation to improve overall quality and readability. **It is well-suited for scenarios like literary works, where semantic understanding and stylistic expression are especially important.**

4. Three-Principle Translation：

> "林间歧路分，吾择幽径以行，此别天地宽。"
>
> ——罗伯特·弗罗斯特

Based on the traditional Chinese translation principles of "faithfulness, expressiveness, and elegance", this approach ensures that the translation is accurate, fluent, and stylistically refined. **It is ideal for literary content such as poetry and prose**, supporting high-register styles like Classical Chinese, and aims to balance form and meaning.

5. Reflective Translation：

> "林中两条路岔开，而我选择了那条少有人走的，而这，让一切截然不同。"
>
> ——罗伯特·弗罗斯特

Building on an initial literal translation, this approach introduces an expert role to reflect on and provide feedback, followed by refinement and optimization of the output. It enhances accuracy and stylistic consistency, making **it suitable for high-quality formal texts** such as publications, legal documents, or professional materials.

6. COT Translation：

```python theme={null}
<think>
First, we need to understand the basic meaning of the poem. Robert Frost’s poem describes a person faced with a fork in the woods, symbolizing a moment of choice. The “two roads” represent different paths in life, and the speaker chooses “the one less traveled by,” indicating a preference for uniqueness and individuality. This choice, which “has made all the difference,” suggests that it had a significant impact on the speaker’s life.
When translating, we need to pay attention to the following points:
1. “Two roads diverged in a wood”中的“diverged”表示分叉，可以翻译为“分叉”或“岔开”。
2. “I took the one less traveled by”中的“less traveled”表示少有人走，可以翻译为“少有人走的”。
3. “And that has made all the difference”中的“made all the difference”表示造成了巨大的影响，可以翻译为“造成了所有的不同”或“这一切都不同了”。
In addition, to preserve the poem’s rhythm and aesthetic appeal, we should strive to make the translation concise and poetic.
</think>
```

> 树林中两条路分叉，而我选择了那条少有人走的，这一切都不同了。
>
> —— 罗伯特·弗罗斯特

By following a “first understand, then translate” reasoning process, this approach enhances the model’s ability to comprehend complex concepts and specialized contexts. **It is well-suited for high-expertise domains such as science and medicine**, resulting in more logical and transparent translations that help users grasp the underlying knowledge structure behind the language.

## Price

Pay-as-you-go by Token usage, \$3 per 1 million Tokens


# Video Effect Template Agent
Source: https://docs.z.ai/guides/agents/video-template



## Overview

Includes three popular special effects video templates: french\_kiss, bodyshake, and sexy\_me. Users can upload a single image to generate a professional-grade special effects video with a single click. The app is easy to use, produces smooth images, and effectively improves the efficiency and quality of short video creation.

## Examples

<table>
  <tr>
    <th>
      Templates
    </th>

    <th>
      Prompt
    </th>

    <th>
      Example Input
    </th>

    <th>
      Example Output
    </th>
  </tr>

  <tr>
    <td>
      french\_kiss
    </td>

    <td>
      The two figures in the painting gradually draw closer, then passionately kiss, alternating between deep and intense moments.
    </td>

    <td>
      ![](https://cdn.bigmodel.cn/static/platform/agent/french-kiss-template.png)
    </td>

    <td>
      <video />
    </td>
  </tr>

  <tr>
    <td>
      bodyshake
    </td>

    <td>
      Video content: The character performs a rhythmic dance sequence in an indoor setting. She first sways her hips, then turns to the other side, briefly shaking her hips in a playful manner. Her movements are fluid and confident, consistently emphasizing body rhythm and expressiveness.  Requirements: Movement Level: High
    </td>

    <td>
      ![](https://cdn.bigmodel.cn/static/platform/agent/dance-template.png)
    </td>

    <td>
      <video />
    </td>
  </tr>

  <tr>
    <td>
      sexy\_me
    </td>

    <td>
      The woman's attire undergoes a seamless transformation, her original clothing smoothly transitioning into a fashionable bikini. In the final moment, she confidently places her hands on her hips, exuding elegance and poise.
    </td>

    <td>
      ![](https://cdn.bigmodel.cn/static/platform/agent/transformation-template.png)
    </td>

    <td>
      <video />
    </td>
  </tr>
</table>

## Core Features

* **Excellent Dynamic Performance**：Natural and smooth movements with a wide range of motion, eliminating the sluggishness commonly found in AI videos, and suitable for a variety of creative special effects scenarios.
* **Accurate Semantic Understanding**: Able to efficiently generate a variety of images based on prompts, providing a smoother card-drawing experience.
* **Excellent Anime Style**: Stable video generation quality without sudden changes.
* **Strong Subject Consistency**:No need for additional first frame images; generate highly consistent videos with a single click, greatly simplifying the creation process.

## Usage

| **Target Users** | **Application Scenarios** |
| :- | :- |
| Individual Creators | Quickly produce popular special effects short videos (such as costume changes/dance challenges) and efficiently produce creative content at low cost. |
| MCN Agencies | Batch generate standardized special effects videos (such as popular costume change templates) to meet the large-scale content needs of matrix accounts. |
| Short Video Platforms | Provide an integrated special effects template library to lower user creation barriers, enhance platform content diversity, and boost user engagement. |

## Use Video Templates

1. **French Kiss Template:**

* Image Quantity: Required; Only 1 image can be uploaded.
* Number of People in Image: Only supports two-person collages or group photos.
* Guidelines: Subjects should be shown from the front and upper body only, with no props in hand for best results.
* Prompt: The two figures in the image gradually move closer, then passionately kiss with alternating deep and firm intensity.
* Effect Limitations: If a full-body image is uploaded, the result may show no kiss, a weaker kiss, or a shorter duration—failing to convey the intended passionate effect.

2. **Bodyshake Template:**

* Image Quantity: Required; Only 1 image can be uploaded.
* Number of People in Image: Supports single-person photos in either realistic or anime style.
* Guidelines: Best results are achieved when the character is alone and the image shows at least above-thigh framing.
* Prompt: "Video content: The character performs a rhythmic dance sequence in an indoor setting. She starts by swaying her hips, then turns to the other side, briefly shaking her hips in a playful manner. Her movements are smooth and confident, consistently emphasizing rhythm and expressiveness. Requirement: Movement intensity – high."
* Effect Limitations: During the turning motion, body structure distortion or hand clipping may occasionally occur. For hip-shaking actions, motion may sometimes appear disconnected or visually unnatural.

3. **Sexy\_Me Template：**

* Image Quantity: Required; Only 1 image can be uploaded.
* Number of People in Image: Only supports single-person photos (female or male) in either realistic or anime style.
* Guidelines: Best results are achieved with half-body or full-body images. If the subject in the input image is female, the output will feature a bikini transformation video showcasing the body. If the subject is male, the output will show a shirtless transformation video highlighting muscular physique.
* Prompt: "Video content: The transformation varies depending on the subject's gender. If the image shows a female: 'The woman's clothing transforms seamlessly—her outfit smoothly changes into a stylish bikini. At the final moment, she confidently places her hands on her waist, exuding grace and poise.' If the image shows a male: 'The man swiftly removes his shirt, revealing a muscular physique matching his skin tone. He then steps forward.' Requirements: 1. If the image is a close-up or medium shot, set the camera motion to "zoom out." 2. Movement intensity: high."
* Effect Limitations: If the subject's gender features are unclear, incorrect gender identification may occur, resulting in mixed or incorrect bikini/muscle transformations. In bikini transformations, clothing may not be fully removed. If a two-person image is uploaded, one subject may be missing in the output.

## Price

Pay-as-you-go based on number of videos, \$ 0.2 per video


# Context Caching
Source: https://docs.z.ai/guides/capabilities/cache



Context caching functionality significantly reduces token consumption and response latency by caching repeated context content. When you repeatedly use the same system prompts or conversation history in dialogues, the caching mechanism automatically identifies and reuses this content, thereby improving performance and reducing costs.

## Features

* **Automatic Cache Recognition**: Implicit caching that intelligently identifies repeated context content without manual configuration
* **Significant Cost Reduction**: Cached tokens are billed at lower prices, dramatically saving costs
* **Improved Response Speed**: Reduces processing time for repeated content, accelerating model responses
* **Transparent Billing**: Detailed display of cached token counts in response field `usage.prompt_tokens_details.cached_tokens`
* **Wide Compatibility**: Supports all mainstream models, including GLM-5, GLM-4.7, GLM-4.6 series, etc.

> Context caching works by computing input message content and identifying content that is identical or highly similar to previous requests. When repeated content is detected, the system reuses previous computation results, avoiding redundant token processing.

This mechanism is particularly suitable for the following scenarios:

* System prompt reuse: In multi-turn conversations, system prompts usually remain unchanged, and caching can significantly reduce token consumption for this part.
* Repetitive tasks: For tasks that process similar content with consistent instructions multiple times, caching can improve efficiency.
* Multi-turn conversation history: In complex conversations, historical messages often contain a lot of repeated information, and caching can effectively reduce token usage for this part.

## Code Examples

**cURL**

**Basic Caching Example**

```bash theme={null}
# First request - establish cache
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "system",
            "content": "You are a professional data analyst, skilled at explaining data trends and providing business insights."
        },
        {
            "role": "user",
            "content": "How to analyze user retention rate?"
        }
    ]
}'
```

**Cache Reuse Example**

```bash theme={null}
# Second request - reuse system prompt cache
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "system",
            "content": "You are a professional data analyst, skilled at explaining data trends and providing business insights."
        },
        {
            "role": "user",
            "content": "What is funnel analysis?"
        }
    ]
}'
```

**Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Basic Conversation Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API key')

# First request - establish cache
response1 = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a professional technical documentation assistant, skilled at explaining complex technical concepts. Please answer user questions with clear and concise language, and provide practical code examples."
        },
        {
            "role": "user",
            "content": "What is RESTful API?"
        }
    ]
)

print("First request result:")
print(f"Reply: {response1.choices[0].message.content}")
print(f"Total tokens: {response1.usage.total_tokens}")
print(f"Cached tokens: {response1.usage.prompt_tokens_details.cached_tokens if hasattr(response1.usage, 'prompt_tokens_details') else 0}")

# Second request - reuse system prompt cache
response2 = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a professional technical documentation assistant, skilled at explaining complex technical concepts. Please answer user questions with clear and concise language, and provide practical code examples."  # Same system prompt
        },
        {
            "role": "user",
            "content": "What are the differences between GraphQL and RESTful API?"
        }
    ]
)

print("\nSecond request result:")
print(f"Reply: {response2.choices[0].message.content}")
print(f"Total tokens: {response2.usage.total_tokens}")
print(f"Cached tokens: {response2.usage.prompt_tokens_details.cached_tokens if hasattr(response2.usage, 'prompt_tokens_details') else 0}")
```

**Long Document Analysis Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API key')

# Long document content (simulated)
long_document = """
This is a detailed technical specification document that includes system architecture, API design, database structure, and many other aspects.
The document is very long and contains a lot of technical details and implementation instructions...
[Large amount of document content omitted here]
"""

# First analysis - establish document cache
response1 = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": f"Please answer user questions based on the following technical document:\n\n{long_document}"
        },
        {
            "role": "user",
            "content": "What is the main architecture of this system?"
        }
    ]
)

print("First analysis:")
print(f"Total tokens: {response1.usage.total_tokens}")
print(f"Cached tokens: {response1.usage.prompt_tokens_details.cached_tokens if hasattr(response1.usage, 'prompt_tokens_details') else 0}")

# Second analysis - reuse document cache
response2 = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": f"Please answer user questions based on the following technical document:\n\n{long_document}"  # Same document content
        },
        {
            "role": "user",
            "content": "What are the characteristics of the API design?"
        }
    ]
)

print("\nSecond analysis:")
print(f"Total tokens: {response2.usage.total_tokens}")
print(f"Cached tokens: {response2.usage.prompt_tokens_details.cached_tokens if hasattr(response2.usage, 'prompt_tokens_details') else 0}")
print(f"Cache savings: {response2.usage.prompt_tokens_details.cached_tokens / response2.usage.total_tokens * 100:.1f}%")
```

**Multi-turn Conversation Caching Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API key')

# Build conversation history
conversation_history = [
    {"role": "system", "content": "You are a Python programming assistant, helping users solve programming problems."},
    {"role": "user", "content": "How to create a simple Flask application?"},
    {"role": "assistant", "content": "Creating a Flask application is simple, first install Flask..."},
    {"role": "user", "content": "How to add routes?"},
    {"role": "assistant", "content": "In Flask, add routes using the @app.route decorator..."},
]

# Continue conversation - reuse conversation history cache
response = client.chat.completions.create(
    model="glm-5.3",
    messages=conversation_history + [
        {"role": "user", "content": "How to handle POST requests?"}
    ]
)

print("Conversation reply:")
print(f"Content: {response.choices[0].message.content}")
print(f"Total tokens: {response.usage.total_tokens}")
print(f"Cached tokens: {response.usage.prompt_tokens_details.cached_tokens if hasattr(response.usage, 'prompt_tokens_details') else 0}")

# Calculate cache efficiency
if hasattr(response.usage, 'prompt_tokens_details') and response.usage.prompt_tokens_details.cached_tokens:
    cache_ratio = response.usage.prompt_tokens_details.cached_tokens / response.usage.prompt_tokens * 100
    print(f"Cache hit rate: {cache_ratio:.1f}%")
```

**Batch Processing Optimization Example**

````python theme={null}
from zai import ZaiClient
import time

# Initialize client
client = ZaiClient(api_key='Your API key')

# Common system prompt
system_prompt = """
You are a professional code review assistant. Please analyze the provided code from the following aspects:
1. Code quality and readability
2. Performance optimization suggestions
3. Security considerations
4. Best practice recommendations
Please provide specific improvement suggestions.
"""

# List of code snippets to review
code_snippets = [
    "def calculate_sum(numbers): return sum(numbers)",
    "class User: def __init__(self, name): self.name = name",
    "for i in range(len(items)): print(items[i])",
    "if user_input == 'yes' or user_input == 'y': return True"
]

results = []
total_cached_tokens = 0

for i, code in enumerate(code_snippets):
    start_time = time.time()
    
    response = client.chat.completions.create(
        model="glm-5.3",
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": f"Please review the following code:\n```python\n{code}\n```"}
        ]
    )
    
    end_time = time.time()
    
    # Count cache effects
    cached_tokens = 0
    if hasattr(response.usage, 'prompt_tokens_details') and response.usage.prompt_tokens_details.cached_tokens:
        cached_tokens = response.usage.prompt_tokens_details.cached_tokens
        total_cached_tokens += cached_tokens
    
    results.append({
        'code': code,
        'review': response.choices[0].message.content,
        'total_tokens': response.usage.total_tokens,
        'cached_tokens': cached_tokens,
        'response_time': end_time - start_time
    })
    
    print(f"Code snippet {i+1} review completed:")
    print(f"  Response time: {end_time - start_time:.2f}s")
    print(f"  Cached tokens: {cached_tokens}")
    print(f"  Total tokens: {response.usage.total_tokens}")
    print()

print(f"Batch processing completed, total cached tokens: {total_cached_tokens}")
````

Response contains context cache token usage information:

```json theme={null}
{
  "usage": {
    "prompt_tokens": 1200,
    "completion_tokens": 300,
    "total_tokens": 1500,
    "prompt_tokens_details": {
      "cached_tokens": 800
    }
  }
}
```

## Best Practices

**System Prompt Optimization**

Use stable system prompts

```python theme={null}
# Recommended: Use stable system prompts
system_prompt = """
You are a professional technical consultant with the following characteristics:
- Deep technical background and rich project experience
- Able to provide accurate and practical technical advice
- Good at explaining complex concepts in clear and concise language
Please provide professional technical guidance based on user questions.
"""
```

**Document Content Reuse**

Use long documents as system messages

```python theme={null}
# Recommended: Use long documents as system messages
def create_document_based_chat(document_content, user_question):
    return client.chat.completions.create(
        model="glm-5.3",
        messages=[
            {
                "role": "system",
                "content": f"Please answer user questions based on the following document content:\n\n{document_content}"
            },
            {
                "role": "user",
                "content": user_question
            }
        ]
    )

# Multiple calls with the same document, system prompts will be cached
questions = ["What is the main content of the document?", "What are the key points?", "How to implement these suggestions?"]
for question in questions:
    response = create_document_based_chat(document_content, question)
    # Second and subsequent calls will hit the cache
```

**Conversation History Management**

Manage conversation history to improve cache efficiency

```python theme={null}
class ConversationManager:
    def __init__(self, client, system_prompt):
        self.client = client
        self.system_prompt = system_prompt
        self.history = [{"role": "system", "content": system_prompt}]
    
    def add_message(self, role, content):
        self.history.append({"role": role, "content": content})
    
    def get_response(self, user_message):
        # Add user message
        self.add_message("user", user_message)
        
        # Get reply (conversation history will be cached)
        response = self.client.chat.completions.create(
            model="glm-5.3",
            messages=self.history
        )
        
        # Add assistant reply to history
        assistant_message = response.choices[0].message.content
        self.add_message("assistant", assistant_message)
        
        return response
    
    def get_cache_stats(self, response):
        """Get cache statistics"""
        if hasattr(response.usage, 'prompt_tokens_details'):
            cached = response.usage.prompt_tokens_details.cached_tokens or 0
            total = response.usage.prompt_tokens
            return f"Cache hit: {cached}/{total} ({cached/total*100:.1f}%)"
        return "No cache information"

# Usage example
manager = ConversationManager(client, "You are a programming assistant...")
response1 = manager.get_response("How to learn Python?")
response2 = manager.get_response("Recommend some learning resources")  # Will reuse previous conversation cache
```

## Use Cases

> **Multi-turn Conversations**
>
> * Intelligent customer service systems
> * Personal assistant services

> **Batch Processing**
>
> * Code review batch processing
> * Content batch analysis

> **Template Applications**
>
> * Report generation templates
> * Standardized process handling

> **Education and Training**
>
> * Homework grading assistance
> * Learning material analysis

## Important Notes

**Understanding Cache Mechanism**

* Caching is automatically triggered based on content similarity, no manual configuration required
* Identical content has the highest cache hit rate
* Minor formatting differences may affect cache effectiveness
* Cache has reasonable time limits, will recalculate after expiration

**Cost Optimization Suggestions**

* Cached tokens are billed at lower prices
* Long documents and repeated content have the most significant cache effects
* Design system prompts reasonably to improve reuse rates
* Monitor cache hit rates and optimize usage patterns

**Performance Considerations**

* Caching can significantly improve response speed
* First request to establish cache may be slightly slower
* Manage conversation history length reasonably
* Avoid overly frequent content changes

**Best Practices**

* Use stable system prompt templates
* Process long documents as system messages
* Organize conversation history structure reasonably
* Regularly analyze cache effectiveness and optimize

## Billing Information

Context caching uses a differentiated billing strategy:

* New content tokens: Billed at standard prices
* Cache hit tokens: Billed at discounted prices (usually 50% of standard price)
* Output tokens: Billed at standard prices

Billing example:

```
Assuming standard price is 0.01 /1K tokens:

Request details:
- Total input tokens: 2000
- Cache hit tokens: 1200
- New content tokens: 800
- Output tokens: 500

Billing calculation:
- New content cost: 800 × 0.01/1000 = 0.008
- Cache cost: 1200 × 0.005/1000 = 0.006
- Output cost: 500 × 0.01/1000 = 0.005
- Total cost: 0.019

Compared to no cache (2500 × 0.01/1000 = 0.025), saves 24%
```


# Function Calling
Source: https://docs.z.ai/guides/capabilities/function-calling



Function Calling allows AI models to call external functions and APIs, greatly expanding the capability boundaries of intelligent agents, enabling them to perform specific operations and obtain real-time data.

## Features

Function calling provides AI models with the ability to interact with external systems, supporting various complex application scenarios and integration requirements.

### Core Parameter Description

* **`tools`**: Defines the list of callable functions, including function names, descriptions, and parameter specifications
* **`tool_choice`**: Controls function calling strategy, default is `auto` (only supports `auto`)

### Response Parameter Description

Key fields in function calling responses:

* **`tool_calls`**: Contains information about functions the model decides to call
* **`function.name`**: Name of the called function
* **`function.arguments`**: Function call parameters (JSON format string)
* **`id`**: Unique identifier for the tool call

## Code Examples

By defining function tools and handling function calls, AI models can perform various external operations:

**Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

```python theme={null}
import json
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='your_api_key')

# Define weather query function
def get_weather(city: str) -> dict:
    """Get weather information for specified city"""
    # This should call a real weather API
    weather_data = {
        "city": city,
        "temperature": "22°C",
        "condition": "Sunny",
        "humidity": "65%",
        "wind_speed": "5 km/h"
    }
    return weather_data

# Define function tools
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "Get current weather information for specified city",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {
                        "type": "string",
                        "description": "City name, e.g.: Beijing, Shanghai"
                    }
                },
                "required": ["city"]
            }
        }
    }
]

# Make conversation request
response = client.chat.completions.create(
    model="glm-5.3",  # Use model that supports function calling
    messages=[
        {"role": "user", "content": "How's the weather in Beijing today?"}
    ],
    tools=tools,         # Pass function tools
    tool_choice="auto"   # Automatically choose whether to call functions
)

# Handle function calls
message = response.choices[0].message
messages = [{"role": "user", "content": "How's the weather in Beijing today?"}]
messages.append(message.model_dump())

if message.tool_calls:
    for tool_call in message.tool_calls:
        if tool_call.function.name == "get_weather":
            # Parse parameters and call function
            args = json.loads(tool_call.function.arguments)
            weather_result = get_weather(args.get("city"))
            
            # Return function result to model
            messages.append({
                "role": "tool",
                "content": json.dumps(weather_result, ensure_ascii=False),
                "tool_call_id": tool_call.id
            })
    
    # Get final answer
    final_response = client.chat.completions.create(
        model="glm-5.3",
        messages=messages,
        tools=tools
    )
    
    print(final_response.choices[0].message.content)
else:
    print(message.content)
```

## Scenario Examples

> [!WARNING]
>
> When using function calling, please ensure proper security validation and permission control for external APIs and database operations.

> **Multi-function Assistant**
>
> ```python theme={null}
> import json
> import requests
> from datetime import datetime
> from zai import ZaiClient
>
> class FunctionAgent:
>     def __init__(self, api_key):
>         self.client = ZaiClient(api_key=api_key)
>         self.tools = self._define_tools()
>
>     def _define_tools(self):
>         return [
>             {
>                 "type": "function",
>                 "function": {
>                     "name": "get_current_time",
>                     "description": "Get current time",
>                     "parameters": {
>                         "type": "object",
>                         "properties": {},
>                         "required": []
>                     }
>                 }
>             },
>             {
>                 "type": "function",
>                 "function": {
>                     "name": "calculate",
>                     "description": "Perform mathematical calculations",
>                     "parameters": {
>                         "type": "object",
>                         "properties": {
>                             "expression": {
>                                 "type": "string",
>                                 "description": "Mathematical expression, e.g.: 2+3*4"
>                             }
>                         },
>                         "required": ["expression"]
>                     }
>                 }
>             },
>             {
>                 "type": "function",
>                 "function": {
>                     "name": "search_web",
>                     "description": "Search web information",
>                     "parameters": {
>                         "type": "object",
>                         "properties": {
>                             "query": {
>                                 "type": "string",
>                                 "description": "Search keywords"
>                             }
>                         },
>                         "required": ["query"]
>                     }
>                 }
>             }
>         ]
>
>     def get_current_time(self):
>         """Get current time"""
>         return {
>             "current_time": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
>             "timezone": "Asia/Shanghai"
>         }
>
>     def calculate(self, expression: str):
>         """Safe mathematical calculation"""
>         try:
>             # Simple security check
>             allowed_chars = set('0123456789+-*/().')
>             if not all(c in allowed_chars or c.isspace() for c in expression):
>                 return {"error": "Expression contains disallowed characters"}
>
>             result = eval(expression)
>             return {
>                 "expression": expression,
>                 "result": result
>             }
>         except Exception as e:
>             return {"error": f"Calculation error: {str(e)}"}
>
>     def search_web(self, query: str):
>         """Simulate web search"""
>         # This should call a real search API
>         return {
>             "query": query,
>             "results": [
>                 {"title": f"Search result 1 about {query}", "url": "https://example1.com"},
>                 {"title": f"Search result 2 about {query}", "url": "https://example2.com"}
>             ]
>         }
>
>     def execute_function(self, function_name: str, arguments: dict):
>         """Execute function call"""
>         if function_name == "get_current_time":
>             return self.get_current_time()
>         elif function_name == "calculate":
>             return self.calculate(arguments.get("expression", ""))
>         elif function_name == "search_web":
>             return self.search_web(arguments.get("query", ""))
>         else:
>             return {"error": f"Unknown function: {function_name}"}
>
>     def chat(self, user_message: str):
>         """Handle user message"""
>         messages = [{"role": "user", "content": user_message}]
>
>         response = self.client.chat.completions.create(
>             model="glm-5.3",
>             messages=messages,
>             tools=self.tools,
>             tool_choice="auto"
>         )
>
>         message = response.choices[0].message
>         messages.append(message.model_dump())
>
>         # Handle function calls
>         if message.tool_calls:
>             for tool_call in message.tool_calls:
>                 function_name = tool_call.function.name
>                 arguments = json.loads(tool_call.function.arguments)
>
>                 # Execute function
>                 result = self.execute_function(function_name, arguments)
>
>                 # Add function result
>                 messages.append({
>                     "role": "tool",
>                     "content": json.dumps(result, ensure_ascii=False),
>                     "tool_call_id": tool_call.id
>                 })
>
>             # Get final answer
>             final_response = self.client.chat.completions.create(
>                 model="glm-5.3",
>                 messages=messages,
>                 tools=self.tools
>             )
>
>             return final_response.choices[0].message.content
>         else:
>             return message.content
>
> # Usage example
> agent = FunctionAgent("your_api_key")
>
> # Test different types of requests
> print(agent.chat("What time is it now?"))
> print(agent.chat("Help me calculate 15 * 23 + 7"))
> print(agent.chat("Search for the latest developments in artificial intelligence"))
> ```

> **Database Query**
>
> ```python theme={null}
> import sqlite3
>
> def query_database(sql: str) -> dict:
>     """Execute database query"""
>     try:
>         conn = sqlite3.connect('example.db')
>         cursor = conn.cursor()
>         cursor.execute(sql)
>         results = cursor.fetchall()
>         conn.close()
>
>         return {
>             "success": True,
>             "data": results,
>             "row_count": len(results)
>         }
>     except Exception as e:
>         return {
>             "success": False,
>             "error": str(e)
>         }
>
> # Function definition
> db_tool = {
>     "type": "function",
>     "function": {
>         "name": "query_database",
>         "description": "Execute SQL query",
>         "parameters": {
>             "type": "object",
>             "properties": {
>                 "sql": {
>                     "type": "string",
>                     "description": "SQL query statement"
>                 }
>             },
>             "required": ["sql"]
>         }
>     }
> }
> ```

> **File Operations**
>
> ```python theme={null}
> import os
> import json
>
> def file_operations(operation: str, file_path: str, content: str = None) -> dict:
>     """File operation function"""
>     try:
>         if operation == "read":
>             with open(file_path, 'r', encoding='utf-8') as f:
>                 content = f.read()
>             return {"success": True, "content": content}
>
>         elif operation == "write":
>             with open(file_path, 'w', encoding='utf-8') as f:
>                 f.write(content)
>             return {"success": True, "message": "File written successfully"}
>
>         elif operation == "list":
>             files = os.listdir(file_path)
>             return {"success": True, "files": files}
>
>         else:
>             return {"success": False, "error": "Unsupported operation"}
>
>     except Exception as e:
>         return {"success": False, "error": str(e)}
>
> # Function definition
> file_tool = {
>     "type": "function",
>     "function": {
>         "name": "file_operations",
>         "description": "Execute file operations",
>         "parameters": {
>             "type": "object",
>             "properties": {
>                 "operation": {
>                     "type": "string",
>                     "enum": ["read", "write", "list"],
>                     "description": "Operation type"
>                 },
>                 "file_path": {
>                     "type": "string",
>                     "description": "File path"
>                 },
>                 "content": {
>                     "type": "string",
>                     "description": "Content to write (only required for write operation)"
>                 }
>             },
>             "required": ["operation", "file_path"]
>         }
>     }
> }
> ```

> **API Integration**
>
> ```python theme={null}
> import requests
>
> def call_external_api(url: str, method: str = "GET", headers: dict = None, data: dict = None) -> dict:
>     """Call external API"""
>     try:
>         if method.upper() == "GET":
>             response = requests.get(url, headers=headers, params=data)
>         elif method.upper() == "POST":
>             response = requests.post(url, headers=headers, json=data)
>         else:
>             return {"success": False, "error": "Unsupported HTTP method"}
>
>         return {
>             "success": True,
>             "status_code": response.status_code,
>             "data": response.json() if response.headers.get('content-type', '').startswith('application/json') else response.text
>         }
>
>     except Exception as e:
>         return {"success": False, "error": str(e)}
>
> # Function definition
> api_tool = {
>     "type": "function",
>     "function": {
>         "name": "call_external_api",
>         "description": "Call external API",
>         "parameters": {
>             "type": "object",
>             "properties": {
>                 "url": {
>                     "type": "string",
>                     "description": "API endpoint URL"
>                 },
>                 "method": {
>                     "type": "string",
>                     "enum": ["GET", "POST"],
>                     "description": "HTTP method"
>                 },
>                 "headers": {
>                     "type": "object",
>                     "description": "Request headers"
>                 },
>                 "data": {
>                     "type": "object",
>                     "description": "Request data"
>                 }
>             },
>             "required": ["url"]
>         }
>     }
> }
> ```

## Best Practices

> **Function Design Principles**
>
> * Single responsibility: Each function should do one thing
> * Clear naming: Function and parameter names should be meaningful
> * Complete description: Provide detailed function and parameter descriptions

> **Security Considerations**
>
> * Input validation: Strictly validate all input parameters
> * Permission control: Limit function access permissions
> * Logging: Record function call logs

### Parameter Design

```python theme={null}
# Good parameter design
{
    "type": "object",
    "properties": {
        "city": {
            "type": "string",
            "description": "City name, supports Chinese and English, e.g.: Beijing, Shanghai, New York",
            "examples": ["Beijing", "Shanghai", "New York"]
        },
        "unit": {
            "type": "string",
            "enum": ["celsius", "fahrenheit"],
            "description": "Temperature unit",
            "default": "celsius"
        }
    },
    "required": ["city"]
}
```

### Error Handling

```python theme={null}
def robust_function(param: str) -> dict:
    """Robust function implementation"""
    try:
        # Parameter validation
        if not param or not isinstance(param, str):
            return {
                "success": False,
                "error": "Invalid parameter",
                "error_code": "INVALID_PARAM"
            }
        
        # Business logic
        result = process_data(param)
        
        return {
            "success": True,
            "data": result,
            "timestamp": datetime.now().isoformat()
        }
    
    except ValueError as e:
        return {
            "success": False,
            "error": f"Data error: {str(e)}",
            "error_code": "DATA_ERROR"
        }
    except Exception as e:
        return {
            "success": False,
            "error": f"System error: {str(e)}",
            "error_code": "SYSTEM_ERROR"
        }
```

### Input Validation

```python theme={null}
def secure_function(user_input: str) -> dict:
    """Secure function implementation"""
    # Input length limit
    if len(user_input) > 1000:
        return {"error": "Input too long"}
    
    # Dangerous character filtering
    dangerous_chars = ['<', '>', '&', '"', "'"]
    if any(char in user_input for char in dangerous_chars):
        return {"error": "Input contains dangerous characters"}
    
    # SQL injection protection
    sql_keywords = ['DROP', 'DELETE', 'UPDATE', 'INSERT']
    if any(keyword in user_input.upper() for keyword in sql_keywords):
        return {"error": "Input contains dangerous keywords"}
    
    return {"success": True, "processed_input": user_input}
```

### Permission Control

```python theme={null}
def check_permissions(user_id: str, operation: str) -> bool:
    """Check user permissions"""
    user_permissions = get_user_permissions(user_id)
    return operation in user_permissions

def protected_function(user_id: str, operation: str, data: dict) -> dict:
    """Function requiring permission validation"""
    if not check_permissions(user_id, operation):
        return {
            "success": False,
            "error": "Insufficient permissions",
            "error_code": "PERMISSION_DENIED"
        }
    
    # Execute operation
    return perform_operation(operation, data)
```

> [!TIP]
>
> It is recommended to provide detailed documentation and examples for each function to help the model better understand the function's purpose and usage.

> [!WARNING]
>
> Function calling involves code execution. Please ensure appropriate security measures are implemented, including input validation, permission control, and error handling.


# Tool Streaming Output
Source: https://docs.z.ai/guides/capabilities/stream-tool



Stream Tool Call is a unique feature of Z.ai's latest model, allowing real-time access to reasoning processes, response content, and tool call information during tool invocation, providing better user experience and real-time feedback.

## Features

Streaming output for tool call responses allows developers to stream tool usage parameters without buffering or JSON validation when calling `chat.completions`, reducing call latency and providing better user experience.

### Core Parameter Description

* **`stream=True`**: Enable streaming output, must be set to `True`
* **`tool_stream=True`**: Enable tool call streaming output
* **`model`**: Use a model that supports tool calling, limited to `glm-5`

### Response Parameter Description

The `delta` object in streaming responses contains the following fields:

* **`reasoning_content`**: Text content of the model's reasoning process
* **`content`**: Text content of the model's response
* **`tool_calls`**: Tool call information, including function names and parameters

## Code Examples

By setting the `tool_stream=True` parameter, you can enable streaming tool call functionality:

**Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Create streaming tool call request
response = client.chat.completions.create(
    model="glm-5.3",  # Use model that supports tool calling
    messages=[
        {"role": "user", "content": "How's the weather in Beijing?"},
    ],
    tools=[
        {
            "type": "function",
            "function": {
                "name": "get_weather",
                "description": "Get current weather conditions for a specified location",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "location": {"type": "string", "description": "City, e.g.: Beijing, Shanghai"},
                        "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
                    },
                    "required": ["location"]
                }
            }
        }
    ],
    stream=True,        # Enable streaming output
    tool_stream=True    # Enable tool call streaming output
)

# Initialize variables to collect streaming data
reasoning_content = ""      # Reasoning process content
content = ""               # Response content
final_tool_calls = {}      # Tool call information
reasoning_started = False  # Reasoning process start flag
content_started = False    # Content output start flag

# Process streaming response
for chunk in response:
    if not chunk.choices:
        continue

    delta = chunk.choices[0].delta

    # Handle streaming reasoning process output
    if hasattr(delta, 'reasoning_content') and delta.reasoning_content:
        if not reasoning_started and delta.reasoning_content.strip():
            print("\n🧠 Thinking Process:")
            reasoning_started = True
        reasoning_content += delta.reasoning_content
        print(delta.reasoning_content, end="", flush=True)

    # Handle streaming response content output
    if hasattr(delta, 'content') and delta.content:
        if not content_started and delta.content.strip():
            print("\n\n💬 Response Content:")
            content_started = True
        content += delta.content
        print(delta.content, end="", flush=True)

    # Handle streaming tool call information
    if delta.tool_calls:
        for tool_call in delta.tool_calls:
            index = tool_call.index
            if index not in final_tool_calls:
                # New tool call
                final_tool_calls[index] = tool_call
                final_tool_calls[index].function.arguments = tool_call.function.arguments
            else:
                # Append tool call parameters (streaming construction)
                final_tool_calls[index].function.arguments += tool_call.function.arguments

# Output final tool call information
if final_tool_calls:
    print("\n📋 Function Calls Triggered:")
    for index, tool_call in final_tool_calls.items():
        print(f"  {index}: Function Name: {tool_call.function.name}, Parameters: {tool_call.function.arguments}")
```

## Application Scenarios

> **Intelligent Customer Service**
>
> * Real-time query progress display
> * Improved waiting experience

> **Code Assistant**
>
> * Real-time code analysis process
> * Display tool call chains


# Streaming Messages
Source: https://docs.z.ai/guides/capabilities/streaming



Streaming Messages allow real-time content retrieval while the model generates responses, without waiting for the complete response to be generated. This approach can significantly improve user experience, especially when generating long text content, as users can immediately see output beginning to appear.

## Features

Streaming messages use an incremental generation mechanism, transmitting content in chunks in real-time during the generation process, rather than waiting for the complete response to be generated before returning it all at once.

### Core Parameter Description

* **`stream=True`**: Enable streaming output, must be set to `True`
* **`model`**: Models that support streaming output, such as `glm-5.3`,  `glm-5.2`, etc.

### Response Format Description

Streaming responses use Server-Sent Events (SSE) format, with each event containing:

* `choices[0].delta.content`: Incremental text content
* `choices[0].delta.reasoning_content`: Incremental reasoning content
* `choices[0].finish_reason`: Completion reason (only appears in the last chunk)
* `usage`: Token usage statistics (only appears in the last chunk)

## Code Examples

**cURL**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Write a poem about spring"
        }
    ],
    "stream": true
}'
```

**Python**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Create streaming message request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Write a poem about spring"}
    ],
    stream=True  # Enable streaming output
)

# Process streaming response
full_content = ""
for chunk in response:
    if not chunk.choices:
        continue
    
    delta = chunk.choices[0].delta
    
    # Handle incremental content
    if hasattr(delta, 'content') and delta.content:
        full_content += delta.content
        print(delta.content, end="", flush=True)
    
    # Check if completed
    if chunk.choices[0].finish_reason:
        print(f"\n\nCompletion reason: {chunk.choices[0].finish_reason}")
        if hasattr(chunk, 'usage') and chunk.usage:
            print(f"Token usage: Input {chunk.usage.prompt_tokens}, Output {chunk.usage.completion_tokens}")

print(f"\n\nComplete content:\n{full_content}")
```

### Response Example

The streaming response format is as follows:

```
data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":"Spring"},"finish_reason":null}]}

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":" comes"},"finish_reason":null}]}

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":" with"},"finish_reason":null}]}

...

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"finish_reason":"stop","delta":{"role":"assistant","content":""}}],"usage":{"prompt_tokens":8,"completion_tokens":262,"total_tokens":270,"prompt_tokens_details":{"cached_tokens":0}}}

data: [DONE]
```

## Application Scenarios

> **Chat Applications**
>
> * Real-time conversation experience
> * Character-by-character reply display
> * Reduced waiting time

> **Content Generation**
>
> * Article writing assistant
> * Code generation tools
> * Creative content creation

> **Educational Applications**
>
> * Online Q\&A systems
> * Learning assistance tools
> * Knowledge Q\&A platforms

> **Customer Service Systems**
>
> * Intelligent customer service bots
> * Real-time problem solving
> * User support systems


# Structured Output
Source: https://docs.z.ai/guides/capabilities/struct-output



Structured output (JSON mode) ensures that AI returns JSON data conforming to predefined formats, providing reliable guarantees for programmatic processing of AI outputs.

## Features

The structured output feature provides AI models with strict data format control capabilities, supporting various complex data structures and validation requirements.

### Core Parameters

* **`response_format`**: Specifies the response format, set to `{"type": "json_object"}` to enable JSON mode
* **`model`**: Use models that support structured output, such as `glm-5`, `glm-4.7`, `glm-4.6`, etc.
* **`messages`**: Define the expected JSON structure and field requirements in system messages

## Code Examples

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

The following is a complete structured output example demonstrating how to perform sentiment analysis and return structured JSON results:

```python theme={null}
from zai import ZaiClient
import json

# Initialize client
client = ZaiClient(api_key="your-api-key")

# Basic JSON mode
response = client.chat.completions.create(
    model="glm-5",
    messages=[
        {
            "role": "system",
            "content": """
            You are a sentiment analysis expert. Please return analysis results in the following JSON format:
            {
                "sentiment": "positive/negative/neutral",
                "confidence": 0.95,
                "emotions": ["joy", "excitement"],
                "keywords": ["weather", "mood"],
                "analysis": "Detailed analysis explanation"
            }
            """
        },
        {
            "role": "user",
            "content": "Please analyze the sentiment of this sentence: 'The weather is really nice today, I'm feeling very happy!'"
        }
    ],
    response_format={
        "type": "json_object"
    }
)

# Parse results
result = json.loads(response.choices[0].message.content)
print(f"Sentiment: {result['sentiment']}")
print(f"Confidence: {result['confidence']}")
print(f"Emotions: {result['emotions']}")
```

## Basic Usage

**Simple JSON Output**

**Simple JSON Output**

```python theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="your-api-key")

# Basic JSON mode
response = client.chat.completions.create(
    model="glm-5",
    messages=[
        {
            "role": "user",
            "content": "Please analyze the sentiment of this sentence: 'The weather is really nice today, I'm feeling very happy!'"
        }
    ],
    response_format={
        "type": "json_object"
    }
)

import json
result = json.loads(response.choices[0].message.content)
print(result)
```

**Specify JSON Structure**

### Specify JSON Structure

```python theme={null}
# Specify specific JSON structure
response = client.chat.completions.create(
    model="glm-5",
    messages=[
        {
            "role": "system",
            "content": """
            You are a sentiment analysis expert. Please return analysis results in the following JSON format:
            {
                "sentiment": "positive/negative/neutral",
                "confidence": 0.95,
                "emotions": ["joy", "excitement"],
                "keywords": ["weather", "mood"],
                "analysis": "Detailed analysis explanation"
            }
            """
        },
        {
            "role": "user",
            "content": "Please analyze the sentiment of this sentence: 'The weather is really nice today, I'm feeling very happy!'"
        }
    ],
    response_format={
        "type": "json_object"
    }
)

result = json.loads(response.choices[0].message.content)
print(f"Sentiment: {result['sentiment']}")
print(f"Confidence: {result['confidence']}")
print(f"Emotions: {result['emotions']}")
```

**Schema Validation**

### Using JSON Schema Validation

```python theme={null}
import jsonschema
from jsonschema import validate

# Define JSON Schema
schema = {
    "type": "object",
    "properties": {
        "sentiment": {
            "type": "string",
            "enum": ["positive", "negative", "neutral"]
        },
        "confidence": {
            "type": "number",
            "minimum": 0,
            "maximum": 1
        },
        "emotions": {
            "type": "array",
            "items": {"type": "string"}
        },
        "keywords": {
            "type": "array",
            "items": {"type": "string"}
        },
        "analysis": {
            "type": "string"
        }
    },
    "required": ["sentiment", "confidence", "analysis"]
}

def analyze_sentiment_with_validation(text):
    """Sentiment analysis with validation"""
    response = client.chat.completions.create(
        model="glm-5",
        messages=[
            {
                "role": "system",
                "content": f"""
                Please return sentiment analysis results according to the following JSON Schema format:
                {json.dumps(schema, indent=2, ensure_ascii=False)}
                """
            },
            {
                "role": "user",
                "content": f"Please analyze the sentiment of this sentence: '{text}'"
            }
        ],
        response_format={"type": "json_object"}
    )
    
    try:
        result = json.loads(response.choices[0].message.content)
        # Validate JSON structure
        validate(instance=result, schema=schema)
        return result
    except jsonschema.exceptions.ValidationError as e:
        print(f"JSON validation failed: {e}")
        return None
    except json.JSONDecodeError as e:
        print(f"JSON parsing failed: {e}")
        return None

# Usage example
result = analyze_sentiment_with_validation("The weather is really nice today, I'm feeling very happy!")
if result:
    print("Analysis result:", result)
```

## Scenario Examples

> [!WARNING]
>
> When using JSON mode for data extraction, please ensure the quality and format of input data to achieve the best extraction results.

> **Data Extraction and Structuring Complete Implementation**
>
> ```python theme={null}
> class DataExtractor:
>     def __init__(self, api_key):
>         self.client = ZaiClient(api_key=api_key)
>
>     def extract_contact_info(self, text):
>         """Extract contact information"""
>         schema = {
>             "type": "object",
>             "properties": {
>                 "contacts": {
>                     "type": "array",
>                     "items": {
>                         "type": "object",
>                         "properties": {
>                             "name": {"type": "string"},
>                             "phone": {"type": "string"},
>                             "email": {"type": "string"},
>                             "company": {"type": "string"},
>                             "position": {"type": "string"},
>                             "address": {"type": "string"}
>                         },
>                         "required": ["name"]
>                     }
>                 },
>                 "total_count": {"type": "integer"},
>                 "extraction_confidence": {"type": "number"}
>             },
>             "required": ["contacts", "total_count"]
>         }
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     You are an information extraction expert. Please extract all contact information from the text,
>                     return in the following JSON format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Note:
>                     - If a field has no information, do not include that field
>                     - phone field should be in standardized phone number format
>                     - email field should be a valid email address
>                     - extraction_confidence represents overall extraction confidence (0-1)
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Please extract contact information from the following text:\n\n{text}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)["properties"]
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Extraction failed: {e}")
>             return None
>
>     def extract_product_info(self, product_description):
>         """Extract product information"""
>         schema = {
>             "type": "object",
>             "properties": {
>                 "product_name": {"type": "string"},
>                 "brand": {"type": "string"},
>                 "category": {"type": "string"},
>                 "price": {
>                     "type": "object",
>                     "properties": {
>                         "amount": {"type": "number"},
>                         "currency": {"type": "string"},
>                         "original_price": {"type": "number"},
>                         "discount": {"type": "number"}
>                     }
>                 },
>                 "specifications": {
>                     "type": "object",
>                     "additionalProperties": True
>                 },
>                 "features": {
>                     "type": "array",
>                     "items": {"type": "string"}
>                 },
>                 "availability": {
>                     "type": "object",
>                     "properties": {
>                         "in_stock": {"type": "boolean"},
>                         "quantity": {"type": "integer"},
>                         "shipping_time": {"type": "string"}
>                     }
>                 },
>                 "ratings": {
>                     "type": "object",
>                     "properties": {
>                         "average_rating": {"type": "number"},
>                         "total_reviews": {"type": "integer"}
>                     }
>                 }
>             },
>             "required": ["product_name"]
>         }
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please extract structured information from product description, return in the following format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Note:
>                     - Price information should accurately extract values and currency units
>                     - specifications should include all technical specifications
>                     - features should list main functional features
>                     - Do not guess if information is unclear
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Product description:\n{product_description}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Product information extraction failed: {e}")
>             return None
>
>     def extract_event_info(self, event_text):
>         """Extract event information"""
>         schema = {
>             "type": "object",
>             "properties": {
>                 "events": {
>                     "type": "array",
>                     "items": {
>                         "type": "object",
>                         "properties": {
>                             "title": {"type": "string"},
>                             "description": {"type": "string"},
>                             "start_time": {"type": "string"},
>                             "end_time": {"type": "string"},
>                             "location": {"type": "string"},
>                             "organizer": {"type": "string"},
>                             "participants": {
>                                 "type": "array",
>                                 "items": {"type": "string"}
>                             },
>                             "category": {"type": "string"},
>                             "priority": {
>                                 "type": "string",
>                                 "enum": ["high", "medium", "low"]
>                             },
>                             "status": {
>                                 "type": "string",
>                                 "enum": ["scheduled", "ongoing", "completed", "cancelled"]
>                             }
>                         },
>                         "required": ["title", "start_time"]
>                     }
>                 }
>             },
>             "required": ["events"]
>         }
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please extract all event information from the text, return in the following format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Time format requirements:
>                     - Use ISO 8601 format: YYYY-MM-DDTHH:MM:SS
>                     - If only date available, use: YYYY-MM-DD
>                     - If time is unclear, try to infer reasonable time
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Please extract event information from the following text:\n\n{event_text}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Event information extraction failed: {e}")
>             return None
>
> # Usage example
> extractor = DataExtractor("your_api_key")
>
> # Extract contact information
> contact_text = """
> Zhang San, mobile: 13800138000, email: zhangsan@example.com,
> works as Technical Director at Beijing Technology Co., Ltd.
> Company address: No. 123, Technology Park, Chaoyang District, Beijing.
>
> Li Si, phone: 010-12345678, work email: lisi@company.com,
> is a Product Manager at Shanghai Innovation Company.
> """
>
> contacts = extractor.extract_contact_info(contact_text)
> if contacts:
>     print(f"Extracted {contacts['total_count']} contacts")
>     for contact in contacts['contacts']:
>         print(f"Name: {contact['name']}")
>         if 'phone' in contact:
>             print(f"Phone: {contact['phone']}")
> ```

> **API Response Formatting Complete Implementation**
>
> ```python theme={null}
> class APIResponseFormatter:
>     def __init__(self, api_key):
>         self.client = ZaiClient(api_key=api_key)
>
>     def format_search_results(self, query, raw_results):
>         """Format search results"""
>         schema = {
>             "type": "object",
>             "properties": {
>                 "query": {"type": "string"},
>                 "total_results": {"type": "integer"},
>                 "results": {
>                     "type": "array",
>                     "items": {
>                         "type": "object",
>                         "properties": {
>                             "title": {"type": "string"},
>                             "url": {"type": "string"},
>                             "snippet": {"type": "string"},
>                             "relevance_score": {"type": "number"},
>                             "source_type": {"type": "string"},
>                             "publish_date": {"type": "string"},
>                             "tags": {
>                                 "type": "array",
>                                 "items": {"type": "string"}
>                             }
>                         },
>                         "required": ["title", "url", "snippet"]
>                     }
>                 },
>                 "suggestions": {
>                     "type": "array",
>                     "items": {"type": "string"}
>                 },
>                 "filters": {
>                     "type": "object",
>                     "properties": {
>                         "date_range": {"type": "string"},
>                         "source_types": {
>                             "type": "array",
>                             "items": {"type": "string"}
>                         },
>                         "languages": {
>                             "type": "array",
>                             "items": {"type": "string"}
>                         }
>                     }
>                 }
>             },
>             "required": ["query", "total_results", "results"]
>         }
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please format search results into standard JSON format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Requirements:
>                     - Calculate relevance score for each result (0-1)
>                     - Identify content type (article, video, image, document, etc.)
>                     - Extract publish date (if available)
>                     - Generate relevant tags
>                     - Provide search suggestions
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Query: {query}\n\nRaw results:\n{raw_results}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Formatting failed: {e}")
>             return None
>
>     def format_analytics_data(self, raw_data, metrics):
>         """Format analytics data"""
>         schema = {
>             "type": "object",
>             "properties": {
>                 "summary": {
>                     "type": "object",
>                     "properties": {
>                         "total_records": {"type": "integer"},
>                         "date_range": {
>                             "type": "object",
>                             "properties": {
>                                 "start_date": {"type": "string"},
>                                 "end_date": {"type": "string"}
>                             }
>                         },
>                         "key_insights": {
>                             "type": "array",
>                             "items": {"type": "string"}
>                         }
>                     }
>                 },
>                 "metrics": {
>                     "type": "object",
>                     "additionalProperties": {
>                         "type": "object",
>                         "properties": {
>                             "current_value": {"type": "number"},
>                             "previous_value": {"type": "number"},
>                             "change_percentage": {"type": "number"},
>                             "trend": {
>                                 "type": "string",
>                                 "enum": ["up", "down", "stable"]
>                             },
>                             "unit": {"type": "string"}
>                         }
>                     }
>                 },
>                 "time_series": {
>                     "type": "array",
>                     "items": {
>                         "type": "object",
>                         "properties": {
>                             "timestamp": {"type": "string"},
>                             "values": {
>                                 "type": "object",
>                                 "additionalProperties": {"type": "number"}
>                             }
>                         }
>                     }
>                 },
>                 "segments": {
>                     "type": "array",
>                     "items": {
>                         "type": "object",
>                         "properties": {
>                             "name": {"type": "string"},
>                             "value": {"type": "number"},
>                             "percentage": {"type": "number"},
>                             "color": {"type": "string"}
>                         }
>                     }
>                 }
>             },
>             "required": ["summary", "metrics"]
>         }
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please format analytics data into standard format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Focus indicators:{', '.join(metrics)}
>
>                     Requirements:
>                     - Calculate change percentage and trend
>                     - Provide key insights
>                     - Time series data sorted by time
>                     - Segments data contain percentage
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Raw data: \n{raw_data}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Analytics data formatting failed: {e}")
>             return None
>
> # Usage example
> formatter = APIResponseFormatter("your_api_key")
>
> # Format search results
> raw_search = """
> 1. Python Programming Tutorial - https://example.com/python-tutorial
>    Detailed introduction to Python basic syntax and programming concepts...
>
> 2. Python Data Analysis Practice - https://example.com/python-data
>    Using pandas and numpy for data processing...
> """
>
> formatted_results = formatter.format_search_results("Python Tutorial", raw_search)
> if formatted_results:
>     print(f"Found {formatted_results['total_results']} results")
>     for result in formatted_results['results']:
>         print(f"Title: {result['title']}")
>         print(f"Relevance: {result['relevance_score']}")
> ```

> **Configuration Management and Validation Complete Implementation**
>
> ```python theme={null}
> class ConfigurationManager:
>     def __init__(self, api_key):
>         self.client = ZaiClient(api_key=api_key)
>
>     def parse_config_file(self, config_text, config_type="general"):
>         """Parse configuration file"""
>         schemas = {
>             "database": {
>                 "type": "object",
>                 "properties": {
>                     "connections": {
>                         "type": "array",
>                         "items": {
>                             "type": "object",
>                             "properties": {
>                                 "name": {"type": "string"},
>                                 "host": {"type": "string"},
>                                 "port": {"type": "integer"},
>                                 "database": {"type": "string"},
>                                 "username": {"type": "string"},
>                                 "ssl": {"type": "boolean"},
>                                 "pool_size": {"type": "integer"}
>                             },
>                             "required": ["name", "host", "database"]
>                         }
>                     },
>                     "settings": {
>                         "type": "object",
>                         "properties": {
>                             "timeout": {"type": "integer"},
>                             "retry_attempts": {"type": "integer"},
>                             "log_level": {
>                                 "type": "string",
>                                 "enum": ["DEBUG", "INFO", "WARNING", "ERROR"]
>                             }
>                         }
>                     }
>                 },
>                 "required": ["connections"]
>             },
>             "api": {
>                 "type": "object",
>                 "properties": {
>                     "endpoints": {
>                         "type": "array",
>                         "items": {
>                             "type": "object",
>                             "properties": {
>                                 "name": {"type": "string"},
>                                 "url": {"type": "string"},
>                                 "method": {
>                                     "type": "string",
>                                     "enum": ["GET", "POST", "PUT", "DELETE"]
>                                 },
>                                 "headers": {"type": "object"},
>                                 "timeout": {"type": "integer"},
>                                 "rate_limit": {"type": "integer"}
>                             },
>                             "required": ["name", "url", "method"]
>                         }
>                     },
>                     "authentication": {
>                         "type": "object",
>                         "properties": {
>                             "type": {
>                                 "type": "string",
>                                 "enum": ["bearer", "basic", "api_key"]
>                             },
>                             "credentials": {"type": "object"}
>                         }
>                     }
>                 },
>                 "required": ["endpoints"]
>             }
>         }
>
>         schema = schemas.get(config_type, {
>             "type": "object",
>             "additionalProperties": True
>         })
>
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please parse the configuration file and convert to JSON format:
>                     {json.dumps(schema, indent=2, ensure_ascii=False)}
>
>                     Configuration type: {config_type}
>
>                     Requirements:
>                     - Identify configuration items and values
>                     - Convert data types (string, number, boolean)
>                     - Handle arrays and nested objects
>                     - Validate required fields
>                     - Provide default values (if applicable)
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Configuration file content:\n{config_text}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             validate(instance=result, schema=schema)
>             return result
>         except Exception as e:
>             print(f"Configuration parsing failed: {e}")
>             return None
>
>     def validate_configuration(self, config_data, validation_rules):
>         """Validate configuration"""
>         response = self.client.chat.completions.create(
>             model="glm-5",
>             messages=[
>                 {
>                     "role": "system",
>                     "content": f"""
>                     Please validate configuration data and return validation results:
>
>                     Return format:
>                     {{
>                         "is_valid": true/false,
>                         "errors": [
>                             {{
>                                 "field": "field_name",
>                                 "error": "error_description",
>                                 "severity": "error/warning/info"
>                             }}
>                         ],
>                         "warnings": [
>                             {{
>                                 "field": "field_name",
>                                 "message": "warning_message"
>                             }}
>                         ],
>                         "suggestions": [
>                             "improvement_suggestion_1",
>                             "improvement_suggestion_2"
>                         ]
>                     }}
>
>                     Validation rules: {validation_rules}
>                     """
>                 },
>                 {
>                     "role": "user",
>                     "content": f"Configuration data:\n{json.dumps(config_data, indent=2, ensure_ascii=False)}"
>                 }
>             ],
>             response_format={"type": "json_object"}
>         )
>
>         try:
>             result = json.loads(response.choices[0].message.content)
>             return result
>         except Exception as e:
>             print(f"Configuration validation failed: {e}")
>             return None
>
> # Usage example
> config_manager = ConfigurationManager("your_api_key")
>
> # Parse database configuration
> db_config_text = """
> [database]
> host = localhost
> port = 5432
> database = myapp
> username = admin
> ssl = true
> pool_size = 10
>
> [settings]
> timeout = 30
> retry_attempts = 3
> log_level = INFO
> """
>
> config = config_manager.parse_config_file(db_config_text, "database")
> if config:
>     print("Parsed configuration:", json.dumps(config, indent=2, ensure_ascii=False))
>
>     # Validate configuration
>     validation_rules = [
>         "Port number must be in range 1-65535",
>         "Database name cannot be empty",
>         "Connection pool size should be greater than 0",
>         "Timeout should be reasonable (1-300 seconds)"
>     ]
>
>     validation_result = config_manager.validate_configuration(config, validation_rules)
>     if validation_result:
>         print(f"Configuration validity: {validation_result['is_valid']}")
>         if validation_result['errors']:
>             print("Errors:", validation_result['errors'])
>         if validation_result['warnings']:
>             print("Warnings:", validation_result['warnings'])
> ```

## Best Practices

> **Schema Design Principles**
>
> * Clarity: Field names and types should be clear and explicit
> * Completeness: Include all necessary validation rules
> * Flexibility: Consider future expansion needs

> **Error Handling Strategy**
>
> * Multi-layer validation: Schema validation + business logic validation
> * Fallback plan: Prepare simplified backup Schema
> * Logging: Record detailed error information

> [!WARNING]
>
> JSON mode requires AI to strictly output according to specified format, but in some complex scenarios it may affect the naturalness of responses. It's recommended to find a balance between functionality and user experience.

> [!TIP]
>
> When designing JSON Schema, it's recommended to start with simple structures and gradually increase complexity. Also, providing detailed descriptions and examples for key fields helps AI better understand and generate JSON data that meets requirements.


# Deep Thinking
Source: https://docs.z.ai/guides/capabilities/thinking



Deep Thinking is an advanced reasoning feature that enables Chain of Thought mechanisms, allowing the model to perform deep analysis and reasoning before answering questions. This approach significantly improves the model's accuracy and interpretability in complex tasks, particularly suitable for scenarios requiring multi-step reasoning, logical analysis, and problem-solving.

## Features

The Deep Thinking feature currently supports the latest models in the GLM-5.3, GLM-5.3-FLASH, GLM-5.2, GLM-5.1, GLM-5, GLM-4.7, GLM-4.6 series. By enabling deep thinking, the model can:

* **Multi-step Reasoning**: Break down complex problems into multiple steps for gradual analysis and resolution
* **Logical Analysis**: Provide clear reasoning processes and logical chains
* **Improved Accuracy**: Reduce errors and improve answer quality through deep thinking
* **Enhanced Interpretability**: Display the thinking process to help users understand the model's reasoning logic
* **Intelligent Judgment**: The model automatically determines whether deep thinking is needed to optimize response efficiency

### Core Parameters

Note: GLM-5.3 and GLM-5.3-FLASH no longer support disabling thinking (an error will occur if the `thinking.type` parameter in the API request is set to `disabled`). Please ensure that thinking is enabled.

* **`thinking.type`**: Controls the deep thinking mode
  * `enabled` (default): Enable dynamic thinking. The model automatically determines whether to think: `GLM-5.2`, `GLM-5.1`, `GLM-5`, `GLM-4.6`auto-decide whether to think, while `GLM-5.3`, `GLM-5.3-FLASH`, `GLM-4.7`use forced thinking
  * `disabled`: Disable deep thinking, provide direct answers
* **`reasoning_effort`**: Controls the degree of reasoning within the thought chain, and is only supported by `GLM-5.2` and above.
  * Available values: `max` (default and recommended, deep inference), `high` (enhanced inference), `low` (mild inference, only supported by GLM-5.3 and GLM-5.3-FLASH)
  * In the API request:
  * * For GLM-5.3 and GLM-5.3-FLASH, only `max`, `high` and `low` are supported. Any other input will result in an error.
  * * For GLM-5.2, the supported options are `max` (default and recommended, for deep inference), `xhigh`, `high` (enhanced inference), `medium`, `low`, `minimal`, and `none`. Among them, `none` or `minimal` indicate that the model stops thinking; `low`/`medium` are mapped to `high`; `xhigh` is mapped to `max`.
  * In the Coding Plan request:
  * * For GLM-5.3 and GLM-5.3-FLASH, `none`, `minimal`, and `low` are mapped to `low`; `medium`, `high` are mapped to `high`; `xhigh` and `max` are mapped to `max`.
  * * For GLM-5.2, `none` or `minimal` indicate that the model stops thinking; `low` / `medium` are mapped to `high`; `xhigh` is mapped to `max`.
* **`model`**: A model that enables deep thinking, supported by `GLM-4.6` and above versions.

## Code Examples

**cURL**

**Basic Call (Enable Deep Thinking)**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Explain in detail the basic principles of quantum computing and analyze its potential impact in the field of cryptography"
        }
    ],
    "thinking": {
        "type": "enabled"
    },
    "max_tokens": 4096,
    "temperature": 1.0
}'
```

**Streaming Call (Deep Thinking + Streaming Output)**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Design a recommendation system architecture for an e-commerce website, considering user behavior, product features, and real-time requirements"
        }
    ],
    "thinking": {
        "type": "enabled"
    },
    "stream": true,
    "max_tokens": 4096,
    "temperature": 1.0
}'
```

**Control Reasoning Effort (reasoning\_effort)**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Analyze the solution approach for this math problem"
        }
    ],
    "thinking": {
        "type": "enabled"
    },
    "reasoning_effort": "max"
}'
```

**Disable deep thinking (GLM-5.3 is no longer supported, only GLM-5.2 supports it)**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.2",
    "messages": [
        {
            "role": "user",
            "content": "How is the weather today?"
        }
    ],
    "thinking": {
        "type": "disabled"
    }
}'
```

**Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Basic Call (Enable Deep Thinking)**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='your_api_key')

# Create deep thinking request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Explain in detail the basic principles of quantum computing and analyze its potential impact in the field of cryptography"}
    ],
    thinking={
        "type": "enabled"  # Enable deep thinking mode
    },
    max_tokens=4096,
    temperature=1.0
)

print("Model response:")
print(response.choices[0].message.content)
print("\n---")
print(response.choices[0].message.reasoning_content)
```

**Streaming Call (Deep Thinking + Streaming Output)**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='your_api_key')

# Create streaming deep thinking request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Design a recommendation system architecture for an e-commerce website, considering user behavior, product features, and real-time requirements"}
    ],
    thinking={
        "type": "enabled"  # Enable deep thinking mode
    },
    stream=True,  # Enable streaming output
    max_tokens=4096,
    temperature=1.0
)

# Process streaming response
reasoning_content = ""
thinking_phase = True

for chunk in response:
    if not chunk.choices:
        continue
    
    delta = chunk.choices[0].delta
    
    # Process thinking process (if any)
    if hasattr(delta, 'reasoning_content') and delta.reasoning_content:
        reasoning_content += delta.reasoning_content
        if thinking_phase:
            print("🧠 Thinking...", end="", flush=True)
            thinking_phase = False
        print(delta.reasoning_content, end="", flush=True)
    
    # Process answer content
    if hasattr(delta, 'content') and delta.content:
        if thinking_phase:
            print("\n\n💡 Answer:")
            thinking_phase = False
        print(delta.content, end="", flush=True)

```

**Control Reasoning Effort (reasoning\_effort)**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='your_api_key')

# Use reasoning_effort to control reasoning level
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Analyze the solution approach for this math problem"}
    ],
    thinking={
        "type": "enabled"
    },
    reasoning_effort="high"
)

print(response.choices[0].message.content)
print(response.choices[0].message.reasoning_content)
```

**Disable Deep Thinking**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='your_api_key')

# Disable deep thinking for quick response
response = client.chat.completions.create(
    model="glm-5.2",
    messages=[
        {"role": "user", "content": "How is the weather today?"}
    ],
    thinking={
        "type": "disabled"  # Disable deep thinking mode
    }
)

print(response.choices[0].message.content)
```

### Response Example

Response format with deep thinking enabled:

```json theme={null}
{
  "created": 1677652288,
  "model": "glm-5.3",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Artificial intelligence has tremendous application prospects in medical diagnosis...",
        "reasoning_content": "Let me analyze this question from multiple angles. First, I need to consider the technical advantages of AI in medical diagnosis..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "completion_tokens": 239,
    "prompt_tokens": 8,
    "prompt_tokens_details": {
      "cached_tokens": 0
    },
    "total_tokens": 247
  }
}
```

## Best Practices

**Recommended scenarios to enable:**

* Complex problem analysis and solving
* Multi-step reasoning tasks
* Technical solution design
* Strategy planning and decision
* Academic research and analysis
* Creative writing and content creation

**Can be disabled scenarios:**

* Simple fact query
* Basic translation tasks
* Simple classification judgment
* Quick question and answer requirements

## Application scenarios

> **Academic Research**
>
> * Research method design
> * Data analysis and explanation
> * Theory deduction and proof

> **Technology Consulting**
>
> * System architecture design
> * Technological scheme evaluation
> * Problem diagnosis and solution

> **Business Analysis**
>
> * Market trends analysis
> * Business model design
> * Investment decision support

> **Education Training**
>
> * Complex concept explanation
> * Learning path planning
> * Knowledge system building

## Notes

1. **Response time**：Different reasoning effort levels will increase response time, particularly for complex tasks
2. **Token consumption**：Thinking process will consume extra tokens, please manage your tokens
3. **Model support**：Ensure you're using models that support deep thinking
4. **Task matching**：Choose whether to enable deep thinking according to the task complexity
5. **Streaming output**：Combine streaming output to see the thinking process, improving user experience


# Thinking Mode
Source: https://docs.z.ai/guides/capabilities/thinking-mode



GLM offers multiple thinking modes for different scenarios. The sections below explain how to enable each mode, key considerations, and example usage.

## **Default Thinking Behaviour**

Thinking is activated by default in GLM-5.3 GLM-5.3-FLASH GLM-5.2 GLM-5.1 GLM-5 GLM-4.7 series, different from the default hybrid thinking in GLM-4.6.

> If you want to disable thinking, use the following (note that `GLM-5.3` and `GLM-5.3-FLASH` use forced thinking and cannot be disabled):

```bash theme={null}
"thinking": {
    "type": "disabled"
}
```

## **Interleaved thinking**

We support **interleaved thinking** by default, allowing GLM to think between tool calls and after receiving tool results. This enables more complex, step-by-step reasoning: interpreting each tool output before deciding what to do next, chaining multiple tool calls with reasoning steps, and making finer-grained decisions based on intermediate results.

> [!TIP]
>
> When using interleaved thinking with tools, **thinking blocks should be explicitly preserved and returned together with the tool results.**

The detailed interleaved thinking process is as follows.

![Description](https://cdn.bigmodel.cn/markdown/1766025484368img_v3_02t3_4677ac48-b748-44d8-a56f-8cbd599b51ag.jpg?attname=img_v3_02t3_4677ac48-b748-44d8-a56f-8cbd599b51ag.jpg)

## **Preserved thinking**

**We introduce a new capability** in coding scenarios: the model can retain **reasoning content from previous assistant turns** in the context. This helps preserve reasoning continuity and conversation integrity, improves model performance, and increases cache hit rates—saving tokens in real tasks.

> [!TIP]
>
> This capability is **enabled by default** on the **Coding Plan endpoint** and **disabled by default** on the **standard API endpoint**. If you want to enable **Preserved Thinking** in your product (primarily recommended for coding/agent scenarios), you can turn it on for the API endpoint by setting **"clear\_thinking": false**, and **you must return the complete**, unmodified reasoning\_content back to the API.
>
> All consecutive reasoning\_content blocks must **exactly match the original sequence** generated by the model during the initial request. Do not reorder or edit these blocks; otherwise, performance may degrade and cache hit rates may be affected.

The detailed Preserved thinking process is as follows.

![Description](https://cdn.bigmodel.cn/markdown/176641919972020251222-235942.jpeg?attname=20251222-235942.jpeg)

## Turn-level Thinking

“Turn-level Thinking” is a capability that **lets you control reasoning computation on a per-turn basis**: within the same session, each request can independently choose to enable or disable thinking. This is a new capability introduced in GLM-4.7, with the following advantages:

* **More flexible cost/latency control:** For lightweight turns like “asking a fact” or “tweaking wording,” you can disable thinking to get faster responses; for heavier tasks like “complex planning,” “multi-constraint reasoning,” or “code debugging,” you can enable thinking to improve accuracy and stability.
* **Smoother multi-turn experience:** The thinking switch can be toggled at any point within a session. The model stays coherent across turns and keeps a consistent output style, making it feel “smarter when things are hard, faster when things are simple.”
* **Better for agent/tool-use scenarios:** On turns that require quick tool execution, you can reduce reasoning overhead; on turns that require making decisions based on tool results, you can turn on deeper thinking—dynamically balancing efficiency and quality.

## Example Usage

This applies to both **Interleaved Thinking** and **Preserved Thinking**—no manual differentiation is required. **Remember to return the historical** `reasoning_content`**to keep the reasoning coherent.**

```python theme={null}
""""Interleaved Thinking + Tool Calling Example"""

import json
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://api.z.ai/api/paas/v4/",
)

tools = [{"type": "function", "function": {
    "name": "get_weather",
    "description": "Get weather information",
    "parameters": {"type": "object", "properties": {"city": {"type": "string"}}, "required": ["city"]},
}}]

messages = [
    {"role": "system", "content": "You are an assistant"},
    {"role": "user", "content": "What's the weather like in Beijing?"},
]

# Round 1: the model reasons and then calls a tool
response = client.chat.completions.create(model="glm-4.7", messages=messages, tools=tools, stream=True, extra_body={
        "thinking":{
        "type":"enabled",
        "clear_thinking": False  # False for Preserved Thinking
    }})
reasoning, content, tool_calls = "", "", []
for chunk in response:
    delta = chunk.choices[0].delta
    if hasattr(delta, "reasoning_content") and delta.reasoning_content:
        reasoning += delta.reasoning_content
    if hasattr(delta, "content") and delta.content:
        content += delta.content
    if hasattr(delta, "tool_calls") and delta.tool_calls:
        for tc in delta.tool_calls:
            if tc.index >= len(tool_calls):
                tool_calls.append({"id": tc.id, "function": {"name": "", "arguments": ""}})
            if tc.function.name:
                tool_calls[tc.index]["function"]["name"] = tc.function.name
            if tc.function.arguments:
                tool_calls[tc.index]["function"]["arguments"] += tc.function.arguments

print(f"Reasoning: {reasoning}\nTool calls: {tool_calls}")

# Key: return reasoning_content to keep the reasoning coherent
messages.append({"role": "assistant", "content": content, "reasoning_content": reasoning,
                 "tool_calls": [{"id": tc["id"], "type": "function", "function": tc["function"]} for tc in tool_calls]})
messages.append({"role": "tool", "tool_call_id": tool_calls[0]["id"],
                 "content": json.dumps({"weather": "Sunny", "temp": "25°C"})})

# Round 2: the model continues reasoning based on the tool result and responds
response = client.chat.completions.create(model="glm-4.7", messages=messages, tools=tools, stream=True, extra_body={
        "thinking":{
        "type":"enabled",
        "clear_thinking": False # False for Preserved Thinking
    }})
reasoning, content = "", ""
for chunk in response:
    delta = chunk.choices[0].delta
    if hasattr(delta, "reasoning_content") and delta.reasoning_content:
        reasoning += delta.reasoning_content
    if hasattr(delta, "content") and delta.content:
        content += delta.content

print(f"Reasoning: {reasoning}\nReply: {content}")
```


# HTTP API Calls
Source: https://docs.z.ai/guides/develop/http/introduction



Z.AI provides standard HTTP API interfaces that support multiple programming languages and development environments, allowing you to easily integrate Z.AI's powerful capabilities.

### Core Advantages

> **Cross-platform Compatible**
>
> Supports all programming languages and platforms that support HTTP protocol

> **Standard Protocol**
>
> Based on RESTful design, follows HTTP standards, easy to understand and use

> **Flexible Integration**
>
> Can be integrated into any existing applications and systems

> **Real-time Calls**
>
> Supports synchronous and asynchronous calls to meet different scenario requirements

## Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

## API Basic Information

### General API Endpoint

```
https://api.z.ai/api/paas/v4/
```

> [!WARNING]
>
> When using GLM Coding Plan, please follow the [tutorial](/devpack/quick-start) to configure your dedicated endpoint.

### Request Header Requirements

```http theme={null}
Content-Type: application/json
Authorization: Bearer YOUR_API_KEY
```

### Supported Authentication Methods

**API Key Authentication**

The simplest authentication method, directly using your API Key:

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Accept-Language: en-US,en' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Hello"
        }
    ]
}'
```

**JWT Token Authentication**

Use JWT Token for authentication, suitable for scenarios requiring higher security:
Install PyJWT

```shell theme={null}
pip install PyJWT
```

```python theme={null}
import time
import jwt

def generate_token(apikey: str, exp_seconds: int):
    try:
        id, secret = apikey.split(".")
    except Exception as e:
        raise Exception("invalid apikey", e)

    payload = {
        "api_key": id,
        "exp": int(round(time.time() * 1000)) + exp_seconds * 1000,
        "timestamp": int(round(time.time() * 1000)),
    }

    return jwt.encode(
        payload,
        secret,
        algorithm="HS256",
        headers={"alg": "HS256", "sign_type": "SIGN"},
    )

# Use the generated token
token = generate_token("your-api-key", 3600)  # 1 hour validity
```

## Basic Call Examples

### Simple Conversation

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Accept-Language: en-US,en' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Please introduce the development history of artificial intelligence"
        }
    ],
    "temperature": 1.0,
    "max_tokens": 1024
}'
```

### Streaming Response

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Accept-Language: en-US,en' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Write a poem about spring"
        }
    ],
    "stream": true
}'
```

### Multi-turn Conversation

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Accept-Language: en-US,en' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "system",
            "content": "You are a professional programming assistant"
        },
        {
            "role": "user",
            "content": "What is recursion?"
        },
        {
            "role": "assistant",
            "content": "Recursion is a programming technique where a function calls itself to solve problems..."
        },
        {
            "role": "user",
            "content": "Can you give me an example of Python recursion?"
        }
    ]
}'
```

## Common Programming Language Examples

**Python**

```python theme={null}
import requests
import json

def call_zai_api(messages, model="glm-5.3"):
    url = "https://api.z.ai/api/paas/v4/chat/completions"

    headers = {
        "Authorization": "Bearer YOUR_API_KEY",
        "Content-Type": "application/json",
        "Accept-Language": "en-US,en"
    }

    data = {
        "model": model,
        "messages": messages,
        "temperature": 1.0
    }

    response = requests.post(url, headers=headers, json=data)

    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"API call failed: {response.status_code}, {response.text}")

# Usage example
messages = [
    {"role": "user", "content": "Hello, please introduce yourself"}
]

result = call_zai_api(messages)
print(result['choices'][0]['message']['content'])
```

**JavaScript**

```javascript theme={null}
async function callZAPI(messages, model = 'glm-5.3') {
    const url = 'https://api.z.ai/api/paas/v4/chat/completions';

    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer YOUR_API_KEY',
            'Content-Type': 'application/json',
            'Accept-Language': 'en-US,en'
        },
        body: JSON.stringify({
            model: model,
            messages: messages,
            temperature: 1.0
        })
    });

    if (!response.ok) {
        throw new Error(`API call failed: ${response.status}`);
    }

    return await response.json();
}

// Usage example
const messages = [
    { role: 'user', content: 'Hello, please introduce yourself' }
];

callZAPI(messages)
    .then(result => {
        console.log(result.choices[0].message.content);
    })
    .catch(error => {
        console.error('Error:', error);
    });
```

**Java**

```java theme={null}
import com.fasterxml.jackson.databind.ObjectMapper;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class AgentExample {

    public static void main(String[] args) throws Exception {

        OkHttpClient client = new OkHttpClient();
        ObjectMapper mapper = new ObjectMapper();
        Map<String, String> messages = new HashMap<>(8);
        messages.put("role", "user");
        messages.put("content", "Hello, please introduce yourself");
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", "glm-5.3");
        requestBody.put("messages", Collections.singletonList(messages));
        requestBody.put("temperature", 1.0);

        String jsonBody = mapper.writeValueAsString(requestBody);
        MediaType JSON = MediaType.get("application/json; charset=utf-8");
        RequestBody body = RequestBody.create(JSON, jsonBody);
        Request request = new Request.Builder()
            .url("https://api.z.ai/api/paas/v4/chat/completions")
            .addHeader("Authorization", "Bearer your_api_key")
            .addHeader("Content-Type", "application/json")
            .addHeader("Accept-Language", "en-US,en")
            .post(body)
            .build();
        try (Response response = client.newCall(request).execute()) {
            System.out.println(response.body().string());
        }
    }
}
```

## Best Practices

> **Security**
>
> * Properly secure API Keys, do not hard-code them in your code
> * Use environment variables or configuration files to store sensitive information
> * Regularly rotate API Keys

> **Performance Optimization**
>
> * Implement connection pooling and session reuse
> * Set reasonable timeout values
> * Use asynchronous requests for high-concurrency scenarios

> **Error Handling**
>
> * Implement exponential backoff retry mechanisms
> * Log detailed error information
> * Set reasonable timeout and retry limits

> **Monitoring**
>
> * Monitor API call frequency and success rates
> * Track response times and error rates
> * Set up alerting mechanisms

## Get More

> **[API Documentation](/api-reference)**
>
> View complete API interface documentation and parameter descriptions

> **[Technical Support](https://z.ai/contact)**
>
> Get technical support and assistance

> [!NOTE]
>
> It is recommended to use HTTPS protocol in production environments and implement appropriate security measures to protect your API keys and data transmission.


# Official Java SDK
Source: https://docs.z.ai/guides/develop/java/introduction



Z.AI Java SDK is the official Java development toolkit provided by Z.AI, offering Java developers convenient and efficient AI model integration solutions.

### Core Advantages

> **Enterprise-grade**
>
> Designed for enterprise applications, supports high concurrency and high availability

> **Easy Integration**
>
> Clean API design, comprehensive documentation, quick integration into existing projects

> **Type Safety**
>
> Complete type definitions, compile-time error checking, reducing runtime errors

> **High Performance**
>
> Optimized network request handling, supports connection pooling and asynchronous calls

### Supported Features

* **💬 Chat Conversations**: Support for single-turn and multi-turn conversations, streaming and non-streaming responses
* **🔧 Function Calling**: Enable AI models to call your custom functions
* **👁️ Vision Understanding**: Image analysis, visual understanding
* **🎨 Image Generation**: Generate high-quality images from text descriptions
* **🎬 Video Generation**: Creative content generation from text to video
* **🔊 Speech Processing**: Speech-to-text, text-to-speech
* **📊 Text Embedding**: Text vectorization, supporting semantic search
* **🤖 Intelligent Assistants**: Build professional AI assistant applications

## Technical Specifications

### Environment Requirements

* **Java Version**: Java 1.8 or higher
* **Build Tools**: Maven 3.6+ or Gradle 6.0+
* **Network Requirements**: HTTPS connection support
* **API Key**: Valid Z.AI API key required

### Dependency Management

The SDK adopts a modular design, allowing you to selectively import functional modules as needed:

* **Core Module**: Basic API calling functionality
* **Async Module**: Asynchronous and concurrent processing support
* **Utility Module**: Utility tools and auxiliary functions

## Quick Start

### Environment Requirements

> **Java Version**
>
> Java 1.8 or higher

> **Build Tools**
>
> Maven 3.6+ or Gradle 6.0+

> [!WARNING]
>
> Supports Java 8, 11, 17, 21 versions, cross-platform compatible with Windows, macOS, Linux

### Add Dependencies

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle**

```gradle theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

> [!NOTE]
>
> Domestic Z.AI platform uses ZaiClient
>
> Domestic API URL: [https://api.z.ai/api/paas/v4/](https://api.z.ai/api/paas/v4/)

#### Create Client

**Environment Variable**

```java theme={null}
import ai.z.openapi.ZaiClient;

public class QuickStart {
    public static void main(String[] args) {
        // Read API Key from environment variable
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey(System.getenv("ZAI_API_KEY"))
            .build();
        
        // Or use directly (if environment variable is set)
        ZaiClient client2 = ZaiClient.builder().ofZAI().build();
    }
}
```

**Direct Setting**

```java theme={null}
import ai.z.openapi.ZaiClient;

public class QuickStart {
    public static void main(String[] args) {
        // Set API Key directly
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
    }
}
```

#### Basic Conversation

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.Arrays;

public class BasicChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
        
        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Hello, please introduce yourself")
                    .build()
            ))
            .build();
        
        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        // Get reply
        if (response.isSuccess()) {
            Object reply = response.getData().getChoices().get(0).getMessage().getContent();
            System.out.println("AI Reply: " + reply);
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

#### Streaming Conversation

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.Arrays;

public class StreamingChat {
    public static void main(String[] args) {
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
        
        // Create streaming chat request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Write a poem about spring")
                    .build()
            ))
            .stream(true)
            .build();
        
        // Handle streaming response
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        if (response.isSuccess() && response.getFlowable() != null) {
            response.getFlowable().subscribe(
                data -> {
                    // Handle streaming data chunks
                    if (data.getChoices() != null && !data.getChoices().isEmpty()) {
                        Delta content = data.getChoices().get(0).getDelta();
                        if (content != null) {
                            System.out.print(content);
                        }
                    }
                },
                error -> System.err.println("\nStreaming error: " + error.getMessage()),
                () -> System.out.println("\nStreaming completed")
            );
        }
    }
}
```

### Complete Example

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.*;

public class ChatBot {
    private final ZaiClient client;
    private final List<ChatMessage> conversation;
    
    public ChatBot(String apiKey) {
        this.client = ZaiClient.builder().ofZAI()
            .apiKey(apiKey)
            .build();
        this.conversation = new ArrayList<>();
        // Add system message
        this.conversation.add(ChatMessage.builder()
            .role(ChatMessageRole.SYSTEM.value())
            .content("You are a friendly AI assistant")
            .build());
    }
    
    public Object chat(String userInput) {
        try {
            // Add user message
        conversation.add(ChatMessage.builder()
            .role(ChatMessageRole.USER.value())
            .content(userInput)
            .build());
        
        // Create request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(conversation)
            .temperature(0.6f)
            .maxTokens(1000)
            .build();
        
        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        if (response.isSuccess()) {
            // Get AI response
            Object aiResponse = response.getData().getChoices().get(0).getMessage().getContent();
            
            // Add AI response to conversation history
            conversation.add(ChatMessage.builder()
                .role(ChatMessageRole.ASSISTANT.value())
                .content(aiResponse)
                .build());
            
            return aiResponse;
        } else {
            return "Error occurred: " + response.getMsg();
        }
        
    } catch (Exception e) {
        return "Error occurred: " + e.getMessage();
        }
    }
    
    public static void main(String[] args) {
        ChatBot bot = new ChatBot(System.getenv("ZAI_API_KEY"));
        Scanner scanner = new Scanner(System.in);
        
        System.out.println("Welcome to Z.ai chatbot! Type 'quit' to exit.");
        
        while (true) {
            System.out.print("You: ");
            String input = scanner.nextLine();
            
            if ("quit".equalsIgnoreCase(input)) {
                break;
            }
            
            Object response = bot.chat(input);
            System.out.println("AI: " + response);
        }
        
        System.out.println("Goodbye!");
        scanner.close();
    }
}
```

## Advanced Features

### Function Calling

Function calling allows AI models to call functions you define to get real-time information or perform specific operations.

#### Defining and Using Functions

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.*;

public class FunctionCallingExample {
    
    // Simulate weather API
    public static Map<String, Object> getWeather(String location, String date) {
        Map<String, Object> weather = new HashMap<>();
        weather.put("location", location);
        weather.put("date", date != null ? date : "today");
        weather.put("weather", "sunny");
        weather.put("temperature", "25°C");
        weather.put("humidity", "60%");
        return weather;
    }
    
    // Simulate stock API
    public static Map<String, Object> getStockPrice(String symbol) {
        Map<String, Object> stock = new HashMap<>();
        stock.put("symbol", symbol);
        stock.put("price", 150.25);
        stock.put("change", "+2.5%");
        return stock;
    }
    
    public static void main(String[] args) {
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey(System.getenv("ZAI_API_KEY"))
            .build();
        
        // Define function tools
        Map<String, ChatFunctionParameterProperty> properties = new HashMap<>();
        ChatFunctionParameterProperty locationProperty = ChatFunctionParameterProperty
                .builder().type("string").description("City name, for example: Beijing").build();
        properties.put("location", locationProperty);
        ChatFunctionParameterProperty unitProperty = ChatFunctionParameterProperty
                .builder().type("string").enums(Arrays.asList("celsius", "fahrenheit")).build();
        properties.put("unit", unitProperty);
        ChatTool weatherTool = ChatTool.builder()
                .type(ChatToolType.FUNCTION.value())
                .function(ChatFunction.builder()
                        .name("get_weather")
                        .description("Get weather information for a specified location")
                        .parameters(ChatFunctionParameters.builder()
                                .type("object")
                                .properties(properties)
                                .required(Collections.singletonList("location"))
                                .build())
                        .build())
                .build();
        
        // Create request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
                .model("glm-5.3")
                .messages(Collections.singletonList(
                        ChatMessage.builder()
                                .role(ChatMessageRole.USER.value())
                                .content("How's the weather in Beijing today?")
                                .build()
                ))
                .tools(Collections.singletonList(weatherTool))
                .toolChoice("auto")
                .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        if (response.isSuccess()) {
            // Handle function calling
            ChatMessage assistantMessage = response.getData().getChoices().get(0).getMessage();
            if (assistantMessage.getToolCalls() != null && !assistantMessage.getToolCalls().isEmpty()) {
                for (ToolCalls toolCall : assistantMessage.getToolCalls()) {
                    String functionName = toolCall.getFunction().getName();

                    if ("get_weather".equals(functionName)) {
                        Map<String, Object> result = getWeather("Beijing", null);
                        System.out.println("Weather info: " + result);
                    }
                }
            } else {
                System.out.println(assistantMessage.getContent());
            }
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

## Getting Help

> **[GitHub Repository](https://github.com/zai-org/z-ai-sdk-java)**
>
> View source code, submit issues, contribute

> **[API Reference](/api-reference)**
>
> View complete API documentation

> **[Example Projects](https://github.com/zai-org/z-ai-sdk-java/tree/main/samples)**
>
> Browse more practical application examples

> **[Best Practices](https://github.com/zai-org/z-ai-sdk-java)**
>
> Learn best practices for SDK usage

> [!NOTE]
>
> This SDK is developed based on the latest API specifications from Z.AI, ensuring synchronization with platform features. It is recommended to regularly update to the latest version for the best experience.


# LangChain Integration
Source: https://docs.z.ai/guides/develop/langchain/introduction



LangChain is a framework for developing applications powered by language models. Z.AI's integration with LangChain enables you to:

* Use LangChain's chain calling functionality
* Build intelligent agents and tool calling
* Implement complex conversation memory management

## Core Advantages

> **Framework Ecosystem**
>
> Access to LangChain's rich ecosystem and toolchain

> **Rapid Development**
>
> Quickly build complex AI applications using pre-built components

> **Modular Design**
>
> Flexibly combine different components to meet various needs

> **Community Support**
>
> Enjoy active open source community and rich resources

## Environment Requirements

> **Python Version**
>
> Python 3.8 or higher

> **LangChain Version**
>
> langchain\_community version 0.0.32 or higher

> [!WARNING]
>
> Please ensure langchain\_community version is 0.0.32 or higher for optimal compatibility and feature support.

## Install Dependencies

### Basic Installation

```bash theme={null}
# Install LangChain and related dependencies
pip install langchain langchainhub httpx_sse

# Install OpenAI compatible package
pip install langchain-openai
```

### Complete Installation

```bash theme={null}
# Install all related packages at once
pip install langchain langchain-openai langchainhub httpx_sse

# Verify installation
python -c "import langchain; print(langchain.__version__)"
```

## Quick Start

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

### Basic Configuration

```python theme={null}
import os
from langchain_openai import ChatOpenAI

# Create Z.AI LLM instance
llm = ChatOpenAI(
    temperature=0.6,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Or use environment variables
llm = ChatOpenAI(
    temperature=0.6,
    model="glm-5.3",
    openai_api_key=os.getenv("ZAI_API_KEY"),
    openai_api_base="https://api.z.ai/api/paas/v4/"
)
```

## Basic Usage Examples

### Simple Conversation

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.schema import HumanMessage, SystemMessage

# Create LLM instance
llm = ChatOpenAI(
    temperature=0.7,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create messages
messages = [
    SystemMessage(content="You are a helpful AI assistant"),
    HumanMessage(content="Please introduce the development history of artificial intelligence")
]

# Call the model
response = llm(messages)
print(response.content)
```

### Using Prompt Templates

```python theme={null}
from langchain.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI

# Create LLM
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create prompt template
prompt = ChatPromptTemplate.from_messages([
    ("system", "You are a professional {domain} expert"),
    ("human", "Please explain the concept and applications of {topic}")
])

# Create chain
chain = prompt | llm

# Invoke chain
response = chain.invoke({
    "domain": "machine learning",
    "topic": "deep learning"
})

print(response.content)
```

### Conversation Memory Management

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.prompts import (
    ChatPromptTemplate,
    MessagesPlaceholder,
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
)
from langchain.chains import LLMChain
from langchain.memory import ConversationBufferMemory

# Create LLM
llm = ChatOpenAI(
    temperature=1.0,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create prompt template
prompt = ChatPromptTemplate(
    messages=[
        SystemMessagePromptTemplate.from_template(
            "You are a nice chatbot having a conversation with a human."
        ),
        MessagesPlaceholder(variable_name="chat_history"),
        HumanMessagePromptTemplate.from_template("{question}")
    ]
)

# Create memory
memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)

# Create conversation chain
conversation = LLMChain(
    llm=llm,
    prompt=prompt,
    verbose=True,
    memory=memory
)

# Have conversations
response1 = conversation.invoke({"question": "tell me a joke"})
print("AI:", response1['text'])

response2 = conversation.invoke({"question": "tell me another one"})
print("AI:", response2['text'])
```

## Advanced Features

### Intelligent Agent

```python theme={null}
import os
from langchain import hub
from langchain.agents import AgentExecutor, create_react_agent
from langchain_community.tools.tavily_search import TavilySearchResults
from langchain_openai import ChatOpenAI

# Set search tool API key
os.environ["TAVILY_API_KEY"] = "your-tavily-api-key"

# Create LLM
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create tools
tools = [TavilySearchResults(max_results=2)]

# Get prompt template
prompt = hub.pull("hwchase17/react")

# Create agent
agent = create_react_agent(llm, tools, prompt)
agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

# Execute task
result = agent_executor.invoke({"input": "what is LangChain?"})
print(result['output'])
```

### Streaming Output

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain.schema import HumanMessage

# Create LLM with streaming output
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/",
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()]
)

# Send message (output will be displayed in real-time streaming)
response = llm([HumanMessage(content="Write a poem about spring")])
```

## Best Practices

> **Performance Optimization**
>
> * Enable LangChain caching mechanism
> * Use batch processing to reduce API calls
> * Set reasonable max\_tokens limits
> * Use async processing for better concurrency

> **Error Handling**
>
> * Implement retry mechanisms and exponential backoff
> * Set reasonable timeout values
> * Log detailed error information
> * Provide fallback solutions

> **Memory Management**
>
> * Use ConversationBufferWindowMemory to limit history length
> * Regularly clean unnecessary conversation history
> * Monitor memory usage
> * Implement conversation summarization

> **Security**
>
> * Use environment variables to store API keys
> * Implement input validation and filtering
> * Monitor API usage and costs
> * Rotate API keys regularly

## Getting Help

> **[Z.AI API Documentation](/api-reference)**
>
> View complete Z.AI API documentation

> **[LangChain Official Documentation](https://python.langchain.com/docs/get_started/introduction)**
>
> View LangChain official documentation and tutorials

> [!NOTE]
>
> LangChain is a rapidly evolving framework. It is recommended to update to the latest version regularly for optimal functionality and performance. Meanwhile, Z.AI will continue to optimize integration with LangChain to ensure the best compatibility and user experience.


# OpenAI Python SDK
Source: https://docs.z.ai/guides/develop/openai/python



Z.AI provides interfaces compatible with OpenAI API, which means you can use existing OpenAI SDK code and seamlessly switch to Z.AI's model services by simply modifying the API key and base URL. This compatibility allows you to:

* Quickly migrate existing OpenAI applications
* Use familiar development patterns and tools
* Enjoy the powerful capabilities of Z.AI models
* Maintain code consistency and maintainability

> [!WARNING]
>
> In some scenarios, there are still differences between Z.AI and OpenAI interfaces, but this does not affect overall compatibility.

### Core Advantages

> **Zero Learning Cost**
>
> If you are already familiar with OpenAI SDK, you can start using it immediately

> **Quick Migration**
>
> Existing OpenAI applications can be quickly migrated to Z.AI platform

> **Ecosystem Compatibility**
>
> Compatible with various tools and frameworks in the OpenAI ecosystem

> **Continuous Updates**
>
> Follow OpenAI SDK updates to maintain latest feature support

## Environment Requirements

> **Python Version**
>
> Python 3.7.1 or higher

> **OpenAI SDK**
>
> OpenAI SDK version 1.0.0 or higher

> [!WARNING]
>
> Please ensure using OpenAI SDK 1.0.0 or higher, older versions may have compatibility issues.

## Install OpenAI Python SDK

### Install using pip

```bash theme={null}
# Install or upgrade to latest version
pip install --upgrade 'openai>=1.0'

# Verify installation
python -c "import openai; print(openai.__version__)"
```

### Install using poetry

```bash theme={null}
poetry add "openai>=1.0"
```

## Quick Start

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

### Create Client

**Basic Configuration**

```python theme={null}
from openai import OpenAI

# Create Z.AI client
client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)
```

**Environment Variables**

```python theme={null}
from openai import OpenAI
import os

# Use environment variables
client = OpenAI(
    api_key=os.getenv("ZAI_API_KEY"),
    base_url="https://api.z.ai/api/paas/v4/"
)
```

**Configuration Class**

```python theme={null}
from openai import OpenAI
from dataclasses import dataclass

@dataclass
class Z.AIConfig:
    api_key: str
    base_url: str = "https://api.z.ai/api/paas/v4/"
    timeout: int = 30
    max_retries: int = 3

config = Z.AIConfig(api_key="your-api-key")
client = OpenAI(
    api_key=config.api_key,
    base_url=config.base_url,
    timeout=config.timeout,
    max_retries=config.max_retries
)
```

## Quick Start Examples

### Basic Chat

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)

completion = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "system", "content": "You are a smart and creative novelist"},
        {"role": "user", "content": "Please write a short fairy tale story as a fairy tale master"}
    ]
)

print(completion.choices[0].message.content)
```

### Streaming Response

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)

stream = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Write a poem about artificial intelligence"}
    ],
    stream=True,
    temperature=1.0
)

for chunk in stream:
    if chunk.choices[0].delta.content is not None:
        print(chunk.choices[0].delta.content, end="", flush=True)

print()  # New line
```

### Multi-turn Conversation

```python theme={null}
from openai import OpenAI

class ChatBot:
    def __init__(self, api_key: str):
        self.client = OpenAI(
            api_key=api_key,
            base_url="https://api.z.ai/api/paas/v4/"
        )
        self.conversation = [
            {"role": "system", "content": "You are a helpful AI assistant"}
        ]
    
    def chat(self, user_input: str) -> str:
        # Add user message
        self.conversation.append({"role": "user", "content": user_input})
        
        # Call API
        response = self.client.chat.completions.create(
            model="glm-5.3",
            messages=self.conversation,
            temperature=1.0
        )
        
        # Get AI response
        ai_response = response.choices[0].message.content
        
        # Add to conversation history
        self.conversation.append({"role": "assistant", "content": ai_response})
        
        return ai_response
    
    def clear_history(self):
        """Clear conversation history, keep system prompt"""
        self.conversation = self.conversation[:1]

# Usage example
bot = ChatBot("your-api-key")
print(bot.chat("Hello, please introduce yourself"))
print(bot.chat("Can you help me write code?"))
print(bot.chat("Write a Python quicksort algorithm"))
```

## Advanced Features

### Thinking Mode

In thinking mode, GLM-4.6 can solve complex reasoning problems, including mathematics, science, and logic problems.

The param `thinking.type` can be either `enabled` or `disabled`.

```python theme={null}
import os
from openai import OpenAI

client = OpenAI(api_key='your-api-key', base_url='https://api.z.ai/api/paas/v4/')
response = client.chat.completions.create(
        model='glm-5.3',
        messages=[
            {"role": "system", "content": "you are a helpful assistant"},
            {"role": "user", "content": "what is the revolution of llm?"}
        ],
        stream=True,
        extra_body={
            "thinking": {
                "type": "enabled",
            },
        }
    )
for chunk in response:
    if chunk.choices[0].delta.reasoning_content:
        print(chunk.choices[0].delta.reasoning_content, end='')
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end='')
```

### Function Calling

```python theme={null}
import json
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)

def get_weather(location: str) -> str:
    """Get weather information for specified location"""
    # This should call a real weather API
    return f"Weather in {location}: Sunny, 25°C"

# Define function description
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "Get weather information for specified location",
            "parameters": {
                "type": "object",
                "properties": {
                    "location": {
                        "type": "string",
                        "description": "Location name, e.g.: Beijing, Shanghai"
                    }
                },
                "required": ["location"]
            }
        }
    }
]

# Call conversation with functions
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "How's the weather in Beijing today?"}
    ],
    tools=tools,
    tool_choice="auto"
)

# Handle function calls
message = response.choices[0].message
if message.tool_calls:
    for tool_call in message.tool_calls:
        if tool_call.function.name == "get_weather":
            args = json.loads(tool_call.function.arguments)
            result = get_weather(args["location"])
            print(f"Function call result: {result}")
```

## Parameter Configuration

### Common Parameters

| Parameter | Type | Default | Description |
| - | - | - | - |
| model | string | Required | Model name to use |
| messages | array | Required | List of conversation messages |
| temperature | float | 0.6 | Controls output randomness (0-1) |
| top\_p | float | 0.95 | Nucleus sampling parameter (0-1) |
| max\_tokens | integer | - | Maximum output tokens |
| stream | boolean | false | Whether to use streaming output |
| stop | string/array | - | Stop generation tokens |

> [!NOTE]
>
> Note: The temperature parameter range is (0,1), do\_sample = False (temperature = 0) is not applicable in OpenAI calls.

## Best Practices

> **Performance Optimization**
>
> * Use connection pooling and session reuse
> * Set reasonable timeout values
> * Implement async calls for high concurrency
> * Cache frequently used responses

> **Cost Control**
>
> * Set reasonable max\_tokens limits
> * Use appropriate models (don't overuse powerful models)
> * Implement request deduplication
> * Monitor API usage

> **Security**
>
> * Use environment variables to store API keys
> * Implement input validation and filtering
> * Log and monitor API calls
> * Rotate API keys regularly

> **Reliability**
>
> * Implement retry mechanisms and error handling
> * Set reasonable timeout values
> * Monitor API status and response times
> * Prepare fallback solutions

## Migration Guide

### Migrating from OpenAI

If you're already using OpenAI API, migrating to Z.AI is very simple:

```python theme={null}
# Original OpenAI code
from openai import OpenAI

client = OpenAI(
    api_key="sk-...",  # OpenAI API Key
    # base_url uses default value
)

# Migrate to Z.AI, only need to modify two places
client = OpenAI(
    api_key="your-Z.AI-api-key",  # Replace with Z.AI API Key
    base_url="https://api.z.ai/api/paas/v4/"  # Add Z.AI base_url
)

# Other code remains unchanged
response = client.chat.completions.create(
    model="glm-5.3",  # Use Z.AI model
    messages=[{"role": "user", "content": "Hello!"}]
)
```

## Getting Help

> **[API Documentation](/api-reference)**
>
> View complete API interface documentation

> **[OpenAI Official Documentation](https://platform.openai.com/docs)**
>
> Refer to OpenAI official documentation for more usage

> [!NOTE]
>
> Z.AI is committed to maintaining compatibility with OpenAI API. If you encounter any issues during migration, please contact our technical support team.


# Official Python SDK
Source: https://docs.z.ai/guides/develop/python/introduction



Z.AI Python SDK is the official Python development toolkit provided by Z.AI, offering Python developers convenient and efficient AI model integration solutions.

### Core Advantages

> **Simple and Easy**
>
> Pythonic API design, comprehensive documentation and examples for quick start

> **Complete Features**
>
> Supports Z.AI's full range of models, including language, vision, image generation, etc.

> **High Performance**
>
> Async support, connection pool management, optimized network request handling

> **Type Safety**
>
> Complete type hints, IDE-friendly, reducing development errors

### Supported Features

* **💬 Chat Conversations**: Support for single-turn and multi-turn conversations, streaming and non-streaming responses
* **🔧 Function Calling**: Enable AI models to call your custom functions
* **👁️ Vision Understanding**: Image analysis, visual understanding
* **🎨 Image Generation**: Generate high-quality images from text descriptions
* **🎬 Video Generation**: Creative content generation from text to video
* **🔊 Speech Processing**: Speech-to-text, text-to-speech
* **📊 Text Embedding**: Text vectorization, supporting semantic search
* **🤖 Intelligent Assistants**: Build professional AI assistant applications
* **🛡️ Content Moderation**: Text and image content safety detection

## Technical Specifications

### Environment Requirements

* **Python Version**: Python 3.8 or higher
* **Package Manager**: pip or poetry
* **Network Requirements**: HTTPS connection support
* **API Key**: Valid Z.AI API key required

### Dependency Management

The SDK adopts a modular design, allowing you to selectively install functional modules as needed:

* **Core Module**: Basic API calling functionality
* **Async Module**: Asynchronous and concurrent processing support
* **Utility Module**: Utility tools and auxiliary functions

## Quick Start

### Environment Requirements

> **Python Version**
>
> Python 3.8 or higher

> **Package Manager**
>
> poetry (recommended), uv (recommended), pip

> [!WARNING]
>
> Supports Python 3.8, 3.9, 3.10, 3.11, 3.12 versions, cross-platform compatible with Windows, macOS, Linux

### Install SDK

#### Install using pip

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

#### Verify Installation

```python theme={null}
import zai
print(zai.__version__)
```

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

> [!NOTE]
>
> Domestic Z.AI platform uses ZaiClient
>
> Domestic API URL: [https://api.z.ai/api/paas/v4/](https://api.z.ai/api/paas/v4/)

#### Create Client

**Environment Variable**

```python theme={null}
from zai import ZaiClient
import os

# Read API Key from environment variable
client = ZaiClient(api_key=os.getenv("ZAI_API_KEY"))

# Or use directly (if environment variable is set)
client = ZaiClient()

```

**Direct Setting**

```python theme={null}
from zai import ZaiClient, ZaiClient

# Set API Key directly
client = ZaiClient(api_key="abc123.def456")

```

#### Basic Conversation

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key="your-api-key")

# Create chat completion
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Hello, please introduce yourself, Z.ai!"}
    ]
)
print(response.choices[0].message.content)
```

#### Streaming Conversation

```python theme={null}
# Create streaming chat request
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key="your-api-key")

# Create chat completion
response = client.chat.completions.create(
    model='glm-5.3',
    messages=[
        {'role': 'system', 'content': 'You are an AI writer.'},
        {'role': 'user', 'content': 'Tell a story about AI.'},
    ],
    stream=True,
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end='')
```

#### Multi-turn Conversation

```python theme={null}
from zai import ZaiClient
client = ZaiClient(api_key="your-api-key")
response = client.chat.completions.create(
    model="glm-5.3",  # Please fill in the model name you want to call
    messages=[
        {"role": "user", "content": "As a marketing expert, please create an attractive slogan for my product"},
        {"role": "assistant", "content": "Of course, to create an attractive slogan, please tell me some information about your product"},
        {"role": "user", "content": "Z.AI Open Platform"},
        {"role": "assistant", "content": "Ignite the future, Z.AI draws infinite possibilities, making innovation within reach!"},
        {"role": "user", "content": "Create a more precise and attractive slogan"}
    ],
)
print(response.choices[0].message.content)
```

### Complete Example

```python theme={null}
from zai import ZaiClient
import os

def main():
    # Initialize client
    client = ZaiClient(api_key=os.getenv("ZAI_API_KEY"))
    
    print("Welcome to Z.ai Chatbot! Type 'quit' to exit.")
    
    # Conversation history
    conversation = [
        {"role": "system", "content": "You are a friendly AI assistant"}
    ]
    
    while True:
        # Get user input
        user_input = input("You: ")
        
        if user_input.lower() == 'quit':
            break
        
        try:
            # Add user message
            conversation.append({"role": "user", "content": user_input})
            
            # Create chat request
            response = client.chat.completions.create(
                model="glm-5.3",
                messages=conversation,
                temperature=1.0,
                max_tokens=1000
            )
            
            # Get AI response
            ai_response = response.choices[0].message.content
            print(f"AI: {ai_response}")
            
            # Add AI response to conversation history
            conversation.append({"role": "assistant", "content": ai_response})
            
        except Exception as e:
            print(f"Error occurred: {e}")
    
    print("Goodbye!")

if __name__ == "__main__":
    main()
```

### Error Handling

```python theme={null}
from zai import ZaiClient
import zai

def robust_chat(message):
    client = ZaiClient(api_key="your-api-key")
    
    try:
        response = client.chat.completions.create(
            model="glm-5.3",
            messages=[{"role": "user", "content": message}]
        )
        return response.choices[0].message.content
        
    except zai.core.APIStatusError as err:
        return f"API status error: {err}"
    except zai.core.APITimeoutError as err:
        return f"Request timeout: {err}"
    except Exception as err:
        return f"Other error: {err}"

# Usage example
result = robust_chat("Hello")
print(result)
```

### Advanced Configuration

```python theme={null}
import httpx
from zai import ZaiClient

# Custom HTTP client
httpx_client = httpx.Client(
    limits=httpx.Limits(
        max_keepalive_connections=20,
        max_connections=100
    ),
    timeout=30.0
)

# Create client with custom configuration
client = ZaiClient(
    api_key="your-api-key",
    base_url="https://api.z.ai/api/paas/v4/",
    timeout=httpx.Timeout(timeout=300.0, connect=8.0),
    max_retries=3,
    http_client=httpx_client
)
```

## Advanced Features

### Function Calling

Function calling allows AI models to call functions you define to get real-time information or perform specific operations.

#### Defining and Using Functions

```python theme={null}
from zai import ZaiClient
import json

# Define functions
def get_weather(location, date=None):
    """Get weather information"""
    # Simulate weather API call
    return {
        "location": location,
        "date": date or "today",
        "weather": "sunny",
        "temperature": "25°C",
        "humidity": "60%"
    }

def get_stock_price(symbol):
    """Get stock price"""
    # Simulate stock API call
    return {
        "symbol": symbol,
        "price": 150.25,
        "change": "+2.5%"
    }

# Function descriptions
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "Get weather information for a specified location",
            "parameters": {
                "type": "object",
                "properties": {
                    "location": {
                        "type": "string",
                        "description": "Location name"
                    },
                    "date": {
                        "type": "string",
                        "description": "Date in YYYY-MM-DD format"
                    }
                },
                "required": ["location"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "get_stock_price",
            "description": "Get current stock price",
            "parameters": {
                "type": "object",
                "properties": {
                    "symbol": {
                        "type": "string",
                        "description": "Stock symbol"
                    }
                },
                "required": ["symbol"]
            }
        }
    }
]

# Use function calling
client = ZaiClient(api_key="your-api-key")

response = client.chat.completions.create(
    model='glm-5.3',
    messages=[
        {'role': 'user', 'content': 'How\'s the weather in Beijing today?'}
    ],
    tools=tools,
    tool_choice="auto"
)

# Handle function calling
if response.choices[0].message.tool_calls:
    for tool_call in response.choices[0].message.tool_calls:
        function_name = tool_call.function.name
        function_args = json.loads(tool_call.function.arguments)
        
        if function_name == "get_weather":
            result = get_weather(**function_args)
            print(f"Weather info: {result}")
        elif function_name == "get_stock_price":
            result = get_stock_price(**function_args)
            print(f"Stock info: {result}")
else:
    print(response.choices[0].message.content)
```

### Web Search Tool

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key="your-api-key")

# Use web search tool
response = client.chat.completions.create(
    model='glm-5.3',
    messages=[
        {'role': 'system', 'content': 'You are a helpful assistant.'},
        {'role': 'user', 'content': 'What is artificial intelligence?'},
    ],
    tools=[
        {
            'type': 'web_search',
            'web_search': {
                'search_query': 'What is artificial intelligence?',
                'search_result': True,
            },
        }
    ],
    temperature=0.5,
    max_tokens=2000,
)

print(response)
```

### Streaming Processing

```python theme={null}
class StreamProcessor:
    def __init__(self, client):
        self.client = client
        self.full_response = ""
    
    def stream_chat(self, messages, model="glm-5.3", callback=None):
        """Streaming chat processing"""
        stream = self.client.chat.completions.create(
            model=model,
            messages=messages,
            stream=True
        )
        
        self.full_response = ""
        for chunk in stream:
            if chunk.choices[0].delta.content is not None:
                content = chunk.choices[0].delta.content
                self.full_response += content
                
                if callback:
                    callback(content, self.full_response)
                else:
                    print(content, end="", flush=True)
        
        print()  # New line
        return self.full_response

# Usage example
processor = StreamProcessor(client)

# Custom callback function
def on_token_received(token, full_text):
    # You can implement real-time processing logic here
    print(token, end="", flush=True)

response = processor.stream_chat(
    messages=[{"role": "user", "content": "Write a Python function to calculate Fibonacci sequence"}],
    callback=on_token_received
)
```

## Getting Help

> **[GitHub Repository](https://github.com/zai-org/z-ai-sdk-python)**
>
> View source code, submit issues, contribute

> **[API Reference](/api-reference)**
>
> View complete API documentation

> **[Example Projects](https://github.com/zai-org/z-ai-sdk-python/tree/main/examples)**
>
> Browse more practical application examples

> **[Best Practices](https://github.com/zai-org/z-ai-sdk-python)**
>
> Learn best practices for SDK usage

> [!NOTE]
>
> This SDK is developed based on the latest API specifications from Z.AI, ensuring synchronization with platform features. It is recommended to regularly update to the latest version for the best experience.


# GLM-5.2
Source: https://docs.z.ai/guides/llm/glm-5.2



## Overview

**GLM-5.2** is a flagship model built for the era of long-horizon tasks. With truly usable 1M-token context, it has been tested to handle project-scale engineering context, delivering more stable long-task execution, more reliable adherence to engineering standards, and higher success rates in development scenarios. A single task can complete the full development workflow—from requirements to deployable products across multiple platforms.

> **Positioning**
>
> Flagship Foundation Model

> **Input Modalities**
>
> Text

> **Output Modalitie**
>
> Text

> **Context Length**
>
> 1M

> **Maximum Output Tokens**
>
> 128K

## Capability

> **[Thinking Mode](/guides/capabilities/thinking-mode)**
>
> Offering multiple thinking modes for different scenarios

> **[Streaming Output](/guides/capabilities/streaming)**
>
> Support real-time streaming responses to enhance user interaction experience

> **[Function Call](/guides/capabilities/function-calling)**
>
> Powerful tool invocation capabilities, enabling integration with various external toolsets

> **[Context Caching](/guides/capabilities/cache)**
>
> Intelligent caching mechanism to optimize performance in long conversations

> **[Structured Output](/guides/capabilities/struct-output)**
>
> Support for structured output formats like JSON, facilitating system integration

> **MCP**
>
> Flexibly integrate external MCP tools and data sources to expand application scenarios

## Usage

> **Project-Level Codebase Takeover: Let the Model Understand an Entire Project in One Go**
>
> This is the best starting point to experience the generational leap of GLM-5.2. It can continuously retain module boundaries, architectural constraints, API contracts, directory structures, and historical decisions, significantly reducing the sense of context fragmentation in the later stages of long-running tasks. For complex projects, the key experience is that the model does not merely read more context—it can carry forward the engineering judgments formed earlier into subsequent execution.
>
> **Recommended way to try it**: Choose a real business codebase, preferably one that includes backend, frontend or client-side code, configuration files, tests, documentation, and engineering conventions. First, ask the model to perform a technical audit:
>
> > Please read the current project and output a system architecture map, core module responsibilities, key API contracts, major data flows, core call chains, potential technical debt, and the engineering constraints that must be followed in future refactoring.

> **Long-Horizon Refactoring: Let It Run a Real Engineering Task End to End**
>
> GLM-5.2 is more stable in cross-file, multi-step, long-chain tasks. It first breaks down the goal, identifies dependencies and risks, then implements, verifies, and closes the task in stages. This makes it suitable for tasks that require continuous progress, such as module decoupling, API migration, directory restructuring, SDK adaptation, and cross-language refactoring.
>
> **Recommended way to try it**: Choose a medium-sized refactoring task, define clear boundaries, and enable `/goal` mode:
>
> > Please complete the decoupling and refactoring of the current module without changing the business logic, API signatures, or runtime behavior. First provide the execution plan, impact scope, risk boundaries, and verification method. After completion, run the necessary tests and output the verification results.

> **Production-Grade Standards Stress Test: See Whether It Can Hold the Line on Hard Engineering Constraints**
>
> GLM-5.2 shows stronger consistency in following engineering standards, especially in long-context and multi-round execution. It is better at adhering to code style, architectural boundaries, dependency constraints, build processes, testing requirements, and commit boundaries, reducing risks such as out-of-scope changes, invalid dependencies, skipped verification, or unauthorized commits.
>
> **Recommended way to try it**: Provide the model with your team’s real engineering standards, such as lint rules, build commands, testing requirements, commit conventions, and prohibited actions in `CLAUDE.md` or `Agent.md`. Then give it a real modification task:
>
> > Please strictly follow the engineering standards of the current repository. Do not introduce new dependencies, do not modify API contracts, and do not commit changes proactively. After completing the modification, run the build, lint, and tests, then report the verification results and any uncovered risks.

> **Mobile On-Device Debugging Loop: From Code Implementation to Device Validation**
>
> In mobile development scenarios, GLM-5.2 can cover client-side architecture, streaming messages, long-connection states, local state management, keyboard behavior, scrolling logic, system notifications, permission mechanisms, and background recovery. More importantly, it can use ADB, logcat, screenshots, and runtime logs to locate real-device issues, making it closely aligned with practical mobile engineering workflows.
>
> **Recommended way to try it**: Choose a real Android or Mini Program task and let the model go from implementation to validation:
>
> > Please implement a native Android client in Kotlin that connects to the existing server-side API and supports multi-session conversations, streaming messages, voice input, notifications, and reconnection after disconnection. After completion, install it on a real device using ADB, and debug it with logcat and screenshots.

> **WeChat Mini Program Development: Migrating from a Web App to a WeChat Mini Program**
>
> GLM-5.2 can handle page subpackages, custom components, page-level components, page stack management, `wx.request` wrapping and API-layer adaptation, authentication and login state maintenance (`wx.login` + custom login state), app/page/component lifecycle management, and exception handling in Mini Program development. It is suitable for testing whether the model can reorganize an existing Web page, official website, or backend capability into a runnable project that complies with Mini Program platform requirements.
>
> **Recommended way to try it**: Choose an existing Web project, specify the target technology stack — native Mini Program, Taro, or uni-app — and migrate all Web features into a Mini Program version:
>
> > Please migrate all features of the current Web project into a WeChat Mini Program. Use the \[native/Taro/uni-app] technology stack. First analyze the page structure, core user paths, backend API contracts, and platform constraints, including package size limits, domain allowlists, and HTTPS requirements. Then complete the implementation of pages, components, page navigation, and data flows. After completion, explain how to run the project, which APIs have been integrated, which features remain uncovered, and what can be optimized next.

> **Mini Game Development: From Gameplay Rules to a Playable Loop**
>
> GLM-5.2 is well suited for testing rule understanding, state machine design, level structure, scoring logic, resource loading, interaction feedback, and settlement flows in mini game development. Compared with static pages, this type of task better demonstrates the model’s understanding of complex states, user paths, and product completeness.
>
> **Recommended way to try it**: Provide a complete but not overly detailed gameplay goal, and let the model first design the rules, then implement a runnable version:
>
> > Please develop a lightweight level-based mini game. First design the core gameplay loop, state machine, level structure, scoring rules, failure and settlement logic, then implement basic features including start, pause, resume, settlement, restart, and local save. After completion, explain the project structure, verified features, and possible next-step extensions.

> **Research Reproduction: From Paper and Data to a Runnable Engineering Project**
>
> GLM-5.2 can turn the model architecture, loss functions, data pipelines, and training/inference scripts described in a paper into runnable code that aligns with the paper. It can correctly set up the model structure in one pass, maintain consistency across multiple files, and autonomously run, debug, and fix code and environment issues. What it delivers is not just code snippets, but an engineering project that can truly reproduce the paper’s reported results.
>
> **Recommended way to try it**: Pick a paper with a model and experiments, preferably one with open-source code or public metrics, and provide the paper and data to the model. See whether it can implement the model, run it successfully, and align the results with the paper:
>
> > Please reproduce the experiments based on this paper and dataset. Fill in implementation details not explicitly described in the paper. Use PyTorch to build the model architecture and loss functions, construct the data pipeline and training/inference scripts, and ensure the project runs successfully with consistency across multiple files. Autonomously identify and fix runtime issues, verify the paper’s metrics item by item until they are aligned, and explain the reproduction path, key changes, and any remaining gaps.

> **Code-to-Video Loop: From Natural-Language Ideas to a Demo-Ready Video**
>
> In Code-to-Video scenarios, GLM-5.2 can use the Remotion framework to create videos programmatically with React code, including components, parameters, and animation logic, and then render them into MP4. In simple terms, it treats video creation as writing code. It covers the full workflow from translating natural-language ideas into Remotion React code to rendering video output, enabling code-driven generation of a runnable, demo-ready video.
>
> **Recommended way to try it**: Choose a real video creation task and let the model start from a single natural-language idea, then gradually produce a renderable, playable, and iterable video:
>
> > Please create a new composition in Remotion and add a map. Start from Los Angeles, zoom the camera out while keeping LA in focus. Then draw an animated route from Los Angeles to New York and have the camera follow the route. Add one more stop to the journey — this time, we are going to Paris.

## Introducing GLM-5.2

### 1. 1M Context: Making Long-Horizon Tasks Stable and Practical

The foundation of long-horizon tasks is not having a 1M context, but making 1M context truly usable. GLM-5.2 delivers a Solid 1M lossless context and has undergone months of specialized training for long-horizon Coding Agent scenarios, covering high-value tasks such as large-scale implementation, automated research, and performance optimization.

Compared to solutions that merely extend context length, GLM-5.2 maintains more stable performance at ultra-long context, even surpassing Opus in select real-world benchmarks.

GLM-5.2 delivers state-of-the-art long-horizon coding performance among open-source models. Across FrontierSWE, PostTrainBench, and SWE-Marathon, it consistently ranks among the top models overall—trailing Opus 4.8 by just 1% on FrontierSWE, outperforming GPT-5.5 and Opus 4.7 on multiple benchmarks, and remaining the highest-ranked open-source model across all three. These results demonstrate that GLM-5.2's 1M context window translates into practical long-horizon engineering capability.

![Description](https://cdn.bigmodel.cn/markdown/17816319661261.png?attname=1.png)

### 2. Coding Capabilities Validated by Both Benchmarks and Developers

On standard coding benchmarks, GLM-5.2 is the strongest open-source model, improving on GLM-5.1 by a wide margin: 81.0 vs. 62.0 on Terminal-Bench 2.1 and 62.1 vs. 58.4 on SWE-bench Pro. It also closes much of the gap to the closed-source frontier — on Terminal-Bench 2.1 (81.0) it lands within a few points of Claude Opus 4.8 (85.0) — while staying ahead of Gemini 3.1 Pro.
![Description](https://cdn.bigmodel.cn/markdown/1781632244480plan2.png?attname=plan2.png)

Before its official release, GLM-5.2 was made available in advance to GLM Coding Plan users. Developers reported improvements mainly in the following areas:

* Stronger project-level context capacity, enabling an entire codebase to be placed within a single reasoning workflow;
* More stable long-horizon task execution, allowing complex tasks to progress continuously without easily going off track;
* More reliable adherence to production-grade engineering standards, helping enforce hard constraints in team development workflows;
* Stronger client-side and mobile engineering capabilities, going beyond app generation to support a complete on-device debugging loop.

[↗ Blog](https://z.ai/blog/glm-5.2)

## Resources

* [API Documentation](/api-reference/llm/chat-completion): Learn how to call the API.

## Quick Start

The following is a full sample code to help you onboard GLM-5.2 with ease.

**cURL**

**Basic Call**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your-api-key" \
-d '{
  "model": "glm-5.2",
  "messages": [
    {
      "role": "system",
      "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."
    },
    {
      "role": "user",
      "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack."
    }
  ],
  "thinking": {
    "type": "enabled"
  },
  "reasoning_effort": "max",
  "max_tokens": 4096,
  "temperature": 1.0
}'
```

**Streaming Call**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your-api-key" \
-d '{
  "model": "glm-5.2",
  "messages": [
    {
      "role": "system",
      "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."
    },
    {
      "role": "user",
      "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack."
    }
  ],
  "thinking": {
    "type": "enabled"
  },
  "reasoning_effort": "max",
  "stream": true,
  "max_tokens": 4096,
  "temperature": 1.0
}'
```

**Official Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai

print(zai.__version__)
```

**Basic Call**

```python theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="your-api-key")  # Your API Key

response = client.chat.completions.create(
    model="glm-5.2",
    messages=[
        {
            "role": "system",
            "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.",
        },
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
    thinking={
        "type": "enabled",
    },
    reasoning_effort="max",
    max_tokens=4096,
    temperature=1.0,
)

# Get complete response
print(response.choices[0].message)
```

**Streaming Call**

```python theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="your-api-key")  # Your API Key

response = client.chat.completions.create(
    model="glm-5.2",
    messages=[
        {
            "role": "system",
            "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.",
        },
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
    thinking={
        "type": "enabled",  # Optional: "disabled" or "enabled", default is "enabled"
    },
    reasoning_effort="max",
    stream=True,
    max_tokens=4096,
    temperature=0.6,
)

# Stream response
for chunk in response:
    if chunk.choices[0].delta.reasoning_content:
        print(chunk.choices[0].delta.reasoning_content, end="", flush=True)

    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

**Official Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

**Basic Call**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.ChatCompletionCreateParams;
import ai.z.openapi.service.model.ChatCompletionResponse;
import ai.z.openapi.service.model.ChatMessage;
import ai.z.openapi.service.model.ChatMessageRole;
import ai.z.openapi.service.model.ChatThinking;
import java.util.Arrays;

public class BasicChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI().apiKey("your-api-key").build();

        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.2")
            .messages(
                Arrays.asList(
                    ChatMessage.builder()
                        .role(ChatMessageRole.SYSTEM.value())
                        .content(
                            "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.")
                        .build(),
                    ChatMessage.builder()
                        .role(ChatMessageRole.USER.value())
                        .content(
                            "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.")
                        .build()))
            .thinking(ChatThinking.builder().type("enabled").build())
            .reasoningEffort("max")
            .maxTokens(4096)
            .temperature(1.0f)
            .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        // Get response
        if (response.isSuccess()) {
            Object reply = response.getData().getChoices().get(0).getMessage();
            System.out.println("AI Response: " + reply);
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

**Streaming Call**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.ChatCompletionCreateParams;
import ai.z.openapi.service.model.ChatCompletionResponse;
import ai.z.openapi.service.model.ChatMessage;
import ai.z.openapi.service.model.ChatMessageRole;
import ai.z.openapi.service.model.ChatThinking;
import ai.z.openapi.service.model.Delta;
import java.util.Arrays;

public class StreamingChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI().apiKey("your-api-key").build();

        // Create streaming chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.2")
            .messages(
                Arrays.asList(
                    ChatMessage.builder()
                        .role(ChatMessageRole.SYSTEM.value())
                        .content(
                            "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.")
                        .build(),
                    ChatMessage.builder()
                        .role(ChatMessageRole.USER.value())
                        .content(
                            "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.")
                        .build()))
            .thinking(ChatThinking.builder().type("enabled").build())
            .reasoningEffort("max")
            .stream(true) // Enable streaming output
            .maxTokens(4096)
            .temperature(1.0f)
            .build();

        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        if (response.isSuccess()) {
            response.getFlowable()
                .subscribe(
                    // Process streaming message data
                    data -> {
                        if (data.getChoices() != null && !data.getChoices().isEmpty()) {
                            Delta delta = data.getChoices().get(0).getDelta();
                            System.out.print(delta + "\n");
                        }
                    },
                    // Process streaming response error
                    error -> System.err.println("\nStream error: " + error.getMessage()),
                    // Process streaming response completion event
                    () -> System.out.println("\nStreaming response completed"));
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

**OpenAI Python SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
pip install --upgrade 'openai>=1.0'
```

**Verify Installation**

```python theme={null}
python -c "import openai; print(openai.__version__)"
```

**Usage Example**

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/",
)

completion = client.chat.completions.create(
    model="glm-5.2",
    messages=[
        {"role": "system", "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."},
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
)

print(completion.choices[0].message.content)
```


# GLM-5.3
Source: https://docs.z.ai/guides/llm/glm-5.3



> [!TIP]
>
> GLM-5.3 is now available to all GLM Coding Plan users, with 50% better coding performance! [**Subscribe now**](https://z.ai/subscribe)

## Overview

**GLM-5.3** is Z.ai's latest flagship model, delivering comprehensive advancements in complex software engineering and agent capabilities. It uses the same base model as GLM-5.2, with all improvements driven by post-training. Compared with GLM-5.2, it delivers significantly stronger performance in complex programming and long-horizon tasks:

**Stronger Coding Capability**

GLM-5.3 delivers a significant improvement in coding capabilities, achieving a 50% performance gain over GLM-5.2 on Z.ai Code Bench. It also reaches state-of-the-art (SOTA) performance among open-source models on public benchmarks, including Terminal Bench 3.0 and Agents’ Last Exam (CLI).

**Emergent Cyber Capability**

As the scale of post-training continues to expand, the model’s cybersecurity capabilities have improved at a rate that exceeds expectations. GLM-5.3 achieves the best performance to date on the CyberGym vulnerability discovery benchmark. Moreover, the deeper it progresses along the vulnerability exploitation chain, the more pronounced its gains over GLM-5.2 become, with scores on vulnerability exploitation benchmarks exceeding twice those of GLM-5.2.

[↗ Blog](https://z.ai/blog/glm-5.3)

## Feature Changes

GLM-5.3 currently supports text-only inputs, with a 1M-token context window and a maximum output length of 128K tokens.
GLM-5.3 always operates with reasoning enabled and supports three reasoning effort levels: `low`, `high`, and `max`. Disabling reasoning is no longer supported. The reasoning parameters are described below:

| Parameter | Values | Default | Description |
| :- | :- | :- | :- |
| `thinking.type` | `enabled` | `enabled` | Supports reasoning only; disabling reasoning is not supported. |
| `reasoning_effort` | `low`, `high`, `max` | `max` | `low` – Lightweight Reasoning; `high` – Enhanced Reasoning; `max` – Deep Reasoning |

> [!NOTE]
>
> Migration Notice: If your application currently uses **`thinking.type: "disabled"`**, please change it to **`enabled`** and set **`reasoning_effort`** to **`low`** before updating the model ID to **`glm-5.3`**. Otherwise, the request will fail.

For complex tasks such as coding, we recommend using `max`. An example is provided below:

```json theme={null}
{
  "model": "glm-5.3",
  "thinking": { "type": "enabled" },
  "reasoning_effort": "max"
}
```

For a detailed introduction to GLM-5.3’s advanced reasoning capabilities, please refer to the documentation: [Deep Thinking](https://docs.z.ai/guides/capabilities/thinking)

## How to Use

#### **Model API**

The supported protocols and endpoints are as follows:

| Protocol | Base URL |
| :- | :- |
| OpenAI Chat Completion Protocol | `	https://api.z.ai/api/coding/paas/v4` |
| OpenAI Response Protocol | `	https://api.z.ai/api/v1` |
| Anthropic Message Protocol | `https://api.z.ai/api/anthropic` |

[↗ API Documentation](https://docs.z.ai/api-reference/llm/chat-completion)：Learn how to call the API

> [!NOTE]
>
> If you have previously subscribed to a GLM Coding Plan, including an expired subscription, you can currently access the model API only through the OpenAI Chat Completion-compatible protocol. We will continue to improve and optimize this in upcoming iterations.

#### **GLM Coding Plan**

GLM-5.3 is now fully available, and you can use it with your preferred coding agents. [Guide](https://docs.z.ai/devpack/latest-model)

The new GLM Coding Plan adopts a points-based quota system, providing transparent and predictable usage. Model calls made during off-peak hours, including all day on weekends, consume only 50% of the standard points.
Subscribe now: [Individual Plan](https://z.ai/subscribe)、[Team Plan](https://z.ai/subscribe?plantype=team).

## Capability Support

* [Thinking Mode](https://docs.z.ai/guides/capabilities/thinking-mode)：Provides multiple thinking modes to accommodate different task requirements
* [Streaming Messages](https://docs.z.ai/guides/capabilities/streaming)：Supports real-time streaming responses for an enhanced user interaction experience
* [Function Calling](https://docs.z.ai/guides/capabilities/function-calling)：Provides powerful tool-calling capabilities and supports integration with a wide range of external tools
* [Context Caching](https://docs.z.ai/guides/capabilities/cache)：Uses an intelligent caching mechanism to optimize performance in long-context conversations
* [Structured Output](https://docs.z.ai/guides/capabilities/struct-output)：Supports structured output formats such as JSON for seamless system integration

## Key Advancements

GLM-5.3 delivers significant improvements across complex software engineering, terminal operations, and a broader range of real-world Agent tasks.

![Description](https://cdn.bigmodel.cn/markdown/1786683281203image.png?attname=image.png)

#### Stronger Coding

For GLM-5.3, we pushed environment scaling toward tasks that look less like coding exercises and more like real units of expert work. The environments now cover a much broader range of production workflows, with tasks designed around how engineering and research work is actually carried out in practice. Some represent several days of work for an experienced engineer.

In an ML infrastructure task, for example, the model may be given the same working environment as an engineer, with access to compute clusters, storage systems, internal documentation, codebases, and experiment results. It must diagnose bottlenecks across the training stack, implement optimizations, run experiments, and deliver a measurable end-to-end speedup while preserving correctness. Training on environments at this level pushes the model toward taking ownership of substantial work end to end, rather than relying on users to decompose the problem and supervise each step.

As agent capability improves, much of the difficulty in scaling post-training moves from the model to the environment. A useful task environment has to be executable, verifiable, and close to real professional work — and we need many of them, not a handful of hand-built ones. To scale this process, we built pipelines that synthesize environments end to end, and for a subset of tasks, the RL reward signal as well. Research agents collect task patterns from real work and turn them into runnable long-horizon environments with multi-step dependencies and hidden state; a judge agent then attempts each task to verify that it is actually solvable. Verifiers are synthesized without access to the reference solution, while solver trajectories are used to discover and close reward shortcuts. A verifier that passes oracle, no-op, and unsolved-state checks produces a binary reward reliable enough to train on directly.

It carries over the RL strategies introduced in GLM-5.2, including SAO with compaction, which helps these gains hold on long-horizon tasks rather than only on short ones. The effect shows up across both coding and general agent tasks. GLM-5.3 improves from 4.6 to 28.3 on Terminal-Bench 3.0, from 46.2 to 66.9 on DeepSWE v1.1, and from 23.8 to 28.5 on Agents' Last Exam. These pipelines still require a meaningful amount of human-in-the-loop work; making environment generation and verification more autonomous is one of the next steps.

Beyond public benchmarks, we introduce Z.ai Code Bench, an in-house benchmark designed to evaluate coding agents under realistic user scenarios. It covers diverse task categories and places agents in complex local development environments. At different effort levels, we evaluate agents along two dimensions: end-to-end task completion rate and fine-grained checklist accuracy. As a private benchmark, Z.ai Code Bench also reduces the risk of contamination from public test sets and gives us a more faithful measure of real-world user experience.

![Description](https://cdn.bigmodel.cn/markdown/1786683377150image.png?attname=image.png)

As shown in the figure, GLM-5.3 improves both performance and token efficiency. It delivers markedly stronger agentic coding results than GLM-5.2 at every effort level while consuming fewer output tokens. At Max effort, GLM-5.3 reaches 34.5% at roughly 75K output tokens per task, compared with 23.4% at 96K for GLM-5.2. The same shift holds against closed models. At High effort, GLM-5.3 reaches 31.4% at around 50K output tokens, surpassing Claude Opus 4.8 at 29.5% with 120K. GLM-5.3 remains behind Claude Fable 5, which reaches 39.5% at Max effort.

#### Emergent Cyber Capability

As part of post-training, we introduced vulnerability discovery data and environments into the training mix. We expected this to make the model better at finding and reasoning about vulnerabilities. What surprised us was how quickly the capability continued to develop as training scaled. GLM-5.3 did not simply become better at identifying isolated flaws: it began to reason across multiple stages of exploitation, forming coherent plans for complete exploitation chains.

We evaluate GLM-5.3 across three benchmarks covering different stages of vulnerability analysis and exploitation.

![Description](https://cdn.bigmodel.cn/markdown/1786683405531image.png?attname=image.png)

We evaluate GLM-5.3 across three benchmarks covering different stages of vulnerability analysis and exploitation. On CyberGym, which starts from white-box source code and tests whether the model can identify and validate vulnerabilities by triggering faults, GLM-5.3 scores 84.5%, up from GLM-5.2's 77.2% — the best result on the benchmark, ahead of Mythos 5 (83.8%) and GPT-5.6 Sol (83.6%). On ExploitBench, which requires deeper reasoning about real vulnerabilities and their exploitation, GLM-5.3 reaches 54.4%, more than doubling GLM-5.2's 24.4%, while Mythos 5 and GPT-5.6 Sol score 78.0% and 76.5%, respectively. On ExploitGym, which measures how many exploitation tasks a model can complete under time-normalized budgets, GLM-5.3 completes 105 tasks within two hours and 130 within six hours, compared with 29 and 39 for GLM-5.2; budgets are normalized across models using per-model throughput figures, detailed in the footnotes. Mythos 5 remains well ahead at 181 and 247 tasks.

The pattern across the three is consistent: the further up the exploitation chain a benchmark sits, the larger the gain from GLM-5.2 — and also the wider the remaining gap to the closed frontier. Capability is growing fastest exactly where we are furthest behind.

We then tested whether these capabilities transfer beyond controlled benchmarks. Since GLM-5.2, we have been working with several security teams in China to run our models against real-world codebases. After expert review, screening, and deduplication, the model identified 2,436 vulnerabilities across 269 projects, including 1,097 medium-to-high severity issues. The findings span system kernels, operating systems, browser engines, open-source infrastructure, web applications, and network protocols. Many had remained unnoticed for years or even decades, with the oldest dating back roughly 40 years.

This work has since grown into an ongoing disclosure effort. We built the Z.ai Security Disclosure Ledger to maintain a public record of the findings as they move through the disclosure process. The ledger is continuously updated as new vulnerabilities are reviewed and disclosed, distinguishing issues that have already been made public from those that are still under disclosure. For disclosed issues, it records information including the affected project, severity, CVE where available, and how long the vulnerability had remained in the codebase.

## Quick Start

The following is a full sample code to help you onboard GLM-5.3 with ease.

**cURL**

**Basic Call**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your-api-key" \
-d '{
  "model": "glm-5.3",
  "messages": [
    {
      "role": "system",
      "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."
    },
    {
      "role": "user",
      "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack."
    }
  ],
  "thinking": {
    "type": "enabled"
  },
  "reasoning_effort": "max",
  "max_tokens": 4096,
  "temperature": 1.0
}'
```

**Streaming Call**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your-api-key" \
-d '{
  "model": "glm-5.3",
  "messages": [
    {
      "role": "system",
      "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."
    },
    {
      "role": "user",
      "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack."
    }
  ],
  "thinking": {
    "type": "enabled"
  },
  "reasoning_effort": "max",
  "stream": true,
  "max_tokens": 4096,
  "temperature": 1.0
}'
```

**Official Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai

print(zai.__version__)
```

**Basic Call**

```python theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="your-api-key")  # Your API Key

response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.",
        },
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
    thinking={
        "type": "enabled",  # Optional: "enabled" only, reasoning is always enabled
    },
    reasoning_effort="max",
    max_tokens=4096,
    temperature=1.0,
)

# Get complete response
print(response.choices[0].message)
```

**Streaming Call**

```python theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="your-api-key")  # Your API Key

response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.",
        },
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
    thinking={
        "type": "enabled",  # Optional: "enabled" only, reasoning is always enabled
    },
    reasoning_effort="max",
    stream=True,
    max_tokens=4096,
    temperature=1.0,
)

# Stream response
for chunk in response:
    if chunk.choices[0].delta.reasoning_content:
        print(chunk.choices[0].delta.reasoning_content, end="", flush=True)

    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

**Official Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

**Basic Call**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.ChatCompletionCreateParams;
import ai.z.openapi.service.model.ChatCompletionResponse;
import ai.z.openapi.service.model.ChatMessage;
import ai.z.openapi.service.model.ChatMessageRole;
import ai.z.openapi.service.model.ChatThinking;
import java.util.Arrays;

public class BasicChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI().apiKey("your-api-key").build();

        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(
                Arrays.asList(
                    ChatMessage.builder()
                        .role(ChatMessageRole.SYSTEM.value())
                        .content(
                            "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.")
                        .build(),
                    ChatMessage.builder()
                        .role(ChatMessageRole.USER.value())
                        .content(
                            "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.")
                        .build()))
            .thinking(ChatThinking.builder().type("enabled").build())
            .reasoningEffort("max")
            .maxTokens(4096)
            .temperature(1.0f)
            .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        // Get response
        if (response.isSuccess()) {
            Object reply = response.getData().getChoices().get(0).getMessage();
            System.out.println("AI Response: " + reply);
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

**Streaming Call**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.ChatCompletionCreateParams;
import ai.z.openapi.service.model.ChatCompletionResponse;
import ai.z.openapi.service.model.ChatMessage;
import ai.z.openapi.service.model.ChatMessageRole;
import ai.z.openapi.service.model.ChatThinking;
import ai.z.openapi.service.model.Delta;
import java.util.Arrays;

public class StreamingChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI().apiKey("your-api-key").build();

        // Create streaming chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(
                Arrays.asList(
                    ChatMessage.builder()
                        .role(ChatMessageRole.SYSTEM.value())
                        .content(
                            "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks.")
                        .build(),
                    ChatMessage.builder()
                        .role(ChatMessageRole.USER.value())
                        .content(
                            "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.")
                        .build()))
            .thinking(ChatThinking.builder().type("enabled").build())
            .reasoningEffort("max")
            .stream(true) // Enable streaming output
            .maxTokens(4096)
            .temperature(1.0f)
            .build();

        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        if (response.isSuccess()) {
            response.getFlowable()
                .subscribe(
                    // Process streaming message data
                    data -> {
                        if (data.getChoices() != null && !data.getChoices().isEmpty()) {
                            Delta delta = data.getChoices().get(0).getDelta();
                            System.out.print(delta + "\n");
                        }
                    },
                    // Process streaming response error
                    error -> System.err.println("\nStream error: " + error.getMessage()),
                    // Process streaming response completion event
                    () -> System.out.println("\nStreaming response completed"));
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

**OpenAI Python SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
pip install --upgrade 'openai>=1.0'
```

**Verify Installation**

```python theme={null}
python -c "import openai; print(openai.__version__)"
```

**Usage Example**

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/",
)

completion = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "system", "content": "You are a senior full-stack software engineer, proficient in frontend development, backend architecture design, and modern web technology stacks."},
        {
            "role": "user",
            "content": "Design and build a personal blog website for me, including a homepage, article list page, and article detail page, using React + Node.js technology stack.",
        },
    ],
)

print(completion.choices[0].message.content)
```


# Core Parameters
Source: https://docs.z.ai/guides/overview/concept-param



> [!TIP]
>
> When interacting with models, you can control the model's output by adjusting different parameters to meet the needs of various scenarios. Understanding these core parameters will help you better utilize the model's capabilities.

## Quick Reference

| Parameter | Type | Default Value | Description |
| :- | :- | :- | :- |
| [do\_sample](#do_sample) | Boolean | `true` | Whether to sample the output to increase diversity. |
| [temperature](#temperature) | Float | (Model dependent) | Controls the randomness of output, higher values are more random. |
| [top\_p](#top_p) | Float | (Model dependent) | Controls diversity through nucleus sampling, recommended to use either this or `temperature`. |
| [max\_tokens](#max_tokens) | Integer | (Model dependent) | Limits the maximum number of tokens generated in a single call. |
| [stream](#stream) | Boolean | `false` | Whether to return responses in streaming mode. |
| [thinking](#thinking) | Object | `{"type": "enabled"}` | Whether to enable chain-of-thought deep thinking, only supported by `GLM-4.6` and above. |
| [reasoning\_effort](#reasoning_effort) | String | `max` `high` `low` | Controls the model's reasoning effort level, only supported by `GLM-5.2` and above. |

***

## Parameter Details

### do\_sample

`do_sample` is a boolean value (`true` or `false`) that determines whether to sample the model's output.

* `true` (default): Performs random sampling based on the probability distribution of each token, increasing text diversity and creativity. Suitable for content creation, dialogue, and other scenarios.
* `false`: Uses a greedy strategy, always selecting the token with the highest probability. Provides high deterministic output, suitable for scenarios requiring precise, factual answers.

Best Practices:

* Set to `false` when you need reproducible, deterministic output.
* Set to `true` when you want the model to generate more diverse and interesting content, and use it in combination with `temperature` or `top_p`.

### temperature

The `temperature` parameter controls the randomness of the model's output.

* Lower values (e.g., 0.2): Make the probability distribution more "sharp", resulting in more deterministic and conservative output.
* Higher values (e.g., 0.8): Make the probability distribution more "flat", resulting in more random and diverse output.

Best Practices:

* For scenarios requiring rigor and factual accuracy (such as knowledge Q\&A), it's recommended to use lower `temperature`.
* For scenarios requiring creativity (such as content creation), you can try higher `temperature`.
* It's recommended to use only one of `temperature` and `top_p`.

### top\_p

`top_p` (nucleus sampling) controls diversity by sampling from the smallest set of tokens whose cumulative probability exceeds the threshold.

* Lower values (e.g., 0.2): Limit the sampling range, resulting in more deterministic output.
* Higher values (e.g., 0.9): Expand the sampling range, resulting in more diverse output.

Best Practices:

* If you want to achieve some diversity while ensuring content quality, `top_p` is a good choice (recommended values 0.8-0.95).
* It's generally not recommended to modify both `temperature` and `top_p` simultaneously.

### max\_tokens

`max_tokens` is used to limit the maximum number of tokens the model can generate in a single call. GLM-4.6 supports a maximum output length of 128K, and it's recommended to set it to no less than 1024. Tokens are the basic units of text, typically 1 token equals approximately 0.75 English words or 1.5 Chinese characters. Setting an appropriate max\_tokens can control response length and cost, avoiding overly long outputs. If the model completes its answer before reaching the max\_tokens limit, it will naturally end; if it reaches the limit, the output may be truncated.

* Purpose: Prevents generating overly long text and controls API call costs.
* Note: `max_tokens` limits the length of generated content, not including input.

Best Practices:

* Set `max_tokens` reasonably according to your application scenario. If you need short answers, you can set it to a smaller value (e.g., 50).

Default `max_tokens` and maximum supported `max_tokens` for each model:

| Model Code | Default max\_tokens | Maximum max\_tokens |
| :- | :-: | :-: |
| glm-5.3 | 65536 | 131072 |
| glm-5.3-flash | 65536 | 131072 |
| glm-5.2 | 65536 | 131072 |
| glm-5.1 | 65536 | 131072 |
| glm-5 | 65536 | 131072 |
| glm-4.7 | 65536 | 131072 |
| glm-4.6 | 65536 | 131072 |

### stream

`stream` is a boolean value used to control the API's response method.

* `false` (default): Returns the complete response at once, simple to implement but with long waiting times.
* `true`: Returns content in streaming (SSE) mode, significantly improving the experience of real-time interactive applications.

Best Practices:

* For chatbots, real-time code generation, and other applications, it's strongly recommended to set this to `true`.

### thinking

The `thinking` parameter controls whether the model enables "Chain of Thought" for deeper thinking and reasoning.

* Type: Object
* Supported Models: `GLM-4.6` and above

Properties:

* `type` (string):
  * `enabled` (default): Enable chain of thought. `GLM-5.3` and `GLM-5.3-FLASH` use forced thinking. For other models, `GLM-5.2`, `GLM-5.1`, `GLM-5`, `GLM-4.6`, while `GLM-4.7` use forced thinking.
  * `disabled`: Disable chain of thought.

Best Practices:

* It's recommended to enable this when you need the model to perform complex reasoning and planning.
* For simple tasks, you can disable it to get faster responses.

### reasoning\_effort

The `reasoning_effort` parameter is used to control the model's reasoning effort level when chain-of-thought thinking is enabled.

* Type: String
* Supported Models: `GLM-5.2` and above
* Allowed values: `max`, `high`, `low`
* Default: `max`

Parameter values:

* `low`: Light thinking
* `high`: Enhanced thinking
* `max`: Deep thinking (default)

Notes:

* The `GLM-5.2` model supports `max`, `xhigh`, `high`, `medium`, `low`, `minimal`, `none`; passing `none` or `minimal` will cause the model to skip thinking; `low` and `medium` will be mapped to `high`; `xhigh` will be mapped to `max`.
* This parameter only takes effect when `thinking.type` is set to `enabled`.
* The `GLM-5.3` and `GLM-5.3-FLASH` models only support `max`, `high`, `low`.

***

## Related Concepts

> **Token Usage Calculation**
>
> Tokens are the basic units for model text processing. Usage calculation includes both input and output parts.
>
> * **Input Token Count:** The number of tokens contained in the text you send to the model.
> * **Output Token Count:** The number of tokens contained in the text generated by the model.
> * **Total Token Count:** The sum of input and output, usually used as the billing basis.
>
> You can call the `tokenizer` API to estimate the token count of text.

> **Maximum Output Tokens**
>
> Maximum Output Tokens refers to the maximum number of tokens a model can generate in a single request. It's different from the `max_tokens` parameter - `max_tokens` is the upper limit you set in your request, while Maximum Output Tokens is the architectural limitation of the model itself.
>
> For example, a model's context window might be 8k tokens, but its maximum output capability might be limited to 4k tokens.

> **Context Window**
>
> The Context Window refers to the total number of tokens a model can process in a single interaction, including all tokens from both **input text** and **generated text**.
>
> * **Importance:** The context window determines how much historical information the model can "remember". If the total length of input and expected output exceeds the model's context window, the model will be unable to process it.
> * **Note:** Different models have different context window sizes. When conducting long conversations or processing long documents, special attention should be paid to context window limitations.

> **Concurrency Limits**
>
> Concurrency refers to the number of API requests you can initiate simultaneously. This is set by the platform to ensure service stability and fair resource allocation.
>
> * **Limits:** Different users or subscription plans may have different concurrency quotas.
> * **Overages:** If you exceed the concurrency limit, new requests may fail or need to wait in queue.
>
> If your application requires high concurrency processing, please check your account limits or contact platform support.


# Migrate to GLM-5.3
Source: https://docs.z.ai/guides/overview/migrate-to-glm-new



This guide explains how to migrate your calls from GLM-5.2 or other earlier models to our latest flagship model GLM-5.3, covering samminpling parameter differences, streaming tool calls, reasoning\_effort parameter, and other key points.

## GLM-5.3 Features

* Context and output: Maximum context 1M, maximum output 128K.
* Support for streaming output during tool calling process (`tool_stream=true`), real-time retrieval of tool call parameters.
* Forced deep thinking by default (`thinking={ type: "enabled" }`), with the reasoning depth controlled by `reasoning_effort`.
* Support for the `reasoning_effort` parameter to control the reasoning effort level of the model when chain-of-thought thinking is enabled.
* Superior code performance and advanced reasoning capabilities.

## Migration Checklist

* [ ] Update model identifier to `glm-5.3`
* [ ] Sampling parameters: `temperature` default value `1.0`, `top_p` default value `0.95`, recommend choosing only one for tuning
* [ ] Deep thinking: Forced `thinking={ type: "enabled" }` for complex reasoning/coding
* [ ] Control reasoning effort: Configure `reasoning_effort` to decide between `low` (light thinking), `high` (enhanced thinking), or `max` (deep thinking, default)
* [ ] Streaming response: Enable `stream=true` and properly handle `delta.reasoning_content` and `delta.content`
* [ ] Streaming tool calls: Enable `stream=true` and `tool_stream=true` and stream-concatenate `delta.tool_calls[*].function.arguments`
* [ ] Maximum output and context: Set `max_tokens` appropriately (GLM-5.3 maximum output 128K, context 1M)
* [ ] Development environment verification: Conduct use case testing and regression, focus on randomness, latency, parameter completeness in tool streams

## Start Migration

### 1. Update Model Identifier

* Update `model` to `glm-5.3`.

```python theme={null}
resp = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "Briefly describe the advantages of GLM-5.3"}]
)
```

### 2. Update Sampling Parameters

* `temperature`: Controls randomness; higher values are more divergent, lower values are more stable.
* `top_p`: Controls nucleus sampling; higher values expand candidate set, lower values converge candidate set.
* `temperature` defaults to `1.0`, `top_p` defaults to `0.95`, not recommended to adjust both simultaneously.

```python theme={null}
# Plan A: Use temperature (recommended)
resp = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "Write a more creative brand introduction"}],
    temperature=1.0
)

# Plan B: Use top_p
resp = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "Generate more stable technical documentation"}],
    top_p=0.8
)
```

### 3. Deep Thinking (Forced)

* GLM-5.3 continues to support deep thinking capability, which is forced on by default and returns an error if disabled.
* Recommended to enable for complex reasoning and coding tasks:

```python theme={null}
resp = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "Design a three-tier microservice architecture for me"}],
    thinking={"type": "enabled"}
)
```

* Control reasoning effort, default `max` for deep reasoning:

```python theme={null}
resp = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "Design a three-tier microservice architecture for me"}],
    thinking={"type": "enabled"},
    reasoning_effort="max"
)
```

### 4. Streaming Output and Tool Calls (Optional)

* GLM-5.3 supports real-time streaming construction and output during tool calling process, disabled by default (`False`), requires enabling both:
  * `stream=True`: Enable streaming output for responses
  * `tool_stream=True`: Enable streaming output for tool call parameters

```python theme={null}
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[{"role": "user", "content": "How's the weather in Beijing"}],
    tools=[
        {
            "type": "function",
            "function": {
                "name": "get_weather",
                "description": "Get current weather conditions for a specified location",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "location": {"type": "string", "description": "City, eg: Beijing, Shanghai"},
                        "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
                    },
                    "required": ["location"]
                }
            }
        }
    ],
    stream=True,
    tool_stream=True,
)

# Initialize streaming collection variables
reasoning_content = ""
content = ""
final_tool_calls = {}
reasoning_started = False
content_started = False

# Process streaming response
for chunk in response:
    if not chunk.choices:
        continue

    delta = chunk.choices[0].delta

    # Streaming reasoning process output
    if hasattr(delta, 'reasoning_content') and delta.reasoning_content:
        if not reasoning_started and delta.reasoning_content.strip():
            print("\n🧠 Thinking Process:")
            reasoning_started = True
        reasoning_content += delta.reasoning_content
        print(delta.reasoning_content, end="", flush=True)

    # Streaming answer content output
    if hasattr(delta, 'content') and delta.content:
        if not content_started and delta.content.strip():
            print("\n\n💬 Answer Content:")
            content_started = True
        content += delta.content
        print(delta.content, end="", flush=True)

    # Streaming tool call information (parameter concatenation)
    if delta.tool_calls:
        for tool_call in delta.tool_calls:
            idx = tool_call.index
            if idx not in final_tool_calls:
                final_tool_calls[idx] = tool_call
                final_tool_calls[idx].function.arguments = tool_call.function.arguments
            else:
                final_tool_calls[idx].function.arguments += tool_call.function.arguments

# Output final tool call information
if final_tool_calls:
    print("\n📋 Function Calls Triggered:")
    for idx, tool_call in final_tool_calls.items():
        print(f"  {idx}: Function Name: {tool_call.function.name}, Parameters: {tool_call.function.arguments}")
```

See: [Tool Streaming Output Documentation](/guides/tools/stream-tool)

### 5. Testing and Regression

> First verify in development environment that post-migration calls are stable, focus on:

* Whether responses meet expectations, whether there's excessive randomness or excessive conservatism in output
* Whether tool streaming construction and output work normally
* Latency and cost in long context and deep thinking scenarios

## More Resources

> **[Concept Parameters](/guides/overview/concept-param)**
>
> Common model parameter concepts and sampling recommendations

> **[Tool Streaming Output](/guides/tools/stream-tool)**
>
> View tool streaming output usage details

> **[API Reference](/api-reference/introduction)**
>
> View complete API documentation

> **[Technical Support](https://z.ai/consultation)**
>
> Get technical support and help


# Overview
Source: https://docs.z.ai/guides/overview/overview



> [!NOTE]
>
> Z.AI offers a variety of models and agents to meet the needs of different scenarios. Choosing the right model can help you complete tasks more efficiently.

## Latest Models

> **[GLM-5.3](/guides/llm/glm-5.3)**
>
> **Flagship Model**
>
> * Open-source SOTA coding capabilities, emerging cybersecurity capabilities.

> **[GLM-5.3-Flash](/guides/vlm/glm-5.3-flash)**
>
> **Native Multimodal Model**
>
> * Matching Claude Opus 4.8 in intelligence score at just 1/40 the price.

## Featured Models

The following are Z.ai’s core models, covering a comprehensive range of use cases, including text generation, multimodal understanding.

| Model | Strength | Language | Context | Resource |
| :- | :- | :- | :- | :- |
| GLM-5.3-Flash | Delivering frontier intelligence at radically lower cost | English & Chinese | 1M | [Guide](/guides/vlm/glm-5.3-flash) [API Reference](/api-reference/llm/chat-completion) |
| GLM-5.3-FlashX | Delivering inference speeds of 200 tokens/s | English & Chinese | 1M | [Guide](/guides/vlm/glm-5.3-flash) [API Reference](/api-reference/llm/chat-completion) |
| GLM-5.3 | Claude Fable 5-level coding and agent capabilities  Stronger in long-horizon, complex tasks | English & Chinese | 1M | [Guide](/guides/llm/glm-5.3) [API Reference](/api-reference/llm/chat-completion) |

> [!TIP]
>
> If you need to get pricing information, please go directly to [Pricing](/guides/overview/pricing).

> **Other Models, Agents and Tools**
>
> To help you find the best fit for your use case, we've created a table outlining the core features and strengths of each model in the Z.AI family.
>
> ### Text Models
>
> Our model matrix includes text models with built-in reasoning capabilities, as well as vision-language models (VLMs) that extend the same reasoning power to multimodal understanding.
>
> | Model | Strength | Language | Context | Resource |
> | :- | :- | :- | :- | :- |
> | GLM-5.3-FlashX-Legal | A Model Optimized for Legal Applications Enabling End-to-End Workflows for Legal Research | English & Chinese | 1M | [Guide](/guides/vlm/glm-5.3-flash#glm-5-3-flashx-legal) [API Reference](/api-reference/llm/chat-completion) |
> | GLM-5.2 | Open-source SOTA coding, from generation to delivery | English & Chinese | 1M | [Guide](/guides/llm/glm-5.2) [API Reference](/api-reference/llm/chat-completion) |
> | GLM-5.1 | Coding proficiency aligned with Opus 4.6 Ability to work independently and consistently for up to 8 hours on a single task | English & Chinese | 200K | [Guide](/guides/llm/glm-5.1) [API Reference](/api-reference/llm/chat-completion) |
> | GLM-5 | Agentic Long-Term Planning and Execution Backend refactoring and in-depth debugging | English & Chinese | 200K | [Guide](/guides/llm/glm-5) [API Reference](/api-reference/llm/chat-completion) |
> | GLM-4.7 | SOTA Performance Optimized Agentic Coding | English & Chinese | 200K | [Guide](/guides/llm/glm-4.7) [API Reference](/api-reference/llm/chat-completion) |
> | GLM-4.6 | Strong Coding More Versatile | English & Chinese | 200K | [Guide](/guides/llm/glm-4.6) [API Reference](/api-reference/llm/chat-completion) |
>
> ### Built-in Tools
>
> A suite of built-in tools designed to streamline workflows and boost productivity.
>
> | Tool | Capability |
> | :- | :- |
> | Web Search | Provide real-time, concise, direct answers Accurately parse complex HTML and converts it into clean Markdown or JSON |
>
> ### Image Generation Models
>
> Image Generation Models learn from massive image data to automatically generate high-quality images from text.
>
> | Model | Strength | Language | Resolution | Resource |
> | :- | :- | :- | :- | :- |
> | GLM-Image | Stronger in complex instruction and knowledge-intensive scenarios - Open-source SOTA in text rendering | English & Chinese | multiple resolutions | [Guide](/guides/image/glm-image) [API Reference](/api-reference/image/generate-image) |
>
> ### Agents
>
> A set of ready-made agents empower users to create and communicate effortlessly.
>
> | Tool | Capability | Resource |
> | :- | :- | :- |
> | GLM Slide/Poster Agent(beta) | Combine content generation with professional design | [Guide](/guides/agents/slide) |
> | General-Purpose Translation | Support 40+ languages, flexible strategies, and terminology customization | [Guide](/guides/agents/slide) |
> | Popular Special Effects Video Templates | Special effects video templates like French\_Kiss, BodyShake, and Sexy\_Me | [Guide](/guides/agents/slide) |


# Pricing
Source: https://docs.z.ai/guides/overview/pricing

This page provides pricing information for Z.AI’s models and tools. All prices are in USD.

## Models

### Latest Models

Prices per 1M tokens.

| Model | Input | Cached Input | Cached Input Storage | Output |
| :- | :- | :- | :- | :- |
| GLM-5.3 | \$1.4 | \$0.26 | Limited-time Free | \$4.4 |
| GLM-5.3-Flash | \$0.15 | \$0.03 | Limited-time Free | \$0.50 |
| GLM-5.3-FlashX | \$0.37 | \$0.075 | Limited-time Free | \$1.25 |

### Text Models

Prices per 1M tokens.

| Model | Input | Cached Input | Cached Input Storage | Output |
| :- | :- | :- | :- | :- |
| GLM-5.3-FlashX-Legal | \$0.37 | \$0.075 | Limited-time Free | \$1.25 |
| GLM-5.2 | \$1.4 | \$0.26 | Limited-time Free | \$4.4 |
| GLM-5.1 | \$1.4 | \$0.26 | Limited-time Free | \$4.4 |
| GLM-5 | \$1 | \$0.2 | Limited-time Free | \$3.2 |
| GLM-4.7 | \$0.6 | \$0.11 | Limited-time Free | \$2.2 |
| GLM-4.6 | \$0.6 | \$0.11 | Limited-time Free | \$2.2 |

### Built-in Tools

| Tool | Cost |
| :- | :- |
| Web Search | \$0.01 / use |

### Image Generation Models

Prices per image.

| Model | Price |
| :- | :- |
| GLM-Image | \$0.015 |

### Agents

| Agent | Price |
| :- | :- |
| GLM Slide/Poster Agent(beta) | \$0.7 / MTok |
| General-Purpose Translation | \$3 / MTok |
| Popular Special Effects Video Templates | \$0.2 / video |

> **Other Pricing**
>
> ### Text Models
>
> Prices per 1M tokens.
>
> | Model | Input | Cached Input | Cached Input Storage | Output |
> | :- | :- | :- | :- | :- |
> | GLM-4.7-FlashX | \$0.07 | \$0.01 | Limited-time Free | \$0.4 |
> | GLM-4.5 | \$0.6 | \$0.11 | Limited-time Free | \$2.2 |
> | GLM-4.5-X | \$2.2 | \$0.45 | Limited-time Free | \$8.9 |
> | GLM-4.5-Air | \$0.2 | \$0.03 | Limited-time Free | \$1.1 |
> | GLM-4.5-AirX | \$1.1 | \$0.22 | Limited-time Free | \$4.5 |
> | GLM-4-32B-0414-128K | \$0.1 | - | - | \$0.1 |
> | GLM-4.7-Flash | Free | Free | Free | Free |
> | GLM-4.5-Flash | Free | Free | Free | Free |
>
> ### Vision Models
>
> Prices per 1M tokens.
>
> | Model | Input | Cached Input | Cached Input Storage | Output |
> | :- | :- | :- | :- | :- |
> | GLM-4.6V | \$0.3 | \$0.05 | Limited-time Free | \$0.9 |
> | GLM-OCR | \$0.03 | \\ | \\ | \$0.03 |
> | GLM-4.6V-FlashX | \$0.04 | \$0.004 | Limited-time Free | \$0.4 |
> | GLM-4.5V | \$0.6 | \$0.11 | Limited-time Free | \$1.8 |
> | GLM-4.6V-Flash | Free | Free | Free | Free |
>
> ### Image Generation Models
>
> Prices per image.
>
> | Model | Price |
> | :- | :- |
> | CogView-4 | \$0.01 |
>
> ### Video Generation Models
>
> Prices per video.
>
> | Model | Price |
> | :- | :- |
> | CogVideoX-3 | \$0.2 |
>
> ### Audio Models
>
> | Model | Price |
> | :- | :- |
> | GLM-ASR-2512 | \$0.03 / MTok (equivalent to approximately \$0.0024/minute) |


# Quick Start
Source: https://docs.z.ai/guides/overview/quick-start



> [!NOTE]
>
> Tired of limits? GLM Coding Plan — monthly access to world-class models for coding and office tasks, compatible with 20+ top agent tools. All from just \$18/month.  [Try it now →](https://z.ai/subscribe?utm_campaign=Platform_Ops&_channel_track_key=DaprgHIc)

## Getting Started

**1. Get API Key**

* Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
* Access [Billing Page](https://z.ai/manage-apikey/billing) to top up if needed.
* Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
* Copy your API Key for use.

  > **[Z.AI Open Platform](https://z.ai/model-api)**
  >
  > Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.

  > **[API Keys Management](https://z.ai/manage-apikey/apikey-list)**
  >
  > Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.

**2. Choose Model**

> The platform offers multiple models, and you can select the appropriate model based on your needs. For detailed model introductions, please refer to the [Models & Agents](/guides/overview/overview).

> **[GLM-5.3](/guides/llm/glm-5.3)**
>
> Z.ai’s latest flagship model, delivering comprehensive advancements in software engineering and agent capabilities.

> **[GLM-5.3-FLASH](/guides/vlm/glm-5.3-flash)**
>
> Multimodal Coding model, specializing in visual programming.

**3. Choose the Calling Method**

Our platform provides various development approaches; you can select the best fit for your project needs and tech stack.

> **[HTTP API](/guides/develop/http/introduction)**
>
> Standard RESTful API, compatible with all programming languages.

> **[Z.AI Python SDK](/guides/develop/python/introduction)**
>
> Official Python SDK, featuring full type hints and async support.

> **[Z.AI Java SDK](/guides/develop/java/introduction)**
>
> Official Java SDK, designed for high concurrency and availability.

> **[OpenAI Python SDK](/guides/develop/openai/python)**
>
> OpenAI SDK Compatibility, quickly migrating from OpenAI.

> **[API Reference](/api-reference)**
>
> Complete API documentation with parameter descriptions.

**4. Make API Call**

After preparing your `API Key` and selecting a model, you can start making API calls. Here are examples using `curl`, `Python SDK`, and `Java SDK`:

> [!WARNING]
>
> When using GLM Coding Plan, please follow the [tutorial](/devpack/quick-start) to configure your dedicated endpoint.

**cURL**

```bash theme={null}
curl -X POST "https://api.z.ai/api/paas/v4/chat/completions" \
-H "Content-Type: application/json" \
-H "Accept-Language: en-US,en" \
-H "Authorization: Bearer YOUR_API_KEY" \
-d '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful AI assistant."
        },
        {
            "role": "user",
            "content": "Hello, please introduce yourself."
        }
    ]
}'
```

**Official Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Usage Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key="YOUR_API_KEY")

# Create chat completion request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {
            "role": "system",
            "content": "You are a helpful AI assistant."
        },
        {
            "role": "user",
            "content": "Hello, please introduce yourself."
        }
    ]
)

# Get response
print(response.choices[0].message.content)
```

**Official Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

**Usage Example**

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import java.util.Arrays;

public class QuickStart {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();

        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Hello, who are you?")
                    .build()
            ))
            .stream(false)
            .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        // Get response
        System.out.println(response.getData().getChoices().get(0).getMessage().getContent());
    }
}
```

**OpenAI Python SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
pip install --upgrade 'openai>=1.0'
```

**Verify Installation**

```python theme={null}
python -c "import openai; print(openai.__version__)"
```

**Usage Example**

```python theme={null}
from openai import OpenAI

client = OpenAI(
    api_key="your-Z.AI-api-key",
    base_url="https://api.z.ai/api/paas/v4/"
)

completion = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "system", "content": "You are a smart and creative novelist"},
        {"role": "user", "content": "Please write a short fairy tale story as a fairy tale master"}
    ]
)

print(completion.choices[0].message.content)
```

**OpenAI NodeJs SDK**

**Install SDK**

```bash theme={null}
# Install or upgrade to latest version
npm install openai

# Or using yarn
yarn add openai
```

**Usage Example**

```javascript theme={null}
import OpenAI from "openai";

const client = new OpenAI({
    apiKey: "your-Z.AI-api-key",
    baseURL: "https://api.z.ai/api/paas/v4/"
});

async function main() {
    const completion = await client.chat.completions.create({
        model: "glm-5.3",
        messages: [
            { role: "system", content: "You are a helpful AI assistant." },
            { role: "user", content: "Hello, please introduce yourself." }
        ]
    });

    console.log(completion.choices[0].message.content);
}

main();
```

**OpenAI Java SDK**

**Install SDK**

**Maven**

```xml theme={null}
<dependency>
    <groupId>com.openai</groupId>
    <artifactId>openai-java</artifactId>
    <version>2.20.1</version>
</dependency>
```

**Gradle (Groovy)**

```groovy theme={null}
implementation 'com.openai:openai-java:2.20.1'
```

**Usage Example**

```java theme={null}
import com.openai.client.OpenAIClient;
import com.openai.client.okhttp.OpenAIOkHttpClient;
import com.openai.models.chat.completions.ChatCompletion;
import com.openai.models.chat.completions.ChatCompletionCreateParams;

public class QuickStart {
    public static void main(String[] args) {
        // Initialize client
        OpenAIClient client = OpenAIOkHttpClient.builder()
            .apiKey("your-Z.AI-api-key")
            .baseUrl("https://api.z.ai/api/paas/v4/")
            .build();

        // Create chat completion request
        ChatCompletionCreateParams params = ChatCompletionCreateParams.builder()
            .addSystemMessage("You are a helpful AI assistant.")
            .addUserMessage("Hello, please introduce yourself.")
            .model("glm-5.3")
            .build();

        // Send request and get response
        ChatCompletion chatCompletion = client.chat().completions().create(params);
        Object response = chatCompletion.choices().get(0).message().content();

        System.out.println(response);
    }
}
```

### Get More

> **[API Reference](/api-reference)**
>
> Access API Reference.

> **[Python SDK](https://github.com/zai-org/z-ai-sdk-python)**
>
> Access Python SDK Github

> **[Java SDK](https://github.com/zai-org/z-ai-sdk-java)**
>
> Access Java SDK Github


# Stream Tool Call
Source: https://docs.z.ai/guides/tools/stream-tool



Stream Tool Call is a unique feature of Z.ai's latest GLM-4.6 model, allowing real-time access to reasoning processes, response content, and tool call information during tool invocation, providing better user experience and real-time feedback.

## Features

Tool calling in the latest GLM model now supports streaming output for responses. This allows developers to stream tool usage parameters without buffering or JSON validation when calling `chat.completions`, thereby reducing call latency and providing a better user experience.

### Core Parameter Description

* **`stream=True`**: Enable streaming output, must be set to `True`
* **`tool_stream=True`**: Enable tool call streaming output
* **`model`**: Use a model that supports tool calling, limited to `glm-4.6` `glm-4.7` `glm-5`

### Response Parameter Description

The `delta` object in streaming responses contains the following fields:

* **`reasoning_content`**: Text content of the model's reasoning process
* **`content`**: Text content of the model's response
* **`tool_calls`**: Tool call information, including function names and parameters

## Code Example

By setting the `tool_stream=True` parameter, you can enable streaming tool call functionality:

**Python SDK**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API key')

# Create streaming tool call request
response = client.chat.completions.create(
    model="glm-4.6",  # Use model that supports tool calling
    messages=[
        {"role": "user", "content": "How's the weather in Beijing?"},
    ],
    tools=[
        {
            "type": "function",
            "function": {
                "name": "get_weather",
                "description": "Get current weather conditions for a specified location",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "location": {"type": "string", "description": "City, e.g.: Beijing, Shanghai"},
                        "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
                    },
                    "required": ["location"]
                }
            }
        }
    ],
    stream=True,        # Enable streaming output
    tool_stream=True    # Enable tool call streaming output
)

# Initialize variables to collect streaming data
reasoning_content = ""      # Reasoning process content
content = ""               # Response content
final_tool_calls = {}      # Tool call information
reasoning_started = False  # Reasoning process start flag
content_started = False    # Content output start flag

# Process streaming response
for chunk in response:
    if not chunk.choices:
        continue

    delta = chunk.choices[0].delta

    # Handle streaming reasoning process output
    if hasattr(delta, 'reasoning_content') and delta.reasoning_content:
        if not reasoning_started and delta.reasoning_content.strip():
            print("\n🧠 Thinking Process:")
            reasoning_started = True
        reasoning_content += delta.reasoning_content
        print(delta.reasoning_content, end="", flush=True)

    # Handle streaming response content output
    if hasattr(delta, 'content') and delta.content:
        if not content_started and delta.content.strip():
            print("\n\n💬 Response Content:")
            content_started = True
        content += delta.content
        print(delta.content, end="", flush=True)

    # Handle streaming tool call information
    if delta.tool_calls:
        for tool_call in delta.tool_calls:
            index = tool_call.index
            if index not in final_tool_calls:
                # New tool call
                final_tool_calls[index] = tool_call
                final_tool_calls[index].function.arguments = tool_call.function.arguments
            else:
                # Append tool call parameters (streaming construction)
                final_tool_calls[index].function.arguments += tool_call.function.arguments

# Output final tool call information
if final_tool_calls:
    print("\n📋 Function Calls Triggered:")
    for index, tool_call in final_tool_calls.items():
        print(f"  {index}: Function Name: {tool_call.function.name}, Parameters: {tool_call.function.arguments}")
```

## Use Cases

> **Intelligent Customer Service System**
>
> * Real-time display of query progress
> * Improve waiting experience

> **Code Assistant**
>
> * Real-time code analysis process
> * Display tool call chain


# Web Search
Source: https://docs.z.ai/guides/tools/web-search



Z.AI provides developers with a full suite of AI search tools, covering three core services: **basic search (Web Search API)**, **retrieval-augmented generation (Web Search in Chat)**, and **intelligent search agents (Search Agent)**. Through a unified API, offering end-to-end capabilities—from raw web data retrieval and fusion of search results with LLM output, to multi-turn dialogue context management. This empowers developers to build trustworthy, real-time, and traceable AI applications at **lower cost**.

* View [Product Price](/guides/overview/pricing)
* View Your [API Key](https://chat.z.ai)

## Product Overview

| Service Module | Developer Value | Technical Features |
| :- | :- | :- |
| Web Search API | Directly obtain **structured search results** (title/summary/link, etc.) | Multi-Search Engine Support |
| Web Search in Chat | Incorporate search results into large model-generated **answers with cited**  **web sources**. | Seamless Integration of Real-Time Retrieval and LLM Generation |

## Web Search API

Web Search API is a search engine specifically designed for large language models. Building on the traditional capabilities of webpage crawling and ranking, it enhances intent recognition and returns results that are more suitable for LLM processing (webpage titles, URLs, summaries, site names, and favicons).

1. Intent-Enhanced Retrieval: Intelligently identifies the user’s query intent and automatically determines whether web search is needed.
2. Structured Output: Returns data formats optimized for LLMs, including titles, URLs, summaries, site names, and icons.
3. Customizable Search Scope: Allows developers to specify the number of results, domain constraints, and time ranges. It also supports adjustable summary lengths for fine-grained search behavior control.
4. Time-Aware Output Control: The response can include the webpage’s publication time, facilitating timeliness analysis and result ranking.

**API Call**

* API Docs: [Web Search API](/api-reference/tools/web-search)
* Example: Search Financial News

```python Search Financial News theme={null}
from zai import ZaiClient

client = ZaiClient(api_key="")  # Fill in your own APIKey

response = client.web_search.web_search(
   search_engine="search-prime",
   search_query="search economic events",
   count=15, # The number of results to return, ranging from 1-50, default 10
   search_domain_filter="www.sohu.com", # Only access content from specified domain names.
   search_recency_filter="noLimit", # Search for content within specified date ranges
)
print(response)
```

```json theme={null}
WebSearchResp(
{
    "created": 1748261757,
    "id": "20250526201557dda85ca6801b467b",
    "request_id": "20250526201557dda85ca6801b467b",
    "search_result": [
        {
            "content": "1. China's outward direct investment from January to April reached $57.54 billion, a year-on-year increase of 7.5%. The trade-in policy continues to show results, with retail sales of home appliances maintaining double-digit growth for eight consecutive months.\n2. China has implemented multiple measures to accelerate the construction of a sci-tech financial system. The CSRC will support tech companies breaking through key core technologies by granting them 'green channel' access. The PBOC stated that nearly 100 institutions have issued over 250 billion yuan in sci-tech innovation bonds.\n3. In May, China's passenger car retail sales are estimated at around 1.85 million units, up 8.5% year-on-year and 5.4% month-on-month. New energy vehicle sales are expected to reach 980,000, maintaining a penetration rate of approximately 52.9%.\n4. The Beijing Humanoid Robot Innovation Center released the world's first 'Humanoid Robot Intelligence Grading' standard, establishing a 'four-dimensional, five-level' evaluation framework covering 'perception and cognition, decision-making and learning, execution performance, and collaborative interaction.'\n5. Federal Reserve Governor Waller stated that if the Trump administration's tariffs stabilize around 10%, the Fed may cut interest rates in the second half of 2025.\n6. The U.S. manufacturing PMI rose to a three-month high of 52.3 in May. The services PMI preliminary reading was 52.3, hitting a two-month high. New orders grew at their fastest pace in over a year, while price indicators reached a nearly three-year peak.\n7. The eurozone manufacturing PMI preliminary reading improved slightly to 49.2 in May, but the services PMI unexpectedly dropped sharply to 48.9, marking the worst performance in 16 months. This dragged the eurozone composite PMI down to 49.5 from 50.4 in April. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/sohu_icon_new.jpg",
            "link": "https://www.sohu.com/a/897879632_121123890",
            "media": "Sohu",
            "publish_date": "2025-05-23",
            "refer": "ref_1",
            "title": "Financial Morning Briefing for May 23, 2025"
        },
        {
            "content": "1. The CSRC and seven other departments issued 'Several Measures to Support Financing for Small and Micro Enterprises,' proposing 23 specific initiatives across eight areas: increasing financing supply, reducing comprehensive financing costs, improving financing efficiency, enhancing support precision, implementing regulatory policies, strengthening risk management, improving policy safeguards, and ensuring effective implementation.\nComment: Increased support for small and micro enterprise financing, with 23 measures targeting the real economy to alleviate financing difficulties.\n2. Ministry of Commerce: From January to April 2025, China's outward direct investment across all sectors totaled $57.54 billion, up 7.5% year-on-year.\nComment: Steady growth in China's outward investment, with 'Belt and Road' cooperation likely a key driver.\nIndustry Insights\n3. Ministry of Commerce: From September 2024 to April 2025, retail sales of home appliances maintained double-digit growth for eight consecutive months. In April, retail sales of household appliances and audio-visual equipment in above-quota units surged 38.8% year-on-year, ranking first among 16 major consumer goods categories.\nComment: Strong demand for home appliances, driven by policy incentives and consumption upgrades.\n4. National Energy Administration: As of end-April, China's total installed power generation capacity reached 3.49 billion kW, up 15.9% year-on-year. From January to April, major power enterprises completed 193.3 billion yuan in power source engineering investments, up 1.6% year-on-year, while grid engineering investments totaled 140.8 billion yuan, up 14.6%.\nComment: Accelerated energy infrastructure development, with grid investment growth outpacing power source investment.\nGlobal Perspective\n5. The U.S. House Rules Committee approved Trump's comprehensive tax cut bill, paving the way for a full House vote.\nComment: Progress on Trump's tax cuts, which could stimulate U.S. consumption but exacerbate fiscal deficit pressures if enacted.\nFinancial Markets\n6. Most domestic commodity futures closed lower, with weak energy performance, widespread pressure on chemical futures, and declines in the black commodities sector. Agricultural products, particularly apple futures, also saw notable drops.\nComment: Broad-based commodity market declines reflect weak demand expectations.\n7. On Thursday, the onshore yuan closed at 7.2040 against the dollar at 16:30, up 25 basis points from the previous session. The central parity rate was adjusted up 34 basis points to 7.1903. Treasury futures mostly flatlined, with the 30-year contract down 0.04% and the 10-year contract up 0.01%.\nComment: Stabilizing yuan exchange rate, with parity rate adjustments signaling policy support.\n8. On May 22, the three major indices closed lower. Defense stocks led gains, while innovative drug concepts were active. Bank stocks bucked the trend, while pet economy concepts lagged, beauty care stocks plunged, and solid-state battery concepts retreated.\nComment: Structural divergence in A-shares as risk aversion rises, with caution advised over June U.S. debt maturity risks and potential global market corrections. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/sohu_icon_new.jpg",
            "link": "https://www.sohu.com/a/897874861_121123901",
            "media": "Sohu",
            "publish_date": "2025-05-23",
            "refer": "ref_2",
            "title": "Financial Morning Briefing for May 23"
        },
        {
            "content": "According to Zhitong Finance APP, Macao's Composite Consumer Price Index (CPI) rose 0.23% year-on-year and 0.16% month-on-month in April 2025. The data, released by Macao's Statistics and Census Service on May 23, reflects subtle shifts in Macao's consumer market.\n1. Changes in Macao's CPI\nData shows that the 12-month average composite CPI as of April 2025 rose 0.42% compared to the same period a year earlier (May 2023–April 2024). Among major categories, the price index for Recreation, Sports, and Culture rose 2.47% year-on-year, while Miscellaneous Goods and Services increased 1.82%, signaling potential economic recovery.\n2. Impact of Food and Rent\nFood and non-alcoholic beverage prices rose 0.48% year-on-year, driven by higher dining-out and takeaway costs. Housing and fuel prices also increased 0.21% due to rising rents. As tourism recovers, dining demand has pushed up service fees.\nHowever, prices for Communication (-3.01%), Clothing (-2.35%), and Transport (-1.67%) declined, complicating the overall CPI trend.\n3. Month-on-Month Changes\nCompared to March 2025, April's CPI rose 0.16%. Recreation, Sports, and Culture saw the largest increase (2.6%) due to higher hotel rates. Summer clothing launches and rebounding airfares also lifted Clothing (+0.78%) and Transport (+0.42%) prices.\nDespite this, Food and non-alcoholic beverages dipped slightly (-0.04%) as lower prices for fresh fish, seafood, vegetables, and fruit offset higher dining costs.\n4. CPI Classification\nMacao's CPI is compiled into three series to reflect price impacts on households with different spending patterns. Category A (50% of households) covers monthly expenditures of 11,000–35,999 MOP, while Category B (30%) covers 36,000–71,999 MOP.\nHousing and fuel (34.47%), food (29.47%), and transport (8.33%) carry the highest weights, ensuring precise reflection of consumption trends across income groups.\n5. Outlook\nMacao's consumer market will continue to face multiple influences. While some categories show price declines, the overall mild CPI rise suggests recovery potential. With tourism rebounding and demand growing, CPI may maintain a gradual upward trend.\nIn summary, April 2025's CPI changes reflect market dynamics and new cost-of-living challenges. Authorities must monitor price impacts to craft effective policies ensuring stable living standards. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/sohu_icon_new.jpg",
            "link": "https://www.sohu.com/a/898091955_121956424",
            "media": "Sohu",
            "publish_date": "2025-05-23",
            "refer": "ref_3",
            "title": "Macao's April 2025 CPI Rises Slightly: Key Drivers Explained"
        },
        {
            "content": "Zhitong Finance APP learned that on May 23, the Asset Management Association of China (AMAC) released its monthly report on private fund manager registrations and product filings. In April 2025, 1,606 new private funds were filed nationwide, with a total size of 643.72 billion yuan. This included 1,189 private securities investment funds (366.95 billion yuan), 143 private equity funds (161.01 billion yuan), and 274 venture capital funds (115.76 billion yuan).\n1. Private Fund Manager Registration Overview\n(1) Monthly Registrations\nIn April 2025, 16 institutions passed AMAC's AMBERS system, including 4 private securities investment fund managers and 12 private equity/venture capital fund managers. During the same period, 76 private fund managers were deregistered.\n(2) Existing Managers\nAs of end-April 2025, there were 19,891 private fund managers overseeing 141,579 funds with 20.22 trillion yuan in assets. This included 7,827 private securities investment fund managers, 11,867 private equity/venture capital managers, 6 asset allocation managers, and 191 others.\n(3) Geographic Distribution\nBy registered location (across 36 regions), 72.20% of managers were concentrated in Shanghai (3,714), Beijing (3,244), Shenzhen (3,012), Guangdong (ex-Shenzhen, 1,594), Zhejiang (ex-Ningbo, 1,590), and Jiangsu (1,207), accounting for 18.67%, 16.31%, 15.14%, 8.01%, 7.99%, and 6.07% respectively.\nBy AUM, the top regions were Shanghai (5.08 trillion yuan, 25.12%), Beijing (4.68 trillion yuan, 23.13%), Shenzhen (1.97 trillion yuan, 9.74%), Guangdong (ex-Shenzhen, 1.30 trillion yuan, 6.43%), Jiangsu (1.18 trillion yuan, 5.82%), and Zhejiang (ex-Ningbo, 967.16 billion yuan, 4.78%), totaling 75.02%.\n2. Private Fund Filings Overview\n(1) Monthly Filings\nApril 2025 saw 1,606 new private funds (643.72 billion yuan), including 1,189 private securities funds (366.95 billion yuan), 143 private equity funds (161.01 billion yuan), and 274 venture capital funds (115.76 billion yuan).\n(2) Existing Funds\nAs of end-April 2025, there were 141,579 private funds (20.22 trillion yuan), including 84,673 private securities funds (5.51 trillion yuan), 30,205 private equity funds (10.96 trillion yuan), and 25,830 venture capital funds (3.41 trillion yuan).",
            "icon": "https://sfile.chatglm.cn/searchImage/sohu_icon_new.jpg",
            "link": "https://www.sohu.com/a/898137965_323087",
            "media": "Sohu",
            "publish_date": "2025-05-23",
            "refer": "ref_4",
            "title": "AMAC: 1,606 New Private Funds Filed in April, Totaling 643.72 Billion Yuan"
        }
    ]
}
)
```

### MCP Server

Access the [official MCP documentation](https://modelcontextprotocol.io/introduction) to learn more about the protocol.

**Installation Guide**

* Use clients that support the MCP protocol, such as Cursor and Cherry Studio.
* Obtain a [Standard API Key](https://z.ai/manage-apikey/apikey-list) from the Z.AI Platform.

**Using in Cursor**

Cursor 0.45.6 includes MCP functionality. Cursor acts as an MCP service client and can connect to the MCP service with simple configuration.

Navigation Path: Cursor Settings → \[Features] → \[MCP Servers]

![370b139e 8201 40a0 9c44 3faf4b1b9655 Pn](https://mintcdn.com/zhipu-32152247/fQm1SxNtD2jBDQ3i/images/370b139e-8201-40a0-9c44-3faf4b1b9655.png?fit=max\&auto=format\&n=fQm1SxNtD2jBDQ3i\&q=85\&s=ea0a9a85b3c412cfcaf112ce29b66be7)

**Configure MCP Server**

```python theme={null}
{
  "mcpServers": {
    "z.ai-web-search-sse": {
      "url": "https://api.z.ai/api/mcp/web_search/sse?Authorization=YOUR API Key"
    }
  } 
}
```

**Cursor MCP Usage**

Cursor MCP must be used in Composer's Agent mode.

![Bb97f60a B887 4e13 92d7 8156227676c7 Pn](https://mintcdn.com/zhipu-32152247/fQm1SxNtD2jBDQ3i/images/bb97f60a-b887-4e13-92d7-8156227676c7.png?fit=max\&auto=format\&n=fQm1SxNtD2jBDQ3i\&q=85\&s=03cca96249a7376deb2e3e316f969771)

## Web Search in Chat

Web Search in Chat allows the Completions API to call search engines, combining real-time web retrieval results with GLM's generative capabilities to provide up-to-date and verifiable answers.

* API Docs: [Web Search in Chat](/api-reference/tools/web-search)
* Example: Financial Analysis Summary

```python Financial Analysis Summary theme={null}
from zai import ZaiClient

# Initialize the ZaiClient client
client = ZaiClient(api_key="YourAPIKey")

# Define tool parameters
tools = [{
    "type": "web_search",
    "web_search": {
         "enable": "True",
         "search_engine": "search-prime",
         "search_result": "True",
         "search_prompt": "You are a financial analyst. Please use concise language to summarize the key information in {{search_result}} from the web search, ranked by importance and citing the source date. Today's date is April 11, 2025.",
         "count": "5",
         "search_domain_filter": "www.sohu.com",
         "search_recency_filter": "noLimit",
         "content_size": "high"
    }
}]

# Define user message
messages = [{
    "role": "user",
    "content": "Key financial events, policy changes, and market data in April 2025"
}]

# Call the API to get a response
response = client.chat.completions.create(
    model="glm-4-air",  # Model identifier
    messages=messages,  # User messages
    tools=tools         # Tool parameters
)

# Print the response result
print(response)
```

```json theme={null}
{
    "choices": [
        {
            "finish_reason": "stop",
            "index": 0,
            "message": {
                "content": "Based on the documents you provided, here are the important financial events for April 2025, ranked by importance:\n\n1. **G20 Finance and Central Bank Ministers Meeting** - Date to be determined. The G20 meeting will discuss key issues such as global economic recovery, financial stability, and sustainable development. This will have a profound impact on the coordination of global economic policies and financial market sentiment. [Source: ref_1]\n\n2. **Preliminary Manufacturing PMI Releases for Multiple Countries and Regions** - Including France, Germany, the Eurozone, and the UK. These data will reveal the activity status of their respective manufacturing sectors, providing critical insights for investors. [Source: ref_1]\n\n3. **Eurozone Seasonally Adjusted Trade Balance for February** - Released at 17:00. The performance of the trade balance will significantly impact the euro and related markets, reflecting the region's trade conditions. [Source: ref_1]\n\n4. **Speeches by Fed Officials Moussalem and Waller** - At 21:30. Speeches by Fed officials will play a crucial role in driving market volatility and sparking discussions about the future direction of U.S. monetary policy. [Source: ref_1]\n\n5. **U.S. April S&P Global Manufacturing and Services PMI Preliminaries** - Released at 21:45. The PMI data for manufacturing and services will indicate the overall trend of economic activity, influencing investment decisions. [Source: ref_1]\n\n6. **U.S. March New Home Sales Annualized Total** - Released at 22:00. As an important indicator of economic health, new home sales will signal changes in market demand. [Source: ref_1]\n\n7. **U.S. EIA Crude Oil Inventory for the Week Ending April 18** - Released at 22:30. Crude oil inventory data will provide a key reference point for the market, affecting oil price fluctuations. [Source: ref_1]\n\n8. **Fed Releases the Beige Book on Economic Conditions** - At 02:00 the next day. The Fed's Beige Book report will summarize economic development trends, impacting investor confidence and strategies. [Source: ref_1]\n\n9. **Speech by Fed Governor Lisa Cook** - At 07:10. Governor Cook's speech may reflect the Fed's current thinking on the economic situation, influencing market sentiment and asset price fluctuations. [Source: ref_2]\n\n10. **Release of China's Core Economic Data** - Including the monthly report on residential sales prices in 70 large and medium-sized cities, Q1 GDP annual rate, total retail sales of consumer goods, and industrial added value above a designated size. These will have a significant impact on global economic growth and the healthy development of China's economy. [Source: ref_2]\n\n11. **Release of European Economic Data** - Including the UK's March CPI monthly rate, the Eurozone's seasonally adjusted current account for February, and the final annual CPI rate for March. These will provide investors with a comprehensive view of the European economy. [Source: ref_2]\n\n12. **WTO Global Trade Outlook Report** - Released at 20:00. The report may reveal the future direction of global trade and changes among countries. [Source: ref_2]\n\n13. **Release of Key U.S. Economic Indicators** - Including the March retail sales monthly rate, industrial production monthly rate, and NAHB housing market index. These will influence investors' assessment of the strength of the U.S. economic recovery. [Source: ref_2]\n\n14. **Bank of Canada Interest Rate Decision** - Released at 21:45. The direction of Canada's monetary policy is also worth noting. [Source: ref_2]\n\n15. **China's State Council Information Office Releases Q1 Import and Export Data** - Released at 10:00. The growth or decline in import and export data will directly affect GDP growth expectations and investor confidence. [Source: ref_3]\n\n16. **Canada's February Wholesale Sales Monthly Rate Data** - Released at 20:30. Wholesale sales are an important indicator of retail demand expectations, directly reflecting the health of the economy. [Source: ref_3]\n\n17. **U.S. One-Year Inflation Expectations Data for March** - Released at 23:00. This data is a key reference for assessing U.S. consumers' and businesses' views on future economic conditions. [Source: ref_3]\n\n18. **OPEC Monthly Crude Oil Market Report** - Release time to be determined. The report typically influences global crude oil market price trends. [Source: ref_3]\n\n19. **Speeches by Fed Officials** - At 00:00 and 06:00 the next day. Speeches by Fed officials will outline the future direction of Fed policy, affecting global investor expectations. [Source: ref_3]\n\nThe above events and data will have a significant impact on global market trends. Investors should closely monitor them and adjust their strategies accordingly.",
                "role": "assistant"
            }
        }
    ],
    "created": 1748311718,
    "id": "20250527100811da2f8f7243f94b02",
    "model": "glm-4-air",
    "request_id": "20250527100811da2f8f7243f94b02",
    "usage": {
        "completion_tokens": 868,
        "prompt_tokens": 4199,
        "total_tokens": 5067
    },
    "web_search": [
        {
            "content": "In the current volatile global economic climate, April 23, 2025, is set to be a major focus day for the financial world. Several important financial events and economic data will be released, keeping market participants on their toes as they anticipate the upcoming dynamics. Below is an in-depth analysis of the key events and data scheduled for the day to help you gain insights into future market trends.\n1. G20 Finance and Central Bank Ministers Meeting (Time to be determined)\nToday, the G20 has become a crucial platform for global economic governance. Finance and central bank ministers will discuss issues such as economic recovery, financial stability, and sustainable development. This meeting will not only impact the coordination of global economic policies but also have a profound effect on financial market sentiment.\n2. France April Preliminary Manufacturing PMI (15:15)\nAs the economic engine of the Eurozone, the current state of France's manufacturing sector will directly influence the region's recovery. The preliminary data will indicate the fragility or stability of France's economic recovery, providing key insights for investors.\n3. Germany April Preliminary Manufacturing PMI (15:30)\nGermany's manufacturing performance is always closely watched, serving as a barometer of the country's economic health and a significant indicator of the Eurozone's overall economic performance. The preliminary data will offer valuable references for investors to determine future market strategies.\n4. Eurozone April Preliminary Manufacturing PMI (16:00)\nFollowing this, the Eurozone's preliminary manufacturing PMI will also be released. This data is expected to reveal the activity status of the region's manufacturing sector, with markets keenly observing whether the upward trend can be sustained. This will also affect investor confidence in the euro.\n5. UK April Preliminary Manufacturing PMI (16:30)\nAfter facing multiple challenges, the recovery of the UK's manufacturing sector is critical. The preliminary PMI will further clarify whether the UK economy is on a positive trajectory, providing economic context for the upcoming general election.\n6. UK April Preliminary Services PMI (16:30)\nAs a key pillar of the UK's national economy, the services PMI directly impacts overall GDP trends. Changes in this data will offer practical suggestions for the flexibility of economic policies, especially during a period when the UK government faces numerous reform challenges.\n7. Eurozone Seasonally Adjusted Trade Balance for February (17:00)\nThe performance of the trade balance will significantly impact the euro and related markets. The positive or negative figures will directly reflect the region's trade conditions, influencing central bank monetary policies and market liquidity.\n8. Speeches by Fed Officials Moussalem and Waller (21:30)\nSpeeches by Fed officials often play a crucial role in driving market volatility. The remarks by Moussalem and Waller will spark discussions about the future direction of U.S. monetary policy, affecting interest rate decisions and economic guidance.\n9. U.S. April S&P Global Manufacturing PMI Preliminary (21:45)\nKey data on the recovery of U.S. industries is coming. Will the market see strong economic signals that boost global market sentiment? The preliminary manufacturing PMI will be the focal point of both skepticism and anticipation.\n10. U.S. April S&P Global Services PMI Preliminary (21:45)\nThe services PMI reflects the overall trend of economic activity, indicating consumer confidence and spending power. Changes in this value will also provide more guidance for investment decisions.\n11. U.S. March New Home Sales Annualized Total (22:00)\nAs an important indicator of economic health, the changes in new home sales this month may signal shifts in market demand, especially under current fluctuating interest rates.\n12. U.S. EIA Crude Oil Inventory for the Week Ending April 18 (22:30)\nAmid the recovery of global energy demand and geopolitical uncertainties, crude oil inventory data will provide a key reference point for the market, influencing oil price fluctuations.\n13. Fed Releases the Beige Book on Economic Conditions (02:00 the next day)\nThe Fed's Beige Book report will summarize economic development trends, often sparking widespread market discussions and trading volatility upon release. Whether this report brings new economic insights will directly affect investor confidence and strategies.\nSummary\nIn summary, the financial events and economic data on April 23, 2025, will profoundly impact global market trends, potentially offering rare opportunities and warnings for investors. Whether you are a policy analyst or a market participant, it is essential to closely monitor these key events and data, adjusting strategies promptly to seize market opportunities. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/www_sohu_com_icon_new.jpg",
            "link": "https://www.sohu.com/a/887836650_121956424",
            "media": "Sohu",
            "publish_date": "",
            "refer": "ref_1",
            "title": "April 23 Financial Highlights: Analyzing Major Global Market Events and New Data Trends (Published: 2025-04-23 07:07:00)"
        },
        {
            "content": "In the current ever-changing global economy, the impact of financial events on the market cannot be ignored. April 16, 2025, will witness a series of highly anticipated financial activities and important data releases, forming a focal point for investors and stirring the nerves of international markets.\n1. Market Volatility and Investor Expectations\nRecently, global markets have faced numerous uncertainties, particularly regarding inflation, interest rate policies, and the winding path of economic recovery, all of which have immersed investors in a sea of analysis and predictions.\nOn April 16, the U.S., Europe, China, and other regions will release a series of heavyweight economic data, which investors eagerly anticipate. These data may provide crucial clues for future market trends. Just as every fluctuation in the stock market is a process of information transmission, only by deeply understanding this information can one make more effective investment decisions.\n2. Prelude to Key Data\n1. Important Fed Speeches\nAt 07:10, Fed Governor Lisa Cook's speech will undoubtedly attract market attention. As a member of the Fed, her views and remarks may reflect the Fed's current thinking on the economic situation. If she mentions the outlook on interest rates or views on economic recovery, this will directly impact market sentiment and the price fluctuations of various asset classes, becoming a key factor influencing the opening of the New York market.\n2. Release of China's Core Economic Data\nFollowing this, at 09:30, China will release the monthly report on residential sales prices in 70 large and medium-sized cities, outlining the current state of China's real estate market. At 10:00, a series of heavyweight economic data will be released, including Q1 GDP annual rate, total retail sales of consumer goods, and industrial added value above a designated size. The release of these data will have a significant impact on global economic growth and the healthy development of China's economy, especially against the backdrop of complex and changing international economics. Analyzing these data will provide a profound understanding of the direction of China's market.\n3. Importance of European Data\nAt 14:00, the release of the UK's March CPI monthly rate will reveal the current inflationary pressures in the UK economy. At 16:00, the Eurozone's seasonally adjusted current account for February and the final annual CPI rate for March will be released, offering investors a more comprehensive view of the European economy. In recent years, the policies of the European Central Bank have also garnered widespread market attention, and changes in inflation trends may influence future monetary policy directions.\n3. Global Trade Outlook from a Global Perspective\nAt 20:00, the World Trade Organization (WTO) will release its Global Trade Outlook Report. Faced with geopolitical risks and challenges to global supply chains, the WTO report may reveal the future direction of global trade and changes among countries. If the trends mentioned in the report diverge from investor expectations, it is bound to trigger strong market fluctuations.\n4. Release of Key U.S. Economic Indicators\nNews from the U.S. market is equally unmissable. At 20:30, the U.S. March retail sales monthly rate will be released, followed by the industrial production monthly rate at 21:15 and the NAHB housing market index before 22:00. These will directly influence investors' assessment of the strength of the U.S. economic recovery. These data, reflecting consumption and production activities, are vital indicators of economic health and can guide investors toward future trends.\nAt 22:30, the U.S. EIA crude oil inventory data for the week ending April 11 will reveal the dynamics of the global energy market. As the global economy remains volatile, fluctuations in energy prices will have a profound impact on prices, inflation, and more.\n5. Canadian Household Confidence and Interest Rate Decision\nFinally, at 21:45, the Bank of Canada will announce its interest rate decision. Faced with global economic uncertainties, the direction of Canada's monetary policy is also worth noting.\nSummary and Reflections\nBehind every data release lies profound economic significance. Although the market is full of uncertainties, for investors, in-depth analysis of each event can help build more comprehensive and scientific investment strategies.\nOn the upcoming April 16, paying attention to the changes in these financial events and data will be a shared responsibility for all market participants. We hope everyone can find their footing in this turbulent market, rationally respond to various information, and make wise choices. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/www_sohu_com_icon_new.jpg",
            "link": "https://www.sohu.com/a/884753377_122066678",
            "media": "Sohu",
            "publish_date": "",
            "refer": "ref_2",
            "title": "April 16 Financial Events and Economic Data: Focus Points for Global Markets (Published: 2025-04-16 06:32:00)"
        },
        {
            "content": "As the global economic landscape continues to evolve, adjustments in macroeconomic policies and data releases by various countries are having a profound impact on the market. On April 14, 2025, multiple important financial events and data will emerge, which not only relate to the health of national economies but may also trigger chain reactions in global markets. This article provides a detailed analysis of the day's major financial events, offering deep insights for readers.\n1. Introduction to Market Phenomena: The Fluid Changes in the Global Economic Landscape\nIn recent years, the global economic situation has been highly volatile, with financial markets reacting sensitively. The economic data for Q1 2025 is about to be released, and its impact on foreign exchange, stock, and commodity markets cannot be underestimated. Especially amid the mutual influences of major economies like China, the U.S., and Europe, market investors are increasingly focusing on real-time data and changes in economic policies.\n2. State Council Information Office Releases Q1 Import and Export Data\nTime: 10:00\nChina's State Council Information Office will hold a press conference to release Q1 import and export data. This information is not just about numbers but also reveals the close ties between China and the global economic chain. The growth or decline in import and export data will directly affect GDP growth expectations, thereby influencing investor confidence.\nBy analyzing recent trade relations, it is evident that China plays a pivotal role in the global supply chain. Especially in the face of international economic fluctuations, China's import and export situation will serve as a market barometer, closely linked to trends in international commodity prices and the foreign exchange market. Achieving strong import and export performance will not only boost market confidence but also inject more vitality into domestic economic development.\n3. Canadian Wholesale Sales Data and Economic Analysis\nTime: 20:30\nOn the same day, Canada will release its February wholesale sales monthly rate data. Wholesale sales are an important indicator of retail demand expectations, directly reflecting the health of the economy.\nAnalysts point out that if the data shows strong sales growth, it may indicate a recovery in Canadian consumer confidence, while the opposite would suggest potential economic risks. Particularly during a rate hike cycle, changes in consumer spending will be especially important, influencing the Bank of Canada's monetary policy decisions and, to some extent, global market liquidity.\n4. Importance of U.S. Inflation Expectations\nTime: 23:00\nNext, the New York Fed will release the one-year inflation expectations data for March. This data is a key reference for assessing U.S. consumers' and businesses' views on future economic conditions.\nCurrently, the Fed is at a crossroads in addressing high inflation and reducing interest rate pressures. Rising inflation expectations may signal the need for further rate hikes, affecting capital market pricing. Conversely, if inflation shows signs of easing, it would indicate that policies are taking effect, potentially boosting market risk appetite.\n5. OPEC Crude Oil Market Report: A Barometer for Energy Prices\nTime: To be determined\nAround the same time, OPEC will release its monthly crude oil market report. This report typically influences global crude oil price trends.\n·\nExpand full text\nIn 2025, global energy agencies are generally facing the dual challenges of balancing market supply and demand with emission reduction policies. As global warming issues intensify, the balance between new energy development and traditional energy consumption in various countries will become increasingly important. Especially for countries reliant on oil and gas resources, this report will be a key reference for observing market adjustments.\n6. Views from Fed Officials: A Critical Moment for Market Interpretation\nTime: 00:00 and 06:00 the next day\nFed officials will deliver speeches on navigating economic uncertainties and sharing the Fed's role. These speeches will not only outline the future direction of Fed policy but also cast a shadow over economic uncertainties.\nThe market's focus on the Fed will influence global investor expectations, especially under the current intertwined scenarios of high inflation and high interest rates. Investors need to closely monitor the Fed's stance, which will not only affect the dollar's trajectory but also directly relate to the stability of commodity markets and the global financial system.\nConclusion: Financial Signals Investors Cannot Ignore\nIn summary, the major financial events and economic data on April 14, 2025, herald new changes in the global economy, warranting close attention from investors. On this day, the release of economic data and policy interpretations will intertwine, forming a barometer for market trends.\nWhile the global economy remains uncertain, investors must stay vigilant, monitor real-time information changes, and respond promptly. At the same time, diversifying investments and controlling risks will be wise choices to cope with market fluctuations. How to find stable investment opportunities in a changing economic environment is a crucial question every investor must consider.\nFinally, we hope that under the impetus of these financial events and data, a clearer and more stable development environment can emerge. Let us stay tuned and seize every possible opportunity. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/www_sohu_com_icon_new.jpg",
            "link": "https://www.sohu.com/a/883713425_122066678",
            "media": "Sohu",
            "publish_date": "",
            "refer": "ref_3",
            "title": "Global Financial Hotspots: Key Events and Data Analysis Not to Be Missed on April 14 (Published: 2025-04-14 07:25:00)"
        },
        {
            "content": "As 2024 approaches, the winds of change in financial markets will continue to shape the global economic landscape. In this upcoming new year, we need to closely monitor several key events that will influence financial developments in 2025. This concerns not only investors but also every individual who cares about their financial future. This article presents the financial calendar for 2025, covering key events, policy changes, and investment opportunities in major global economies, helping you navigate your financial direction.\n1. Introduction: Where is the Future of Financial Markets?\nFinancial markets are a barometer of economic operations, with national fiscal policies, monetary policies, trade wars, and economic growth directly influencing market trends. In 2025, the economic situations in the U.S., China, and Europe will be the focus of many investors. Against the backdrop of globalization, the relationships between these regions have become particularly complex. Let’s start with the key financial events of 2025.\n2. Key Events in the 2025 Financial Calendar\nChina's Two Sessions (March)\nThe annual National People's Congress and the Chinese People's Political Consultative Conference are important indicators of China's economic policies. The 2025 Two Sessions will discuss new national economic strategies and policies. This is a critical time for investors to watch, as new policy directions will directly affect investment opportunities in the Chinese market.\nFed Interest Rate Decision Meetings (Quarterly)\nThe Fed's interest rate decisions impact not only the U.S. but also cause ripples in the global economy. In 2025, the Fed's monetary policy may undergo significant changes, and investors need to stay alert to how interest rate changes affect capital flows and stock markets.\nEuropean Central Bank Policy Meetings (Quarterly)\nAs another major economy, Europe's policy direction will directly influence global markets. The ECB may adopt tightening or easing policies during the year, and investors should monitor its relationship with inflation and economic growth.\nGlobal Climate Summit (November)\nAs climate change issues become increasingly severe, the international community has recognized the importance of environmental policies for economic growth. The 2025 Climate Summit will discuss global emission reduction targets and their impact on national economies, offering opportunities for sustainable investments.\nG20 Summit (October)\nLeaders from various countries will gather to discuss global economic and trade cooperation issues. The G20 Summit is not only a platform for policy dialogue but also a key signal for investments among economies. New cooperation frameworks and economic policies may be introduced during the summit.\nChina's Two Sessions (March)\nThe annual National People's Congress and the Chinese People's Political Consultative Conference are important indicators of China's economic policies. The 2025 Two Sessions will discuss new national economic strategies and policies. This is a critical time for investors to watch, as new policy directions will directly affect investment opportunities in the Chinese market.\nFed Interest Rate Decision Meetings (Quarterly)\nThe Fed's interest rate decisions impact not only the U.S. but also cause ripples in the global economy. In 2025, the Fed's monetary policy may undergo significant changes, and investors need to stay alert to how interest rate changes affect capital flows and stock markets.\n·\nExpand full text\nEuropean Central Bank Policy Meetings (Quarterly)\nAs another major economy, Europe's policy direction will directly influence global markets. The ECB may adopt tightening or easing policies during the year, and investors should monitor its relationship with inflation and economic growth.\nGlobal Climate Summit (November)\nAs climate change issues become increasingly severe, the international community has recognized the importance of environmental policies for economic growth. The 2025 Climate Summit will discuss global emission reduction targets and their impact on national economies, offering opportunities for sustainable investments.\nG20 Summit (October)\nLeaders from various countries will gather to discuss global economic and trade cooperation issues. The G20 Summit is not only a platform for policy dialogue but also a key signal for investments among economies. New cooperation frameworks and economic policies may be introduced during the summit.\nIn 2025, we need to conduct in-depth analysis of policy trends in various countries to assess their impact on the market.\nChina\nIn 2025, China is expected to further promote \"digital economy\" and \"green economy\" policies, thereby boosting the development of technology and environmental industries. Investors can seek to increase investments in these areas.\nUnited States\nThe Fed may continue to raise interest rates while seeking a moderate inflation target, potentially causing short-term stock market volatility. In the long run, the technology and healthcare sectors will continue to receive strong support, making them suitable for investor focus.\nEuropean Union\nFacing ongoing economic growth pressures, the ECB may adopt moderately accommodative policies, with manufacturing and services sectors receiving more support, making them suitable for related investments.\nChina\nIn 2025, China is expected to further promote \"digital economy\" and \"green economy\" policies, thereby boosting the development of technology and environmental industries. Investors can seek to increase investments in these areas.\nUnited States\nThe Fed may continue to raise interest rates while seeking a moderate inflation target, potentially causing short-term stock market volatility. In the long run, the technology and healthcare sectors will continue to receive strong support, making them suitable for investor focus.\nEuropean Union\nFacing ongoing economic growth pressures, the ECB may adopt moderately accommodative policies, with manufacturing and services sectors receiving more support, making them suitable for related investments.\nIn 2025, investors will still face opportunities and challenges. Identifying potential investment opportunities is crucial, and we can start from the following aspects:\nTechnology Industry\nAreas such as artificial intelligence, blockchain, and new energy will remain hot. It is recommended to focus on related companies and investment funds.\nGreen Investments\nIn response to climate change, countries are injecting more capital into renewable energy. Investors should consider companies related to green energy to capture investment opportunities.\nConsumer Market\nWith the expansion of the middle class, the consumer market will continue to grow, especially in China, where household consumption will become a key driver of economic growth.\nTechnology Industry\nAreas such as artificial intelligence, blockchain, and new energy will remain hot. It is recommended to focus on related companies and investment funds.\nGreen Investments\nIn response to climate change, countries are injecting more capital into renewable energy. Investors should consider companies related to green energy to capture investment opportunities.\nConsumer Market\nWith the expansion of the middle class, the consumer market will continue to grow, especially in China, where household consumption will become a key driver of economic growth.\nHowever, investors must also be wary of potential risks, including policy risks and market volatility. Thorough research before investing is essential.\n5. Summary and Expectations\nLooking ahead to 2025, the global economy still faces many uncertainties, but opportunities and challenges coexist. As investors, we must continuously learn and adapt to new situations, track global economic dynamics, and adjust our investment portfolios accordingly. We hope this 2025 financial calendar provides valuable insights, and we look forward to a promising future.\nIn the coming year, we remind everyone to stay alert, keep up with financial news, especially policy changes, international relations, and market trends. We hope every reader can find suitable investment opportunities in 2025 and reap abundant rewards!\nReturn to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/www_sohu_com_icon_new.jpg",
            "link": "https://www.sohu.com/a/844117122_122066678",
            "media": "Sohu",
            "publish_date": "",
            "refer": "ref_4",
            "title": "2025 Financial Calendar: Major Financial Events and Investment Opportunities Not to Be Missed (Published: 2025-01-01 10:37:00)"
        },
        {
            "content": "April 16, 2025, will be an important financial day, with investors needing to closely monitor a series of key economic indicators and events. First, at 07:10, Fed Governor Lisa Cook will deliver a speech, expected to provide guidance on future monetary policy directions.\nIn the morning, China will release multiple heavyweight economic data. At 09:30, the monthly report on residential sales prices in 70 large and medium-sized cities will be released, revealing trends in the real estate market. At 10:00, China's Q1 GDP annual rate, March total retail sales of consumer goods, and industrial added value above a designated size will be released, once again drawing market attention to the strength of economic recovery.\nIn the afternoon, the UK will release its March Consumer Price Index (CPI) monthly rate at 14:00, an indicator that will influence market views on inflation and central bank decisions. Following this, at 16:00, the Eurozone's seasonally adjusted current account and the final annual CPI rate for March will also be released, further affecting assessments of the Eurozone's economic situation.\nIn the evening, the World Trade Organization will release its global trade outlook report at 20:00. Subsequently, at 20:30, the U.S. March retail sales monthly rate and at 21:15 the industrial production monthly rate will be released, drawing high market attention to the health of the U.S. economy.\nAt 21:45, the Bank of Canada will announce its latest interest rate decision, with investors closely watching its impact on future monetary policy. Finally, at 22:00, the U.S. April NAHB housing market index will debut, along with the release of the U.S. February business inventory monthly rate data, providing investors with a more comprehensive market understanding.\nAt 22:30, the U.S. EIA crude oil inventory data for the week ending April 11 will be released, potentially significantly impacting oil prices and related market sentiment. Overall, April 16 will be a busy day for financial activities, holding great importance for market analysis and decision-making. Return to Sohu for more.",
            "icon": "https://sfile.chatglm.cn/searchImage/www_sohu_com_icon_new.jpg",
            "link": "https://www.sohu.com/a/884753440_122006510",
            "media": "Sohu",
            "publish_date": "",
            "refer": "ref_5",
            "title": "April 16 Financial Focus: Overview of Global Economic Data and Monetary Policy Decisions (Published: 2025-04-16 06:32:00)"
        }
    ]
}
```


# GLM-5.3-Flash/FlashX
Source: https://docs.z.ai/guides/vlm/glm-5.3-flash



> [!TIP]
>
> GLM-5.3-Flash is now fully available on the [GLM Coding Plan](https://z.ai/subscribe). With native multimodal capabilities and 3× the quota, it delivers a smoother and more cost-effective coding experience.

## Model Overview

**GLM-5.3-Flash/GLM-5.3-FlashX** is the first native multimodal model in the GLM-5 series, delivering stronger intelligence than GLM-5.2 at an exceptionally low cost.

* **Highly Efficient Hybrid Architecture**

GLM-5.3-Flash has 320B total parameters with 18B activated. As the first open-source frontier model to combine sparse and linear attention, it significantly cuts computation and serving costs while preserving long-context quality — reducing attention computation and KV cache by 3.01× and 4.44× versus GLM-5.3.

* **Native Multimodal Visual Coding**

Visual capabilities are built into the coding loop: the model observes interfaces, rendered results, and interaction feedback to continuously test and improve its work. It coordinates tasks across code, browsers, and GUIs — from frontend and game development to Blender 3D scenes and real-world operation via BUA and CUA.

* **A Professional Work Partner Beyond Coding**

Beyond coding, GLM-5.3-Flash supports professional workflows such as Office tasks, financial research, and document processing. It autonomously breaks down goals, invokes tools, and reviews outputs — from research and analysis to finished PPTX, PDF, DOCX, and XLSX deliverables.

**GLM-5.3-FlashX** is now live, delivering inference speeds of **200 tokens/s** for faster responses and a smoother experience.

[↗ blog](https://z.ai/blog/glm-5.3-flash)

> **Input Modality**
>
> Video / Image / Text / File

> **Output Modality**
>
> Text

> **Context Length**
>
> 1M

> **Maximum Output Tokens**
>
> 128K

## How to Use

#### **Model API**

* **Model Code**：`glm-5.3-flash`/`glm-5.3-flashx`/`glm-5.3-flashx-legal`
* **API Documentation**：[Chat Completion API](https://docs.z.ai/api-reference/introduction)
* **Parameter Settings**：Text parameters are consistent with GLM-5.3, with support for a 1M-token context window.
* **Image Parameters**：Add a content block with `type: image_url` to` messages[].content[]`, and pass the image URL (recommended) or a Base64 Data URL through `image_url.url`. Multiple images can be added by including multiple `image_url` content blocks.
* **Recommended Settings**： `temperature: 1`, `top_p: 0.95`, and `reasoning_effort: max`. `thinking.type` only supports `enabled`; we recommend setting `thinking.clear_thinking: false`. For streaming requests, we recommend enabling both stream: true and `tool_stream: true`.

#### **GLM Coding Plan**

* Now fully available, GLM-5.3-Flash can be used with your preferred tools, with 3× the available quota compared with GLM-5.3. (GLM-5.3-FlashX is not yet available on the plan)
* The new GLM Coding Plan adopts a points-based quota system with transparent usage limits. Model calls made during off-peak hours, including all day on weekends, consume only 50% of the standard points.\
  Subscribe now: [Personal Plan](https://z.ai/subscribe?plantype=individual)、[Team Plan](https://z.ai/subscribe?plantype=team)

## Capabilities

* [Thinking Mode](https://docs.z.ai/guides/capabilities/thinking-mode): Provides multiple thinking modes to address different task requirements. `thinking.type` only supports `enabled`; thinking cannot be disabled.
* [Streaming Output](https://docs.z.ai/guides/capabilities/stream-tool): Supports real-time streaming responses for an enhanced user interaction experience.
* [Function Calling](https://docs.z.ai/guides/capabilities/function-calling): Provides powerful tool-calling capabilities and supports integration with a wide range of external tools.
* [Context Caching](https://docs.z.ai/guides/capabilities/cache): Uses an intelligent caching mechanism to optimize long-conversation performance.
* [Structured Output](https://docs.z.ai/guides/capabilities/struct-output): Supports structured output formats such as JSON for seamless system integration.
* **Visual Understanding**: Native multimodal input supporting images, videos, and files.

## Best Practices

> **Vision-Driven UI Coding: From Reference Assets to a Complete Application**
>
> GLM-5.3-Flash can directly transform screenshots, multiple page images, website URLs, or screen recordings into polished, functional applications. It goes beyond reproducing colors and layouts by understanding page relationships, shared components, design systems, interaction states, and animation logic, enabling an end-to-end workflow from visual analysis to a complete frontend project.
>
> **Recommended Approach**:
> Provide a set of page screenshots from the same product, or a website with complex animations, and have the model create a high-fidelity reproduction:
>
> > Based on the page screenshots I provide, fully recreate this product using Next.js and TypeScript. First analyze the design system, page relationships, shared components, navigation structure, interaction states, and animation logic, then build a fully functional project. After implementation, launch the project and compare screenshots of each page against the reference images. Continuously refine differences in layout, typography, spacing, colors, image cropping, and interactions. Finally, explain which pages have been covered, how to run the project, and any remaining discrepancies.

> **Office Deliverables: From Content Organization to Visual Validation**
>
> GLM-5.3-Flash's Office capabilities go beyond content generation to simultaneously handle information structure, visual styling, charts, image cropping, and page layout. Whether creating PPTX, PDF, DOCX, or XLSX files from scratch or reproducing existing files and reference images, it can proactively identify issues such as overflow, misalignment, overlapping elements, and inconsistent styling through rendering and visual inspection.
>
> **Recommended Approach**:
> Select a real presentation or reporting task, clearly specify the audience, page count, content structure, and visual style, and have the model deliver a ready-to-use file:
>
> > Using the materials in the current directory, create a 15-page business presentation for management. First extract the key conclusions and narrative structure, then create editable charts, page layouts, image selections, and speaker notes. Do not fabricate any business data; cite the source for all referenced information. After completion, render and inspect each page, fixing text overflow, image cropping, element overlap, alignment issues, and visual inconsistencies. Finally, deliver the PPTX and PDF, and explain what has been verified and what risks remain uncovered.

> **Financial Professional Workflows: From Research Traceability to Models and Reports**
>
> GLM-5.3-Flash covers the complete workflow from financial research and analysis to valuation modeling and report generation. It can integrate information from multiple sources, retain the supporting evidence for key conclusions, distinguish disclosed facts from analytical assumptions and derived results, and translate research insights into auditable financial models and formal reports.
>
> **Recommended Approach**:
> Select a listed company that has just released its latest financial results, and provide its announcements, financial statements, and an existing model:
>
> > Based on the company's latest financial statements, announcements, and verifiable public information, conduct a comprehensive earnings analysis. Break down revenue, profit, cash flow, business segments, and key operating metrics, clearly distinguishing between company disclosures, analytical assumptions, and your own conclusions. Update the earnings forecast and valuation model, and verify the consistency of the three financial statements, formula references, and accounting definitions. Finally, deliver a PDF research report with source citations and an editable, formula-driven Excel model, and list the key risks and unverified information.

> **Video Understanding and Editing: From Long-Form Footage to Publish-Ready Videos**
>
> In an Agent environment, GLM-5.3-Flash can simultaneously understand visuals, speech, subtitles, relationships between people, and timelines, reorganizing long-form videos or multiple pieces of footage into coherent content. It can also leverage name tags, on-screen appearances, and contextual cues to continuously identify different speakers, handling tasks such as speaker attribution, story structuring, and visual matching that are difficult to accomplish with ASR alone, significantly improving the efficiency of video editing Agents.
>
> **Recommended Approach**:
> Provide a set of interview, event, or product footage and have the model complete the editing, subtitles, and final quality checks:
>
> > Edit the videos in the footage directory into a 90-second product launch recap. First create a media inventory, identify people, events, key statements, and usable shots, then design the pacing of the opening, main section, and ending. Distinguish between speakers and generate accurate subtitles, ensuring that the visuals correspond to the spoken content. After completing the editing, music, transitions, and basic color grading, check for typos, speaker attribution, audio-video synchronization, black frames, and duplicate shots. Finally, output the MP4, SRT, and editing notes.

> **3D Scene Creation: From a Single Description to a Complete Blender Project**
>
> GLM-5.3-Flash can transform spatial requirements, visual styles, and functional constraints into editable Blender projects, and continuously develop them from blockout to assets, materials, lighting, camera setup, and final rendering. More importantly, it can repeatedly render from fixed camera positions, inspect the actual visuals, identify issues, and iterate, rather than merely generating a modeling script.
>
> **Recommended Approach**:
> Provide a complete scene task that specifies spatial planning, design style, and delivery requirements:
>
> > Create a complete, editable Blender scene in the current directory depicting a restaurant and bar located on a high floor in a city. First define the artistic direction, spatial layout, asset list, and fixed camera shots, then build the blockout and render previews as early as possible. Perform at least four rounds of "build → fixed-camera render → inspect → refine → re-render," focusing on spatial scale, circulation, materials, lighting, mesh intersections, and camera composition. Finally, reopen the project from a clean environment and render the main shot. Deliver the .blend file, final images, and reproduction instructions.

> **Game Development: From Gameplay Rules to a Fully Playable Loop**
>
> GLM-5.3-Flash can extract visual language, core mechanics, state machines, and win/loss conditions from reference images and gameplay descriptions, efficiently adapt them to professional cross-platform game engines such as Godot, and implement them as genuinely playable game prototypes. Rather than simply generating scenes, it is better suited to testing whether movement, collisions, AI, scoring, levels, results, and save systems form a complete gameplay loop.
>
> **Recommended Approach**:
> Use your own or properly licensed assets and choose a focused, time-bounded Godot game task:
>
> > Using the properly licensed assets in the current directory, develop a playable cooperative cooking game prototype with Godot 4. First review the gameplay specifications, asset mappings, and reference screenshots, then implement character movement, item pickup, food preparation, serving, scoring, countdown timers, and the results flow according to milestones. Run tests after each stage and compare the result against the reference images using the same map and camera positions. Finally, ensure that a complete game session can run from start to finish, and deliver the Web build, test logs, runtime instructions, and a list of incomplete items.

> **Computer Use Closed Loop: Operating, Testing, and Refining in Real Interfaces**
>
> GLM-5.3-Flash can visually understand software interfaces and perform clicks, text input, decision-making, and continuous actions even when structured APIs are unavailable. It can test games and websites, observe existing applications, reproduce key functionality, and then operate the reproduced version to identify discrepancies, forming a closed loop of "observe → implement → use → refine."
>
> **Recommended Approach**:
> Select a local application or web product, enable `/goal` mode, and have the model reproduce it while using it:
>
> > `/goal` Use Computer Use to inspect the currently open application, identify its page structure, core functionality, primary user flows, and interaction feedback, and recreate a functional version in the current directory. After completion, operate both the original application and the recreated version, comparing their layouts, functionality, state changes, and interaction paths. Record the differences and continuously improve the implementation. Finally, complete all major workflows except login, and provide the verification results, remaining discrepancies, and instructions for running the application.

> **CAD Visual Reproduction: From Design Blueprints to Editable 3D Models**
>
> GLM-5.3-Flash can understand the main structure, hole positions, fillets, chamfers, surfaces, and assembly relationships from part photographs, sketches, or CAD blueprints, and generate parametric 3D models using code-based tools such as build123d. After generating the model, it can render multiple views, continuously compare them against the reference images, refine discrepancies, and deliver editable project assets for further development.
>
> **Recommended Approach**:
> Select a mechanical part with a clear structure and multiple recognizable features, and provide dimensional or proportional references:
>
> > Based on the part blueprint and multi-angle reference images I provide, use build123d to write parametric code that recreates the part. First identify the main structure, key dimensions, hole positions, fillets, chamfers, and symmetry relationships, and list assumptions for any dimensions that cannot be determined with certainty. After completion, render the model from angles matching the reference images and compare the proportions and structural differences item by item, continuously refining the model. Finally, deliver the Python source code, STEP file, STL file, dimensional specifications, and an interactive viewing page.

## Key Advancements

GLM-5.3-Flash incorporates several architectural improvements over GLM-5. For the first time, we introduce a hybrid architecture combining sparse and linear attention, sharply reducing long-context serving costs while preserving precise long-context capabilities. It also adopts Manifold-Constrained Hyper-Connections (mHC) to further improve scaling efficiency. Combined with our latest 30T-token multimodal pre-training corpus, these changes let GLM-5.3-Flash produce more intelligence with less compute.

Before release, we tested GLM-5.3-Flash anonymously as ox-alpha on OpenCode and OpenRouter to gather user feedback. It quickly became the most popular model of the week — with all of this traffic served on Chinese AI chips.

#### Competitive Performance at Flash Cost

GLM-5.3-Flash pushes the Pareto frontier of the Artificial Analysis Intelligence Index v4.1.1, scoring 57 at just \$0.045 per task (discounted) — a level of intelligence previously only available at roughly 10× the cost. This makes it a highly competitive default choice for a broad range of workloads.

![](https://cloud-document-converter.oss-cn-beijing.aliyuncs.com/feishu2md/20260826/1787756367779-b9g6ji.png)

Across six coding and agentic benchmarks, GLM-5.3-Flash consistently outperforms GLM-5.2, often by a wide margin — 63.4 vs. 46.2 on DeepSWE v1.1 and 48.8 vs. 26.2 on AutomationBench — while approaching Claude Opus 4.8 overall.This holds on our in-house coding evaluation as well: on Z.ai Code Bench v1.0 (run on Claude Code 2.1.207), GLM-5.3-Flash clearly outperforms GLM-5.2 at every effort level, and at max effort nearly matches Claude Opus 4.8 (29.0 vs. 29.5).

![](https://cloud-document-converter.oss-cn-beijing.aliyuncs.com/feishu2md/20260826/1787756174815-hakryv.png)

#### Architecture for Extreme Efficiency

![](https://cloud-document-converter.oss-cn-beijing.aliyuncs.com/feishu2md/20260826/1787756367780-5fet8o.png)

GLM-5.3-Flash is specifically designed for ultra-low-cost inference. Despite a similar total parameter count (320B vs. 355B), it nearly halves both the activated parameter count (18B vs. 32B) and the number of layers (45 vs. 92).

To minimize attention costs in long-context scenarios, we use a hybrid architecture combining linear and sparse attention. Linear attention captures local dependencies through state modeling, while sparse attention retrieves relevant global context through a lightweight indexer. To further reduce the latency and memory overhead of the indexer at a 1M-token context length, we introduce IndexPool, which compresses four indexer key vectors into one through weighted pooling.

To illustrate the efficiency of our architecture, we compare the per-token compute and KV cache size of GLM-5.3-Flash against GLM-5.3 and two recent open models DeepSeek-V4-Flash and Kimi-K3. For a fair comparison among different scales, we calculate the attention compute per head per layer and average KV cache size per layer (BF16). Compared with GLM-5.3, GLM-5.3-Flash reduces the attention compute and KV cache size by factors of 3.0x and 4.4x. GLM-5.3-Flash has the lowest attention compute among all models compared. The KV cache size is still slightly larger than Kimi-K3 and DeepSeek-V4-Flash, leaving further room for improvement.

The overall architecture improvements, combined with optimized pre-training corpus, enable GLM-5.3-Flash to produce more intelligence with less compute. In the table below we show the evaluation results of the base model of GLM-5.3-Flash, comparing with our previous base models and DeepSeek-V4-Flash-Base. The results show that GLM-5.3-Flash-Base remains competitive with GLM-5-Base across most benchmarks.

#### Visual Intelligence in the Coding Loop

Visual coding is not just about processing images. It expands the boundary of what coding can reach. For tasks such as frontend development, game development, and 3D simulation, the final output is not code alone, but an interface, an interaction, or a world experienced by the user. Many failures only surface through rendering, interaction, or playtesting. CUA further extends coding beyond programmable systems into visible and interactive environments. Vision therefore needs to be natively integrated into the model, enabling it to decide when to observe and use visual feedback to guide subsequent actions.

We develop data synthesis pipelines for visual coding, with a focus on self-visual judgment and test-time improvement. The resulting trajectories require the model to interact with environments, inspect its own outputs, and refine them iteratively. For frontend coding, we also explored reinforcement learning with environment feedback and further strengthened GUI judgment through agent-based verification grounded in real user flows. This extends validation beyond functional correctness to the rendered and interactive product.

#### Beyond Coding -- Your Partner at Work

Coding capabilities provide an important foundation for intelligent knowledge work, while visual intelligence extends these capabilities to a broader range of professional tasks. A substantial portion of professional activities involves interpreting heterogeneous visual and structured information, including documents, spreadsheets, presentations, dashboards, interfaces, and meeting artifacts.

Visual intelligence extends the model’s capabilities beyond code-centric environments by enabling it to jointly reason over textual, visual, and structural context. Rather than requiring users to explicitly translate their working environment into textual instructions, the model can directly interpret the artifacts associated with a task and identify relevant information. It can also assess its own outputs against the visual context and intended outcome, enabling more effective self-verification and refinement — including stronger judgments of presentation quality and aesthetics.

#### Serving at Scale on Chinese AI Chips

Over the past week, we have served GLM-5.3-Flash on a large-scale cluster of Chinese AI chips, supported by a high-bandwidth interconnect and a serving stack optimized for the underlying hardware.

To overcome the relatively limited compute and memory capacity of individual chips, we built a dedicated inference engine for this architecture on top of SGLang. Notably, this effort was accelerated by our GLM-5.3-powered infrastructure agent, which assisted engineers in developing and optimizing kernels, diagnosing performance bottlenecks, and improving the serving stack — creating a feedback loop in which the model helped optimize the system serving the model itself.

These chips are primarily constrained by memory capacity and bandwidth, especially when supporting context lengths of up to one million tokens. This calls for aggressive memory optimization, including compute-for-bandwidth and communication-for-bandwidth techniques tailored to the underlying architecture. Our stack combines intra-node tensor parallelism for Linear Attention and the LM head, ReplaySSM, W8A8 quantization, hybrid INT8/FP8/BF16 cache quantization, and Layer Split.

At cluster scale, our production-grade Encode–Prefill–Decode (EPD) disaggregated architecture separates multimodal encoding, prompt prefill, and token-by-token decoding into independently scheduled and scalable worker pools, enabling efficient and reliable serving across tens of thousands of domestically developed accelerators.

Compared with our initial baseline on the same hardware, we achieved a 3× improvement in end-to-end serving performance, reaching hardware efficiency and per-token cost comparable to mainstream NVIDIA GPUs. This demonstrates that Chinese chips can support frontier-model inference efficiently and economically at scale.

#### Conclusion

GLM-5.3-Flash shows that frontier intelligence does not have to come at frontier cost. This is not the result of any single trick, but of three layers working together: an architecture that delivers stronger capability from less compute, a richer multimodal pre-training corpus, and infrastructure co-designed with inference hardware. We are now scaling this recipe to larger models — GLM-5.3-Flash pushes the cost-performance frontier, and the lessons from building it are already shaping our next frontier model.

## GLM-5.3-FlashX-Legal

### Overview

**GLM-5.3-FlashX-Legal (GLM-Legal)** is a model optimized for professional legal workflows, with a focus on document understanding, legal reasoning, and task execution. It covers scenarios such as legal research, contract review, and legal document drafting, helping legal professionals work from source material analysis through to initial drafts.

**API access is available to selected enterprise customers by application.** [Apply for Trial Access and Business Collaboration](https://z.ai/consultation/legal)

### Specialized MCP Integrations

For a more complete legal workflow experience, we recommend using the model in an agent environment connected to specialized MCP services. The model handles document understanding, legal reasoning, and task planning, while MCP services provide external data such as laws, cases, and company information, as well as tools for creating and revising Word documents.

Choose the following MCP services based on your task requirements: [Recommended Installation List by Country/Region](/guides/legal-mcp).

### AutoClaw Legal Assistant

For a ready-to-use legal workspace, explore [AutoClaw Legal Assistant](https://autoclaw.z.ai/). It covers workflows such as legal research, contract review, and document drafting, and provides legal AI capabilities through plugins for Word and WPS.

### In Detail

#### Specialized Training for Real-World Legal Workflows

GLM-5.3-FlashX-Legal undergoes legal-domain mid-training and post-training. The model first builds a strong professional foundation using large-scale, high-quality legal data, then undergoes reinforcement learning in tens of thousands of task environments modeled on real-world legal workflows. This improves its ability to understand documents, apply legal knowledge, and produce deliverables.

Training focuses on multi-step tasks in actual legal practice, including reading materials, analyzing issues, calling tools, and drafting documents, enabling the model to better assist professionals with complex work.

#### Performance on Legal Benchmarks

On **Harvey LAB**, a benchmark for legal agents, GLM's legal capabilities achieve an **all-criteria pass rate of 19.2%**, an improvement of **12.5 percentage points** over the base model, achieving state-of-the-art performance among open-source models in this evaluation.

![](https://cloud-document-converter.oss-cn-beijing.aliyuncs.com/feishu2md/20261021/1792588276974-biln8h.png)

#### From Professional Analysis to Document Delivery

Legal tasks often require working with multiple documents, external sources, and document versions at the same time. GLM-5.3-FlashX-Legal combines document analysis with specialized tool use to organize issues requiring clarification, supporting references, and revision suggestions, then turn them into due diligence working papers, defense strategies, or initial contract revisions.

For tasks involving Word files, the Document Drafting and Document Revision MCP services support formatting, tracked changes, and comments, making it easier for legal professionals to review and continue editing.

#### Use Cases

* **Document understanding and cross-checking**: Extract key information from contracts, corporate registration records, and litigation materials, and check consistency across multiple documents.
* **Legal issue analysis**: Use facts, evidence, and research results to identify legal relationships, issues in dispute, and lines of argument.
* **Contract review and negotiation support**: Identify clause-level risks and potential conflicts, compare multiple rounds of revisions, and propose amendments based on different negotiation positions.
* **Legal document drafting**: Produce initial drafts based on task requirements and reference materials, and use document processing tools to create or revise Word files.
* **Specialized tool integration**: Work with legal data and document processing MCP services to support research, analysis, and document delivery.


# FAQ
Source: https://docs.z.ai/help/faq



#### How does cache works in GLM models?

The caching mechanism is currently undergoing open beta testing. Therefore, the underlying logic of the caching mechanism (including but not limited to: how cache hits are triggered, how long cached content is retained, etc.) has not yet been announced. A part of your request content may be cached on our cloud platform, and the storage of cached content is currently completely free. During the request process, if your request triggers a cache hit, the cached content will be calculated at approximately 1/5 of the original price.

#### Why hasn’t my account balance changed after I used the API?

* Please be advised that there is currently a processing delay in our billing system. You do not need to worry too much if you see your account balance kept the same after using api port.
* The billing history reflect daily consumption records, and therefore display the billing status from the previous day (n-1).
* Current day consumption will not be immediately visible in the billing details.

#### How can I check the rate limits?

You can check the rate limits from here: [https://z.ai/manage-apikey/rate-limits](https://z.ai/manage-apikey/rate-limits)

#### How should I recharge?

You can recharge from the billing page: [https://z.ai/manage-apikey/billing](https://z.ai/manage-apikey/billing)

#### Why can't I recharge when using credit card?

When using a credit card to recharge, please ensure that you are not using 3DS verification. 3DS verification is not supported in our platform at this moment.


# Privacy Policy
Source: https://docs.z.ai/legal-agreement/privacy-policy



**Last Update:** September 29, 2025

Thank you for using Z.ai!

JINGSHENG HENGXING TECHNOLOGY PTE.LTD, along with its affiliates (collectively referred to as "we", "our", or "us"), is committed to respecting your privacy and diligently protecting any information we gather from or about you. This Privacy Policy outlines how we process Personal Data that we collect from or about you in connection with your use of our websites, applications, and other services on Z.ai (together referred to as the "Services").

This Privacy Policy only applies to individual users and does not apply to content that we process on behalf of customers of our business offerings. If you are enterprises or developers using the API Services ("API Services") available on Z.ai, please refer to the Data Processing Addendum for API Services.

## 1. Data Controller

The Services are provided and controlled by JINGSHENG HENGXING TECHNOLOGY PTE.LTD with its registered address in 10 ANSON ROAD, #26-03, INTERNATIONAL PLAZA, SINGAPORE 0759903. JINGSHENG HENGXING TECHNOLOGY PTE.LTD is the controller and is responsible for the processing of your Personal Data as described in this Privacy Policy. If you have any questions about how we use your personal data, please contact [user\_feedback@z.ai](mailto:user_feedback@z.ai).

## 2. What Personal Data We Collect

### 2.1 We collect personal data relating to you ("Personal Data") as follows:

**Personal Data You Provide to Us Directly**

We collect Personal Data if you create an account to use our Services or communicate with us as follows:

* **Account information.** We collect information that you provide when you set up an account, such as your date of birth (where applicable), username (where applicable), email address, and password.
* **Communication Information.** When you contact us, we collect the information you send us, such as proof of identity or age, contact details, feedback or inquiries and the contents of any messages you send.
* **User Content.** This includes any text prompts, images, or other data you input. This information is processed in real-time to provide you with the Service.

**Personal Data We Receive from Your Use of the Services:**

We automatically collect certain information from you when you use the Services, including internet or other network activity information such as your IP address, unique device identifiers, and cookies.

* **Device and Network Information.** In accordance with the permissions set on your device or browser, your device or browser automatically provides us with information regarding the timing and manner in which you install, access, or use our Services. This information may include details such as your device type, operating system, browser data and referring web pages, mobile network and connection details, mobile carrier or internet service provider (ISP), time zone configuration, IP address (including location information inferred from your IP address), and various identifiers (such as device or advertising identifiers, probabilistic identifiers, and other unique personal or online identifiers).
* **Usage Information.** We collect information about your use of the Services, such as the dates and times of access, browsing history, information about the links you click, pages you view, and other information about how you use the Services, and technology on the devices you use to access the Services.
* **Log Information.** We collect information about how our Services are performing when you use them. This information includes log files. In the event that you or your device encounters an error, we may collect details about the error, the time it occurred, the feature in use at that moment, the application's state when the error happened, and any communications or content present at the time of the error.
* **Cookies & Similar Technologies.** We use cookies and similar tracking technologies to operate and provide the Service. For example, we use cookies to remember your language preferences, and for security purposes. We will obtain your consent to our use of cookies where required by law. For more information about cookie and instructions on how to set up your browser to accept, delete or disable cookies, see Chapter 10.

**Personal Data We Receive from Other Sources**

We may receive the information described in this Privacy Policy from other sources, such as:

* **Log-in, Sign-up, or Linked Services.** Where available, if you choose to sign-up or log-in to the Services using a third-party service such as Google or GitHub, or link your account to a third-party service, we may collect information from the service, such as access token.
* **Security Information.** We receive information from our trusted partners, such as security partners, to protect against fraud, abuse, and other security threats to our Services.
* **Public Information.** We may obtain publicly available information via the Internet sources in order to train our models and provide services.

## 3. How We Use Your Personal Data

We use your information to operate, provide, develop, and improve the Service, including for the following purposes set out in the table below. We may also aggregate or anonymize Personal Data in such a way that it no longer identifies you and use this data for the purposes mentioned above, such as analyzing how our Services are being utilized, enhancing and adding new features, and conducting research. We will keep and use the anonymized data in its de-identified form and will not attempt to re-identify it, unless required by law.

| Purpose | Type of Personal Data processed, depending on the processing activity | Legal basis, depending on the process activity |
| - | - | - |
| To provide, analyze, and maintain our Services | Account Information Communication Information Data We Receive from Other Sources Log Information Usage Information Device and Network Information Cookies and Similar Technologies Payment Personal Data | Where necessary to perform a contract with you, such as processing a user's prompts to provide a response. |
| To improve and develop our Services and conduct research | Account Information Communication Information Data We Receive from Other Sources Log Information Usage Information Device and Network Information Cookies and Similar Technologies | Where necessary for our legitimate interests and those of third parties and broader society, including in developing, improving, or promoting our Services, such as when we train and improve our models. |
| To communicate with you, including to send you information about our Services and events | Account Information Communication Information Log Information Usage Information Device and Network Information Cookies and Similar Technologies | Where necessary to perform a contract with you, such as processing your contact information to send you a technical announcement about the Services. Your consent when we ask for it to process your Personal Data for a specific purpose that we communicate to you, such as processing your contact information to send you certain forms of marketing communications. |
| To prevent fraud, illegal activity, or misuses of our Services, and to protect the security of our systems and Services | Account Information Communication Information Data We Receive from Other Sources Log Information Usage Information Device and Network Information Cookies and Similar Technologies Payment Personal Data | Where necessary to comply with a legal obligation. Where we are not under a specific legal obligation, where necessary for our legitimate interests and those of third parties, including in protecting our Services from abuse, fraud, or security risks, such as processing data from security partners to protect against fraud, abuse and security threats in our Services. |
| To comply with legal obligations and to protect the rights, privacy, safety, or property of our users, Z.ai, or third parties | Account Information Communication Information Data We Receive from Other Sources Log Information Usage Information Device and Network Information Cookies and Similar Technologies Payment Personal Data | Where necessary to comply with a legal obligation, such as retaining transaction information to comply with record-keeping obligations. Where we are not under a specific legal obligation, where necessary for our legitimate interests and those of third parties and broader society, including in protecting our or our affiliates', users', or third parties' rights, safety, and property, such as analyzing log data to identify fraud and abuse in our Services. |

## 4. How We Share Your Personal Data

We may disclose your personal data to the following categories of recipients:

* **Affiliates:** our affiliated companies in the Z.ai Group and/or their designated service providers, who provide data processing services necessary to provide you with our Services;
* **Vendors and Service Providers:** third party vendors and service providers partners who provide data processing services to us as necessary to provide you with our Services, or assist us in delivering our products, services, websites and platforms as well as improving and optimising the same, or who otherwise process personal data for purposes described in this Privacy Policy. Specifically, we may use analytics service providers, such as Google Analytics or other analytics platforms, that collect pseudonymized or anonymized technical and interaction data, such as device type, operating system, browser version, IP address (with anonymization applied, where applicable), page views, session duration, and on-site event interactions (e.g., link or button clicks). Such information helps us understand usage patterns, enhance user experience, and improve our website and platform functionality.
* **Business Transfers:** any business partner, investor, assignee or transferee (actual or prospective) to facilitate business asset transactions (which may extend to any merger, acquisition or any debt or asset sale) involving any of the Companies;
* **Government Authorities or Other Third Parties:** any competent law enforcement body, regulatory, government agency, court or other third party (such as our professional advisers) where we believe disclosure is necessary (i) as a matter of applicable law or regulation, (ii) to exercise, establish or defend our legal rights or so a third party can defend theirs, or (iii) to protect your vital interests or those of any other person; and
* **With your consent:** any other person, with your consent to the disclosure.

## 5. Your Rights and Choices

Depending on where you live, you may have certain rights with respect to your personal information, such as the right to know how we collect and use your personal information. You may also have the right to access, change, oppose, request a copy of your authorization, file complaints before the competent authorities, withdraw your consent, or limit our collection and use of your personal information as well as to request that we delete it, and potentially others. In certain circumstances, you can also ask us to provide additional information about our collection and use of your personal information. Please note that your exercise of certain rights may impact your ability to use some or all of Services' features and functionalities.

* **Right to know:** the right to know what personal data we process about you, including the categories of personal data, the categories of sources from which it is collected, the business or commercial purposes for collection, and the categories of third parties to whom we disclose it.
* **Access & data portability:** the right to request a copy of the personal data we process about you, subject to certain exceptions and conditions. In certain cases and subject to applicable law, you have the right to port your information.
* **Deletion:** the right to request that we delete personal data collected from you when you use our Services, subject to certain exceptions. You also are able to delete individual conversations, which will be removed immediately from your conversation history and automatically deleted from our back-end.
* **Correction:** the right to request that we correct inaccurate personal data we retain about you, subject to certain exceptions. Please note that we cannot guarantee the factual accuracy of Outputs. If Outputs contain factually inaccurate personal data relating to you, you can submit a correction request and we will make a reasonable effort to correct this information—but due to the technical complexity of our large language models, it may not always be possible for us to do so.
* **Objection:** the right to object to processing of your personal data, including profiling conducted on grounds of public or legitimate interest. In places where such a right applies, we will no longer process the personal data in case of such objection unless we demonstrate compelling legitimate grounds for the processing which override your interests, rights, and freedoms, or for the establishment, exercise or defense of legal claims.
* **Restriction:** the right to restrict our processing of your personal data in certain circumstances.
* **Withdrawal of consent.** Where our processing of your personal data is based on consent, you have the right to withdraw your consent. The withdrawal of consent will not affect the lawfulness of processing based on consent before its withdrawal.
* **Automated decision-making:** We do not engage in decision making based solely on automated processing or profiling in a manner which produces a legal effect (i.e., impacts your legal rights) or significantly affects you in a similar way (e.g., significantly affects your financial circumstances or ability to access essential goods or services).

Please contact us by using the contact information provided in this Privacy Policy if you would like to exercise any of your rights. We will respond to your request consistent with applicable law and subject to proper verification.

You can control and access some of your personal information directly through settings. For example, you can manage your chat history. Should you choose to do so, you may also delete your chat history via your settings. If you choose to delete your account, you will not be able to reactivate your account or retrieve any of the content or information in connection with your account.

## 6. The Security of Your Personal Data

The security of your personal data is a priority for us. We implement commercially reasonable technical, administrative, and physical security measures to protect your personal data from unauthorized access, theft, disclosure, alteration, or loss. We regularly assess our security practices to incorporate new technologies and methods. However, please note that no Internet or email transmission is completely secure or error-free. Specifically, emails sent to or from us may not be fully secure. As a result, we advise you to exercise caution when deciding what information to send to us via the Services or email. Additionally, we are not responsible for any bypassing of privacy settings or security measures on the Service or third-party websites.

## 7. How Long Do We Keep Your Personal Data

We retain personal data for as long as it is needed to provide our Services and for the other purposes outlined in this Privacy Policy. Additionally, we retain personal data when required to fulfill contractual and legal obligations, when we have a legitimate business interest (such as improving and developing our Services and enhancing their safety, security, and stability), and for the exercise or defense of legal claims.

Retention periods vary based on the volume, type, and sensitivity of the personal data, the purposes for which we use it, and any legal requirements, among other factors. For instance, when we process your personal data to provide you with the Services, we retain it as long as you have an account. This includes your account information, input. If you violate any of our terms, policies, or guidelines, we may immediately remove your account information, input from public view, but we will keep the personal data as needed to address the violation.

When the personal data collected is no longer required by us, we and our service providers will perform the necessary procedures for destroying, deleting, erasing, or converting it into an anonymous form as permitted or required under applicable laws.

## 8. Where We Store Your Personal Data

Depending on your location and the features of our Services that you access or use, your personal data may be transferred and processed outside the jurisdiction in which you reside, including to an affiliated company within the Z.ai Group or a third party (as outlined in section 4) located overseas.

We generally provide the Services from Singapore, and our group companies and their designated service providers are typically located in Singapore. As a result, your personal data is generally processed in Singapore.

We implement appropriate safeguards to ensure that your personal data remains protected in line with this Privacy Policy and applicable laws, and we only transfer your data abroad in compliance with applicable laws and legally recognized transfer mechanisms.

## 9. Personal Data relating to children

Our Services are not directed to, or intended for, the individual under 18. We do not knowingly collect Personal Data from children under 13. If you have reason to believe that a child under 13 has provided Personal Data to us through the Services, please contact us via the contact details set out below. We will investigate any notification and, if appropriate, delete the Personal Data from our systems. Users under 18 must have permission from their parent or guardian to use our Services.

## 10. How Do We Use Cookies and Similar Technologies

### 10.1 Cookies and similar technologies we used

We may use cookies and similar technologies to gather information resulting from access and/or use of our Services. The technologies we use may include:

I. **Cookies or Browser Cookies:** Cookies are small text files stored in your computing or other electronic devices when you visit our website and platforms for record keeping purposes. Cookies can allow us to provide you with a personalized user experience, such as by allowing us to retrieve information which you previously provided us. Cookies can be "Persistent" or "Session" Cookies. Persistent Cookies may remain on your personal computer or mobile device across browsing sessions, while Session Cookies are generally deleted as soon as you close your web browser.

We may also use cookies and similar tracking technologies to analyze usage trends, administer the website, monitor user interactions within the Services, and gain insights into how users engage with our offerings. To support these purposes, we may automatically collect certain technical information, which helps us understand aggregate usage patterns and optimize our Services. This information may include your Internet Protocol (IP) address, browser type, Internet Service Provider (ISP), referring and exit pages, files accessed on our website (e.g., HTML pages and graphics), operating system, date and time stamps, and clickstream data.

II. **Web Beacons:** Certain sections of our services may contain small electronic files known as web beacons (also referred to as clear GIFs, pixel tags, and single-pixel GIFs) that permit us, for example, to count users who have visited those pages and understand how our Services are used (for example, to assess the popularity of certain pages and services) and to improve the performance and reliability of our Services.

### 10.2 Clear/disable cookies

You can manage your cookie preferences or enable cookie-related features through the settings available in commonly used web browsers. We will respect your choices regarding the use of cookies and similar technologies in accordance with applicable laws and your expressed preferences.

Please note that certain features of our Services rely on the use of cookies. If you choose to disable or delete cookies, or do not consent to their use, some parts of the Services may not function properly or as intended. For guidance on how to manage cookies, please refer to the help documentation provided by your browser. For your convenience, we have included links to cookie management instructions for popular browsers. However, we are not responsible for the accuracy, reliability, or completeness of any third-party content linked therein.

* [Microsoft Internet Explorer](#)
* [Microsoft Edge](#)
* [Mozilla Firefox](#)
* [Google Chrome](#)
* [Safari for macOS](#)
* [Safari for iOS](#)

### 10.3 More information about cookies

For more information about cookie and instructions on how to set up your browser to accept, delete or disable cookies, see [www.allaboutcookies.org](http://www.allaboutcookies.org).

## 11. Privacy Policy Updates

We may update this Privacy Policy from time to time. When we do, we will publish an updated version and effective date on this page, unless another type of notice is required by applicable law. We recommend that you review the Privacy Policy each time you access our Services to stay informed of our privacy practices.

## 12. Contact Us

Questions, comments and requests regarding this policy should be addressed to [user\_feedback@z.ai](mailto:user_feedback@z.ai).

## Data Processing Addendum for API Services

This Data Processing Addendum ("DPA") governs the processing of the data provided by business and enterprises users ("Customers") through API Services offered by JINGSHENG HENGXING TECHNOLOGY PTE.LTD, along with its affiliates (collectively referred to as "Company") on Z.ai platform ("Z.ai") and is hereby incorporated into the Additional Terms for API Services.

Company and Customer each agree to comply with their respective obligations under applicable data privacy and data protection laws (collectively, "Data Protection Laws") in connection with the API Services.

Customer is the person that determines the purposes and means for which Customer Data (as defined below) is processed (a "Data Controller"), whereas Company processes Customer Data in accordance with the Data Controller's instructions and on behalf of the Data Controller (as a "Data Processor"). "Data Controller" and "Data Processor" also mean the equivalent concepts under Data Protection Laws.

For the purposes of this DPA, (i) "Personal Data" has the meaning assigned to the term "personal data" or "personal information" under applicable Data Protection Laws; and (ii) "Customer Data" means Personal Data that Customer provides to Company that Company processes on behalf of Customer to provide the Services. Company will process Customer Data as Customer's Data Processor to provide or maintain the Services and for the purposes set forth in this DPA, the Terms and/or in any other applicable agreements between Customer and Company.

### 1. Obligations of the Company

As a Data Processor, the Company agrees to:

a) Process Customer Data only 1) on Customer's behalf for the purpose of providing and supporting API Services; (ii) in compliance with the written instructions received from Customer; and (iii) in a manner that provides no less than the level of privacy protection required by Data Protection Laws.

b) Adopt reasonable and appropriate organizational and institutional measures to safeguard the security of Customer Data in accordance with the requirements under Data Protection Laws.

c) Provide reasonable assistance to the requests from Customer to fulfill its legal obligations in accordance with Data Protection Laws.

d) Inform Customer if Company receives any complaint or request (in particular, requests for access to, rectification or blocking of Customer Data) directly from Customer's data subjects.

e) Inform Customer if Company receives any legally binding request to disclose Customer Data to a law enforcement authority unless otherwise prohibited by laws.

### 2. Obligations of the Customer

a) Customer shall comply with all Data Protection Laws.

b) Customer represents, warrants and covenants that it has and shall maintain throughout the term all necessary rights, consents and authorizations to provide the Customer Data to Company and to authorize Company to use, disclose, retain and otherwise process Customer Data as contemplated by this DPA, the Terms and/or other processing instructions provided to Company.

c) Customer shall be solely responsible for performing any of its obligations with regard to any requests from Customer's data subjects.

d) Customer shall not provide Customer Data to Company except through agreed mechanisms.

### 3. International Data Transfer

a) Company generally provide the Services from Singapore, and our group companies and their designated service providers are typically located in Singapore. As a result, Customer Data is generally processed in Singapore.

b) Company implements appropriate safeguards to ensure that the Customer Data will only transfer Customer Data abroad in compliance with applicable Data Protection Laws and legally recognized transfer mechanisms.

### 4. Data Return and Deletion

a) This DPA shall remain in effect as long as Company carries out Customer Data processing operations on Customer's behalf or until the termination of the Terms.

b) The Company do not store any of the content the Customer or its End Users provide or generate while using our Services. This includes any texts, or other data you input. This information is processed in real-time to provide the Customer and End Users with the API Service and is not saved on our servers.

c) For Customer Data other than those provided under Section 4(b), Company will temporarily store such data for the purposes of providing the API Services or in compliance with applicable laws. The Company will delete such data after the termination of the Terms unless otherwise required by applicable laws.


# Subscriptions, Fees, and Payment
Source: https://docs.z.ai/legal-agreement/subscription-terms



## 1. Fees and Billing

* Access to certain Z.ai services or features may require payment. Applicable fees are listed on the **Z.ai**, unless otherwise communicated in writing.
* To purchase services, you must provide complete and accurate billing information (“Payment Method”). You authorize us to charge your Payment Method for the applicable fees and any taxes.
* If you proactively check the box, confirm, or otherwise select to enable automatic renewal as presented on the relevant page, you will be deemed to have authorized Z.ai or its partners (if any), upon expiration of the current subscription term, to automatically charge your payment method for the next renewal term in accordance with the fee standard applicable to such next renewal term, as then displayed on the platform or service page, including the then-effective pricing and service content. If you have not enabled automatic renewal, no automatic charge will be made.
* If you subscribe through a third-party distributor (e.g., an app store), your payment, billing, and refund terms are governed by the distributor’s policies.
* Except where required by law or explicitly stated, all payments are **non-refundable**. Please review your order carefully before confirming.

## 2. Subscriptions and Renewal

* Some services are available by subscription (“Subscription”). The scope, features, and duration of your Subscription will be described during checkout. Z.ai may update the content or features from time to time.
* Paid Subscriptions **renew automatically** at the end of each subscription term, and your Payment Method will be charged for the renewal term.
* If your Subscription has an initial minimum term (“Initial Term”), it will automatically renew for the same duration (“Renewal Term”) unless canceled.
* You may cancel your Subscription at any time using the tools provided by Z.ai. If purchased through a distributor, cancellation must follow the distributor’s process.
* To avoid renewal charges, cancel at least **24 hours before the end of the current term**. Cancellation does not entitle you to a refund, but you will continue to have access until the end of the paid term.
* For services with automatic renewal enabled, you may unilaterally stop automatic renewal at any time before the service validity period of the current subscription term expires. Z.ai may also unilaterally discontinue providing automatic renewal services to you as required by its operational strategy, in which case the subscription service term will terminate upon expiration of the current service period. In such event, Z.ai will notify you by announcement, in-site notice, or other reasonable means. Under these circumstances, the current subscription term that has already been paid for and has taken effect will remain unaffected; however, whether the service will continue to be available in the next term, how it may be renewed, and which rules will apply shall be subject to the page display or notice in effect at that time. **If an adjustment to the fee standard has a materially adverse effect on you, you should promptly discontinue the service. If you continue to use the service, you will be deemed to have accepted such adjustment.**

## 3. Service and Fee Adjustments

Z.ai reserves the right to adjust its products and pricing standards due to business adjustments, operational strategy needs, unforeseeable circumstances, changed circumstances, or other reasonable causes, including but not limited to changes in content, price adjustments, product removal, service cancellation, reduction or addition of content, or termination of all or part of the service. Z.ai will notify you by page display, announcement, in-site notice, or other reasonable means before any new pricing takes effect. Unless otherwise expressly stated by Z.ai, the price for automatic renewal shall be the price, applicable rules, and promotional arrangements displayed on the page and actually applied by the system on the date the charge is made, rather than the original price at the time of your initial subscription. If you do not agree to such adjustments, you may turn off automatic renewal before the renewal takes place. If you continue to keep automatic renewal enabled or continue using the relevant service after the new price takes effect, you will be deemed to have agreed to the new subscription price.

## 4. Usage Rules

1. **Basic rules**

* During your use, you should comply with applicable laws, regulations, policies, and industry standards. You shall not use the subscription service for any unlawful purposes, nor for engaging in any activities that are contrary to public policy and morals or infringe upon the rights and interests of others.
* You shall comply with Terms of Use, Service Agreements, usage rules and other relevant agreementsf [Z.ai](http://Z.ai).
* Violation of this agreement may result in the suspension or termination of services, and paid fees will not be refunded.

2. **Usage Scenario Restrictions**

* You understand and agree that the usage quota under GLM Coding Plan is only used within [officially supported tools](https://docs.z.ai/devpack/tool/others). If the system detects usage through unauthorized or unsupported tools (such as SDK-based access or other third-party integrations), some subscription benefits may be restricted to ensure fairness and service stability.
* You shall not use the GLM Coding Plan quota for general-purpose API access or any scenarios outside such tools, including but not limited to directly invoking model APIs from your own applications, bots, websites, SaaS products or other systems, unless you have entered into a separate written agreement with [Z.ai](http://Z.ai).
* Unless otherwise agreed in writing, you may not resell, sub-resell, repackage, aggregate, proxy or otherwise provide the GLM Coding Plan to any third party, whether on a paid or free basis, nor may you use the GLM Coding Plan to provide model capabilities as a service to third parties.

3. **Personal-Use Only**

* The GLM Coding Plan subscription is tied to a single account and is licensed only to the individual natural person associated with such account.
* You shall not share your account or subscription, or allow any other person (including but not limited to your colleagues, friends, customers or any organization) to use your GLM Coding Plan quota, whether by sharing credentials, renting, lending, transferring, sublicensing or by any other means.
* If [Z.ai](http://Z.ai) reasonably suspects that you are engaging in account sharing, bulk or automated usage on behalf of others, resale of access, or any other conduct that may harm the platform’s fair-use order or the legitimate rights and interests of third parties, [Z.ai](http://Z.ai) is entitled to take measures including, without limitation, restricting certain features, reducing or limiting your usage quota, suspending or terminating the service, reclaiming any remaining quota, and pursuing further liability in accordance with this Agreement and applicable laws.

## 5. Limitation of Liability

* Z.ai shall exert reasonable efforts to ensure the stability and security of the subscription service. However, Z.ai **does not warrant** that the service will be entirely error-free or meet your expectations. Z.ai **shall not be liable** for service interruptions caused by **force majeure**, network failures, or third-party actions.
* Unless otherwise expressly required by applicable laws or regulations, our total direct liability to you, regardless of the cause of action or form of claim, shall not exceed **the total amount you have spent in the most recent calendar month.**

## 6. Amendment and Termination of the Agreement

* Z.ai reserves the right to amend this Agreement in accordance with changes in laws, regulations, and regulatory policies, market conditions, service needs, or business strategies. Any amended version of this Agreement will be published on the platform and will take effect upon publication.
* If you continue to use the subscription services after the Agreement has been amended, you will be deemed to have accepted the amended Agreement.
* If you violate any applicable laws, regulations, policies, the User Agreement, the Usage Guidelines, or any other platform rules or notices, Z.ai shall have the right to suspend or terminate your access to and use of the services.

## 7. Other Terms

* Additional services or features may require separate fees, which you will have the chance to review and accept before being charged.
* Z.ai is not liable for errors or issues caused by third-party payment processors.


# Terms of Use
Source: https://docs.z.ai/legal-agreement/terms-of-use



**Last Update:** April 14, 2026

Welcome to Z.ai!

**These terms of use (“Terms”) form an agreement between JINGSHENG HENGXING TECHNOLOGY PTE.LTD (“we”, “our”, or “us”) and you (“you”) as a user of the Z.ai accessible via z.ai/model-api (“Services”), governing your access and use of Z.ai for services or tools for your own purposes, internal organizational use, or for the benefit of end users. If you are enterprises or developers using the API services (“API Services”) available on Z.ai, you are also subject to the Additional Terms for API Services (“Additional Terms”).**

Please ensure that you carefully read and fully understand all the terms and conditions herein. These Terms, along with our Privacy Policy, guidelines, rules, interpretations, announcements, notices, and other content that have been published or may be published in the future by Z.ai (**“Platform Rules”**), constitute the terms and conditions you must adhere to when using Z.ai. If you are a minor under the laws of your country or state of residence, you should review these Terms with your legal guardian.

By using our Services or confirming your acceptance of these Terms (whichever occurs first), you acknowledge and agree to be legally bound by these Terms, which are subject to change and may be updated from time to time. Notification of changes to these Terms may be provided by us through posting the updated Terms on our website at Z.ai, by email to you, or via any other method we deem suitable. Should you object to any such changes, you must promptly notify us in writing to terminate your account and immediately cease all access to and use of the Services. Your continued account maintenance or use of the Services after such notification will be deemed as your acceptance of the amended Terms. If you have any questions regarding these Terms, you may consult Z.ai customer service.

Except as otherwise expressly provided, any new functions or products provided by Z.ai and their associated services shall be unconditionally governed by these Terms.

## I. Definition

**"Z.ai"** shall mean the platform operated by JINGSHENG HENGXING TECHNOLOGY PTE.LTD that provides generative AI services (such as API services) to individual and business users, including new service types we may develop from time to time. (also referred to as the **“Services”**).

**"Platform Rules"** shall mean all rules, interpretations, announcements, notices, and other content that have been published or may be published in the future by the Z.ai.

## II. Registration and Security

The Z.ai may be accessed without registration. However, to use the core services or features of Z.ai, you may be required to apply for an account and log into our Services using your account.

Please keep your created key(s) securely, prevent any form of leakage, do not share or publicly disclose your key(s) with others, and do not expose your key(s) in browser or other client-side code. You will be solely responsible for any fees incurred and losses caused by sharing or leaking your key(s).

We have the right to restrict, suspend or terminate your access to or use of the Z.ai or certain Services (or any features within the Services) if in our sole and absolute opinion, you are in breach of the Terms.

## III. Using Z.ai Services

1. You have the right to use the Services provided by the Z.ai in accordance with the terms and conditions of these Terms and Platform Rules. As generative artificial intelligence technology, models, and products continue to evolve, along with changes in laws and regulations, we may add, upgrade, modify, suspend, or terminate services, or make necessary adjustments to the technology, method, and performance of the Services. If such changes may affect your rights, we will promptly notify you through in-site notifications, website announcements, or via email. If you continue to use the Services after the changes, it indicates your acceptance of the adjusted service content.

2. As a user of our Service, you should comply with the requirements of these Terms, adhere to principles of voluntariness, equality, fairness, and good faith, not use this service to infringe upon the legitimate rights and interests of others or to seek improper gains, and not disrupt the normal order of the internet platform. You should comply with any and all applicable laws, regulations and policies of jurisdictions including but not limited to the United States, the European Union and/or other government authorities.

3. When using Z.ai services , you should comply with relevant laws, regulations, and platform rules, ensuring that you possess the necessary qualifications and capacities required by laws and regulations in the jurisdiction where the services are used. You are prohibited from engaging in or facilitating any of the following actions or any other illegal activities, including but not limited to:

   a) Using the Services in any way that compromises or threatens critical infrastructure or national security.

   b) Harming national interests of any country or jurisdiction.

   c) Inciting racial discrimination or hate crimes, or inciting, facilitating, or promoting violent extremism, terrorism, or hateful conduct.

   d) Undermining freedom of religion or inciting discrimination based on religious beliefs.

   e) Using the Service to create, publish or disseminate false information that may mislead or harm the public or any individual.

   f) Creating or disseminating obscene, pornographic, violent, murderous, terroristic, or criminal incitement content.

   g) Submitting, creating or disseminating content that is sexually explicit, suggestive, visually shocking, or otherwise disturbing.

   h) Using the Z.ai service in any manner that violates the rights of any third party, including but not limited to fraud, scams, spamming (SMS/email), or engaging in disrespectful, aggressive, hateful, abusive, insulting, or defamatory conduct toward any third party.

   i) Without proper authorization, you may not provide to Z.ai or collect, process, disclose, infer, or generate the personal information of others through Z.ai. You are strictly prohibited from using the Service in any manner that may infringe upon the personal information protection rights and privacy of any third party, including but not limited to the unauthorized collection, storage, or dissemination of any individual’s personal data, home address, phone number, email address, identification details (such as national identification number, social security number, passport number), or credit card information.

   j) Using the Services in any manner that endangers or exploits children, including any form of child sexual abuse or exploitation, as well as developing applications, tools, products, or services that are inappropriate for minors.

   k) Promote, encourage, or provide instructions on self-harm, suicide, or other dangerous activities that could lead to serious harm or death.

   l) Directly or indirectly, manipulate, redirect, divert, or hijack traffic, readership, or engagement related to Z.ai’s products and services, whether independently or in collaboration with third parties.

   m) Conducting any other prohibited behavior by applicable laws or regulations.

4. You should not engage in any activities that endanger the network and system security of Z.ai during your use of the Services, including but not limited to:

   a) Unauthorized access to networks, interference with normal network functions, data theft, intentional dissemination of malicious programs or viruses, or any other activities that endanger network security and order are strictly prohibited.

   b) The use of deep linking, page scraping, social bots, spiders, or other automated means to access this service or third-party services through this service, for the purpose of accessing, obtaining, or monitoring any unauthorized data, content, or servers/accounts, is strictly forbidden.

   c) Providing technical support, advertising, payment processing, or any other assistance to facilitate activities that endanger network security is prohibited.

   d) Reverse engineering, decompiling, disassembling, or attempting to extract data from Z.ai algorithms, source code, mechanisms, or any related components is strictly prohibited.

   e) You may not bypass, circumvent, or disseminate methods or techniques designed to bypass or circumvent the security measures implemented by Z.ai to prevent actions prohibited under these Terms.

   f) You may not use Z.ai to develop, train, or enhance any algorithms, models, or technologies that directly or indirectly compete with us is prohibited.

   g) Any other usage that may harm our interests is strictly forbidden.

5. You should comply with the following rules when using the service:

   a) You may not use this service dishonestly, including but not limited to falsely claiming that AI generated content (“Outputs”) is human-created, disguising human-created works as Outputs, engaging in deceptive advertising, marketing, transactions, fabricating discussions, generating fake public opinion trends, or manipulating engagement metrics such as readership, views, likes, plays, shares, or inducing interactions such as tipping, liking, or voting.

   b) You may not distribute or publish inappropriate Outputs using Z.ai through any means, including but not limited to copying, taking screenshots, screen recording, screen sharing, or live streaming.

   c) Unless otherwise agreed, you may not remove, modify, or obscure any AI identifiers added to Outputs by Z.ai, regardless of the form in which such identifiers are presented.

   d) You shall adopt technical measures to add identifiers and maintain logs in compliance with applicable laws. AI-generated Outputs  shall be prominently marked at reasonable locations to indicate that it is generated by AI.

   e) You may not violate these Terms, any applicable law, or any rights of third parties, including but not limited to intellectual property rights, confidentiality obligations, or prior contractual rights.

6. You should follow these rules in the scenarios listed below:

   a) You shall not use Z.ai to conduct any decision-making activities, including but not limited to high-risk automated decision-making in areas that may materially affect the safety, rights and interests, or well‑being of individuals and society, such as health, education, credit, and critical infrastructure management.

   b) You may not use Z.ai for any services that require subject qualification or professional review, or as a substitute for professional services, including but not limited to professional fields such as medical care, legal services, news reporting, education, and investment and financial management.

   c) You may not use the Services to generate content intended for or related to political campaigns.

   d) If you create virtual characters using the Z.ai, you shall ensure that the virtual characters’ images, avatars, names, character settings, and data and content in the creation process do not infringe third-party rights, and shall not use virtual characters to engage in any acts that may constitute unfair competition.

   e) You may not steal or misappropriate domain names, websites, links, text, images, videos, data or any other content from Z.ai’s products and services, whether directly or indirectly, through methods including but not limited to link hijacking, redundant extraction, unauthorized scraping, simulated downloads, deep linking, or fraudulent registration. Additionally, you may not delete, obscure, or alter domain names, platform-specific identifiers, trademarks, or other proprietary notices in any manner.

7. If you discover any vulnerabilities or defects in the services provided by Z.ai, you are encouraged to report them through the feedback feature available within the interface or through the designated contact methods specified in these Terms. You must not publicly disclose, distribute, or publish any vulnerabilities or defects through the internet or other channels.

8. If you add new data, during your use of the Z.ai and models within the platform, the resulting responsibilities shall be borne by yourself. If you suffer damage after training, fine-tuning and development and claim that we should bear the responsibility, you shall prove that the damage is unrelated to your training, fine-tuning and development, otherwise we shall be exempted from liability for the damage.

9. If you use our Services to provide services to any third party, you assume full responsibility for all associated rights, obligations, and liabilities. We shall not be liable for any losses resulting from such use.

10. Third Party Services

    a) The Services may include third-party services, such as third-party plugins, APIs, or other applications. You acknowledge and agree that these third-party services are provided to you directly by the respective third-party service providers. You must comply with the terms and conditions set forth by these third parties (such as Terms of Use, Privacy Policy, Content Policy etc.), and any disputes arising from your use of such services must be resolved independently between you and the third-party provider.

    b) You acknowledge and agree that your use of third-party data,models, plug-ins, APIs, or other applications is subject to the terms and agreements of the respective third-party providers. If your Input contains illegal or inappropriate content, certain third-party services may refuse to generate Outputs for you.

11. If you use our service in the United States, you, together with your affiliates, employees, service providers, end users and any third parties under your direction or control, shall not use the Services to:process, store, or transmit any technical data, software, or content subject to the ITAR administered by the U.S. Department of State, or the EAR administered by the U.S. Department of Commerce;

    a) process, store, or transmit any Protected Health Information (PHI) or other data regulated under the Health Insurance Portability and Accountability Act （HIPAA）, its implementing regulations, or any amendments thereto;

    b) process, store, or transmit any Nonpublic Personal Information (NPI) subject to the GLBA, or any regulations promulgated thereunder; or

    c) process or store personal information of children under the age of 13 in a manner that violates COPPA or any regulations promulgated thereunder.

## IV. Content

1. **User Content.** From time to time, and at our sole discretion, we may permit you to submit text, documents, or other materials to our Services for processing (**“Prompts”**). Prompts and Outputs are collectively referred to herein as “User Content”.

2. You assume full responsibility for all Prompts submitted to our Services. By submitting Prompts, you represent and warrant that:

   a) You own or hold all necessary rights, licenses, consents, and permissions to submit your Prompts to us. You hereby grant us all rights required to process your Prompts for the purposes set forth in these Terms, including the right to process on such Prompts in compliance with applicable laws and these Terms.

   b) All content included in your Prompts is either owned by you or that you have obtained all valid authorizations, licenses, and permissions from the relevant rights holders to submit, share, and allow us to process such content as provided under these Terms.

   c) Any Prompts provided by you to us shall not infringe any third party’s intellectual property rights, trade secrets, or other legitimate rights and interests. You further warrant that you shall not use such User content for any purpose that may constitute an infringement of third party intellectual property rights.

   d) You have not and will not, in any way, violate data protection laws or regulations, infringe upon the legal rights of any individual or entity, or use User Content for any unlawful purposes. You also ensure compliance with all applicable laws and regulations in protecting the personal data of your end-users, including issuing legally effective privacy notices to end-users. You confirm that you have obtained and retained any required consent from end-users for authorizing us to process your prompts.

3. You further agree and acknowledge that:

   a) Service Improvement and Development. For individual users, we reserve the right to process any User Content to improve our existing Services and/or to develop new products and services, including for our internal business operations and for the benefit of other customers. You hereby explicitly authorize and consent to our: (i) utilization and storage of User Content that does not constitute personal data for the purpose of developing and improving our machine learning and artificial intelligence technologies; and (ii) processing and storage of such User Content in locations outside of the jurisdiction where you access or use the Services.

   b) Data Collection and Usage. We may collect performance and usage metrics and data pertaining to your access and/or utilization of the Services, which may include, but is not limited to, model versions, inferences, upload times, as well as processing, diagnostic, and other technical data. We retain the right to use such metrics and data to enhance the quality, functionality, and other aspects of the Services and their underlying technologies.

   c) Public Sharing of User Content. Should you activate any features or functionalities of the Services (as applicable, and which we may, at our sole discretion, make available from time to time) that enable the public sharing of any User Content, you hereby agree and acknowledge that, notwithstanding any other provision within these Terms, any individual or entity with whom you have shared such content shall have the unrestricted right to access, use, save, reproduce, distribute, display, and transmit such content without limitation. You further agree that we shall not be held liable to you for any consequences arising from such public sharing and subsequent use of the User Content by third parties.

4. **Ownership of User Content.** As between you and us, and to the fullest extent permitted by applicable law, you retain all rights, title, and interest in the Prompts you submit and the Outputs generated specifically at your request and provided to you as a response to your submitted Prompts. You acknowledge and agree that the Outputs generated may lack uniqueness and could be similar or identical to Outputs generated for other users or any Third Party. Consequently, your rights in specific Outputs, if any, may not extend to Outputs generated at the request of other users or any Third Party.

5. **Our Use of User Content.** For individual users, we may use User Content to provide, maintain, develop, and improve our Services, comply with laws, enforce our terms, and ensure Service safety. Contact us with any concerns about our use of User Content. You grant us and our affiliates an unconditional, irrevocable, non-exclusive, royalty-free, fully transferable, sub-licensable, perpetual, worldwide license to access, use, host, modify, communicate, reproduce, adapt, create derivative works from, publish, perform, and distribute your User Content, and to authorize others to do so, in any format and on any platform. You also grant us a similar license to use your username and profile picture to identify you as the source of User Content. We may disclose your identity to third parties, including advisors, subcontractors, and those claiming User Content infringes their rights.

6. You waive any rights to prior review or approval of marketing materials related to User Content, and all privacy and publicity rights. You also waive any moral rights under the Copyright Act 2021 of Singapore, the rights described in Article 6bis of the Berne Convention for Protection of Literary and Artistic Works 1886, or any other rights of a similar laws. If any such rights cannot be waived, you agree not to assert them.

7. **For enterprises and developers using API Services, we will not use your User Content for developing or improving Services unless you explicitly agree to such use. Please refer to Additional Terms for further details.**

## V. IP Rights and Usage Terms

1. Each party retains ownership of their intellectual property as it stood prior to using the Services. Neither party’s rights will be transferred to the other, nor will they become jointly owned, as a result of fulfilling these Terms. To avoid any doubt, the ownership and intellectual property rights of our models, including but not limited to its parameters, algorithms, code, and framework structure, belong to us.

2. Without our prior written consent, you shall not at any time, in any manner or through any means (whether directly or indirectly, expressly or impliedly, including but not limited to any marketing, advertising, promotion, marketing or other purposes):  use, refer to, publish, copy, display (publicly or non-publicly), mention, register or apply to register: (i) any trademarks (whether registered or not), trade names, names, logos or identifiers owned or authorized for use by us or our affiliates, or enable third parties to identify us or any of our affiliates through specific descriptions; or  (ii) any names, trademarks, logos, words, terms, symbols or any combinations thereof in any language that are similar to the foregoing. If you violate the foregoing agreement, we have the right to require you to immediately delete the relevant content and be responsible for eliminating the impact and compensating all losses of us.

3. You shall not use Z.ai in any manner that infringes the copyrights or other intellectual property rights of any third party. Any content you input that is owned by a third party shall remain the property of its respective owner.

4. We are not responsible for any intellectual property infringement resulting from your use of Z.ai in the following ways:

   a) Using infringing or potentially infringing content when you knew or should have known about its infringement.

   b) Using unauthorized content within products or services on Z.ai.

   c) Inducing others to infringe on intellectual property rights through your use of Z.ai.

   d) Using content that explicitly belongs to a third party without permission.

   e) Removing or modifying copyright notices or other explicit attributions displayed by Z.ai.

   f) Using, modifying or otherwise altering AI-generated content in a way that gives rise to the infringement of any third party’s intellectual property rights.

   g) Combining generated content with external, content of third party, products, or services in a manner that leads to infringement.

   h) Using generated content for trademark or brand identity purposes without authorization.

   By using Z.ai, you acknowledge and agree to these terms and accept full responsibility for any violations.

5. You agree to defend and indemnify us from all losses, damages, and expenses (including reasonable attorneys' fees) arising from any third-party claims alleging that:

   a) your Prompts or any content you submit infringe third-party intellectual property rights;

   b) you have used the Services in combination with unauthorized third-party products, services, or content in a manner that results in infringement; or

   c) you continued to use the Services after receiving notice from us to cease such use due to potential infringement.

## VI. Copyright Protection Statement

1. **Reporting intellectual property or legal rights violations:** We have established measures and procedures to protect the legitimate rights of rights holders. If you believe that any part of Z.ai or its generated content infringes upon your intellectual property rights or other legal rights, you may submit an infringement notice in the form prescribed by and contain the information required under the applicable laws along with supporting evidence via email to [user\_feedback@z.ai](mailto:user_feedback@z.ai). We will take appropriate action, including removing or disabling relevant content that we believe violates these Terms. Additionally, we may terminate the visits of users who repeatedly infringe upon the rights of others, as deemed appropriate.

2. **Legal responsibility for false claims.** Please note that if a rights notice contains false statements, the submitter of the notice shall bear full legal responsibility for any resulting consequences, including but not limited to damages incurred by us and associated legal fees. If you are uncertain whether the content in question infringes your rights, we strongly recommend seeking professional legal advice before submitting a notice.

3. **Required documentation for a rights notice:** To ensure an efficient review of your claim, please provide the following information:

   a) Your detailed contact information, including your full name, a copy of your ID card or passport (if you are an individual), a copy of your business registration certificate (if you are an organization), your mailing address, telephone number, and email address.

   b) Proof of ownership of the allegedly infringed trademark, copyright, or any other legally enforceable rights.

   c) A clear and detailed description of the content that allegedly infringes upon your legal rights, along with screenshots or other evidence.

   d) A specific identification of the content on the alleged infringing webpage, indicating which part violates your rights.

   e) The following declaration in your notice:“I hereby confirm that the information in this notice is true and accurate. I am the legitimate rights holder of the reported content, or I have been authorized to act on its behalf. The reported content infringes my legal rights. If any information contained herein is false, I will bear all legal responsibilities resulting from such false statement.”

   f) If you are an individual, you must sign the notice. If you are an organization, the notice must bear your official seal.

4. **Processing and response time.** Upon verifying the authenticity of your claim, we will remove or disable access to the reported content. Please be aware that due to the time required for artificial intelligence model training, our removal or disabling measures may not take effect immediately. We appreciate your understanding.

5. We will only process English-language copyright notices that strictly comply with the applicable laws, and these Terms. Non-compliant notices will not be processed.

6. We will duly consider all compliant notices. You agree that before pursuing any legal action or remedy against us regarding infringing material, you must first provide us with a compliant notice and a reasonable opportunity to remove the material after our receipt. If we remove the material in response to your notice, you irrevocably waive any right of action against us under applicable law for any infringing material appearing on our Services before our removal under this section.

7. You expressly agree and acknowledge that we have no control over, and are not responsible or liable for, any content or material on third-party platforms or services.

## VII. Paid Accounts

1. When you use our fee-based Services, you may need to prepay for such Services. Your balance will be consumed if you use our paid services, and we reserve the right to suspend services if your balance are insufficient.

2. We have the right to determine the charging standards for fee-based services. The specific charging standards including but not limited to recharge methods, billing prices, consumption methods, and free quota for paid services are subject to the standards published on the relevant fee-based service pages. You may choose and pay for the corresponding Services according to your needs.

3. We may adjust our pricing model based on our business strategy. If such adjustments significantly affect your rights, we will promptly notify you through prominent methods such as in-site notifications, website announcements, or emails. If you continue to use the Services after the adjustments, it signifies your acceptance of the adjustments.

4. When recharging, please carefully confirm your order, account and information. You shall bear full responsibility for any losses incurred due to improper handling, inadequate comprehension, or insufficient understanding of the charge and billing procedures. You may cancel your paid subscription at any time. However, all payments are non-refundable unless otherwise required by applicable law. These Terms do not override any mandatory local laws that grant you specific cancellation rights.

5. You promise and guarantee the legality of the funds used to recharge. Otherwise, We have the right to cooperate with judicial or other government authorities' requirements to take corresponding measures on your account, including but not limited to locking, sealing, or restricting its use.

## VIII. Disclaimer of Warranties

TO THE FULLEST EXTENT PERMITTED BY APPLICABLE LAW:

1. OUR SERVICES ARE PROVIDED STRICTLY ON AN “AS IS” AND “AS AVAILABLE” BASIS;

2. WE, ALONG WITH OUR AFFILIATES AND LICENSORS, MAKE NO REPRESENTATIONS, WARRANTIES, GUARANTEES, OR UNDERTAKINGS (WHETHER EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE) REGARDING THE SERVICES, INCLUDING ANY OUTPUTS GENERATED. WE EXPRESSLY DISCLAIM ALL SUCH WARRANTIES, INCLUDING BUT NOT LIMITED TO THOSE RELATED TO MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, SATISFACTORY QUALITY, NON-INFRINGEMENT, RELIABILITY, IMPARTIALITY, CONSISTENCY, INTEROPERABILITY, CONFORMITY TO DESCRIPTION, SPECIFICATIONS, OR CRITERIA, COMPLIANCE WITH APPLICABLE LAWS AND THIRD-PARTY RIGHTS (INCLUDING BUT NOT LIMITED TO DATA PROTECTION, EXPORT CONTROL, AND INTELLECTUAL PROPERTY LAWS), AND UNINTERRUPTED ENJOYMENT. WE FURTHER DISCLAIM ANY GUARANTEES THAT THE SERVICES OR OUTPUTS WILL BE FREE FROM MALICIOUS, ERRONEOUS, MANIPULATED, INAPPROPRIATE, COMPROMISED, OR OTHERWISE HARMFUL CODE OR DATA, AS WELL AS ANY WARRANTIES ARISING FROM TRADE USAGE OR PRIOR DEALINGS;

3. WE DO NOT GUARANTEE THAT THE SERVICES WILL BE UNINTERRUPTED, ACCURATE, ERROR-FREE, OR SECURE, NOR DO WE WARRANT THAT ANY CONTENT WILL REMAIN INTACT, UNALTERED, OR FREE FROM LOSS; AND

4. YOU ACKNOWLEDGE AND AGREE THAT YOUR USE OF ANY OUTPUTS GENERATED BY OUR SERVICES IS ENTIRELY AT YOUR OWN RISK. YOU MUST NOT RELY ON SUCH OUTPUTS AS YOUR SOLE SOURCE OF TRUTH OR FACTUAL INFORMATION, NOR SHOULD THEY BE USED AS A SUBSTITUTE FOR PROFESSIONAL ADVICE.

## IX. Limitation of Liability

1. **Service limitations and no guarantees.** You acknowledge that despite our best efforts, due to inherent technical limitations in the technology supporting Z.ai’s services, we cannot guarantee the following:

   a) Z.ai may not meet your specific needs or expectations, and any products, services, or content obtained through we may not fulfill your intended purposes.

   b) We do not guarantee that the content generated by Z.ai will be 100% accurate, reliable, complete, functional, timely, secure, error-free, uninterrupted, or continuously stable. You should not rely on Z.ai’s Outputs as your sole source of factual information or as a substitute for professional advice.

   c) Although Z.ai applies AI-driven filtering mechanisms, some Outputs may still contain inaccuracies, inconsistencies, or material that could be deemed inappropriate. The Output of our Services does not represent the views or opinions of us, and you are responsible for assessing its accuracy and suitability, including conducting human review where necessary.

   d) Z.ai does not comprehend user inputs in the same way as a human, nor does it recognize underlying risks or ethical implications. The Outputs does not constitute advice or recommendations. You must make independent decisions based on your circumstances, and we are not liable for any risks or direct or indirect losses arising from reliance on Z.ai.

2. **Exemption from liability for service disruptions.** Z.ai fulfills its basic obligations as required by law; however, we are not liable for any failure, defect, delay, or modification in service performance due to:

   a) Force majeure events, including but not limited to natural disasters, strikes, riots, wars, governmental actions, and judicial or administrative orders.

   b) Failures in public utilities or third-party services, such as power outages or network disruptions.

   c) Situations where Z.ai has exercised reasonable management and due diligence, but disruptions occur due to routine or emergency system maintenance, system failures, software malfunctions, or cybersecurity threats.

3. **"As is" and "as available" disclaimer.** Unless otherwise explicitly stated in a written agreement, all information, content, materials, products, and services provided through Z.ai, including its generated content, are offered strictly on an “as is” and “as available” basis without any warranties of any kind.

4. **Exclusion of certain damages.** Neither we nor our affiliates, licensors, suppliers, or distributors shall be liable for any indirect, incidental, special, consequential, or exemplary damages. This includes, but is not limited to, damages for lost profits, loss of goodwill, loss of data, loss of use, or other intangible losses, even if we were advised of the possibility of such damages.

5. **Cap on liability.** Our total liability under these Terms shall not exceed the greater of:

   a) The amount you paid for the service that led to the claim within the 6 months preceding the liability, or

   b) One hundred dollars (\$100).

6. **Jurisdictional variations.** Some jurisdictions do not allow certain warranty disclaimers or limitations on liability. In such cases, the limitations set forth in this section will apply only to the maximum extent permitted by the laws of your country or state of residence, and you may have additional rights under local law.

7. **User commitments and liability.** You acknowledge the limitations of Z.ai and agree to comply with these Terms and all applicable laws and regulations. You agree not to engage in any actions that would result in legal claims, arbitration, compensation, or losses against Z.ai, its affiliates, or its employees in connection with your use of the platform. In the event of any such claims or disputes, you will be solely responsible for any resulting liabilities.

8. **Minor protection and liability disclaimer.** You acknowledge and agree that, despite our implementation of comprehensive measures to protect minors, we cannot fully guarantee that their use of this Service will always be controlled due to technical limitations, intentional misuse, or improper use by users. We bear no legal responsibility for any violations of laws, these Terms, or Platform Rules by minors, nor for any loss or harm resulting from a minor’s actions or a guardian’s failure to provide proper supervision. We strongly recommend that minors use this service under the guidance of their guardians, who are responsible for supervising their activities.

## X. Export Control and Sanctions

1. You understand that your use of Services, providing Inputs to and obtaining Outputs via Services, might be subject to all applicable laws and regulations of export controls and sanctions laws (collectively **“Export Control and Sanctions Laws”**) .You confirm that you and any person or entity you represent are not listed on any export control or economic sanction lists, nor are you or persons and entities you represent deemed to be included (including due to equity, control, or other affiliation relationships) in any such lists.

2. Such lists include but are not limited to the United Nations Security Council Sanctions List; the U.S. Specially Designated Nationals and Blocked Persons List (SDN List); the U.S. Entity List; the EU Consolidated List of Persons, Groups and Entities Subject to Financial Sanctions; China’s Unreliable Entities List, Export Control List, and Countermeasures List; and all other applicable control and sanctions lists issued by competent authorities with jurisdiction over you.

3. You confirm that you are not located in the following regions: Iran, North Korea, Cuba, Crimea, Donetsk, or Zaporizhzhia.

4. You understand that your use of the service, the inputs you provide to the service provider, and the Outputs you obtain through the service may be subject to all applicable export control and sanctions laws (collectively, "Export Control and Sanctions Laws").

5. You confirm that you are solely responsible for complying with all applicable Export Control and Sanctions Laws related to your use of the service and your end users’ access to and use of the services.

6. You represent and warrant that the service shall not be used for the benefit of, nor exported, re-exported, or transferred to:

   (a) any person or entity located in Iran, North Korea, Cuba, Crimea, Donetsk, or Zaporizhzhia;

   (b) any party subject to restricted lists under applicable Export Control and Sanctions Laws that would prohibit your use of the service.

7. Our service shall not be used for any prohibited end use, including military purposes or end uses involving human rights violations. You and your end users shall not provide inputs containing materials or information that require export control licenses.

## XI. Termination

1. **Termination by you.** You may terminate these Terms through any of the following methods:

   a) By closing your account;

   b) By ceasing to use the services and explicitly rejecting any proposed changes before they take effect;

   c) By expressing your intent to discontinue using our Services and meeting the platform's termination conditions.

   If you submit a request to delete your account, we will review and process your request. Once approved, your account will be deleted, and these Terms or any other terms as between you and us will be terminated.

2. **Termination by us.** We reserve the right to terminate your access to the services by deleting your account under the following circumstances:

   a) **Violation of the service agreement** – If you breach the terms of these Terms, we may terminate your access in accordance with its enforcement policies. You will be notified upon termination. If you attempt to re-register under your own or any third party’s name after termination, we reserves the right to terminate your access again.

   b) **Inability to contact you** – If the email or phone number provided during registration is no longer valid or cannot receive messages, and you fail to update your contact details within three business days after being notified, or if we cannot reach you through any means provided at registration.

   c) **False registration information** – If we determine that the key details in your registration information are false.

   d) **Rejection of updated terms** – If these Terms are terminated or updated, and you explicitly refuse to accept the new terms.

   e) **Illegal or unauthorized activity** – If you post illegal content, infringe on others' legal rights, or engage in other serious violations of laws or Platform Rules.

   f) **Repeated violations of Platform Rules** – If you repeatedly violate Z.ai’s platform policies, and the violations are deemed serious, leading to account suspension under Platform Rules.

   g) **Other grounds for termination** – If we determine that your account should be terminated for any other justified reason.

## XII. General

1. **Entire agreement.** These Terms represent the complete and final agreement between you and us regarding the Services and any related matters, replacing any previous written or verbal agreements on the same subject.

2. **Updates and changes.** We continually develop and improve our Services and may update these Terms or the Services themselves from time to time. Changes may occur due to:

   a) Legal or regulatory updates,

   b) Security or safety concerns,

   c) Events beyond our reasonable control,

   d) Routine service improvements, or

   e) Adaptation to new technologies.

   We will notify you of these changes by updating our website, sending an email, or using other appropriate communication methods with at least 30 days advance notice.

   Any changes take effect as soon as they are posted on our website. If you disagree with the updated Terms, you must stop using the Services immediately and notify us in writing to close your account. Continuing to use the Services after notification means you accept the updated Terms.

3. **Independent Relationship.** You and us are independent entities. These Terms do not create an agency, partnership, joint venture, employment, or franchise relationship between us.

4. **Service Discontinuation.** We reserve the right to discontinue, suspend or modify the Services at any time, and may provide you with prior notice where reasonably practicable.

5. **Severability.** If any part of these Terms is found to be illegal, invalid, or unenforceable under the law of any jurisdiction, the remaining provisions will still be valid and enforceable to the fullest extent possible.

6. **No waiver of rights.** If we delay or choose not to enforce any part of these Terms, it does not mean we waive our rights. We retain the right to act on any future or similar breaches. Our rights and remedies under these Terms are in addition to any other legal rights we have.

7. **No assignment without consent.** You cannot transfer, assign, sublicense, or subcontract your rights or obligations under these Terms to anyone else without our prior written consent.

8. **Breach and remedial actions.** If you engage in any of the following actions, it will be considered a breach of contract, and we reserve the right to take appropriate measures:

   a) We have the right to review your registration information and service usage. If we identify or suspect any issues, we may issue inquiries or require corrective action. If you fail to provide a reasonable response within a reasonable timeframe, we may take measures such as restricting, blocking, or deleting your registration information or service access.

   b) We have the right, without prior notice, to delete or impose restrictions on content you publish through our services if we reasonably believe that:

   i. It contains fraudulent, malicious, or harmful false information.

   ii. It disrupts the normal order of the platform.

   iii. It violates public interest or may seriously harm the legitimate rights of us or other users.

   If your breach causes us any loss, we have the right to seek compensation from you. Such losses include, but are not limited to, damage to our reputation, lost business opportunities, legal costs (such as litigation, arbitration, attorney fees, travel expenses, investigation costs, notarization, and appraisal fees), as well as any compensation, penalties, or fines paid to third parties, including those arising from administrative or criminal sanctions.

## XIII. Governing Law and Dispute Resolution

1. You understand and agree that when using Z.ai products in different regions, you must comply with all applicable laws, regulations, and policy requirements of those regions.

2. You may not resell, export, or transfer Z.ai products, or any products developed using Z.ai, to specific individuals or countries subject to regulatory restrictions. Furthermore, you may not directly or indirectly export Z.ai products or any technology containing Z.ai products without obtaining all necessary prior approvals from the relevant government authorities and the explicit prior consent of the Z.ai team.

3. These Terms, along with any disputes, claims, or legal actions arising from or relating to their subject matter, interpretation, enforcement, or formation—including those based on tort or other non-contractual grounds—shall be governed by the laws of Singapore.

4. Any disputes arising from or related to your use of the Services shall be resolved through mutual negotiation between you and us. You may contact us using the details provided at the bottom of these Terms.

5. If a resolution cannot be reached through negotiation, any dispute, controversy, difference or claim arising out of or relating to these Terms, including the existence, validity, interpretation, performance, breach or termination thereof, shall be referred to and finally resolved by arbitration administered by the Singapore International Arbitration Centre (“SIAC”) in accordance with the Arbitration Rules of the Singapore International Arbitration Centre (“SIAC Rules”) for the time being in force, which rules are deemed to be incorporated by reference in this clause. The seat of the arbitration shall be Singapore.

6. If any provision of these Terms is deemed null, invalid, or unenforceable, such provision shall be considered severable and shall not affect the validity or enforceability of the remaining provisions of these Terms.

## XIV. Contact

You have the right to submit complaints and reports by emailing [user\_feedback@z.ai](mailto:user_feedback@z.ai). We will assess your complaint or report and take appropriate action.

***

## Additional Terms for API Services

The use of API services (**“API Services”**) by business and developers on Z.ai platform is also governed by these Additional Terms for API Services (**“Additional Terms”**). Terms that are not defined in these Additional Terms have the meanings given in the Terms. For business and developers’ use of API Services, in case any inconsistency exists between the Terms and Additional Terms, the Additional Terms shall prevail.

### 1. Use API Services

a) We grant you a non-exclusive right to access and use the API Services during the valid term, which includes the right to use Z.ai’s API to integrate the Services into your applications or to develop downstream systems, applications or functions to your end users (**“End Users”**).

b) You are solely responsible for the systems, applications, or functions developed by using API Services. You must assume the corresponding legal responsibilities and obligations and establish agreements with your End Users regarding their rights and obligations. It is also your duty to regulate and manage the usage behavior of End Users.

c) You shall establish organizational and technical measures in accordance with the requirements provided under applicable laws and regulations, including but not limited to user management, data security, monitoring, warning, and emergency disposal, to ensure the integrity, confidentiality, and availability of your systems, networks, information, and data, preventing risks such as data security, public opinion, or misuse and abuse of products or services.

d) You must truthfully and accurately disclose the use of models from the Z.ai (and other models, if applicable) in accordance with relevant regulations and supervisory requirements. You may not make false or misleading disclosures, labels or displays regarding products and services that are not making use of models from the Z.ai

e) Given that API-returned information may not be flawless, you shall promptly take deletion, blocking, withdrawal and other disposal measures in accordance with the law for illegal and harmful information discovered or notified, to ensure lawful and compliant operation of applications.

f) You and End Users must not engage in any of the following improper uses of the Z.ai during your use of our API Services, including but not limited to:

i. Use the Services and End User Content (as defined in Section 3) in a way that violate applicable laws or our policies;

ii. Any attempt to extract data, parameters, or code from the Z.ai;

iii. Any use of the Z.ai for any services requiring specific qualifications, including but not limited to healthcare, finance, investments, insurance, credit, education, employment, housing, legal affairs, medical decisions, or other significant decisions;

iv. Any employment of the Z.ai for decision-making activities;

v. Any removal of copyright information from the Z.ai or its internal models;

vi. Any reverse engineering, decompiling, disassembling, or attempting to extract data from Z.ai algorithms, source code, mechanisms, or any related components;

vii. Any exploitation of our intellectual property through use, lease, lending, copying, modification, linking, reposting, compilation, publication, or establishment of mirror sites;

viii. Except for specific functionalities expressly authorized by us, any unauthorized reproduction, modification, augmentation, deletion, interception, or creation of derivative works from: (i) the Z.ai; (ii) data cached in terminal memory during platform operation; (iii) client-server interaction data during software execution; or (iv) essential system operation data, including but not limited to access through unauthorized plugins, extensions, or third-party tools/services;

ix. Any tampering with or falsification of software instructions or operational data to modify, enhance, reduce, or otherwise alter platform functionality or performance characteristics, or the commercialization or public dissemination of software/methods designed for such purposes, regardless of commercial intent;

x. Any access to or use of our software and services through unauthorized third-party software, plugins, add-ons, or systems, or any development/publication/distribution of such tools;

xi. Any direct or authorized interference with the Z.ai, including its components, modules, or data, whether conducted personally or through third-party software;

xii. Except for authorized integration with your specific business scenarios, any use of the Z.ai's models, prompts, or model-generated content for the development, training, labeling, fine-tuning, optimization, iteration, or similar activities related to external models is strictly prohibited. Furthermore, any utilization of platform models or exported data to develop, train, or enhance algorithms or models that compete with us;

xiii. Any generation of malicious code by using the Z.ai or its internal models;

xiv. Any other unauthorized actions or usage potentially detrimental to our interests.

### 2. Data Protection

a) We will collect and process data you provide as set forth in the Data Processing Addendum for API Services.

### 3. End User Content

a) You and End User might provide Input and generate Output (collectively **"End User Content"**) through the applications or services you provided/developed by using our API Services. We won't claim ownership over your Input and End User Content.

b) We will only use End User Content as necessary to provide you with the API Services, comply with applicable law, enforce our policies, and prevent abuse. \[We will not use End User Content to develop or improve Services, unless you explicitly agree to such use.]

c) You are responsible for all End User Content. You warrant that the Input has all rights, licenses and permission required to provide Input to the Services. You are responsible for managing the accuracy, safety, legality, and compliance of End User Content, establishing mechanisms including but not limited to content review, user management, data security, monitoring and early warning, and emergency response.

d) Due to the nature of the Services and artificial intelligence generally, Outputs may not be unique, and other users may receive similar content from the Services. Responses that are requested by and generated for other users are not considered your Outputs. You acknowledge that we may generate the same or similar content for others and that we reserve all rights to do so.

### 4. IP Rights

a) Except as expressly set forth herein, the Additional Terms does not grant: (a) us any IP Rights in End User Content; or (b) you any IP Rights in the Services. You only have a limited right to use the API Services, and no ownership rights are transferred to you or End Users under the Additional Terms.

b) You grant us only the limited rights that are reasonably necessary for us to deliver the API Services. This limited permission also extends to subcontractors or sub-processors.

c) If you provide feedback, you grant us the right to use and exploit feedback without restriction or compensation.

### 5. Termination and Effect of Termination

a) The Additional Terms will be terminated once the conditions set under Section XI of the Terms are met.

b) If the Terms and the Additional Terms terminate, the rights granted by us to you will cease immediately. Termination or expiration will not affect any rights or obligations, including the payment of amounts due, which have accrued under these Terms and separate order forms (if any) up to the date of termination.

c) Section 4 of the Additional Terms shall survive the termination of the Additional Terms.


# New Released
Source: https://docs.z.ai/release-notes/new-released

Follow along with updates across Z.AI’s models

## Models

**2026-08-26**

* Native visual capabilities enable the model to observe interfaces, rendering results, and interaction feedback—creating a closed loop across code, browsers, and GUIs.

* Efficient hybrid architecture: Combines linear and sparse attention with 320B total parameters and 18B activated, significantly reducing compute and KV-cache requirements.

* Beyond coding: Supports office document and financial research workflows, autonomously breaking down goals, using tools, and refining outputs. Learn more in our [documentation](/guides/vlm/glm-5.3-flash).\*

**2026-08-18**

* Stronger Coding Capabilities: GLM-5.3 delivers a significant improvement in coding capabilities, achieving a 50% gain over GLM-5.2 on Z.ai Code Bench and reaching state-of-the-art (SOTA) performance among open-source models on public benchmarks, including Terminal Bench 3.0.

* Emergent Cybersecurity Capabilities: GLM-5.3 matches Mythos 5 in white-box code review and vulnerability discovery. In collaboration with multiple cybersecurity teams, it has been tested on real-world targets and has identified a total of 2,436 vulnerabilities, including 1,097 medium- and high-severity vulnerabilities.Learn more in our [documentation](https://docs.z.ai/guides/llm/glm-5.3).\*

**2026-06-16**

* Supports 1M lossless context, significantly improving long-horizon task capabilities and reducing context drift and goal forgetting in complex tasks

* Achieves open-source SOTA performance on coding and long-horizon task benchmarks, delivering more stable results in complex system engineering and deep debugging

* Significantly improves the real-world developer experience, with more reliable project-level context handling, adherence to engineering standards, and multi-platform development. Learn more in our [documentation](/guides/llm/glm-5.2).\*

**2026-04-07**

* Designed for long-horizon tasks, GLM-5.1 can work independently for up to 8 hours in a single run, enabling a full loop from planning and execution to iterative refinement and final delivery.

* It demonstrates stronger engineering intelligence across autonomous planning, sustained execution, bug fixing, and strategy iteration, while achieving comprehensive capability alignment with Claude Opus 4.6. Built with multi-turn SFT, RL, and a process-quality evaluation framework, GLM-5.1 further improves stability, consistency, and tool use over extended tasks. Learn more in our [documentation](/guides/llm/glm-5.1).\*

**2026-02-12**

* Designed for complex system engineering and long-range Agent tasks, GLM-5 shifts the paradigm from coding to engineering, demonstrating strong deep-reasoning performance in backend architecture, complex algorithms, and stubborn bug fixing.

* It directly benchmarks against Claude Opus 4.5 in code-logic density and systems-engineering capability, and integrates DeepSeek Sparse Attention for higher token efficiency while preserving long-context quality.Learn more in our [documentation](/guides/llm/glm-5).\*

**2026-01-14**

* We’ve launched GLM-Image, a state-of-the-art image generation model built on a multimodal architecture and fully trained on domestic chips, combining autoregressive semantic understanding with diffusion-based decoding to deliver high-quality, controllable visual generation.
* The update significantly enhances performance in knowledge-intensive scenarios, with more stable and accurate text rendering inside images, making GLM-Image especially well suited for commercial design, educational illustrations, and content-rich visual applications.Learn more in our [documentation](/guides/image/glm-image).\*

**2025-12-22**

* We’ve released GLM-4.7, our foundation model with significant improvements in coding, reasoning, and agentic capabilities. It delivers more reliable code generation, stronger long-context understanding, and improved end-to-end task execution across real-world development workflows.
* The update brings open-source SOTA performance on major coding and reasoning benchmarks, enhanced agentic coding for goal-driven, multi-step tasks, and improved front-end and document generation quality. Learn more in our [documentation](/guides/llm/glm-4.7).\*

**2025-09-30**

* We’ve launched GLM-4.6, the flagship coding model, showcasing enhanced performance in both public benchmarks and real-world programming tasks, making it the leading coding model in China.
* The update also expands the context window to 200K, improving its ability to handle longer code and complex agent tasks. Learn more in our [documentation](/guides/llm/glm-4.6).\*


