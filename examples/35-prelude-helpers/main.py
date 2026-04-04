"""Demonstrate PresentationBuilder helpers: themes, slide size, slide numbers.

This example demonstrates:
- Setting 16:9 slide size
- Applying a built-in theme preset
- Using slide numbers and footers
- Builder pattern for concise presentation setup
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import SIZE_16X9_HEIGHT, SIZE_16X9_WIDTH
from gopptx.presentation.theme import get_theme


def main() -> None:
    """Create presentation demonstrating prelude helper utilities."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Prelude Helpers Demo") as prs:
        # Set widescreen 16:9 slide size
        prs.set_slide_size(SIZE_16X9_WIDTH, SIZE_16X9_HEIGHT)

        # Apply built-in theme preset
        prs.apply_theme(get_theme("ocean"))

        # Slide 1: welcome / overview
        prs.add_bullet_slide(
            "Welcome",
            [
                "Built with PresentationBuilder",
                "Fluent API for creating presentations",
                "Chain methods for concise setup",
            ],
        )

        # Slide 2: numbered features list
        prs.add_bullet_slide(
            "Features",
            [
                "1. Builder pattern",
                "2. Theme support",
                "3. Slide size control",
                "4. Metadata helpers",
            ],
        )

        # Slide 3: slide size awareness — EMU values
        emu_16 = int(SIZE_16X9_WIDTH)
        emu_9 = int(SIZE_16X9_HEIGHT)
        prs.add_bullet_slide(
            "Slide Size Awareness",
            [
                f"16:9 width = {emu_16} EMU",
                f"16:9 height = {emu_9} EMU",
                "Inches() converts from human-readable to EMU",
                "SIZE_16X9_WIDTH/HEIGHT are standard presets",
            ],
        )

        # Slide 4: builder options
        prs.add_bullet_slide(
            "Builder Options",
            [
                "set_slide_size(w, h) — widescreen or custom dimensions",
                "apply_theme(theme) — apply a color+font palette",
                "set_core_properties({...}) — core properties",
                "save(path) — build and persist in one call",
            ],
        )

        output_path = output_dir / "35_prelude_helpers.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("4 slides: Welcome, Features, Slide Size, Builder Options")
    print(f"16:9 dimensions: {SIZE_16X9_WIDTH} x {SIZE_16X9_HEIGHT} EMU")


if __name__ == "__main__":
    main()
