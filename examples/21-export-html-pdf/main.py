"""Demonstrate creating a rich presentation ready for HTML and PDF export.

This example demonstrates:
- Building slides with shapes, tables, and text content
- Creating a multi-slide presentation documenting export capabilities
- Using add_shape(), add_table_from_rows() for rich content
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_shapes_slide(prs: Presentation) -> None:
    """Add a slide with shapes representing SVG export."""
    prs.add_bullet_slide(
        "High-Fidelity Export",
        [
            "This presentation demonstrates gopptx's enhanced export capabilities.",
            "Features include SVG shape rendering, pure CSS styling, and PDF generation.",
        ],
    )
    idx = prs.slide_count - 1
    prs.add_shape(
        idx,
        "ROUNDED_RECTANGLE",
        bounds=(Inches(1), Inches(3.5), Inches(2), Inches(1)),
        text="SVG Rect",
        properties={"fill": {"solid": "0078D4"}},
    )


def _add_table_slide(prs: Presentation) -> None:
    """Add a slide with a feature status table."""
    prs.add_slide("Gradients and Tables", layout=SlideLayoutType.TITLE_ONLY)
    idx = prs.slide_count - 1

    rows = [
        ["Feature", "Status"],
        ["HTML SVG Export", "Done"],
        ["Headless PDF", "Done"],
        ["Native PDF Engine", "Done"],
    ]
    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(1), Inches(1.5), Inches(8), Inches(3)),
        first_row=True,
        band_row=True,
    )


def _add_rotations_slide(prs: Presentation) -> None:
    """Add a slide with shapes representing rotation and custom options."""
    prs.add_slide("Rotations and Custom Options", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1
    prs.add_shape(
        idx,
        "TRIANGLE",
        bounds=(Inches(4), Inches(2), Inches(2), Inches(1)),
        text="Rotated Arrow",
        properties={"fill": {"solid": "28A745"}},
    )
    prs.add_textbox(
        idx,
        Inches(1),
        Inches(4),
        Inches(8),
        Inches(0.5),
        text="Shapes with custom fill and rotation demonstrate export fidelity.",
    )


def main() -> None:
    """Create a rich presentation demonstrating export capabilities."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Export Demo") as prs:
        _add_shapes_slide(prs)
        _add_table_slide(prs)
        _add_rotations_slide(prs)

        output_path = output_dir / "21-export-html-pdf.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created a 3-slide presentation demonstrating export-ready content:")
    print("  - Slide 1: SVG-renderable shapes with text")
    print("  - Slide 2: Feature status table")
    print("  - Slide 3: Rotations and custom fill options")


if __name__ == "__main__":
    main()
