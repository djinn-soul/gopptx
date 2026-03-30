"""Demonstrate layout helpers including unit conversion and precise positioning.

This example demonstrates:
- Using Inches() for precise EMU-based coordinate calculations
- Positioning shapes with exact inch measurements
- Documenting unit conversion values (Inches, Centimeters, Points)
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_unit_conversion_slide(prs: Presentation) -> None:
    """Add a slide showing unit conversion helper values."""
    prs.add_bullet_slide(
        "Unit Conversion Helpers",
        [
            f"Inches(1)          = {Inches(1)} EMU",
            f"Inches(2)          = {Inches(2)} EMU",
            "Centimeters(2.54)  = ~914400 EMU  (= 1 inch)",
            "Points(72)         = ~914400 EMU  (= 1 inch)",
            "All coordinates passed to gopptx use Inches() values.",
        ],
    )


def _add_positioned_shapes_slide(prs: Presentation) -> None:
    """Add a slide with shapes placed at precise inch coordinates."""
    prs.add_slide("Precise Positioning with Inches()", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(1), Inches(1), Inches(2), Inches(1)),
        text="(1in, 1in) 2x1in",
        properties={"fill": "4472C4"},
    )
    prs.add_shape(
        idx,
        "ROUNDED_RECTANGLE",
        bounds=(Inches(4), Inches(2), Inches(2), Inches(1)),
        text="(4in, 2in) 2x1in",
        properties={"fill": "ED7D31"},
    )


def _add_emu_helpers_slide(prs: Presentation) -> None:
    """Add a slide demonstrating value retrieval from inch measurements."""
    one_inch = Inches(1)
    prs.add_bullet_slide(
        "Retrieving Values from Inches()",
        [
            f"Inches(1) value = {one_inch}",
            "Pass Inches(n) directly to x, y, w, h parameters.",
            "1 inch = 914,400 EMU (English Metric Units).",
            "Use consistent units across all shape positioning calls.",
        ],
    )


def main() -> None:
    """Create a presentation demonstrating layout helper usage."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Layout Helpers Demo") as prs:
        _add_unit_conversion_slide(prs)
        _add_positioned_shapes_slide(prs)
        _add_emu_helpers_slide(prs)

        output_path = output_dir / "18-layout-helpers.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 3 slides demonstrating layout helpers:")
    print("  - Unit conversion reference (Inches, Centimeters, Points)")
    print("  - Shapes positioned with precise Inches() coordinates")
    print("  - EMU value retrieval documentation")


if __name__ == "__main__":
    main()
