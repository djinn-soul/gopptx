"""Chart axis ergonomics helpers and validation tests."""

from __future__ import annotations

import pytest

from gopptx import Presentation


def test_chart_axis_aliases_and_crosses_helpers() -> None:
    with Presentation.new("Chart Axis Ergonomics") as prs:
        slide = prs.add_slide("Chart")
        _ = slide.add_chart(
            "bar",
            ["A", "B"],
            [1.0, 2.0],
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        chart = slide.charts[0]
        axis = chart.category_axis

        assert axis.axis_kind == "category"
        assert axis.is_category_axis is True
        assert axis.is_value_axis is False

        axis.has_major_gridlines = True
        assert axis.has_major_gridlines is True

        axis.set_crosses_at_maximum()
        assert axis.crosses_at_maximum is True
        axis.set_crosses_auto_zero()
        assert axis.crosses_auto_zero is True
        axis.set_crosses_at_minimum()
        assert axis.crosses_at_minimum is True


def test_chart_axis_validation_errors() -> None:
    with Presentation.new("Chart Axis Validation") as prs:
        slide = prs.add_slide("Chart")
        _ = slide.add_chart(
            "bar",
            ["A", "B"],
            [1.0, 2.0],
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        axis = slide.charts[0].value_axis
        with pytest.raises(ValueError):
            axis.tick_label_position = "middle"
        with pytest.raises(ValueError):
            axis.crosses = "zero"
