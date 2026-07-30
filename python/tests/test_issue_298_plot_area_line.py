"""Regression coverage for issue #298 plot-area outlines."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import CategoryChartData, Presentation
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path


def _chart_xml(path: Path) -> str:
    with zipfile.ZipFile(path) as archive:
        return archive.read("ppt/charts/chart1.xml").decode()


def test_plot_area_format_line_round_trips(tmp_path: Path) -> None:
    """The plot area's direct c:spPr carries the exact requested outline."""
    output_path = tmp_path / "plot_area_outline.pptx"
    data = CategoryChartData()
    for category in ("A", "B", "C"):
        data.add_category(category)
    data.add_series("Sales", [4, 7, 5])

    with Presentation.new("Issue 298") as presentation:
        slide = presentation.slides[0]
        _ = slide.add_chart(
            ChartType.BAR,
            data,
            bounds=(Inches(1), Inches(1.5), Inches(8), Inches(5)),
        )
        chart = slide.charts[0]
        chart.plot_area.format.line.color = "#C00000"
        chart.plot_area.format.line.width = 38100
        chart.plot_area.format.line.dash_style = "dash"
        presentation.save(str(output_path))

    xml = _chart_xml(output_path)
    expected = '<a:ln w="38100"><a:solidFill><a:srgbClr val="C00000"/>'
    assert expected in xml
    assert '<a:prstDash val="dash"/>' in xml

    with Presentation(str(output_path)) as presentation:
        line = presentation.slides[0].charts[0].plot_area.format.line
        assert line.color == "C00000"
        assert line.width == 38100
        assert line.dash_style == "dash"
