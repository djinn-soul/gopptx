"""Chart data-table write, read-back, and validation tests."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.presentation.charts import ChartType

if TYPE_CHECKING:
    from pathlib import Path

    from gopptx.slide.chart import Chart

_FONT_SIZE_PT = 9


def _bar_chart(prs: Presentation) -> Chart:
    slide = prs.add_slide("Chart")
    _ = slide.add_chart(
        ChartType.BAR,
        ["A", "B", "C"],
        [1.0, 4.0, 2.0],
        bounds=(1000000, 1000000, 5000000, 3000000),
    )
    return slide.charts[0]


def test_data_table_round_trips() -> None:
    with Presentation.new("Chart Data Table") as prs:
        chart = _bar_chart(prs)
        assert chart.data_table.visible is False

        chart.data_table.show(vertical_border=False, font_size_pt=_FONT_SIZE_PT)
        table = chart.data_table
        assert table.visible is True
        assert table.show_vertical_border is False
        # Unspecified flags default to on.
        assert table.show_horizontal_border is True
        assert table.show_outline is True
        assert table.font_size_pt == _FONT_SIZE_PT

        # Setting one flag must leave the others alone.
        chart.data_table.show_legend_keys = False
        assert chart.data_table.show_vertical_border is False
        assert chart.data_table.show_legend_keys is False

        chart.data_table.hide()
        assert chart.data_table.visible is False


def test_data_table_survives_save_and_reopen(tmp_path: Path) -> None:
    output = tmp_path / "data_table.pptx"
    with Presentation.new("Chart Data Table Persist") as prs:
        chart = _bar_chart(prs)
        chart.data_table.show(legend_keys=True)
        prs.save(output)

    with Presentation(output) as prs:
        table = prs.slides[1].charts[0].data_table
        assert table.visible is True
        assert table.show_legend_keys is True


def test_data_table_validation_errors() -> None:
    with Presentation.new("Chart Data Table Validation") as prs:
        chart = _bar_chart(prs)

        with pytest.raises(ValueError, match="font_size_pt must be between 1 and 400"):
            chart.data_table.show(font_size_pt=500)
        # Options only make sense once the table exists.
        with pytest.raises(ValueError, match="no data table"):
            chart.data_table.show_outline = False
