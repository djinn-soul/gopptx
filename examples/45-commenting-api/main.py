"""Demonstrate the commenting API: adding authors and review comments to slides.

This example demonstrates:
- Registering comment authors with names and initials
- Adding positioned comments from multiple authors on different slides
- Comments support a review workflow with author attribution
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create presentation demonstrating the commenting API."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Comments API Smoke Test") as prs:
        # Slide 1: team review content
        prs.add_bullet_slide(
            "Slide 1 (Team Review)",
            ["Team member A should comment here."],
        )
        # Slide 2: client review content
        prs.add_bullet_slide(
            "Slide 2 (Client Review)",
            ["Client feedback wanted here."],
        )

        # Register authors
        author_a = prs.add_author("Alice Reviewer", "AR")
        author_b = prs.add_author("Bob Architect", "BA")
        print(f"Registered authors: Alice (ID={author_a}), Bob (ID={author_b})")

        # Add comments to Slide 1 from both authors
        prs.add_comment(
            0, author_a, "Looks good, but check the font size.", 500000, 500000
        )
        prs.add_comment(0, author_b, "Agreed. Let's make it 24pt.", 600000, 600000)
        print("Added comments to Slide 1 from Alice and Bob")

        # Add comment to Slide 2 from Alice
        prs.add_comment(1, author_a, "Need comparison chart here.", 1000000, 1000000)
        print("Added comment to Slide 2 from Alice")

        # Explanation slide
        prs.add_bullet_slide(
            "Commenting API Features",
            [
                "add_author(name, initials) — register a comment author",
                "add_comment(slide, author_id, text, x, y) — place a comment",
                "Supports multi-author review workflows",
                "Comments visible in PowerPoint Review pane",
                "Positioned by EMU coordinates on the slide",
            ],
        )

        output_path = output_dir / "45_commenting_api.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("3 slides with comments from Alice Reviewer and Bob Architect")


if __name__ == "__main__":
    main()
