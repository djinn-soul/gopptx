"""Document hyperlinks and click-action patterns on shapes.

This example demonstrates:
- Adding shapes with add_shape() and labels explaining action concepts
- Documenting URL hyperlinks, slide navigation, email links, and hover actions
- Using add_bullet_slide() and add_textbox() to explain each action type
- Full click/hover action wiring is available in the Go API
- Python slides document concepts and show shape placement for those actions
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ShapeType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_url_hyperlink_slide(prs: Presentation) -> None:
    """Slide describing URL hyperlink actions."""
    prs.add_bullet_slide(
        "URL Hyperlink on Shape",
        [
            "Shapes can carry click actions (click-to-URL).",
            "URL: https://github.com/djinn-soul/gopptx",
            "In Go: shape.WithClickAction(action.NewHyperlink(...))",
            "Tooltip and highlight-on-click are configurable.",
        ],
    )


def _add_navigation_slide(prs: Presentation) -> None:
    """Slide describing slide navigation actions."""
    idx = prs.add_slide(
        "Slide Navigation Actions", layout=SlideLayoutType.TITLE_ONLY
    ).index
    labels = [
        ("Next Slide", "4CAF50", Inches(0.5), Inches(1.5)),
        ("Prev Slide", "F44336", Inches(3.0), Inches(1.5)),
        ("First Slide", "FF9800", Inches(5.5), Inches(1.5)),
        ("Last Slide", "9C27B0", Inches(0.5), Inches(3.0)),
        ("End Show", "455A64", Inches(3.0), Inches(3.0)),
    ]
    for label, color, x, y in labels:
        prs.add_shape(
            idx,
            ShapeType.RECTANGLE,
            (x, y, Inches(2.2), Inches(0.8)),
            text=label,
            properties={"fill_color": color},
        )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(4.2),
        Inches(9),
        Inches(0.6),
        text="In Go: WithClickAction(action.NewHyperlink(HyperlinkNextSlide()))",
    )


def _add_jump_to_slide_slide(prs: Presentation) -> None:
    """Slide describing jump-to-specific-slide action."""
    idx = prs.add_slide(
        "Jump to Specific Slide", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(1), Inches(2), Inches(4), Inches(1)),
        text="Jump to Slide 5",
        properties={"fill_color": "4BACC6"},
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.5),
        Inches(9),
        Inches(0.6),
        text="In Go: WithClickAction(action.NewHyperlink(action.HyperlinkSlide(5)))",
    )


def _add_email_hyperlink_slide(prs: Presentation) -> None:
    """Slide describing email hyperlink action."""
    idx = prs.add_slide("Email Hyperlink", layout=SlideLayoutType.TITLE_ONLY).index
    prs.add_shape(
        idx,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(1), Inches(2), Inches(6), Inches(1)),
        text="Click to email hello@example.com",
        properties={"fill_color": "EBF1DE"},
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.5),
        Inches(9),
        Inches(0.8),
        text="In Go: HyperlinkEmailWithSubject('hello@example.com', 'Subject')",
    )


def _add_hover_action_slide(prs: Presentation) -> None:
    """Slide describing hover action."""
    idx = prs.add_slide("Hover Action", layout=SlideLayoutType.TITLE_ONLY).index
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(1), Inches(2), Inches(6), Inches(1)),
        text="Hover over me!",
        properties={"fill_color": "FDE9D9"},
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.5),
        Inches(9),
        Inches(0.6),
        text="In Go: shape.WithHoverAction(action.NewHyperlink(...))",
    )


def main() -> None:
    """Create presentation documenting hyperlink and action concepts."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "71_action_api.pptx"

    with Presentation.new("Action API Demo") as prs:
        _add_url_hyperlink_slide(prs)
        _add_navigation_slide(prs)
        _add_jump_to_slide_slide(prs)
        _add_email_hyperlink_slide(prs)
        _add_hover_action_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: URL hyperlinks, slide navigation buttons, jump-to-slide,")
    print("  email hyperlinks, and hover actions via shape placement + labels")


if __name__ == "__main__":
    main()
