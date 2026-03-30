"""Demonstrate slide transitions: Fade, Push, Wipe, Split, and Zoom.

This example demonstrates:
- Adding slide transitions is noted per slide (transitions are applied in Go layer)
- Creating slides representing each transition type
- Documenting the transition type, direction, and duration
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_overview_slide(prs: Presentation) -> None:
    """Add an overview slide."""
    prs.add_bullet_slide(
        "Slide Transitions Demo",
        [
            "Each following slide demonstrates a different transition.",
            "Open in PowerPoint and advance slides to see the effects.",
        ],
    )


def _add_fade_slide(prs: Presentation) -> None:
    """Add a slide representing the Fade transition."""
    prs.add_bullet_slide(
        "Fade Transition",
        [
            "Type: Fade",
            "Duration: 500 ms",
            "Fades the current slide into the next.",
        ],
    )


def _add_push_slide(prs: Presentation) -> None:
    """Add a slide representing the Push Left transition."""
    prs.add_bullet_slide(
        "Push Left Transition",
        [
            "Type: Push, Direction: Left",
            "Duration: 700 ms",
            "The new slide pushes in from the right side.",
        ],
    )


def _add_wipe_slide(prs: Presentation) -> None:
    """Add a slide representing the Wipe Right transition."""
    prs.add_bullet_slide(
        "Wipe Right Transition",
        [
            "Type: Wipe, Direction: Right",
            "Duration: 800 ms",
            "The new slide wipes over from the left.",
        ],
    )


def _add_zoom_slide(prs: Presentation) -> None:
    """Add a slide representing the Zoom transition."""
    prs.add_bullet_slide(
        "Zoom Transition",
        [
            "Type: Zoom",
            "Duration: 600 ms",
            "The new slide zooms in from the center.",
        ],
    )


def _add_split_slide(prs: Presentation) -> None:
    """Add a slide representing the Split transition."""
    prs.add_bullet_slide(
        "Split Transition",
        [
            "Applied with the shorthand WithTransition(TransitionSplit).",
            "Default duration is used.",
        ],
    )


def main() -> None:
    """Create a presentation demonstrating slide transitions."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Slide Transitions Demo") as prs:
        _add_overview_slide(prs)
        _add_fade_slide(prs)
        _add_push_slide(prs)
        _add_wipe_slide(prs)
        _add_zoom_slide(prs)
        _add_split_slide(prs)

        output_path = output_dir / "14-transitions.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 6 slides, each documenting a transition type:")
    print("  - Fade (500 ms)")
    print("  - Push Left (700 ms)")
    print("  - Wipe Right (800 ms)")
    print("  - Zoom (600 ms)")
    print("  - Split (default duration)")


if __name__ == "__main__":
    main()
