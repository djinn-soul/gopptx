"""Embedded-fonts-oriented example in Python.

Demonstrates explicit font run settings in slide content.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "30_embedded_fonts.pptx"

    with Presentation.new("Embedded Fonts (Python)") as prs:
        slide = prs.slides[0]
        slide.title = "Font Styling Demo"
        shape_id = slide.add_shape(
            "rect",
            (914400, 1371600, 7315200, 1828800),
            properties={"fill": {"solid": "FFFFFF"}, "line": {"color": "CBD5E1"}},
        )
        slide.set_shape_runs(
            shape_id,
            [
                {
                    "text": "Calibri Light Heading\n",
                    "font": "Calibri Light",
                    "size_pt": 30,
                    "bold": True,
                },
                {"text": "Aptos Body text\n", "font": "Aptos", "size_pt": 20},
                {
                    "text": "Courier New code sample",
                    "font": "Courier New",
                    "size_pt": 18,
                    "code": True,
                },
            ],
        )
        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
