"""VBA macros example in Python.

Embeds a VBA project blob and saves a macro-enabled presentation.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "26_vba_macros.pptm"

    vba_blob = b"dummy vba project data for api smoke"

    with Presentation.new("VBA Macros (Python)") as prs:
        prs.update_slide(
            0,
            title="VBA Macro Example",
            bullets=[
                "Presentation has embedded VBA project data",
                "Saved with .pptm extension",
            ],
        )
        prs.add_vba_project(vba_blob)
        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
