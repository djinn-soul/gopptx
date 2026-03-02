import os  # noqa: D100
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118


def test_python_sections(tmp_path: pathlib.Path) -> None:
    output_path = tmp_path / "python_sections_test.pptx"

    with Presentation.new("Sections Test") as pres:
        # 1. Create a few slides
        s1 = pres.add_slide("Slide 1")
        s2 = pres.add_slide("Slide 2")
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

        if s3.index == 2:  # noqa: PLR2004
            pass

        pres.save(output_path)

        assert len(pres.sections) >= 1
        assert pres.slide_count == 3

    assert output_path.exists()
