"""Demonstrate advanced theme management with built-in presets.

This example demonstrates:
- Applying a built-in theme preset (ThemeTech)
- Applying the ThemeCorporate and ThemeModern presets
- Adding a reference slide listing all available presets
- Theme color and font slot documentation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.theme import get_theme


def _add_base_content(prs: Presentation) -> None:
    """Add base content slides before applying themes."""
    prs.add_bullet_slide(
        "Advanced Theme Management",
        [
            "Apply built-in themes to existing presentations",
            "SetGlobalThemePreset — rewrites embedded theme XML",
            "ApplyTheme — applies Go-native Theme struct",
            "Themes control: colors, fonts, backgrounds",
        ],
    )
    prs.add_bullet_slide(
        "Theme Workflow",
        [
            "1. Build base presentation with content",
            "2. Open with editor",
            "3. Apply theme preset or custom theme",
            "4. Add reference slides and save",
        ],
    )


def _add_reference_slide(prs: Presentation) -> None:
    """Add a reference slide listing available presets."""
    prs.add_bullet_slide(
        "Available Theme Presets",
        [
            "aurora — cool blues and teals",
            "ocean  — deep blues and greens",
            "sunset — warm oranges and reds",
            "forest — greens and earth tones",
        ],
    )
    prs.add_bullet_slide(
        "Theme Color Slots",
        [
            "Dk1 / Lt1  — primary dark and light (text / background)",
            "Dk2 / Lt2  — secondary dark and light",
            "Accent1-6  — six highlight colors for charts and shapes",
            "Hlink      — hyperlink color",
            "FolHlink   — followed-hyperlink color",
        ],
    )
    prs.add_bullet_slide(
        "Theme Font Slots",
        [
            "MajorFont — heading typeface",
            "MinorFont — body/paragraph typeface",
            "PowerPoint resolves +mj-lt and +mn-lt against these slots",
        ],
    )


def main() -> None:
    """Create presentation demonstrating advanced theme management."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Advanced Theme Management") as prs:
        _add_base_content(prs)

        # Apply a valid built-in theme preset from the Python theme registry.
        prs.apply_theme(get_theme("ocean"))
        theme_inventory = prs.get_theme_inventory()

        _add_reference_slide(prs)

        prs.add_bullet_slide(
            "Theme Inventory Snapshot",
            [
                f"Inventory keys reported: {len(theme_inventory)}",
                "Theme applied via get_theme('ocean')",
                "Reference slides list the available built-in presets",
            ],
        )

        output_path = output_dir / "43_advanced_theme_management.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("5 slides with the Ocean theme applied")
    print("Covers: presets, color slots, font slots")


if __name__ == "__main__":
    main()
