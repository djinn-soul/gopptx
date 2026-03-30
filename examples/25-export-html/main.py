"""Demonstrate creating content for HTML export.

This example demonstrates:
- Building slides with shapes, images, and code-like text content
- Documenting the HTML export workflow
- Creating content that maps well to HTML/CSS rendering
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

# Minimal 1x1 red pixel PNG
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
    0x06,
    0x00,
    0x00,
    0x00,
    0x1F,
    0x15,
    0xC4,
    0x89,
    0x00,
    0x00,
    0x00,
    0x0A,
    0x49,
    0x44,
    0x41,
    0x54,
    0x78,
    0x9C,
    0x63,
    0xF8,
    0xCF,
    0x00,
    0x00,
    0x02,
    0x03,
    0x01,
    0x01,
    0x24,
    0x95,
    0x8C,
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


def _add_html_export_slide(prs: Presentation) -> None:
    """Add an intro slide about HTML export."""
    prs.add_bullet_slide(
        "HTML Export Demo",
        [
            "This presentation was exported to HTML via gopptx.",
            "It mimics the ppt-rs export style.",
        ],
    )


def _add_visual_content_slide(prs: Presentation, image_path: str) -> None:
    """Add a slide with a shape and image."""
    prs.add_slide("Visual Content", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(1), Inches(1), Inches(3), Inches(1)),
        text="This is a shape text.",
        properties={"fill": {"solid": "0078D4"}},
    )
    prs.add_image(idx, image_path, (Inches(5), Inches(1), Inches(1), Inches(1)))


def _add_code_simulation_slide(prs: Presentation) -> None:
    """Add a slide simulating a code block."""
    prs.add_slide("Code Simulation", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(1), Inches(2), Inches(3), Inches(2)),
        text='func main() {\n    fmt.Println("Hello World")\n}',
        properties={"fill": {"solid": "F0F0F0"}},
    )


def main() -> None:
    """Create a presentation demonstrating HTML export content."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
        tmp.write(_RED_PNG)
        image_path = tmp.name

    try:
        with Presentation.new("Export Demo") as prs:
            _add_html_export_slide(prs)
            _add_visual_content_slide(prs, image_path)
            _add_code_simulation_slide(prs)

            output_path = output_dir / "25-export-html.pptx"
            prs.save(str(output_path))
            print(f"Saved: {output_path}")
    finally:
        Path(image_path).unlink(missing_ok=True)

    print("\n=== SUMMARY ===")
    print("Created a 3-slide presentation for HTML export:")
    print("  - Intro slide with export description")
    print("  - Visual content with shape and image")
    print("  - Code simulation slide")


if __name__ == "__main__":
    main()
