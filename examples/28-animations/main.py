"""Animations and transitions example in Python.

Applies entrance/emphasis animations and slide transitions.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation, ShapeType
from gopptx.animations import (
    ANIMATION_AFTER_PREVIOUS,
    ANIMATION_EMPHASIS_PULSE,
    ANIMATION_ENTRANCE_FADE,
    ANIMATION_ENTRANCE_FLY_IN,
    ANIMATION_ON_CLICK,
)
from gopptx.schemas import Inches
from gopptx.transitions import TRANSITION_FADE, TRANSITION_PUSH


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "28_animations.pptx"

    with Presentation.new("Animations (Python)") as prs:
        slide1 = prs.slides[0]
        slide1.title = "Animations & Transitions"
        slide1.set_transition(TRANSITION_FADE, duration_ms=800)
        slide1.add_shape(
            ShapeType.ROUNDED_RECTANGLE,
            (Inches(0.9), Inches(1.5), Inches(2.2), Inches(1.2)),
            text="Fade In",
            properties={"fill": {"solid": "FEE2E2"}, "line": {"color": "DC2626"}},
        )
        slide1.add_shape(
            ShapeType.ELLIPSE,
            (Inches(3.5), Inches(1.5), Inches(2.2), Inches(1.2)),
            text="Fly In",
            properties={"fill": {"solid": "DCFCE7"}, "line": {"color": "16A34A"}},
        )
        slide1.add_shape(
            ShapeType.TRIANGLE,
            (Inches(6.1), Inches(1.5), Inches(2.2), Inches(1.2)),
            text="Pulse",
            properties={"fill": {"solid": "DBEAFE"}, "line": {"color": "1D4ED8"}},
        )
        slide1.add_animation(
            1,
            ANIMATION_ENTRANCE_FADE,
            trigger=ANIMATION_ON_CLICK,
            duration_ms=800,
        )
        slide1.add_animation(
            2,
            ANIMATION_ENTRANCE_FLY_IN,
            trigger=ANIMATION_AFTER_PREVIOUS,
            duration_ms=900,
        )
        slide1.add_animation(
            3,
            ANIMATION_EMPHASIS_PULSE,
            trigger=ANIMATION_AFTER_PREVIOUS,
            duration_ms=700,
        )

        slide2 = prs.add_slide("Transition Push")
        slide2.set_transition(TRANSITION_PUSH, duration_ms=1000)
        slide2.add_textbox(
            Inches(1.0),
            Inches(2.0),
            Inches(7.5),
            Inches(1.0),
            text="This slide uses a push transition.",
        )

        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
