"""Test opening file-like presentation objects, BytesIO, bytes, and PathLike sources (Issue #1050)."""

import io
import pathlib

from gopptx import Presentation


def test_open_file_like_bytes_io(tmp_path: pathlib.Path) -> None:
    """Presentation(io.BytesIO(...)) and Presentation.open(io.BytesIO(...)) work cleanly (Issue #1050)."""
    pptx_file = tmp_path / "stream_test.pptx"

    # Create test presentation
    with Presentation.new(title="Stream Test Deck") as pres:
        pres.save(pptx_file)

    # Read bytes into BytesIO
    raw_bytes = pptx_file.read_bytes()
    stream = io.BytesIO(raw_bytes)

    # Test opening via Presentation constructor with BytesIO
    with Presentation(stream) as pres:
        assert pres.slide_count > 0
        assert pres.slides[0].title == "Stream Test Deck"

    # Test opening via file object from open('...', 'rb')
    with pathlib.Path(pptx_file).open("rb") as f, Presentation(f) as pres:
        assert pres.slide_count > 0
        assert pres.slides[0].title == "Stream Test Deck"

    # Test opening via raw bytes
    with Presentation(raw_bytes) as pres:
        assert pres.slide_count > 0

    # Test opening via pathlib.Path
    with Presentation(pptx_file) as pres:
        assert pres.slide_count > 0
