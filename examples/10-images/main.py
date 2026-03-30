"""Demonstrate image embedding from raw bytes and from a file path.

This example demonstrates:
- Embedding an image from raw bytes with add_image_from_bytes()
- Embedding an image loaded from a file path with add_image()
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches

# Minimal 1x1 white PNG
_WHITE_PNG = bytes([
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


def _add_image_from_bytes_slide(prs: Presentation) -> None:
    """Add a slide that embeds an image from raw bytes."""
    prs.add_bullet_slide(
        "Image from Bytes",
        [
            "This slide embeds a 1x1 white PNG supplied as raw bytes.",
            "No file is required at runtime.",
        ],
    )
    slide = prs.slides[prs.slide_count - 1]
    slide.add_image(
        None,
        (Inches(1), Inches(3), Inches(4), Inches(3)),
        data=_WHITE_PNG,
        image_format="png",
    )


def _add_image_from_file_slide(prs: Presentation, image_path: str) -> None:
    """Add a slide that embeds an image from a file path."""
    prs.add_bullet_slide(
        "Image from File Path",
        [
            "This slide loads the same PNG from a temporary file on disk.",
            "Use add_image(slide_idx, path, x, y, w, h) for file-based images.",
        ],
    )
    slide = prs.slides[prs.slide_count - 1]
    slide.add_image(
        image_path,
        (Inches(1), Inches(3), Inches(4), Inches(3)),
    )


def main() -> None:
    """Create a presentation demonstrating image embedding."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
        tmp.write(_WHITE_PNG)
        tmp_path = tmp.name

    try:
        with Presentation.new("Image Embedding Demo") as prs:
            _add_image_from_bytes_slide(prs)
            _add_image_from_file_slide(prs, tmp_path)

            output_path = output_dir / "10-images.pptx"
            prs.save(str(output_path))
            print(f"Saved: {output_path}")
    finally:
        Path(tmp_path).unlink(missing_ok=True)

    print("\n=== SUMMARY ===")
    print("Created 2 slides demonstrating image embedding:")
    print("  - Image embedded from raw bytes")
    print("  - Image loaded from a temporary file on disk")


if __name__ == "__main__":
    main()
