"""Test combo chart creation and rendering (Issue #1101)."""

import pathlib

from gopptx import Presentation
from gopptx.schemas import Inches


def test_add_combo_chart(tmp_path: pathlib.Path) -> None:
    """Combo chart (bar + line series) creation works on slide proxies (Issue #1101)."""
    output_path = tmp_path / "combo_chart_test.pptx"

    with Presentation.new(title="Combo Chart Test Deck") as pres:
        slide = pres.slides[0]

        # Add a combo chart
        chart_id = slide.add_combo_chart(
            categories=["Q1", "Q2", "Q3", "Q4"],
            bar_series=[{"name": "Revenue ($M)", "values": [12.5, 18.0, 22.4, 30.1]}],
            line_series=[{"name": "Growth (%)", "values": [8.5, 12.1, 15.3, 20.0]}],
            title="Quarterly Performance",
            bounds=(Inches(1), Inches(1), Inches(8), Inches(5)),
        )
        assert chart_id > 0

        pres.save(output_path)

    # Verify presentation can be re-opened
    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        charts = slide.list_charts()
        assert len(charts) > 0
