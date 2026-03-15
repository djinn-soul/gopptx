"""Shape-media operations for the presentation facade."""

from __future__ import annotations

import os
from typing import TYPE_CHECKING, cast

from ... import ops
from .shape_payload_mixin import PresentationShapePayloadMixin

if TYPE_CHECKING:
    from ...schemas import ImageMetadata


class PresentationShapeMediaMixin(PresentationShapePayloadMixin):
    """Methods that add and inspect image/video/audio/OLE shapes."""

    def add_image(
        self,
        slide_index: int,
        source: str | bytes | None = None,
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
        **kwargs: object,
    ) -> int:
        """Add an image to a slide and return the created shape ID."""
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
        """Get dimensions and format metadata for an image shape."""
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
        """Add a video to a slide and return the created shape ID."""
        name = kwargs.get("name")
        poster_frame = kwargs.get("poster_frame")
        mime_type = kwargs.get("mime_type")
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if isinstance(name, str) and name:
            payload["name"] = name
        if isinstance(mime_type, str) and mime_type:
            payload["mime_type"] = mime_type

        if isinstance(poster_frame, (str, bytes, os.PathLike)):
            poster_source = cast("str | bytes | os.PathLike[str]", poster_frame)
            self._set_source_payload(
                payload,
                poster_source,
                path_key="poster_path",
                data_key="poster_data",
            )

        result = self.execute(ops.OP_ADD_VIDEO, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_audio(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Add an audio file to a slide and return the created shape ID."""
        name = kwargs.get("name")
        icon = kwargs.get("icon", kwargs.get("poster_frame"))
        mime_type = kwargs.get("mime_type")
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if isinstance(name, str) and name:
            payload["name"] = name
        if isinstance(mime_type, str) and mime_type:
            payload["mime_type"] = mime_type

        if isinstance(icon, (str, bytes, os.PathLike)):
            icon_source = cast("str | bytes | os.PathLike[str]", icon)
            self._set_source_payload(
                payload,
                icon_source,
                path_key="icon_path",
                data_key="icon_data",
            )

        result = self.execute(ops.OP_ADD_AUDIO, payload)
        return int(cast("int", result.get("shape_id", -1)))

    def add_ole_object(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Add an OLE object to a slide and return the created shape ID."""
        name = kwargs.get("name")
        prog_id = kwargs.get("prog_id")
        icon = kwargs.get("icon")
        payload = self._init_bounds_payload(slide_index, bounds)
        self._set_source_payload(payload, source)

        if isinstance(name, str) and name:
            payload["name"] = name
        if isinstance(prog_id, str) and prog_id:
            payload["prog_id"] = prog_id

        if isinstance(icon, (str, bytes, os.PathLike)):
            icon_source = cast("str | bytes | os.PathLike[str]", icon)
            self._set_source_payload(
                payload,
                icon_source,
                path_key="icon_path",
                data_key="icon_data",
            )

        result = self.execute(ops.OP_ADD_OLE_OBJECT, payload)
        return int(cast("int", result.get("shape_id", -1)))
