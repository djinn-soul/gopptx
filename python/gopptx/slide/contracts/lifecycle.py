"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ...presentation.slides.master import SlideMaster
    from ...slide.slide import Slide


class SlideLifecycleProtocol(Protocol):
    """Slide lifecycle management."""

    def slide_index_for_id(self, slide_id: int) -> int:
        """Protocol member."""
        ...

    @property
    def slides(self) -> list[Slide]:
        """Protocol member."""
        ...

    @property
    def slide_masters(self) -> list[SlideMaster]:
        """Protocol member."""
        ...

    def remove_slide(self, index: int) -> None:
        """Protocol member."""
        ...

    def update_slide(
        self,
        index: int,
        *,
        title: str | None = None,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None:
        """Protocol member."""
        ...

    def set_slide_title(self, index: int, title: str) -> None:
        """Protocol member."""
        ...

    def duplicate_slide(self, index: int, insert_at: int | None = None) -> int:
        """Protocol member."""
        ...

    def get_slide_layout_ref(self, slide_index: int) -> tuple[str, str]:
        """Protocol member."""
        ...
