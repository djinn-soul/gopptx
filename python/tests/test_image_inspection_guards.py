"""Guards on image source sniffing, resolution and axis state coercion."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx.presentation.shapes.image_inspection import (
    is_svg_source,
    picture_bounds,
    resolve_picture_source,
)

if TYPE_CHECKING:
    import pathlib

SVG_WITH_PROLOG = (
    b'<?xml version="1.0" encoding="UTF-8"?>\n'
    b"<!-- a comment before the root element -->\n"
    b'<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"/>'
)
NON_SVG_XML = b'<?xml version="1.0"?><rels><Relationship Id="rId1"/></rels>'


def test_xml_without_svg_root_is_not_svg() -> None:
    """An XML prolog alone must not take the SVG fallback."""
    assert not is_svg_source(NON_SVG_XML)


def test_svg_after_prolog_is_detected() -> None:
    assert is_svg_source(SVG_WITH_PROLOG)


def test_svg_path_is_detected(tmp_path: pathlib.Path) -> None:
    svg_file = tmp_path / "logo.SVG"
    svg_file.write_bytes(SVG_WITH_PROLOG)
    assert is_svg_source(str(svg_file))


def test_non_svg_xml_bytes_fail_to_decode() -> None:
    """Arbitrary XML must fail rather than silently size as a 2.54m EMU SVG."""
    with pytest.raises(Exception, match=r".*"):
        picture_bounds(NON_SVG_XML, 0, 0, 0, 0)


def test_empty_bytes_source_is_rejected() -> None:
    with pytest.raises(ValueError, match="must not be empty"):
        resolve_picture_source(b"", {})


def test_empty_data_option_is_rejected() -> None:
    with pytest.raises(ValueError, match="must not be empty"):
        resolve_picture_source(None, {"data": b""})


def test_missing_path_raises_file_not_found(tmp_path: pathlib.Path) -> None:
    missing = tmp_path / "nope.png"
    with pytest.raises(FileNotFoundError):
        picture_bounds(str(missing), 0, 0, 0, 0)
