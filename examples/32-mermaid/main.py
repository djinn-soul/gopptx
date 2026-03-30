"""Mermaid diagram example in Python.

Renders Mermaid flowchart/sequence/pie diagrams into slides.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "32_mermaid.pptx"

    with Presentation.new("Mermaid Diagrams (Python)") as prs:
        flow = prs.slides[0]
        flow.title = "Mermaid Flowchart"
        flow.add_mermaid(
            "flowchart LR\nA[Start] --> B{Decision}\nB --> C[Action]\nC --> D[End]"
        )

        sequence = prs.add_slide("Mermaid Sequence")
        sequence.add_mermaid("sequenceDiagram\nAlice->>Bob: Hello\nBob-->>Alice: Hi")

        pie = prs.add_slide("Mermaid Pie")
        pie.add_mermaid(
            'pie title Browser Share\n"Chrome" : 65\n"Firefox" : 15\n"Safari" : 12\n"Edge" : 8'
        )

        prs.save(str(output_path))

    print(f"Created: {output_path}")


if __name__ == "__main__":
    main()
