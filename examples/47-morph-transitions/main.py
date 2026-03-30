"""Demonstrate morph transitions with matching shape names across slides.

This example demonstrates:
- Creating two slides where a shape moves and resizes between them
- Shapes on both slides share the same name to enable PowerPoint Morph
- The Morph transition animates position, size, and color changes
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def main() -> None:
    """Create presentation demonstrating morph transitions."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Morph Transition Demo") as prs:
        # Slide 1: shape at start position (red square, top-left)
        prs.add_slide("Slide 1: Start", layout=SlideLayoutType.BLANK)
        idx1 = prs.slide_count - 1
        prs.add_shape(
            idx1,
            "RECTANGLE",
            (Inches(1), Inches(1), Inches(2), Inches(2)),
            text="Morphed-Rect",
            properties={"fill_color": "FF0000", "name": "Morphed-Rect"},
        )

        # Slide 2: same shape at end position (blue rectangle, bottom-right)
        # Morph transition animates the change between slides
        prs.add_slide("Slide 2: End", layout=SlideLayoutType.BLANK)
        idx2 = prs.slide_count - 1
        prs.add_shape(
            idx2,
            "RECTANGLE",
            (Inches(5), Inches(3), Inches(4), Inches(1)),
            text="Morphed-Rect",
            properties={"fill_color": "0000FF", "name": "Morphed-Rect"},
        )

        # Explanation slide
        prs.add_bullet_slide(
            "Morph Transition Guide",
            [
                "Both slides must have a shape with the same name",
                "PowerPoint Morph animates position, size, and color",
                "Set Morph transition on the destination slide (Slide 2)",
                "Works with rectangles, ellipses, text boxes, and images",
                "Enable in PowerPoint: Transitions -> Morph",
            ],
        )

        output_path = output_dir / "47_morph_transition.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("2 slides with matching 'Morphed-Rect' shapes for Morph transition.")
    print("Verify in PowerPoint:")
    print("  1. Go to slide 2.")
    print("  2. Apply Transitions -> Morph.")
    print("  3. The red square should animate to the blue rectangle.")


if __name__ == "__main__":
    main()
