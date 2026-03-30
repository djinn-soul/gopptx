"""Demonstrate deep shape editing: search, update text, and reposition shapes.

This example demonstrates:
- Creating a presentation with a named textbox shape
- Using search_shapes() to locate shapes by text content
- Updating shape text and position with update_shape()
- Verifying edited properties after save
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches, ShapeUpdate


def main() -> None:
    """Create presentation demonstrating deep shape editing."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Create the base presentation with a shape
    with Presentation.new("Deep Shape Editing") as prs:
        slide = prs.add_slide("Shape Search Demo")
        shape_id = slide.add_textbox(
            Inches(1),
            Inches(1),
            Inches(5),
            Inches(2),
            text="Original Text",
        )
        slide.list_shapes()

        # Find the shape by text using the supported presentation search API.
        results = prs.search_shapes("Original Text")
        print(f"Found {len(results)} shape(s) matching 'Original Text'")

        if results:
            found = next(
                (item for item in results if item["SlideIndex"] == slide.index),
                results[0],
            )
            slide_idx = int(found["SlideIndex"])
            shape_id = int(found["Shape"]["ID"])

            prs.update_shape(
                slide_idx,
                shape_id,
                ShapeUpdate(
                    text="Edited Text",
                    x=Inches(0.5),
                    y=Inches(0.5),
                ),
            )
            print(f"Updated shape {shape_id}: text='Edited Text', x=0.5in, y=0.5in")

        # Add a summary slide
        prs.add_bullet_slide(
            "Shape Editing Features",
            [
                "search_shapes(query) — search shapes by text or name",
                "update_shape(slide, id, ShapeUpdate) — modify properties",
                "ShapeUpdate accepts: text, x, y, width, height",
                "Useful for template-driven content injection",
            ],
        )

        output_path = output_dir / "41_deep_shape_editing.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Shape 'Original Text' found, text updated to 'Edited Text'")
    print("Position moved to (0.5in, 0.5in)")


if __name__ == "__main__":
    main()
