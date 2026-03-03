"""Placeholder mixin for gopptx slide library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .placeholder import Placeholder

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation


class SlidePlaceholderMixin:
    """Mixin providing placeholder access methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def placeholders(self) -> list[Placeholder]:
        """List all placeholders on this slide."""
        ph_data = self._presentation.list_placeholders(self.index)
        return [
            Placeholder(
                self,  # type: ignore[arg-type]
                ph["index"],  # type: ignore[dict-item]
                ph.get("type", ""),  # type: ignore[dict-item]
                ph.get("name", ""),  # type: ignore[dict-item]
            )
            for ph in ph_data
        ]

    def get_placeholder(self, idx: int) -> Placeholder | None:
        """Get a placeholder by its index.

        Args:
            idx: The integer index of the placeholder.

        Returns:
            The Placeholder object, or None if not found.
        """
        phs = self.placeholders()
        for ph in phs:
            if ph.idx == idx:
                return ph
        return None

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
