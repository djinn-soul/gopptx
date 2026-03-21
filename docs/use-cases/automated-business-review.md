# Automated Business Review (QBR)

Generate repeatable QBR decks from business data instead of manual slide editing.

## Workflow

1. Pull quarter/department KPI data from your source system.
2. Build slides with `Presentation` helpers.
3. Export one deck per team/region.

## Python Example

```python
from gopptx import Presentation

qbr_data = {
    "quarter": "Q1 2026",
    "department": "Engineering",
    "kpis": [
        {"metric": "System Uptime", "value": "99.98%"},
        {"metric": "Feature Velocity", "value": "12 pts/week"},
        {"metric": "Tech Debt Reduction", "value": "15%"},
    ],
}

with Presentation.new() as pres:
    pres.add_slide(f"{qbr_data['quarter']} Review: {qbr_data['department']}")
    perf = pres.add_slide("Performance Summary")
    for i, kpi in enumerate(qbr_data["kpis"]):
        y = 2 + (i * 1.2)
        perf.add_shape(0, "rect", (1, y, 8, 1), text=f"{kpi['metric']}: {kpi['value']}")
    pres.add_slide("Strategic Outlook")
    pres.save("qbr_engineering_q1.pptx")
```

## Result Screenshot

![QBR result](../assets/images/showcase/qbr-result.png)