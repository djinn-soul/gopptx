"""Demonstrate advanced image source options: bytes and base64.

This example demonstrates:
- Embedding an image from raw bytes with add_image_from_bytes()
- Embedding an image from a base64-encoded string (decoded to bytes first)
"""

from __future__ import annotations

import base64
from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches

# Minimal 1x1 white PNG as raw bytes
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


def _add_bytes_slide(prs: Presentation) -> None:
    """Add a slide with an image sourced from raw bytes."""
    prs.add_bullet_slide(
        "Image from Bytes",
        [
            "add_image(None, bounds, data=bytes, image_format='png')",
            "Pass any bytes PNG/JPEG directly — no file needed.",
        ],
    )
    slide = prs.slides[prs.slide_count - 1]
    slide.add_image(
        None,
        (Inches(1), Inches(2.5), Inches(3), Inches(2)),
        data=_WHITE_PNG,
        image_format="png",
    )


def _add_base64_slide(prs: Presentation) -> None:
    """Add a slide with an image sourced from a base64 string."""
    b64_str = base64.b64encode(_WHITE_PNG).decode("ascii")
    image_bytes = base64.b64decode(b64_str)

    prs.add_bullet_slide(
        "Image from Base64",
        [
            "Decode a base64 string and pass to add_image().",
            "Useful when images arrive via JSON APIs or config files.",
        ],
    )
    slide = prs.slides[prs.slide_count - 1]
    slide.add_image(
        None,
        (Inches(1), Inches(2.5), Inches(3), Inches(2)),
        data=image_bytes,
        image_format="png",
    )


def main() -> None:
    """Create a presentation demonstrating advanced image source options."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Advanced Image Sources Demo") as prs:
        _add_bytes_slide(prs)
        _add_base64_slide(prs)

        output_path = output_dir / "11-image-advanced-sources.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 2 slides demonstrating advanced image sources:")
    print("  - Image from raw bytes")
    print("  - Image decoded from a base64 string")


if __name__ == "__main__":
    main()
