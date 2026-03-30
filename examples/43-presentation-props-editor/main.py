"""Demonstrate editing presentation properties: theme, slide size, and core metadata.

This example demonstrates:
- Applying a corporate theme to an existing presentation
- Setting slide size to 16:9 widescreen
- Writing all core document properties in one call
- Brand reskin workflow: swap theme + update metadata
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import SIZE_16X9_HEIGHT, SIZE_16X9_WIDTH
from gopptx.presentation.theme import get_theme


def _add_props_demo_slides(prs: Presentation) -> None:
    """Add slides explaining presentation properties."""
    prs.add_bullet_slide(
        "Presentation Properties",
        [
            "Theme — color palette and font scheme",
            "Slide Size — 16:9, 4:3, or custom dimensions",
            "Core Properties — title, author, subject, keywords",
            "All editable after presentation creation",
        ],
    )
    prs.add_bullet_slide(
        "Core Properties",
        [
            "title — shown in file browsers and title bars",
            "subject — topical category",
            "creator / author — attribution",
            "description — longer summary",
            "keywords — for search indexing",
            "lastModifiedBy — audit trail",
        ],
    )
    prs.add_bullet_slide(
        "Brand Reskin Workflow",
        [
            "1. Open existing presentation",
            "2. apply_theme() — swap color/font palette",
            "3. set_slide_size() — resize to target format",
            "4. set_metadata() — update title and author",
            "5. save() — write updated file",
        ],
    )


def main() -> None:
    """Create presentation demonstrating presentation props editor."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Presentation Props Base") as prs:
        _add_props_demo_slides(prs)

        # Apply corporate theme
        prs.apply_theme(get_theme("corporate"))

        # Set 16:9 slide size
        prs.set_slide_size(SIZE_16X9_WIDTH, SIZE_16X9_HEIGHT)

        # Update all core properties
        prs.set_metadata(
            title="Presentation Properties Example",
            author="gopptx example",
        )

        output_path = output_dir / "43_presentation_props_editor.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    # Brand reskin variant
    with Presentation.new("Brand Reskin Base") as prs:
        prs.add_bullet_slide(
            "Brand Identity",
            [
                "Theme swapped to 'dark' preset",
                "Updated author and title",
                "16:9 widescreen format",
            ],
        )
        prs.apply_theme(get_theme("dark"))
        prs.set_slide_size(SIZE_16X9_WIDTH, SIZE_16X9_HEIGHT)
        prs.set_metadata(title="Brand Reskin Demo", author="Design Team")

        reskin_path = output_dir / "43_brand_reskin_theme_swap.pptx"
        prs.save(str(reskin_path))
        print(f"Saved: {reskin_path}")

    print("\n=== SUMMARY ===")
    print("2 files: props editor demo + brand reskin")
    print("Theme: corporate / dark | Size: 16:9")


if __name__ == "__main__":
    main()
