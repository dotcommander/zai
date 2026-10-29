> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Memory-mechanism

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


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.