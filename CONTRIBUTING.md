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
| Go | `1.25.12` (from `go.mod`) |
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
- No source file should exceed ~400 lines of code.
- Never silence a warning, lint or type error to go green — fix the cause or migrate the callers.

---

## Things That Will Bite You

### Go declarations are the source of truth

Parts of the Python surface are **generated**, not written:

| Generated | From |
|---|---|
| `python/gopptx/_ops_constants.py`, `ops.pyi` | `pkg/pptx/editor/opspec.go` |
| The `ChartType` enum | `XLChartType` in `pkg/pptx/enums` |
| The `ShapeType` enum | The shape preset declarations |
| The slide builder surface | `elements.SlideContent` |
| `python/gopptx/gopptx.h` | `bindings/c/bridge.go` (cgo) |

Edit the Go source and run `task generate`. Never hand-edit a generated file — `task
check:generated` fails the build when they drift.

`python/gopptx/__init__.pyi` is not generated but is checked: a stub-parity test fails if its
`__all__` and the runtime `__all__` disagree. Update both together.

### Stage the cgo header with the bridge

Any change to `bindings/c/bridge.go` regenerates `python/gopptx/gopptx.h`. Stage both together
or the pre-commit hooks fail with a confusing error.

### Rebuild the shared library after Go changes

The Python package binds a compiled library. Go edits have no effect in Python until:

```bash
task build:go
```

### Geometry is EMU

914 400 per inch. A shape at `(40, 120, 600, 220)` is smaller than a pixel — the file saves,
PowerPoint opens it, and nothing is visible. Use `Inches()` / `Point()` / `Emu()` in examples,
tests and docs. This is the single most common defect in contributed samples.

### Verify by rendering, not only by testing

A deck can pass every structural test and still be one PowerPoint refuses to open, or one that
opens blank. For anything touching output bytes, open the result in real PowerPoint. For PDF
changes, rasterise with `pypdfium2` and compare — LibreOffice drops clips and substitutes
fonts, so it disagrees with itself.

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
- **Run every code sample you add.** A sample that does not compile or raises on the first call
  is worse than no sample — most of the documentation defects found in this repository were
  snippets nobody executed.
- Wrap coordinates in `Inches()` / `Point()` / `Emu()` in every sample.
- When you close a gap listed in `PPT_RS_PARITY_2026-08-04.md`, update that document **and**
  `docs/reference/feature-matrix.md`. When you find a new gap, record it there rather than
  leaving the comparison flattering.
- `docs/reference/python-presentation-api.md` is written from introspection of the installed
  `Presentation` class. Regenerate it rather than hand-editing entries.

### Building the docs

```bash
task docs:serve      # http://localhost:8000
task docs:build      # strict; fails on broken links
```

---

## Security

Do not introduce insecure patterns or commit secrets.
If you find a security issue, report it **privately** to the maintainers before any public disclosure.

---

## Questions

If something is unclear, open a draft PR or issue with a minimal reproducer and the expected outcome.
