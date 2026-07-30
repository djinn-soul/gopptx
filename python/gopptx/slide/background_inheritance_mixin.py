"""Slide background-inheritance control."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .. import ops

if TYPE_CHECKING:
    from .contracts import SlidePresentationProtocol


class SlideBackgroundInheritanceMixin:
    """Expose python-pptx's follow_master_background property."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Current slide index."""
            ...

    @property
    def follow_master_background(self) -> bool:
        """Whether this slide inherits its slide-master background."""
        result = self._presentation.execute(
            ops.OP_GET_SLIDE_BACKGROUND,
            {"slide_index": self.index},
        )
        return not bool(result.get("background_xml", ""))

    @follow_master_background.setter
    def follow_master_background(self, follow: bool) -> None:
        self._presentation.execute(
            ops.OP_SET_SLIDE_FOLLOW_MASTER_BACKGROUND,
            {"slide_index": self.index, "follow": bool(follow)},
        )
