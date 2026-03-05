"""Paragraph-level text facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping
from typing import cast

_ALLOWED_PARAGRAPH_FIELDS = frozenset({"indent", "hanging"})
_PARAGRAPH_FIELD_ALIASES = {
    "left_margin": "indent",
    "hanging_indent": "hanging",
}


class ParagraphProps:
    """Paragraph-level controls mapped to the bridge `paragraph` payload."""

    __slots__ = ("hanging", "indent")

    def __init__(
        self,
        *,
        indent: int | None = None,
        hanging: int | None = None,
        left_margin: int | None = None,
        hanging_indent: int | None = None,
    ) -> None:
        """Initialize paragraph controls with alias support."""
        super().__init__()
        self.indent = indent
        self.hanging = hanging
        if left_margin is not None:
            self.indent = left_margin
        if hanging_indent is not None:
            self.hanging = hanging_indent

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
