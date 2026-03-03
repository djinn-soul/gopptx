import os
import pathlib

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118


def test_textbox_and_connector_compat() -> None:
    if not pathlib.Path(input_deck).exists():
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

        assert textbox_id > 0  # noqa: S101
        assert connector_id > 0  # noqa: S101

        shape_ids = {int(s["ID"]) for s in slide.list_shapes()}
        assert textbox_id in shape_ids  # noqa: S101
        assert connector_id in shape_ids  # noqa: S101


def test_group_and_freeform_creation() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.slides[0]
        group_id = slide.add_group_shape()
        assert group_id > 0  # noqa: S101

        builder = slide.build_freeform(914400, 914400)
        freeform_id = (
            builder
            .add_line_to(1828800, 914400)
            .add_line_to(1828800, 1828800)
            .convert_to_shape(close=True)
        )
        assert freeform_id > 0  # noqa: S101

        shape_ids = {int(s["ID"]) for s in slide.list_shapes()}
        assert group_id in shape_ids  # noqa: S101
        assert freeform_id in shape_ids  # noqa: S101


def test_group_creation_with_members() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity grouped members")
        shape_a = slide.add_shape("rect", (914400, 914400, 914400, 914400))
        shape_b = slide.add_shape("ellipse", (2743200, 914400, 914400, 914400))
        before_ids = {int(s["ID"]) for s in slide.list_shapes()}
        group_id = slide.add_group_shape([shape_a, shape_b])

        assert group_id > 0  # noqa: S101
        after_ids = {int(s["ID"]) for s in slide.list_shapes()}
        assert group_id in after_ids  # noqa: S101
        assert len(after_ids) <= len(before_ids) + 1  # noqa: S101


def test_freeform_builder_scale_and_segments() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity freeform scale")
        builder = slide.build_freeform(100, 100, scale=(2.0, 3.0))
        freeform_id = builder.add_line_segments([
            (200, 100),
            (200, 200),
        ]).convert_to_shape(close=False, text="freeform text")

        assert freeform_id > 0  # noqa: S101
        shapes = slide.list_shapes()
        ids = {int(s["ID"]) for s in shapes}
        assert freeform_id in ids  # noqa: S101
        assert any(
            int(shape["ID"]) == freeform_id and shape.get("Text") == "freeform text"
            for shape in shapes
        )


def test_freeform_builder_move_to_and_validation() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("python-pptx parity freeform validation")
        builder = slide.build_freeform(100, 100)

        with pytest.raises(ValueError, match="at least one line segment"):
            builder.convert_to_shape()

        freeform_id = builder.move_to(300, 300).add_line_to(400, 350).convert_to_shape()
        assert freeform_id > 0  # noqa: S101

        builder2 = slide.build_freeform(0, 0).add_line_to(10, 10)
        with pytest.raises(ValueError, match="only allowed before line segments"):
            builder2.move_to(20, 20)
