"""Shape object proxies with python-pptx-like ergonomics."""
# ruff: noqa: D102,D105,D107,SLF001,TC003
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from collections.abc import Iterator
from typing import TYPE_CHECKING, cast

from .text_model import ShapeTextFrame

if TYPE_CHECKING:
    from ..schemas import FillFormat, LineFormat, ShadowFormat, Shape, ShapeUpdate
    from .slide import Slide


class ShapeFillProxy:
    """Live fill proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> FillFormat:
        record = self._shape._shape_record()
        raw = cast("object", record.get("fill", record.get("Fill", {})))
        return cast("FillFormat", raw if raw is not None else {})

    def _apply(self, payload: FillFormat) -> None:
        self._shape._slide.update_shape(
            self._shape.id,
            cast("ShapeUpdate", {"fill": payload}),
        )

    @property
    def solid_color(self) -> str | None:
        payload = self._payload()
        value = payload.get("solid")
        return str(value) if isinstance(value, str) else None

    @solid_color.setter
    def solid_color(self, value: str | None) -> None:
        if value is None:
            self._apply(cast("FillFormat", {"background": True}))
            return
        self._apply(cast("FillFormat", {"solid": value}))

    def background(self) -> None:
        self._apply(cast("FillFormat", {"background": True}))


class ShapeLineProxy:
    """Live line proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> LineFormat:
        record = self._shape._shape_record()
        raw = cast("object", record.get("line", record.get("Line", {})))
        return cast("LineFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape._slide.update_shape(
            self._shape.id,
            cast("ShapeUpdate", {"line": payload}),
        )

    @property
    def color(self) -> str | None:
        value = self._payload().get("color")
        return str(value) if isinstance(value, str) else None

    @color.setter
    def color(self, value: str) -> None:
        self._apply({"color": value})

    @property
    def width(self) -> int | None:
        value = self._payload().get("width_emu")
        return int(value) if isinstance(value, int) else None

    @width.setter
    def width(self, value: int) -> None:
        self._apply({"width_emu": value})

    @property
    def dash_style(self) -> str | None:
        value = self._payload().get("dash_style")
        return str(value) if isinstance(value, str) else None

    @dash_style.setter
    def dash_style(self, value: str) -> None:
        self._apply({"dash_style": value})


class ShapeShadowProxy:
    """Live shadow proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> ShadowFormat:
        record = self._shape._shape_record()
        raw = cast("object", record.get("shadow", record.get("Shadow", {})))
        return cast("ShadowFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape._slide.update_shape(
            self._shape.id,
            cast("ShapeUpdate", {"shadow": payload}),
        )

    @property
    def color(self) -> str | None:
        value = self._payload().get("color")
        return str(value) if isinstance(value, str) else None

    @color.setter
    def color(self, value: str) -> None:
        self._apply({"color": value})

    @property
    def blur_radius(self) -> int | None:
        value = self._payload().get("blur_emu")
        return int(value) if isinstance(value, int) else None

    @blur_radius.setter
    def blur_radius(self, value: int) -> None:
        self._apply({"blur_emu": value})

    @property
    def distance(self) -> int | None:
        value = self._payload().get("distance_emu")
        return int(value) if isinstance(value, int) else None

    @distance.setter
    def distance(self, value: int) -> None:
        self._apply({"distance_emu": value})

    @property
    def angle(self) -> float | None:
        value = self._payload().get("angle_deg")
        return float(value) if isinstance(value, int | float) else None

    @angle.setter
    def angle(self, value: float) -> None:
        self._apply({"angle_deg": value})

    @property
    def inherit(self) -> bool | None:
        value = self._payload().get("inherit")
        return bool(value) if isinstance(value, bool) else None

    @inherit.setter
    def inherit(self, value: bool) -> None:
        self._apply({"inherit": value})


class ShapeProxy:
    """Live shape proxy object."""

    def __init__(self, slide: Slide, shape_id: int) -> None:
        self._slide = slide
        self._presentation = slide._presentation
        self._shape_id = shape_id
        self._fill = ShapeFillProxy(self)
        self._line = ShapeLineProxy(self)
        self._shadow = ShapeShadowProxy(self)
        self._text_frame: ShapeTextFrame | None = None

    def _shape_record(self) -> Shape:
        return self._slide._shape_record_by_id(self._shape_id)

    @property
    def id(self) -> int:
        return self._shape_id

    @property
    def name(self) -> str:
        shape = self._shape_record()
        value = shape.get("Name", shape.get("name", ""))
        return str(value)

    @property
    def shape_type(self) -> str:
        shape = self._shape_record()
        value = shape.get("Type", shape.get("type", ""))
        return str(value)

    @property
    def text(self) -> str:
        shape = self._shape_record()
        value = shape.get("Text", shape.get("text", ""))
        return str(value)

    @text.setter
    def text(self, value: str) -> None:
        self._slide.update_shape(self._shape_id, cast("ShapeUpdate", {"text": value}))

    @property
    def text_frame(self) -> ShapeTextFrame:
        if self._text_frame is None:
            self._text_frame = ShapeTextFrame(self._slide, self._shape_id)
        return self._text_frame

    @property
    def fill(self) -> ShapeFillProxy:
        return self._fill

    @property
    def line(self) -> ShapeLineProxy:
        return self._line

    @property
    def shadow(self) -> ShapeShadowProxy:
        return self._shadow

    @property
    def left(self) -> int:
        shape = self._shape_record()
        value = shape.get("X", shape.get("x", 0))
        return int(str(value))

    @left.setter
    def left(self, value: int) -> None:
        self._slide.update_shape(self._shape_id, cast("ShapeUpdate", {"x": value}))

    @property
    def top(self) -> int:
        shape = self._shape_record()
        value = shape.get("Y", shape.get("y", 0))
        return int(str(value))

    @top.setter
    def top(self, value: int) -> None:
        self._slide.update_shape(self._shape_id, cast("ShapeUpdate", {"y": value}))

    @property
    def width(self) -> int:
        shape = self._shape_record()
        value = shape.get("W", shape.get("w", 0))
        return int(str(value))

    @width.setter
    def width(self, value: int) -> None:
        self._slide.update_shape(self._shape_id, cast("ShapeUpdate", {"w": value}))

    @property
    def height(self) -> int:
        shape = self._shape_record()
        value = shape.get("H", shape.get("h", 0))
        return int(str(value))

    @height.setter
    def height(self, value: int) -> None:
        self._slide.update_shape(self._shape_id, cast("ShapeUpdate", {"h": value}))


class ShapeCollection:
    """python-pptx-style slide shapes collection."""

    def __init__(self, slide: Slide) -> None:
        self._slide = slide

    def _shape_ids(self) -> list[int]:
        out: list[int] = []
        for shape in self._slide.list_shapes():
            shape_id = shape.get("ID", shape.get("id"))
            if shape_id is None:
                continue
            out.append(int(str(shape_id)))
        return out

    def __len__(self) -> int:
        return len(self._shape_ids())

    def __getitem__(self, index: int) -> ShapeProxy:
        ids = self._shape_ids()
        if index < 0:
            index += len(ids)
        if index < 0 or index >= len(ids):
            raise IndexError("shape index out of range")
        return self._slide.shape(ids[index])

    def __iter__(self) -> Iterator[ShapeProxy]:
        for shape_id in self._shape_ids():
            yield self._slide.shape(shape_id)
