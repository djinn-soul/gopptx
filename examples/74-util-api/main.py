"""Demonstrate unit conversion helpers and the Inches() utility.

This example demonstrates:
- Inches() for positioning and sizing shapes
- Placement of shapes using Inches() for different coordinates
- get_slide_size() to read the current slide dimensions
- set_slide_size() to change the presentation canvas size
- The SIZE_16X9_WIDTH / SIZE_16X9_HEIGHT constants
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import (
    SIZE_16X9_HEIGHT,
    SIZE_16X9_WIDTH,
    ShapeType,
)
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_unit_info_slide(prs: Presentation) -> None:
    """Slide listing unit conversion facts."""
    prs.add_bullet_slide(
        "Unit Conversion Reference",
        [
            "Inches(1)       = 914400 EMU",
            "Inches(2.54 cm) = 914400 EMU (1 inch)",
            "1 point         = 12700 EMU",
            "Inches() returns an integer EMU value for OOXML",
            f"SIZE_16X9_WIDTH  = {SIZE_16X9_WIDTH} EMU",
            f"SIZE_16X9_HEIGHT = {SIZE_16X9_HEIGHT} EMU",
        ],
    )


def _add_shapes_with_units_slide(prs: Presentation) -> None:
    """Shapes placed with Inches() for each coordinate system."""
    idx = prs.add_slide(
        "Shapes Placed with Inches()", layout=SlideLayoutType.TITLE_ONLY
    )
    placements = [
        (
            "Inches - top-left",
            "4472C4",
            Inches(0.5),
            Inches(1.5),
            Inches(2.5),
            Inches(1.0),
        ),
        (
            "Inches - centre",
            "C0504D",
            Inches(3.5),
            Inches(1.5),
            Inches(2.5),
            Inches(1.0),
        ),
        (
            "Inches - right",
            "9BBB59",
            Inches(6.5),
            Inches(1.5),
            Inches(2.5),
            Inches(1.0),
        ),
        (
            "Inches - row 2",
            "F79646",
            Inches(0.5),
            Inches(3.0),
            Inches(2.5),
            Inches(1.0),
        ),
    ]
    for label, color, x, y, w, h in placements:
        prs.add_shape(
            idx,
            ShapeType.RECTANGLE,
            (x, y, w, h),
            text=label,
            properties={"fill_color": color},
        )


def _add_slide_size_slide(prs: Presentation) -> None:
    """Demonstrate reading the slide size."""
    size = prs.get_slide_size()
    prs.add_bullet_slide(
        "Slide Size API",
        [
            f"prs.get_slide_size() = {size}",
            "prs.set_slide_size(width_emu, height_emu) changes the canvas",
            f"SIZE_16X9_WIDTH  = {SIZE_16X9_WIDTH}",
            f"SIZE_16X9_HEIGHT = {SIZE_16X9_HEIGHT}",
            "from gopptx.constants import SIZE_16X9_WIDTH, SIZE_16X9_HEIGHT",
        ],
    )


def main() -> None:
    """Create presentation demonstrating Inches() and unit conversion helpers."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "74_util_api.pptx"

    with Presentation.new("Util API Demo") as prs:
        _add_unit_info_slide(prs)
        _add_shapes_with_units_slide(prs)
        _add_slide_size_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: Inches(), SIZE_16X9_WIDTH/HEIGHT, get_slide_size(),")
    print("  shape placement with consistent Inches() coordinates")


if __name__ == "__main__":
    main()
