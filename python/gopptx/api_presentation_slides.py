"""Presentation slides mixin for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, cast

from . import ops
from .api_master import SlideMasters
from .api_slide import Slide

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .schemas import (
        CoreProperties,
        Section,
        SlideLayoutInfo,
        SlideMasterCloneResult,
    )


class PresentationLayoutMixin:
    """Mixin providing slide layout management methods."""

    def list_slide_layouts(self) -> list[SlideLayoutInfo]:
        """List all available slide layouts."""
        result = self.execute(ops.OP_LIST_SLIDE_LAYOUTS, {})
        layouts = cast("list[dict]", result.get("layouts", []))
        for item in layouts:
            if "Name" in item and "name" not in item:
                item["name"] = item["Name"]
            if "Part" in item and "part" not in item:
                item["part"] = item["Part"]
            if "MasterPart" in item and "master_part" not in item:
                item["master_part"] = item["MasterPart"]
        return cast("list[SlideLayoutInfo]", layouts)

    def rebind_slide_layout(self, slide_index: int, layout_part: str) -> None:
        """Rebind a slide to a different layout.

        Args:
            slide_index: Index of the slide.
            layout_part: Part name of the layout.
        """
        target_layout = layout_part
        if "/" not in target_layout:
            for layout in self.list_slide_layouts():
                if layout.get("name") == target_layout:
                    target_layout = cast("str", layout.get("part", target_layout))
                    break
        self.execute(
            ops.OP_REBIND_SLIDE_LAYOUT,
            {"slide_index": slide_index, "layout_part": target_layout},
        )
        self.invalidate_cache()

    def clone_layout_master_family(self, layout_part: str) -> SlideMasterCloneResult:
        """Clone a layout and its master family.

        Args:
            layout_part: Part name of the layout to clone.

        Returns:
            Result containing the new layout and master parts.
        """
        result = self.execute(
            ops.OP_CLONE_LAYOUT_MASTER_FAMILY, {"layout_part": layout_part}
        )
        self.invalidate_cache()
        return cast("SlideMasterCloneResult", result)


class PresentationThemeMixin:
    """Mixin providing theme and slide size configuration methods."""

    def apply_theme(self, theme_name: str) -> None:
        """Apply a theme to the presentation.

        Args:
            theme_name: Name of the theme to apply.
        """
        theme = theme_name
        if theme_name.lower() == "office":
            theme = "Corporate"
        self.execute(ops.OP_APPLY_THEME, {"theme_name": theme})
        self.invalidate_cache()

    def set_slide_size(self, width: int, height: int) -> None:
        """Set the slide size.

        Args:
            width: Width in EMUs.
            height: Height in EMUs.
        """
        self.execute(ops.OP_SET_SLIDE_SIZE, {"width": width, "height": height})
        self.invalidate_cache()


class PresentationSectionMixin:
    """Mixin providing section management methods."""

    @property
    def sections(self) -> list[Section]:
        """Get all sections in the presentation."""
        result = self.execute(ops.OP_GET_SECTIONS, {})
        raw_sections = result.get("sections")
        sections = cast(
            "list[dict]", raw_sections if isinstance(raw_sections, list) else []
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


class PresentationPropertiesMixin:
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
        self.execute(ops.OP_SET_CORE_PROPERTIES, props)

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


class PresentationSlidesMixin(
    PresentationPropertiesMixin,
    PresentationSectionMixin,
    PresentationThemeMixin,
    PresentationLayoutMixin,
):
    """Mixin providing slide-related methods for Presentation."""

    @property
    def slide_masters(self) -> SlideMasters:
        """Get the slide masters collection."""
        if getattr(self, "_slide_masters_obj", None) is None:
            self._slide_masters_obj = SlideMasters(self)
        return self._slide_masters_obj

    def add_slide(
        self, title: str, layout: str | None = None, bullets: list[str] | None = None
    ) -> Slide:
        """Add a new slide to the presentation."""
        payload: dict[str, Any] = {"title": title}
        if layout:
            payload["layout"] = layout
        if bullets:
            payload["bullets"] = bullets
        result = self.execute(ops.OP_ADD_SLIDE, payload)
        if result.get("_batched", False):
            placeholder = {
                "Index": -1,
                "SlideID": -1,
                "RelationshipID": "",
                "PartName": "",
                "Title": title,
                "index": -1,
                "slide_id": -1,
                "relationship_id": "",
                "part_name": "",
                "title": title,
            }
            return Slide(self, cast("Any", placeholder))
        self.invalidate_cache()
        return self.slides[int(result.get("index", -1))]

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
        return int(result.get("new_index", -1))

    def update_slide(
        self,
        index: int,
        title: str | None = None,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None:
        """Update slide properties."""
        payload: dict[str, Any] = {"slide_index": index}
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

    def __getitem__(self, index: int) -> Slide:
        """Get a slide by index."""
        return self.slides[index]

    def __len__(self) -> int:
        """Return the number of slides."""
        return self.slide_count

    def __iter__(self) -> Iterator[Slide]:
        """Iterate over all slides."""
        return iter(self.slides)
