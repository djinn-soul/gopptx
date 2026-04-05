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
            "get_core_properties() / set_core_properties()",
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
        props = prs.get_core_properties()
        props["title"] = "Presentation API Demo"
        props["creator"] = "Python Developer"
        prs.set_core_properties(props)
        _add_overview_slides(prs)
        prs.save(str(output_path))
        print(f"Created initial presentation: {output_path}")

    # --- Part 2: Open, read, and update metadata ---
    with Presentation(str(output_path)) as prs:
        meta = prs.get_core_properties()
        print(f"Slide count  : {prs.slide_count}")
        print(f"Title        : {meta.get('title', '')!r}")
        print(f"Author       : {meta.get('creator', '')!r}")

        props = prs.get_core_properties()
        props["title"] = "Updated Presentation API Demo"
        props["creator"] = "Python Developer"
        prs.set_core_properties(props)
        prs.save(str(output_path))
        print("Saved with updated metadata.")

    # --- Part 3: Verify round-trip ---
    with Presentation(str(output_path)) as prs:
        meta = prs.get_core_properties()
        print(f"Verified title  : {meta.get('title', '')!r}")
        print(f"Verified author : {meta.get('creator', '')!r}")
        print("Slide size      : use SIZE_16X9_* constants with set_slide_size(...)")

    print("\n=== SUMMARY ===")
    print(
        "Demonstrated: new(), open(), save(), get_core_properties(), set_core_properties()"
    )
    print(f"Output: {output_path}")


if __name__ == "__main__":
    main()
