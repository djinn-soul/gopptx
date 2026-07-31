"""Regression coverage for issue #845 master-background graphics."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation

if TYPE_CHECKING:
    from pathlib import Path


def _slide_xml(path: Path) -> str:
    with zipfile.ZipFile(path) as archive:
        return archive.read("ppt/slides/slide1.xml").decode()


def test_show_master_shapes_round_trips_exact_ooxml(tmp_path: Path) -> None:
    """False writes showMasterSp=0 and True restores the schema default."""
    hidden_path = tmp_path / "hidden_master_graphics.pptx"
    shown_path = tmp_path / "shown_master_graphics.pptx"

    with Presentation.new("Issue 845") as presentation:
        slide = presentation.slides[0]
        assert slide.show_master_shapes is True
        slide.show_master_shapes = False
        assert slide.show_master_shapes is False
        presentation.save(str(hidden_path))

    assert 'showMasterSp="0"' in _slide_xml(hidden_path)

    with Presentation(str(hidden_path)) as presentation:
        slide = presentation.slides[0]
        assert slide.show_master_shapes is False
        slide.show_master_shapes = True
        assert slide.show_master_shapes is True
        presentation.save(str(shown_path))

    assert "showMasterSp" not in _slide_xml(shown_path)
