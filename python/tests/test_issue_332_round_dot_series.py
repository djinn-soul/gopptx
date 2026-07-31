"""Regression coverage for chart-series round-dot density (issue #332)."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import MSO_LINE, Presentation
from gopptx.presentation.charts.chart_types import ChartType

if TYPE_CHECKING:
    from pathlib import Path


def test_series_round_dot_uses_system_dot_token(tmp_path: Path) -> None:
    """MSO_LINE.ROUND_DOT serializes as the denser DrawingML sysDot preset."""
    output = tmp_path / "round-dot-series.pptx"
    with Presentation.new("Issue 332") as presentation:
        presentation.slides[0].add_chart(
            ChartType.LINE_MARKERS,
            ["Q1", "Q2", "Q3", "Q4"],
            [{"name": "Trend", "values": [2.0, 5.0, 3.0, 7.0]}],
            bounds=(500000, 1200000, 7000000, 3500000),
        )
        chart = presentation.slides[0].charts[0]
        chart.series[0].format.line.dash_style = MSO_LINE.ROUND_DOT
        presentation.save(output)

    with zipfile.ZipFile(output) as package:
        chart_xml = package.read("ppt/charts/chart1.xml").decode()
    assert '<a:prstDash val="sysDot"/>' in chart_xml
    assert '<a:prstDash val="dot"/>' not in chart_xml

    reopened = Presentation()
    reopened.open(output)
    try:
        assert reopened.slides[0].charts[0].series[0].format.line.dash_style == "sysDot"
    finally:
        reopened.close()
