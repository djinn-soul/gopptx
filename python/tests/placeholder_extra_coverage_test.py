import pathlib

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.slide.chart.data import CategoryChartData, XyChartData
from gopptx.slide.placeholders.placeholder import PlaceholderFormat


def test_placeholder_format_props():
    pf = PlaceholderFormat("body", 10)
    assert pf.type == "body"
    assert pf.idx == 10
    assert str(pf) == "body"


def test_insert_text_style_aliases(presentation: Presentation):
    slide = presentation.slides[0]
    # Find a placeholder that supports text (usually any)
    ph = next(iter(slide.placeholders))
    ph.insert_text(
        "Styled Text", size=24, font_size=20, font_name="Arial", colour="FF0000"
    )


def test_insert_chart_with_builders(presentation: Presentation, tmp_path):
    # Add a fresh slide for each chart type to avoid clobbering the placeholder
    def _body_ph(pres: Presentation, slide_idx: int):
        slide = pres.slides[slide_idx]
        return next(
            (
                p
                for p in slide.placeholders
                if p.placeholder_format.type in {"body", "obj", "chart"}
            ),
            None,
        )

    # CategoryChartData — use the existing slide
    ph = _body_ph(presentation, 0)
    if ph is None:
        pytest.skip("No suitable placeholder for chart in testdata")
    cat_data = CategoryChartData()
    cat_data.add_category("C1")
    cat_data.add_series("S1", [1.0])
    ph.insert_chart("bar", cat_data)

    # XyChartData — add a fresh slide so the body placeholder is intact
    presentation.add_slide(
        "Chart Slide 2", layout="title_and_content", bullets=["placeholder"]
    )
    xy_ph = _body_ph(presentation, presentation.slide_count - 1)
    if xy_ph is None:
        pytest.skip("No body placeholder on newly added slide")
    xy_data = XyChartData()
    xy_data.add_series("S1", [1.0], [2.0])
    xy_ph.insert_chart("scatter", xy_data)


def test_placeholder_layout_master_refs(presentation: Presentation):
    slide = presentation.slides[0]
    ph = slide.placeholders[0]

    _ = ph.layout_placeholder
    _ = ph.master_placeholder


def test_insert_picture_with_bounds(presentation: Presentation, tmp_path: pathlib.Path):
    slide = presentation.slides[0]
    ph = next(
        (
            p
            for p in slide.placeholders
            if p.placeholder_format.type in {"body", "obj", "pic"}
        ),
        None,
    )
    if ph is None:
        pytest.skip("No suitable placeholder for picture")
    img = tmp_path / "test.png"
    img.write_bytes(b"data")
    ph.insert_picture(str(img), bounds=(10, 10, 100, 100))


def test_insert_table_with_bounds(presentation: Presentation):
    slide = presentation.slides[0]
    ph = next(
        (
            p
            for p in slide.placeholders
            if p.placeholder_format.type in {"body", "obj", "tbl"}
        ),
        None,
    )
    if ph is None:
        pytest.skip("No suitable placeholder for table")
    ph.insert_table(2, 2, bounds=(0, 0, 200, 200))
