from __future__ import annotations

from typing import TYPE_CHECKING, Any

from . import ops
from .types import BatchItemResult

if TYPE_CHECKING:
    from .api_presentation import Presentation


class _BatchContext:
    """Context manager for buffering operations and executing them as a batch."""

    _READ_OPS = {
        ops.OP_SLIDE_COUNT,
        ops.OP_GET_METADATA,
        ops.OP_LIST_SLIDES,
        ops.OP_GET_SECTIONS,
        ops.OP_GET_CORE_PROPERTIES,
        ops.OP_SEARCH_SHAPES,
        ops.OP_GET_AUTHORS,
        ops.OP_GET_COMMENTS,
        ops.OP_LIST_SHAPES,
        ops.OP_GET_NOTES,
        ops.OP_LIST_SLIDE_CHARTS,
        ops.OP_LIST_SLIDE_LAYOUTS,
    }

    def __init__(self, presentation: Presentation, stop_on_error: bool = False):
        self._presentation = presentation
        self._stop_on_error = stop_on_error
        self._results: list[BatchItemResult] = []

    def __enter__(self) -> _BatchContext:
        self._presentation._begin_batch(self._stop_on_error)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        if exc_type is None:
            self._results = self._presentation._end_batch()
        else:
            self._presentation._abort_batch()

    @property
    def results(self) -> list[BatchItemResult]:
        """Returns the results of the batch execution."""
        return self._results

    def __getattr__(self, name: str):
        # Forward mutating API calls to Presentation while batch mode is active.
        target = getattr(self._presentation, name)
        if callable(target):
            return target
        raise AttributeError(name)
