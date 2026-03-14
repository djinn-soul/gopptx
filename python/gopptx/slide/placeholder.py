"""Placeholder proxy class for gopptx library."""
# ruff: noqa: I001,PLR1702
# pyright: reportPrivateUsage=false

from __future__ import annotations

from collections import UserString
from typing import TYPE_CHECKING, cast

from typing_extensions import override

if TYPE_CHECKING:
    from ..schemas import PlaceholderInfo
    from ..presentation.master import SlideLayout, SlideMaster
    from .slide import Slide
    from .chart_data import CategoryChartData, XyChartData


class PlaceholderFormat(UserString):
    """String-compatible placeholder format with python-pptx-like attributes."""

    def __init__(self, value: str, idx: int) -> None:
        """Initialize the placeholder format payload."""
        super().__init__(value)
        self.idx_value = idx

    @property
    def type(self) -> str:
        """Return the placeholder type token."""
        return str(self.data)

    @property
    def idx(self) -> int:
        """Return the placeholder index."""
        return self.idx_value


class Placeholder:
    """Proxy object for a placeholder within a slide."""

    def __init__(self, slide: Slide, index: int, ph_type: str, name: str) -> None:
        """Initialize the placeholder proxy.

        Args:
            slide: The Parent Slide proxy object.
            index: The zero-based index of the placeholder.
            ph_type: The placeholder type (e.g., 'body', 'title', 'pic').
            name: The human-readable name of the placeholder.
        """
        super().__init__()
        self._slide = slide
        self._index = index
        self._type = ph_type
        self._name = name

    @property
    def idx(self) -> int:
        """The index of this placeholder."""
        return self._index

    @property
    def placeholder_format(self) -> PlaceholderFormat:
        """The placeholder format object (string-compatible)."""
        return PlaceholderFormat(self._type, self._index)

    @property
    def name(self) -> str:
        """The name of this placeholder."""
        return self._name

    def _lookup_layout_and_master(
        self,
    ) -> tuple[SlideLayout | None, SlideMaster | None]:
        layout_part, master_part = self._slide._presentation.get_slide_layout_ref(  # noqa: SLF001
            self._slide.index
        )
        if not layout_part:
            return None, None

        layout_obj: SlideLayout | None = None
        master_obj: SlideMaster | None = None
        for master in self._slide._presentation.slide_masters:  # noqa: SLF001
            if master.part == master_part:
                master_obj = master
            for layout in master.slide_layouts:
                if layout.part == layout_part:
                    layout_obj = layout
                    if master_obj is None:
                        master_obj = master
                    break
            if layout_obj is not None and master_obj is not None:
                break
        return layout_obj, master_obj

    def _find_placeholder_info(
        self, entries: list[PlaceholderInfo]
    ) -> PlaceholderInfo | None:
        for entry in entries:
            idx = int(entry.get("Index", entry.get("index", -1)))
            ph_type = str(entry.get("Type", entry.get("type", "")))
            if idx == self._index and (ph_type == self._type or not ph_type):
                return entry
        for entry in entries:
            idx = int(entry.get("Index", entry.get("index", -1)))
            if idx == self._index:
                return entry
        return None

    @property
    def layout_placeholder(self) -> PlaceholderInfo | None:
        """Layout placeholder metadata matching this slide placeholder."""
        layout_obj, _ = self._lookup_layout_and_master()
        if layout_obj is None:
            return None
        return self._find_placeholder_info(layout_obj.placeholders)

    @property
    def master_placeholder(self) -> PlaceholderInfo | None:
        """Master placeholder metadata matching this slide placeholder."""
        _, master_obj = self._lookup_layout_and_master()
        if master_obj is None:
            return None
        return self._find_placeholder_info(master_obj.placeholders)

    def insert_text(self, text: str, **style_kwargs: object) -> None:
        """Replace the placeholder with text.

        Args:
            text: The text to insert.
            **style_kwargs: Optional text style properties (size_pt, bold, italic, color, font).
        """
        # Normalize style keys (e.g. size -> size_pt, colour -> color)
        text_style = {}
        for k, v in style_kwargs.items():
            key = k
            if k in {"size", "font_size"}:
                key = "size_pt"
            elif k == "font_name":
                key = "font"
            elif k == "colour":  # Handle British spelling
                key = "color"
            text_style[key] = v

        self._slide.set_placeholder_content(
            self.idx,
            self._type,
            text=text,
            text_style=cast("dict[str, object]", text_style),
        )

    def insert_picture(
        self,
        image_path: str,
        bounds: tuple[float, float, float, float] | None = None,
    ) -> None:
        """Replace the placeholder with a picture.

        Args:
            image_path: Path to the image file.
            bounds: Optional (x, y, width, height) in points relative to the placeholder.
        """
        self._slide.set_placeholder_content(
            self.idx, self._type, image_path=image_path, bounds=bounds
        )

    def insert_table(
        self,
        rows: int,
        cols: int,
        bounds: tuple[float, float, float, float] | None = None,
    ) -> None:
        """Replace the placeholder with a table.

        Args:
            rows: Number of rows in the table.
            cols: Number of columns in the table.
            bounds: Optional (x, y, width, height) in points relative to the placeholder.
        """
        self._slide.set_placeholder_content(
            self.idx,
            self._type,
            table_rows=rows,
            table_cols=cols,
            bounds=bounds,
        )

    def insert_chart(
        self,
        chart_type: str,
        categories: list[str] | CategoryChartData | XyChartData,
        values_or_series: list[float]
        | list[dict[str, str | list[float]]]
        | None = None,
        bounds: tuple[float, float, float, float] | None = None,
        **kwargs: object,
    ) -> None:
        """Replace the placeholder with a chart.

        Args:
            chart_type: Type of chart (e.g., 'bar', 'column', 'line', 'pie').
            categories: List of category labels for the chart.
            values_or_series: Data values or series data for the chart.
            bounds: Optional (x, y, width, height) in points relative to the placeholder.
            **kwargs: Additional chart options (title, legend, etc.).
        """
        if hasattr(categories, "to_add_chart_args"):
            chart_builder = cast("CategoryChartData | XyChartData", categories)
            builder_categories, builder_values = chart_builder.to_add_chart_args()
            categories = builder_categories
            values_or_series = cast(
                "list[float] | list[dict[str, str | list[float]]]",
                builder_values,
            )
        if values_or_series is None:
            values_or_series = []

        self._slide.set_placeholder_content(
            self.idx,
            self._type,
            chart_type=chart_type,
            chart_categories=cast("list[str]", categories),
            chart_values=values_or_series,
            bounds=bounds,
            chart_options=kwargs or None,
        )

    @override
    def __repr__(self) -> str:
        """Return a string representation of this placeholder."""
        return f"<Placeholder idx={self.idx} type='{self.placeholder_format}' name='{self.name}'>"


class TitlePlaceholder(Placeholder):
    """Placeholder subtype for title-like placeholders."""


class BodyPlaceholder(Placeholder):
    """Placeholder subtype for body/content placeholders."""


class PicturePlaceholder(Placeholder):
    """Placeholder subtype for picture placeholders."""


class ChartPlaceholder(Placeholder):
    """Placeholder subtype for chart placeholders."""


class TablePlaceholder(Placeholder):
    """Placeholder subtype for table placeholders."""


_PLACEHOLDER_TYPE_TO_CLASS: dict[str, type[Placeholder]] = {
    "title": TitlePlaceholder,
    "ctrTitle": TitlePlaceholder,
    "body": BodyPlaceholder,
    "obj": BodyPlaceholder,
    "pic": PicturePlaceholder,
    "chart": ChartPlaceholder,
    "tbl": TablePlaceholder,
}


def create_placeholder(
    slide: Slide, index: int, ph_type: str, name: str
) -> Placeholder:
    """Create a placeholder proxy using the most-specific subtype mapping."""
    cls = _PLACEHOLDER_TYPE_TO_CLASS.get(ph_type, Placeholder)
    return cls(slide, index, ph_type, name)
