"""Table cell border tests (upstream python-pptx issue #71)."""

import pathlib
import re
import zipfile

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches

TCPR = re.compile(r"<a:tcPr.*?</a:tcPr>", re.DOTALL)


def _table_deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "borders.pptx"
    with Presentation.new("Cell borders") as pres:
        pres.slides[0].add_table(
            2,
            2,
            (Inches(1), Inches(1.5), Inches(6), Inches(2)),
            data=[["A1", "B1"], ["A2", "B2"]],
        )
        pres.save(output_path)
    return output_path


def _table(pres: Presentation):
    slide = pres.slides[0]
    shape_id = next(
        s.id for s in slide.shapes if s.shape_type in {"tbl", "graphicFrame"}
    )
    return slide.table(shape_id)


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    return _table_deck(tmp_path)


def test_cell_border_round_trips_cap_join_and_inset(deck: pathlib.Path) -> None:
    """Issue #71 asks for cap, join and inset pen beyond colour/width/style."""
    pres = Presentation()
    pres.open(deck)
    _table(pres).cell(0, 0).border_left = {
        "width": 38100,
        "color": "C00000",
        "dash": "dash",
        "cap": "square",
        "join": "miter",
        "compound": "double",
        "inset": True,
        "miter_limit": 400000,
    }
    pres.save(deck)
    pres.close()

    slide_xml = zipfile.ZipFile(deck).read("ppt/slides/slide1.xml").decode("utf-8")
    assert 'cap="sq"' in slide_xml
    assert 'cmpd="dbl"' in slide_xml
    assert 'algn="in"' in slide_xml
    assert '<a:miter lim="400000"/>' in slide_xml

    reopened = Presentation()
    reopened.open(deck)
    try:
        border = _table(reopened).cell(0, 0).border_left
        assert border is not None
        assert border["color"] == "C00000"
        assert border["dash"] == "dash"
        assert border["cap"] == "sq"
        assert border["compound"] == "dbl"
        assert border["join"] == "miter"
        assert border["inset"] is True
    finally:
        reopened.close()


def test_cell_border_precedes_cell_fill(deck: pathlib.Path) -> None:
    """CT_TableCellProperties puts the border lines before the fill."""
    pres = Presentation()
    pres.open(deck)
    table = _table(pres)
    table.update_cell(0, 0, {"background_color": "FFE699"})
    table.cell(0, 0).border_bottom = {
        "width": 38100,
        "color": "0070C0",
        "dash": "solid",
    }
    pres.save(deck)
    pres.close()

    slide_xml = zipfile.ZipFile(deck).read("ppt/slides/slide1.xml").decode("utf-8")
    for block in TCPR.findall(slide_xml):
        if "<a:lnB" not in block or "solidFill" not in block:
            continue
        fill_index = block.find('<a:solidFill><a:srgbClr val="FFE699"/>')
        if fill_index < 0:
            continue
        assert block.index("<a:lnB") < fill_index


def test_cell_border_can_be_cleared(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    _table(pres).cell(1, 1).border_top = {"width": 12700, "color": "000000"}
    pres.save(deck)
    pres.close()

    pres = Presentation()
    pres.open(deck)
    _table(pres).cell(1, 1).border_top = None
    pres.save(deck)
    pres.close()

    reopened = Presentation()
    reopened.open(deck)
    try:
        assert _table(reopened).cell(1, 1).border_top is None
    finally:
        reopened.close()
