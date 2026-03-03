import os
import pathlib

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118


def test_notes_slide_is_none_until_created() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.slides[0]
        assert slide.notes_slide is None  # noqa: S101

        slide.notes = "hello notes"
        notes_slide = slide.notes_slide
        assert notes_slide is not None  # noqa: S101
        assert notes_slide.text == "hello notes"  # noqa: S101
