"""Shared payload helpers for shape-related presentation mixins."""

from __future__ import annotations

import base64
import os
from typing import TYPE_CHECKING

from ...slide.text.text_frame import serialize_text_frame_for_payload
from ...slide.text.text_paragraph import serialize_paragraph_for_payload
from ...slide.text.text_run import serialize_runs_for_payload
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from collections.abc import Mapping


class PresentationShapePayloadMixin(PresentationMixinBase):
    """Utilities for building normalized shape command payloads."""

    _RECT_BOUNDS_COMPONENTS = 4

    @staticmethod
    def _init_bounds_payload(
        slide_index: int, bounds: tuple[float, float, float, float]
    ) -> dict[str, object]:
        """Build a base payload carrying slide index and rectangle bounds."""
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
        """Map source path/bytes/path-like input into payload fields."""
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
        options: Mapping[str, object],
        *,
        include_text: bool,
    ) -> None:
        """Apply typed payload options and serializers for shape command calls."""
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
