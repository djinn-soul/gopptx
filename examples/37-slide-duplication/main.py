"""Demonstrate slide duplication, reordering, and removal operations.

This example demonstrates:
- Duplicating a slide to a target position
- Moving slides to different indices
- Using duplicate_slide and move_slide on an open presentation
- Final slide order after a sequence of operations
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create presentation demonstrating slide duplication and reordering."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Slide Duplication Base") as prs:
        prs.add_bullet_slide(
            "Original Slide A",
            ["This slide will be duplicated."],
        )
        prs.add_bullet_slide(
            "Original Slide B",
            ["This slide will stay as is."],
        )
        prs.add_bullet_slide(
            "Original Slide C",
            ["This slide will be moved to the beginning."],
        )

        # Duplicate Slide A (index 0) — places copy after index 2
        prs.duplicate_slide(0)  # [A, B, C, A(copy)]

        # Move Slide C (index 2) to the beginning
        prs.move_slide(2, 0)  # [C, A, B, A(copy)]

        # Duplicate Slide B (now at index 2)
        prs.duplicate_slide(2)  # [C, A, B, B(copy), A(copy)]

        # Add a summary slide describing the final order
        prs.add_bullet_slide(
            "Final Slide Order",
            [
                "Operations performed:",
                "1. Duplicate Slide A -> appended as copy",
                "2. Move Slide C to index 0 (front)",
                "3. Duplicate Slide B -> inserted after original",
                "Final order: C, A, B, B(copy), A(copy)",
            ],
        )

        output_path = output_dir / "37_slide_duplication.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Slide operations: duplicate x2, move x1")
    print("Final slide count: 6 (3 originals + 2 duplicates + 1 summary)")


if __name__ == "__main__":
    main()
