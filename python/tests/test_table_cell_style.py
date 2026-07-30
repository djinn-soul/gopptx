"""update_table_cell style fields reach the XML (found while fixing #1037)."""

import pathlib
import zipfile

import pytest
from gopptx.api_errors import GopptxError
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches


def _table(pres: Presentation):
    slide = pres.slides[0]
    shape_id = next(
        s.id for s in slide.shapes if s.shape_type in {"tbl", "graphicFrame"}
    )
    return slide.table(shape_id)


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "cells.pptx"
    with Presentation.new("Cell styling") as pres:
        pres.slides[0].add_table(
            2,
            2,
            (Inches(0.8), Inches(1.5), Inches(8.0), Inches(2.0)),
            data=[["a", "b"], ["c", "d"]],
        )
        pres.save(output_path)
    return output_path


def _slide_xml(deck_path: pathlib.Path) -> str:
    with zipfile.ZipFile(deck_path) as archive:
        return archive.read("ppt/slides/slide1.xml").decode("utf-8")


def test_bold_colour_and_fill_are_applied(deck: pathlib.Path) -> None:
    """These three were accepted and silently dropped before."""
    pres = Presentation()
    pres.open(deck)
    _table(pres).update_cell(
        0, 0, {"bold": True, "color": "C00000", "background_color": "FFE699"}
    )
    pres.save(deck)
    pres.close()

    xml = _slide_xml(deck)
    assert 'b="1"' in xml
    assert "C00000" in xml
    assert "FFE699" in xml


def test_italic_and_underline_are_applied(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    _table(pres).update_cell(0, 1, {"italic": True, "underline": True})
    pres.save(deck)
    pres.close()

    xml = _slide_xml(deck)
    assert 'i="1"' in xml
    assert 'u="sng"' in xml


def test_style_and_text_apply_together(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    _table(pres).update_cell(1, 0, {"text": "styled", "bold": True, "size_pt": 24})
    pres.save(deck)
    pres.close()

    xml = _slide_xml(deck)
    assert "styled" in xml
    assert 'b="1"' in xml
    assert 'sz="2400"' in xml

    reopened = Presentation()
    reopened.open(deck)
    try:
        assert _table(reopened).cell(1, 0).text == "styled"
    finally:
        reopened.close()


def test_wrong_type_is_rejected_rather_than_ignored(deck: pathlib.Path) -> None:
    """A bad value must fail loudly, not vanish."""
    pres = Presentation()
    pres.open(deck)
    try:
        with pytest.raises(GopptxError):
            _table(pres).update_cell(0, 0, {"bold": "yes"})
    finally:
        pres.close()
