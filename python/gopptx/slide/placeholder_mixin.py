"""Placeholder mixin for gopptx slide library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .placeholder import Placeholder
from .placeholder_collection import PlaceholderCollection

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation


class SlidePlaceholderMixin:
    """Mixin providing placeholder access methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    @property
    def placeholders(self) -> PlaceholderCollection:
        """Collection of placeholders on this slide."""
        return PlaceholderCollection(self)  # type: ignore[arg-type]

    def get_placeholder(self, idx: int) -> Placeholder | None:
        """Get a placeholder by its index.

        Args:
            idx: The integer index of the placeholder.

        Returns:
            The Placeholder object, or None if not found.
        """
        return self.placeholders.get(idx)

    def set_placeholder_content(
        self,
        ph_index: int,
        ph_type: str = "",
        text: str | None = None,
        image_path: str | None = None,
        bounds: tuple[float, float, float, float] | None = None,
        text_style: dict[str, object] | None = None,
    ) -> None:
        """Set content (text or image) on a placeholder. Internal method used by Placeholder."""
        self._presentation.set_placeholder_content(
            self.index, ph_index, ph_type, text, image_path, bounds, text_style
        )
