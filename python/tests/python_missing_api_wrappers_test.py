"""Regression tests for newly exposed Python bridge wrapper methods."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation
from gopptx.smartart import SMARTART_BASIC_LIST, SMARTART_BASIC_PROCESS

if TYPE_CHECKING:
    import pathlib


def _shape_ids(slide: object) -> list[int]:
    records = slide.list_shapes()  # type: ignore[attr-defined]
    out: list[int] = []
    for record in records:
        raw = record.get("ID", record.get("id"))
        if raw is not None:
            out.append(int(str(raw)))
    return out


def test_slide_master_and_layout_wrappers() -> None:
    """Presentation exposes functional add/remove wrapper methods for masters/layouts."""
    with Presentation.new("Master/Layout Wrapper Test") as pres:
        masters_before = list(pres.slide_masters)
        if not masters_before:
            raise AssertionError("expected at least one slide master")

        base_master_part = masters_before[0].part

        # Layout add/remove on existing master.
        layout_part = pres.add_slide_layout(base_master_part, "Wrapper Layout")
        if "slideLayout" not in layout_part:
            raise AssertionError(f"unexpected layout part path: {layout_part}")
        pres.remove_slide_layout(layout_part)

        # Master add/remove round-trip (regression for remove failure on fresh masters).
        master_part = pres.add_slide_master()
        if "slideMaster" not in master_part:
            raise AssertionError(f"unexpected master part path: {master_part}")
        pres.remove_slide_master(master_part)

        masters_after = list(pres.slide_masters)
        if len(masters_after) != len(masters_before):
            raise AssertionError(
                "expected slide master count to return to baseline after add/remove"
            )


def test_set_slide_hidden_wrapper_writes_show_flag(tmp_path: pathlib.Path) -> None:
    """set_slide_hidden writes the schema-valid hidden marker on the slide root."""
    output_path = tmp_path / "slide_hidden_wrapper.pptx"
    with Presentation.new("Hide Wrapper Test") as pres:
        pres.add_slide("Second Slide")
        pres.set_slide_hidden(1, hidden=True)
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        slide_xml = zf.read("ppt/slides/slide2.xml").decode("utf-8")
    if 'show="0"' not in slide_xml:
        raise AssertionError(
            'expected hidden slide marker show="0" on slide root XML'
        )


def test_smartart_wrapper_methods_cover_style_layout_nodes_delete() -> None:
    """Slide exposes wrapper methods for SmartArt style/layout/node/delete ops."""
    quick_style = "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"
    color_style = "urn:microsoft.com/office/officeart/2005/8/colors/colorful1"

    with Presentation.new("SmartArt Wrapper Test") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_smartart(SMARTART_BASIC_LIST, ["One", "Two"])
        if shape_id <= 0:
            raise AssertionError("expected positive SmartArt shape id")

        before_delete_ids = _shape_ids(slide)
        if shape_id not in before_delete_ids:
            raise AssertionError("expected SmartArt shape id to exist before delete")

        slide.set_smartart_style(
            shape_id,
            quick_style=quick_style,
            color_style=color_style,
        )
        slide.change_smartart_layout(shape_id, SMARTART_BASIC_PROCESS)
        slide.set_smartart_nodes(shape_id, ["Plan", "Execute", "Review"])
        slide.delete_smartart(shape_id)

        after_delete_ids = _shape_ids(slide)
        if shape_id in after_delete_ids:
            raise AssertionError("expected SmartArt shape to be deleted")
