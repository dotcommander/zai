> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Team Plan Benefits

> Learn about GLM Coding Team Plan usage quotas, benefits, and usage rules

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

For details on how API Keys are classified and billed, see the [API Key Management Guide](https://docs.z.ai/devpack/notice/api-key-management.md).

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

To protect subscriber rights, ensure system fairness, and maintain service stability, GLM Coding Plan must be used in [officially supported tools and products](https://docs.z.ai/devpack/tool/others.md#1-coding-agent-tool), and must comply with the [Subscriptions, Fees, and Payment](https://docs.z.ai/legal-agreement/subscription-terms.md) and related usage rules.

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

* [Quick Start](https://docs.z.ai/devpack/quick-start.md): Complete the basic integration process in just a few minutes and get started quickly
* [Tool Integration](https://docs.z.ai/devpack/tool/others.md): View the agent tools supported by the plan and their configuration methods, and choose the working environment that best suits your needs
* [How to Switch Models](https://docs.z.ai/devpack/latest-model.md): Make sure your coding tool is using your target model version


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.