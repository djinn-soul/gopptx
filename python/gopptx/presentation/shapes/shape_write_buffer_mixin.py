"""Buffered textbox insertion for slide-level shape creation."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from typing_extensions import override

from ... import ops
from ..helpers import PresentationMixinBase
from ..runtime import PresentationRuntimeMixin
from ..text.text_write_buffer_mixin import PresentationTextWriteBufferMixin

if TYPE_CHECKING:
    from collections.abc import Mapping


class PresentationShapeWriteBufferMixin(PresentationMixinBase):
    """Coalesce simple slide textbox inserts into bulk bridge writes."""

    _RESERVED_SHAPE_ID_BLOCK = 64
    _FLUSH_ALL_OPS = frozenset({
        ops.OP_ADD_SLIDE,
        ops.OP_BATCH_EXECUTE,
        ops.OP_DUPLICATE_SLIDE,
        ops.OP_MERGE_FROM_FILE,
        ops.OP_MOVE_SLIDE,
        ops.OP_REMOVE_SLIDE,
    })

    def __init__(self) -> None:
        """Initialize buffered textbox state."""
        super().__init__()
        self._pending_textboxes: dict[int, list[dict[str, object]]] = {}
        self._reserved_shape_ids: dict[int, list[int]] = {}

    def queue_textbox_add(
        self,
        slide_index: int,
        textbox: Mapping[str, object],
    ) -> int:
        """Queue a simple textbox insert and return its real reserved shape ID."""
        shape_id = self._next_reserved_shape_id(slide_index)
        payload = dict(textbox)
        payload["shape_id"] = shape_id
        self._pending_textboxes.setdefault(slide_index, []).append(payload)
        return shape_id

    def has_pending_textbox_adds(self, slide_index: int) -> bool:
        """Return whether the slide has queued textbox inserts."""
        return bool(self._pending_textboxes.get(slide_index))

    def flush_pending_textbox_adds(self, slide_index: int) -> list[int]:
        """Flush queued textboxes for one slide using the bulk bridge op."""
        pending = self._pending_textboxes.get(slide_index)
        if not pending:
            return []
        runtime_self = cast("PresentationRuntimeMixin", self)
        result = PresentationRuntimeMixin.execute(
            runtime_self,
            ops.OP_ADD_TEXTBOXES,
            {
                "slide_index": slide_index,
                "textboxes": [dict(textbox) for textbox in pending],
            },
        )
        self._pending_textboxes.pop(slide_index, None)
        self._reserved_shape_ids.pop(slide_index, None)
        return cast("list[int]", result.get("shape_ids", []))

    def flush_all_pending_textbox_adds(self) -> None:
        """Flush queued textboxes for every slide."""
        for slide_index in sorted(self._pending_textboxes):
            self.flush_pending_textbox_adds(slide_index)

    @override
    def execute(
        self, op: str, payload: dict[str, object] | None = None
    ) -> dict[str, object]:
        """Flush queued textbox inserts before incompatible bridge operations."""
        normalized_payload = payload or {}
        if op in self._FLUSH_ALL_OPS:
            self.flush_all_pending_textbox_adds()
        elif op != ops.OP_RESERVE_SHAPE_IDS:
            slide_index = normalized_payload.get("slide_index")
            if isinstance(slide_index, int) and self.has_pending_textbox_adds(
                slide_index
            ):
                self.flush_pending_textbox_adds(slide_index)
        runtime_self = cast("PresentationRuntimeMixin", self)
        return PresentationRuntimeMixin.execute(runtime_self, op, normalized_payload)

    def open(self, path: str) -> None:
        """Discard pending textbox state before opening another deck."""
        self._reset_shape_write_buffer()
        runtime_self = cast("PresentationRuntimeMixin", self)
        PresentationRuntimeMixin.open(runtime_self, path)

    def save(self, path: str) -> None:
        """Flush queued textboxes before the regular save pipeline runs."""
        self.flush_all_pending_textbox_adds()
        text_self = cast("PresentationTextWriteBufferMixin", self)
        PresentationTextWriteBufferMixin.save(text_self, path)

    def close(self) -> None:
        """Discard queued textbox state when closing the deck."""
        self._reset_shape_write_buffer()
        runtime_self = cast("PresentationRuntimeMixin", self)
        PresentationRuntimeMixin.close(runtime_self)

    def _next_reserved_shape_id(self, slide_index: int) -> int:
        reserved_ids = self._reserved_shape_ids.get(slide_index)
        if not reserved_ids:
            if self.has_pending_textbox_adds(slide_index):
                self.flush_pending_textbox_adds(slide_index)
            reserved_ids = self._reserve_shape_ids(
                slide_index,
                self._RESERVED_SHAPE_ID_BLOCK,
            )
            self._reserved_shape_ids[slide_index] = reserved_ids
        return reserved_ids.pop(0)

    def _reserve_shape_ids(self, slide_index: int, count: int) -> list[int]:
        runtime_self = cast("PresentationRuntimeMixin", self)
        result = PresentationRuntimeMixin.execute(
            runtime_self,
            ops.OP_RESERVE_SHAPE_IDS,
            {"slide_index": slide_index, "count": count},
        )
        return cast("list[int]", result.get("shape_ids", []))

    def _reset_shape_write_buffer(self) -> None:
        self._pending_textboxes = {}
        self._reserved_shape_ids = {}
