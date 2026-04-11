# 24 - Missing Text and Formatting Features

Scope: track text/formatting capabilities marked missing for `gopptx` in the provided feature audit screenshots.

## Missing in gopptx (Screenshot-Derived)

| Feature | Status | TODO |
| --- | --- | --- |
| Strikethrough | Missing | [ ] Add/verify strikethrough support in run-format APIs and bridge payloads. |
| Superscript / subscript | Missing | [ ] Add/verify superscript and subscript support in run-format APIs. |
| Font highlight color | Missing | [ ] Add/verify highlight-color controls in run formatting. |
| Text transparency | Missing | [ ] Add text transparency controls in run/paragraph styling payloads. |
| Text alignment (V / anchor) | Missing | [ ] Add vertical-anchor controls for text frames/paragraphs. |
| Line spacing | Missing | [ ] Add/verify line-spacing controls (points and/or multiplier) and serialization. |
| Character spacing | Missing | [ ] Add/verify character-spacing controls in run formatting. |
| Paragraph space before/after | Missing | [ ] Add/verify paragraph spacing-before/after controls. |
| Bullets (standard) | Missing | [ ] Add/verify standard bullet controls. |
| Bullets (numbered) | Missing | [ ] Add/verify numbered bullet/list controls. |
| Bullets (custom / Unicode) | Missing | [ ] Add custom-bullet character controls. |
| Hierarchical sub-bullets | Missing | [ ] Add/verify nested bullet level controls and rendering behavior. |
| RTL text | Missing | [ ] Add RTL text controls for runs/paragraphs/text frames. |
| Language tag | Missing | [ ] Add language-tag controls for text runs/paragraphs and ensure round-trip persistence. |
| Vertical text | Missing | [ ] Add vertical-text orientation modes and serialization. |
| Auto-fit (shrink/resize) | Missing | [ ] Add/verify text auto-fit controls and persistence. |
| Text wrap control | Missing | [ ] Add/verify text-wrap controls in text frame APIs. |
| Tab stops | Missing | [ ] Add/verify paragraph tab-stop support. |
| Code block with syntax highlighting | Missing | [ ] Add code-block formatting helper/API for syntax-highlighted text insertion. |
| Text shadow / glow / outline | Missing | [ ] Add text shadow/glow/outline effect controls. |

## Verification Follow-up

- [ ] Verify every listed gap against current runtime behavior and deduplicate with existing text trackers where already supported.
- [ ] Add regression tests for each truly-missing item once implemented.
