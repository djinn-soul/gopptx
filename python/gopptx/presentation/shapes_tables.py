"""Shape, text, notes, and table mixins for the Presentation API."""

from __future__ import annotations

import base64
import os
from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from ..slide.freeform_builder import FreeformBuilder
from ..slide.text_frame import serialize_text_frame_for_payload
from ..slide.text_paragraph import serialize_paragraph_for_payload
from ..slide.text_run import serialize_runs_for_payload
from ..utils import normalize_table_index
from .helpers import PresentationProtocol

if TYPE_CHECKING:
    from ..schemas import (
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

    def get_notes_payload(self, slide_index: int) -> dict[str, object]:
        """Fetch raw notes payload from bridge with null-safe fallback."""
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        if not isinstance(result, dict):
            return {"text": "", "notes_slide": None}
        return result

    def get_notes(self, slide_index: int) -> str:
        """Get speaker notes for a slide."""
        result = self.get_notes_payload(slide_index)
        return str(cast("str", result.get("text", "")))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Set speaker notes for a slide."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})


class PresentationShapeMixin(PresentationProtocol):
    """Mixin providing shape manipulation methods."""

    _RECT_BOUNDS_COMPONENTS = 4

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

    @staticmethod
    def _apply_shape_payload_options(
        payload: dict[str, object],
        options: dict[str, object],
        *,
        include_text: bool,
    ) -> None:
        serializers = {
            "runs": serialize_runs_for_payload,
            "text_frame": serialize_text_frame_for_payload,
            "paragraph": serialize_paragraph_for_payload,
        }
        keys = (
            "text",
            *serializers.keys(),
            "click_action",
            "hover_action",
            "properties",
        )
        for key in keys:
            value = options.get(key)
            if value is None:
                continue
            if key == "text" and (not include_text or not isinstance(value, str)):
                continue
            serializer = serializers.get(key)
            payload[key] = serializer(value) if serializer is not None else value

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
        **kwargs: str | ShapeProps,
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
        **kwargs: str | ShapeProps,
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

    def add_image(
        self,
        slide_index: int,
        source: str | bytes | None = None,
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
        **kwargs: object,
    ) -> int:
        """Add an image to a slide.

        Args:
            slide_index: Slide index (0-based).
            path: Path to the image file, or None if data is provided.
            bounds: (x, y, w, h) in EMU.
            data: Raw image bytes.
            source: Direct path/bytes source when not using path/data args.
            image_format: Image format (e.g., 'png', 'jpeg') when source/data is bytes.
            img_format: Backward-compatible alias for image_format.
            crop: ImageCrop dictionary.
            rotation: Rotation in degrees.
            flip_h: Flip horizontally.
            flip_v: Flip vertically.
            **kwargs: Optional source/path/data/options payload members.
        """
        path = kwargs.get("path")
        data = kwargs.get("data")
        image_format = kwargs.get("image_format")
        img_format = kwargs.get("img_format")
        crop = kwargs.get("crop")
        rotation = kwargs.get("rotation")
        flip_h = kwargs.get("flip_h")
        flip_v = kwargs.get("flip_v")
        payload = self._init_bounds_payload(slide_index, bounds)
        if source:
            self._set_source_payload(payload, source)
        elif isinstance(path, str):
            self._set_source_payload(payload, path)
        elif isinstance(data, bytes):
            self._set_source_payload(payload, data)
            resolved_format = (
                image_format
                if isinstance(image_format, str)
                else img_format
                if isinstance(img_format, str)
                else None
            )
            if resolved_format:
                payload["format"] = resolved_format

        options: dict[str, object] = {}
        if isinstance(crop, dict):
            options["crop"] = cast("dict[str, object]", crop)
        if isinstance(rotation, int | float):
            options["rotation"] = rotation
        if isinstance(flip_h, bool):
            options["flip_h"] = flip_h
        if isinstance(flip_v, bool):
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
        **kwargs: object,
    ) -> int:
        """Add a video to a slide."""
        name = kwargs.get("name")
        poster_frame = kwargs.get("poster_frame")
        mime_type = kwargs.get("mime_type")
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if isinstance(name, str) and name:
            payload["name"] = name
        if isinstance(mime_type, str) and mime_type:
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
        **kwargs: object,
    ) -> int:
        """Add an OLE object to a slide."""
        name = kwargs.get("name")
        prog_id = kwargs.get("prog_id")
        icon = kwargs.get("icon")
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if isinstance(name, str) and name:
            payload["name"] = name
        if isinstance(prog_id, str) and prog_id:
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

    def set_table_style(self, slide_index: int, shape_id: int, style_guid: str) -> None:
        """Apply a table style to a table.

        The style_guid must be a valid PowerPoint table style GUID, e.g.:
            "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}" - Medium Style 2 - Accent 1
            "{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}" - Light Style 1
            "{5940675A-B579-460E-94D1-54222C63F5DA}" - Medium Style 1 - Accent 1
        """
        self.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "style_guid": style_guid,
            },
        )

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
