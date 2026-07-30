"""Grayscale-operation typed schema definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .constants import PlaceholderType


class GrayscaleShapeRef(TypedDict):
    """Shape target for grayscale conversion."""

    slide_index: int
    shape_id: int


class GrayscaleTextRef(TypedDict):
    """Text target for grayscale conversion."""

    slide_index: int
    shape_id: int
    run_indices: NotRequired[list[int]]


class GrayscalePlaceholderRef(TypedDict):
    """Placeholder target for grayscale conversion."""

    slide_index: int
    type: NotRequired[PlaceholderType]
    index: NotRequired[int | None]


class GrayscaleScope(TypedDict, total=False):
    """Which content types to convert to grayscale."""

    colors: bool
    images: bool
    backgrounds: bool
