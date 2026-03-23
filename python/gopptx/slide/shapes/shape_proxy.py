"""Shape object proxies with python-pptx-like ergonomics."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from ..tables.table import Table
from ..text.text_model import ShapeTextFrame
from .shape_format_proxies import (
    _ShapeFillProxy,
    _ShapeLineProxy,
    _ShapeShadowProxy,
)

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...constants import ShapeType
    from ...schemas import Shape, ShapeProps, ShapeUpdate
    from ...shapes import ShapeBuilder
    from ...text.run_builder import RunBuilder
    from ..contracts import SlidePresentationProtocol


class _ShapeProxySlideProto(Protocol):
    @property
    def presentation(self) -> SlidePresentationProtocol: ...

    @property
    def index(self) -> int: ...

    def list_shapes(self) -> list[Shape]: ...

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None: ...

    def add_shape(
        self,
        shape_type: ShapeType,
        bounds: tuple[float, float, float, float],
        **kwargs: str | ShapeProps,
    ) -> int: ...

    def shape(self, shape_id: int) -> ShapeProxy: ...


class ShapeProxy:
    """Live shape proxy object."""

    def __init__(self, slide: _ShapeProxySlideProto, shape_id: int) -> None:
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

    def set_text_runs(self, builders: list[RunBuilder]) -> None:
        """Replace the shape's text runs from a list of :class:`~gopptx.text.RunBuilder` instances.

        Convenience wrapper that delegates to :meth:`text_frame.set_runs`.

        Example::

            from gopptx.text import RunBuilder

            proxy.set_text_runs([
                RunBuilder("Hello").bold(),
                RunBuilder(" world").italic(),
            ])
        """
        self.text_frame.set_runs(builders)

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

    def __init__(self, slide: _ShapeProxySlideProto) -> None:
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

    def add(self, builder: ShapeBuilder) -> ShapeProxy:
        """Add a shape from a :class:`~gopptx.shapes.ShapeBuilder` and return its proxy.

        Example::

            proxy = slide.shapes.add(
                ShapeBuilder.rectangle(1.0, 1.0, 4.0, 2.0)
                .with_text("Hello")
                .with_fill("4472C4")
            )
        """
        shape_id = self._slide.add_shape(
            builder.shape_type,
            builder.bounds,
            **cast("dict[str, ShapeProps]", builder.to_kwargs()),
        )
        return self._slide.shape(shape_id)


ShapeFillProxy = _ShapeFillProxy
ShapeLineProxy = _ShapeLineProxy
ShapeShadowProxy = _ShapeShadowProxy
