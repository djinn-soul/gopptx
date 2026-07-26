"""Misspelled shape keyword arguments must fail loudly, not be dropped.

add_shape/add_textbox/add_connector take **kwargs, so a typo used to be
silently discarded and the shape was created without the option the caller
asked for. These check the option name is validated instead.
"""

from __future__ import annotations

import pytest
from gopptx import GopptxError, Presentation
from gopptx.constants import ShapeType
from gopptx.schemas import Inches

BOUNDS = (Inches(1), Inches(1), Inches(2), Inches(1))


def test_add_shape_accepts_supported_options() -> None:
    """The supported option names still work."""
    with Presentation.new("Strict kwargs") as prs:
        shape_id = prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS, text="hello")
        assert shape_id > 0


@pytest.mark.parametrize(
    "bad_kwarg",
    ["tex", "clik_action", "propertie", "text_frme", "bounds"],
)
def test_add_shape_rejects_unknown_option(bad_kwarg: str) -> None:
    """A misspelled option raises TypeError naming the offending keyword."""
    with (
        Presentation.new("Strict kwargs") as prs,
        pytest.raises(TypeError, match=bad_kwarg),
    ):
        prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS, **{bad_kwarg: "x"})


def test_add_textbox_rejects_unknown_option() -> None:
    """add_textbox validates its option names too."""
    with (
        Presentation.new("Strict kwargs") as prs,
        pytest.raises(TypeError, match="txt"),
    ):
        prs.add_textbox(0, *BOUNDS, txt="oops")


def test_rejection_message_lists_supported_options() -> None:
    """The error tells the caller what it could have written instead."""
    with (
        Presentation.new("Strict kwargs") as prs,
        pytest.raises(TypeError) as excinfo,
    ):
        prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS, tex="oops")

    message = str(excinfo.value)
    assert "supported:" in message
    for supported in ("text", "runs", "click_action", "properties"):
        assert supported in message


def test_update_shape_still_forwards_unrecognized_fields() -> None:
    """update_shape passes through fields the shared helper does not handle.

    The strictness is opt-in for exactly this reason: ShapeUpdate carries many
    keys that _apply_shape_payload_options does not know about (rotation, x,
    fill, ...), and those must keep reaching the Go bridge.
    """
    with Presentation.new("Strict kwargs") as prs:
        shape_id = prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS, text="before")
        prs.update_shape(0, shape_id, {"text": "after", "rotation": 45, "x": 100})
        shapes = prs.list_shapes(0)

    updated = next(s for s in shapes if s.get("ID") == shape_id)
    assert updated.get("Text") == "after"
    assert updated.get("Rotation") == 45
    assert updated.get("X") == 100


def test_rotation_applies_to_a_non_picture_shape() -> None:
    """Rotating an autoshape works; rot lives on <a:xfrm>, not on pictures only."""
    with Presentation.new("Rotation") as prs:
        shape_id = prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS, text="spin")
        prs.update_shape(0, shape_id, {"rotation": 45, "flip_h": True})
        rotated = next(s for s in prs.list_shapes(0) if s.get("ID") == shape_id)

    assert rotated.get("Rotation") == 45


def test_crop_is_still_rejected_on_a_non_picture_shape() -> None:
    """Cropping needs <p:blipFill>, so it stays picture-only."""
    with Presentation.new("Crop") as prs:
        shape_id = prs.add_shape(0, ShapeType.RECTANGLE, BOUNDS)
        with pytest.raises(GopptxError, match="not a picture shape"):
            prs.update_shape(
                0,
                shape_id,
                {"crop": {"left": 0.1, "right": 0.1, "top": 0.1, "bottom": 0.1}},
            )
