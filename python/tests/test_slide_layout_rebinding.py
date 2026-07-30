"""Test cross-master slide layout rebinding (Issue #1109)."""

import pathlib

from gopptx import Presentation


def test_slide_layout_rebinding(tmp_path: pathlib.Path) -> None:
    """slide.rebind_layout(layout_part) rebinds slide relationship to different layout across slide masters (Issue #1109)."""
    output_path = tmp_path / "layout_rebind_test.pptx"

    with Presentation.new(title="Layout Rebind Test Deck") as pres:
        slide = pres.slides[0]

        # List layouts
        layouts = pres.list_slide_layouts()
        assert len(layouts) > 0

        # Choose a target layout name or part
        target_layout = layouts[-1]["name"]

        # Rebind slide layout via slide proxy
        slide.rebind_layout(target_layout)

        pres.save(output_path)

    # Re-open presentation and verify slide saves without corruption
    with Presentation(output_path) as pres:
        assert pres.slide_count == 1
