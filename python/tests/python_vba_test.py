import zipfile
import os
import pathlib
from pathlib import Path

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")

def test_add_vba_project(tmp_path: Path) -> None:
    """Test adding a VBA project to a presentation."""
    if not pathlib.Path(input_deck).exists():
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
    assert out_path.exists()
    
    with zipfile.ZipFile(out_path, "r") as zf:
        # Check if the file exists in the zip
        assert "ppt/vbaProject.bin" in zf.namelist()
        
        # Check the content
        with zf.open("ppt/vbaProject.bin") as f:
            assert f.read() == mock_vba_data
        
        # Check content types
        with zf.open("[Content_Types].xml") as f:
            content_types = f.read().decode("utf-8")
            assert "application/vnd.ms-office.vbaProject" in content_types
