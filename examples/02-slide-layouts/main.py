"""Demonstrate all available slide layout types in gopptx.

This example demonstrates:
- Title-and-content (default) layout
- Title-only layout
- Blank layout
- Two-column layout
- Centered title layout
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType


def _add_default_slide(prs: Presentation) -> None:
    """Add a slide using the default title-and-content layout."""
    prs.add_bullet_slide(
        "Title and Content (Default)",
        [
            "This is the default layout.",
            "Used when no layout is specified.",
        ],
    )


def _add_title_only_slide(prs: Presentation) -> None:
    """Add a slide using the title-only layout."""
    prs.add_slide("Title Only Layout", layout=SlideLayoutType.TITLE_ONLY)


def _add_blank_slide(prs: Presentation) -> None:
    """Add a blank slide with no placeholders."""
    prs.add_slide("", layout=SlideLayoutType.BLANK)


def _add_two_column_slide(prs: Presentation) -> None:
    """Add a two-column layout slide with bullet content."""
    prs.add_bullet_slide(
        "Two Column Layout",
        [
            "First item in the content area.",
            "Second item in the content area.",
        ],
    )


def _add_centered_slide(prs: Presentation) -> None:
    """Add a slide with a centered title layout."""
    prs.add_slide("Centered Title Layout", layout=SlideLayoutType.CENTERED_TITLE)


def main() -> None:
    """Create a presentation demonstrating all slide layout types."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Slide Layout Types") as prs:
        _add_default_slide(prs)
        _add_title_only_slide(prs)
        _add_blank_slide(prs)
        _add_two_column_slide(prs)
        _add_centered_slide(prs)

        output_path = output_dir / "02-slide-layouts.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 5 slides, each with a different layout type:")
    print("  - Default (title + content)")
    print("  - Title Only")
    print("  - Blank")
    print("  - Two Column")
    print("  - Centered Title")


if __name__ == "__main__":
    main()
