"""Chart error-bar write, read-back, and validation tests."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.presentation.charts import ChartType

if TYPE_CHECKING:
    from pathlib import Path

    from gopptx.slide.chart import Chart


def _bar_chart(prs: Presentation) -> Chart:
    slide = prs.add_slide("Chart")
    _ = slide.add_chart(
        ChartType.BAR,
        ["A", "B", "C", "D"],
        [1.0, 4.0, 2.0, 8.0],
        bounds=(1000000, 1000000, 5000000, 3000000),
    )
    return slide.charts[0]


def test_error_bars_round_trip() -> None:
    with Presentation.new("Chart Error Bars") as prs:
        chart = _bar_chart(prs)
        assert len(chart.error_bars) == 0

        chart.add_error_bars(
            "both", "fixedVal", value=1.5, no_end_cap=True, line_color="0000FF"
        )
        bars = chart.error_bars
        assert len(bars) == 1
        assert bars[0].bar_type == "both"
        assert bars[0].value_type == "fixedVal"
        assert bars[0].value == pytest.approx(1.5)
        assert bars[0].no_end_cap is True
        assert bars[0].series_index == 0

        chart.clear_error_bars()
        assert len(chart.error_bars) == 0


def test_custom_and_directional_error_bars() -> None:
    with Presentation.new("Chart Custom Error Bars") as prs:
        chart = _bar_chart(prs)
        chart.set_error_bars([
            {
                "bar_type": "both",
                "value_type": "cust",
                "direction": "y",
                "plus_reference": "Sheet1!$D$2:$D$5",
                "minus_reference": "Sheet1!$E$2:$E$5",
            },
            {"bar_type": "plus", "value_type": "stdDev", "direction": "x"},
        ])
        bars = chart.error_bars
        assert len(bars) == 2
        assert bars[0].plus_reference == "Sheet1!$D$2:$D$5"
        assert bars[0].minus_reference == "Sheet1!$E$2:$E$5"
        assert bars[0].direction == "y"
        assert bars[1].value_type == "stdDev"
        assert bars[1].direction == "x"


def test_error_bars_survive_save_and_reopen(tmp_path: Path) -> None:
    output = tmp_path / "error_bars.pptx"
    with Presentation.new("Chart Error Bars Persist") as prs:
        chart = _bar_chart(prs)
        chart.add_error_bars("both", "percentage", value=5.0)
        prs.save(output)

    with Presentation(output) as prs:
        bars = prs.slides[1].charts[0].error_bars[0]
        assert bars.value_type == "percentage"
        assert bars.value == pytest.approx(5.0)


def test_error_bars_validation_errors() -> None:
    with Presentation.new("Chart Error Bars Validation") as prs:
        chart = _bar_chart(prs)

        with pytest.raises(ValueError, match="bar_type must be one of"):
            chart.add_error_bars("up", "fixedVal")
        with pytest.raises(ValueError, match="value_type must be one of"):
            chart.add_error_bars("both", "guess")
        with pytest.raises(ValueError, match="direction must be one of"):
            chart.add_error_bars("both", "stdDev", direction="z")
        with pytest.raises(ValueError, match="requires plus_reference"):
            chart.add_error_bars("both", "cust")
        with pytest.raises(ValueError, match="require value_type 'cust'"):
            chart.add_error_bars("both", "stdDev", plus_reference="Sheet1!$D$2")
        with pytest.raises(ValueError, match="value must be a non-negative"):
            chart.add_error_bars("both", "fixedVal", value=-1.0)
        with pytest.raises(ValueError, match="series_index must not be negative"):
            chart.add_error_bars("both", "stdDev", series_index=-1)
        with pytest.raises(ValueError, match="unknown error bar option"):
            chart.add_error_bars("both", "stdDev", spread=2)
