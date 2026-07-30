"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ...schemas import ImageMetadata, SlideMediaRef


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

    def list_slide_media(self, slide_index: int) -> list[SlideMediaRef]:
        """Protocol member."""
        ...

    def extract_media(self, part_path: str) -> bytes:
        """Protocol member."""
        ...

    def swap_image_by_index(
        self, slide_index: int, image_index: int, data: bytes, img_format: str
    ) -> None:
        """Protocol member."""
        ...

    def swap_image_by_rel_id(
        self, slide_index: int, rel_id: str, data: bytes, img_format: str
    ) -> None:
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
