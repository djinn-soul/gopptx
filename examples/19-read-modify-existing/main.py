"""Demonstrate opening and modifying an existing PPTX file.

This example demonstrates:
- Creating a base presentation and saving it
- Opening the saved file with Presentation.open()
- Adding new slides to the opened presentation
- Moving and duplicating slides
- Saving the modified presentation
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation


def _create_base_presentation(path: str) -> None:
    """Create a simple base presentation and save it."""
    with Presentation.new("Base Presentation") as prs:
        prs.add_bullet_slide(
            "Original Slide 1",
            ["This slide was in the original file.", "It will be preserved."],
        )
        prs.add_bullet_slide(
            "Original Slide 2",
            ["Second original slide.", "More content here."],
        )
        prs.save(path)


def _modify_presentation(base_path: str, output_path: str) -> None:
    """Open the base presentation and add/modify content."""
    with Presentation(base_path) as prs:
        print(f"  Opened presentation with {prs.slide_count} slides")

        # Add a new slide
        prs.add_bullet_slide(
            "Newly Added Slide",
            [
                "This slide was added after opening the file.",
                "Demonstrates read-modify-save round-trip.",
            ],
        )

        # Duplicate the first slide
        prs.duplicate_slide(0)

        print(f"  After modifications: {prs.slide_count} slides")
        prs.save(output_path)


def main() -> None:
    """Demonstrate opening and modifying an existing PPTX presentation."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with tempfile.NamedTemporaryFile(suffix=".pptx", delete=False) as tmp:
        base_path = tmp.name

    try:
        print("Creating base presentation...")
        _create_base_presentation(base_path)

        output_path = output_dir / "19-read-modify-existing.pptx"
        print("Opening and modifying presentation...")
        _modify_presentation(base_path, str(output_path))

        print(f"Saved: {output_path}")
    finally:
        Path(base_path).unlink(missing_ok=True)

    print("\n=== SUMMARY ===")
    print("Demonstrated the read-modify-save workflow:")
    print("  1. Created a 2-slide base presentation")
    print("  2. Opened it with Presentation(path)")
    print("  3. Added a new slide and duplicated the first")
    print("  4. Saved the modified result")


if __name__ == "__main__":
    main()
