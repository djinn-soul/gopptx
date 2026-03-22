"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ...schemas import ImageMetadata


class MediaOperationsProtocol(Protocol):
    """Image and media management."""

    def add_image(
        self,
        slide_index: int,
        path: str | None,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def get_image_metadata(self, slide_index: int, shape_id: int) -> ImageMetadata:
        """Protocol member."""
        ...

    def add_video(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        """Protocol member."""
        ...

    def add_audio(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        """Protocol member."""
        ...

    def add_ole_object(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        prog_id: str | None = None,
        icon: str | bytes | None = None,
    ) -> int:
        """Protocol member."""
        ...

    def add_mermaid(
        self, slide_index: int, diagram: str, *, theme: str = ""
    ) -> tuple[int, int]:
        """Protocol member."""
        ...
