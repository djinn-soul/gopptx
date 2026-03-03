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


def test_python_command_api(tmp_path: pathlib.Path) -> None:
    input_deck = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"
    if not input_deck.exists():
        pytest.skip("Smoke sample missing for python_command_api_test")

    with Presentation(str(input_deck)) as pres:
        initial_count = pres.slide_count
        assert initial_count > 0

        meta = pres.metadata
        assert meta is not None

        new_idx = pres.duplicate_slide(0, 1)
        assert new_idx >= 0

        pres.move_slide(0, pres.slide_count - 1)
        pres.remove_slide(1)

        output_path = tmp_path / "python_management_output.pptx"
        pres.save(output_path)

    assert output_path.exists()
