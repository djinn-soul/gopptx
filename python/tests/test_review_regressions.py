"""Regressions found by the 2026-07-29 differential review."""

import pathlib
import re
import zipfile

import pytest
from gopptx import Presentation
from gopptx.presentation.charts.chart_types import ChartType


def test_shape_chart_proxy_targets_its_own_graphic_frame(
    tmp_path: pathlib.Path,
) -> None:
    path = tmp_path / "shape_charts.pptx"
    with Presentation.new("Shape charts") as pres:
        slide = pres.slides[0]
        text_id = slide.add_textbox(100000, 100000, 1500000, 500000)
        first_id = slide.add_chart(
            ChartType.COLUMN,
            ["A", "B"],
            [1.0, 2.0],
            bounds=(100000, 800000, 3000000, 2000000),
        )
        second_id = slide.add_chart(
            ChartType.LINE,
            ["A", "B"],
            [3.0, 4.0],
            bounds=(3300000, 800000, 3000000, 2000000),
        )
        pres.save(path)

    with Presentation(path) as pres:
        slide = pres.slides[0]
        assert slide.shape(text_id).has_chart is False
        with pytest.raises(AttributeError, match="does not contain a chart"):
            _ = slide.shape(text_id).chart
        first = slide.shape(first_id)
        second = slide.shape(second_id)
        assert first.has_chart is True
        assert second.has_chart is True
        assert first.chart.index == 0
        assert second.chart.index == 1
        assert first.chart.rel_id != second.chart.rel_id


def test_shape_paragraphs_persist_independently(tmp_path: pathlib.Path) -> None:
    path = tmp_path / "paragraphs.pptx"
    with Presentation.new("Paragraphs") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(100000, 100000, 4000000, 1500000)
        frame = slide.shape(shape_id).text_frame
        frame.paragraphs[0].text = "First"
        second = frame.add_paragraph("Second")
        second.alignment = "r"
        second.runs[0].bold = True
        assert len(frame.paragraphs) == 2
        pres.save(path)

    with Presentation(path) as pres:
        paragraphs = pres.slides[0].shape(shape_id).text_frame.paragraphs
        assert len(paragraphs) == 2
        assert [paragraph.text for paragraph in paragraphs] == ["First", "Second"]
        assert paragraphs[1].alignment == "r"
        assert paragraphs[1].runs[0].bold is True


def test_axis_title_can_be_removed_and_reopened(tmp_path: pathlib.Path) -> None:
    path = tmp_path / "axis_title.pptx"
    with Presentation.new("Axis title") as pres:
        pres.slides[0].add_chart(
            ChartType.COLUMN,
            ["A", "B"],
            [1.0, 2.0],
            bounds=(100000, 500000, 4000000, 2500000),
        )
        pres.save(path)

    with Presentation(path) as pres:
        axis = pres.slides[0].charts[0].value_axis
        axis.title = "Revenue"
        assert axis.has_title is True
        axis.has_title = False
        pres.save(path)

    with Presentation(path) as pres:
        assert pres.slides[0].charts[0].value_axis.has_title is False
    with zipfile.ZipFile(path) as archive:
        chart_xml = archive.read("ppt/charts/chart1.xml").decode("utf-8")
    value_axis = re.search(r"<c:valAx>.*?</c:valAx>", chart_xml, re.DOTALL)
    assert value_axis is not None
    assert "<c:title>" not in value_axis.group(0)


def test_shape_metadata_can_be_cleared(tmp_path: pathlib.Path) -> None:
    path = tmp_path / "shape_metadata.pptx"
    with Presentation.new("Shape metadata") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(100000, 100000, 2000000, 500000)
        shape = slide.shape(shape_id)
        shape.alt_text = "Accessible text"
        shape.title = "Text box title"
        shape.alt_text = ""
        shape.title = ""
        pres.save(path)

    with Presentation(path) as pres:
        shape = pres.slides[0].shape(shape_id)
        assert not shape.alt_text
        assert not shape.title
    with zipfile.ZipFile(path) as archive:
        slide_xml = archive.read("ppt/slides/slide1.xml").decode("utf-8")
    assert 'descr="Accessible text"' not in slide_xml
    assert 'title="Text box title"' not in slide_xml


def test_remove_belongs_to_run_not_font(tmp_path: pathlib.Path) -> None:
    with Presentation.new("Run remove") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(100000, 100000, 2000000, 500000)
        paragraph = slide.shape(shape_id).text_frame.paragraphs[0]
        paragraph.text = "Remove me"
        run = paragraph.runs[0]
        assert not hasattr(run.font, "remove")
        run.remove()
        assert len(paragraph.runs) == 0
