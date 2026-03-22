"""Parity tests for text-frame facade objects and alias handling."""

from pathlib import Path

import pytest
from gopptx import Presentation, TextFrameProps
from gopptx.slide.text._utils import as_optional_int


def test_text_frame_props_aliases_round_trip(tmp_path: Path) -> None:
    """Ensure alias-based text-frame settings serialize through shape payloads."""
    out_path = tmp_path / "text_frame_props_facade.pptx"

    with Presentation.new("TextFrame Facade") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1500000),
            text="frame",
            text_frame=TextFrameProps(
                margin_top=90000,
                margin_bottom=90000,
                margin_left=180000,
                margin_right=180000,
                word_wrap=True,
                vertical_anchor="middle",
                auto_size="shape_to_fit_text",
                text_direction="vertical_270",
                column_count=2,
                text_rotation=45.0,
            ),
        )
        assert shape_id > 0
        prs.save(out_path)

    with Presentation(out_path) as prs:
        texts = [shape["Text"] for shape in prs.slides[0].list_shapes()]
        assert "frame" in texts


def test_text_frame_props_rejects_out_of_range_rotation() -> None:
    """Ensure invalid text-frame rotation values fail fast."""
    with Presentation.new("TextFrame Unsupported") as prs:
        slide = prs.slides[0]
        with pytest.raises(ValueError, match="between -360 and 360"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                text="unsupported",
                text_frame={"rotation": 720},
            )


def test_text_frame_props_init_validation() -> None:
    """Ensure TextFrameProps validates values in __init__."""
    with pytest.raises(ValueError, match="unsupported vertical alignment"):
        TextFrameProps(vertical_align="invalid")

    with pytest.raises(ValueError, match="unsupported auto_fit_type"):
        TextFrameProps(auto_fit_type="invalid")

    with pytest.raises(ValueError, match="unsupported orientation"):
        TextFrameProps(orientation="invalid")


def test_as_optional_int_with_floats() -> None:
    """Ensure as_optional_int handles integer-like floats."""
    assert as_optional_int(1.0) == 1
    assert as_optional_int(1.5) is None
    assert as_optional_int("1") is None
