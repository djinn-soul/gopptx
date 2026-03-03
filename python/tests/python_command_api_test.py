"""Smoke tests for core presentation command API operations."""

import pathlib
import sys

import pytest
from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = (pathlib.Path(__file__).parent / "../..").resolve()
sys.path.append(str(project_root / "python"))


def test_python_command_api(tmp_path: pathlib.Path) -> None:
    """Run a minimal open/modify/save workflow via Python bindings."""
    input_deck = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"
    if not input_deck.exists():
        pytest.skip("Smoke sample missing for python_command_api_test")

    with Presentation(str(input_deck)) as pres:
        initial_count = pres.slide_count
        if initial_count <= 0:
            msg = "expected initial deck to have at least one slide"
            raise AssertionError(msg)

        meta = pres.metadata
        if meta is None:
            raise AssertionError("expected presentation metadata to be available")

        new_idx = pres.duplicate_slide(0, 1)
        if new_idx < 0:
            raise AssertionError(
                "expected duplicate_slide to return non-negative index"
            )

        pres.move_slide(0, pres.slide_count - 1)
        pres.remove_slide(1)

        output_path = tmp_path / "python_management_output.pptx"
        pres.save(output_path)

    if not output_path.exists():
        raise AssertionError("expected output file to be saved to disk")
