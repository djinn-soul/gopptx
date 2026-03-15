"""Notes shape proxy for notes object model."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from typing_extensions import override

from .notes_text_model import NotesTextFrame

if TYPE_CHECKING:
    from collections.abc import Callable

    from .notes_slide import NotesSlide


_TEXT_PLACEHOLDER_TYPES = {"body", "title", "ctrTitle", "subTitle"}


class _NotesShapeStyleMixin:
    """Style mutation helpers for notes shape proxies."""

    @property
    def shape_id(self) -> int:
        raise NotImplementedError

    def _apply_shape_props(self, updates: dict[str, object]) -> None:
        raise NotImplementedError

    def set_fill_solid(self, color: str) -> None:
        """Set a solid fill color on the shape."""
        self._apply_shape_props({"fill": {"solid": color}})

    def set_line_color(self, color: str) -> None:
        """Set the shape line color."""
        self._apply_shape_props({"line": {"color": color}})

    def set_line_width(self, width_emu: int) -> None:
        """Set the shape line width in EMU units."""
        self._apply_shape_props({"line": {"width_emu": int(width_emu)}})


class NotesShape(_NotesShapeStyleMixin):
    """Proxy for one notes placeholder/shape entry."""

    _SHAPE_PROPS_WRITER_NAME = "_set_shape_props"
    _SHAPE_TEXT_WRITER_NAME = "_set_shape_text"

    def __init__(self, notes_slide: NotesSlide, payload: dict[str, object]) -> None:
        """Initialize a notes-shape proxy with source payload."""
        super().__init__()
        self._notes_slide = notes_slide
        self._payload = payload

    @property
    def placeholder_type(self) -> str:
        """Return placeholder kind (for example body/title) when present."""
        value = self._payload.get("placeholder_type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def idx(self) -> int:
        """Return placeholder index or -1 for non-placeholder shapes."""
        value = self._payload.get("placeholder_index", self._payload.get("Index"))
        if isinstance(value, int):
            return value
        return -1

    @property
    @override
    def shape_id(self) -> int:
        """Return unique shape identifier or -1 when unavailable."""
        value = self._payload.get("id", self._payload.get("ID"))
        if isinstance(value, int):
            return value
        return -1

    @property
    def name(self) -> str:
        """Return shape name."""
        value = self._payload.get("name", self._payload.get("Name"))
        return str(value) if isinstance(value, str) else ""

    @property
    def shape_type(self) -> str:
        """Return shape type string."""
        value = self._payload.get("type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def is_placeholder(self) -> bool:
        """Return whether this shape is a placeholder."""
        return self.idx >= 0

    @property
    def has_text_frame(self) -> bool:
        """Return whether this shape exposes text-frame behavior."""
        value = self._payload.get("has_text_frame", self._payload.get("HasTextFrame"))
        if isinstance(value, bool):
            return value
        return self.placeholder_type in _TEXT_PLACEHOLDER_TYPES

    @property
    def x(self) -> float:
        """Return x position in EMU units."""
        value = self._payload.get("x", self._payload.get("X"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def y(self) -> float:
        """Return y position in EMU units."""
        value = self._payload.get("y", self._payload.get("Y"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cx(self) -> float:
        """Return width in EMU units."""
        value = self._payload.get("cx", self._payload.get("CX"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cy(self) -> float:
        """Return height in EMU units."""
        value = self._payload.get("cy", self._payload.get("CY"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def left(self) -> float:
        """Return left position in EMU units."""
        return self.x

    @left.setter
    def left(self, value: float) -> None:
        """Set left position in EMU units."""
        self._apply_shape_props({"x": int(value)})
        self._payload["x"] = float(value)

    @property
    def top(self) -> float:
        """Return top position in EMU units."""
        return self.y

    @top.setter
    def top(self, value: float) -> None:
        """Set top position in EMU units."""
        self._apply_shape_props({"y": int(value)})
        self._payload["y"] = float(value)

    @property
    def width(self) -> float:
        """Return width in EMU units."""
        return self.cx

    @width.setter
    def width(self, value: float) -> None:
        """Set width in EMU units."""
        self._apply_shape_props({"w": int(value)})
        self._payload["cx"] = float(value)

    @property
    def height(self) -> float:
        """Return height in EMU units."""
        return self.cy

    @height.setter
    def height(self, value: float) -> None:
        """Set height in EMU units."""
        self._apply_shape_props({"h": int(value)})
        self._payload["cy"] = float(value)

    def _shape_props_writer(self) -> Callable[[int, dict[str, object]], None]:
        method_name = self._SHAPE_PROPS_WRITER_NAME
        return cast(
            "Callable[[int, dict[str, object]], None]",
            getattr(self._notes_slide, method_name),
        )

    def _shape_text_writer(self) -> Callable[[int, str], None]:
        method_name = self._SHAPE_TEXT_WRITER_NAME
        return cast(
            "Callable[[int, str], None]",
            getattr(self._notes_slide, method_name),
        )

    @override
    def _apply_shape_props(self, updates: dict[str, object]) -> None:
        self._shape_props_writer()(self.shape_id, updates)

    @property
    def text(self) -> str:
        """Return shape text."""
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            return self._notes_slide.text
        value = self._payload.get("text", self._payload.get("Text"))
        return str(value) if isinstance(value, str) else ""

    @text.setter
    def text(self, value: str) -> None:
        """Set shape text."""
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            self._notes_slide.text = value
            return
        if not self.has_text_frame:
            raise ValueError("target notes shape has no text frame")
        if self.shape_id < 0:
            raise ValueError("notes shape id is unavailable for mutation")
        self._shape_text_writer()(self.shape_id, value)
        self._payload["text"] = value

    @property
    def text_frame(self) -> NotesTextFrame | None:
        """Return text-frame proxy for text placeholders, else None."""
        if self.placeholder_type not in _TEXT_PLACEHOLDER_TYPES:
            return None
        return NotesTextFrame(self)

    @override
    def __repr__(self) -> str:
        """Return debug representation for this shape proxy."""
        return (
            f"<NotesShape id={self.shape_id} idx={self.idx}"
            f" type={self.shape_type!r} ph={self.placeholder_type!r}"
            f" name={self.name!r}>"
        )
