"""Cached slide proxy access for presentation objects."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ...slide.slide import Slide
from ..helpers import PresentationMixinBase
from .collection import Slides

if TYPE_CHECKING:
    import threading

    from ...schemas import SlideMetadata
    from ...slide.contracts import SlidePresentationProtocol


class PresentationSlideProxyMixin(PresentationMixinBase):
    """Mixin providing cached slide proxy lists."""

    if TYPE_CHECKING:
        _lock: threading.RLock

        @property
        def slides_metadata(self) -> list[SlideMetadata]:
            """Slide metadata snapshot."""
            ...

    def __init__(self) -> None:
        """Initialize slide-proxy cache state."""
        super().__init__()
        self._slides_cache: list[Slide] | None = None
        self._slides_collection_cache: Slides | None = None

    @property
    def slides(self) -> Slides:
        """List of cached slide proxies for all slides in the presentation."""
        with self._lock:
            if self._slides_cache is None:
                self._slides_cache = [
                    Slide(cast("SlidePresentationProtocol", self), metadata)
                    for metadata in self.slides_metadata
                ]
            self._slides_collection_cache = Slides(
                cast("SlidePresentationProtocol", self),
                list(self._slides_cache),
            )
            return self._slides_collection_cache

    def _invalidate_slide_proxy_cache(self) -> None:
        """Clear cached slide proxies."""
        self._slides_cache = None
        self._slides_collection_cache = None
