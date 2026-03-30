"""Notes master example in Python.

Configures global notes master settings and writes slides with notes.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "33_notes_master.pptx"

    with Presentation.new("Notes Master (Python)") as prs:
        prs.update_notes_master(
            header="CONFIDENTIAL - internal use only",
            footer="Notes Master Smoke Test",
            show_date_time=True,
            show_slide_num=True,
        )
        prs.update_slide(0, title="Notes Master Demo", bullets=["Slide 1 body content"])
        prs.set_notes(
            0, "This is level 1 notes text.\nThis is additional notes content."
        )
        second = prs.add_slide("Another Slide", bullets=["Slide 2 body content"])
        second.notes = "Check the header and footer on the notes page."
        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
