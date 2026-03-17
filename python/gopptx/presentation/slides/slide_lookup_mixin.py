"""Slide ID lookup cache for presentation proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING

from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    import threading

    from ...schemas import SlideMetadata


class PresentationSlideLookupMixin(PresentationMixinBase):
    """Mixin providing cached slide-ID to index lookup."""

    if TYPE_CHECKING:
        _lock: threading.RLock

        @property
        def slides_metadata(self) -> list[SlideMetadata]:
            """Slide metadata snapshot."""
            ...

    def __init__(self) -> None:
        """Initialize slide-lookup cache state."""
        super().__init__()
        self._slide_index_map_cache: dict[int, int] | None = None

    def slide_index_for_id(self, slide_id: int) -> int:
        """Resolve a slide index by stable slide ID using a cached lookup map."""
        with self._lock:
            if self._slide_index_map_cache is None:
                self._slide_index_map_cache = {
                    int(item["SlideID"]): int(item["Index"])
                    for item in self.slides_metadata
                }
            return self._slide_index_map_cache.get(slide_id, -1)

    def _invalidate_slide_lookup_cache(self) -> None:
        """Clear the slide lookup cache."""
        self._slide_index_map_cache = None
