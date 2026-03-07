"""Presentation slides mixin for gopptx library."""
# ruff: noqa: D102

from __future__ import annotations

import uuid
from typing import TYPE_CHECKING, cast

from .. import ops
from ..slide.slide import Slide
from ..utils import is_four_number_bounds
from .helpers import PresentationMixinBase
from .layout_theme import PresentationLayoutMixin, PresentationThemeMixin
from .master import SlideMasters
from .placeholder_payloads import (
    build_placeholder_chart_payload,
    build_placeholder_table_payload,
)

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ..schemas import (
        CoreProperties,
        Section,
        SlideMetadata,
    )
    from .helpers import PresentationProtocol
    from .presentation import Presentation


class PresentationSectionMixin(PresentationMixinBase):
    """Mixin providing section management methods."""

    @property
    def sections(self) -> list[Section]:
        """Get all sections in the presentation."""
        result = self.execute(ops.OP_GET_SECTIONS, {})
        raw_sections = result.get("sections")
        sections = cast(
            "list[dict[str, object]]",
            raw_sections if isinstance(raw_sections, list) else [],
        )
        for s in sections:
            if "Name" in s and "name" not in s:
                s["name"] = s["Name"]
            if "GUID" in s and "guid" not in s:
                s["guid"] = s["GUID"]
            if "SlideIDs" in s and "slide_ids" not in s:
                s["slide_ids"] = s["SlideIDs"]
            if "id" not in s and "name" in s:
                s["id"] = s["name"]
        return cast("list[Section]", sections)

    def get_sections(self) -> list[Section]:
        """Get all sections in the presentation."""
        return self.sections

    def add_section(self, name: str, slide_indices: list[int]) -> None:
        """Add a section to the presentation.

        Args:
            name: Name of the section.
            slide_indices: List of slide indices to include in the section.
        """
        self.execute(ops.OP_ADD_SECTION, {"name": name, "slide_indices": slide_indices})

    def remove_section(self, name: str) -> None:
        """Remove a section from the presentation.

        Args:
            name: Name of the section to remove.
        """
        self.execute(ops.OP_REMOVE_SECTION, {"name": str(name)})

    def rename_section(self, old_name: str, new_name: str) -> None:
        """Rename a section.

        Args:
            old_name: Current name of the section.
            new_name: New name for the section.
        """
        self.execute(
            ops.OP_RENAME_SECTION,
            {"old_name": str(old_name), "new_name": str(new_name)},
        )


class PresentationPropertiesMixin(PresentationMixinBase):
    """Mixin providing document properties and protection methods."""

    @property
    def core_properties(self) -> CoreProperties:
        """Get the core properties of the presentation."""
        return cast("CoreProperties", self.execute(ops.OP_GET_CORE_PROPERTIES, {}))

    def get_core_properties(self) -> CoreProperties:
        """Get the core properties of the presentation."""
        return self.core_properties

    @core_properties.setter
    def core_properties(self, props: CoreProperties) -> None:
        self.execute(ops.OP_SET_CORE_PROPERTIES, cast("dict[str, object]", props))

    def set_core_properties(self, props: CoreProperties) -> None:
        """Set the core properties of the presentation.

        Args:
            props: Dictionary of core properties to set.
        """
        self.core_properties = props

    @property
    def title(self) -> str:
        """The title of the presentation."""
        return self.core_properties.get("title", "")

    @title.setter
    def title(self, value: str) -> None:
        props = self.core_properties
        props["title"] = value
        self.core_properties = props

    def set_modify_password(self, password: str) -> None:
        """Set the modify password for the presentation.

        Args:
            password: Password to set.
        """
        self.execute(ops.OP_SET_MODIFY_PASSWORD, {"password": password})

    def set_mark_as_final(self, *, final: bool = True) -> None:
        """Mark the presentation as final.

        Args:
            final: True to mark as final, False to unmark.
        """
        self.execute(ops.OP_SET_MARK_AS_FINAL, {"final": final})

    # Additional python-pptx compatible core properties

    @property
    def author(self) -> str:
        """The author/creator of the presentation (python-pptx: author)."""
        return self.core_properties.get("creator", "")

    @author.setter
    def author(self, value: str) -> None:
        props = self.core_properties
        props["creator"] = value
        self.core_properties = props

    @property
    def comments(self) -> str:
        """The comments/description of the presentation (python-pptx: comments)."""
        return self.core_properties.get("description", "")

    @comments.setter
    def comments(self, value: str) -> None:
        props = self.core_properties
        props["description"] = value
        self.core_properties = props

    @property
    def identifier(self) -> str:
        """The identifier of the presentation."""
        return self.core_properties.get("identifier", "")

    @identifier.setter
    def identifier(self, value: str) -> None:
        props = self.core_properties
        props["identifier"] = value
        self.core_properties = props

    @property
    def language(self) -> str:
        """The language of the presentation."""
        return self.core_properties.get("language", "")

    @language.setter
    def language(self, value: str) -> None:
        props = self.core_properties
        props["language"] = value
        self.core_properties = props

    @property
    def last_printed(self) -> str:
        """The last printed date of the presentation."""
        return self.core_properties.get("lastPrinted", "")

    @last_printed.setter
    def last_printed(self, value: str) -> None:
        props = self.core_properties
        props["lastPrinted"] = value
        self.core_properties = props

    @property
    def version(self) -> str:
        """The version of the presentation."""
        return self.core_properties.get("version", "")

    @version.setter
    def version(self, value: str) -> None:
        props = self.core_properties
        props["version"] = value
        self.core_properties = props


class PresentationSlidesMixin(
    PresentationPropertiesMixin,
    PresentationSectionMixin,
    PresentationThemeMixin,
    PresentationLayoutMixin,
):
    """Mixin providing slide-related methods for Presentation."""

    _BOUNDS_COMPONENTS = 4

    if TYPE_CHECKING:

        @property
        def slides(self) -> list[Slide]: ...

        @property
        def slide_count(self) -> int: ...

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
        """Add a new slide to the presentation."""
        payload: dict[str, object] = {"title": title}
        if layout:
            payload["layout"] = layout
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
            payload["layout"] = layout
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

    def __getitem__(self, index: int | slice) -> Slide | list[Slide]:
        """Get a slide by index or slice."""
        return self.slides[index]

    def __len__(self) -> int:
        """Return the number of slides."""
        return self.slide_count

    def __iter__(self) -> Iterator[Slide]:
        """Iterate over all slides."""
        return iter(self.slides)
