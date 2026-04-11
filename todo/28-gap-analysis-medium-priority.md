# 28 - Gap Analysis (Medium Priority)

Scope: prioritized medium-impact gaps for `gopptx` based on the provided `6.2 Medium Priority Gaps` list.

## Prioritized Gaps

| Gap | Impact | Notes | Tracking |
| --- | --- | --- | --- |
| Color manipulation API (`.lighter` / `.darker` / `.opacity`) | Medium | Useful for algorithmic theming | [19-missing-masters-layouts-themes-features.md](19-missing-masters-layouts-themes-features.md) |
| Embedded fonts | Medium | Consistent rendering in locked environments | [19-missing-masters-layouts-themes-features.md](19-missing-masters-layouts-themes-features.md) |
| Sub-bullets / hierarchical bullets | Medium | Complex list content | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| RTL text support | Medium | Arabic, Hebrew, Persian content | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Language tag on runs | Medium | Spell check and accessibility | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Code block with syntax highlighting | Low | Technical presentations | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| 3D model embedding | Low | Modern PowerPoint feature | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Digital signature creation | Medium | Compliance workflows | [25-missing-document-features.md](25-missing-document-features.md) |
| Image from URL | Medium | Simplifies image pipelines | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Image reflection / soft edges / blur | Low | Visual polish | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| Shape layout helpers (grid/center) | Low | Ergonomic positioning | [23-missing-shapes-advanced-operations-features.md](23-missing-shapes-advanced-operations-features.md) |
| Hidden slide flag | Low | Presenter-mode specific | [22-missing-slide-operations-features.md](22-missing-slide-operations-features.md) |

## Execution Order Suggestion

- [ ] Phase 1: RTL + language tags + embedded fonts + digital signatures.
- [ ] Phase 2: image URL/effects + shape layout helpers + hidden slide behavior.
- [ ] Phase 3: code-block helper and remaining polish.
