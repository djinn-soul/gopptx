"""Notes shape proxy for notes object model."""
# ruff: noqa: D102, D107, SLF001, PLR0904
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

from .notes_text_model import NotesTextFrame

if TYPE_CHECKING:
    from .notes_slide import NotesSlide


_TEXT_PLACEHOLDER_TYPES = {"body", "title", "ctrTitle", "subTitle"}


class NotesShape:
    """Proxy for one notes placeholder/shape entry."""

    def __init__(self, notes_slide: NotesSlide, payload: dict[str, object]) -> None:
        super().__init__()
        self._notes_slide = notes_slide
        self._payload = payload

    @property
    def placeholder_type(self) -> str:
        value = self._payload.get("placeholder_type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def idx(self) -> int:
        value = self._payload.get("placeholder_index", self._payload.get("Index"))
        if isinstance(value, int):
            return value
        return -1

    @property
    def shape_id(self) -> int:
        value = self._payload.get("id", self._payload.get("ID"))
        if isinstance(value, int):
            return value
        return -1

    @property
    def name(self) -> str:
        value = self._payload.get("name", self._payload.get("Name"))
        return str(value) if isinstance(value, str) else ""

    @property
    def shape_type(self) -> str:
        value = self._payload.get("type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def is_placeholder(self) -> bool:
        return self.idx >= 0

    @property
    def has_text_frame(self) -> bool:
        value = self._payload.get("has_text_frame")
        if isinstance(value, bool):
            return value
        return self.placeholder_type in _TEXT_PLACEHOLDER_TYPES

    @property
    def x(self) -> float:
        value = self._payload.get("x", self._payload.get("X"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def y(self) -> float:
        value = self._payload.get("y", self._payload.get("Y"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cx(self) -> float:
        value = self._payload.get("cx", self._payload.get("CX"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cy(self) -> float:
        value = self._payload.get("cy", self._payload.get("CY"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def left(self) -> float:
        return self.x

    @left.setter
    def left(self, value: float) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"x": int(value)})
        self._payload["x"] = float(value)

    @property
    def top(self) -> float:
        return self.y

    @top.setter
    def top(self, value: float) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"y": int(value)})
        self._payload["y"] = float(value)

    @property
    def width(self) -> float:
        return self.cx

    @width.setter
    def width(self, value: float) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"w": int(value)})
        self._payload["cx"] = float(value)

    @property
    def height(self) -> float:
        return self.cy

    @height.setter
    def height(self, value: float) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"h": int(value)})
        self._payload["cy"] = float(value)

    def set_fill_solid(self, color: str) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"fill": {"solid": color}})

    def set_line_color(self, color: str) -> None:
        self._notes_slide._set_shape_props(self.shape_id, {"line": {"color": color}})

    def set_line_width(self, width_emu: int) -> None:
        self._notes_slide._set_shape_props(
            self.shape_id,
            {"line": {"width_emu": int(width_emu)}},
        )

    @property
    def text(self) -> str:
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            return self._notes_slide.text
        value = self._payload.get("text", self._payload.get("Text"))
        return str(value) if isinstance(value, str) else ""

    @text.setter
    def text(self, value: str) -> None:
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            self._notes_slide.text = value
            return
        if not self.has_text_frame:
            raise ValueError("target notes shape has no text frame")
        if self.shape_id < 0:
            raise ValueError("notes shape id is unavailable for mutation")
        self._notes_slide._set_shape_text(self.shape_id, value)
        self._payload["text"] = value

    @property
    def text_frame(self) -> NotesTextFrame | None:
        if self.placeholder_type not in _TEXT_PLACEHOLDER_TYPES:
            return None
        return NotesTextFrame(self)

    @override
    def __repr__(self) -> str:
        return (
            f"<NotesShape id={self.shape_id} idx={self.idx}"
            f" type={self.shape_type!r} ph={self.placeholder_type!r}"
            f" name={self.name!r}>"
        )
