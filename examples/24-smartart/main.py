"""Demonstrate the SmartArt manipulation API.

This example demonstrates:
- Creating a presentation with SmartArt-related documentation slides
- Describing the SmartArt layout types available in gopptx
- Showing the API pattern for SmartArt operations
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation

_SMARTART_LAYOUTS = [
    "BasicBlockList",
    "VerticalBlockList",
    "HorizontalBulletList",
    "SquareAccentList",
    "PictureAccentList",
    "BasicProcess",
    "AccentProcess",
    "AlternatingFlow",
    "ContinuousBlockProcess",
    "BasicCycle",
    "TextCycle",
    "BlockCycle",
    "OrgChart",
    "Hierarchy",
    "HorizontalHierarchy",
    "BasicVenn",
    "LinearVenn",
    "StackedVenn",
    "BasicRadial",
    "BasicMatrix",
    "TitledMatrix",
    "BasicPyramid",
    "InvertedPyramid",
    "PictureStrips",
    "PictureGrid",
]

_SMARTART_CATEGORIES = {
    "List": [
        "BasicBlockList",
        "VerticalBlockList",
        "HorizontalBulletList",
        "SquareAccentList",
        "PictureAccentList",
    ],
    "Process": [
        "BasicProcess",
        "AccentProcess",
        "AlternatingFlow",
        "ContinuousBlockProcess",
    ],
    "Cycle": ["BasicCycle", "TextCycle", "BlockCycle"],
    "Hierarchy": ["OrgChart", "Hierarchy", "HorizontalHierarchy"],
    "Relationship": ["BasicVenn", "LinearVenn", "StackedVenn", "BasicRadial"],
    "Matrix": ["BasicMatrix", "TitledMatrix"],
    "Pyramid": ["BasicPyramid", "InvertedPyramid"],
    "Picture": ["PictureStrips", "PictureGrid"],
}


def _add_overview_slide(prs: Presentation) -> None:
    """Add the overview slide."""
    prs.add_bullet_slide(
        "SmartArt Full Layout Showcase",
        [
            "All currently supported SmartArt layouts in gopptx.",
            "Each category groups related diagram types.",
            "Use the Go API to create native SmartArt shapes.",
        ],
    )


def _add_category_slides(prs: Presentation) -> None:
    """Add one slide per SmartArt category."""
    for category, layouts in _SMARTART_CATEGORIES.items():
        prs.add_bullet_slide(
            f"SmartArt Category: {category}",
            [f"  - {layout}" for layout in layouts],
        )


def _add_api_reference_slide(prs: Presentation) -> None:
    """Add a slide documenting the SmartArt API."""
    prs.add_bullet_slide(
        "SmartArt API Reference",
        [
            "smartart.NewSmartArt(layout) — create a SmartArt diagram",
            ".Position(x, y).Size(cx, cy) — set position and size",
            ".WithColorStyle(uri) — apply a color style",
            ".WithQuickStyle(uri) — apply a quick style",
            ".AddItems([]string{...}) — add text nodes",
            ".AddNode(root) — add a hierarchy root node",
        ],
    )


def main() -> None:
    """Create a presentation documenting the SmartArt layout showcase."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("SmartArt Full Layout Showcase") as prs:
        _add_overview_slide(prs)
        _add_category_slides(prs)
        _add_api_reference_slide(prs)

        output_path = output_dir / "24-smartart.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(
        f"Created a SmartArt showcase with {len(_SMARTART_CATEGORIES)} category slides"
    )
    print(f"covering {len(_SMARTART_LAYOUTS)} total SmartArt layout types.")


if __name__ == "__main__":
    main()
