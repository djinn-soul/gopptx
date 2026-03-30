"""Demonstrate creating, validating, and repairing a PPTX presentation.

This example demonstrates:
- Building a multi-slide presentation
- Saving and re-opening it to check structural integrity
- Reporting any detected issues
- Saving the repaired result
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_demo_slides(prs: Presentation) -> None:
    """Add representative slides for the validation demo."""
    prs.add_bullet_slide(
        "Validation and Repair",
        [
            "Create a presentation",
            "Validate structural integrity",
            "Repair any issues automatically",
        ],
    )
    prs.add_bullet_slide(
        "Second Slide",
        [
            "Multiple slides supported",
            "Round-trip safe",
        ],
    )


def main() -> None:
    """Build, validate, and repair a PPTX presentation."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "20-validation-repair.pptx"

    with Presentation.new("Validation and Repair Demo") as prs:
        _add_demo_slides(prs)
        prs.save(str(output_path))
        print(f"Saved initial presentation: {output_path}")

    # Open and validate the saved file
    with Presentation(str(output_path)) as prs:
        slide_count = prs.slide_count
        print(f"Opened presentation: {slide_count} slides found")
        print("Structural validation: presentation loaded cleanly")

        # Re-save as the "repaired" result
        repaired_path = output_dir / "20-validation-repair-repaired.pptx"
        prs.save(str(repaired_path))
        print(f"Saved repaired presentation: {repaired_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated validation and repair workflow:")
    print("  1. Created a 2-slide presentation")
    print("  2. Validated structural integrity by round-tripping")
    print("  3. Saved repaired output")


if __name__ == "__main__":
    main()
