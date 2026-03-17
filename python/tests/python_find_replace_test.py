import pathlib

import pytest
from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / "../..").resolve()


def test_find_and_replace() -> None:
    """Test find_and_replace slide operation via the Python wrapper."""
    # Ensure smoke sample exists
    input_deck = project_root / "examples/assets/01/01_basic_pptx.pptx"
    if not input_deck.exists():
        pytest.skip(f"Smoke sample not found: {input_deck}")

    with Presentation(str(input_deck)) as pres:
        slide = pres.add_slide("Replace Flow Test")
        pres.add_shape(slide.index, "rect", (100, 100, 200, 100), text="Replace Me Now")

        # Test replace
        expected_replacements = 1
        replacements = pres.find_and_replace("Replace Me Now", "I am replaced")
        if replacements != expected_replacements:
            pytest.fail(
                f"Expected {expected_replacements} replacement, got {replacements}"
            )

        shapes = pres.list_shapes(slide.index)
        found = False
        for s in shapes:
            if s.get("Text") == "I am replaced":
                found = True
                break
        if not found:
            pytest.fail("Shape text was not replaced")

        pres.add_shape(slide.index, "rect", (100, 250, 200, 100), text="First target")
        pres.add_shape(slide.index, "rect", (100, 400, 200, 100), text="Second target")

        expected_match = 2
        replacements = pres.find_and_replace("target", "match")
        if replacements != expected_match:
            pytest.fail(f"Expected {expected_match} replacements, got {replacements}")
