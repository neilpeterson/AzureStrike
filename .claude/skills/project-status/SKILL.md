---
name: project-status
description: Provides project status by reading the PLAN.md file. Use this when asked about project progress or status.
---

When asked about project status or progress, look for a `PLAN.md` file in the current project root.

## Process

1. Search for a `PLAN.md` file in the workspace root
2. Read the contents of the file
3. Summarize the project status including:
   - Current phase or milestone
   - Completed items
   - Remaining items
   - **Percent done** (calculate based on completed vs total tasks if not explicitly stated)

## Response Format

Provide a concise status update that includes:

- **Project**: Name of the project
- **Status**: Current state (e.g., In Progress, On Track, Blocked)
- **Percent Complete**: X%
- **Summary**: Brief description of what's done and what remains

If no `PLAN.md` file is found, inform the user that no status file exists and suggest creating one.