from __future__ import annotations

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class SlideSize(TypedDict):
    width: int
    height: int


class PresentationMetadata(TypedDict):
    title: str
    slide_count: int
    size: SlideSize

