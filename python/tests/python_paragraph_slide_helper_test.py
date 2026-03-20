"""Smoke tests for the high-level paragraph slide helper."""

from __future__ import annotations

from gopptx import Presentation


def test_add_paragraph_slide_uses_default_textbox_bounds() -> None:
    """Paragraph slide helper should create a text-bearing shape by default."""
    paragraph = (
        "Paragraph text helps explain context clearly. "
        "Use one concise block for intent."
    )

    with Presentation.new("Paragraph Helper") as prs:
        slide = prs.add_paragraph_slide("Paragraph Text", paragraph)
        texts = [str(shape.get("Text", "")) for shape in slide.list_shapes()]
        if not any(paragraph in text for text in texts):
            raise AssertionError(f"expected paragraph text shape, got: {texts!r}")


def test_add_paragraph_allows_multiple_blocks_on_same_slide() -> None:
    """Slide-level helper should append a second paragraph on the same slide."""
    first = "First paragraph block."
    second = "Second paragraph block."

    with Presentation.new("Paragraph Helper Multi") as prs:
        slide = prs.add_paragraph_slide("Paragraph Text", first)
        slide.add_paragraph(second)
        texts = [str(shape.get("Text", "")) for shape in slide.list_shapes()]
        if not any(first in text for text in texts):
            raise AssertionError(f"missing first paragraph in shapes: {texts!r}")
        if not any(second in text for text in texts):
            raise AssertionError(f"missing second paragraph in shapes: {texts!r}")
