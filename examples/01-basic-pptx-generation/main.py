"""Create a minimal single-slide PPTX presentation.

This example demonstrates:
- Creating a new blank presentation with Presentation.new()
- Adding a single blank slide
- Saving the result to disk
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType


def main() -> None:
    """Create a basic hello-world PPTX presentation."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Hello World") as prs:
        prs.add_slide("", layout=SlideLayoutType.BLANK)
        output_path = output_dir / "01-basic-pptx-generation.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created a minimal single-slide PPTX presentation.")
    print("Demonstrated: Presentation.new(), add_slide(), save()")


if __name__ == "__main__":
    main()
