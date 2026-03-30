"""Demonstrate presentation-protection concepts from Python.

This example demonstrates:
- Creating a presentation with confidential-content slides
- Documenting protection features available in the Go API
- Referencing password-to-modify / mark-as-final concepts
- Preparing a deck you can use for manual follow-up validation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    """Create a conceptual protection demo deck (no protection applied in Python)."""
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
                "Go API: WithModifyPassword - prevent edits without a password",
                "Go API: WithMarkAsFinal - display read-only banner in PowerPoint",
                "Go API: WithSignaturesEnabled - enable digital signature support",
                "This Python example documents the feature set only",
            ],
        )

        output_path = output_dir / "46_presentation_protection.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Conceptual presentation-protection overview generated.")
    print("Note: this Python script does not apply protection flags.")
    print("Use Go protection APIs for enforceable protection metadata.")


if __name__ == "__main__":
    main()
