"""Demonstrate presentation protection: modify password and mark-as-final.

This example demonstrates:
- Creating a presentation with confidential content
- Adding password-to-modify protection
- Marking a presentation as final (read-only banner in PowerPoint)
- Protection metadata embedded in the PPTX file
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create a protected presentation with password and mark-as-final."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Protected Presentation") as prs:
        prs.add_title_slide("Confidential Content")

        prs.add_bullet_slide(
            "Features",
            [
                "Password to Modify: 'test'",
                "Marked as Final",
                "Digital signatures enabled",
            ],
        )

        prs.add_bullet_slide(
            "Protection Features",
            [
                "WithModifyPassword — prevent edits without a password",
                "WithMarkAsFinal — display read-only banner in PowerPoint",
                "WithSignaturesEnabled — enable digital signature support",
                "Protection is stored in presentation.xml and app.xml",
            ],
        )

        output_path = output_dir / "46_presentation_protection.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Protected presentation with password='test' and mark-as-final.")
    print("Verify in PowerPoint:")
    print("  1. Should prompt for password 'test' to modify.")
    print("  2. Should show a 'Marked as Final' banner.")


if __name__ == "__main__":
    main()
