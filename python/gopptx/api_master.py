"""Slide master and layout classes for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from . import ops

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .api_presentation_base import PresentationBase
    from .types import SlideLayoutInfo


class SlideLayout:
    """Represents a slide layout within a slide master."""

    def __init__(self, master: SlideMaster, info: SlideLayoutInfo) -> None:
        """Initialize the slide layout."""
        self._master = master
        self._part = info.get("part", info.get("Part", ""))
        self._name = info.get("name", info.get("Name", ""))

    @property
    def part(self) -> str:
        """The part path of this layout."""
        return self._part

    @property
    def name(self) -> str:
        """The name of this layout."""
        return self._name

    @property
    def slide_master(self) -> SlideMaster:
        """The parent slide master."""
        return self._master


class SlideLayouts:
    def __init__(self, master: SlideMaster, layouts: list[SlideLayoutInfo]) -> None:
        self._master = master
        self._layouts = [SlideLayout(master, info) for info in layouts]

    def __len__(self) -> int:
        return len(self._layouts)

    def __getitem__(self, idx: int) -> SlideLayout:
        return self._layouts[idx]

    def __iter__(self) -> Iterator[SlideLayout]:
        return iter(self._layouts)

    def get_by_name(self, name: str) -> SlideLayout | None:
        for layout in self._layouts:
            if layout.name == name:
                return layout
        return None


class SlideMaster:
    def __init__(self, prs: PresentationBase, part: str) -> None:
        self._prs = prs
        self._part = part
        self._slide_layouts: SlideLayouts | None = None

    @property
    def part(self) -> str:
        return self._part

    @property
    def slide_layouts(self) -> SlideLayouts:
        if self._slide_layouts is None:
            # Rebind OP_LIST_MASTER_LAYOUTS to the presentation base execute
            result = self._prs.execute(
                ops.OP_LIST_MASTER_LAYOUTS, {"master_part": self._part}
            )
            layouts = cast("list[dict]", result.get("layouts", []))
            self._slide_layouts = SlideLayouts(
                self, cast("list[SlideLayoutInfo]", layouts)
            )
        return self._slide_layouts


class SlideMasters:
    def __init__(self, prs: PresentationBase) -> None:
        self._prs = prs
        self._masters: list[SlideMaster] | None = None

    def _load(self) -> None:
        if self._masters is not None:
            return
        result = self._prs.execute(ops.OP_LIST_SLIDE_MASTERS, {})
        master_infos = cast("list[dict]", result.get("masters", []))
        self._masters = []
        for info in master_infos:
            part = info.get("part", info.get("Part", ""))
            self._masters.append(SlideMaster(self._prs, part))

    def __len__(self) -> int:
        self._load()
        assert self._masters is not None
        return len(self._masters)

    def __getitem__(self, idx: int) -> SlideMaster:
        self._load()
        assert self._masters is not None
        return self._masters[idx]

    def __iter__(self) -> Iterator[SlideMaster]:
        self._load()
        assert self._masters is not None
        return iter(self._masters)
