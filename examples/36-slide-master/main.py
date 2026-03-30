"""Document multi-master slide styling concepts in Python.

This example demonstrates:
- Slides that describe two visual theme families (blue and warm)
- Conceptual title/body styling notes per master family
- Multi-level bullet indentation patterns
- Guidance on master assignment behavior in the Go API
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_blue_master_slides(prs: Presentation) -> None:
    """Add slides styled with the blue (tech) master."""
    prs.add_bullet_slide(
        "Slide 1 (Blue Master)",
        [
            "Bullet Level 1",
            "  Bullet Level 2 (sub-item)",
            "Background: light blue (#E3F2FD)",
        ],
    )
    prs.add_bullet_slide(
        "Slide 3 (Blue Master)",
        [
            "Arial title font, 44pt",
            "Body color: dark navy (#1A237E)",
            "Tech-focused color palette",
        ],
    )


def _add_warm_master_slides(prs: Presentation) -> None:
    """Add slides styled with the warm (Calibri) master."""
    prs.add_bullet_slide(
        "Slide 2 (Warm Master)",
        [
            "Second master visual family",
            "Background: warm orange (#FFF3E0)",
            "Calibri title font in deep orange",
        ],
    )
    prs.add_bullet_slide(
        "Slide 4 (Warm Master)",
        [
            "Calibri font, 44pt bold title",
            "Body color: burnt orange (#E65100)",
            "Warm enterprise color palette",
        ],
    )


def main() -> None:
    """Create a conceptual presentation for multi-master styling patterns."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Multi-Master Demo") as prs:
        _add_blue_master_slides(prs)
        _add_warm_master_slides(prs)

        # Reference slide explaining the master pattern
        prs.add_bullet_slide(
            "About Slide Masters",
            [
                "Masters define background, font, and color for a group of slides",
                "Multiple masters allow distinct visual families in one deck",
                "Assign slides to masters via layout index",
                "Title style: Font, size, bold, and color",
                "Body style: Multi-level indent and color",
            ],
        )

        output_path = output_dir / "36_multi_master_smoke.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("5 slides across two visual master families (blue and warm)")


if __name__ == "__main__":
    main()
