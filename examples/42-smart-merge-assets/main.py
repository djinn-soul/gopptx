"""Demonstrate smart merging of assets between presentations.

This example demonstrates:
- Creating a source presentation with an embedded image
- Creating a target presentation
- Merging slides and preserving image assets across presentations
- Duplicate_slide to copy a slide with all its assets
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches

# Minimal 1x1 PNG used as a source image asset
_PNG_DATA = bytes([
    0x89,
    0x50,
    0x4E,
    0x47,
    0x0D,
    0x0A,
    0x1A,
    0x0A,
    0x00,
    0x00,
    0x00,
    0x0D,
    0x49,
    0x48,
    0x44,
    0x52,
    0x00,
    0x00,
    0x00,
    0x01,
    0x00,
    0x00,
    0x00,
    0x01,
    0x08,
    0x06,
    0x00,
    0x00,
    0x00,
    0x1F,
    0x15,
    0xC4,
    0x89,
    0x00,
    0x00,
    0x00,
    0x0A,
    0x49,
    0x44,
    0x41,
    0x54,
    0x78,
    0x9C,
    0x63,
    0x00,
    0x01,
    0x00,
    0x00,
    0x05,
    0x00,
    0x01,
    0x0D,
    0x0A,
    0x2D,
    0xB4,
    0x00,
    0x00,
    0x00,
    0x00,
    0x00,
    0x49,
    0x45,
    0x4E,
    0x44,
    0xAE,
    0x42,
    0x60,
    0x82,
])


def _build_source_deck(prs: Presentation, image_path: Path) -> None:
    """Build a source deck with image assets."""
    slide = prs.add_bullet_slide(
        "Source Slide with Image",
        ["This slide contains an embedded PNG asset."],
    )
    slide.add_image(
        str(image_path), (Inches(0.9), Inches(1.5), Inches(1.2), Inches(1.2))
    )


def _build_target_deck(prs: Presentation) -> None:
    """Build a target deck to receive merged slides."""
    prs.add_title_slide("Title Slide")
    prs.add_bullet_slide(
        "Smart Merge Assets",
        [
            "Merge slides from source into target presentation",
            "Image assets are carried along with merged slides",
            "Deduplication prevents duplicate image files",
        ],
    )


def main() -> None:
    """Create presentation demonstrating smart asset merge."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with tempfile.TemporaryDirectory() as tmpdir:
        tmpdir_path = Path(tmpdir)
        image_path = tmpdir_path / "source_asset.png"
        image_path.write_bytes(_PNG_DATA)

        source_path = tmpdir_path / "42_smart_merge_assets_source.pptx"
        with Presentation.new("Source Deck") as source:
            _build_source_deck(source, image_path)
            source.save(str(source_path))

        with Presentation.new("Target Deck") as prs:
            _build_target_deck(prs)
            prs.merge_from_file(str(source_path))

            merged_idx = prs.slide_count - 1
            image_refs = prs.list_slide_images(merged_idx)

            prs.add_bullet_slide(
                "Merge Features",
                [
                    "merge_from_file() imports full slides from another deck",
                    "Embedded images remain available after the merge",
                    f"Merged slide image count: {len(image_refs)}",
                    "Use duplicate_slide_after() to clone merged content if needed",
                ],
            )

            output_path = output_dir / "42_smart_merge_assets.pptx"
            prs.save(str(output_path))
            print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Target deck with merged slide and preserved image asset")


if __name__ == "__main__":
    main()
