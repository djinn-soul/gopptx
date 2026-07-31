"""Issue #309: access slide shapes by name.

Names come from p:cNvPr/@name, so the tests check the facade against the names
actually written into the saved slide part.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

_CNVPR_NAME_RE = re.compile(r'<p:cNvPr[^>]*\bname="([^"]*)"')


def test_names_match_the_saved_slide_part(tmp_path: Path) -> None:
    """names() reports the p:cNvPr names present in the package."""
    output_path = tmp_path / "named.pptx"

    with Presentation.new(title="Issue 309") as prs:
        slide = prs.slides[0]
        slide.shapes.add_textbox(Inches(1), Inches(1), Inches(4), Inches(1), text="A")
        facade_names = slide.shapes.names()
        prs.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
    # The first cNvPr belongs to the spTree group, not to a shape.
    xml_names = _CNVPR_NAME_RE.findall(xml)[1:]

    assert facade_names
    assert facade_names == xml_names


def test_lookup_by_name_returns_the_same_shape() -> None:
    """shapes[name] resolves to the shape carrying that name."""
    with Presentation.new(title="Issue 309 hit") as prs:
        slide = prs.slides[0]
        box = slide.shapes.add_textbox(
            Inches(1), Inches(1), Inches(4), Inches(1), text="A"
        )

        found = slide.shapes[box.name]

        assert found.id == box.id
        assert slide.shapes.get_by_name(box.name) is not None


def test_unknown_name_raises_key_error() -> None:
    """A missing name raises KeyError rather than returning None from [] access."""
    with (
        Presentation.new(title="Issue 309 miss") as prs,
        pytest.raises(KeyError, match="no shape named"),
    ):
        _ = prs.slides[0].shapes["definitely-not-here"]


def test_get_by_name_returns_none_for_unknown_name() -> None:
    """get_by_name is the non-raising form."""
    with Presentation.new(title="Issue 309 none") as prs:
        assert prs.slides[0].shapes.get_by_name("definitely-not-here") is None


def test_integer_indexing_still_works() -> None:
    """Adding string keys must not break positional access or negative indices."""
    with Presentation.new(title="Issue 309 int") as prs:
        slide = prs.slides[0]
        last = slide.shapes.add_textbox(
            Inches(1), Inches(1), Inches(4), Inches(1), text="LAST"
        )

        assert slide.shapes[0].id > 0
        assert slide.shapes[-1].id == last.id
        with pytest.raises(IndexError):
            _ = slide.shapes[999]


def test_lookup_returns_the_first_match_in_document_order() -> None:
    """Names are not unique in OOXML, so lookup resolves the first match."""
    with Presentation.new(title="Issue 309 order") as prs:
        slide = prs.slides[0]
        slide.shapes.add_textbox(Inches(1), Inches(1), Inches(3), Inches(1), text="ONE")
        slide.shapes.add_textbox(Inches(1), Inches(3), Inches(3), Inches(1), text="TWO")

        names = slide.shapes.names()
        target = names[0]

        assert slide.shapes[target].id == slide.shapes[names.index(target)].id
