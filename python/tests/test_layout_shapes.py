"""Adding shapes to slide layouts and masters (upstream python-pptx issue #1044)."""

import pathlib

import pytest
from gopptx.api_errors import GopptxError
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "layouts.pptx"
    with Presentation.new("Layout shapes") as pres:
        pres.add_slide("Slide A", layout="title_and_content", bullets=["a"])
        pres.add_slide("Slide B", layout="title_and_content", bullets=["b"])
        pres.save(output_path)
    return output_path


def test_add_textbox_to_layout(deck: pathlib.Path) -> None:
    """A layout text box is what the issue asks for and survives a round trip."""
    pres = Presentation()
    pres.open(deck)
    layout_part, _ = pres.get_slide_layout_ref(1)

    shape_id = pres.add_layout_textbox(
        layout_part,
        "Footer on every slide",
        (Inches(0.6), Inches(4.6), Inches(6.0), Inches(0.6)),
    )
    assert shape_id > 0
    pres.save(deck)
    pres.close()

    reopened = Presentation()
    reopened.open(deck)
    try:
        assert any(
            name.startswith("TextBox")
            for name in reopened.get_layout_shapes(layout_part) or []
        )
    finally:
        reopened.close()


def test_add_shape_to_layout_and_master(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    try:
        layout_part, master_part = pres.get_slide_layout_ref(1)

        layout_shape = pres.add_layout_shape(
            layout_part, "rect", (Inches(7.2), Inches(4.6), Inches(2.0), Inches(0.6))
        )
        master_shape = pres.add_master_shape(
            master_part, "rect", (Inches(0.4), Inches(0.4), Inches(1.0), Inches(0.4))
        )
        assert layout_shape > 0
        assert master_shape > 0
        assert layout_shape != master_shape or layout_part != master_part

        assert any(
            name.startswith("rect") for name in pres.get_layout_shapes(layout_part) or []
        )
        assert any(
            name.startswith("rect") for name in pres.get_master_shapes(master_part) or []
        )
    finally:
        pres.close()


def test_layout_shape_rejects_non_layout_part(deck: pathlib.Path) -> None:
    """A master path passed as a layout is refused rather than silently written."""
    pres = Presentation()
    pres.open(deck)
    try:
        _, master_part = pres.get_slide_layout_ref(1)
        with pytest.raises(GopptxError):
            pres.add_layout_textbox(
                master_part, "nope", (Inches(1), Inches(1), Inches(1), Inches(1))
            )
    finally:
        pres.close()


def test_layout_shape_ids_do_not_collide(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    try:
        layout_part, _ = pres.get_slide_layout_ref(1)
        ids = [
            pres.add_layout_textbox(
                layout_part, f"box {i}", (Inches(1), Inches(1 + i), Inches(2), Inches(0.4))
            )
            for i in range(3)
        ]
        assert len(set(ids)) == len(ids)
    finally:
        pres.close()
