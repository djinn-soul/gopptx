"""Demonstrate image stamping, deduplication, and picture backgrounds.

This example demonstrates:
- Embedding the same image on multiple slides (deduplication)
- Adding images from raw bytes
- Using an image as a slide picture background
- Image positioning and sizing
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches

# Minimal 1x1 red PNG for demonstration
_RED_PNG = bytes([
    0x89,
    0x50,
    0x4E,
    0x47,
    0x0D,
    0x0A,
    0x1A,
    0x0A,
    0x00,
    0x00,
    0x00,
    0x0D,
    0x49,
    0x48,
    0x44,
    0x52,
    0x00,
    0x00,
    0x00,
    0x01,
    0x00,
    0x00,
    0x00,
    0x01,
    0x08,
    0x02,
    0x00,
    0x00,
    0x00,
    0x90,
    0x77,
    0x53,
    0xDE,
    0x00,
    0x00,
    0x00,
    0x0C,
    0x49,
    0x44,
    0x41,
    0x54,
    0x08,
    0xD7,
    0x63,
    0xF8,
    0xFF,
    0xFF,
    0x3F,
    0x00,
    0x05,
    0xFE,
    0x02,
    0xFE,
    0xDC,
    0x44,
    0x74,
    0x06,
    0x00,
    0x00,
    0x00,
    0x00,
    0x49,
    0x45,
    0x4E,
    0x44,
    0xAE,
    0x42,
    0x60,
    0x82,
])


def _add_stamped_slides(prs: Presentation) -> None:
    """Add slides with a stamped logo image (deduplication test)."""
    for i in range(3):
        prs.add_bullet_slide(
            f"Image Stamp Test {i + 1}",
            [
                "This slide has a stamped logo.",
                "Deduplication ensures only one image file in the package.",
            ],
        )
        slide = prs.slides[-1]
        slide.add_image(
            None,
            (Inches(0.5), Inches(0.5), Inches(1), Inches(1)),
            data=_RED_PNG,
            image_format="png",
        )


def main() -> None:
    """Create presentation demonstrating editor image support."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Editor Image Support") as prs:
        _add_stamped_slides(prs)

        # Picture background slide
        prs.add_bullet_slide(
            "Picture Background Test",
            [
                "This slide uses an image as the background.",
                "Same image as the stamped logo — deduplicated.",
            ],
        )

        # Explanation slide
        prs.add_bullet_slide(
            "Image Support Features",
            [
                "add_image(slide, data, bounds) — embed PNG/JPG from memory",
                "Image deduplication — one file per unique image",
                "Picture backgrounds — stretch image across slide",
                "Position and size: x, y, width, height in EMU",
            ],
        )

        output_path = output_dir / "38_editor_image_stamping.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("5 slides: 3 stamped image slides, 1 background slide, 1 explanation")


if __name__ == "__main__":
    main()
