---
name: pull-request
description: Steps to complete prior to creating a pull request.
---

## Pre-PR Workflow

### Step 1: Run Project Scrub

Before creating any pull request, run the `/scrub` command to:
- Remove unnecessary code
- Update documentation
- Clean up dependencies
- Ensure all tests pass

```
/scrub
```

Wait for scrub to complete before proceeding.

### Step 2: Create Pull Request

After scrub completes successfully:

1. Stage and commit any scrub changes
2. Push branch to remote
3. Create PR using `gh pr create`

### Step 3: Monitor CI Status

After PR is created, monitor the CI checks:

```bash
# Watch PR checks status
gh pr checks <PR_NUMBER> --watch
```

Or poll periodically:

```bash
# Get current check status
gh pr checks <PR_NUMBER>
```

### Step 4: Handle Results

**If all checks pass:**
- Prompt user: "All CI checks passed. Would you like to merge the PR?"
- If yes, merge using: `gh pr merge <PR_NUMBER> --squash --delete-branch`

**If any checks fail:**
1. Get the failed check details:
   ```bash
   gh pr checks <PR_NUMBER>
   ```

2. View the workflow run logs:
   ```bash
   gh run view <RUN_ID> --log-failed
   ```

3. Automatically fix the identified issues:
   - For lint errors: Fix the code issues reported
   - For test failures: Fix the failing tests or underlying code
   - For build errors: Resolve compilation issues

4. Commit fixes and push:
   ```bash
   git add -A
   git commit -m "Fix CI failures"
   git push
   ```

5. Return to Step 3 to monitor the new CI run

## Commands Reference

| Command | Purpose |
|---------|---------|
| `gh pr create` | Create pull request |
| `gh pr checks <PR>` | View CI check status |
| `gh pr checks <PR> --watch` | Watch checks until complete |
| `gh run view <ID> --log-failed` | View failed workflow logs |
| `gh pr merge <PR> --squash` | Squash merge the PR |
| `gh pr merge <PR> --delete-branch` | Delete branch after merge |

## Example Flow

```
User: /pull-request

1. Run /scrub
2. Commit scrub changes if any
3. gh pr create --title "..." --body "..."
4. gh pr checks 123 --watch
5. [If pass] "All checks passed. Merge PR?" → gh pr merge 123 --squash --delete-branch
   [If fail] Fix issues → git commit → git push → repeat step 4
```