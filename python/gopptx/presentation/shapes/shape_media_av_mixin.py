"""Video, audio, and OLE-object operations for the presentation facade."""

from __future__ import annotations

import os
from typing import cast

from ... import ops
from ..helpers import get_required_int
from .shape_payload_mixin import PresentationShapePayloadMixin


class PresentationShapeAVMixin(PresentationShapePayloadMixin):
    """Methods that add video, audio, and embedded OLE object shapes."""

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
        return get_required_int(result, "shape_id")

    def add_online_video(
        self,
        slide_index: int,
        url: str,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Link a slide to a hosted video and return the created shape ID.

        This is PowerPoint's Insert > Media > Video > Online Video: the media is
        not embedded, the slide keeps an external relationship to ``url``
        (upstream python-pptx issue #1034).

        Args:
            slide_index: Zero-based slide index.
            url: http(s) URL of the video.
            bounds: (left, top, width, height) in EMU.
            **kwargs: ``poster_frame`` (image path or bytes shown as the click
                target), ``poster_format`` (extension for poster bytes, e.g.
                "png") and ``alt_text`` (accessible description).
        """
        payload = self._init_bounds_payload(slide_index, bounds)
        payload["url"] = url

        poster_frame = kwargs.get("poster_frame")
        if isinstance(poster_frame, (str, bytes, os.PathLike)):
            poster_source = cast("str | bytes | os.PathLike[str]", poster_frame)
            self._set_source_payload(
                payload,
                poster_source,
                path_key="poster_path",
                data_key="poster_data",
            )
        poster_format = kwargs.get("poster_format")
        if isinstance(poster_format, str) and poster_format:
            payload["poster_format"] = poster_format
        alt_text = kwargs.get("alt_text")
        if isinstance(alt_text, str) and alt_text:
            payload["alt_text"] = alt_text

        result = self.execute(ops.OP_ADD_ONLINE_VIDEO, payload)
        return get_required_int(result, "shape_id")

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
        return get_required_int(result, "shape_id")

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
        return get_required_int(result, "shape_id")
