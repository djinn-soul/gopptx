"""Build a quick-reference guide for the gopptx CLI subcommands.

This example demonstrates:
- Creating a multi-slide reference presentation in code
- Using add_bullet_slide() to document CLI commands
- No CLI execution is performed — this generates a PPTX reference guide
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_overview_slide(prs: Presentation) -> None:
    """Add a title and overview slide."""
    prs.add_bullet_slide(
        "gopptx CLI Reference",
        [
            "The gopptx binary exposes several subcommands for common PPTX tasks.",
            "Each subsequent slide documents one subcommand.",
            "Run `gopptx --help` to see all available commands.",
        ],
    )


def _add_create_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'create' subcommand."""
    prs.add_bullet_slide(
        "gopptx create",
        [
            "Creates a new blank PPTX presentation.",
            'Usage:  gopptx create --title "My Deck" --slides 5 -o out.pptx',
            "Flags:",
            "  --title   Presentation title (required)",
            "  --slides  Number of blank slides to generate",
            "  -o        Output file path",
        ],
    )


def _add_md2ppt_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'md2ppt' subcommand."""
    prs.add_bullet_slide(
        "gopptx md2ppt",
        [
            "Converts a Markdown file into a PPTX presentation.",
            "Usage:  gopptx md2ppt input.md -o out.pptx",
            "Each top-level heading becomes a slide title.",
            "Bullet lists under a heading become slide bullets.",
            "Code blocks are rendered as monospace text boxes.",
        ],
    )


def _add_info_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'info' subcommand."""
    prs.add_bullet_slide(
        "gopptx info",
        [
            "Displays metadata and structural information about a PPTX file.",
            "Usage:  gopptx info presentation.pptx",
            "Output includes:",
            "  - Title, author, and creation date",
            "  - Slide count and slide sizes",
            "  - Embedded media and chart counts",
        ],
    )


def _add_validate_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'validate' subcommand."""
    prs.add_bullet_slide(
        "gopptx validate",
        [
            "Validates the structural integrity of a PPTX file.",
            "Usage:  gopptx validate presentation.pptx",
            "Checks for:",
            "  - Missing or malformed XML parts",
            "  - Broken relationship references",
            "  - Invalid media or chart entries",
            "Exits with a non-zero code if issues are found.",
        ],
    )


def _add_merge_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'merge' subcommand."""
    prs.add_bullet_slide(
        "gopptx merge",
        [
            "Merges two or more PPTX files into a single presentation.",
            "Usage:  gopptx merge a.pptx b.pptx -o merged.pptx",
            "Slides are appended in the order the input files are given.",
            "Themes and masters from the first file are preserved.",
            "Assets (images, charts) are deduplicated automatically.",
        ],
    )


def _add_version_slide(prs: Presentation) -> None:
    """Add a slide documenting the 'version' subcommand."""
    prs.add_bullet_slide(
        "gopptx version",
        [
            "Prints the current gopptx version and build information.",
            "Usage:  gopptx version",
            "Output example:",
            "  gopptx v1.2.3 (commit abc1234, built 2025-01-01)",
            "Useful for confirming the installed binary in CI pipelines.",
        ],
    )


def main() -> None:
    """Create a presentation documenting the gopptx CLI subcommands."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("gopptx CLI Reference") as prs:
        _add_overview_slide(prs)
        _add_create_slide(prs)
        _add_md2ppt_slide(prs)
        _add_info_slide(prs)
        _add_validate_slide(prs)
        _add_merge_slide(prs)
        _add_version_slide(prs)

        output_path = output_dir / "15-cli-commands.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created a 7-slide CLI reference guide covering:")
    print("  create, md2ppt, info, validate, merge, version")


if __name__ == "__main__":
    main()
