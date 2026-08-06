"""Slide mixin for reading SmartArt diagrams and editing individual nodes.

A node is addressed by its path: the 0-based index of each step down from the
top level, so ``[1, 0]`` is the first child of the second entry.
"""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops

if TYPE_CHECKING:
    from ..contracts import SlidePresentationProtocol


class SlideSmartArtNodeMixin:
    """Mixin adding per-node SmartArt edits and reads to Slide."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

        def _invalidate_shape_and_text_caches_if_present(self) -> None: ...

    def add_smartart_node(
        self,
        shape_id: int,
        text: str,
        *,
        parent_path: list[int] | None = None,
        index: int | None = None,
        color: str | None = None,
        image: str | None = None,
    ) -> None:
        """Insert one node into an existing SmartArt diagram.

        Args:
            shape_id: The shape ID of the SmartArt graphic frame.
            text: Caption for the new node.
            parent_path: Path of the node to add under, as the 0-based index of
                each step down from the top level -- ``[1, 0]`` is the first
                child of the second entry.  Omit for a top-level node.
            index: Position among its siblings.  Omit to append.
            color: RGB hex fill for this node alone, e.g. ``"C00000"``.
            image: Path to a picture, for layouts that draw a placeholder.
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "text": text,
        }
        if parent_path is not None:
            payload["parent_path"] = parent_path
        if index is not None:
            payload["index"] = index
        if color is not None:
            payload["color"] = color
        if image is not None:
            payload["image"] = image
        self._presentation.execute(ops.OP_ADD_SMART_ART_NODE, payload)
        self._invalidate_shape_and_text_caches_if_present()

    def remove_smartart_node(self, shape_id: int, path: list[int]) -> None:
        """Delete the SmartArt node at ``path``, along with its children."""
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "path": path,
        }
        self._presentation.execute(ops.OP_REMOVE_SMART_ART_NODE, payload)
        self._invalidate_shape_and_text_caches_if_present()

    def update_smartart_node(
        self,
        shape_id: int,
        path: list[int],
        *,
        text: str | None = None,
        color: str | None = None,
        image: str | None = None,
    ) -> None:
        """Change the text, color or picture of one SmartArt node.

        Omitted arguments are left as they are.
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "path": path,
        }
        if text is not None:
            payload["text"] = text
        if color is not None:
            payload["color"] = color
        if image is not None:
            payload["image"] = image
        self._presentation.execute(ops.OP_UPDATE_SMART_ART_NODE, payload)
        self._invalidate_shape_and_text_caches_if_present()

    def get_smartart(self, shape_id: int) -> dict[str, object]:
        """Read back one SmartArt diagram.

        Returns:
            ``{"shape_id", "layout", "quick_style", "color_style", "nodes"}``,
            where ``nodes`` is the nested tree ``add_smartart`` and
            ``set_smartart_nodes`` accept.  The tree comes from the data model
            PowerPoint draws from, so it reflects what the slide shows.
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
        }
        return self._presentation.execute(ops.OP_GET_SMART_ART, payload)

    def list_smartart(self) -> list[dict[str, object]]:
        """Read back every SmartArt diagram on the slide.

        Returns:
            A list of the same dicts ``get_smartart`` returns, one per diagram.
        """
        payload: dict[str, object] = {"slide_index": self.index}
        result = self._presentation.execute(ops.OP_LIST_SMART_ART, payload)
        diagrams = result.get("diagrams", [])
        if not isinstance(diagrams, list):
            return []
        return [
            cast("dict[str, object]", entry)
            for entry in cast("list[object]", diagrams)
            if isinstance(entry, dict)
        ]
