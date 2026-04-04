"""Demonstrate and document slide background techniques.

This example demonstrates:
- Solid color backgrounds by embedding a full-slide colored shape
- Gradient background concepts (Go API reference)
- Picture (image) backgrounds via add_image() covering the full slide
- Theme-default background (no explicit background shape)
- Background type constants via bullet slides
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ShapeType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

# Minimal 1x1 blue PNG for picture-background demo.
_BLUE_PNG = bytes([
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
    0x60,
    0x98,
    0xFF,
    0xFF,
    0x00,
    0x00,
    0x02,
    0x00,
    0x01,
    0xE5,
    0x27,
    0xDE,
    0xFC,
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

# Slide canvas dimensions (16:9) in Inches.
_SLIDE_W = Inches(10)
_SLIDE_H = Inches(7.5)


def _add_solid_bg_slide(
    prs: Presentation, color: str, title: str, bullets: list[str]
) -> None:
    """Simulate a solid background by placing a full-slide colored rectangle."""
    idx = prs.add_slide(title, layout=SlideLayoutType.TITLE_ONLY).index
    # Full-slide background rectangle
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(0), Inches(0), _SLIDE_W, _SLIDE_H),
        properties={"fill_color": color},
    )
    # Content textbox on top
    for i, text in enumerate(bullets):
        prs.add_textbox(
            idx,
            Inches(0.5),
            Inches(1.5 + i * 0.6),
            Inches(9),
            Inches(0.55),
            text=text,
        )


def _add_gradient_reference_slide(prs: Presentation) -> None:
    """Gradient background concept (Go API reference)."""
    prs.add_bullet_slide(
        "Gradient Backgrounds (Go API)",
        [
            "In Go: NewShapeGradientFill('linear', stops).WithLinearAngle(135)",
            "Then:  NewGradientBackground(gradFill) and slide.Background = &bg",
            "Stops: NewShapeGradientStop(0, '4472C4'), NewShapeGradientStop(100, 'FFFFFF')",
            "Angles: 90° = top-to-bottom, 0° = left-to-right, 135° = diagonal",
        ],
    )


def _add_picture_bg_slide(prs: Presentation) -> None:
    """Picture background using a full-slide image."""
    idx = prs.add_slide("Picture Background", layout=SlideLayoutType.TITLE_ONLY).index
    # Full-slide image covers the background
    prs.add_image(
        idx,
        None,
        (Inches(0), Inches(0), _SLIDE_W, _SLIDE_H),
        data=_BLUE_PNG,
        image_format="png",
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(1.2),
        Inches(9),
        Inches(0.6),
        text="Full-slide image acts as a picture background.",
    )


def _add_default_bg_slide(prs: Presentation) -> None:
    """Default theme background (no explicit background)."""
    prs.add_bullet_slide(
        "Default Theme Background",
        [
            "No explicit background shape - uses the presentation theme default.",
            "slide.Background = nil in Go (or simply not set).",
        ],
    )


def _add_constants_slide(prs: Presentation) -> None:
    """Enumerate Go API background type constants."""
    prs.add_bullet_slide(
        "Background Type Constants (Go API)",
        [
            'SlideBackgroundSolid    = "solid"',
            'SlideBackgroundGradient = "gradient"',
            'SlideBackgroundPicture  = "picture"',
        ],
    )


def main() -> None:
    """Create presentation demonstrating slide background options."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "77_background_api.pptx"

    with Presentation.new("Background API Demo") as prs:
        _add_solid_bg_slide(
            prs,
            "1565C0",
            "Solid Color Background - Deep Blue",
            [
                "Background: deep blue (1565C0)",
                "Simulated via full-slide rectangle shape.",
            ],
        )
        _add_solid_bg_slide(
            prs,
            "F0F4FF",
            "Solid Color Background - Light",
            [
                "Background: light blue-grey (F0F4FF)",
                "Light backgrounds work well with dark content.",
            ],
        )
        _add_gradient_reference_slide(prs)
        _add_picture_bg_slide(prs)
        _add_default_bg_slide(prs)
        _add_constants_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: solid-color background shape, picture background image,")
    print("  gradient concept reference, default theme background, constant listing")


if __name__ == "__main__":
    main()
