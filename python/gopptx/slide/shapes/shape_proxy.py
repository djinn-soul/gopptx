"""Shape object proxies with python-pptx-like ergonomics."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ..tables.table import Table
from ..text.text_model import ShapeTextFrame

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import FillFormat, LineFormat, ShadowFormat, Shape, ShapeUpdate
    from ..slide import Slide


class _ShapeFillProxy:
    """Live fill proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> FillFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("fill", record.get("Fill", {})))
        return cast("FillFormat", raw if raw is not None else {})

    def _apply(self, payload: FillFormat) -> None:
        self._shape.apply_update(cast("ShapeUpdate", {"fill": payload}))

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


class _ShapeLineProxy:
    """Live line proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> LineFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("line", record.get("Line", {})))
        return cast("LineFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape.apply_update(cast("ShapeUpdate", {"line": payload}))

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


class _ShapeShadowProxy:
    """Live shadow proxy."""

    def __init__(self, shape: ShapeProxy) -> None:
        self._shape = shape

    def _payload(self) -> ShadowFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("shadow", record.get("Shadow", {})))
        return cast("ShadowFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape.apply_update(cast("ShapeUpdate", {"shadow": payload}))

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
        """Create a proxy around the specified slide shape."""
        self._slide = slide
        self._shape_id = shape_id
        self._fill = _ShapeFillProxy(self)
        self._line = _ShapeLineProxy(self)
        self._shadow = _ShapeShadowProxy(self)
        self._text_frame: ShapeTextFrame | None = None

    def shape_record(self) -> Shape:
        """Return the current shape payload from slide state."""
        for shape in self._slide.list_shapes():
            shape_id = shape.get("ID", shape.get("id"))
            if shape_id is None:
                continue
            if int(str(shape_id)) == self._shape_id:
                return shape
        raise KeyError(f"shape id not found: {self._shape_id}")

    def apply_update(self, patch: ShapeUpdate) -> None:
        """Apply a shape patch to the backing slide."""
        self._slide.update_shape(self._shape_id, patch)

    @property
    def id(self) -> int:
        """Return the shape id."""
        return self._shape_id

    @property
    def name(self) -> str:
        """Return the shape name."""
        shape = self.shape_record()
        value = shape.get("Name", shape.get("name", ""))
        return str(value)

    @property
    def shape_type(self) -> str:
        """Return the shape type token."""
        shape = self.shape_record()
        value = shape.get("Type", shape.get("type", ""))
        return str(value)

    @property
    def text(self) -> str:
        """Return shape text content."""
        shape = self.shape_record()
        value = shape.get("Text", shape.get("text", ""))
        return str(value)

    @text.setter
    def text(self, value: str) -> None:
        self.apply_update(cast("ShapeUpdate", {"text": value}))

    @property
    def text_frame(self) -> ShapeTextFrame:
        """Return the lazy text-frame proxy."""
        if self._text_frame is None:
            self._text_frame = ShapeTextFrame(self._slide, self._shape_id)
        return self._text_frame

    @property
    def fill(self) -> _ShapeFillProxy:
        """Return fill formatting facade."""
        return self._fill

    @property
    def line(self) -> _ShapeLineProxy:
        """Return line formatting facade."""
        return self._line

    @property
    def shadow(self) -> _ShapeShadowProxy:
        """Return shadow formatting facade."""
        return self._shadow

    @property
    def left(self) -> int:
        """Return shape left position in EMU."""
        shape = self.shape_record()
        value = shape.get("X", shape.get("x", 0))
        return int(str(value))

    @left.setter
    def left(self, value: int) -> None:
        self.apply_update(cast("ShapeUpdate", {"x": value}))

    @property
    def top(self) -> int:
        """Return shape top position in EMU."""
        shape = self.shape_record()
        value = shape.get("Y", shape.get("y", 0))
        return int(str(value))

    @top.setter
    def top(self, value: int) -> None:
        self.apply_update(cast("ShapeUpdate", {"y": value}))

    @property
    def width(self) -> int:
        """Return shape width in EMU."""
        shape = self.shape_record()
        value = shape.get("W", shape.get("w", 0))
        return int(str(value))

    @width.setter
    def width(self, value: int) -> None:
        self.apply_update(cast("ShapeUpdate", {"w": value}))

    @property
    def height(self) -> int:
        """Return shape height in EMU."""
        shape = self.shape_record()
        value = shape.get("H", shape.get("h", 0))
        return int(str(value))

    @height.setter
    def height(self, value: int) -> None:
        self.apply_update(cast("ShapeUpdate", {"h": value}))

    @property
    def has_table(self) -> bool:
        """True if this shape is a table."""
        return self.shape_type in {"tbl", "graphicFrame"}

    @property
    def table(self) -> Table:
        """Return a table proxy if this shape is a table."""
        if not self.has_table:
            raise AttributeError(f"shape {self._shape_id} has no table")
        return Table(self._slide.presentation, self._slide.index, self._shape_id)


class ShapeCollection:
    """python-pptx-style slide shapes collection."""

    def __init__(self, slide: Slide) -> None:
        """Create a shape collection for a slide."""
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
        """Return shape count."""
        return len(self._shape_ids())

    def __getitem__(self, index: int) -> ShapeProxy:
        """Return shape proxy at index."""
        ids = self._shape_ids()
        if index < 0:
            index += len(ids)
        if index < 0 or index >= len(ids):
            raise IndexError("shape index out of range")
        return self._slide.shape(ids[index])

    def __iter__(self) -> Iterator[ShapeProxy]:
        """Iterate shape proxies."""
        for shape_id in self._shape_ids():
            yield self._slide.shape(shape_id)


ShapeFillProxy = _ShapeFillProxy
ShapeLineProxy = _ShapeLineProxy
ShapeShadowProxy = _ShapeShadowProxy
