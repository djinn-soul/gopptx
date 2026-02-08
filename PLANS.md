# Plan: Project Renaming (gopptx)

The project has been renamed from `goppt` to `gopptx`.

## Purpose

The name `goppt` was already in use on GitHub. Renaming to `gopptx` ensures a unique package name and repository.

## Progress

- [x] Rename module in `go.mod` <!-- id: 0 -->
- [x] Update all internal imports in `.go` files <!-- id: 1 -->
- [x] Update project name in `roadmap.md` <!-- id: 2 -->
- [x] Update project name in `cmd/pptcli/README.md` <!-- id: 3 -->
- [x] Update project name in `tasks/README.md` <!-- id: 4 -->
- [x] Update project name in GitHub Action workflows (`.github/workflows/*.yml`) <!-- id: 5 -->
- [x] Update project name in `.github/commands/*.toml` <!-- id: 6 -->
- [x] Research root causes of "needs repair" errors (2026-02-08)
- [x] Analyze existing `validate` command and XML generation logic (2026-02-08)
- [x] Update any other references found during scouting <!-- id: 7 -->

## Surprises & Discoveries

- The `gopptx` library uses manual string building for XML, which requires careful tag ordering and namespace management.
- Recent fixes in `CONTINUITY.md` highlight real-world repair issues (invalid attributes, incomplete themes).
- The `validate` command currently only checks XML syntax, not OOXML schema compliance.

## Decision Log

- **D001 ACTIVE**: Focus on providing a diagnostic guide and highlighting existing library safeguards instead of a code-heavy "fix" since the user's question is conceptual.

## Decisions

- **Module Name**: `github.com/djinn09/gopptx`
- **Internal Imports**: Replace all occurrences of `github.com/djinn09/goppt` with `github.com/djinn09/gopptx`.
- **Consistency**: All documentation and configurations now reflect the new name.

## Validation

- `go build ./...` passes.
- `go test ./...` passes.
- Grep for `goppt` confirmed no stray references remain (except in this plan's history).