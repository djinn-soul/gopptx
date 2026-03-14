"""Slide master and layout classes for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ..schemas import PlaceholderInfo, SlideLayoutInfo
    from .base import PresentationProtocol


class SlideLayout:
    """Represents a slide layout within a slide master."""

    def __init__(self, master: SlideMaster, info: SlideLayoutInfo) -> None:
        """Initialize the slide layout.

        Args:
            master: The parent slide master.
            info: The layout information dictionary.
        """
        super().__init__()
        self._master = master
        self._info = info
        self._part = str(info.get("part", info.get("Part", "")))
        self._name = str(info.get("name", info.get("Name", "")))

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

    @property
    def shapes(self) -> list[str]:
        """Layout shape names snapshot."""
        shapes = self._info.get("shapes", self._info.get("Shapes", []))
        if isinstance(shapes, list):
            return [str(item) for item in cast("list[object]", shapes)]
        return []

    @property
    def placeholders(self) -> list[PlaceholderInfo]:
        """Layout placeholder records snapshot."""
        placeholders = self._info.get(
            "placeholders", self._info.get("Placeholders", [])
        )
        if isinstance(placeholders, list):
            return cast("list[PlaceholderInfo]", placeholders)
        return []


class SlideLayouts:
    """Collection of slide layouts within a slide master."""

    def __init__(self, master: SlideMaster, layouts: list[SlideLayoutInfo]) -> None:
        """Initialize the slide layouts collection.

        Args:
            master: The parent slide master.
            layouts: List of layout information dictionaries.
        """
        super().__init__()
        self._master = master
        self._layouts = [SlideLayout(master, info) for info in layouts]

    def __len__(self) -> int:
        """Return the number of layouts."""
        return len(self._layouts)

    def __getitem__(self, idx: int) -> SlideLayout:
        """Get a layout by index."""
        return self._layouts[idx]

    def __iter__(self) -> Iterator[SlideLayout]:
        """Iterate over all layouts."""
        return iter(self._layouts)

    def get_by_name(self, name: str) -> SlideLayout | None:
        """Get a layout by name.

        Args:
            name: The name of the layout to find.

        Returns:
            The layout with the given name, or None if not found.
        """
        for layout in self._layouts:
            if layout.name == name:
                return layout
        return None


class SlideMaster:
    """Represents a slide master in the presentation."""

    def __init__(
        self,
        prs: PresentationProtocol,
        part: str,
        info: dict[str, object] | None = None,
    ) -> None:
        """Initialize the slide master.

        Args:
            prs: The presentation base instance.
            part: The part path of this slide master.
            info: Optional master metadata payload from the bridge.
        """
        super().__init__()
        self._prs = prs
        self._part = part
        self._info = info or {}
        self._slide_layouts: SlideLayouts | None = None

    @property
    def part(self) -> str:
        """The part path of this slide master."""
        return self._part

    @property
    def slide_layouts(self) -> SlideLayouts:
        """Get the slide layouts for this master."""
        if self._slide_layouts is None:
            # Rebind OP_LIST_MASTER_LAYOUTS to the presentation base execute
            result = self._prs.execute(
                ops.OP_LIST_MASTER_LAYOUTS, {"master_part": self._part}
            )
            layouts = cast("list[dict[str, object]]", result.get("layouts", []))
            self._slide_layouts = SlideLayouts(
                self, cast("list[SlideLayoutInfo]", layouts)
            )
        return self._slide_layouts

    @property
    def shapes(self) -> list[str]:
        """Master shape names snapshot."""
        shapes = self._info.get("shapes", self._info.get("Shapes", []))
        if isinstance(shapes, list):
            return [str(item) for item in cast("list[object]", shapes)]
        return []

    @property
    def placeholders(self) -> list[PlaceholderInfo]:
        """Master placeholder records snapshot."""
        placeholders = self._info.get(
            "placeholders", self._info.get("Placeholders", [])
        )
        if isinstance(placeholders, list):
            return cast("list[PlaceholderInfo]", placeholders)
        return []


class SlideMasters:
    """Collection of slide masters in the presentation."""

    def __init__(self, prs: PresentationProtocol) -> None:
        """Initialize the slide masters collection.

        Args:
            prs: The presentation base instance.
        """
        super().__init__()
        self._prs = prs
        self._masters: list[SlideMaster] | None = None

    def _load(self) -> None:
        """Load the slide masters from the presentation."""
        if self._masters is not None:
            return
        result = self._prs.execute(ops.OP_LIST_SLIDE_MASTERS, {})
        master_infos = cast("list[dict[str, object]]", result.get("masters", []))
        self._masters = []
        for info in master_infos:
            part = str(info.get("part", info.get("Part", "")))
            self._masters.append(SlideMaster(self._prs, part, info))

    def __len__(self) -> int:
        """Return the number of slide masters."""
        self._load()
        if self._masters is None:  # pragma: no cover
            msg = "Masters not loaded"
            raise RuntimeError(msg)
        return len(self._masters)

    def __getitem__(self, idx: int) -> SlideMaster:
        """Get a slide master by index."""
        self._load()
        if self._masters is None:  # pragma: no cover
            msg = "Masters not loaded"
            raise RuntimeError(msg)
        return self._masters[idx]

    def __iter__(self) -> Iterator[SlideMaster]:
        """Iterate over all slide masters."""
        self._load()
        if self._masters is None:  # pragma: no cover
            msg = "Masters not loaded"
            raise RuntimeError(msg)
        return iter(self._masters)
