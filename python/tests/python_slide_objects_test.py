"""Object-proxy smoke tests for Python slide API."""

import pathlib
import sys

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
