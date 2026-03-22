"""Slide shape mixin scoped to shape-domain operations."""

from __future__ import annotations

from typing import TYPE_CHECKING

from ... import ops

if TYPE_CHECKING:
    from ...constants import ConnectorType, ShapeType
    from ...schemas import ImageMetadata, Shape, ShapeProps, ShapeUpdate
    from ..contracts import SlidePresentationProtocol
    from .freeform_builder import FreeformBuilder


class SlideShapeMixin:
    """Mixin providing shape manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

        def _invalidate_text_state_cache_if_present(self) -> None: ...

    def _invalidate_shape_cache_if_present(self) -> None:
        """Invalidate shape cache when slide implementation exposes it."""
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()

    def add_shape(
        self,
        shape_type: ShapeType,
        bounds: tuple[float, float, float, float],
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a shape to the slide and invalidate shape/text caches."""
        shape_id = self._presentation.add_shape(
            self.index, shape_type, bounds, **kwargs
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_textbox(
        self,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a textbox and invalidate shape/text caches."""
        shape_id = self._presentation.add_textbox(
            self.index, left, top, width, height, text=text, **kwargs
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_connector(
        self,
        connector_type: ConnectorType,
        begin_x: float,
        begin_y: float,
        end_x: float,
        end_y: float,
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a connector and invalidate shape/text caches."""
        shape_id = self._presentation.add_connector(
            self.index,
            connector_type,
            begin_x,
            begin_y,
            end_x,
            end_y,
            **kwargs,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_group_shape(self, shapes: list[int] | None = None) -> int:
        """Add a group shape and invalidate shape/text caches."""
        shape_id = self._presentation.add_group_shape(self.index, shapes=shapes)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def build_freeform(
        self,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> FreeformBuilder:
        """Build a freeform path builder anchored to this slide."""
        return self._presentation.build_freeform(
            self.index, start_x=start_x, start_y=start_y, scale=scale
        )

    def add_image(
        self,
        path: str | None,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Add an image and invalidate shape/text caches."""
        shape_id = self._presentation.add_image(self.index, path, bounds, **kwargs)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def get_image_metadata(self, shape_id: int) -> ImageMetadata:
        """Return image metadata for an image shape on this slide."""
        return self._presentation.get_image_metadata(self.index, shape_id)

    def add_video(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        """Add a video and invalidate shape/text caches."""
        shape_id = self._presentation.add_video(
            self.index,
            source,
            bounds,
            name=name,
            poster_frame=poster_frame,
            mime_type=mime_type,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_audio(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        """Add audio media and invalidate shape/text caches."""
        shape_id = self._presentation.add_audio(
            self.index,
            source,
            bounds,
            name=name,
            poster_frame=poster_frame,
            mime_type=mime_type,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_ole_object(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        prog_id: str | None = None,
        icon: str | bytes | None = None,
    ) -> int:
        """Add an OLE object and invalidate shape/text caches."""
        shape_id = self._presentation.add_ole_object(
            self.index, source, bounds, name=name, prog_id=prog_id, icon=icon
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def remove_shape(self, shape_id: int) -> None:
        """Remove a shape and invalidate shape/text caches."""
        self._presentation.remove_shape(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def group_shapes(self, shape_ids: list[int]) -> int:
        """Group existing shapes and return the new group shape id."""
        shape_id = self._presentation.group_shapes(self.index, shape_ids)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def ungroup_shapes(self, shape_id: int) -> int:
        """Ungroup a group shape and return the first child group id."""
        group_id = self._presentation.ungroup_shapes(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return group_id

    def move_shape_to_front(self, shape_id: int) -> None:
        """Move a shape to front and invalidate shape/text caches."""
        self._presentation.move_shape_to_front(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def move_shape_to_back(self, shape_id: int) -> None:
        """Move a shape to back and invalidate shape/text caches."""
        self._presentation.move_shape_to_back(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None:
        """Update a shape and invalidate shape/text caches."""
        self._presentation.update_shape(self.index, shape_id, updates)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def add_mermaid(self, diagram: str, *, theme: str = "") -> tuple[int, int]:
        r"""Render a Mermaid diagram onto this slide using the Go engine.

        The diagram code is parsed and converted into native PowerPoint shapes
        and connectors which are appended to the slide.

        Args:
            diagram: Mermaid diagram source code (e.g. ``"flowchart LR\nA-->B"``).
            theme: Optional Mermaid theme name (reserved for future use).

        Returns:
            A ``(shape_count, connector_count)`` tuple indicating how many
            shapes and connectors were added to the slide.

        Example::

            count, conns = slide.add_mermaid("flowchart LR\nA-->B-->C")
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "diagram": diagram,
        }
        if theme:
            payload["theme"] = theme
        result = self._presentation.execute(ops.OP_ADD_MERMAID_SHAPE, payload)
        self._invalidate_shape_cache_if_present()
        shape_count = result.get("shape_count")
        if not isinstance(shape_count, int):
            raise TypeError("bridge response shape_count must be an int")
        connector_count = result.get("connector_count")
        if not isinstance(connector_count, int):
            raise TypeError("bridge response connector_count must be an int")
        return shape_count, connector_count

    def list_shapes(self) -> list[Shape]:
        """List shapes on the slide."""
        return self._presentation.list_shapes(self.index)
