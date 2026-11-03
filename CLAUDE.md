# Working rules for Claude

These rules apply to every task in this repo and override default behavior.

## 1. Plan first, never execute without approval

- For every request, produce a **change plan** first: which files change, what changes in each, why, the order, and how to verify.
- "Execute" means anything that changes state: creating, editing, or deleting files; installing packages; running builds, formatters, or migrations; any git write (commit, branch, push, stash, reset) — see rule 4: git writes are never Claude's to run.
- Read-only work is allowed while planning: reading files, searching, `git status` / `git diff` / `git log`.

## 2. Ask questions freely

- Ask as many clarifying questions as needed before or during planning. Prefer asking over guessing on any decision about code, structure, naming, or behavior.

## 3. The approval phrase: "execute order 66"

- Do not execute anything until the user's message contains the exact phrase **"execute order 66"**.
- Approval covers only the plan it was given for. A new or changed plan needs the phrase again.
- At the end of every plan, remind the user: _"Say **execute order 66** when you've reviewed this and want me to apply it."_
- If the user asks for a change without the phrase, reply with the plan and the reminder — don't apply it.
- Why: the user wants to review every code change and decision before it happens.

## 4. Never run git write commands — the user does all git themselves

- Never commit, push, branch, checkout, merge, rebase, stash, reset, tag, `git add`, or `git rm` (including `--cached`), even after "execute order 66".
- When a plan needs git, list the exact commands for the user to run themselves instead.
- Read-only git is still fine: `git status`, `git diff`, `git log`, `git show`, `git blame`.
