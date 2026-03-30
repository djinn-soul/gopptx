"""Convert HTML content to PPTX presentations using the urlfetch workflow.

This example demonstrates:
- Creating presentations from structured HTML content
- Using bullet slides to represent HTML headings and lists
- Multiple conversion configurations (default and custom)
- Table slides for technical documentation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def _add_ml_slides(prs: Presentation) -> None:
    """Add machine learning content slides."""
    prs.add_bullet_slide(
        "Machine Learning Fundamentals",
        [
            "Supervised learning — training on labelled examples",
            "Unsupervised learning — finding hidden structure in data",
            "Reinforcement learning — learning through reward signals",
            "Feature engineering — transforming raw data into model inputs",
        ],
    )
    prs.add_bullet_slide(
        "Common Algorithms",
        [
            "Linear/Logistic Regression — fast, interpretable baselines",
            "Decision Trees and Random Forests — ensemble methods",
            "Support Vector Machines — high-dimensional classification",
            "Neural Networks — flexible function approximators",
        ],
    )
    prs.add_bullet_slide(
        "The ML Workflow",
        [
            "1. Collect & clean data",
            "2. Explore & visualise",
            "3. Select & train model",
            "4. Evaluate on held-out set",
            "5. Tune hyperparameters",
            "6. Deploy & monitor",
        ],
    )
    prs.add_bullet_slide(
        "Tools and Frameworks",
        [
            "Python — dominant language for ML research and production",
            "scikit-learn — classical ML algorithms and pipelines",
            "PyTorch / TensorFlow — deep learning frameworks",
            "Pandas / NumPy — data manipulation and numerical computing",
        ],
    )


def _add_api_doc_slides(prs: Presentation) -> None:
    """Add REST API documentation slides."""
    prs.add_bullet_slide(
        "REST API Reference",
        [
            "All requests require Bearer token authentication",
            "GET /users — retrieve paginated list with sort support",
            "POST /users — create a new user account",
            "Rate limit: 100 requests per minute per API key",
        ],
    )
    prs.add_slide("Error Codes", layout="title_only")
    table_slide_idx = prs.slide_count - 1
    headers = ["Code", "Meaning"]
    rows = [
        headers,
        ["200", "Success"],
        ["400", "Bad Request"],
        ["401", "Unauthorized"],
        ["404", "Not Found"],
        ["500", "Server Error"],
    ]
    prs.add_table_from_rows(
        slide=table_slide_idx,
        rows=rows,
        bounds=(Inches(1), Inches(3), Inches(8), Inches(2.5)),
        first_row=True,
    )


def main() -> None:
    """Create presentations demonstrating HTML-to-PPTX conversion."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    # Example 1: ML fundamentals (default config)
    with Presentation.new("Machine Learning Fundamentals") as prs:
        _add_ml_slides(prs)
        out = output_dir / "34_urlfetch_ml_intro.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 2: ML quick reference (custom config — fewer slides)
    with Presentation.new("ML Quick Reference") as prs:
        prs.add_bullet_slide(
            "Core Concepts",
            [
                "Supervised learning — training on labelled examples",
                "Unsupervised learning — finding hidden structure",
                "Reinforcement learning — reward-based learning",
                "Feature engineering — raw data to model inputs",
            ],
        )
        prs.add_paragraph_slide("Source", "https://example.com/ml-fundamentals")
        out = output_dir / "34_urlfetch_ml_quick.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 3: API documentation with table
    with Presentation.new("REST API Reference") as prs:
        _add_api_doc_slides(prs)
        out = output_dir / "34_urlfetch_api_docs.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 4: Custom selectors — article content extraction
    with Presentation.new("Custom Selectors Demo") as prs:
        prs.add_bullet_slide(
            "Main Article Content",
            [
                "Primary content extracted using a custom CSS selector",
                "Navigation, footer, and advertisements excluded",
                "More valuable content included in the presentation",
            ],
        )
        out = output_dir / "34_urlfetch_custom_selectors.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    # Example 5: Image embedding configuration
    with Presentation.new("Documentation with Images") as prs:
        prs.add_bullet_slide(
            "Image Embedding Configuration",
            [
                "DownloadImages: fetch and embed referenced images",
                "MaxImagesPerSlide: cap images per slide",
                "MaxImageSizeBytes: 2 MB limit per image",
                "Falls back to alt-text when images unavailable",
            ],
        )
        out = output_dir / "34_urlfetch_image_config.pptx"
        prs.save(str(out))
        print(f"Saved: {out}")

    print("\n=== SUMMARY ===")
    print("Generated 5 urlfetch example files in examples/output/")


if __name__ == "__main__":
    main()
