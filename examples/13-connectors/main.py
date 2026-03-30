"""Demonstrate the three connector types: straight, elbow, and curved.

This example demonstrates:
- Straight connector between two points
- Elbow (right-angle) connector
- Curved (bezier) connector
- Shapes placed alongside connectors to show connection points
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_overview_slide(prs: Presentation) -> None:
    """Add an overview slide listing the connector types."""
    prs.add_bullet_slide(
        "Connectors Demo",
        [
            "gopptx supports three connector types:",
            "  Straight — direct line between two points",
            "  Elbow   — right-angle (bent) connector",
            "  Curved  — smooth bezier-curve connector",
        ],
    )


def _add_straight_connector_slide(prs: Presentation) -> None:
    """Add a slide demonstrating a straight connector."""
    prs.add_slide("Straight Connector", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(0.25), Inches(1.75), Inches(0.75), Inches(0.5)),
        text="A",
        properties={"fill": {"solid": "4472C4"}},
    )
    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(5), Inches(1.75), Inches(0.75), Inches(0.5)),
        text="B",
        properties={"fill": {"solid": "4472C4"}},
    )
    prs.add_connector(idx, "STRAIGHT", Inches(1), Inches(2), Inches(5), Inches(2))
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3),
        Inches(6),
        Inches(0.5),
        text="add_connector(idx, 'STRAIGHT', x1, y1, x2, y2)",
    )


def _add_elbow_connector_slide(prs: Presentation) -> None:
    """Add a slide demonstrating an elbow connector."""
    prs.add_slide("Elbow Connector", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(0.25), Inches(1.75), Inches(0.75), Inches(0.5)),
        text="C",
        properties={"fill": {"solid": "ED7D31"}},
    )
    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(5), Inches(3.75), Inches(0.75), Inches(0.5)),
        text="D",
        properties={"fill": {"solid": "ED7D31"}},
    )
    prs.add_connector(idx, "ELBOW", Inches(1), Inches(2), Inches(5), Inches(4))
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(5),
        Inches(6),
        Inches(0.5),
        text="add_connector(idx, 'ELBOW', x1, y1, x2, y2)",
    )


def _add_curved_connector_slide(prs: Presentation) -> None:
    """Add a slide demonstrating a curved connector."""
    prs.add_slide("Curved Connector", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "ELLIPSE",
        bounds=(Inches(0.25), Inches(2.75), Inches(0.75), Inches(0.5)),
        text="E",
        properties={"fill": {"solid": "A9D18E"}},
    )
    prs.add_shape(
        idx,
        "ELLIPSE",
        bounds=(Inches(5), Inches(4.75), Inches(0.75), Inches(0.5)),
        text="F",
        properties={"fill": {"solid": "A9D18E"}},
    )
    prs.add_connector(idx, "CURVED", Inches(1), Inches(3), Inches(5), Inches(5))
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(6),
        Inches(6),
        Inches(0.5),
        text="add_connector(idx, 'CURVED', x1, y1, x2, y2)",
    )


def main() -> None:
    """Create a presentation demonstrating all three connector types."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Connectors Demo") as prs:
        _add_overview_slide(prs)
        _add_straight_connector_slide(prs)
        _add_elbow_connector_slide(prs)
        _add_curved_connector_slide(prs)

        output_path = output_dir / "13-connectors.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 4 slides demonstrating connector types:")
    print("  - Overview")
    print("  - Straight connector")
    print("  - Elbow (right-angle) connector")
    print("  - Curved connector")


if __name__ == "__main__":
    main()
