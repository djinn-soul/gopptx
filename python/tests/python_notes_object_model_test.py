"""Notes object-model parity tests."""

from __future__ import annotations

from gopptx import Presentation


def test_notes_shape_collection_and_body_proxy() -> None:
    """Notes slide exposes typed shape proxies and body text mutation."""
    with Presentation.new("Notes Object Model") as prs:
        slide = prs.add_slide("S1")
        slide.notes = "speaker notes"
        notes_slide = slide.notes_slide
        assert notes_slide is not None

        assert len(notes_slide.shapes) > 0
        first = notes_slide.shapes[0]
        assert isinstance(first.name, str)
        assert isinstance(first.placeholder_type, str)
        assert isinstance(first.is_placeholder, bool)
        assert isinstance(first.x, float)
        assert isinstance(first.y, float)
        assert isinstance(first.cx, float)
        assert isinstance(first.cy, float)
        assert notes_slide.shapes.get(name=first.name) is not None

        body = notes_slide.body_shape
        if body is not None:
            assert body.text == "speaker notes"
            body_tf = body.text_frame
            assert body_tf is not None
            assert body_tf.text == "speaker notes"
            body_tf.text = "updated via text frame"
            assert slide.notes == "updated via text frame"
            body.text = "updated notes"
            assert slide.notes == "updated notes"
            notes_tf = notes_slide.notes_text_frame
            assert notes_tf is not None
            notes_tf.text = "updated via notes_text_frame"
            assert slide.notes == "updated via notes_text_frame"
