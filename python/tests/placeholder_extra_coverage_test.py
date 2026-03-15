import pathlib

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.slide.chart.data import CategoryChartData, XyChartData
from gopptx.slide.placeholders.placeholder import PlaceholderFormat

project_root = pathlib.Path(__file__).parent.parent.parent.resolve()


@pytest.fixture
def presentation() -> Presentation:
    p = Presentation()
    p.open(str(project_root / "testdata" / "placeholders.pptx"))
    return p


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


def test_insert_chart_with_builders(presentation: Presentation):
    slide = presentation.slides[0]
    # Find a placeholder that supports charts (obj or body or chart)
    ph = next(
        (
            p
            for p in slide.placeholders
            if p.placeholder_format.type in {"body", "obj", "chart"}
        ),
        None,
    )
    if ph is None:
        pytest.skip("No suitable placeholder for chart in testdata")

    # CategoryChartData
    cat_data = CategoryChartData()
    cat_data.add_category("C1")
    cat_data.add_series("S1", [1.0])
    ph.insert_chart("bar", cat_data)

    # XyChartData
    xy_data = XyChartData()
    xy_data.add_series("S1", [1.0], [2.0])
    ph.insert_chart("scatter", xy_data)


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
