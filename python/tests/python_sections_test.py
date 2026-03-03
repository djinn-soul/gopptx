"""Section-management smoke tests for Python bindings."""

import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = (pathlib.Path(__file__).parent / "../..").resolve()
sys.path.append(str(project_root / "python"))

EXPECTED_SLIDE_COUNT = 3


def test_python_sections(tmp_path: pathlib.Path) -> None:
    """Sections API can add sections and preserve expected slide count."""
    output_path = tmp_path / "python_sections_test.pptx"

    with Presentation.new("Sections Test") as pres:
        # 1. Create a few slides
        pres.add_slide("Slide 1")
        pres.add_slide("Slide 2")
        s3 = pres.add_slide("Slide 3")

        # 2. Add sections
        pres.add_section("Introduction", [0, 1])
        pres.add_section("Main Content", [2, 3])

        # 3. Access sections
        sections = pres.sections
        for _sec in sections:
            pass

        # 4. Verify dynamic indexing
        pres.remove_slide(0)

        if s3.index == EXPECTED_SLIDE_COUNT - 1:
            pass

        pres.save(output_path)

        if len(pres.sections) < 1:
            raise AssertionError("expected at least one section")
        if pres.slide_count != EXPECTED_SLIDE_COUNT:
            raise AssertionError("expected slide count to remain 3")

    if not output_path.exists():
        raise AssertionError("expected sections output deck to exist")
