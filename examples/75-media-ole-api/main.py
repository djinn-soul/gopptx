"""Demonstrate media embedding and document OLE embedding patterns.

This example demonstrates:
- Embedding audio from raw bytes with add_audio()
- Embedding video from a file path with add_video() (requires fixture)
- Reference slides covering OLE embedding concepts
- Graceful skip when optional fixture files are not present
"""

from __future__ import annotations

import tempfile
from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

# Minimal valid WAV (PCM mono, 8-bit, 8000 Hz, 1 sample).
_TINY_WAV = bytes([
    ord("R"),
    ord("I"),
    ord("F"),
    ord("F"),
    0x25,
    0x00,
    0x00,
    0x00,
    ord("W"),
    ord("A"),
    ord("V"),
    ord("E"),
    ord("f"),
    ord("m"),
    ord("t"),
    ord(" "),
    0x10,
    0x00,
    0x00,
    0x00,
    0x01,
    0x00,
    0x01,
    0x00,
    0x40,
    0x1F,
    0x00,
    0x00,
    0x40,
    0x1F,
    0x00,
    0x00,
    0x01,
    0x00,
    0x08,
    0x00,
    ord("d"),
    ord("a"),
    ord("t"),
    ord("a"),
    0x01,
    0x00,
    0x00,
    0x00,
    0x80,
])

ASSETS_DIR = Path("examples/assets/75")


def _add_overview_slide(prs: Presentation) -> int:
    """Overview slide listing media embedding methods."""
    return prs.add_bullet_slide(
        "Media & OLE Embedding",
        [
            "Audio from bytes  - add_audio(slide_idx, path, x, y, w, h)",
            "Video from file   - add_video(slide_idx, path, x, y, w, h)",
            "OLE objects       - available via PresentationEditor (Go API)",
            "Fixture files in  examples/assets/75/ for video/OLE demos",
        ],
    )


def _add_audio_slide(prs: Presentation, wav_path: Path) -> None:
    """Embed audio from a temp WAV file."""
    idx = prs.add_slide("Audio Embedding", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_textbox(
        idx,
        Inches(0.5),
        Inches(1.2),
        Inches(9),
        Inches(0.6),
        text=f"Audio source: {wav_path.name}  (tiny WAV, 1 sample)",
    )
    prs.add_audio(
        idx,
        str(wav_path),
        Inches(4.5),
        Inches(1.5),
        Inches(1.8),
        Inches(0.9),
    )
    print(f"  Inserted WAV audio from bytes (via temp file: {wav_path})")


def _add_video_slide(prs: Presentation) -> None:
    """Embed video if fixture files exist, otherwise log a skip notice."""
    video_path = ASSETS_DIR / "sample.mp4"
    if not video_path.exists():
        prs.add_bullet_slide(
            "Video Embedding (fixture not found)",
            [
                f"Expected fixture: {video_path}",
                "Place sample.mp4 in examples/assets/75/ to enable this demo.",
                "add_video(slide_idx, path, x, y, w, h)",
            ],
        )
        print(f"  Video fixture not found ({video_path}); skipping.")
        return

    idx = prs.add_slide("Video Embedding", layout=SlideLayoutType.TITLE_ONLY)
    prs.add_video(
        idx,
        str(video_path),
        Inches(0.6),
        Inches(1.4),
        Inches(3.6),
        Inches(2.1),
    )
    print(f"  Inserted video from fixture: {video_path}")


def _add_ole_reference_slide(prs: Presentation) -> None:
    """Reference slide for OLE object embedding concepts."""
    prs.add_bullet_slide(
        "OLE Object Embedding (Go API)",
        [
            "OLE embedding requires the PresentationEditor API in Go.",
            "ed.AddOLEObject(slideIdx, data, icon, progID, x, y, cx, cy)",
            "Supported ProgIDs: Excel.Sheet.12, Word.Document.12, etc.",
            "Place sample_ole.bin + poster.png in examples/assets/75/ for Go demo.",
        ],
    )


def main() -> None:
    """Create presentation demonstrating media and OLE embedding."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "75_media_ole_api.pptx"

    # Write tiny WAV to a temp file so add_audio() (file-path API) can read it.

    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp:
        tmp.write(_TINY_WAV)
        wav_path = Path(tmp.name)

    try:
        with Presentation.new("Media & OLE API Demo") as prs:
            _add_overview_slide(prs)
            _add_audio_slide(prs, wav_path)
            _add_video_slide(prs)
            _add_ole_reference_slide(prs)

            prs.save(str(output_path))
            print(f"Saved: {output_path}")
    finally:
        wav_path.unlink(missing_ok=True)

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_audio (from bytes via temp file),")
    print("  add_video (fixture-conditional), OLE concept reference slide")


if __name__ == "__main__":
    main()
