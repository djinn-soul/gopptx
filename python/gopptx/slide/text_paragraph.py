"""Paragraph-level text facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping
from typing import cast

_ALLOWED_PARAGRAPH_FIELDS = frozenset({"indent", "hanging", "tab_stops"})
_PARAGRAPH_FIELD_ALIASES = {
    "left_margin": "indent",
    "hanging_indent": "hanging",
    "tabs": "tab_stops",
    "tabStops": "tab_stops",
}


class ParagraphProps:
    """Paragraph-level controls mapped to the bridge `paragraph` payload."""

    __slots__ = ("hanging", "indent", "tab_stops")

    def __init__(
        self,
        *,
        indent: int | None = None,
        hanging: int | None = None,
        left_margin: int | None = None,
        hanging_indent: int | None = None,
        **kwargs: object,
    ) -> None:
        """Initialize paragraph controls with alias support."""
        super().__init__()
        self.indent = indent
        self.hanging = hanging
        tab_stops = kwargs.get("tab_stops")
        tabs = kwargs.get("tabs")
        self.tab_stops = _coerce_optional_int_list(tab_stops)
        if left_margin is not None:
            self.indent = left_margin
        if hanging_indent is not None:
            self.hanging = hanging_indent
        if tabs is not None:
            self.tab_stops = _coerce_optional_int_list(tabs)

    @classmethod
    def from_payload(
        cls, payload: Mapping[str, object] | ParagraphProps
    ) -> ParagraphProps:
        """Build paragraph props from payload and validate keys."""
        if isinstance(payload, ParagraphProps):
            return payload
        normalized = _normalize_paragraph_mapping(payload)
        return cls(
            indent=_as_optional_int(normalized.get("indent")),
            hanging=_as_optional_int(normalized.get("hanging")),
            tab_stops=_as_optional_int_list(normalized.get("tab_stops")),
        )

    def to_payload(self) -> dict[str, object]:
        """Serialize paragraph controls for the bridge payload."""
        payload: dict[str, object] = {}
        if self.indent is not None:
            payload["indent"] = self.indent
        if self.hanging is not None:
            if self.hanging < 0:
                raise ValueError("paragraph.hanging must be >= 0")
            payload["hanging"] = self.hanging
        if self.tab_stops is not None:
            payload["tab_stops"] = _normalize_tab_stops(self.tab_stops)
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


def _as_optional_int(value: object) -> int | None:
    if value is None:
        return None
    if isinstance(value, int):
        return value
    return None


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
