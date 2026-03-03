"""Shape, text, notes, and table mixins for the Presentation API."""

from __future__ import annotations

import base64
import os
from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from ..slide.freeform_builder import FreeformBuilder
from ..utils import normalize_table_index
from .helpers import PresentationProtocol

if TYPE_CHECKING:
    from ..schemas import (
        ImageCrop,
        ImageMetadata,
        Shape,
        ShapeProps,
        ShapeSearchQuery,
        ShapeSearchResult,
        ShapeUpdate,
        TableCellInfo,
        TableInfo,
    )
else:

    class PresentationProtocol:
        """Runtime placeholder to avoid Protocol abstract behavior."""


class PresentationNotesMixin(PresentationProtocol):
    """Mixin providing speaker notes methods."""

    def _has_notes_slide(self, slide_index: int) -> bool:
        """Return whether a notes slide currently exists for a slide."""
        result = self.execute(ops.OP_HAS_NOTES_SLIDE, {"slide_index": slide_index})
        return bool(result.get("has_notes_slide", False))

    def get_notes(self, slide_index: int) -> str:
        """Get speaker notes for a slide."""
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(cast("str", result.get("text", "")))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Set speaker notes for a slide."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})


class PresentationShapeMixin(PresentationProtocol):
    """Mixin providing shape manipulation methods."""

    @staticmethod
    def _init_bounds_payload(
        slide_index: int, bounds: tuple[float, float, float, float]
    ) -> dict[str, object]:
        x, y, w, h = bounds
        return {"slide_index": slide_index, "x": x, "y": y, "w": w, "h": h}

    @staticmethod
    def _set_source_payload(
        payload: dict[str, object],
        source: str | bytes | os.PathLike[str] | None,
        *,
        path_key: str = "path",
        data_key: str = "data",
    ) -> None:
        if source is None:
            return
        if isinstance(source, str):
            payload[path_key] = source
            return
        if isinstance(source, os.PathLike):
            payload[path_key] = os.fspath(source)
            return
        payload[data_key] = base64.b64encode(source).decode("ascii")

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
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a shape to a slide."""
        x, y, w, h = bounds
        text = kwargs.get("text")
        runs = kwargs.get("runs")
        text_frame = kwargs.get("text_frame")
        click_action = kwargs.get("click_action")
        properties = kwargs.get("properties")
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "type": shape_type,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        }
        if text is not None:
            payload["text"] = text
        if runs is not None:
            payload["runs"] = runs
        if text_frame is not None:
            payload["text_frame"] = cast("dict[str, object]", text_frame)
        if click_action is not None:
            payload["click_action"] = cast("dict[str, object]", click_action)
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_ADD_SHAPE, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_textbox(
        self,
        slide_index: int,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a textbox-like shape to a slide."""
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "left": left,
            "top": top,
            "width": width,
            "height": height,
        }
        if text:
            payload["text"] = text
        for key in (
            "runs",
            "text_frame",
            "click_action",
            "hover_action",
            "properties",
        ):
            if key in kwargs and kwargs[key] is not None:
                payload[key] = cast("object", kwargs[key])
        result = self.execute(ops.OP_ADD_TEXTBOX, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_connector(
        self,
        slide_index: int,
        connector_type: str,
        begin_x: float,
        begin_y: float,
        end_x: float,
        end_y: float,
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a connector-like shape to a slide."""
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "connector_type": connector_type,
            "begin_x": begin_x,
            "begin_y": begin_y,
            "end_x": end_x,
            "end_y": end_y,
        }
        for key in (
            "text",
            "runs",
            "text_frame",
            "click_action",
            "hover_action",
            "properties",
        ):
            if key in kwargs and kwargs[key] is not None:
                payload[key] = cast("object", kwargs[key])
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

    def _commit_freeform(
        self,
        slide_index: int,
        points: list[tuple[float, float]],
        *,
        close: bool,
        text: str | None = None,
        runs: object | None = None,
        text_frame: object | None = None,
        click_action: object | None = None,
        hover_action: object | None = None,
        properties: object | None = None,
    ) -> int:
        """Create a freeform shape from prepared points."""
        payload: dict[str, object] = {
            "slide_index": slide_index,
            "points": [[x, y] for x, y in points],
            "close": close,
        }
        if text is not None:
            payload["text"] = text
        if runs is not None:
            payload["runs"] = runs
        if text_frame is not None:
            payload["text_frame"] = text_frame
        if click_action is not None:
            payload["click_action"] = click_action
        if hover_action is not None:
            payload["hover_action"] = hover_action
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_BUILD_FREEFORM, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_image(
        self,
        slide_index: int,
        source: str | bytes | None = None,
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
        *,
        path: str | None = None,
        data: bytes | None = None,
        format: str | None = None,
        img_format: str | None = None,
        crop: ImageCrop | None = None,
        rotation: float | None = None,
        flip_h: bool | None = None,
        flip_v: bool | None = None,
    ) -> int:
        """Add an image to a slide.

        Args:
            slide_index: Slide index (0-based).
            path: Path to the image file, or None if data is provided.
            bounds: (x, y, w, h) in EMU.
            data: Raw image bytes.
            format: Image format (e.g., 'png', 'jpeg') required if data is bytes.
            crop: ImageCrop dictionary.
            rotation: Rotation in degrees.
            flip_h: Flip horizontally.
            flip_v: Flip vertically.
        """
        payload = self._init_bounds_payload(slide_index, bounds)
        if source:
            self._set_source_payload(payload, source)
        elif path:
            self._set_source_payload(payload, path)
        elif data:
            self._set_source_payload(payload, data)
            resolved_format = format or img_format
            if resolved_format:
                payload["format"] = resolved_format

        options: dict[str, object] = {}
        if crop:
            options["crop"] = cast("dict[str, object]", crop)
        if rotation is not None:
            options["rotation"] = rotation
        if flip_h is not None:
            options["flip_h"] = flip_h
        if flip_v is not None:
            options["flip_v"] = flip_v

        if options:
            payload["options"] = options

        result = self.execute(ops.OP_ADD_IMAGE, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def get_image_metadata(self, slide_index: int, shape_id: int) -> ImageMetadata:
        """Get dimensions and format for an image shape."""
        result = self.execute(
            ops.OP_GET_IMAGE_METADATA,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("ImageMetadata", result)

    def add_video(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        """Add a video to a slide."""
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if name:
            payload["name"] = name
        if mime_type:
            payload["mime_type"] = mime_type

        if poster_frame:
            self._set_source_payload(
                payload,
                poster_frame,
                path_key="poster_path",
                data_key="poster_data",
            )

        result = self.execute(ops.OP_ADD_VIDEO, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_ole_object(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        prog_id: str | None = None,
        icon: str | bytes | None = None,
    ) -> int:
        """Add an OLE object to a slide."""
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if name:
            payload["name"] = name
        if prog_id:
            payload["prog_id"] = prog_id

        if icon:
            self._set_source_payload(
                payload,
                icon,
                path_key="icon_path",
                data_key="icon_data",
            )

        result = self.execute(ops.OP_ADD_OLE_OBJECT, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        """Remove a shape from a slide."""
        self.execute(
            ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id}
        )

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

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Update shape properties."""
        self.execute(
            ops.OP_UPDATE_SHAPE,
            {"slide_index": slide_index, "shape_id": shape_id, "updates": updates},
        )


class PresentationTextMixin(PresentationProtocol):
    """Mixin providing text search and replace methods."""

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        """Find and replace text in the presentation."""
        result = self.execute(
            ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text}
        )
        return int(cast("int", result.get("replacements", 0)))


class PresentationTableMixin(PresentationProtocol):
    """Mixin providing table creation and manipulation methods."""

    def add_table(
        self,
        slide_index: int,
        rows: int,
        cols: int,
        bounds: tuple[int, int, int, int],
    ) -> int:
        """Add a table to a slide."""
        x, y, cx, cy = bounds
        result = self.execute(
            ops.OP_ADD_TABLE,
            {
                "slide_index": slide_index,
                "rows": rows,
                "cols": cols,
                "x": x,
                "y": y,
                "cx": cx,
                "cy": cy,
            },
        )
        return int(cast("int", result.get("shape_id", 0)))

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Get table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE, {"slide_index": slide_index, "shape_id": shape_id}
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))

    def set_table_flags(
        self, slide_index: int, shape_id: int, flags: dict[str, bool]
    ) -> None:
        """Set table style flags."""
        self.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {"slide_index": slide_index, "shape_id": shape_id, "flags": flags},
        )

    def set_table_cell_text(
        self, slide_index: int, shape_id: int, row: int, col: int, text: str
    ) -> None:
        """Set the text of a table cell."""
        self.execute(
            ops.OP_UPDATE_TABLE_CELL,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row": row,
                "col": col,
                "updates": {"text": text},
            },
        )

    def get_table_cell(
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> TableCellInfo:
        """Get information about a table cell."""
        table = self.get_table(slide_index, shape_id)
        cells = table.get("cells", [])
        cell_map: dict[tuple[int, int], dict[str, object]] = {}
        for cell in cells:
            try:
                row_idx = normalize_table_index(cell["row"])
                col_idx = normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            cell_map[row_idx, col_idx] = cast("dict[str, object]", cell)
        found = cell_map.get((row, col))
        if found is not None:
            return cast("TableCellInfo", found)
        raise GopptxError(f"table cell [{row},{col}] not found", code="OP_FAILED")

    def merge_table_cells(
        self,
        slide_index: int,
        shape_id: int,
        cell_range: tuple[int, int, int, int],
    ) -> None:
        """Merge a range of table cells."""
        row1, col1, row2, col2 = cell_range
        self.execute(
            ops.OP_MERGE_TABLE_CELLS,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row1": row1,
                "col1": col1,
                "row2": row2,
                "col2": col2,
            },
        )

    def split_table_cell(
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> None:
        """Split a merged table cell."""
        self.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {"slide_index": slide_index, "shape_id": shape_id, "row": row, "col": col},
        )
