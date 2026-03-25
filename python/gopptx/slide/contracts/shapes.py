"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from collections.abc import Mapping, Sequence

    from ...constants import ConnectorType, ShapeType
    from ...schemas import Shape, ShapeUpdate
    from ..shapes.freeform_builder import FreeformBuilder


class ShapeOperationsProtocol(Protocol):
    """Shape management operations."""

    def list_shapes(self, slide_index: int) -> list[Shape]:
        """Protocol member."""
        ...

    def add_shape(
        self,
        slide_index: int,
        shape_type: ShapeType,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def add_connector(
        self,
        slide_index: int,
        connector_type: ConnectorType,
        begin_x: float,
        begin_y: float,
        end_x: float,
        end_y: float,
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def add_connectors(
        self, slide_index: int, connectors: Sequence[Mapping[str, object]]
    ) -> list[int]:
        """Protocol member."""
        ...

    def add_group_shape(self, slide_index: int, shapes: list[int] | None = None) -> int:
        """Protocol member."""
        ...

    def build_freeform(
        self,
        slide_index: int,
        *,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> FreeformBuilder:
        """Protocol member."""
        ...

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        """Protocol member."""
        ...

    def clear_shapes(self, slide_index: int) -> int:
        """Protocol member."""
        ...

    def group_shapes(self, slide_index: int, shape_ids: list[int]) -> int:
        """Protocol member."""
        ...

    def ungroup_shapes(self, slide_index: int, shape_id: int) -> int:
        """Protocol member."""
        ...

    def move_shape_to_front(self, slide_index: int, shape_id: int) -> None:
        """Protocol member."""
        ...

    def move_shape_to_back(self, slide_index: int, shape_id: int) -> None:
        """Protocol member."""
        ...

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None:
        """Protocol member."""
        ...

    def list_placeholders(self, slide_index: int) -> list[dict[str, object]]:
        """Protocol member."""
        ...

    def set_placeholder_content(
        self,
        slide_index: int,
        ph_index: int,
        ph_type: str = "",
        **kwargs: object,
    ) -> None:
        """Protocol member."""
        ...
