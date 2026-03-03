"""VBA project insertion smoke test."""

import pathlib
import zipfile
from pathlib import Path

import pytest
from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / "../..").resolve()
input_deck = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"


def test_add_vba_project(tmp_path: Path) -> None:
    """Test adding a VBA project to a presentation."""
    if not input_deck.exists():
        pytest.skip("smoke sample missing")

    out_path = tmp_path / "vba_test.pptm"

    with Presentation(input_deck) as prs:
        # Synthetic minimal vbaProject.bin (just some bytes for the test)
        mock_vba_data = b"fake_vba_bin_content"

        # Add the VBA project
        prs.add_vba_project(mock_vba_data)

        # Save the presentation
        prs.save(str(out_path))

    # Verify the VBA project was saved correctly
    if not out_path.exists():
        raise AssertionError("expected .pptm output with VBA project")

    with zipfile.ZipFile(out_path, "r") as zf:
        # Check if the file exists in the zip
        if "ppt/vbaProject.bin" not in zf.namelist():
            raise AssertionError("expected vba project stream in output package")

        # Check the content
        with zf.open("ppt/vbaProject.bin") as f:
            if f.read() != mock_vba_data:
                raise AssertionError("expected exact VBA payload bytes in package")

        # Check content types
        with zf.open("[Content_Types].xml") as f:
            content_types = f.read().decode("utf-8")
            if "application/vnd.ms-office.vbaProject" not in content_types:
                raise AssertionError("expected VBA content type in package manifest")
