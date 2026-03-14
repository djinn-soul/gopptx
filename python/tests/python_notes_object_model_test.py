"""Notes object-model parity tests."""

from __future__ import annotations

from gopptx import Presentation
from gopptx.slide.notes_slide import NotesShape


def test_notes_shape_collection_and_body_proxy() -> None:
    """Notes slide exposes typed shape proxies and body text mutation."""
    with Presentation.new("Notes Object Model") as prs:
        slide = prs.add_slide("S1")
        slide.notes = "speaker notes"
        notes_slide = slide.notes_slide
        assert notes_slide is not None

        assert len(notes_slide.shapes) > 0
        first = notes_slide.shapes[0]
        assert isinstance(first.shape_id, int)
        assert isinstance(first.name, str)
        assert isinstance(first.shape_type, str)
        assert isinstance(first.placeholder_type, str)
        assert isinstance(first.is_placeholder, bool)
        assert isinstance(first.has_text_frame, bool)
        assert isinstance(first.x, float)
        assert isinstance(first.y, float)
        assert isinstance(first.cx, float)
        assert isinstance(first.cy, float)
        assert isinstance(first.left, float)
        assert isinstance(first.top, float)
        assert isinstance(first.width, float)
        assert isinstance(first.height, float)
        assert notes_slide.shapes.get(name=first.name) is not None
        if first.shape_id >= 0:
            assert notes_slide.shapes.get(shape_id=first.shape_id) is not None
            assert notes_slide.shapes.by_id(first.shape_id) is not None
            assert notes_slide.shape(first.shape_id) is not None
        assert notes_slide.shapes.get(shape_type=first.shape_type) is not None
        assert isinstance(notes_slide.text_shapes, list)
        assert isinstance(notes_slide.placeholder_shapes, list)
        assert all(shape.has_text_frame for shape in notes_slide.text_shapes)
        body_matches = notes_slide.shapes.find_all(placeholder_type="body")
        if body_matches:
            assert body_matches[0].placeholder_type == "body"

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
            assert len(notes_tf.paragraphs) == 1
            para0 = notes_tf.paragraphs[0]
            assert para0.text == "updated via notes_text_frame"
            run0 = para0.runs[0]
            assert run0.text == "updated via notes_text_frame"
            run0.text = "updated via notes run"
            assert slide.notes == "updated via notes run"
            para0.runs.add_run(" + appended")
            assert slide.notes == "updated via notes run + appended"
            para0.add_run(" via paragraph")
            assert slide.notes == "updated via notes run + appended via paragraph"
            para1 = notes_tf.add_paragraph("second paragraph")
            assert para1.text == "second paragraph"
            assert (
                slide.notes
                == "updated via notes run + appended via paragraph\nsecond paragraph"
            )
            notes_tf.clear()
            assert not slide.notes


def test_notes_shape_non_placeholder_text_routes_to_shape_setter() -> None:
    class _DummyNotesSlide:
        def __init__(self) -> None:
            self.text = ""
            self.called_shape_id: int | None = None
            self.called_text: str | None = None

        def _set_shape_text(self, shape_id: int, text: str) -> None:
            self.called_shape_id = shape_id
            self.called_text = text

    dummy = _DummyNotesSlide()
    shape = NotesShape(
        dummy,  # type: ignore[arg-type]
        {
            "id": 42,
            "name": "Aux text",
            "type": "sp",
            "has_text_frame": True,
            "placeholder_type": "",
            "text": "before",
        },
    )
    shape.text = "after"
    assert dummy.called_shape_id == 42
    assert dummy.called_text == "after"
    assert shape.text == "after"
