import math
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


def test_insert_chart_with_multi_series_category_builder_fails_fast(
    presentation: Presentation,
):
    slide = presentation.slides[0]
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

    cat_data = CategoryChartData()
    cat_data.add_category("Q1")
    cat_data.add_series("Revenue", [10.0])
    cat_data.add_series("Cost", [8.0])

    with pytest.raises(
        ValueError,
        match="placeholder charts currently support exactly one series",
    ):
        ph.insert_chart("bar", cat_data)


def test_insert_chart_with_multi_series_xy_builder_fails_fast(
    presentation: Presentation,
):
    presentation.add_slide(
        "Chart Slide Multi XY", layout="title_and_content", bullets=["placeholder"]
    )
    slide = presentation.slides[presentation.slide_count - 1]
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

    xy_data = XyChartData()
    xy_data.add_series("S1", [1.0], [2.0])
    xy_data.add_series("S2", [3.0], [4.0])

    with pytest.raises(
        ValueError,
        match="placeholder charts currently support exactly one series",
    ):
        ph.insert_chart("scatter", xy_data)


def test_insert_chart_rejects_reserved_chart_option_values(presentation: Presentation):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="chart_options may not include reserved keys: values",
    ):
        ph.insert_chart("bar", ["A"], [1.0], values=[2.0])


def test_insert_chart_rejects_reserved_chart_option_bounds_key(
    presentation: Presentation,
):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="chart_options may not include reserved keys: x",
    ):
        ph.insert_chart("bar", ["A"], [1.0], x=123)


def test_insert_chart_rejects_nan_flat_values(presentation: Presentation):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="values must be an array of finite numbers",
    ):
        ph.insert_chart("bar", ["A"], [math.nan])


def test_insert_chart_rejects_infinite_scatter_series_values(
    presentation: Presentation,
):
    presentation.add_slide(
        "Chart Slide Scatter Inf", layout="title_and_content", bullets=["placeholder"]
    )
    slide = presentation.slides[presentation.slide_count - 1]
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

    xy_data = XyChartData()
    xy_data.add_series("S1", [1.0], [math.inf])
    with pytest.raises(
        ValueError,
        match="y_values must be an array of finite numbers",
    ):
        ph.insert_chart("scatter", xy_data)


def test_insert_chart_rejects_boolean_values(presentation: Presentation):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="values must be an array of finite numbers",
    ):
        ph.insert_chart("bar", ["A"], [True])


def test_insert_chart_rejects_non_string_categories(presentation: Presentation):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="chart_categories must be an array of strings",
    ):
        ph.insert_chart("bar", [1], [1.0])  # type: ignore[arg-type]


def test_insert_chart_rejects_non_string_builder_categories(presentation: Presentation):
    slide = presentation.slides[0]
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

    with pytest.raises(
        ValueError,
        match="categories must be an array of strings",
    ):
        ph.insert_chart(
            "bar",
            ["A"],
            [{"name": "S1", "values": [1.0], "categories": [2]}],  # type: ignore[list-item]
        )


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


def test_insert_table_rejects_non_positive_dimensions(presentation: Presentation):
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

    with pytest.raises(
        ValueError,
        match="table_rows must be a positive integer",
    ):
        ph.insert_table(0, 2)


def test_insert_table_rejects_boolean_dimensions(presentation: Presentation):
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

    with pytest.raises(
        ValueError,
        match="table_rows must be a positive integer",
    ):
        ph.insert_table(True, 2)


def test_set_placeholder_content_rejects_partial_table_dimensions(
    presentation: Presentation,
):
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

    with pytest.raises(
        ValueError,
        match="table_rows and table_cols must both be provided as positive integers",
    ):
        slide.set_placeholder_content(ph.idx, ph.placeholder_format.type, table_rows=2)


def test_insert_picture_rejects_nan_bounds(
    presentation: Presentation, tmp_path: pathlib.Path
):
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
    img = tmp_path / "nan-test.png"
    img.write_bytes(b"data")

    with pytest.raises(
        ValueError,
        match="bounds must be a tuple of four finite numbers",
    ):
        ph.insert_picture(str(img), bounds=(math.nan, 10, 100, 100))


def test_presentation_set_placeholder_content_rejects_infinite_bounds(
    presentation: Presentation,
):
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

    with pytest.raises(
        ValueError,
        match="bounds must be a tuple of four finite numbers",
    ):
        presentation.set_placeholder_content(
            slide.index,
            ph.idx,
            ph.placeholder_format.type,
            text="x",
            bounds=(10, math.inf, 100, 100),
        )
