# HELMOR Skills

HELMOR Agent ships with 14 built-in skills. They help AI coding agents choose the right behavior before spending tokens, editing code, running checks, or closing out a task.

## Lifecycle Map

| Lifecycle | Skills | Outcome |
|---|---|---|
| Plan | Product Planning, Architecture, API Contracts | Turn intent into scoped implementation work. |
| Build | Frontend, Backend, Data, Infrastructure, UI Design | Keep implementation aligned with the project shape. |
| Verify | Testing, Security, Launch Readiness | Make checks, review, and release blockers explicit. |
| Remember | Project Memory, Token Reduction, Docs & Handoff | Preserve context and reduce repeated discovery. |

## Skill Reference

| Skill | Focus |
|---|---|
| Product Planning | Scope, acceptance criteria, task breakdown, and user-facing intent. |
| Architecture | System boundaries, module ownership, dependencies, and implementation shape. |
| API Contracts | Routes, schemas, clients, mocks, and integration expectations. |
| Frontend | Components, routes, state, forms, responsiveness, and user workflows. |
| Backend | Services, jobs, integrations, validation, and operational behavior. |
| Data | Models, migrations, seeds, fixtures, and data integrity expectations. |
| Infrastructure | CI, package runners, deployment checks, environments, and config hygiene. |
| UI Design | Design drift, visual polish, accessibility basics, and product-specific interface rules. |
| Testing | Required checks, failure evidence, closeout gates, and regression coverage. |
| Security | Secrets, destructive commands, unsafe deploys, and sensitive code paths. |
| Launch Readiness | Release blockers, final review, production approval, and ship/no-ship signals. |
| Project Memory | Context cards, handoffs, prior decisions, and local task state. |
| Token Reduction | RTK-first output, SQZ fallback, graph-first discovery, and fewer repeated scans. |
| Docs & Handoff | Truth files, READMEs, task summaries, and next-agent continuity. |

## How Skills Are Used

Skills are declared in `.helmor/project.json` and surfaced through session context, hooks, and closeout guidance. The watcher routes prompts and tool-use events toward relevant skills so agents do not need to rediscover the same project rules repeatedly.
