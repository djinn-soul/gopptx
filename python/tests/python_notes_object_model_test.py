"""Notes object-model parity tests."""

from __future__ import annotations

from typing import Protocol

from gopptx import Presentation
from gopptx.slide.notes.notes_slide import NotesShape, NotesSlide


class _HasNotes(Protocol):
    notes: str


def _assert_first_shape_contract(notes_slide: NotesSlide, first: NotesShape) -> None:
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
    assert notes_slide.shape_by_name(first.name) is not None
    assert isinstance(notes_slide.text_shapes, list)
    assert isinstance(notes_slide.placeholder_shapes, list)
    assert all(shape.has_text_frame for shape in notes_slide.text_shapes)
    body_matches = notes_slide.shapes.find_all(placeholder_type="body")
    if body_matches:
        assert body_matches[0].placeholder_type == "body"


def _exercise_body_proxy(
    slide: _HasNotes, notes_slide: NotesSlide, body: NotesShape
) -> None:
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


def test_notes_shape_collection_and_body_proxy() -> None:
    """Notes slide exposes typed shape proxies and body text mutation."""
    with Presentation.new("Notes Object Model") as prs:
        slide = prs.add_slide("S1")
        slide.notes = "speaker notes"
        notes_slide = slide.notes_slide
        assert notes_slide is not None

        assert len(notes_slide.shapes) > 0
        first = notes_slide.shapes[0]
        _assert_first_shape_contract(notes_slide, first)

        body = notes_slide.body_shape
        if body is not None:
            _exercise_body_proxy(slide, notes_slide, body)


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


def test_notes_shape_geometry_and_style_route_to_props_setter() -> None:
    class _DummyNotesSlide:
        def __init__(self) -> None:
            self.calls: list[tuple[int, dict[str, object]]] = []
            self.text = ""

        def _set_shape_props(self, shape_id: int, updates: dict[str, object]) -> None:
            self.calls.append((shape_id, dict(updates)))

    dummy = _DummyNotesSlide()
    shape = NotesShape(
        dummy,  # type: ignore[arg-type]
        {
            "id": 7,
            "name": "Aux",
            "type": "sp",
            "has_text_frame": True,
            "placeholder_type": "",
            "x": 10.0,
            "y": 20.0,
            "cx": 30.0,
            "cy": 40.0,
        },
    )
    shape.left = 111
    shape.top = 222
    shape.width = 333
    shape.height = 444
    shape.set_fill_solid("FF0000")
    shape.set_line_color("00FF00")
    shape.set_line_width(12700)

    assert len(dummy.calls) == 7
    assert dummy.calls[0] == (7, {"x": 111})
    assert dummy.calls[1] == (7, {"y": 222})
    assert dummy.calls[2] == (7, {"w": 333})
    assert dummy.calls[3] == (7, {"h": 444})
    assert dummy.calls[4] == (7, {"fill": {"solid": "FF0000"}})
    assert dummy.calls[5] == (7, {"line": {"color": "00FF00"}})
    assert dummy.calls[6] == (7, {"line": {"width_emu": 12700}})


def test_notes_slide_shape_bounds_and_fill_background_helpers() -> None:
    class _DummyPresentation:
        def __init__(self) -> None:
            self.calls: list[tuple[int, int, dict[str, object]]] = []

        def get_notes_payload(self, _slide_index: int) -> dict[str, object]:
            return {
                "notes_shapes": [
                    {"id": 5, "name": "Body 1", "placeholder_type": "body"},
                ]
            }

        def set_notes_shape_props(
            self, slide_index: int, shape_id: int, updates: dict[str, object]
        ) -> None:
            self.calls.append((slide_index, shape_id, updates))

    class _DummySlide:
        def __init__(self, presentation: _DummyPresentation) -> None:
            self.index = 2
            self._presentation = presentation
            self.notes = ""

    dummy_presentation = _DummyPresentation()
    dummy_slide = _DummySlide(dummy_presentation)
    notes_slide = NotesSlide(dummy_slide)  # type: ignore[arg-type]
    assert notes_slide.shape_by_name("Body 1") is not None
    notes_slide.set_shape_bounds(5, left=10, top=20, width=30, height=40)
    notes_slide.set_shape_fill_background(5)
    assert dummy_presentation.calls[0] == (2, 5, {"x": 10, "y": 20, "w": 30, "h": 40})
    assert dummy_presentation.calls[1] == (2, 5, {"fill": {"background": True}})


def test_notes_slide_shape_style_helpers_route_updates() -> None:
    class _DummyPresentation:
        def __init__(self) -> None:
            self.calls: list[tuple[int, int, dict[str, object]]] = []

        def get_notes_payload(self, _slide_index: int) -> dict[str, object]:
            return {"notes_shapes": [{"id": 9, "name": "Aux"}]}

        def set_notes_shape_props(
            self, slide_index: int, shape_id: int, updates: dict[str, object]
        ) -> None:
            self.calls.append((slide_index, shape_id, updates))

    class _DummySlide:
        def __init__(self, presentation: _DummyPresentation) -> None:
            self.index = 4
            self._presentation = presentation
            self.notes = ""

    dummy_presentation = _DummyPresentation()
    notes_slide = NotesSlide(_DummySlide(dummy_presentation))  # type: ignore[arg-type]
    notes_slide.set_shape_fill_gradient(
        9,
        angle_deg=45.0,
        stops=[
            {"position_pct": 0.0, "color": "FF0000"},
            {"position_pct": 1.0, "color": "00FF00"},
        ],
    )
    notes_slide.set_shape_fill_pattern(
        9,
        preset="pct10",
        fg_color="FFFFFF",
        bg_color="000000",
    )
    notes_slide.set_shape_line_dash(
        9,
        dash_style="sysDot",
        color="FF00FF",
        width_emu=12700,
    )

    calls = dummy_presentation.calls
    assert calls[0][0:2] == (4, 9)
    assert calls[0][2] == {
        "fill": {
            "gradient": {
                "angle_deg": 45.0,
                "stops": [
                    {"position_pct": 0.0, "color": "FF0000"},
                    {"position_pct": 1.0, "color": "00FF00"},
                ],
            }
        }
    }
    assert calls[1][2] == {
        "fill": {
            "pattern": {
                "preset": "pct10",
                "fg_color": "FFFFFF",
                "bg_color": "000000",
            }
        }
    }
    assert calls[2][2] == {
        "line": {"dash_style": "sysDot", "color": "FF00FF", "width_emu": 12700}
    }
