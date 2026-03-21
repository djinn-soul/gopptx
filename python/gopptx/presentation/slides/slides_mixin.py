"""Presentation slides mixin for gopptx library."""

from __future__ import annotations

import uuid
from typing import TYPE_CHECKING, cast

from ... import ops
from ...slide.slide import Slide
from ...utils import is_four_number_bounds
from ..placeholders import (
    build_placeholder_chart_payload,
    build_placeholder_table_payload,
)
from .layout_theme_mixin import PresentationLayoutMixin, PresentationThemeMixin
from .master import SlideMasters
from .properties_mixin import PresentationPropertiesMixin
from .sections_mixin import PresentationSectionMixin
from .slide_layout_enum import SlideLayoutType

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import (
        SlideMetadata,
    )
    from ..helpers import PresentationProtocol
    from ..presentation import Presentation


class PresentationSlidesMixin(
    PresentationPropertiesMixin,
    PresentationSectionMixin,
    PresentationThemeMixin,
    PresentationLayoutMixin,
):
    """Mixin providing slide-related methods for Presentation."""

    _BOUNDS_COMPONENTS = 4
    _EMU_PER_INCH = 914400

    if TYPE_CHECKING:

        @property
        def slides(self) -> list[Slide]:
            """Return materialized slide objects for the current presentation."""
            ...

        @property
        def slide_count(self) -> int:
            """Return the total number of slides in the presentation."""
            ...

    def list_placeholders(self, slide_index: int) -> list[dict[str, object]]:
        """Bridge op: list all placeholders on a slide."""
        result = self.execute(
            ops.OP_LIST_PLACEHOLDERS,
            {"slide_index": slide_index},
        )
        return cast("list[dict[str, object]]", result.get("placeholders", []))

    def get_slide_layout_ref(self, slide_index: int) -> tuple[str, str]:
        """Return (layout_part, master_part) for a slide."""
        result = self.execute(
            ops.OP_GET_SLIDE_LAYOUT_REF,
            {"slide_index": slide_index},
        )
        layout_part = str(result.get("layout_part", ""))
        master_part = str(result.get("master_part", ""))
        return layout_part, master_part

    def set_placeholder_content(
        self,
        slide_index: int,
        ph_index: int,
        ph_type: str = "",
        **kwargs: object,
    ) -> None:
        """Bridge op: insert rich content (text, image, table, or chart) into a placeholder."""
        bounds = kwargs.get("bounds")

        payload: dict[str, object] = {
            "slide_index": slide_index,
            "index": ph_index,
            "ph_type": ph_type,
        }
        text = kwargs.get("text")
        if isinstance(text, str):
            payload["text"] = text
        image_path = kwargs.get("image_path")
        if isinstance(image_path, str):
            payload["image_path"] = image_path
        if is_four_number_bounds(bounds):
            payload["bounds"] = list(bounds)
        text_style = kwargs.get("text_style")
        if isinstance(text_style, dict):
            payload["text_style"] = text_style
        force_rect = kwargs.get("force_rect_geometry")
        if isinstance(force_rect, bool):
            payload["force_rect_geometry"] = force_rect
        table_payload = build_placeholder_table_payload(
            kwargs.get("table"),
            kwargs.get("table_rows"),
            kwargs.get("table_cols"),
        )
        if table_payload is not None:
            payload["table"] = table_payload
        chart_payload = build_placeholder_chart_payload(
            kwargs,
            bounds,
        )
        if chart_payload is not None:
            payload["chart"] = chart_payload

        self.execute(ops.OP_SET_PLACEHOLDER_CONTENT, payload)
        self.invalidate_cache()

    _slide_masters_obj: SlideMasters | None = None

    @property
    def slide_masters(self) -> SlideMasters:
        """Get the slide masters collection."""
        if self._slide_masters_obj is None:
            self._slide_masters_obj = SlideMasters(cast("PresentationProtocol", self))
        return self._slide_masters_obj

    def add_slide(
        self, title: str, layout: str | None = None, bullets: list[str] | None = None
    ) -> Slide:
        """Add a new slide to the presentation.

        Args:
            title: Title for the slide.
            layout: Slide layout type (use SlideLayoutType constants).
                   Examples: SlideLayoutType.BLANK, SlideLayoutType.TITLE_AND_CONTENT
            bullets: Optional list of bullet points for content.

        Returns:
            The newly created Slide object.

        Raises:
            ValueError: If layout is invalid.
        """
        payload: dict[str, object] = {"title": title}
        if layout:
            # Validate layout using enum
            try:
                validated_layout = SlideLayoutType.validate(layout)
                payload["layout"] = validated_layout
            except ValueError as e:
                raise ValueError(f"Invalid slide layout: {e}") from e
        if bullets:
            payload["bullets"] = bullets
        result = self.execute(ops.OP_ADD_SLIDE, payload)
        if result.get("_batched", False):
            # Batch mode optimization: return a Slide object with a dummy index
            # This allows safe construction of slide layouts/designs within a batch
            placeholder_metadata = {
                "Title": title,
                "SlideID": -1,
                "RelationshipID": str(uuid.uuid4()),
                "PartName": "/ppt/slides/slide_placeholder.xml",
                "Index": -1,
            }
            return Slide(
                cast("Presentation", self),
                cast("SlideMetadata", placeholder_metadata),
            )
        self.invalidate_cache()
        return self.slides[int(cast("int", result.get("index", -1)))]

    def remove_slide(self, index: int) -> None:
        """Remove a slide from the presentation."""
        self.execute(ops.OP_REMOVE_SLIDE, {"index": index})
        self.invalidate_cache()

    def move_slide(self, from_index: int, to_index: int) -> None:
        """Move a slide to a new position."""
        self.execute(ops.OP_MOVE_SLIDE, {"from": from_index, "to": to_index})
        self.invalidate_cache()

    def duplicate_slide(self, index: int, insert_at: int | None = None) -> int:
        """Duplicate a slide and return the new slide index."""
        if insert_at is None:
            insert_at = index + 1
        result = self.execute(
            ops.OP_DUPLICATE_SLIDE, {"index": index, "insert_at": insert_at}
        )
        self.invalidate_cache()
        return int(cast("int", result.get("new_index", -1)))

    def duplicate_slide_after(self, index: int) -> int:
        """Duplicate slide at *index* and insert it immediately after the original.

        Returns the zero-based index of the newly inserted slide.
        """
        result = self.execute(ops.OP_DUPLICATE_SLIDE_AFTER, {"slide_index": index})
        self.invalidate_cache()
        return int(cast("int", result.get("new_index", -1)))

    def update_slide(
        self,
        index: int,
        title: str | None = None,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None:
        """Update slide properties.

        Args:
            index: Zero-based slide index.
            title: New title for the slide.
            layout: New slide layout type (use SlideLayoutType constants).
                   Examples: SlideLayoutType.BLANK, SlideLayoutType.TITLE_AND_CONTENT
            bullets: New bullet points for content.

        Raises:
            ValueError: If layout is invalid.
        """
        payload: dict[str, object] = {"slide_index": index}
        if title is not None:
            payload["title"] = title
        if layout is not None:
            # Validate layout using enum
            try:
                validated_layout = SlideLayoutType.validate(layout)
                payload["layout"] = validated_layout
            except ValueError as e:
                raise ValueError(f"Invalid slide layout: {e}") from e
        if bullets is not None:
            payload["bullets"] = bullets
        self.execute(ops.OP_UPDATE_SLIDE, payload)
        self.invalidate_cache()

    def set_slide_title(self, index: int, title: str) -> None:
        """Set the title of a slide.

        Args:
            index: Index of the slide.
            title: New title for the slide.
        """
        self.execute(ops.OP_SET_SLIDE_TITLE, {"slide_index": index, "title": title})
        self.invalidate_cache()

    def merge_from_file(self, path: str) -> None:
        """Merge slides from another presentation file.

        Args:
            path: Path to the presentation file to merge.
        """
        self.execute(ops.OP_MERGE_FROM_FILE, {"path": path})
        self.invalidate_cache()

    def merge_from_editor(self, other: PresentationProtocol) -> None:
        """Merge all slides from *other* into this presentation."""
        self.execute(ops.OP_MERGE_FROM_EDITOR, {"source_handle": other.handle})
        self.invalidate_cache()

    def add_title_slide(self, title: str) -> Slide:
        """Add a title slide to the presentation.

        Args:
            title: Title for the slide.

        Returns:
            The newly created slide.
        """
        return self.add_slide(title, layout="title_only")

    def add_bullet_slide(self, title: str, bullets: list[str]) -> Slide:
        """Add a bullet slide to the presentation.

        Args:
            title: Title for the slide.
            bullets: List of bullet points.

        Returns:
            The newly created slide.
        """
        return self.add_slide(title, layout="title_and_content", bullets=bullets)

    def add_paragraph_slide(
        self,
        title: str,
        paragraph: str,
        *,
        bounds: tuple[float, float, float, float] | None = None,
        layout: str | None = None,
    ) -> Slide:
        """Add a slide with one paragraph textbox using sensible default bounds.

        Args:
            title: Slide title text.
            paragraph: Paragraph body content.
            bounds: Optional textbox bounds in EMU (left, top, width, height).
                Defaults to a readable content region when omitted.
            layout: Optional slide layout name.

        Returns:
            The newly created slide.
        """
        slide = self.add_slide(title, layout=layout)
        slide.add_paragraph(paragraph, bounds=bounds)
        return slide

    def add_slide_from_markdown(self, markdown: str, *, layout: str = "") -> int:
        """Append slides generated from a Markdown string.

        Converts the Markdown document into one or more slides using the Go
        markdown engine and appends them to the presentation.

        Args:
            markdown: Markdown content to convert.
            layout: Optional layout name to apply to all generated slides.

        Returns:
            Zero-based index of the first slide that was added, or -1 when no
            slides were produced (empty markdown).
        """
        payload: dict[str, object] = {"markdown": markdown}
        if layout:
            payload["layout"] = layout
        result = self.execute(ops.OP_MARKDOWN_TO_SLIDES, payload)
        self.invalidate_cache()
        return int(cast("int", result.get("first_index", -1)))

    def add_slide_from_url(self, url: str, *, layout: str = "") -> int:
        """Fetch a web page and append slides generated from its content.

        Uses the Go URL-fetch engine to download the page, parse it, and
        convert it into one or more slides that are appended to the
        presentation.

        Args:
            url: URL of the web page to fetch.
            layout: Optional layout name to apply to all generated slides.

        Returns:
            Zero-based index of the first slide that was added.
        """
        payload: dict[str, object] = {"url": url}
        result = self.execute(ops.OP_URL_FETCH_TO_SLIDES, payload)
        self.invalidate_cache()
        _ = layout  # layout not yet forwarded for URL fetch
        return int(cast("int", result.get("first_index", -1)))

    def __getitem__(self, index: int | slice) -> Slide | list[Slide]:
        """Get a slide by index or slice."""
        return self.slides[index]

    def __len__(self) -> int:
        """Return the number of slides."""
        return self.slide_count

    def __iter__(self) -> Iterator[Slide]:
        """Iterate over all slides."""
        return iter(self.slides)
