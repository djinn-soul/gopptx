"""Demonstrate basic image embedding and document advanced image options.

This example demonstrates:
- Embedding images from raw bytes with add_image_from_bytes()
- Embedding images from a file path with add_image()
- Image placement using Inches() for position and size
- Reference slides for rotation, flip, crop, shadow, reflection, and
  decorative-image options available in the Go API
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

# Minimal 1x1 white PNG for embedding demos (no external file needed).
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
    """Embed the minimal PNG from raw bytes."""
    idx = prs.add_slide("Image from Bytes", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(1.2),
        Inches(9),
        Inches(0.5),
        text="add_image_from_bytes(slide_idx, data, x, y, w, h) - no file on disk required.",
    )
    prs.add_image_from_bytes(
        idx,
        _WHITE_PNG,
        Inches(1),
        Inches(1.8),
        Inches(4),
        Inches(3),
    )


def _add_image_from_file_slide(prs: Presentation, tmp_png: str) -> None:
    """Embed image loaded from a file path."""
    idx = prs.add_slide("Image from File Path", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(1.2),
        Inches(9),
        Inches(0.5),
        text=f"add_image(slide_idx, path, x, y, w, h) - reads from: {tmp_png}",
    )
    prs.add_image(
        idx,
        tmp_png,
        Inches(1),
        Inches(1.8),
        Inches(4),
        Inches(3),
    )


def _add_rotation_flip_slide(prs: Presentation) -> None:
    """Document rotation and flip image options."""
    prs.add_bullet_slide(
        "Image Rotation & Flip (Go API)",
        [
            "image.WithRotation(degrees)  - rotate the image by N degrees",
            "image.WithFlip(horiz, vert)  - mirror horizontally or vertically",
            "Both options are set on the Image object before embedding",
            "Available via Go pptx.NewImageFromBytes(...).WithRotation(30)",
        ],
    )


def _add_crop_slide(prs: Presentation) -> None:
    """Document image crop options."""
    prs.add_bullet_slide(
        "Image Crop (Go API)",
        [
            "image.WithCrop(left, right, top, bottom)",
            "Values are 0.0-1.0 fractions of the image dimension",
            "WithCrop(0.1, 0.1, 0.1, 0.1) crops 10% from each edge",
        ],
    )


def _add_shadow_reflection_slide(prs: Presentation) -> None:
    """Document shadow and reflection options."""
    prs.add_bullet_slide(
        "Shadow & Reflection (Go API)",
        [
            "image.WithShadow(true)      - adds an outer shadow effect",
            "image.WithReflection(true)  - adds a reflection below the image",
            "image.WithDecorative(true)  - marks image as decorative (accessibility)",
            "Decorative images are skipped by screen readers",
        ],
    )


def main() -> None:
    """Create presentation demonstrating image embedding."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "73_image_api.pptx"

    with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
        tmp.write(_WHITE_PNG)
        tmp_png = tmp.name

    with Presentation.new("Image API Demo") as prs:
        _add_image_from_bytes_slide(prs)
        _add_image_from_file_slide(prs, tmp_png)
        _add_rotation_flip_slide(prs)
        _add_crop_slide(prs)
        _add_shadow_reflection_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    Path(tmp_png).unlink(missing_ok=True)

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_image_from_bytes, add_image (from file),")
    print("  documented: WithRotation, WithFlip, WithCrop, WithShadow, WithReflection")


if __name__ == "__main__":
    main()
