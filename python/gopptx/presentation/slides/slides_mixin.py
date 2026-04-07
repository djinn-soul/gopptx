"""Presentation slides mixin for gopptx library."""

from __future__ import annotations

import uuid
from typing import TYPE_CHECKING, cast

from ... import ops
from ...slide.slide import Slide
from .layout_theme_mixin import PresentationLayoutMixin, PresentationThemeMixin
from .master import SlideMasters
from .properties_mixin import PresentationPropertiesMixin
from .sections_mixin import PresentationSectionMixin
from .slide_layout_enum import SlideLayoutType
from .slides_extras_mixin import PresentationSlidesExtrasMixin
from .slides_placeholder_mixin import PresentationPlaceholderMixin

if TYPE_CHECKING:
    from ...schemas import SlideMetadata
    from ...slide.contracts import SlidePresentationProtocol
    from ..helpers import PresentationProtocol
    from .collection import Slides


class PresentationSlidesMixin(
    PresentationPlaceholderMixin,
    PresentationSlidesExtrasMixin,
    PresentationPropertiesMixin,
    PresentationSectionMixin,
    PresentationThemeMixin,
    PresentationLayoutMixin,
):
    """Mixin providing slide-related methods for Presentation."""

    _BOUNDS_COMPONENTS = 4
    _EMU_PER_INCH = 914400

    if TYPE_CHECKING:
        header_footer_defaults: dict[str, object]

        @property
        def slides(self) -> Slides:
            """Return materialized slide objects for the current presentation."""
            ...

        @property
        def slide_count(self) -> int:
            """Return the total number of slides in the presentation."""
            ...

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
        """Add a new slide to the presentation."""
        payload: dict[str, object] = {"title": title}
        if layout:
            try:
                validated_layout = SlideLayoutType.validate(layout)
                payload["layout"] = validated_layout
            except ValueError as e:
                raise ValueError(f"Invalid slide layout: {e}") from e
        if bullets:
            payload["bullets"] = bullets
        result = self.execute(ops.OP_ADD_SLIDE, payload)
        if result.get("_batched", False):
            placeholder_metadata = {
                "Title": title,
                "SlideID": -1,
                "RelationshipID": str(uuid.uuid4()),
                "PartName": "/ppt/slides/slide_placeholder.xml",
                "Index": -1,
            }
            return Slide(
                cast("SlidePresentationProtocol", self),
                cast("SlideMetadata", placeholder_metadata),
            )
        self.invalidate_cache()
        slide_index = int(cast("int", result.get("index", -1)))

        if hasattr(self, "header_footer_defaults"):
            defaults = self.header_footer_defaults
            if any([
                defaults.get("show_footer"),
                defaults.get("show_slide_num"),
                defaults.get("show_date_time"),
            ]):
                hf_payload: dict[str, object] = {
                    "slide_index": slide_index,
                    "footer": cast("str", defaults.get("footer", "")),
                    "show_footer": cast("bool", defaults.get("show_footer", False)),
                    "show_slide_num": cast(
                        "bool", defaults.get("show_slide_num", False)
                    ),
                    "show_date_time": cast(
                        "bool", defaults.get("show_date_time", False)
                    ),
                    "date_time_text": cast("str", defaults.get("date_time_text", "")),
                }
                self.execute(ops.OP_SET_SLIDE_HEADER_FOOTER, hf_payload)

        return self.slides[slide_index]

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
        """Duplicate slide at *index* and insert it immediately after the original."""
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
        """Update slide properties."""
        payload: dict[str, object] = {"slide_index": index}
        if title is not None:
            payload["title"] = title
        if layout is not None:
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
        """Set the title of a slide."""
        self.execute(ops.OP_SET_SLIDE_TITLE, {"slide_index": index, "title": title})
        self.invalidate_cache()

    def set_slide_hidden(self, index: int, *, hidden: bool) -> None:
        """Mark or unmark a slide as hidden."""
        self.execute(ops.OP_SET_SLIDE_HIDDEN, {"slide_index": index, "hidden": hidden})
        self.invalidate_cache()

    def merge_from_file(self, path: str) -> None:
        """Merge slides from another presentation file."""
        self.execute(ops.OP_MERGE_FROM_FILE, {"path": path})
        self.invalidate_cache()

    def merge_from_editor(self, other: PresentationProtocol) -> None:
        """Merge all slides from *other* into this presentation."""
        self.execute(ops.OP_MERGE_FROM_EDITOR, {"source_handle": other.handle})
        self.invalidate_cache()
