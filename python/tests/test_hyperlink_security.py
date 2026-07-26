"""Test hyperlink target URI sanitization and security defense (Issue #1140)."""

import pathlib
import zipfile

from gopptx import Presentation


def test_hyperlink_uri_sanitization(tmp_path: pathlib.Path) -> None:
    """Malicious hyperlink targets (javascript:, vbscript:, data:text/html) are sanitized to '#' (Issue #1140)."""
    output_path = tmp_path / "hyperlink_security_test.pptx"

    with Presentation.new(title="Hyperlink Security Test Deck") as pres:
        slide = pres.slides[0]

        # Add textbox with a safe URL hyperlink
        tb1_id = slide.add_textbox(1000000, 1000000, 4000000, 1000000, text="Safe Link")
        assert tb1_id > 0

        pres.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        slide_rels = zf.read("ppt/slides/_rels/slide1.xml.rels").decode("utf-8")
        assert "javascript:" not in slide_rels
        assert "vbscript:" not in slide_rels
