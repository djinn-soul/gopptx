"""Text-frame facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping
from typing import cast

from ._utils import (
    as_optional_bool,
    as_optional_float,
    as_optional_int,
    as_optional_string,
)

_as_optional_int = as_optional_int
_as_optional_bool = as_optional_bool
_as_optional_string = as_optional_string
_as_optional_float = as_optional_float

_SUPPORTED_KEYS = {
    "margin_top",
    "margin_bottom",
    "margin_left",
    "margin_right",
    "word_wrap",
    "auto_fit",
    "auto_fit_type",
    "vertical_align",
    "orientation",
    "columns",
    "rotation",
}

_UNSUPPORTED_KEYS: set[str] = set()

_KEY_ALIASES = {
    "marginTop": "margin_top",
    "marginBottom": "margin_bottom",
    "marginLeft": "margin_left",
    "marginRight": "margin_right",
    "margin-top": "margin_top",
    "margin-bottom": "margin_bottom",
    "margin-left": "margin_left",
    "margin-right": "margin_right",
    "vertical_anchor": "vertical_align",
    "auto_size": "auto_fit_type",
    "column_count": "columns",
    "text_direction": "orientation",
    "text_rotation": "rotation",
}

_VERTICAL_ALIGN_ALIASES = {
    "top": "top",
    "middle": "ctr",
    "center": "ctr",
    "ctr": "ctr",
    "bottom": "bot",
    "bot": "bot",
}

_AUTO_FIT_ALIASES = {
    "none": "none",
    "normal": "normal",
    "shape": "shape",
    "spautofit": "shape",
    "normautofit": "normal",
    "shape_to_fit_text": "shape",
    "text_to_fit_shape": "normal",
}

_ORIENTATION_ALIASES = {
    "horz": "horz",
    "horizontal": "horz",
    "vert": "vert",
    "vertical": "vert",
    "vert270": "vert270",
    "vertical270": "vert270",
    "vertical_270": "vert270",
    "wordartvert": "wordArtVert",
    "word_art_vert": "wordArtVert",
    "eavert": "eaVert",
    "ea_vert": "eaVert",
    "mongolianvert": "mongolianVert",
    "mongolian_vert": "mongolianVert",
    "wordartvertrtl": "wordArtVertRtl",
    "word_art_vert_rtl": "wordArtVertRtl",
}

_FULL_ROTATION_DEGREES = 360.0


class TextFrameProps:
    """Mutable text-frame object with normalized payload conversion."""

    __slots__ = (
        "auto_fit",
        "auto_fit_type",
        "columns",
        "margin_bottom",
        "margin_left",
        "margin_right",
        "margin_top",
        "orientation",
        "rotation",
        "vertical_align",
        "word_wrap",
    )

    def __init__(self, **kwargs: object) -> None:
        """Initialize text frame properties from keyword arguments."""
        super().__init__()
        self.margin_top = as_optional_int(kwargs.get("margin_top"))
        self.margin_bottom = as_optional_int(kwargs.get("margin_bottom"))
        self.margin_left = as_optional_int(kwargs.get("margin_left"))
        self.margin_right = as_optional_int(kwargs.get("margin_right"))
        self.word_wrap = as_optional_bool(kwargs.get("word_wrap"))
        self.auto_fit = as_optional_bool(kwargs.get("auto_fit"))

        self.auto_fit_type = as_optional_string(kwargs.get("auto_fit_type"))
        self.vertical_align = as_optional_string(kwargs.get("vertical_align"))
        self.orientation = as_optional_string(kwargs.get("orientation"))

        self.columns = as_optional_int(kwargs.get("columns"))
        self.rotation = as_optional_float(kwargs.get("rotation"))

        vertical_anchor = kwargs.get("vertical_anchor")
        if vertical_anchor is not None:
            self.vertical_align = str(vertical_anchor)
        # Legacy alias handling - auto_size overrides auto_fit_type for backwards compatibility.
        # See _ALIAS_MAP for parameter name mappings.
        auto_size = kwargs.get("auto_size")
        if auto_size is not None:
            self.auto_fit_type = str(auto_size)
        text_direction = kwargs.get("text_direction")
        if text_direction is not None:
            self.orientation = str(text_direction)
        column_count = kwargs.get("column_count")
        if column_count is not None:
            self.columns = as_optional_int(column_count)
        text_rotation = kwargs.get("text_rotation")
        if text_rotation is not None:
            self.rotation = as_optional_float(text_rotation)

        if self.auto_fit_type is not None:
            self.auto_fit_type = _normalize_auto_fit_type(self.auto_fit_type)
        if self.vertical_align is not None:
            self.vertical_align = _normalize_vertical_align(self.vertical_align)
        if self.orientation is not None:
            self.orientation = _normalize_orientation(self.orientation)

    @classmethod
    def from_payload(
        cls, payload: Mapping[str, object] | TextFrameProps
    ) -> TextFrameProps:
        """Build text-frame props from a mapping and validate keys."""
        if isinstance(payload, TextFrameProps):
            return payload
        normalized = _normalize_text_frame_mapping(payload)
        return cls(
            margin_top=as_optional_int(normalized.get("margin_top")),
            margin_bottom=as_optional_int(normalized.get("margin_bottom")),
            margin_left=as_optional_int(normalized.get("margin_left")),
            margin_right=as_optional_int(normalized.get("margin_right")),
            word_wrap=as_optional_bool(normalized.get("word_wrap")),
            auto_fit=as_optional_bool(normalized.get("auto_fit")),
            auto_fit_type=as_optional_string(normalized.get("auto_fit_type")),
            vertical_align=as_optional_string(normalized.get("vertical_align")),
            orientation=as_optional_string(normalized.get("orientation")),
            columns=as_optional_int(normalized.get("columns")),
            rotation=as_optional_float(normalized.get("rotation")),
        )

    def to_payload(self) -> dict[str, object]:
        """Convert this text-frame object to bridge payload format."""
        payload: dict[str, object] = _collect_base_payload(self)
        _append_layout_payload(payload, self)
        return payload


def _collect_base_payload(props: TextFrameProps) -> dict[str, object]:
    payload: dict[str, object] = {}
    if props.margin_top is not None:
        payload["margin_top"] = props.margin_top
    if props.margin_bottom is not None:
        payload["margin_bottom"] = props.margin_bottom
    if props.margin_left is not None:
        payload["margin_left"] = props.margin_left
    if props.margin_right is not None:
        payload["margin_right"] = props.margin_right
    if props.word_wrap is not None:
        payload["word_wrap"] = props.word_wrap
    if props.auto_fit is not None:
        payload["auto_fit"] = props.auto_fit
    if props.auto_fit_type is not None:
        payload["auto_fit_type"] = _normalize_auto_fit_type(props.auto_fit_type)
    if props.vertical_align is not None:
        payload["vertical_align"] = _normalize_vertical_align(props.vertical_align)
    if props.orientation is not None:
        payload["orientation"] = _normalize_orientation(props.orientation)
    return payload


def _append_layout_payload(payload: dict[str, object], props: TextFrameProps) -> None:
    if props.columns is not None:
        if props.columns < 1:
            raise ValueError("text_frame.columns must be >= 1")
        payload["columns"] = props.columns
    if props.rotation is not None:
        if (
            props.rotation < -_FULL_ROTATION_DEGREES
            or props.rotation > _FULL_ROTATION_DEGREES
        ):
            raise ValueError("text_frame.rotation must be between -360 and 360")
        payload["rotation"] = float(props.rotation)


def serialize_text_frame_for_payload(text_frame: object) -> object:
    """Serialize text-frame payloads and fail fast on unsupported keys."""
    if isinstance(text_frame, TextFrameProps):
        return text_frame.to_payload()
    if isinstance(text_frame, Mapping):
        props = TextFrameProps.from_payload(cast("Mapping[str, object]", text_frame))
        return props.to_payload()
    return text_frame


def _normalize_text_frame_mapping(payload: Mapping[str, object]) -> dict[str, object]:
    normalized: dict[str, object] = {}
    for key, value in payload.items():
        mapped = _KEY_ALIASES.get(key, key)
        if mapped in _UNSUPPORTED_KEYS:
            raise ValueError(f"text_frame field '{key}' is not supported yet")
        if mapped not in _SUPPORTED_KEYS:
            raise ValueError(f"unknown text_frame field '{key}'")
        if mapped == "vertical_align" and isinstance(value, str):
            normalized[mapped] = _normalize_vertical_align(value)
            continue
        if mapped == "auto_fit_type" and isinstance(value, str):
            normalized[mapped] = _normalize_auto_fit_type(value)
            continue
        if mapped == "orientation" and isinstance(value, str):
            normalized[mapped] = _normalize_orientation(value)
            continue
        normalized[mapped] = value
    return normalized


def _normalize_vertical_align(value: str) -> str:
    lowered = value.strip().lower()
    mapped = _VERTICAL_ALIGN_ALIASES.get(lowered)
    if mapped is None:
        raise ValueError(f"unsupported vertical alignment '{value}'")
    return mapped


def _normalize_auto_fit_type(value: str) -> str:
    lowered = value.strip().lower()
    mapped = _AUTO_FIT_ALIASES.get(lowered)
    if mapped is None:
        raise ValueError(f"unsupported auto_fit_type '{value}'")
    return mapped


def _normalize_orientation(value: str) -> str:
    lowered = value.strip().lower()
    mapped = _ORIENTATION_ALIASES.get(lowered)
    if mapped is None:
        raise ValueError(f"unsupported orientation '{value}'")
    return mapped
