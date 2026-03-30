"""Demonstrate all available bullet list styles.

This example demonstrates:
- Standard bullet points
- Numbered (ordered) lists
- Lettered lists
- Sub-bullets at multiple indent levels
- Mixed bullet styles on one slide
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_standard_bullets(prs: Presentation) -> None:
    """Add a slide with standard bullet points."""
    prs.add_bullet_slide(
        "Standard Bullet Points",
        [
            "First standard bullet item.",
            "Second standard bullet item.",
            "Third standard bullet item.",
        ],
    )


def _add_numbered_list(prs: Presentation) -> None:
    """Add a slide with a numbered list."""
    prs.add_bullet_slide(
        "Numbered List",
        [
            "1. First numbered item.",
            "2. Second numbered item.",
            "3. Third numbered item.",
        ],
    )


def _add_lettered_list(prs: Presentation) -> None:
    """Add a slide with a lettered list."""
    prs.add_bullet_slide(
        "Lettered List",
        [
            "a. Item alpha.",
            "b. Item beta.",
            "c. Item gamma.",
        ],
    )


def _add_sub_bullets(prs: Presentation) -> None:
    """Add a slide with sub-bullets at multiple indent levels."""
    prs.add_bullet_slide(
        "Sub-Bullet Levels",
        [
            "Top-level bullet (level 0).",
            "  Sub-bullet at level 1.",
            "    Sub-bullet at level 2.",
            "Back to top-level.",
            "  Another level 1 sub-bullet.",
        ],
    )


def _add_mixed_styles(prs: Presentation) -> None:
    """Add a slide mixing different bullet styles."""
    prs.add_bullet_slide(
        "Mixed Bullet Styles",
        [
            "Standard bullet item.",
            "1. Numbered item follows.",
            "a. Lettered item next.",
            "  Indented sub-point.",
        ],
    )


def main() -> None:
    """Create a presentation demonstrating all bullet list styles."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Bullet Styles Demo") as prs:
        _add_standard_bullets(prs)
        _add_numbered_list(prs)
        _add_lettered_list(prs)
        _add_sub_bullets(prs)
        _add_mixed_styles(prs)

        output_path = output_dir / "05-bullet-styles.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 5 slides demonstrating bullet list styles:")
    print("  - Standard bullets")
    print("  - Numbered list")
    print("  - Lettered list")
    print("  - Sub-bullets at multiple levels")
    print("  - Mixed styles")


if __name__ == "__main__":
    main()
