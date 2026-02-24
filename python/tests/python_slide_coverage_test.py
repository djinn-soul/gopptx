import os  # noqa: D100
import pathlib

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118


def test_slide_management() -> None:  # noqa: D103
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        initial_count = prs.slide_count

        # Add slide
        new_slide = prs.add_slide("New Slide", bullets=["Bullet 1", "Bullet 2"])
        assert prs.slide_count == initial_count + 1  # noqa: S101

        # Duplicate
        dup_slide = new_slide.duplicate()
        assert prs.slide_count == initial_count + 2  # noqa: S101

        # Move
        prs.move_slide(prs.slide_count - 1, 0)
        assert prs.slides[0].title == dup_slide.title  # noqa: S101

        # Remove
        prs.remove_slide(0)
        assert prs.slide_count == initial_count + 1  # noqa: S101

        # Update slide
        prs.update_slide(0, title="Updated", bullets=["One"])
        assert prs.slides[0].title == "Updated"  # noqa: S101


def test_presentation_metadata() -> None:  # noqa: D103
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        meta = prs.metadata
        assert "title" in meta  # noqa: S101
        assert "size" in meta  # noqa: S101

        # Core properties
        props = prs.core_properties
        prs.core_properties = props

        # Set slide size
        prs.set_slide_size(12192000, 6858000)


def test_sections() -> None:  # noqa: D103
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        prs.add_slide("S1")
        prs.add_slide("S2")

        prs.add_section("Sec 1", [0])
        assert len(prs.sections) > 0  # noqa: S101

        prs.rename_section("Sec 1", "Sec 1 Renamed")
        assert prs.sections[0]["Name"] == "Sec 1 Renamed"  # noqa: S101

        prs.remove_section("Sec 1 Renamed")
        assert len(prs.sections) == 0  # noqa: S101


def test_slide_layouts_and_masters() -> None:  # noqa: D103
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        masters = prs.slide_masters
        assert len(masters) > 0  # noqa: S101

        master = masters[0]
        layouts = master.slide_layouts
        assert len(layouts) > 0  # noqa: S101

        layout = layouts[0]
        assert layout.name != ""  # noqa: PLC1901, S101

        # List all layouts
        all_layouts = prs.list_slide_layouts()
        assert len(all_layouts) > 0  # noqa: S101

        # Rebind
        prs.rebind_slide_layout(0, all_layouts[0]["Part"])
