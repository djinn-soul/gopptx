"""Demonstrate embedding media (video, audio) into a PPTX presentation.

This example demonstrates:
- Adding a video file to a slide with add_video()
- Adding an audio file to a slide with add_audio()
- Falling back gracefully when media fixture files are not present
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches

_ASSETS_DIR = Path("examples/assets/23")
_VIDEO_PATH = _ASSETS_DIR / "sample.mp4"
_AUDIO_PATH = _ASSETS_DIR / "sample.mp3"


def _add_overview_slide(prs: Presentation) -> None:
    """Add an overview slide."""
    prs.add_bullet_slide(
        "Media Embed Demo",
        [
            "Video, audio, and OLE media can be embedded in PPTX.",
            "Provide real fixture files in examples/assets/23 to enable inserts.",
            "add_video(slide_idx, path, x, y, w, h)",
            "add_audio(slide_idx, path, x, y, w, h)",
        ],
    )


def _try_add_video(prs: Presentation) -> None:
    """Add a video slide if the fixture file exists."""
    prs.add_slide("Video Embed")
    idx = prs.slide_count - 1

    if _VIDEO_PATH.exists():
        prs.add_video(
            idx,
            str(_VIDEO_PATH),
            (Inches(1), Inches(1.5), Inches(8), Inches(4.5)),
        )
        print(f"  Added video from: {_VIDEO_PATH}")
    else:
        prs.add_textbox(
            idx,
            Inches(1),
            Inches(2),
            Inches(8),
            Inches(1),
            text=f"Video fixture not found: {_VIDEO_PATH}",
        )
        print(f"  Video fixture missing ({_VIDEO_PATH}); placeholder added.")


def _try_add_audio(prs: Presentation) -> None:
    """Add an audio slide if the fixture file exists."""
    prs.add_slide("Audio Embed")
    idx = prs.slide_count - 1

    if _AUDIO_PATH.exists():
        prs.add_audio(
            idx,
            str(_AUDIO_PATH),
            (Inches(4.5), Inches(1.5), Inches(1.8), Inches(0.9)),
        )
        print(f"  Added audio from: {_AUDIO_PATH}")
    else:
        prs.add_textbox(
            idx,
            Inches(1),
            Inches(2),
            Inches(8),
            Inches(1),
            text=f"Audio fixture not found: {_AUDIO_PATH}",
        )
        print(f"  Audio fixture missing ({_AUDIO_PATH}); placeholder added.")


def main() -> None:
    """Create a presentation demonstrating media embedding."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Media Embed Demo") as prs:
        _add_overview_slide(prs)
        _try_add_video(prs)
        _try_add_audio(prs)

        output_path = output_dir / "23-media-embed.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated media embedding:")
    print("  - Video with add_video()")
    print("  - Audio with add_audio()")
    print("  - Graceful fallback when fixtures are missing")


if __name__ == "__main__":
    main()
