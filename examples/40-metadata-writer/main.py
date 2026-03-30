"""Demonstrate reading and updating core document metadata properties.

This example demonstrates:
- Creating a presentation with initial metadata
- Reading title, author, and other core properties with get_metadata()
- Updating all core properties via set_metadata()
- Verifying metadata persists after save
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create presentation demonstrating metadata writer functionality."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Initial Title") as prs:
        prs.set_metadata(title="Initial Title", author="Initial Creator")
        prs.add_slide("Metadata Base")

        # Read initial metadata
        meta = prs.get_metadata()
        print(f"Initial title: {meta.get('title', '')}")

        # Update all core properties
        prs.set_metadata(
            title="Updated Title",
            author="Updated Creator",
        )

        # Add an explanatory slide
        prs.add_bullet_slide(
            "Metadata Properties",
            [
                "title — presentation title shown in file browsers",
                "author / creator — document author",
                "subject — document subject/category",
                "description — longer summary text",
                "keywords — comma-separated search terms",
            ],
        )

        output_path = output_dir / "40_metadata_output.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

        # Verify updated metadata
        updated_meta = prs.get_metadata()
        print(f"Updated title: {updated_meta.get('title', '')}")

    print("\n=== SUMMARY ===")
    print("Metadata updated: Initial Title -> Updated Title")
    print("Core properties: title, author, subject, description, keywords")


if __name__ == "__main__":
    main()
