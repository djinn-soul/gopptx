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
        """Set content (text, image, table, or chart) on a placeholder. Internal method used by Placeholder."""
        text = kwargs.get("text")
        image_path = kwargs.get("image_path")
        bounds = kwargs.get("bounds")
        text_style = kwargs.get("text_style")
        table_rows = kwargs.get("table_rows")
        table_cols = kwargs.get("table_cols")
        chart_type = kwargs.get("chart_type")
        chart_categories = kwargs.get("chart_categories")
        chart_values = kwargs.get("chart_values")
        chart_options = kwargs.get("chart_options")

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
            table_rows=int(table_rows) if isinstance(table_rows, int) else None,
            table_cols=int(table_cols) if isinstance(table_cols, int) else None,
            chart_type=chart_type if isinstance(chart_type, str) else None,
            chart_categories=chart_categories if isinstance(chart_categories, list) else None,
            chart_values=chart_values if isinstance(chart_values, list) else None,
            chart_options=cast("dict[str, object] | None", chart_options)
            if isinstance(chart_options, dict)
            else None,
        )
