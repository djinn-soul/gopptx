"""Regression coverage for issue #933 background inheritance setter."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation

if TYPE_CHECKING:
    from pathlib import Path


def _slide_xml(path: Path) -> str:
    with zipfile.ZipFile(path) as archive:
        return archive.read("ppt/slides/slide1.xml").decode()


def test_follow_master_background_is_settable(tmp_path: Path) -> None:
    """False adds no-fill p:bg and True removes the explicit background."""
    independent_path = tmp_path / "independent_background.pptx"
    inherited_path = tmp_path / "inherited_background.pptx"

    with Presentation.new("Issue 933") as presentation:
        slide = presentation.slides[0]
        assert slide.follow_master_background is True
        slide.follow_master_background = False
        assert slide.follow_master_background is False
        presentation.save(str(independent_path))

    independent_xml = _slide_xml(independent_path)
    expected = "<p:bg><p:bgPr><a:noFill/><a:effectLst/></p:bgPr></p:bg>"
    assert expected in independent_xml
    assert independent_xml.index("<p:bg>") < independent_xml.index("<p:spTree")

    with Presentation(str(independent_path)) as presentation:
        slide = presentation.slides[0]
        assert slide.follow_master_background is False
        slide.follow_master_background = True
        assert slide.follow_master_background is True
        presentation.save(str(inherited_path))

    assert "<p:bg" not in _slide_xml(inherited_path)


def test_false_preserves_custom_background(tmp_path: Path) -> None:
    """Disabling inheritance does not replace an existing custom background."""
    output_path = tmp_path / "custom_background.pptx"
    with Presentation.new("Issue 933 custom") as presentation:
        slide = presentation.slides[0]
        slide.background.set_solid("123456")
        slide.follow_master_background = False
        presentation.save(str(output_path))

    xml = _slide_xml(output_path)
    assert '<a:srgbClr val="123456"/>' in xml
    background_xml = xml[xml.index("<p:bg>") : xml.index("</p:bg>") + len("</p:bg>")]
    assert "<a:noFill/>" not in background_xml
