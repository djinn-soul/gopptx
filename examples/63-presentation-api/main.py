"""Demonstrate core Presentation API operations.

This example demonstrates:
- Creating presentations with Presentation.new()
- Opening existing presentations with Presentation.open()
- Reading and writing metadata (title, author, subject, keywords, etc.)
- Saving with save() and verifying a round-trip
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_overview_slides(prs: Presentation) -> None:
    """Add overview and core-properties slides."""
    prs.add_bullet_slide(
        "Overview",
        [
            "Create presentations with Presentation.new()",
            "Open existing files with Presentation.open()",
            "Save with prs.save(path)",
            "Manage core properties: title, author, subject, keywords",
        ],
    )
    prs.add_bullet_slide(
        "Core Properties",
        [
            "Title, Subject, Author",
            "Keywords, Description, Category",
            "get_metadata() / set_metadata()",
            "Metadata survives save/open round-trips",
        ],
    )


def main() -> None:
    """Create presentation demonstrating core Presentation API operations."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "63_presentation_api.pptx"

    # --- Part 1: Create and populate a presentation ---
    with Presentation.new("Presentation API Demo") as prs:
        prs.set_metadata(
            title="Presentation API Demo",
            author="Python Developer",
        )
        _add_overview_slides(prs)
        prs.save(str(output_path))
        print(f"Created initial presentation: {output_path}")

    # --- Part 2: Open, read, and update metadata ---
    with Presentation.open(str(output_path)) as prs:
        meta = prs.get_metadata()
        print(f"Slide count  : {prs.slide_count}")
        print(f"Title        : {meta.get('title', '')!r}")
        print(f"Author       : {meta.get('author', '')!r}")

        prs.set_metadata(
            title="Updated Presentation API Demo",
            author="Python Developer",
        )
        prs.save(str(output_path))
        print("Saved with updated metadata.")

    # --- Part 3: Verify round-trip ---
    with Presentation.open(str(output_path)) as prs:
        meta = prs.get_metadata()
        print(f"Verified title  : {meta.get('title', '')!r}")
        print(f"Verified author : {meta.get('author', '')!r}")
        size = prs.get_slide_size()
        print(f"Slide size      : {size}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: new(), open(), save(), get_metadata(), set_metadata()")
    print(f"Output: {output_path}")


if __name__ == "__main__":
    main()
