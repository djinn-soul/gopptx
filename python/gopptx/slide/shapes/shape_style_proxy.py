"""Read-only style views for a shape: resolved style and custom geometry."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import EffectiveShapeStyle, FreeformGeometry, Shape
    from ..contracts import SlidePresentationProtocol


class _StyleSlideProto(Protocol):
    """The slide members a style view needs to reach the bridge."""

    @property
    def presentation(self) -> SlidePresentationProtocol: ...

    @property
    def index(self) -> int: ...


class _StyledShapeProto(Protocol):
    """Structural protocol for the shape object a style view reads from."""

    @property
    def id(self) -> int: ...

    @property
    def slide(self) -> _StyleSlideProto: ...

    def shape_record(self) -> Shape: ...


class _ShapeStyleProxy:
    """Live, read-only view of how a shape is styled."""

    def __init__(self, shape: _StyledShapeProto) -> None:
        self._shape = shape

    @property
    def effective(self) -> EffectiveShapeStyle:
        """Return how the shape actually looks, after inheritance.

        A placeholder that sets no colour of its own reports ``None`` from
        ``fill.solid_color`` and from its runs, because the value lives on the
        layout, the master or the theme (upstream #1013). Every value here
        carries the ``source`` it was resolved from: ``shape``, ``layout``,
        ``master`` or ``theme``.
        """
        return self._shape.slide.presentation.get_effective_shape_style(
            self._shape.slide.index, self._shape.id
        )

    @property
    def freeform(self) -> FreeformGeometry | None:
        """Return the shape's custom geometry, or ``None`` for preset geometry.

        Freeform paths used to be write-only: ``<a:custGeom>`` was emitted and
        never parsed back (upstream #1020).
        """
        raw = cast("object", self._shape.shape_record().get("freeform"))
        if not isinstance(raw, dict):
            return None
        return cast("FreeformGeometry", raw)


ShapeStyleProxy = _ShapeStyleProxy
