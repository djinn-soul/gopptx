"""Hyperlinks example in Python.

Adds shape-level click actions and run-level hyperlink text.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "31_hyperlinks.pptx"

    with Presentation.new("Hyperlinks (Python)") as prs:
        slide = prs.slides[0]
        slide.title = "Hyperlinks"

        slide.add_shape(
            "roundRect",
            (Inches(0.8), Inches(1.4), Inches(3.8), Inches(1.1)),
            text="Open README.md (relative path)",
            click_action={"address": "README.md", "tooltip": "Open project README"},
            properties={"fill": {"solid": "DBEAFE"}, "line": {"color": "1D4ED8"}},
        )

        run_shape = slide.add_shape(
            "rect",
            (Inches(0.8), Inches(3.0), Inches(6.2), Inches(1.0)),
            properties={"fill": {"solid": "FFFFFF"}, "line": {"color": "CBD5E1"}},
        )
        slide.set_shape_runs(
            run_shape,
            [
                {"text": "Visit the project site: "},
                {
                    "text": "OpenAI",
                    "color": "2563EB",
                    "underline": "single",
                    "hyperlink": {"address": "https://openai.com"},
                },
            ],
        )

        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
