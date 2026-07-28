"""Per-point chart formatting and negative-value inversion tests."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.presentation.charts import ChartType

if TYPE_CHECKING:
    from pathlib import Path

    from gopptx.slide.chart import Chart

_NEGATIVE_POINT_INDEX = 2


def _bar_chart(prs: Presentation) -> Chart:
    slide = prs.add_slide("Chart")
    _ = slide.add_chart(
        ChartType.BAR,
        ["A", "B", "C", "D"],
        [4.0, 7.0, -3.0, 9.0],
        bounds=(1000000, 1000000, 5000000, 3000000),
    )
    return slide.charts[0]


def test_data_point_formatting_round_trips() -> None:
    with Presentation.new("Chart Data Points") as prs:
        chart = _bar_chart(prs)
        assert len(chart.data_points) == 0

        chart.format_data_point(3, fill_color="FF0000")
        chart.format_data_point(1, line_color="00FF00", line_width_emu=12700)

        points = chart.data_points
        assert [point.point_index for point in points] == [1, 3]
        assert points[1].fill_color == "FF0000"
        assert points[0].line_color == "00FF00"
        assert points[0].line_width_emu == 12700

        # Merging must keep the properties set by the earlier call.
        chart.format_data_point(3, line_color="0000FF")
        merged = chart.data_points.for_series(0)[1]
        assert merged.fill_color == "FF0000"
        assert merged.line_color == "0000FF"

        chart.clear_data_point_formatting()
        assert len(chart.data_points) == 0


def test_format_data_points_bulk() -> None:
    with Presentation.new("Chart Data Points Bulk") as prs:
        chart = _bar_chart(prs)
        chart.format_data_points([
            {"point_index": 0, "fill_color": "111111"},
            {"point_index": 2, "fill_color": "222222"},
        ])
        assert [point.point_index for point in chart.data_points] == [0, 2]


def test_invert_if_negative_colours_negative_points() -> None:
    with Presentation.new("Chart Invert") as prs:
        chart = _bar_chart(prs)
        chart.set_invert_if_negative(negative_fill_color="C00000")

        points = chart.data_points
        assert len(points) == 1
        assert points[0].point_index == _NEGATIVE_POINT_INDEX
        assert points[0].fill_color == "C00000"
        # The point turns its own inversion off so the fill is what is drawn.
        assert points[0].invert_if_negative is False


def test_data_point_formatting_survives_save_and_reopen(tmp_path: Path) -> None:
    output = tmp_path / "data_points.pptx"
    with Presentation.new("Chart Data Points Persist") as prs:
        chart = _bar_chart(prs)
        chart.format_data_point(1, fill_color="ABCDEF")
        prs.save(output)

    with Presentation(output) as prs:
        point = prs.slides[1].charts[0].data_points[0]
        assert point.point_index == 1
        assert point.fill_color == "ABCDEF"


def test_data_point_validation_errors() -> None:
    with Presentation.new("Chart Data Point Validation") as prs:
        chart = _bar_chart(prs)

        with pytest.raises(ValueError, match="point_index must not be negative"):
            chart.format_data_point(-1)
        with pytest.raises(ValueError, match="series_index must not be negative"):
            chart.format_data_point(0, series_index=-1)
        with pytest.raises(ValueError, match="explosion must be between 0 and 400"):
            chart.format_data_point(0, explosion=500)
        with pytest.raises(ValueError, match="line_width_emu must not be negative"):
            chart.format_data_point(0, line_width_emu=-1)
        with pytest.raises(ValueError, match="unknown data point option"):
            chart.format_data_point(0, glow="red")
