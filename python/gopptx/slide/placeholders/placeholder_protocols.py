"""Structural protocols for placeholder proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

if TYPE_CHECKING:
    from collections.abc import Iterable

    from ...presentation.slides.master import SlideMaster


class PlaceholderPresentationProto(Protocol):
    """Structural protocol for presentation APIs used by placeholders."""

    def get_slide_layout_ref(self, slide_index: int) -> tuple[str | None, str | None]:
        """Return ``(layout_rel_id, master_rel_id)`` for the given slide index."""
        ...

    @property
    def slide_masters(self) -> Iterable[SlideMaster]:
        """Expose slide masters collection for placeholder inheritance lookups."""
        ...


class PlaceholderSlideProto(Protocol):
    """Structural protocol for slide APIs used by placeholder wrappers."""

    @property
    def presentation(self) -> PlaceholderPresentationProto:
        """Return the parent presentation facade."""
        ...

    @property
    def index(self) -> int:
        """Return the 0-based slide index."""
        ...

    def set_placeholder_content(
        self, ph_index: int, ph_type: str, **kwargs: object
    ) -> None:
        """Set placeholder content by placeholder index/type with typed payload kwargs."""
        ...


class PlaceholderCollectionSlideProto(PlaceholderSlideProto, Protocol):
    """Structural protocol for slide APIs used by placeholder collections."""

    def list_placeholders(self) -> list[dict[str, object]]:
        """List raw placeholder records for the slide."""
        ...
