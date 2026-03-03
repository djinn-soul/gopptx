"""Smoke test for adding an image through Python bindings."""

import pathlib
import sys

import pytest
from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = (pathlib.Path(__file__).parent / "../..").resolve()
sys.path.append(str(project_root / "python"))


def test_python_image(tmp_path: pathlib.Path) -> None:
    """Image insertion writes a valid output deck."""
    image_path = (
        project_root
        / "examples"
        / "assets"
        / "55"
        / "repository-open-graph-template.png"
    )
    if not image_path.exists():
        pytest.skip("Image asset not available for python_image_test")

    output_path = tmp_path / "python_image_test.pptx"

    with Presentation.new("Image Test") as pres:
        slide = pres.slides[0]

        # Add image to the first slide
        shape_id = slide.add_image(
            image_path,
            (1000000, 2000000, 4000000, 2000000),
        )
        if shape_id <= 0:
            raise AssertionError("expected positive shape id for inserted image")

        pres.save(output_path)

    if not output_path.exists():
        raise AssertionError("expected output deck to exist")
