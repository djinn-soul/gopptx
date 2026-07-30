"""Chart trendline write, read-back, and validation tests."""

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


def test_trendline_round_trips() -> None:
    with Presentation.new("Chart Trendline") as prs:
        chart = _bar_chart(prs)
        assert len(chart.trendlines) == 0

        chart.add_trendline(
            "poly",
            order=3,
            name="Fit",
            display_r_squared=True,
            display_equation=True,
            line_color="FF0000",
            line_dash="dash",
        )
        trendlines = chart.trendlines
        assert len(trendlines) == 1
        line = trendlines[0]
        assert line.trendline_type == "poly"
        assert line.order == 3
        assert line.name == "Fit"
        assert line.display_r_squared is True
        assert line.display_equation is True
        assert line.series_index == 0

        # A second trendline on the same series must not evict the first.
        chart.add_trendline("movingAvg", period=2)
        kinds = sorted(item.trendline_type or "" for item in chart.trendlines)
        assert kinds == ["movingAvg", "poly"]

        chart.clear_trendlines()
        assert len(chart.trendlines) == 0


def test_set_trendlines_replaces_and_clears() -> None:
    with Presentation.new("Chart Trendline Replace") as prs:
        chart = _bar_chart(prs)
        chart.add_trendline("linear")
        chart.set_trendlines([{"type": "exp"}])
        assert [item.trendline_type for item in chart.trendlines] == ["exp"]

        chart.set_trendlines([])
        assert len(chart.trendlines) == 0


def test_trendline_survives_save_and_reopen(tmp_path: Path) -> None:
    output = tmp_path / "trendline.pptx"
    with Presentation.new("Chart Trendline Persist") as prs:
        chart = _bar_chart(prs)
        chart.add_trendline("linear", display_r_squared=True)
        prs.save(output)

    with Presentation(output) as prs:
        line = prs.slides[1].charts[0].trendlines[0]
        assert line.trendline_type == "linear"
        assert line.display_r_squared is True


def test_trendline_validation_errors() -> None:
    with Presentation.new("Chart Trendline Validation") as prs:
        chart = _bar_chart(prs)

        with pytest.raises(ValueError, match="trendline type must be one of"):
            chart.add_trendline("spline")
        with pytest.raises(ValueError, match="only valid for the poly"):
            chart.add_trendline("linear", order=2)
        with pytest.raises(ValueError, match="only valid for the movingAvg"):
            chart.add_trendline("linear", period=2)
        with pytest.raises(ValueError, match="order must be between 2 and 6"):
            chart.add_trendline("poly", order=9)
        with pytest.raises(ValueError, match="period must be greater than"):
            chart.add_trendline("movingAvg", period=1)
        with pytest.raises(ValueError, match="series_index must not be negative"):
            chart.add_trendline("linear", series_index=-1)
        with pytest.raises(ValueError, match="unknown trendline option"):
            chart.add_trendline("linear", slope=2)
        for field in ("forward", "backward", "intercept"):
            with pytest.raises(ValueError, match=rf"{field} must be a finite number"):
                chart.add_trendline("linear", **{field: float("nan")})
