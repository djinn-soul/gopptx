"""Placeholder mixin for gopptx slide library."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ..utils import is_four_number_bounds
from .placeholder_collection import PlaceholderCollection

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from .placeholder import Placeholder


class SlidePlaceholderMixin:
    """Mixin providing placeholder access methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Zero-based slide index."""
            ...

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
        typed_bounds: tuple[float, float, float, float] | None = None
        if is_four_number_bounds(bounds):
            typed_bounds = bounds
        self._presentation.set_placeholder_content(
            self.index,
            ph_index,
            ph_type,
            text=text if isinstance(text, str) else None,
            image_path=image_path if isinstance(image_path, str) else None,
            bounds=typed_bounds,
            text_style=cast("dict[str, object] | None", text_style)
            if isinstance(text_style, dict)
            else None,
        )
