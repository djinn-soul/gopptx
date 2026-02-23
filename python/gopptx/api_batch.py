"""Batch execution context for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, ClassVar

from typing_extensions import Self

from . import ops

if TYPE_CHECKING:
    from .api_presentation import Presentation
    from .schemas import BatchItemResult


class BatchContext:
    """Context manager for buffering operations and executing them as a batch."""

    READ_OPS: ClassVar[set[str]] = {
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
        ops.OP_GET_TABLE,
    }

    def __init__(
        self, presentation: Presentation, *, stop_on_error: bool = False
    ) -> None:
        """Initialize the batch context.

        Args:
            presentation: The presentation to batch operations for.
            stop_on_error: If True, stop batch execution on first error.
        """
        self._presentation = presentation
        self._stop_on_error = stop_on_error
        self._results: list[BatchItemResult] = []

    def __enter__(self) -> Self:
        """Enter the batch context and begin buffering operations."""
        self._presentation.begin_batch(stop_on_error=self._stop_on_error)
        return self

    def __exit__(
        self,
        exc_type: type[BaseException] | None,
        exc_val: BaseException | None,
        exc_tb: object,
    ) -> None:
        """Exit the batch context and execute buffered operations."""
        if exc_type is None:
            self._results = self._presentation.end_batch()
        else:
            self._presentation.abort_batch()

    @property
    def results(self) -> list[BatchItemResult]:
        """Returns the results of the batch execution."""
        return self._results

    def __getattr__(self, name: str) -> object:
        """Forward attribute access to the presentation."""
        # Forward mutating API calls to Presentation while batch mode is active.
        target = getattr(self._presentation, name)
        if callable(target):
            return target
        raise AttributeError(name)
