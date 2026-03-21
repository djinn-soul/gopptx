"""Regression tests for slide background payload coercion."""

from __future__ import annotations

from gopptx import ops
from gopptx.slide.shapes.smartart_anim_mixin import SlideSmartArtAnimMixin


class _DummyPresentation:
    def __init__(self) -> None:
        self.calls: list[tuple[str, dict[str, object]]] = []

    def execute(self, op: str, payload: dict[str, object]) -> dict[str, object]:
        self.calls.append((op, payload))
        return {}


class _DummySlide(SlideSmartArtAnimMixin):
    def __init__(self) -> None:
        self._presentation = _DummyPresentation()
        self._index = 2

    @property
    def index(self) -> int:
        return self._index

    def _invalidate_shape_cache_if_present(self) -> None:
        return None


def test_set_background_none_optionals_are_not_stringified() -> None:
    """None optionals should be serialized as empty strings, not 'None'."""
    slide = _DummySlide()
    slide.set_background(
        "image",
        color=None,
        image_path=None,
        image_data=None,
        color_ref=None,
    )

    calls = slide._presentation.calls  # noqa: SLF001
    assert len(calls) == 1
    op, payload = calls[0]
    assert op == ops.OP_SET_SLIDE_BACKGROUND
    assert not payload["color"]
    assert not payload["image_path"]
    assert not payload["image_data"]
    assert not payload["color_ref"]


def test_set_background_falsey_non_none_values_are_preserved() -> None:
    """Falsey but non-None values should still be stringified."""
    slide = _DummySlide()
    slide.set_background("theme", color_ref=0)

    _, payload = slide._presentation.calls[0]  # noqa: SLF001
    assert payload["color_ref"] == "0"
