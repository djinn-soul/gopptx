"""Per-label chart formatting tests: number format, font and display flags."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.presentation.charts import ChartType

if TYPE_CHECKING:
    from gopptx.slide.chart import Chart


def _bar_chart(prs: Presentation) -> Chart:
    slide = prs.add_slide("Chart")
    _ = slide.add_chart(
        ChartType.BAR,
        ["North", "South", "East", "West"],
        [1200.0, 0.35, 940.0, 1580.0],
        bounds=(1000000, 1000000, 5000000, 3000000),
    )
    return slide.charts[0]


def test_data_label_number_format_round_trips() -> None:
    """Upstream #638 and #803: one point's label takes its own format."""
    with Presentation.new("Chart Label Formats") as prs:
        chart = _bar_chart(prs)
        assert len(chart.data_label_points) == 0

        chart.set_data_label_number_format(1, "0.0%")

        label = chart.data_label_points.for_point(1)
        assert label is not None
        assert label.number_format == "0.0%"
        # A linked label ignores its own code, so asking for one unlinks it.
        assert label.format_linked is False


def test_data_label_font_keeps_category_flag() -> None:
    """Upstream #650: colouring a label must not drop the category name."""
    with Presentation.new("Chart Label Fonts") as prs:
        chart = _bar_chart(prs)
        chart.apply_format({
            "show_data_labels": True,
            "data_label_show_value": True,
            "data_label_show_category": True,
        })

        chart.format_data_label(3, font_color="0070C0", font_size_pt=16, font_bold=True)

        label = chart.data_label_points.for_point(3)
        assert label is not None
        assert label.font_color == "0070C0"
        assert label.font_size_pt == 16
        assert label.font_bold is True
        assert label.show_category is True
        assert label.show_value is True


def test_format_data_labels_bulk_and_delete() -> None:
    with Presentation.new("Chart Labels Bulk") as prs:
        chart = _bar_chart(prs)
        chart.format_data_labels([
            {"point_index": 0, "number_format": "#,##0"},
            {"point_index": 2, "font_color": "C00000"},
        ])
        assert [label.point_index for label in chart.data_label_points] == [0, 2]

        # A later patch merges rather than replacing the label.
        chart.format_data_label(0, font_bold=True)
        merged = chart.data_label_points.for_point(0)
        assert merged is not None
        assert merged.number_format == "#,##0"
        assert merged.font_bold is True

        chart.format_data_label(2, delete=True)
        deleted = chart.data_label_points.for_point(2)
        assert deleted is not None
        assert deleted.deleted is True


@pytest.mark.parametrize(
    "options",
    [
        {"font_color": "not-a-colour"},
        {"font_size_pt": 401},
        {"unknown_option": True},
    ],
)
def test_format_data_label_rejects_bad_options(options: dict[str, object]) -> None:
    with Presentation.new("Chart Label Guards") as prs:
        chart = _bar_chart(prs)
        with pytest.raises(ValueError):
            chart.format_data_label(0, **options)


def test_data_point_marker_colours_round_trip() -> None:
    """Upstream #825: a scatter point's colour lives on its marker."""
    with Presentation.new("Chart Markers") as prs:
        chart = _bar_chart(prs)
        chart.format_data_point(
            1,
            marker_fill_color="C00000",
            marker_line_color="7F0000",
            marker_symbol="circle",
            marker_size=10,
        )

        point = chart.data_points.for_series(0)[0]
        assert point.marker_fill_color == "C00000"
        assert point.marker_line_color == "7F0000"
        assert point.marker_symbol == "circle"
        assert point.marker_size == 10
        # The marker colours must not be reported as the point's own fill.
        assert point.fill_color is None


@pytest.mark.parametrize(
    "options",
    [
        {"marker_symbol": "hexagon"},
        {"marker_size": 100},
    ],
)
def test_format_data_point_rejects_bad_markers(options: dict[str, object]) -> None:
    with Presentation.new("Chart Marker Guards") as prs:
        chart = _bar_chart(prs)
        with pytest.raises(ValueError):
            chart.format_data_point(0, **options)
