"""Cached slide proxy access for presentation objects."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ...slide.slide import Slide
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    import threading

    from ...schemas import SlideMetadata
    from ..presentation import Presentation


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

    @property
    def slides(self) -> list[Slide]:
        """List of cached slide proxies for all slides in the presentation."""
        with self._lock:
            if self._slides_cache is None:
                self._slides_cache = [
                    Slide(cast("Presentation", self), metadata)
                    for metadata in self.slides_metadata
                ]
            return list(self._slides_cache)

    def _invalidate_slide_proxy_cache(self) -> None:
        """Clear cached slide proxies."""
        self._slides_cache = None
