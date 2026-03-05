"""Placeholder mixin for gopptx slide library."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from typing_extensions import TypeGuard

from .placeholder_collection import PlaceholderCollection

_FOUR_BOUNDS_COMPONENTS = 4

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from .placeholder import Placeholder


class SlidePlaceholderMixin:
    """Mixin providing placeholder access methods for Slide objects."""

    _BOUNDS_COMPONENTS = _FOUR_BOUNDS_COMPONENTS

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
        if _is_four_number_bounds(bounds):
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


def _is_four_number_bounds(
    value: object,
) -> TypeGuard[tuple[float, float, float, float]]:
    if not isinstance(value, tuple):
        return False
    components = cast("tuple[object, ...]", value)
    if len(components) != _FOUR_BOUNDS_COMPONENTS:
        return False
    return all(isinstance(component, int | float) for component in components)
