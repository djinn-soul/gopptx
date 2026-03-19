# Contributing to gopptx

Thanks for contributing.

## Before You Start

- Check open issues and discussions before starting large work.
- For non-trivial changes, open a short proposal first.
- Keep changes modular and avoid large monolithic files.
- Do not start broad refactors without a scoped problem statement.

## Development Setup

## System Requirements

- Git.
- Go `1.25.8` (from `go.mod`).
- Python `3.10+` (from `pyproject.toml`).
- Docker (recommended for container-first workflows and docs build).
- PowerShell (required for `scripts/build_python.ps1` on Windows).

### Optional but useful

- `uv` for Python dependency/task workflows.
- `prek` for running project quality gates.

### Preferred: container workflow

Use the repository's container/dev-task workflow when available.

### Local workflow

1. Clone the repository.
2. Use repo-managed tooling (`Taskfile.yml`, scripts, and project config).
3. Build/run only through project commands where possible.
4. Run checks before opening a PR.

## Project Structure (high level)

- `pkg/` Go implementation and core editor modules.
- `internal/` internal Go packages and XML utilities.
- `python/` Python package and bridge bindings.
- `bindings/c/` C bridge interface docs and examples.
- `docs/` user and architecture documentation.
- `scripts/` build and smoke-test scripts.

## Coding Guidelines

- Make the smallest safe change that solves the problem.
- Preserve existing style and naming conventions.
- Prefer focused patch-style diffs over rewrites.
- Avoid silent fallbacks during development.
- Do not add empty `try/catch` blocks.
- Keep entry points stable and isolate new logic into focused modules.
- Keep code files modular and avoid oversized files.

## Testing and Validation

Run relevant checks for your change set:

- Go tests:
  `go test ./... -count=1`
- Go lint:
  `prek run golangci-lint --all-files`
- Architectural guardrails:
  `prek run architectural-guardrails --all-files`
- Python lint/type checks (if Python files changed):
  `prek run ruff --all-files`
  `prek run basedpyright --all-files`
- Docs build (if docs changed):
  `docker run --rm -v "${PWD}:/docs" squidfunk/mkdocs-material:9.7.5 build --strict`

If a check is environment-dependent, include details in the PR description.

## Branch and Commit Hygiene

- Use descriptive branch names (`feat/...`, `fix/...`, `docs/...`).
- Keep commits focused; avoid mixing unrelated changes.
- Write clear commit messages with intent and impact.

## Pull Request Checklist

- Explain what changed and why.
- Link related issue(s).
- Include verification commands and outcomes.
- Update docs when behavior or APIs changed.
- Keep PR scope focused; split unrelated work.
- Note any known limits, follow-ups, or deferred work.

## Documentation Expectations

- Update `README.md` when onboarding or usage changes.
- Update API/architecture docs when command contracts change.
- Prefer concise examples that can be copied and run.
- Keep examples aligned with current CLI/API behavior.

## Security

Do not introduce insecure patterns or commit secrets.
If you find a security issue, report it privately to maintainers and avoid public disclosure until triaged.

## Questions

If something is unclear, open a draft PR or issue with a minimal reproducer and expected outcome.
