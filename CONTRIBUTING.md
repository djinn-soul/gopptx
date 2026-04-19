# Contributing to gopptx

Thanks for taking the time to contribute!
Before diving in, please read this guide — it keeps the review cycle short for everyone.

---

## Before You Start

- Check open issues and discussions before starting large work.
- For non-trivial changes, open a short proposal (issue or draft PR) first.
- Keep changes modular; avoid large monolithic files.
- Do not start broad refactors without a scoped problem statement.

---

## Development Setup

### System Requirements

| Tool | Version |
|---|---|
| Git | any recent |
| Go | `1.25.9` (from `go.mod`) |
| Python | `3.10+` (from `pyproject.toml`) |
| Docker | recommended — used for the container workflow and docs builds |
| PowerShell | required on Windows for `scripts/build_python.ps1` |

### Optional but useful

- `uv` — Python dependency and task management.
- `prek` — runs project quality gates (pre-commit wrapper).

### Preferred: container workflow

Use the repository's container/dev-task workflow when available.
Check `Taskfile.yml` for available targets:

```bash
task --list
```

### Local workflow

1. Clone the repository.
2. Use repo-managed tooling (`Taskfile.yml`, scripts, project config).
3. Build and run only through project commands where possible.
4. Run all checks before opening a PR (see below).

---

## Project Structure

```
pkg/          Go implementation and core editor modules
internal/     Internal Go packages and XML utilities
python/       Python package and bridge bindings
bindings/c/   C bridge interface and examples
docs/         User and architecture documentation
scripts/      Build and smoke-test scripts
```

---

## Coding Guidelines

- Make the **smallest safe change** that solves the problem.
- Preserve existing style and naming conventions.
- Prefer focused, patch-style diffs over full rewrites.
- Avoid silent fallbacks — let failures surface so they can be fixed.
- Do not add empty `try/catch` (or `recover`) blocks.
- Keep entry points stable; isolate new logic into focused modules.
- No source file should exceed ~300 lines of code.

---

## Running Checks

Run the full gate before opening a PR:

```bash
prek run --all-files
```

Or run individual checks as needed:

| Check | Command |
|---|---|
| Go tests | `go test ./... -count=1` |
| Go lint | `prek run golangci-lint --all-files` |
| Architectural guardrails | `prek run architectural-guardrails --all-files` |
| Python lint | `prek run ruff --all-files` |
| Python type check | `prek run basedpyright --all-files` |
| Docs build | `docker run --rm -v "${PWD}:/docs" squidfunk/mkdocs-material:9.7.5 build --strict` |

If a check is environment-dependent, include details in the PR description.

---

## Branch and Commit Hygiene

- Use descriptive branch names: `feat/...`, `fix/...`, `docs/...`.
- Keep commits focused; avoid mixing unrelated changes.
- Write clear commit messages — state the intent and impact, not just what changed.

---

## Pull Request Checklist

- [ ] Explain what changed and why.
- [ ] Link related issue(s).
- [ ] Include verification commands and their output.
- [ ] Update docs when behavior or APIs changed.
- [ ] Keep PR scope focused; split unrelated work into separate PRs.
- [ ] Note any known limits, follow-ups, or intentionally deferred work.

---

## Documentation Expectations

- Update `README.md` when onboarding steps or usage change.
- Update API/architecture docs when command contracts change.
- Prefer concise, runnable examples.
- Keep examples aligned with current CLI/API behavior.

---

## Security

Do not introduce insecure patterns or commit secrets.
If you find a security issue, report it **privately** to the maintainers before any public disclosure.

---

## Questions

If something is unclear, open a draft PR or issue with a minimal reproducer and the expected outcome.
