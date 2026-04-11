# 29 - Gap Analysis (Low Priority / Niche)

Scope: prioritized low-impact/niche gaps for `gopptx` based on the provided `6.3 Low Priority / Niche Gaps` list.

## Prioritized Gaps

| Gap | Impact | Notes | Tracking |
| --- | --- | --- | --- |
| YouTube embed | Low | M365-only feature | [18-missing-images-media-features.md](18-missing-images-media-features.md) |
| HTML table -> PPTX | Low | Browser-specific use case | [17-missing-table-features.md](17-missing-table-features.md) |
| Logarithmic chart scale | Low | Specialized chart need | [16-missing-chart-features.md](16-missing-chart-features.md) |
| Vertical text (7 modes) | Low | East Asian language use | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Web URL -> PPTX | Low | Niche import feature | [20-missing-advanced-unique-features.md](20-missing-advanced-unique-features.md) |
| PDF -> PPTX import | Low | Conversion workflow | [20-missing-advanced-unique-features.md](20-missing-advanced-unique-features.md) |
| HTML -> PPTX | Low | Conversion workflow | [20-missing-advanced-unique-features.md](20-missing-advanced-unique-features.md) |
| Tab stops in text | Low | Advanced typography | [24-missing-text-formatting-features.md](24-missing-text-formatting-features.md) |
| Data table inside chart | Low | Chart annotation | [16-missing-chart-features.md](16-missing-chart-features.md) |
| TableMergeMap advanced API | Low | Complex table merging | [17-missing-table-features.md](17-missing-table-features.md) |
| HandoutLayout type | Low | Handout-specific | [19-missing-masters-layouts-themes-features.md](19-missing-masters-layouts-themes-features.md) |
| IBM Carbon / Material Design color constants | Low | Design system shortcuts | [19-missing-masters-layouts-themes-features.md](19-missing-masters-layouts-themes-features.md) |

## Execution Order Suggestion

- [ ] Keep low-priority gaps queued behind high/medium priorities unless a customer use case depends on one.
- [ ] Batch related gaps by domain (charts/text/tables/themes/import) to reduce implementation churn.
