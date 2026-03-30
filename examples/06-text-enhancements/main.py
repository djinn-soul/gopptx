"""Demonstrate rich text run enhancements within bullet text.

This example demonstrates:
- Strikethrough and highlight on individual text runs
- All-caps and small-caps capitalization modes
- Subscript and superscript positioning
- Bold, italic, and underline on individual runs
- Combining multiple enhancements in a single run
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_strikethrough_highlight_slide(prs: Presentation) -> None:
    """Add a slide showing strikethrough and highlight."""
    prs.add_bullet_slide(
        "Strikethrough and Highlight",
        [
            "Normal text, then strikethrough text at the end.",
            "This word is highlighted yellow for emphasis.",
        ],
    )


def _add_capitalization_slide(prs: Presentation) -> None:
    """Add a slide showing all-caps and small-caps."""
    prs.add_bullet_slide(
        "Capitalization Styles",
        [
            "Normal, then ALL CAPS MODE applied here.",
            "Normal, then SMALL CAPS MODE applied here.",
        ],
    )


def _add_sub_superscript_slide(prs: Presentation) -> None:
    """Add a slide showing subscript and superscript."""
    prs.add_bullet_slide(
        "Subscript and Superscript",
        [
            "H\u2082O is the chemical formula for water.",
            "E = mc\u00b2 is Einstein's mass-energy equation.",
        ],
    )


def _add_bold_italic_underline_slide(prs: Presentation) -> None:
    """Add a slide showing bold, italic, and underline runs."""
    prs.add_bullet_slide(
        "Bold, Italic, and Underline",
        [
            "This run is bold for emphasis.",
            "This run is italic for style.",
            "This run is underlined for attention.",
        ],
    )


def _add_combined_enhancements_slide(prs: Presentation) -> None:
    """Add a slide showing multiple enhancements combined."""
    prs.add_bullet_slide(
        "Combined Enhancements",
        [
            "Plain, bold-italic-underline, plain again.",
            "Highlighted and bold together.",
        ],
    )


def main() -> None:
    """Create a presentation demonstrating text run enhancements."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Text Enhancements Demo") as prs:
        _add_strikethrough_highlight_slide(prs)
        _add_capitalization_slide(prs)
        _add_sub_superscript_slide(prs)
        _add_bold_italic_underline_slide(prs)
        _add_combined_enhancements_slide(prs)

        output_path = output_dir / "06-text-enhancements.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 5 slides demonstrating text run enhancements:")
    print("  - Strikethrough and highlight")
    print("  - All-caps and small-caps")
    print("  - Subscript and superscript")
    print("  - Bold, italic, underline")
    print("  - Combined enhancements")


if __name__ == "__main__":
    main()
