# Project Profile

HELMOR stores project rules in `.helmor/project.json`.

Important fields:

- `mode`: `observe`, `guard`, or `strict`
- `packageRunner`: detected package runner
- `framework`: detected app/framework type
- `truthFiles`: project authority docs
- `checks`: commands agents should run before closeout
- `agents`: enabled adapters
- `tools`: token-saving tool preferences
- `policies`: detector and enforcement settings
- `skills`: built-in HELMOR skill modules

The profile is intentionally local and plain JSON so developers can review and
edit it without a service account or cloud dependency.

