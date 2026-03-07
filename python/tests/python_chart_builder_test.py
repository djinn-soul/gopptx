"""Tests for chart data builder classes."""
# ruff: noqa: D103

from __future__ import annotations

from gopptx import CategoryChartData, Presentation, XyChartData


def test_category_chart_builder_add_chart() -> None:
    with Presentation.new("Chart Builder") as prs:
        slide = prs.add_slide("Chart")
        data = CategoryChartData()
        data.add_category("A")
        data.add_category("B")
        data.add_series("Series 1", [1.0, 2.0])

        _shape_id = slide.add_chart("bar", data, title="Builder Chart")
        charts = slide.list_charts()
        assert len(charts) >= 1  # noqa: S101


def test_xy_chart_builder_update_payload() -> None:
    data = XyChartData()
    data.add_series("S1", [1.0, 2.0], [3.0, 4.0], sizes=[5.0, 6.0])
    payload = data.to_update_payload()
    assert "series" in payload  # noqa: S101
    series = payload["series"]
    assert isinstance(series, list) and len(series) == 1  # noqa: S101
