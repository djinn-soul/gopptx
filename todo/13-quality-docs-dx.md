# 13 - Quality, Docs, and DX

Scope: track repo-quality gaps for discoverability and developer experience (README, docs depth, examples, IDE hints, and error clarity).

## Investigation Snapshot (2026-04-11)

| Area | Score | Current finding |
| --- | --- | --- |
| README | ⭐⭐⭐⭐ | `README.md` is strong and task-oriented, with good quickstart and architecture references. |
| API Docs | ⭐⭐⭐ | Python has broad docstring coverage overall (`functions: 855/1149`, `classes: 208/237`), but public API narrative depth and consistency are uneven. |
| Examples | ⭐⭐⭐⭐⭐ | `examples/` is extensive and diverse (many feature-focused directories plus API examples). |
| IDE Integration | ⭐⭐⭐ | `python/gopptx/api.pyi` exists and typing is substantial, but discoverability and richer type-guided docs can still improve. |
| Error Messages | ⭐⭐ | Bridge failures can be opaque: runtime errors often surface only message/code and omit context like op name, request_id, and payload hints. |

## TODO

- [ ] Expand Python API docs with richer public-method examples and parameter/return semantics.
- [ ] Add a docs-quality checklist for new API methods (docstring quality + reference-page update requirement).
- [ ] Improve IDE affordances: add usage snippets and typed examples near key `.pyi` surfaces.
- [ ] Improve bridge error reporting to include `op`, `request_id`, and actionable context in `GopptxError`.
- [ ] Add error-message regression tests for common bridge failure scenarios.

## Evidence Pointers

- `README.md`
- `python/gopptx/api.pyi`
- `python/gopptx/presentation/runtime.py`
- `python/gopptx/api_errors.py`
- `examples/`
