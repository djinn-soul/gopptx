"""Demonstrate document infrastructure: sections, comments, and metadata.

This example demonstrates:
- add_section() to group slides into named sections
- get_sections() to list sections after creation
- add_comment() / add_author() for slide-level comments
- set_core_properties() / get_core_properties() for document properties
- Slide manipulation: duplicate_slide(), move_slide()
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ShapeType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _build_section_slides(prs: Presentation) -> None:
    """Add slides that will be grouped into sections."""
    prs.add_bullet_slide(
        "Section A - Slide 1",
        ["First slide in Section A", "Introduction content"],
    )
    prs.add_bullet_slide(
        "Section A - Slide 2",
        ["Second slide in Section A", "More introduction content"],
    )
    prs.add_bullet_slide(
        "Section B - Slide 3",
        ["First slide in Section B", "Core content starts here"],
    )
    prs.add_bullet_slide(
        "Section B - Slide 4",
        ["Second slide in Section B", "Supporting details"],
    )


def _build_appendix_slide(prs: Presentation) -> int:
    """Add appendix slide with overlapping shapes for z-order demo."""
    idx = prs.add_slide(
        "Appendix - Shape Demo", layout=SlideLayoutType.TITLE_ONLY
    ).index
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(1), Inches(1.5), Inches(3), Inches(1.5)),
        text="Bottom shape",
        properties={"fill_color": "4472C4"},
    )
    prs.add_shape(
        idx,
        ShapeType.RECTANGLE,
        (Inches(2), Inches(2.0), Inches(3), Inches(1.5)),
        text="Overlapping shape",
        properties={"fill_color": "C0504D"},
    )
    return idx


def _add_comments_slide(prs: Presentation) -> int:
    """Add a slide with author comments."""
    idx = prs.add_bullet_slide(
        "Slide with Comments",
        [
            "This slide has author comments attached.",
            "Comments reference author objects added via add_author().",
        ],
    ).index
    author_id = prs.add_author("Alice", "A")
    prs.add_comment(idx, author_id, "Great point - add more data here.")
    author_id2 = prs.add_author("Bob", "B")
    prs.add_comment(idx, author_id2, "Needs a chart to illustrate this.")
    return idx


def _create_sections(prs: Presentation) -> None:
    """Group slides into named sections."""
    prs.add_section("Introduction", [0, 1])
    prs.add_section("Core Content", [2, 3])
    prs.add_section("Back Matter", [4, 5])


def _log_sections(prs: Presentation) -> None:
    """Print section names and slide counts."""
    sections = prs.get_sections()
    print(f"Sections: {len(sections)}")
    for sec in sections:
        name = sec.get("name", "?")
        slides = sec.get("slide_indices", sec.get("slideIndices", []))
        print(f"  - {name!r}  (slides: {len(slides)})")


def main() -> None:
    """Create presentation demonstrating document infrastructure APIs."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "79_document_infra.pptx"

    with Presentation.new("Document Infrastructure Demo") as prs:
        # Set core metadata
        props = prs.get_core_properties()
        props["title"] = "Document Infrastructure Demo"
        props["creator"] = "gopptx Python API"
        prs.set_core_properties(props)

        # Build content slides (indices 0-3)
        _build_section_slides(prs)

        # Appendix slide (index 4)
        _build_appendix_slide(prs)

        # Comments slide (index 5)
        _add_comments_slide(prs)

        # Create sections grouping the slides above
        _create_sections(prs)
        _log_sections(prs)

        # Demonstrate slide manipulation
        before = prs.slide_count
        prs.duplicate_slide(0)
        print(f"After duplicate_slide(0): {prs.slide_count} slides (was {before})")
        prs.remove_slide(prs.slide_count - 1)
        print(f"After remove_slide:       {prs.slide_count} slides")

        # Verify metadata
        meta = prs.get_core_properties()
        print(f"Title  : {meta.get('title', '')!r}")
        print(f"Author : {meta.get('creator', '')!r}")

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    # Verify sections survive round-trip
    with Presentation(str(output_path)) as prs2:
        sections = prs2.get_sections()
        print(f"Round-trip sections: {len(sections)}")
        for sec in sections:
            print(f"  - {sec.get('name', '?')!r}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_section, get_sections, add_comment, add_author,")
    print("  set_core_properties, get_core_properties, duplicate_slide, remove_slide")


if __name__ == "__main__":
    main()
