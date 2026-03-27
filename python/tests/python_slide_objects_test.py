"""Object-proxy smoke tests for Python slide API."""

import pathlib
import sys
import zipfile

from gopptx import SHAPE_ROUNDED_RECTANGLE, Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = (pathlib.Path(__file__).parent / "../..").resolve()
sys.path.append(str(project_root / "python"))

MIN_SLIDE_COUNT = 2


def test_python_slide_objects(tmp_path: pathlib.Path) -> None:
    """Slide proxy helpers can update, duplicate, and save decks."""
    output_path = tmp_path / "python_slide_objects.pptx"

    with Presentation.new("Object Oriented Test") as pres:
        # 1. Access first slide
        slide = pres.slides[0]

        # 2. Update title via property
        slide.title = "Enhanced via Python"

        # 3. Add shape via slide object
        shape_id = slide.add_shape(
            SHAPE_ROUNDED_RECTANGLE,
            (1000000, 1000000, 3000000, 1500000),
            text="I am a proxy",
        )
        if shape_id <= 0:
            raise AssertionError("expected positive shape id")

        # 4. Set notes via property
        slide.notes = "Proxy objects make API much cleaner."

        # 5. Duplicate slide
        new_slide = slide.duplicate()
        new_slide.title = "I am the copy"

        # 6. Verify newly exposed master/layout object proxies
        masters = pres.slide_masters
        if masters:
            first_master = masters[0]
            master_layouts = first_master.slide_layouts
            if master_layouts:
                _ = master_layouts[0]

        # 7. Save
        pres.save(output_path)
        if pres.slide_count < MIN_SLIDE_COUNT:
            raise AssertionError("expected at least two slides after duplication")

    if not output_path.exists():
        raise AssertionError("expected object-proxy output deck to exist")


def test_python_slide_proxy_cache_invalidation() -> None:
    """Cached slide proxies are rebuilt when slide structure changes."""
    with Presentation.new("Slide Cache Test") as pres:
        initial_slides = pres.slides
        if len(initial_slides) != 1:
            raise AssertionError(
                f"expected one initial slide, got {len(initial_slides)}"
            )

        second_slide = pres.add_slide("Second")
        refreshed_slides = pres.slides

        if len(refreshed_slides) != 2:
            raise AssertionError(
                f"expected two slides after add_slide, got {len(refreshed_slides)}"
            )
        if refreshed_slides[1].slide_id != second_slide.slide_id:
            raise AssertionError("expected refreshed slide cache to expose new slide")


def test_slides_find_by_slide_id() -> None:
    """Slides collection can resolve one slide by stable slide ID."""
    with Presentation.new("Slide ID Lookup") as pres:
        first_slide = pres.slides[0]
        second_slide = pres.add_slide("Second", layout="blank")
        second_slide.add_textbox(0.5, 0.5, 4.0, 0.5, text="Second")

        found = pres.slides.find_by_slide_id(second_slide.slide_id)
        if found is None:
            raise AssertionError("expected to find slide by stable slide ID")
        if found.slide_id != second_slide.slide_id:
            raise AssertionError("find_by_slide_id returned wrong slide")
        if pres.slides.find_by_slide_id(999999999) is not None:
            raise AssertionError("find_by_slide_id should return None for unknown ID")
        if first_slide.slide_id == second_slide.slide_id:
            raise AssertionError("expected unique slide IDs")


def test_slide_transition_disable_advance_on_click(tmp_path: pathlib.Path) -> None:
    """set_transition can disable click-advance and emit advClick=0."""
    output_path = tmp_path / "transition_disable_click.pptx"

    with Presentation.new("Transition click-disable") as pres:
        slide = pres.slides[0]
        slide.set_transition(
            "fade",
            advance_ms=20000,
            disable_advance_on_click=True,
        )
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
    if 'advTm="20000"' not in slide_xml:
        raise AssertionError("expected advTm to be written")
    if 'advClick="0"' not in slide_xml:
        raise AssertionError("expected advClick=0 to be written")
