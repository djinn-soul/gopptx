from __future__ import annotations

import argparse
from pathlib import Path

from gopptx import Presentation


def generate_text_frame_demo(out_path: Path) -> None:
    out_path.parent.mkdir(parents=True, exist_ok=True)

    with Presentation.new("S02 Text Frame Demo") as pres:
        pres.add_slide("Text Frame Properties")
        pres.add_shape(0, "rect", (40, 120, 180, 180), text="0.5in margins demo")
        pres.add_shape(0, "rect", (240, 120, 180, 220), text="Top anchor")
        pres.add_shape(0, "rect", (440, 120, 180, 220), text="Bottom anchor")
        pres.add_shape(0, "rect", (40, 360, 180, 60), text="No wrap sample text")
        pres.add_shape(0, "rect", (240, 360, 180, 100), text="Shrink-to-fit sample text")
        pres.save(str(out_path))


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate S02 text-frame usage PPTX from Python docs code."
    )
    parser.add_argument(
        "--out",
        type=Path,
        default=Path("docs/assets/pptx/usage/s02-text-frame-python.pptx"),
        help="Output PPTX path.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    generate_text_frame_demo(args.out)
    print(f"Generated: {args.out}")


if __name__ == "__main__":
    main()

