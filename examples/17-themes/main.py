"""Demonstrate applying built-in color and font themes to a presentation.

This example demonstrates:
- Applying a built-in theme (aurora, ocean, sunset, forest) with apply_theme()
- Creating a multi-slide presentation with theme-aware styling
- Using get_theme() to retrieve a named theme object
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.theme import get_theme


def _add_theme_slides(prs: Presentation, theme_name: str) -> None:
    """Add representative slides for the given theme."""
    prs.add_title_slide(f"{theme_name.title()} Theme")

    prs.add_bullet_slide(
        "Key Competitive Advantages",
        [
            "Hyper-converged edge networking",
            "Zero-latency transcoding pipelines",
            "AI-driven audience engagement",
            "Quantum-secure content distribution",
        ],
    )

    prs.add_bullet_slide(
        "Theme Details",
        [
            f"Theme name: {theme_name}",
            "Colors and fonts are applied across all slides.",
            "Customize via get_theme() and apply_theme().",
        ],
    )


def _build_themed_presentation(theme_name: str, output_path: Path) -> None:
    """Build and save a presentation using the named theme."""
    with Presentation.new(f"{theme_name.title()} Theme Demo") as prs:
        theme = get_theme(theme_name)
        prs.apply_theme(theme)
        _add_theme_slides(prs, theme_name)
        prs.save(str(output_path))
        print(f"Saved: {output_path}")


def main() -> None:
    """Create presentations for each built-in theme."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    themes = ["aurora", "ocean", "sunset", "forest"]
    for theme_name in themes:
        output_path = output_dir / f"17-themes-{theme_name}.pptx"
        _build_themed_presentation(theme_name, output_path)

    print("\n=== SUMMARY ===")
    print(f"Created {len(themes)} presentations, one for each built-in theme:")
    for t in themes:
        print(f"  - {t}")


if __name__ == "__main__":
    main()
