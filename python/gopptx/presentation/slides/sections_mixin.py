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
        # execute() supplies the name/guid/slide_ids aliases; only `id`, which
        # is not a bridge field at all, has to be derived here.
        for section in sections:
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
