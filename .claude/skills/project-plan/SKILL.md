---
name: project-plan
description: How to organize, structure, and keep a project plan up to date.
---

When asked about project plan, planning, spec, specification, project status, or task management, use this skill to maintain structured project documentation.

## When to Use This Skill

- User asks to "create a plan", "plan this project", or "break down this task"
- User asks to "create a spec" or create a "specification document"
- User asks about "project status", "what's left to do", or "next steps"
- User wants to track progress on a multi-step implementation
- Starting a new feature or project that requires multiple steps

## Required Files

Look for these files in the docs/project-plan directory. If they don't exist, create them:

1. **`SPEC.md`** - Project specification document
2. **`PLAN.md`** - Project plan with trackable tasks

## Process

1. **Search** for `SPEC.md` and `PLAN.md` in the docs/project-plan directory
2. **Read** the contents to understand current project state
3. **Create** the files if they don't exist
4. **Update** the files when:
   - New requirements are identified
   - Tasks are completed (check them off)
   - The project scope changes
   - New phases or features are added

## File Formats

### SPEC.md Structure

```markdown
# Project Specification

## Overview
Brief description of the project purpose and goals.

## Scope
What is included and explicitly excluded from this project.

## Requirements

### Functional Requirements
Features and capabilities the system must provide.

### Non-Functional Requirements
Performance, security, scalability, reliability, and accessibility requirements.

### Constraints and Assumptions
Technical limitations and assumptions made during design.

## Architecture
High-level design decisions and component structure.

## Schema
Data models, database schema, or API contracts.

## Implementation Details
Specifics about technologies, frameworks, and patterns to be used.

## Dependencies
External libraries, services, or resources needed.

## Glossary
Definitions of domain-specific terms (optional, for complex projects).
```

### PLAN.md Structure

This should be a simple document that outlines features, tasks, TODOs, and blockers with checkboxes for tracking progress. Do no inlcude dates, timeline, specifiction, or executions commands. 

```markdown
# Project Plan

## Status
Current project phase and overall progress summary.

## Phases

### Feature 1: [Phase Name]

- [x] Completed task
- [ ] Pending task
- [ ] Another pending task

### Feature 2: [Phase Name]

- [ ] Future task
```

## Best Practices

1. **Keep tasks granular** - Each task should be completable in a single session
2. **Update immediately** - Mark tasks complete as soon as they're done
3. **Add context** - Include file paths, method names, or specific details in task descriptions
4. **Track blockers** - Note any dependencies or blockers for pending tasks
5. **Version milestones** - Use phases to group related work into logical milestones