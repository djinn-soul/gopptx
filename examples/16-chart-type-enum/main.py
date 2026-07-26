"""Chart type selection with the ChartType enum.

ChartType is a real enum generated from the Go XLChartType constants, so the two
languages can never drift. Members subclass str, so a member equals its wire
value (ChartType.PIE == "pie"), but passing a bare string is deprecated: the
enum is what a type checker can verify and an IDE can complete.

This example uses the slide object returned by add_slide rather than a slide
index, so a chart cannot land on the wrong slide.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation, Slide
from gopptx.presentation.charts import ChartType
from gopptx.schemas import Inches


def _add_single_chart_slides(prs: Presentation) -> None:
    """One slide per chart type, each built from the ChartType member."""
    samples: list[tuple[str, ChartType, list[str], list[float], str]] = [
        (
            "Column Chart",
            ChartType.COLUMN,
            ["Q1", "Q2", "Q3", "Q4"],
            [12.0, 18.0, 15.0, 24.0],
            "Revenue by Quarter",
        ),
        (
            "Line Chart",
            ChartType.LINE,
            ["Jan", "Feb", "Mar", "Apr", "May"],
            [4.0, 7.0, 6.0, 9.0, 8.0],
            "Monthly Trend",
        ),
        (
            "Pie Chart",
            ChartType.PIE,
            ["Direct", "Search", "Referral", "Social"],
            [45.0, 35.0, 12.0, 8.0],
            "Traffic Sources",
        ),
    ]

    for slide_title, chart_type, categories, values, chart_title in samples:
        slide = prs.add_slide(slide_title, layout="title_only")
        slide.add_chart(
            chart_type,
            categories,
            values,
            title=chart_title,
            bounds=(Inches(0.5), Inches(1.5), Inches(9.0), Inches(5.0)),
        )
        print(f"  {slide_title:<14} -> ChartType.{chart_type.name}")


def _add_same_data_slide(prs: Presentation) -> None:
    """The same data drawn three ways, selected by enum member."""
    categories = ["A", "B", "C", "D"]
    values = [30.0, 40.0, 20.0, 10.0]
    views: list[ChartType] = [ChartType.COLUMN, ChartType.LINE, ChartType.PIE]

    slide: Slide = prs.add_slide("Same Data, Different Views", layout="title_only")
    for position, chart_type in enumerate(views):
        slide.add_chart(
            chart_type,
            categories,
            values,
            title=f"{chart_type.name.title()} View",
            bounds=(
                Inches(0.4 + position * 3.1),
                Inches(1.5),
                Inches(3.0),
                Inches(4.0),
            ),
        )
    print(f"  Same data as {', '.join(t.name for t in views)} on one slide")


def _show_discovery() -> None:
    """ChartType is discoverable without reading the docs."""
    all_types = ChartType.get_all()
    print("\nSupported chart types:")
    for name, value in sorted(all_types.items()):
        print(f"  ChartType.{name:<17} = {value!r}")
    print(f"\n{len(all_types)} names, {len(ChartType.values())} distinct chart types")
    print("(COLUMN and BAR are the same chart; COLUMN is the python-pptx spelling.)")


def _describe_rejection(candidate: str) -> str:
    """Return how ChartType reacts to a candidate value.

    add_chart runs this same check on whatever it is handed, so calling
    ChartType.validate directly shows the rejection without having to pass a
    deliberately wrong argument through the typed add_chart signature.
    """
    try:
        return f"accepted as {ChartType.validate(candidate)!r}"
    except ValueError as exc:
        return f"ValueError: {str(exc)[:70]}..."


def _show_rejected_values() -> None:
    """Invalid chart types are rejected with a message that names the fix."""
    print("\nHow values are checked:")
    for candidate in ("COLUMN", "not_a_chart_type", "pie"):
        print(f"  {candidate!r:<20} -> {_describe_rejection(candidate)}")
    print("  (COLUMN is a member *name*, not a value; ChartType.COLUMN is 'bar'.)")


def main() -> None:
    """Create a presentation demonstrating the ChartType enum."""
    output_dir = Path("examples/output")
    output_dir.mkdir(parents=True, exist_ok=True)

    with Presentation.new("Chart Type Examples") as prs:
        print("Building chart slides:")
        _add_single_chart_slides(prs)
        _add_same_data_slide(prs)
        _show_discovery()
        _show_rejected_values()

        output_path = output_dir / "16-chart-type-enum.pptx"
        prs.save(str(output_path))
        print(f"\nSaved: {output_path}")

    print(
        """
CHART TYPE API

    from gopptx.presentation.charts import ChartType

    slide = prs.add_slide("Sales", layout="title_only")
    slide.add_chart(ChartType.COLUMN, categories, values, title="Sales")

Why the enum:
    - a type checker rejects ChartType typos; a bare string it cannot check
    - IDE autocomplete and inline docs on every member
    - ChartType.get_all() lists everything available
    - generated from the Go XLChartType constants, so Go and Python agree

Bare strings such as "bar" still work but raise a DeprecationWarning and will
stop being accepted in a future release.

Why slide.add_chart over prs.add_chart(slide_index, ...):
    Presentation.new() already creates slide 0, so the first slide you add is
    index 1. Using the Slide object returned by add_slide removes the question.
"""
    )


if __name__ == "__main__":
    main()
