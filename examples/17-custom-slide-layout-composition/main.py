"""Example showing custom slide layout composition with reusable helper functions.

This example demonstrates building structured slide layouts using:
- Reusable layout functions
- Consistent positioning and styling
- Multiple shape types (textbox, rounded rectangles, images)
"""

from gopptx import Presentation, ShapeType
from gopptx.schemas import Inches


def build_structured_slide(slide, title, tag, points, image_path):
    """Build a structured slide with reusable layout components.

    Layout structure:
    - Title bar at top (full width)
    - Left panel: Tag + bullet points (blue background)
    - Right panel: Image in white container
    - Bottom: Action buttons (Plan, Build, Review)

    Args:
        slide: Slide object to populate
        title: Slide title text
        tag: Category/section label
        points: List of bullet point strings
        image_path: Path to image file
    """
    body = "\n".join(points)

    # ===== Title Bar (full width) =====
    slide.add_textbox(
        Inches(0.8),
        Inches(0.35),
        Inches(6.4),
        Inches(0.45),
        text=title,
    )

    # ===== Left Panel: Summary Box (blue background) =====
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(0.75), Inches(1.15), Inches(3.45), Inches(4.85)),
        text=f"{tag}\n\n{body}",
        properties={
            "fill": {"solid": "EEF4FB"},  # Light blue background
            "line": {"color": "B7CBE3", "width_emu": 12700},  # Blue border
        },
    )

    # ===== Right Panel: Image Container (white background) =====
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(4.45), Inches(1.15), Inches(4.0), Inches(3.15)),
        properties={
            "fill": {"solid": "FFFFFF"},  # White background
            "line": {"color": "C8D0DA", "width_emu": 12700},  # Gray border
        },
    )

    # ===== Image (inside container) =====
    slide.add_image(
        image_path,
        (Inches(4.63), Inches(1.32), Inches(3.84), Inches(2.78)),
    )

    # ===== Footer: Action Buttons =====
    # Plan button
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(4.45), Inches(4.45), Inches(1.2), Inches(0.48)),
        text="Plan",
        properties={
            "fill": {"solid": "DCE6F2"},  # Light blue
            "line": {"color": "DCE6F2", "width_emu": 12700},
        },
    )

    # Build button
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(5.78), Inches(4.45), Inches(1.2), Inches(0.48)),
        text="Build",
        properties={
            "fill": {"solid": "E2F0D9"},  # Light green
            "line": {"color": "E2F0D9", "width_emu": 12700},
        },
    )

    # Review button
    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(7.11), Inches(4.45), Inches(1.2), Inches(0.48)),
        text="Review",
        properties={
            "fill": {"solid": "FCE4D6"},  # Light orange
            "line": {"color": "FCE4D6", "width_emu": 12700},
        },
    )

    return slide


def main():
    """Create presentation with custom slide layouts."""
    with Presentation.new("I06 Custom Slide Layout Composition") as prs:
        # ===== Slide 1: Overview =====
        prs.update_slide(0, layout="blank")
        build_structured_slide(
            prs.slides[0],
            "Custom Slide Layout Composition",
            "Overview",
            [
                "Reusable title and summary region",
                "Dedicated image card for visuals",
                "Consistent footer chips for structure",
            ],
            "examples/assets/test_image.png",
        )

        # ===== Slide 2: Execution =====
        slide2 = prs.add_slide("Custom Slide Layout Composition", layout="blank")
        build_structured_slide(
            slide2,
            "Custom Slide Layout Composition",
            "Execution",
            [
                "Same helper builds another slide",
                "Only the content tags change",
                "Layout stays stable across pages",
            ],
            "examples/assets/test_image.png",
        )

        # Save presentation
        prs.save("examples/output/17-custom-slide-layout-composition.pptx")
        print("Presentation created: examples/output/17-custom-slide-layout-composition.pptx")

        # Print summary
        print("\n" + "=" * 70)
        print("CUSTOM SLIDE LAYOUT COMPOSITION")
        print("=" * 70)
        print("\nPATTERN: Reusable Slide Layout Functions")
        print("\nBenefits:")
        print("  [+] Consistent design across slides")
        print("  [+] Easy to maintain and update layout")
        print("  [+] Reusable for different content")
        print("  [+] Clear separation of concerns")
        print("\nLayout Structure:")
        print("  - Title Bar (full width)")
        print("  - Left Panel: Summary with bullet points")
        print("  - Right Panel: Image in white container")
        print("  - Footer: Action buttons (Plan, Build, Review)")
        print("\nAPI Usage:")
        print("  slide.add_textbox() - Title")
        print("  slide.add_shape()   - Panels, buttons")
        print("  slide.add_image()   - Image content")
        print("\nTemplate Pattern:")
        print("  1. Create helper function with layout logic")
        print("  2. Accept reusable parameters (title, tag, points, image)")
        print("  3. Call multiple times with different content")
        print("  4. Consistent design maintained automatically")


if __name__ == "__main__":
    main()
