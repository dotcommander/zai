> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# API Key Management Guide

> Starting October 16, 2026, API keys for the standard model API and GLM Coding Plan are managed separately. Learn which key to use, how existing keys are classified, and how billing works for each key type.

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

      <td rowSpan="3">
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

Check the key configured in your tool. **Standard keys do not use Coding Plan quota. An existing personal Compatible key also uses standard API billing when calling the standard API endpoint.** To use plan quota, get a key from the corresponding Coding Plan entry point and replace the key in your tool. For setup instructions, see the [Quick Start](https://docs.z.ai/devpack/quick-start.md).

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


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.