# Agent Adapters

HELMOR v1 supports four first-class agent targets:

- Codex
- Claude Code
- Cursor
- Windsurf

Codex and Claude Code can consume `helmor hook --event <EventName>` command
hooks. Cursor and Windsurf receive generated rule files that instruct agents to
use HELMOR status, doctor, handoff, and token-reduction commands.

Adapter event names:

- `SessionStart`
- `UserPromptSubmit`
- `PreToolUse`
- `PostToolUse`
- `Stop`
- `PreCompact`
- `SessionEnd`

Unsupported agents should use the same CLI commands and can implement adapters
around the hook interface.

