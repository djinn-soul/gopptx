"""Placeholder mixin for gopptx slide library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .placeholder_collection import PlaceholderCollection

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from .placeholder import Placeholder


class SlidePlaceholderMixin:
    """Mixin providing placeholder access methods for Slide objects."""

    _BOUNDS_COMPONENTS = 4

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

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

    def list_placeholders(self) -> list[dict[str, object]]:
        """Return raw placeholder records from the bridge."""
        return self._presentation.list_placeholders(self.index)

    def set_placeholder_content(
        self,
        ph_index: int,
        ph_type: str = "",
        **kwargs: object,
    ) -> None:
        """Set content (text or image) on a placeholder. Internal method used by Placeholder."""
        text = kwargs.get("text")
        image_path = kwargs.get("image_path")
        bounds = kwargs.get("bounds")
        text_style = kwargs.get("text_style")
        self._presentation.set_placeholder_content(
            self.index,
            ph_index,
            ph_type,
            text=text if isinstance(text, str) else None,
            image_path=image_path if isinstance(image_path, str) else None,
            bounds=(
                bounds
                if isinstance(bounds, tuple) and len(bounds) == self._BOUNDS_COMPONENTS
                else None
            ),
            text_style=text_style if isinstance(text_style, dict) else None,
        )
