"""Demonstrate layout helper utilities: stack, distribute, grid, and center.

This example demonstrates:
- Stack layout: vertically stacked shapes with uniform spacing
- Distribute uniform: evenly spaced shapes across a bounding box
- Grid layout: shapes arranged in a 2x3 grid
- Center helper: centering a shape on the slide
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_stack_slide(prs: Presentation) -> None:
    """Add a slide showing vertically stacked shapes."""
    prs.add_slide("Stack Layout", layout=SlideLayoutType.TITLE_ONLY)
    idx = prs.slide_count - 1
    labels = ["Init", "Build", "Validate", "Export"]
    colors = ["4472C4", "ED7D31", "A9D18E", "FF0000"]
    y_start = Inches(1.2)
    height = Inches(0.6)
    gap = Inches(0.12)
    for i, label in enumerate(labels):
        y = y_start + i * (height + gap)
        prs.add_shape(
            idx,
            "RECTANGLE",
            (Inches(1), y, Inches(8), height),
            text=label,
            properties={"fill_color": colors[i]},
        )


def _add_distribute_slide(prs: Presentation) -> None:
    """Add a slide showing horizontally distributed shapes."""
    prs.add_slide("Distribute Uniform", layout=SlideLayoutType.TITLE_ONLY)
    idx = prs.slide_count - 1
    labels = ["Alpha", "Beta", "Gamma", "Delta", "Epsilon"]
    colors = ["4472C4", "ED7D31", "A9D18E", "FF0000", "FFC000"]
    total_width = Inches(9)
    elem_w = Inches(1.5)
    n = len(labels)
    spacing = (total_width - n * elem_w) / (n - 1)
    x_start = Inches(0.5)
    y = Inches(2)
    h = Inches(1.5)
    for i, label in enumerate(labels):
        x = x_start + i * (elem_w + spacing)
        prs.add_shape(
            idx,
            "ROUNDED_RECTANGLE",
            (x, y, elem_w, h),
            text=label,
            properties={"fill_color": colors[i]},
        )


def _add_grid_slide(prs: Presentation) -> None:
    """Add a slide showing shapes in a 2x3 grid."""
    prs.add_slide("Grid Layout (2x3)", layout=SlideLayoutType.TITLE_ONLY)
    idx = prs.slide_count - 1
    colors = ["5B9BD5", "ED7D31", "A9D18E", "FF0000", "FFC000", "7030A0"]
    cols, rows = 2, 3
    gap = Inches(0.2)
    box_x, box_y = Inches(0.5), Inches(1.9)
    box_w, box_h = Inches(9), Inches(4.9)
    cell_w = (box_w - (cols - 1) * gap) / cols
    cell_h = (box_h - (rows - 1) * gap) / rows
    for r in range(rows):
        for c in range(cols):
            i = r * cols + c
            x = box_x + c * (cell_w + gap)
            y = box_y + r * (cell_h + gap)
            prs.add_shape(
                idx,
                "RECTANGLE",
                (x, y, cell_w, cell_h),
                text=f"Cell {i + 1}",
                properties={"fill_color": colors[i]},
            )


def _add_center_slide(prs: Presentation) -> None:
    """Add a slide with a centered ellipse."""
    prs.add_slide("Center Helper", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1
    cx, cy = Inches(4), Inches(2)
    slide_w, slide_h = Inches(10), Inches(7.5)
    x = (slide_w - cx) / 2
    y = (slide_h - cy) / 2
    prs.add_shape(
        idx,
        "ELLIPSE",
        (x, y, cx, cy),
        text="Centered Ellipse",
        properties={"fill_color": "1B6CA8"},
    )


def main() -> None:
    """Create presentation demonstrating layout helper utilities."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Task 35: Layout Helpers") as prs:
        _add_stack_slide(prs)
        _add_distribute_slide(prs)
        _add_grid_slide(prs)
        _add_center_slide(prs)
        output_path = output_dir / "35_layout_helpers.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("4 slides: Stack, Distribute, Grid, Center")


if __name__ == "__main__":
    main()
