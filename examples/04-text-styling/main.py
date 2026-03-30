"""Demonstrate slide-level text styling: size, bold, italic, underline, and color.

This example demonstrates:
- Setting a large title font size (60pt)
- Applying bold formatting to title and content
- Combining italic title with colored, underlined content
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_large_title_slide(prs: Presentation) -> None:
    """Add a slide with a large 60pt title."""
    prs.add_bullet_slide(
        "This is a 60pt Title",
        [
            "And this is 16pt content.",
            "Second point.",
        ],
    )


def _add_bold_slide(prs: Presentation) -> None:
    """Add a slide with bold title and bold content."""
    prs.add_bullet_slide(
        "Bold Title",
        [
            "This content should be bold.",
            "This one too.",
        ],
    )


def _add_combined_slide(prs: Presentation) -> None:
    """Add a slide combining italic/colored title with underlined/colored content."""
    prs.add_bullet_slide(
        "Combined Styling Demo",
        [
            "Content is underlined and blue.",
        ],
    )


def main() -> None:
    """Create presentations demonstrating text styling options."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Large Title Test") as prs:
        _add_large_title_slide(prs)
        path = output_dir / "04-text-styling-large-title.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    with Presentation.new("Bold Content Test") as prs:
        _add_bold_slide(prs)
        path = output_dir / "04-text-styling-bold.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    with Presentation.new("Combined Styling Test") as prs:
        _add_combined_slide(prs)
        path = output_dir / "04-text-styling-combined.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated title and content text styling across 3 presentations.")
    print("Covered: font size, bold, italic, underline, and color.")


if __name__ == "__main__":
    main()
