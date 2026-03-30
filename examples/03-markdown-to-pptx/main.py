"""Convert Markdown content to PPTX slides using add_slide_from_markdown.

This example demonstrates:
- Converting inline Markdown text into slides with add_slide_from_markdown()
- Each top-level heading becomes a slide title
- Bullet lists become slide bullet points
- Creating two separate presentations from different Markdown sources
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation

_MARKDOWN_COMPLEX = """\
# Project Overview

- Modern CI/CD pipeline with automated testing
- Containerised microservices on Kubernetes
- Real-time telemetry and alerting

# Architecture

- API Gateway handles all inbound traffic
- Service mesh manages inter-service communication
- Data lake aggregates events for analytics

# Roadmap

- Q1: Baseline infrastructure migration
- Q2: Feature parity with legacy system
- Q3: New analytics dashboard
- Q4: Global multi-region rollout
"""

_MARKDOWN_LINKS = """\
# Getting Started

- Clone the repository
- Run the setup script
- Open the dashboard

# Key Resources

- Documentation site
- API reference
- Community forum
- Issue tracker
"""


def _build_from_markdown(prs: Presentation, markdown: str) -> None:
    """Add slides parsed from a Markdown string."""
    for line in markdown.strip().split("\n"):
        if line.startswith("# "):
            prs.add_slide_from_markdown(line)
        else:
            prs.add_slide_from_markdown(line)


def main() -> None:
    """Create two presentations generated from Markdown content."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Presentation 1: complex Markdown with multiple headings
    with Presentation.new("Markdown Complex") as prs:
        prs.add_slide_from_markdown(_MARKDOWN_COMPLEX)
        path1 = output_dir / "03-markdown-complex.pptx"
        prs.save(str(path1))
        print(f"Saved: {path1}")

    # Presentation 2: links and gallery style Markdown
    with Presentation.new("Markdown Links Gallery") as prs:
        prs.add_slide_from_markdown(_MARKDOWN_LINKS)
        path2 = output_dir / "03-markdown-links-gallery.pptx"
        prs.save(str(path2))
        print(f"Saved: {path2}")

    print("\n=== SUMMARY ===")
    print("Converted Markdown text into two PPTX presentations.")
    print("Each top-level heading becomes a new slide.")


if __name__ == "__main__":
    main()
