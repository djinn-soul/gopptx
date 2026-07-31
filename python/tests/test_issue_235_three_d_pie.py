"""Regression coverage for issue #235 3D pie charts."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import CategoryChartData, Presentation
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path


def test_three_d_pie_writes_exact_chart_elements(tmp_path: Path) -> None:
    """THREE_D_PIE serializes c:view3D and c:pie3DChart."""
    output_path = tmp_path / "three_d_pie.pptx"
    data = CategoryChartData()
    for category in ("Desktop", "Mobile", "Tablet"):
        data.add_category(category)
    data.add_series("Traffic", [48, 39, 13])

    with Presentation.new("Issue 235") as presentation:
        slide = presentation.slides[0]
        _ = slide.add_chart(
            ChartType.THREE_D_PIE,
            data,
            title="Traffic share",
            bounds=(Inches(1), Inches(1.4), Inches(8), Inches(5.4)),
        )
        presentation.save(str(output_path))

    with zipfile.ZipFile(output_path) as package:
        xml = package.read("ppt/charts/chart1.xml").decode()
    assert "<c:pie3DChart>" in xml
    assert "<c:pieChart>" not in xml
    assert '<c:rotX val="30"/>' in xml
    assert '<c:perspective val="30"/>' in xml
    assert xml.index("<c:view3D>") < xml.index("<c:plotArea>")
