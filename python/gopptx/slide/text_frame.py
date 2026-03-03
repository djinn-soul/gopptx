"""Text-frame facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping

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

_UNSUPPORTED_KEYS = set()

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

    def __init__(  # noqa: PLR0913, D107
        self,
        *,
        margin_top: int | None = None,
        margin_bottom: int | None = None,
        margin_left: int | None = None,
        margin_right: int | None = None,
        word_wrap: bool | None = None,
        auto_fit: bool | None = None,
        auto_fit_type: str | None = None,
        vertical_align: str | None = None,
        orientation: str | None = None,
        columns: int | None = None,
        rotation: float | None = None,
        vertical_anchor: str | None = None,
        auto_size: str | None = None,
        text_direction: str | None = None,
        column_count: int | None = None,
        text_rotation: float | None = None,
    ) -> None:
        self.margin_top = margin_top
        self.margin_bottom = margin_bottom
        self.margin_left = margin_left
        self.margin_right = margin_right
        self.word_wrap = word_wrap
        self.auto_fit = auto_fit
        self.auto_fit_type = auto_fit_type
        self.vertical_align = vertical_align
        self.orientation = orientation
        self.columns = columns
        self.rotation = rotation

        if vertical_anchor is not None:
            self.vertical_align = _normalize_vertical_align(vertical_anchor)
        if auto_size is not None:
            self.auto_fit_type = _normalize_auto_fit_type(auto_size)
        if text_direction is not None:
            self.orientation = _normalize_orientation(text_direction)
        if column_count is not None:
            self.columns = column_count
        if text_rotation is not None:
            self.rotation = text_rotation

    @classmethod
    def from_payload(
        cls, payload: Mapping[str, object] | TextFrameProps
    ) -> TextFrameProps:
        """Build text-frame props from a mapping and validate keys."""
        if isinstance(payload, TextFrameProps):
            return payload
        normalized = _normalize_text_frame_mapping(payload)
        return cls(
            margin_top=_as_optional_int(normalized.get("margin_top")),
            margin_bottom=_as_optional_int(normalized.get("margin_bottom")),
            margin_left=_as_optional_int(normalized.get("margin_left")),
            margin_right=_as_optional_int(normalized.get("margin_right")),
            word_wrap=_as_optional_bool(normalized.get("word_wrap")),
            auto_fit=_as_optional_bool(normalized.get("auto_fit")),
            auto_fit_type=_as_optional_string(normalized.get("auto_fit_type")),
            vertical_align=_as_optional_string(normalized.get("vertical_align")),
            orientation=_as_optional_string(normalized.get("orientation")),
            columns=_as_optional_int(normalized.get("columns")),
            rotation=_as_optional_float(normalized.get("rotation")),
        )

    def to_payload(self) -> dict[str, object]:
        """Convert this text-frame object to bridge payload format."""
        payload: dict[str, object] = {}
        if self.margin_top is not None:
            payload["margin_top"] = self.margin_top
        if self.margin_bottom is not None:
            payload["margin_bottom"] = self.margin_bottom
        if self.margin_left is not None:
            payload["margin_left"] = self.margin_left
        if self.margin_right is not None:
            payload["margin_right"] = self.margin_right
        if self.word_wrap is not None:
            payload["word_wrap"] = self.word_wrap
        if self.auto_fit is not None:
            payload["auto_fit"] = self.auto_fit
        if self.auto_fit_type is not None:
            payload["auto_fit_type"] = _normalize_auto_fit_type(self.auto_fit_type)
        if self.vertical_align is not None:
            payload["vertical_align"] = _normalize_vertical_align(self.vertical_align)
        if self.orientation is not None:
            payload["orientation"] = _normalize_orientation(self.orientation)
        if self.columns is not None:
            if self.columns < 1:
                raise ValueError("text_frame.columns must be >= 1")
            payload["columns"] = self.columns
        if self.rotation is not None:
            if self.rotation < -360.0 or self.rotation > 360.0:
                raise ValueError("text_frame.rotation must be between -360 and 360")
            payload["rotation"] = float(self.rotation)
        return payload


def serialize_text_frame_for_payload(text_frame: object) -> object:
    """Serialize text-frame payloads and fail fast on unsupported keys."""
    if isinstance(text_frame, TextFrameProps):
        return text_frame.to_payload()
    if isinstance(text_frame, Mapping):
        props = TextFrameProps.from_payload(text_frame)
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


def _as_optional_int(value: object) -> int | None:
    if value is None:
        return None
    if isinstance(value, int):
        return value
    return None


def _as_optional_bool(value: object) -> bool | None:
    if value is None:
        return None
    if isinstance(value, bool):
        return value
    return None


def _as_optional_string(value: object) -> str | None:
    if value is None:
        return None
    if isinstance(value, str):
        return value
    return None


def _as_optional_float(value: object) -> float | None:
    if value is None:
        return None
    if isinstance(value, (int, float)):
        return float(value)
    return None
