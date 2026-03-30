"""Handout master example in Python.

Configures handout layouts and writes multiple output variants.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def build_base_slides(prs: Presentation) -> None:
    prs.update_slide(
        0, title="Handout Master Demo", bullets=["Demonstrate handout modes"]
    )
    prs.add_slide("Slide 2", bullets=["Additional content"])
    prs.add_slide("Slide 3", bullets=["Third slide content"])


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)

    default_path = output_dir / "29_handout_master_default.pptx"
    six_up_path = output_dir / "29_handout_master_6up.pptx"
    outline_path = output_dir / "29_handout_master_outline.pptx"

    with Presentation.new("Handout Master Default (Python)") as prs:
        build_base_slides(prs)
        prs.save(str(default_path))

    with Presentation.new("Handout Master 6-Up (Python)") as prs:
        build_base_slides(prs)
        prs.update_handout_master(
            orientation="landscape",
            slides_per_page=6,
        )
        prs.save(str(six_up_path))
        print(f"6-up handout metadata: {prs.get_handout_master()}")

    with Presentation.new("Handout Master Outline (Python)") as prs:
        prs.update_slide(
            0, title="Outline Mode", bullets=["Single-slide outline output"]
        )
        prs.update_handout_master(
            orientation="portrait",
            slides_per_page=1,
        )
        prs.save(str(outline_path))

    print(f"Created: {default_path}")
    print(f"Created: {six_up_path}")
    print(f"Created: {outline_path}")


if __name__ == "__main__":
    main()
