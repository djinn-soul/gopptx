"""Tests for chart state traversal and axis controls."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation

if TYPE_CHECKING:
    from pathlib import Path

_MIN_SERIES_VALUES = 2


def test_chart_series_and_axis_state(tmp_path: Path) -> None:
    """Chart state exposes axes and series values from bridge parsing."""
    out_path = tmp_path / "chart_state_parity.pptx"

    with Presentation.new("Chart State") as prs:
        slide = prs.add_slide("Chart")
        _ = slide.add_chart(
            "bar",
            ["Q1", "Q2"],
            [1.0, 2.0],
            title="Revenue",
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        chart = slide.charts[0]
        assert chart.category_axis.present is True  # noqa: S101
        assert chart.value_axis.present is True  # noqa: S101
        assert len(chart.series) >= 1  # noqa: S101
        assert len(chart.series[0].values) >= _MIN_SERIES_VALUES  # noqa: S101
        prs.save(str(out_path))


def test_chart_axis_tick_label_position_update(tmp_path: Path) -> None:
    """Chart axis tick-label updates are persisted into chart XML."""
    out_path = tmp_path / "chart_axis_tick_label_pos.pptx"

    with Presentation.new("Chart Axis Tick Labels") as prs:
        slide = prs.add_slide("Chart")
        _ = slide.add_chart(
            "bar",
            ["A", "B"],
            [1.0, 2.0],
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        chart = slide.charts[0]
        chart.category_axis.tick_label_position = "low"
        chart.value_axis.tick_label_position = "high"
        prs.save(str(out_path))

    with zipfile.ZipFile(out_path) as zf:
        chart_xml = zf.read("ppt/charts/chart1.xml").decode("utf-8")
    assert '<c:tickLblPos val="low"/>' in chart_xml  # noqa: S101
    assert '<c:tickLblPos val="high"/>' in chart_xml  # noqa: S101
