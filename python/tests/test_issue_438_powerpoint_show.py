"""Regression coverage for valid PowerPoint Show output (issue #438)."""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation

if TYPE_CHECKING:
    from pathlib import Path

_PRESENTATION_MAIN = (
    "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"
)
_SLIDESHOW_MAIN = (
    "application/vnd.openxmlformats-officedocument.presentationml.slideshow.main+xml"
)


def test_save_ppsx_writes_slideshow_main_content_type(tmp_path: Path) -> None:
    """A .ppsx suffix changes the package kind rather than only its filename."""
    output = tmp_path / "show.ppsx"
    with Presentation.new("Issue 438") as presentation:
        presentation.slides[0].title = "PowerPoint Show"
        presentation.save(str(output))

    with zipfile.ZipFile(output) as package:
        content_types = package.read("[Content_Types].xml").decode()
    assert _SLIDESHOW_MAIN in content_types
    assert _PRESENTATION_MAIN not in content_types

    reopened = Presentation()
    reopened.open(output)
    try:
        assert reopened.slides[0].title == "PowerPoint Show"
    finally:
        reopened.close()
