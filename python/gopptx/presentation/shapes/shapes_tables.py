"""Shape-domain mixin for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...slide.shapes.freeform_builder import FreeformBuilder
from ..helpers import PresentationMixinBase
from .shape_media_mixin import PresentationShapeMediaMixin
from .shape_payload_mixin import PresentationShapePayloadMixin
from .shape_text_runs_mixin import PresentationShapeTextRunMixin

if TYPE_CHECKING:
    from ...schemas import (
        Shape,
        ShapeSearchQuery,
        ShapeSearchResult,
        ShapeUpdate,
    )


class PresentationShapeMixin(
    PresentationShapeMediaMixin,
    PresentationShapePayloadMixin,
    PresentationShapeTextRunMixin,
    PresentationMixinBase,
):
    """Mixin providing shape manipulation methods."""

    def search_shapes(self, query: ShapeSearchQuery | str) -> list[ShapeSearchResult]:
        """Search for shapes matching a query."""
        if isinstance(query, str):
            query = {"text_contains": query}
        result = self.execute(ops.OP_SEARCH_SHAPES, cast("dict[str, object]", query))
        return cast("list[ShapeSearchResult]", result.get("results", []))

    def list_shapes(self, slide_index: int) -> list[Shape]:
        """List all shapes on a slide."""
        result = self.execute(ops.OP_LIST_SHAPES, {"slide_index": slide_index})
        return cast("list[Shape]", result.get("shapes", []))

    def add_shape(
        self,
        slide_index: int,
        shape_type: str,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Add a shape to a slide."""
        x, y, w, h = bounds
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "type": shape_type,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        }
        self._apply_shape_payload_options(payload, kwargs, include_text=True)
        result = self.execute(ops.OP_ADD_SHAPE, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_textbox(
        self,
        slide_index: int,
        *bounds: float,
        text: str = "",
        **kwargs: object,
    ) -> int:
        """Add a textbox-like shape to a slide."""
        if len(bounds) != self._RECT_BOUNDS_COMPONENTS:
            msg = "add_textbox expects bounds as left, top, width, height"
            raise ValueError(msg)
        left, top, width, height = bounds
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "left": left,
            "top": top,
            "width": width,
            "height": height,
        }
        if text:
            payload["text"] = text
        self._apply_shape_payload_options(
            payload,
            kwargs,
            include_text=False,
        )
        result = self.execute(ops.OP_ADD_TEXTBOX, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_connector(
        self,
        slide_index: int,
        connector_type: str,
        *points: float,
        **kwargs: object,
    ) -> int:
        """Add a connector-like shape to a slide."""
        if len(points) != self._RECT_BOUNDS_COMPONENTS:
            msg = "add_connector expects points as begin_x, begin_y, end_x, end_y"
            raise ValueError(msg)
        begin_x, begin_y, end_x, end_y = points
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "connector_type": connector_type,
            "begin_x": begin_x,
            "begin_y": begin_y,
            "end_x": end_x,
            "end_y": end_y,
        }
        self._apply_shape_payload_options(
            payload,
            kwargs,
            include_text=True,
        )
        result = self.execute(ops.OP_ADD_CONNECTOR, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_group_shape(
        self,
        slide_index: int,
        shapes: list[int] | None = None,
    ) -> int:
        """Add a group shape to a slide."""
        payload: dict[str, object] = {"slide_index": slide_index}
        if shapes is not None:
            payload["shapes"] = shapes
        result = self.execute(ops.OP_ADD_GROUP_SHAPE, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def build_freeform(
        self,
        slide_index: int,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> FreeformBuilder:
        """Create a freeform builder for this slide."""
        return FreeformBuilder(
            self,
            slide_index,
            start_x=start_x,
            start_y=start_y,
            scale=scale,
        )

    def commit_freeform(
        self,
        slide_index: int,
        points: list[tuple[float, float]],
        *,
        close: bool,
        options: dict[str, object] | None = None,
    ) -> int:
        """Create a freeform shape from prepared points."""
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "points": [[x, y] for x, y in points],
            "close": close,
        }
        opt = options or {}
        self._apply_shape_payload_options(payload, opt, include_text=True)
        result = self.execute(ops.OP_BUILD_FREEFORM, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        """Remove a shape from a slide."""
        self.execute(
            ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id}
        )

    def group_shapes(self, slide_index: int, shape_ids: list[int]) -> int:
        """Group multiple shapes on a slide into a group shape.

        Returns the ID of the created group shape.
        """
        result = self.execute(
            ops.OP_GROUP_SHAPES,
            {"slide_index": slide_index, "shape_ids": shape_ids},
        )
        return int(cast("int", result.get("group_id", -1)))

    def ungroup_shapes(self, slide_index: int, shape_id: int) -> int:
        """Ungroup a group shape, returning the ID of the first member shape."""
        result = self.execute(
            ops.OP_UNGROUP_SHAPES,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return int(cast("int", result.get("group_id", -1)))

    def move_shape_to_front(self, slide_index: int, shape_id: int) -> None:
        """Move a shape to the front of the z-order."""
        self.execute(
            ops.OP_MOVE_SHAPE_TO_FRONT,
            {"slide_index": slide_index, "shape_id": shape_id},
        )

    def move_shape_to_back(self, slide_index: int, shape_id: int) -> None:
        """Move a shape to the back of the z-order."""
        self.execute(
            ops.OP_MOVE_SHAPE_TO_BACK,
            {"slide_index": slide_index, "shape_id": shape_id},
        )

    def move_shape_to_index(
        self, slide_index: int, shape_id: int, target_index: int
    ) -> None:
        """Move a shape to a specific z-index within a slide.

        Args:
            slide_index: Zero-based slide index.
            shape_id: ID of the shape to reorder.
            target_index: Zero-based target z-order position.
        """
        self.execute(
            ops.OP_MOVE_SHAPE_TO_INDEX,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "target_index": target_index,
            },
        )

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Update shape properties."""
        updates_dict = dict(cast("dict[str, object]", updates))
        normalized_updates: dict[str, object] = {}
        self._apply_shape_payload_options(
            normalized_updates, updates_dict, include_text=True
        )
        # Copy remaining update fields not handled by the helper
        for k, v in updates_dict.items():
            if k not in normalized_updates:
                normalized_updates[k] = v

        self.execute(
            ops.OP_UPDATE_SHAPE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "updates": normalized_updates,
            },
        )


__all__ = [
    "PresentationShapeMixin",
]
