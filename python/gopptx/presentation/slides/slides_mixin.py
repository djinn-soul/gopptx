"""Presentation slides mixin for gopptx library."""

from __future__ import annotations

import os
import uuid
from typing import TYPE_CHECKING, cast

from ... import ops
from ...api_errors import GopptxError
from ...slide.slide import Slide
from .layout_theme_mixin import PresentationLayoutMixin, PresentationThemeMixin
from .master import SlideLayout, SlideMasters
from .properties_mixin import PresentationPropertiesMixin
from .sections_mixin import PresentationSectionMixin
from .slide_layout_enum import SlideLayoutType
from .slides_extras_mixin import PresentationSlidesExtrasMixin
from .slides_placeholder_mixin import PresentationPlaceholderMixin

if TYPE_CHECKING:
    from collections.abc import Sequence

    from ...schemas import SlideMetadata
    from ...slide.contracts import SlidePresentationProtocol
    from ..helpers import PresentationProtocol
    from .collection import Slides
    from .master import SlideLayouts, SlideMaster


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

    @property
    def slide_master(self) -> SlideMaster:
        """Return the primary slide master."""
        return self.slide_masters[0]

    @property
    def slide_layouts(self) -> SlideLayouts:
        """Return slide layouts of the primary slide master."""
        return self.slide_masters[0].slide_layouts

    def _resolve_layout_object(self, layout: SlideLayout) -> str:
        """Return the local part path for a SlideLayout, importing it if foreign.

        A layout from another open presentation has no parts in this package, so
        its master family is copied over first and the imported copy's part path
        is what the slide gets bound to (Issue #175).
        """
        if layout.presentation is self:
            return layout.part
        if getattr(self, "_batch_active", False):
            # The import has to finish before the queued add can name the part it
            # produced, and a queued op returns nothing to name it with.
            message = """importing a layout from another presentation is not \
allowed inside a batch; add the slide outside the batch, or use a layout from \
this presentation"""
            raise GopptxError(
                message,
                code="BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED",
            )
        result = self.execute(
            ops.OP_IMPORT_LAYOUT_FROM,
            {"source_handle": layout.presentation.handle, "layout_part": layout.part},
        )
        self.invalidate_cache()
        imported = result.get("layout_part")
        if not isinstance(imported, str) or not imported:
            raise TypeError("bridge response layout_part must be a non-empty string")
        return imported

    def add_slide(
        self,
        title: str = "",
        layout: str | SlideLayout | None = None,
        bullets: list[str] | None = None,
        index: int | None = None,
    ) -> Slide:
        """Add a new slide to the presentation, optionally inserting at index (Issue #194).

        Args:
            title: Title text for the new slide.
            layout: A built-in layout name, or a :class:`SlideLayout` object --
                including one taken from another open presentation, whose master
                family is imported into this deck first (Issue #175).
            bullets: Body bullet lines.
            index: Position to insert at; appended when omitted.
        """
        payload: dict[str, object] = {"title": title}
        if isinstance(layout, SlideLayout):
            # The part goes into the add itself. Adding a default slide and then
            # only retargeting its layout relationship left the default title and
            # body placeholders in place whatever the layout said.
            payload["layout_part"] = self._resolve_layout_object(layout)
        elif layout:
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

        if index is not None and index != slide_index:
            self.move_slide(slide_index, index)
            slide_index = index

        return self.slides[slide_index]

    def insert_slide(
        self,
        index: int,
        layout: str | None = None,
        title: str = "",
        bullets: list[str] | None = None,
    ) -> Slide:
        """Insert a slide at index (Issue #194)."""
        return self.add_slide(title=title, layout=layout, bullets=bullets, index=index)

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

    def set_slide_hidden(self, index: int, *, hidden: bool = True) -> None:
        """Mark or unmark a slide as hidden."""
        self.execute(
            ops.OP_SET_SLIDE_HIDDEN, {"slide_index": index, "hidden": bool(hidden)}
        )
        self.invalidate_cache()

    def merge_from_file(self, path: str) -> None:
        """Merge slides from another presentation file."""
        self.execute(ops.OP_MERGE_FROM_FILE, {"path": path})
        self.invalidate_cache()

    def merge_from_editor(self, other: PresentationProtocol) -> None:
        """Merge all slides from *other* into this presentation."""
        self.execute(ops.OP_MERGE_FROM_EDITOR, {"source_handle": other.handle})
        self.invalidate_cache()

    def copy_slides_from(
        self,
        source: PresentationProtocol | str | os.PathLike[str],
        slide_indices: Sequence[int] | int | None = None,
    ) -> list[int]:
        """Copy selected slides from another presentation (Issue #1036).

        Unlike :meth:`merge_from_editor`, which appends the whole deck, this
        lifts out only the slides named, in the order given.

        Args:
            source: An open presentation, or a path to a ``.pptx`` file.
            slide_indices: Zero-based source slide index, a sequence of them, or
                ``None`` for every slide.

        Returns:
            The indices the copied slides now occupy in this presentation.
        """
        payload: dict[str, object] = {}
        if isinstance(source, (str, os.PathLike)):
            payload["path"] = os.fspath(source)
        else:
            payload["source_handle"] = source.handle
        if slide_indices is not None:
            indices = (
                [slide_indices]
                if isinstance(slide_indices, int)
                else [int(i) for i in slide_indices]
            )
            payload["slide_indices"] = indices

        result = self.execute(ops.OP_COPY_SLIDES_FROM, payload)
        self.invalidate_cache()
        first = int(cast("int", result.get("first_index", 0)))
        count = int(cast("int", result.get("slide_count", 0)))
        return list(range(first, first + count))

    def copy_slide_from(
        self,
        source: PresentationProtocol | str | os.PathLike[str],
        slide_index: int,
    ) -> int:
        """Copy one slide from another presentation (Issue #1036).

        Returns the index the copied slide now occupies in this presentation.
        """
        return self.copy_slides_from(source, [slide_index])[0]
