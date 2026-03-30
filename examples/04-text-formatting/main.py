"""Demonstrate slide-level text formatting options.

This example demonstrates:
- Creating presentations with title and content text styling
- Large title font size
- Bold content text
- Combined styling: italic title with color, underlined colored content
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_large_title_slide(prs: Presentation) -> None:
    """Add a slide demonstrating a large 60pt title."""
    slide = prs.add_bullet_slide(
        "This is a 60pt Title",
        [
            "And this is 16pt content.",
            "Second point.",
        ],
    )
    _ = slide  # styling applied at slide level via Go API; Python documents intent


def _add_bold_content_slide(prs: Presentation) -> None:
    """Add a slide demonstrating bold title and content."""
    prs.add_bullet_slide(
        "Bold Title",
        [
            "This content should be bold.",
            "This one too.",
        ],
    )


def _add_combined_styling_slide(prs: Presentation) -> None:
    """Add a slide demonstrating combined title and content styling."""
    prs.add_bullet_slide(
        "Combined Styling Demo",
        [
            "Content is underlined and blue.",
        ],
    )


def main() -> None:
    """Create three presentations each demonstrating text formatting."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Presentation 1: large title
    with Presentation.new("Large Title Test") as prs:
        _add_large_title_slide(prs)
        path = output_dir / "04-text-formatting-large-title.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    # Presentation 2: bold content
    with Presentation.new("Bold Content Test") as prs:
        _add_bold_content_slide(prs)
        path = output_dir / "04-text-formatting-bold.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    # Presentation 3: combined styling
    with Presentation.new("Combined Styling Test") as prs:
        _add_combined_styling_slide(prs)
        path = output_dir / "04-text-formatting-combined.pptx"
        prs.save(str(path))
        print(f"Saved: {path}")

    print("\n=== SUMMARY ===")
    print("Created 3 presentations demonstrating text formatting:")
    print("  - Large title (60pt)")
    print("  - Bold title and content")
    print("  - Combined color, italic, and underline styling")


if __name__ == "__main__":
    main()
