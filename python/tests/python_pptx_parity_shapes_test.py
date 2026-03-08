"""python-pptx compatibility smoke checks for shape APIs."""

import pathlib

import pytest
from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / "../..").resolve()
input_deck = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"


def test_textbox_and_connector_compat() -> None:
    """Textbox/connector helper calls create real shapes."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity")
        textbox_id = slide.add_textbox(914400, 914400, 1828800, 914400, text="textbox")
        connector_id = slide.add_connector(
            "line",
            914400,
            1828800,
            2743200,
            2743200,
        )

        if textbox_id <= 0 or connector_id <= 0:
            raise AssertionError("expected positive ids for textbox and connector")

        shape_ids = {int(s["ID"]) for s in slide.list_shapes()}
        if textbox_id not in shape_ids or connector_id not in shape_ids:
            raise AssertionError("expected inserted shapes to exist in shape list")


def test_bulk_textbox_creation() -> None:
    """Slide bulk textbox helper creates all requested textboxes."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity bulk textboxes")
        textbox_ids = slide.add_textboxes([
            {
                "left": 914400,
                "top": 914400,
                "width": 1828800,
                "height": 914400,
                "text": "alpha",
            },
            {
                "left": 914400,
                "top": 1828800,
                "width": 1828800,
                "height": 914400,
                "text": "beta",
            },
        ])

        if len(textbox_ids) != 2 or any(shape_id <= 0 for shape_id in textbox_ids):
            raise AssertionError(f"expected two positive ids, got {textbox_ids!r}")

        texts_by_id = {
            int(shape["ID"]): shape.get("Text", "")
            for shape in slide.list_shapes()
            if int(shape["ID"]) in textbox_ids
        }
        if texts_by_id != {textbox_ids[0]: "alpha", textbox_ids[1]: "beta"}:
            raise AssertionError(f"unexpected bulk textbox texts: {texts_by_id!r}")


def test_buffered_textbox_id_remains_addressable() -> None:
    """Simple textbox inserts keep a stable real ID through buffered flush."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity buffered textbox")
        textbox_id = slide.add_textbox(914400, 914400, 1828800, 914400, text="queued")

        if textbox_id <= 0:
            raise AssertionError(f"expected positive textbox id, got {textbox_id!r}")

        shape = slide.shape(textbox_id)
        if shape.text != "queued":
            raise AssertionError(f"expected queued text, got {shape.text!r}")

        shape.text = "updated"
        texts_by_id = {
            int(shape_info["ID"]): shape_info.get("Text", "")
            for shape_info in slide.list_shapes()
            if int(shape_info["ID"]) == textbox_id
        }
        if texts_by_id != {textbox_id: "updated"}:
            raise AssertionError(f"unexpected textbox text state: {texts_by_id!r}")


def test_group_and_freeform_creation() -> None:
    """Group and freeform creation produce addressable shapes."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.slides[0]
        group_id = slide.add_group_shape()
        if group_id <= 0:
            raise AssertionError("expected positive group shape id")

        builder = slide.build_freeform(914400, 914400)
        freeform_id = (
            builder
            .add_line_to(1828800, 914400)
            .add_line_to(1828800, 1828800)
            .convert_to_shape(close=True)
        )
        if freeform_id <= 0:
            raise AssertionError("expected positive freeform shape id")

        shape_ids = {int(s["ID"]) for s in slide.list_shapes()}
        if group_id not in shape_ids or freeform_id not in shape_ids:
            raise AssertionError("expected created group/freeform ids in shape list")


def test_group_creation_with_members() -> None:
    """Grouping existing members yields a new group with stable topology."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity grouped members")
        shape_a = slide.add_shape("rect", (914400, 914400, 914400, 914400))
        shape_b = slide.add_shape("ellipse", (2743200, 914400, 914400, 914400))
        before_ids = {int(s["ID"]) for s in slide.list_shapes()}
        group_id = slide.add_group_shape([shape_a, shape_b])

        if group_id <= 0:
            raise AssertionError("expected positive grouped shape id")
        after_ids = {int(s["ID"]) for s in slide.list_shapes()}
        if group_id not in after_ids:
            raise AssertionError("expected group id in resulting shape ids")
        if len(after_ids) > len(before_ids) + 1:
            raise AssertionError("expected no unexpected shape count growth")


def test_freeform_builder_scale_and_segments() -> None:
    """Scaled freeform builder writes shape text and id correctly."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity freeform scale")
        builder = slide.build_freeform(100, 100, scale=(2.0, 3.0))
        freeform_id = builder.add_line_segments([
            (200, 100),
            (200, 200),
        ]).convert_to_shape(close=False, text="freeform text")

        if freeform_id <= 0:
            raise AssertionError("expected positive freeform id")
        shapes = slide.list_shapes()
        ids = {int(s["ID"]) for s in shapes}
        if freeform_id not in ids:
            raise AssertionError("expected freeform id present in shape list")
        expected_text_present = any(
            int(shape["ID"]) == freeform_id and shape.get("Text") == "freeform text"
            for shape in shapes
        )
        if not expected_text_present:
            raise AssertionError("expected freeform shape to retain provided text")


def test_freeform_builder_move_to_and_validation() -> None:
    """Builder validation for convert/move ordering stays enforced."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity freeform validation")
        builder = slide.build_freeform(100, 100)

        with pytest.raises(ValueError, match="at least one line segment"):
            builder.convert_to_shape()

        freeform_id = builder.move_to(300, 300).add_line_to(400, 350).convert_to_shape()
        if freeform_id <= 0:
            raise AssertionError("expected positive id after valid move_to + line")

        builder2 = slide.build_freeform(0, 0).add_line_to(10, 10)
        with pytest.raises(ValueError, match="only allowed before line segments"):
            builder2.move_to(20, 20)
