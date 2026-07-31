"""Extra shape features (deletion, z-order, visibility, flip aliases)."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import Shape, ShapeUpdate


class _ShapeExtraSlideProto(Protocol):
    @property
    def index(self) -> int: ...

    def list_shapes(self) -> list[Shape]: ...

    def remove_shape(self, shape_id: int) -> None: ...

    def move_shape_to_front(self, shape_id: int) -> None: ...

    def move_shape_to_back(self, shape_id: int) -> None: ...


class _ShapeExtraFeatureHost(Protocol):
    @property
    def slide(self) -> _ShapeExtraSlideProto: ...

    @property
    def id(self) -> int: ...

    def shape_record(self) -> Shape: ...

    def apply_update(self, patch: ShapeUpdate) -> None: ...


class ShapeProxyExtraMixin:
    """Extra shape proxy features kept separate to maintain LOC limits."""

    def delete(self: _ShapeExtraFeatureHost) -> None:
        """Delete this shape from its slide (Issue #41)."""
        self.slide.remove_shape(self.id)

    def move_to_front(self: _ShapeExtraFeatureHost) -> None:
        """Move shape to the front of the slide z-order (Issue #49)."""
        self.slide.move_shape_to_front(self.id)

    def move_to_back(self: _ShapeExtraFeatureHost) -> None:
        """Move shape to the back of the slide z-order (Issue #49)."""
        self.slide.move_shape_to_back(self.id)

    @property
    def z_order(self: _ShapeExtraFeatureHost) -> int:
        """Return zero-based z-order index on the slide (Issue #49)."""
        shapes = self.slide.list_shapes()
        for idx, shape in enumerate(shapes):
            payload = cast("dict[str, object]", shape)
            sid = payload.get("ID")
            if sid is None:
                sid = payload.get("id")
            if sid is not None and int(str(sid)) == self.id:
                return idx
        return -1

    @property
    def hidden(self: _ShapeExtraFeatureHost) -> bool:
        """Return whether shape is hidden (Issue #452)."""
        shape = self.shape_record()
        return bool(shape.get("Hidden", shape.get("hidden", False)))

    @hidden.setter
    def hidden(self: _ShapeExtraFeatureHost, value: bool) -> None:
        """Set shape hidden state (Issue #452)."""
        self.apply_update(cast("ShapeUpdate", {"hidden": bool(value)}))
