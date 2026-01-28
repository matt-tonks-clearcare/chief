---
description: Learn how Chief works as an autonomous coding agent, transforming your requirements into working code through an automated execution loop.
---

# How Chief Works

Chief is an autonomous coding agent that transforms your requirements into working code, without constant back-and-forth prompting.

::: tip Background
For the motivation behind Chief and a deeper exploration of autonomous coding agents, read the blog post: [Introducing Chief: Autonomous PRD Agent](https://minicodemonkey.com/blog/2025/chief)
:::

## The Core Concept

Traditional AI coding assistants hit a wall: the context window. As your conversation grows, the AI loses track of earlier details, makes contradictory decisions, or simply runs out of space. Long coding sessions become unwieldy.

Chief takes a different approach using a [Ralph Wiggum loop](https://ghuntley.com/ralph/): **each iteration starts fresh, but nothing is forgotten.**

You describe what you want to build as a series of user stories. Chief works through them one at a time, spawning a new Claude session for each. Between iterations, Chief persists state to a `progress.md` file: what was built, which files changed, patterns discovered, and context for future work. The next iteration loads this history, giving Claude everything it needs without the baggage of a bloated conversation.

Running `chief` opens a TUI dashboard where you can review your project, then press `s` to start the loop.

## The Execution Loop

Chief works through your stories methodically. Each iteration focuses on a single story:

```
                ┌───────────────────────────────────────┐
                │                                       │
                ▼                                       │
        ┌──────────────┐                                │
        │  Pick Story  │                                │
        │  (next todo) │                                │
        └──────┬───────┘                                │
               │                                        │
               ▼                                        │
        ┌──────────────┐                                │
        │ Invoke Claude│                                │
        │  with prompt │                                │
        └──────┬───────┘                                │
               │                                        │
               ▼                                        │
        ┌──────────────┐                                │
        │    Claude    │                                │
        │ codes & tests│                                │
        └──────┬───────┘                                │
               │                                        │
               ▼                                        │
        ┌──────────────┐                                │
        │    Commit    │                                │
        │   changes    │                                │
        └──────┬───────┘                                │
               │                                        │
               ▼                                        │
        ┌──────────────┐           more stories         │
        │ Mark Complete├────────────────────────────────┘
        └──────┬───────┘
               │ all done
               ▼
           ✓ Finished
```

Here's what happens in each step:

1. **Pick Story**: Chief finds the highest-priority incomplete story
2. **Invoke Claude**: Constructs a prompt with the story details and project context, then spawns Claude Code
3. **Claude Codes**: Claude reads files, writes code, runs tests, and fixes issues until the story is complete
4. **Commit**: Claude commits the changes with a message like `feat: [US-001] - Feature Title`
5. **Mark Complete**: Chief updates the project state and records progress
6. **Repeat**: If more stories remain, the loop continues

This isolation is intentional. If something breaks, you know exactly which story caused it. Each commit represents one complete feature.

## Conventional Commits

Every completed story results in a well-formed commit:

```
feat: [US-003] - Add user authentication

- Implemented login/logout endpoints
- Added JWT token validation
- Created auth middleware
```

Your git history becomes a timeline of features, matching 1:1 with your stories.

## Progress Tracking

The `progress.md` file is what makes fresh context windows possible. After every iteration, Claude appends:

- What was implemented
- Which files changed
- Learnings for future iterations (patterns discovered, gotchas, context)

When the next iteration starts, Claude reads this file and immediately understands the project's history, without needing thousands of tokens of prior conversation. This gives you the benefits of long-running context (consistency, institutional memory) without the downsides (context overflow, degraded performance).

## Staying in Control

Autonomous doesn't mean unattended. The TUI lets you:

- **Start / Pause / Stop**: Press `s` to start, `p` to pause after the current story, `x` to stop immediately
- **Switch projects**: Press `n` to cycle through projects, or `1-9` to jump directly
- **Resume anytime**: Walk away, come back, press `s`. Chief picks up where you left off

## Further Reading

- [The Ralph Loop](/concepts/ralph-loop): Deep dive into the execution loop mechanics
- [PRD Format](/concepts/prd-format): How to structure your project with effective user stories
- [The .chief Directory](/concepts/chief-directory): Understanding where state is stored
