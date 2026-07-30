"""Slide-level control of shapes inherited from the slide master."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .. import ops

if TYPE_CHECKING:
    from .contracts import SlidePresentationProtocol


class SlideMasterShapesMixin:
    """Expose PowerPoint's Hide Background Graphics slide setting."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Current slide index."""
            ...

    @property
    def show_master_shapes(self) -> bool:
        """Whether shapes inherited from the slide master are visible."""
        result = self._presentation.execute(
            ops.OP_GET_SLIDE_SHOW_MASTER_SHAPES,
            {"slide_index": self.index},
        )
        return bool(result["visible"])

    @show_master_shapes.setter
    def show_master_shapes(self, visible: bool) -> None:
        self._presentation.execute(
            ops.OP_SET_SLIDE_SHOW_MASTER_SHAPES,
            {"slide_index": self.index, "visible": bool(visible)},
        )
