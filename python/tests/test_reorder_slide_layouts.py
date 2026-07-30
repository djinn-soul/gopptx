"""Test reordering slide layouts within slide masters (Issue #1080)."""

import pathlib
import zipfile

from gopptx import Presentation


def test_reorder_slide_layouts(tmp_path: pathlib.Path) -> None:
    """pres.reorder_slide_layouts(layout_parts) reorders custom slide layouts in master (Issue #1080)."""
    output_path = tmp_path / "reorder_layouts_test.pptx"

    with Presentation.new(title="Layout Reorder Test Deck") as pres:
        # Get current layout order
        layouts_before = pres.list_slide_layouts()
        assert len(layouts_before) > 1

        names_before = [layout["name"] for layout in layouts_before]

        # Reverse layout order
        reversed_names = list(reversed(names_before))
        pres.reorder_slide_layouts(reversed_names)

        pres.save(output_path)

    # Re-open presentation and verify new layout order
    with Presentation(output_path) as pres:
        layouts_after = pres.list_slide_layouts()
        names_after = [layout["name"] for layout in layouts_after]
        assert names_after == reversed_names


def test_partial_reorder_preserves_omitted_layouts(tmp_path: pathlib.Path) -> None:
    """A partial order moves named layouts without deleting the omitted layouts."""
    output_path = tmp_path / "partial_reorder_layouts_test.pptx"

    with Presentation.new(title="Partial Layout Reorder Test Deck") as pres:
        layouts_before = pres.list_slide_layouts()
        names_before = [layout["name"] for layout in layouts_before]
        assert len(names_before) > 2
        requested_first = names_before[-1]
        master_part = layouts_before[0]["master_part"]
        pres.reorder_slide_layouts([requested_first])
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as archive:
        master_xml = archive.read(master_part).decode()
    assert master_xml.count("<p:sldLayoutId ") == len(names_before)

    with Presentation(output_path) as pres:
        names_after = [layout["name"] for layout in pres.list_slide_layouts()]

    assert names_after == [requested_first, *names_before[:-1]]
