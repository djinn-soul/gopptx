"""Demonstrate section management: adding and organizing slides into sections.

This example demonstrates:
- Creating a multi-slide presentation
- Adding named sections that group slide indices
- Retrieving sections and verifying structure
- Sections visible in PowerPoint Slide Sorter view
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _build_slides(prs: Presentation) -> None:
    """Build the base slide content."""
    prs.add_title_slide("Intro Slide")
    prs.add_bullet_slide("Detail 1", ["First detail point", "Supporting information"])
    prs.add_bullet_slide("Detail 2", ["Second detail point", "More supporting info"])
    prs.add_bullet_slide("Appendix", ["Reference material", "Additional data"])


def main() -> None:
    """Create presentation demonstrating section management."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)
    output_path = output_dir / "44_section_management.pptx"

    with Presentation.new("Section Demo") as prs:
        _build_slides(prs)

        # Add sections grouping slides
        prs.add_section("Introduction", [0])
        prs.add_section("Main Content", [1, 2])
        prs.add_section("Appendix", [3])

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    with Presentation(str(output_path)) as prs:
        sections = prs.sections
        print(f"Created {len(sections)} sections:")
        for section in sections:
            print(
                f"  - {section.get('Name', section.get('name', '?'))}"
                f" -> slide IDs {section.get('SlideIDs', section.get('slide_ids', []))}"
            )

    print("\n=== SUMMARY ===")
    print("3 sections: Introduction, Main Content, Appendix")
    print("Tip: Open in PowerPoint Slide Sorter view to see sections.")


if __name__ == "__main__":
    main()
