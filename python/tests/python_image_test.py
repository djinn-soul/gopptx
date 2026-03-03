import os  # noqa: D100
import pathlib
import sys

import pytest
from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118


def test_python_image(tmp_path: pathlib.Path) -> None:
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
        assert shape_id > 0

        pres.save(output_path)

    assert output_path.exists()
