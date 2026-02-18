from __future__ import annotations

from typing import Any, Dict, Optional, cast

from . import ops
from .api_slide import Slide
from .types import CoreProperties, Section, SlideLayoutInfo, SlideMasterCloneResult


class PresentationSlidesMixin:
    @property
    def sections(self) -> list[Section]:
        result = self.execute(ops.OP_GET_SECTIONS, {})
        return cast(list[Section], result.get("sections", []))

    def add_slide(self, title: str, layout: Optional[str] = None, bullets: Optional[list[str]] = None) -> Slide:
        payload: Dict[str, Any] = {"title": title}
        if layout:
            payload["layout"] = layout
        if bullets:
            payload["bullets"] = bullets
        result = self.execute(ops.OP_ADD_SLIDE, payload)
        self.invalidate_cache()
        return self.slides[int(result.get("index", -1))]

    def remove_slide(self, index: int) -> None:
        self.execute(ops.OP_REMOVE_SLIDE, {"index": index})
        self.invalidate_cache()

    def move_slide(self, from_index: int, to_index: int) -> None:
        self.execute(ops.OP_MOVE_SLIDE, {"from": from_index, "to": to_index})
        self.invalidate_cache()

    def duplicate_slide(self, index: int, insert_at: Optional[int] = None) -> int:
        if insert_at is None:
            insert_at = index + 1
        result = self.execute(ops.OP_DUPLICATE_SLIDE, {"index": index, "insert_at": insert_at})
        self.invalidate_cache()
        return int(result.get("new_index", -1))

    def update_slide(self, index: int, title: Optional[str] = None, layout: Optional[str] = None, bullets: Optional[list[str]] = None) -> None:
        payload: Dict[str, Any] = {"slide_index": index}
        if title is not None:
            payload["title"] = title
        if layout is not None:
            payload["layout"] = layout
        if bullets is not None:
            payload["bullets"] = bullets
        self.execute(ops.OP_UPDATE_SLIDE, payload)
        self.invalidate_cache()

    def set_slide_title(self, index: int, title: str) -> None:
        self.execute(ops.OP_SET_SLIDE_TITLE, {"slide_index": index, "title": title})
        self.invalidate_cache()

    def merge_from_file(self, path: str) -> None:
        self.execute(ops.OP_MERGE_FROM_FILE, {"path": path})
        self.invalidate_cache()

    def add_section(self, name: str, slide_indices: list[int]) -> None:
        self.execute(ops.OP_ADD_SECTION, {"name": name, "slide_indices": slide_indices})

    def remove_section(self, name: str) -> None:
        self.execute(ops.OP_REMOVE_SECTION, {"name": name})

    def rename_section(self, old_name: str, new_name: str) -> None:
        self.execute(ops.OP_RENAME_SECTION, {"old_name": old_name, "new_name": new_name})

    @property
    def core_properties(self) -> CoreProperties:
        return cast(CoreProperties, self.execute(ops.OP_GET_CORE_PROPERTIES, {}))

    @core_properties.setter
    def core_properties(self, props: CoreProperties) -> None:
        self.execute(ops.OP_SET_CORE_PROPERTIES, props)

    @property
    def title(self) -> str:
        return self.core_properties.get("title", "")

    @title.setter
    def title(self, value: str) -> None:
        props = self.core_properties
        props["title"] = value
        self.core_properties = props

    def add_title_slide(self, title: str) -> Slide:
        return self.add_slide(title, layout="title_only")

    def add_bullet_slide(self, title: str, bullets: list[str]) -> Slide:
        return self.add_slide(title, layout="title_and_content", bullets=bullets)

    def __getitem__(self, index: int) -> Slide:
        return self.slides[index]

    def __len__(self) -> int:
        return self.slide_count

    def __iter__(self):
        return iter(self.slides)

    def apply_theme(self, theme_name: str) -> None:
        self.execute(ops.OP_APPLY_THEME, {"theme_name": theme_name})
        self.invalidate_cache()

    def set_slide_size(self, width: int, height: int) -> None:
        self.execute(ops.OP_SET_SLIDE_SIZE, {"width": width, "height": height})
        self.invalidate_cache()

    def list_slide_layouts(self) -> list[SlideLayoutInfo]:
        result = self.execute(ops.OP_LIST_SLIDE_LAYOUTS, {})
        return cast(list[SlideLayoutInfo], result.get("layouts", []))

    def rebind_slide_layout(self, slide_index: int, layout_part: str) -> None:
        self.execute(ops.OP_REBIND_SLIDE_LAYOUT, {"slide_index": slide_index, "layout_part": layout_part})
        self.invalidate_cache()

    def clone_layout_master_family(self, layout_part: str) -> SlideMasterCloneResult:
        result = self.execute(ops.OP_CLONE_LAYOUT_MASTER_FAMILY, {"layout_part": layout_part})
        self.invalidate_cache()
        return cast(SlideMasterCloneResult, result)

