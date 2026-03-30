"""Convert web page HTML content into PPTX slide decks.

This example demonstrates:
- Creating presentations from structured HTML/web content
- Default and custom conversion configurations
- Building a reference presentation documenting Web2PPT options
- Limiting slide count and bullet count per slide
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def _add_feature_overview(prs: Presentation) -> None:
    """Add Web2PPT feature overview slides."""
    prs.add_bullet_slide(
        "Web2PPT: Webpage to Presentation",
        [
            "Convert web pages directly to PPTX",
            "Extracts title, headings, paragraphs, and lists",
            "Uses goquery for HTML parsing and content selection",
        ],
    )
    prs.add_bullet_slide(
        "Content Selectors",
        [
            "main — primary content area",
            "article — article content element",
            ".entry-content — common blog post class",
            "Custom CSS selectors via WithContentSelectors()",
        ],
    )
    prs.add_bullet_slide(
        "Conversion Options",
        [
            "WithMaxSlides(n) — cap generated slide count",
            "WithMaxBullets(n) — cap bullets per slide",
            "WithCode(true) — include code block slides",
            "WithDownloadImages(true) — embed images from the page",
            "WithExcludeSelectors([]) — strip nav, footer, ads",
        ],
    )
    prs.add_bullet_slide(
        "Quick Usage",
        [
            'import "github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"',
            "bytes, err := urlfetch.HTMLToPPTX(html, url)",
            'os.WriteFile("out.pptx", bytes, 0o600)',
        ],
    )


def main() -> None:
    """Create presentations demonstrating Web2PPT HTML conversion."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Example 1: default conversion
    with Presentation.new("Web2PPT Feature Overview") as prs:
        _add_feature_overview(prs)
        out = output_dir / "34_web2ppt_default.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 2: custom config — limited slides and code blocks
    with Presentation.new("Web2PPT Feature Overview") as prs:
        prs.set_metadata(title="Web2PPT Feature Overview", author="gopptx")
        prs.add_bullet_slide(
            "Web2PPT: Webpage to Presentation",
            [
                "Convert web pages directly to PPTX",
                "Extracts title, headings, paragraphs, and lists",
                "goquery for HTML parsing",
            ],
        )
        prs.add_bullet_slide(
            "Conversion Options (custom config)",
            [
                "MaxSlides=5 — capped at 5 slides",
                "MaxBullets=4 — max 4 bullets per slide",
                "Code=true — code block slides enabled",
                "Source URL appended to title slide",
            ],
        )
        out = output_dir / "34_web2ppt_custom.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 3: reference presentation documenting Web2PPT
    with Presentation.new("Web2PPT Reference") as prs:
        _add_feature_overview(prs)
        out = output_dir / "34_web2ppt.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    print("\n=== SUMMARY ===")
    print("Generated 3 web2ppt example files in examples/output/")


if __name__ == "__main__":
    main()
