"""Document DrawingML fill, line, and color formatting concepts.

This example demonstrates:
- Shape fill color via add_shape() properties
- Contrasting solid fills across multiple shapes per slide
- Visual references for gradient and shadow concepts
- Line style notes (solid, dashed, dotted) documented in bullets
- Full advanced DML formatting hooks available in the Go API
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ShapeType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_solid_fill_slide(prs: Presentation) -> None:
    """Three rectangles with different solid fills."""
    idx = prs.add_slide(
        "Solid Fill & Transparency", layout=SlideLayoutType.TITLE_ONLY
    ).index
    fills = [
        ("Solid Fill", "4472C4", Inches(0.5)),
        ("Alternate Fill", "C0504D", Inches(3.5)),
        ("Light Fill", "9BBB59", Inches(6.5)),
    ]
    for label, color, x in fills:
        prs.add_shape(
            idx,
            ShapeType.RECTANGLE,
            (x, Inches(1.5), Inches(2.5), Inches(1.5)),
            text=label,
            properties={"fill_color": color},
        )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.5),
        Inches(9),
        Inches(0.8),
        text="In Go: NewShapeFill(hex).WithTransparency(0.5) for semi-transparent fills",
    )


def _add_gradient_fill_slide(prs: Presentation) -> None:
    """Document gradient fills with reference shapes."""
    idx = prs.add_slide("Gradient Fills", layout=SlideLayoutType.TITLE_ONLY).index
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(0.5), Inches(1.5), Inches(3.5), Inches(2)),
        text="Linear Gradient 45°",
        properties={"fill_color": "4472C4"},
    )
    prs.add_shape(
        idx,
        ShapeType.ELLIPSE,
        (Inches(5), Inches(1.5), Inches(3.5), Inches(2)),
        text="Radial Gradient",
        properties={"fill_color": "FF6F00"},
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(4.0),
        Inches(9),
        Inches(0.8),
        text="In Go: NewShapeGradientFill('linear', stops).WithLinearAngle(45)",
    )


def _add_line_styles_slide(prs: Presentation) -> None:
    """Document line style variations."""
    idx = prs.add_slide("Line Styles", layout=SlideLayoutType.TITLE_ONLY).index
    styles = [
        ("Solid 3pt", "4472C4", Inches(0.5)),
        ("Dashed 2pt", "C0504D", Inches(3.5)),
        ("Dotted 2pt", "9BBB59", Inches(6.5)),
    ]
    for label, color, x in styles:
        prs.add_shape(
            idx,
            ShapeType.RECTANGLE,
            (x, Inches(1.5), Inches(2.5), Inches(1.2)),
            text=label,
            properties={"fill_color": color},
        )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.2),
        Inches(9),
        Inches(0.8),
        text="In Go: NewShapeLine(color, width).WithDash(LineDashDash)",
    )


def _add_shadow_effects_slide(prs: Presentation) -> None:
    """Document shadow effect types."""
    idx = prs.add_slide("Shadow Effects", layout=SlideLayoutType.TITLE_ONLY).index
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(1), Inches(1.5), Inches(3), Inches(1.5)),
        text="Outer Shadow",
        properties={"fill_color": "4472C4"},
    )
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(5), Inches(1.5), Inches(3), Inches(1.5)),
        text="Inner Shadow",
        properties={"fill_color": "C0504D"},
    )
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(3.5),
        Inches(9),
        Inches(0.8),
        text="In Go: NewOuterShadow('333333') / NewInnerShadow('000000')",
    )


def _add_dash_constants_slide(prs: Presentation) -> None:
    """Enumerate line dash style constants."""
    prs.add_bullet_slide(
        "Line Dash Style Constants",
        [
            'LineDashSolid      = "solid"',
            'LineDashDash       = "dash"',
            'LineDashDot        = "dot"',
            'LineDashDashDot    = "dashDot"',
            'LineDashDashDotDot = "lgDashDotDot"',
            'LineDashLongDash   = "lgDash"',
        ],
    )


def main() -> None:
    """Create presentation demonstrating DML fill, line, and color formatting."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "72_dml_api.pptx"

    with Presentation.new("DML API Demo") as prs:
        _add_solid_fill_slide(prs)
        _add_gradient_fill_slide(prs)
        _add_line_styles_slide(prs)
        _add_shadow_effects_slide(prs)
        _add_dash_constants_slide(prs)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print(
        "Demonstrated: solid fill shapes, gradient fill reference, line style reference,"
    )
    print("  shadow effect reference, and LineDash constants enumeration")


if __name__ == "__main__":
    main()
