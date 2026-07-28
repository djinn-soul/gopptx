"""Freeform geometry and picture-fill readback (Issue #1020)."""

import base64
import pathlib

from gopptx import Presentation

# 1x1 PNG.
PNG_BYTES = base64.b64decode(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk"
    "+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
)


def test_freeform_paths_are_readable(tmp_path: pathlib.Path) -> None:
    """A committed freeform reports its paths back through shape.freeform."""
    output_path = tmp_path / "freeform_readback.pptx"

    with Presentation.new(title="Freeform Readback") as pres:
        slide = pres.slides[0]
        builder = slide.build_freeform(start_x=100, start_y=200)
        builder.add_line_segments([(900, 200), (500, 800)])
        shape_id = builder.convert_to_shape(close=True)
        pres.save(output_path)

    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        shape = next(s for s in slide.shapes if s.id == shape_id)

        geometry = shape.freeform
        assert geometry is not None
        paths = geometry["paths"]
        assert len(paths) == 1

        segments = paths[0]["segments"]
        assert [segment["type"] for segment in segments] == [
            "moveTo",
            "lnTo",
            "lnTo",
            "close",
        ]
        # Points are stored local to the shape origin.
        assert segments[0]["points"][0]["x"] == 0
        assert segments[0]["points"][0]["y"] == 0
        assert segments[2]["points"][0]["x"] == 400
        assert segments[2]["points"][0]["y"] == 600


def test_preset_shape_has_no_freeform(tmp_path: pathlib.Path) -> None:
    """A preset-geometry shape reports no custom geometry."""
    output_path = tmp_path / "preset_no_freeform.pptx"

    with Presentation.new(title="Preset Geometry") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape("rect", (1000000, 1000000, 2000000, 1000000))
        pres.save(output_path)

    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        shape = next(s for s in slide.shapes if s.id == shape_id)
        assert shape.freeform is None


def test_picture_fill_is_readable(tmp_path: pathlib.Path) -> None:
    """A picture reports its embedded image through fill.picture or .image."""
    output_path = tmp_path / "picture_fill_readback.pptx"

    with Presentation.new(title="Picture Fill Readback") as pres:
        slide = pres.slides[0]
        picture_id = slide.add_picture(
            PNG_BYTES, left=1000000, top=1000000, width=2000000, height=2000000
        )
        pres.save(output_path)

    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        shape = next(s for s in slide.shapes if s.id == picture_id)
        assert shape.image.blob == PNG_BYTES
        # A picture stores its image in <p:blipFill>, not as a shape fill, so
        # fill.picture stays None here; it reports only <a:blipFill> fills.
        assert shape.fill.picture is None
