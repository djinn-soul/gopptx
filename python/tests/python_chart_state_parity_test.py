"""Tests for chart state traversal and axis controls."""

from __future__ import annotations

import zipfile
from math import isclose
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
        assert chart.category_axis.present is True
        assert chart.value_axis.present is True
        assert len(chart.series) >= 1
        assert len(chart.series[0].values) >= _MIN_SERIES_VALUES
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
    assert '<c:tickLblPos val="low"/>' in chart_xml
    assert '<c:tickLblPos val="high"/>' in chart_xml


def test_chart_axis_title_scale_and_number_format_update(tmp_path: Path) -> None:
    """Chart axes expose python-pptx-style title, scale, unit, and number format."""
    out_path = tmp_path / "chart_axis_details.pptx"

    with Presentation.new("Chart Axis Details") as prs:
        slide = prs.add_slide("Chart")
        _ = slide.add_chart(
            "bar",
            ["A", "B"],
            [1.0, 2.0],
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        chart = slide.charts[0]
        chart.category_axis.title = "Quarter"
        value_axis = chart.value_axis
        value_axis.title = "Revenue"
        value_axis.minimum_scale = 0.0
        value_axis.maximum_scale = 200.0
        value_axis.major_unit = 25.0
        value_axis.minor_unit = 5.0
        value_axis.tick_label_number_format = "$#,##0.00"
        value_axis.tick_label_number_format_is_linked = False

        assert chart.category_axis.title == "Quarter"
        assert value_axis.title == "Revenue"
        minimum_scale = value_axis.minimum_scale
        maximum_scale = value_axis.maximum_scale
        major_unit = value_axis.major_unit
        minor_unit = value_axis.minor_unit
        assert minimum_scale is not None and isclose(minimum_scale, 0.0)
        assert maximum_scale is not None and isclose(maximum_scale, 200.0)
        assert major_unit is not None and isclose(major_unit, 25.0)
        assert minor_unit is not None and isclose(minor_unit, 5.0)
        assert value_axis.tick_label_number_format == "$#,##0.00"
        assert value_axis.tick_label_number_format_is_linked is False
        prs.save(str(out_path))

    with zipfile.ZipFile(out_path) as zf:
        chart_xml = zf.read("ppt/charts/chart1.xml").decode("utf-8")
    assert "<a:t>Quarter</a:t>" in chart_xml
    assert "<a:t>Revenue</a:t>" in chart_xml
    assert '<c:min val="0"/>' in chart_xml
    assert '<c:max val="200"/>' in chart_xml
    assert '<c:majorUnit val="25"/>' in chart_xml
    assert '<c:minorUnit val="5"/>' in chart_xml
    assert '<c:numFmt formatCode="$#,##0.00" sourceLinked="0"/>' in chart_xml
