"""Presentation section-management mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import Section


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
        for section in sections:
            if "Name" in section and "name" not in section:
                section["name"] = section["Name"]
            if "GUID" in section and "guid" not in section:
                section["guid"] = section["GUID"]
            if "SlideIDs" in section and "slide_ids" not in section:
                section["slide_ids"] = section["SlideIDs"]
            if "id" not in section and "name" in section:
                section["id"] = section["name"]
        return cast("list[Section]", sections)

    def get_sections(self) -> list[Section]:
        """Get all sections in the presentation."""
        return self.sections

    def add_section(self, name: str, slide_indices: list[int]) -> None:
        """Add a section to the presentation."""
        self.execute(ops.OP_ADD_SECTION, {"name": name, "slide_indices": slide_indices})

    def remove_section(self, name: str) -> None:
        """Remove a section from the presentation."""
        self.execute(ops.OP_REMOVE_SECTION, {"name": str(name)})

    def rename_section(self, old_name: str, new_name: str) -> None:
        """Rename a section."""
        self.execute(
            ops.OP_RENAME_SECTION,
            {"old_name": str(old_name), "new_name": str(new_name)},
        )
