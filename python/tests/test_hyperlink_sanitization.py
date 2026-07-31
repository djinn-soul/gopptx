"""Hyperlink script protocols must never reach the saved relationship target."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

SLIDE_RELS = "ppt/slides/_rels/slide1.xml.rels"


def _save_run_hyperlink(path: Path, address: str) -> str:
    """Write a deck whose only run carries `address`, return slide1 rels XML."""
    with Presentation.new(title="Hyperlink Sanitization") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(Inches(1), Inches(2), Inches(5), Inches(1))
        paragraph = slide.shape(shape_id).text_frame.paragraphs[0]
        paragraph.text = "click"
        paragraph.runs[0].hyperlink = address
        pres.save(path)

    with zipfile.ZipFile(path) as archive:
        return archive.read(SLIDE_RELS).decode("utf-8")


@pytest.mark.parametrize(
    "address",
    [
        "javascript:alert(1)",
        "JavaScript:alert(1)",
        "  javascript:alert(1)",
        "vbscript:msgbox(1)",
        "data:text/html;base64,PHNjcmlwdD4=",
    ],
)
def test_script_protocol_hyperlinks_are_neutralized(
    tmp_path: Path, address: str
) -> None:
    rels_xml = _save_run_hyperlink(tmp_path / "blocked.pptx", address)

    lowered = rels_xml.lower()
    assert "javascript:" not in lowered
    assert "vbscript:" not in lowered
    assert "data:text/html" not in lowered
    assert 'Target="#"' in rels_xml


@pytest.mark.parametrize(
    "address",
    [
        "https://example.com/a?b=c",
        "mailto:someone@example.com",
    ],
)
def test_safe_hyperlinks_are_preserved(tmp_path: Path, address: str) -> None:
    rels_xml = _save_run_hyperlink(tmp_path / "safe.pptx", address)

    assert f'Target="{address}"' in rels_xml
