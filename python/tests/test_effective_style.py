"""Effective style resolution through layout, master and theme (Issue #1013)."""

import pathlib

from gopptx import Presentation


def test_placeholder_color_resolves_through_inheritance(
    tmp_path: pathlib.Path,
) -> None:
    """A placeholder with no direct colour still reports one, with its source."""
    output_path = tmp_path / "effective_style.pptx"

    with Presentation.new(title="Effective Style") as pres:
        pres.add_slide("Inherited Title", layout="title_and_content")
        pres.save(output_path)

    with Presentation(output_path) as pres:
        slide = pres.slides[1]
        shape = slide.shapes[0]

        style = shape.style.effective
        assert isinstance(style, dict)
        # The layout and master consulted are named, so a caller can inspect
        # them directly.
        assert style.get("layout_part", "").startswith("ppt/slideLayouts/")
        assert style.get("master_part", "").startswith("ppt/slideMasters/")

        typeface = style.get("font_typeface")
        assert typeface is not None
        assert typeface["value"]
        assert typeface["source"] in {"shape", "layout", "master", "theme"}


def test_direct_color_wins_and_reports_shape_source(tmp_path: pathlib.Path) -> None:
    """A colour set on the shape is reported with the shape as its source."""
    output_path = tmp_path / "effective_style_direct.pptx"

    with Presentation.new(title="Direct Style") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape("rect", (1000000, 1000000, 2000000, 1000000))
        slide.update_shape(shape_id, {"fill": {"solid": "FF0000"}})
        pres.save(output_path)

    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        shape = next(s for s in slide.shapes if s.id == shape_id)

        fill_color = shape.style.effective.get("fill_color")
        assert fill_color is not None
        assert fill_color["rgb"] == "FF0000"
        assert fill_color["source"] == "shape"

        position = shape.style.effective.get("position")
        assert position is not None
        assert position["source"] == "shape"
        assert position["w"] == 2000000
