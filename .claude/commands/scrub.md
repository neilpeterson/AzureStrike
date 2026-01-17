---
description: Clean up project by removing unnecessary code and updating documentation
allowed-tools: Glob, Grep, Read, Edit, Write, Bash(git:*), Bash(go:*), Bash(ls:*), Task
---

# Project Scrub Command

Performs comprehensive project cleanup: removes unnecessary code, simplifies where possible, and ensures all documentation is current.

## Workflow:

### 1. **Pre-Scrub Assessment**
   - Read `docs/status.md` to understand current project state
   - Create a checklist of areas to review
   - Ask user if there are specific areas of concern

### 2. **Code Cleanup**

   **A. Find Unused Code:**
   - Run `go mod tidy` to clean dependencies
   - Check for unused imports (golangci-lint should catch these)
   - Look for:
     - Commented-out code blocks
     - Unused variables/functions (use staticcheck)
     - Dead code paths
     - TODO/FIXME comments that need action or removal

   **B. Simplification Opportunities:**
   - Look for overly complex functions (gocyclo findings)
   - Check for repeated code that could be refactored
   - Find magic numbers that should be constants
   - Identify functions that could be simplified

   **C. File Organization:**
   - Check for empty directories
   - Look for duplicate or redundant files
   - Verify examples are still relevant and working
   - Check for temporary or backup files (.bak, .tmp, etc.)

### 3. **Configuration & Build Files**

   **Review:**
   - `.goreleaser.yml` - Ensure config is correct and minimal
   - `.golangci.yml` - Check linter config is appropriate
   - `Makefile` - Verify all targets work and are needed
   - `.github/workflows/` - Ensure CI/CD configs are current
   - `go.mod` - Check for unused dependencies

   **Actions:**
   - Remove any unused configuration
   - Update outdated settings
   - Simplify where possible

### 4. **Documentation Review**

   **A. README.md:**
   - Verify installation instructions work
   - Check all code examples are current and correct
   - Ensure feature list matches implementation
   - Verify links are not broken
   - Update roadmap/status if needed
   - Check that Quick Start examples actually work

   **B. docs/status.md:**
   - Update "Current Focus" section
   - Move completed "In Progress" items to "Completed"
   - Add any new tasks to backlog
   - Update "Recent Changes" with latest work
   - Remove stale entries
   - Verify milestone completion status

   **C. Feature Documentation (docs/sparkline.md, docs/bar-chart.md, etc.):**
   - Verify API examples are correct
   - Check CLI usage examples work
   - Ensure code structure diagrams are current
   - Update any changed function signatures
   - Verify links to other docs

   **D. API Reference (docs/api-reference.md):**
   - Check all exported types are documented
   - Verify option names and behaviors
   - Ensure examples compile and work
   - Update if new features added

   **E. Project Instructions (.claude/CLAUDE.md):**
   - Verify workflow rules are still relevant
   - Check architecture diagram matches current structure
   - Update tech stack if dependencies changed
   - Review commands section for accuracy
   - Update chart types priority if changed

### 5. **Examples & Tests**

   **Examples:**
   - Run all example programs to verify they work
   - Remove any broken or outdated examples
   - Ensure examples demonstrate current best practices
   - Check that examples are referenced in docs

   **Tests:**
   - Run full test suite: `make test`
   - Check test coverage: `make cover`
   - Look for obsolete tests
   - Ensure test names are descriptive

### 6. **Dependencies**

   **Review:**
   - Run `go mod tidy`
   - Check `go.mod` for unused dependencies
   - Look for outdated dependencies (consider `go list -m -u all`)
   - Verify all dependencies are necessary

### 7. **Repository Hygiene**

   **Check for:**
   - Uncommitted changes: `git status`
   - Untracked files that should be committed or gitignored
   - Large files that shouldn't be in repo
   - Sensitive data (API keys, credentials)
   - Branches that can be deleted

### 8. **Create Summary Report**

   **Document all changes made:**
   - What was removed and why
   - What was simplified
   - Documentation updates made
   - Any issues found but not fixed (and why)
   - Recommendations for future cleanup

### 9. **Update Status**

   - Update `docs/status.md` with scrub completion
   - Add entry to "Recent Changes"
   - Note any technical debt identified

### 10. **Commit Changes**

   - Stage all changes
   - Create detailed commit message explaining cleanup
   - Suggest whether changes should be in one commit or split
   - Ask user if they want to create PR or push directly

## Important Notes:

- **Be Conservative**: Only remove code you're certain is unused
- **Preserve History**: Don't remove comments that explain "why" decisions were made
- **Test After Changes**: Run tests after any code modifications
- **Ask Before Big Changes**: If unsure about removing something, ask the user
- **Document Everything**: Keep track of what was removed for the summary

## Tools to Use:

- `go mod tidy` - Clean dependencies
- `make lint` - Find code issues
- `make test` - Verify nothing broke
- `git status` - Check for untracked files
- Task tool with Explore agent - For comprehensive code exploration
- Grep tool - Search for patterns (TODO, FIXME, etc.)
- Read tool - Review documentation files

## Success Criteria:

- ✅ All tests still pass
- ✅ All documentation is current
- ✅ No unnecessary files remain
- ✅ Examples all work
- ✅ Dependencies are minimal
- ✅ Code is as simple as possible
- ✅ Status.md accurately reflects project state

## Example Usage:

```bash
/scrub
```

The command will systematically go through all areas and report findings.