"""Paragraph-level text facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping
from typing import cast

from ._utils import as_optional_int, as_optional_string

_MAX_PARAGRAPH_LEVEL = 8
_ALLOWED_PARAGRAPH_FIELDS = frozenset({
    "alignment",
    "bullet_char",
    "bullet_color",
    "bullet_size_pct",
    "bullet_style",
    "hanging",
    "indent",
    "level",
    "line_spacing_pct",
    "line_spacing_pts",
    "space_after_pts",
    "space_before_pts",
    "tab_stops",
})
_PARAGRAPH_FIELD_ALIASES = {
    "hanging_indent": "hanging",
    "left_margin": "indent",
    "tabStops": "tab_stops",
    "tabs": "tab_stops",
}


class ParagraphProps:
    """Paragraph-level controls mapped to the bridge `paragraph` payload."""

    __slots__ = (
        "alignment",
        "bullet_char",
        "bullet_color",
        "bullet_size_pct",
        "bullet_style",
        "hanging",
        "indent",
        "level",
        "line_spacing_pct",
        "line_spacing_pts",
        "space_after_pts",
        "space_before_pts",
        "tab_stops",
    )

    def __init__(self, **kwargs: object) -> None:
        """Initialize paragraph controls with alias support."""
        super().__init__()
        self.indent = as_optional_int(kwargs.get("indent", kwargs.get("left_margin")))
        self.hanging = as_optional_int(
            kwargs.get("hanging", kwargs.get("hanging_indent"))
        )
        self.tab_stops = _coerce_optional_int_list(
            kwargs.get("tab_stops", kwargs.get("tabs"))
        )
        self.alignment = as_optional_string(kwargs.get("alignment"))
        self.bullet_style = as_optional_string(kwargs.get("bullet_style"))
        self.bullet_char = as_optional_string(kwargs.get("bullet_char"))
        self.bullet_color = as_optional_string(kwargs.get("bullet_color"))
        self.bullet_size_pct = as_optional_int(kwargs.get("bullet_size_pct"))
        self.level = as_optional_int(kwargs.get("level"))
        self.line_spacing_pct = as_optional_int(kwargs.get("line_spacing_pct"))
        self.line_spacing_pts = as_optional_int(kwargs.get("line_spacing_pts"))
        self.space_before_pts = as_optional_int(kwargs.get("space_before_pts"))
        self.space_after_pts = as_optional_int(kwargs.get("space_after_pts"))

    @classmethod
    def from_payload(
        cls, payload: Mapping[str, object] | ParagraphProps
    ) -> ParagraphProps:
        """Build paragraph props from payload and validate keys."""
        if isinstance(payload, ParagraphProps):
            return payload
        normalized = _normalize_paragraph_mapping(payload)
        return cls(**normalized)

    def to_payload(self) -> dict[str, object]:
        """Serialize paragraph controls for the bridge payload."""
        payload: dict[str, object] = {}
        _append_int_field(payload, "indent", self.indent)
        _append_hanging(payload, self.hanging)
        _append_tab_stops(payload, self.tab_stops)
        _append_string_field(payload, "alignment", self.alignment)
        _append_bullet_fields(
            payload,
            self.bullet_style,
            self.bullet_char,
            self.bullet_color,
            self.bullet_size_pct,
        )
        _append_level(payload, self.level)
        _append_line_spacing(payload, self.line_spacing_pct, self.line_spacing_pts)
        _append_non_negative(payload, "space_before_pts", self.space_before_pts)
        _append_non_negative(payload, "space_after_pts", self.space_after_pts)
        return payload


def serialize_paragraph_for_payload(paragraph: object) -> object:
    """Serialize paragraph payloads, supporting dict and `ParagraphProps`."""
    if isinstance(paragraph, ParagraphProps):
        return paragraph.to_payload()
    if isinstance(paragraph, Mapping):
        props = ParagraphProps.from_payload(cast("Mapping[str, object]", paragraph))
        return props.to_payload()
    return paragraph


def _normalize_paragraph_mapping(payload: Mapping[str, object]) -> dict[str, object]:
    normalized: dict[str, object] = {}
    for key, value in payload.items():
        mapped = _PARAGRAPH_FIELD_ALIASES.get(key, key)
        if mapped not in _ALLOWED_PARAGRAPH_FIELDS:
            raise ValueError(f"unknown paragraph field '{key}'")
        normalized[mapped] = value
    return normalized


def _append_int_field(payload: dict[str, object], key: str, value: int | None) -> None:
    if value is not None:
        payload[key] = value


def _append_string_field(
    payload: dict[str, object], key: str, value: str | None
) -> None:
    if value is not None:
        payload[key] = value


def _append_hanging(payload: dict[str, object], hanging: int | None) -> None:
    if hanging is None:
        return
    if hanging < 0:
        raise ValueError("paragraph.hanging must be >= 0")
    payload["hanging"] = hanging


def _append_tab_stops(payload: dict[str, object], tab_stops: list[int] | None) -> None:
    if tab_stops is None:
        return
    payload["tab_stops"] = _normalize_tab_stops(tab_stops)


def _append_level(payload: dict[str, object], level: int | None) -> None:
    if level is None:
        return
    if level < 0 or level > _MAX_PARAGRAPH_LEVEL:
        raise ValueError("paragraph.level must be between 0 and 8")
    payload["level"] = level


def _append_bullet_fields(
    payload: dict[str, object],
    bullet_style: str | None,
    bullet_char: str | None,
    bullet_color: str | None,
    bullet_size_pct: int | None,
) -> None:
    if bullet_style is None:
        return
    normalized = bullet_style.strip().lower()
    supported = {
        "none",
        "bullet",
        "number",
        "letter_lower",
        "letter_upper",
        "roman_lower",
        "roman_upper",
        "custom",
    }
    if normalized not in supported:
        raise ValueError(f"unsupported paragraph.bullet_style '{bullet_style}'")
    if normalized == "custom" and (bullet_char is None or not bullet_char.strip()):
        raise ValueError("paragraph.bullet_char is required when bullet_style='custom'")
    payload["bullet_style"] = normalized
    _append_string_field(payload, "bullet_char", bullet_char)
    _append_string_field(payload, "bullet_color", bullet_color)
    _append_non_negative(payload, "bullet_size_pct", bullet_size_pct)


def _append_line_spacing(
    payload: dict[str, object], pct: int | None, pts: int | None
) -> None:
    if pct is not None and pts is not None:
        raise ValueError(
            "paragraph cannot set both line_spacing_pct and line_spacing_pts"
        )
    _append_non_negative(payload, "line_spacing_pct", pct)
    _append_non_negative(payload, "line_spacing_pts", pts)


def _append_non_negative(
    payload: dict[str, object], key: str, value: int | None
) -> None:
    if value is None:
        return
    if value < 0:
        raise ValueError(f"paragraph.{key} must be >= 0")
    payload[key] = value


def _as_optional_int_list(value: object) -> list[int] | None:
    if value is None:
        return None
    if not isinstance(value, list):
        return None
    value_list = cast("list[object]", value)
    out: list[int] = []
    for item in value_list:
        if not isinstance(item, int):
            return None
        out.append(item)
    return out


def _coerce_optional_int_list(value: object) -> list[int] | None:
    converted = _as_optional_int_list(value)
    if converted is None and value is not None:
        raise ValueError("paragraph.tab_stops must be a list of integers")
    return converted


def _normalize_tab_stops(values: list[int]) -> list[int]:
    normalized_tab_stops: list[int] = []
    for pos in values:
        if pos < 0:
            raise ValueError("paragraph.tab_stops values must be >= 0")
        normalized_tab_stops.append(pos)
    return normalized_tab_stops
