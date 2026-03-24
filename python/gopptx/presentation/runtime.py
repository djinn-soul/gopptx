"""Runtime method mixin for the core Presentation implementation."""

from __future__ import annotations

import ctypes
import threading
from typing import TYPE_CHECKING, cast

from .. import ops
from ..api_errors import GopptxError
from .batch import BatchContext
from .helpers import PresentationProtocol, json_dumps, json_loads, with_key_aliases

if TYPE_CHECKING:
    from ..schemas import BatchItemResult, PresentationMetadata, SlideMetadata
    from ..slide.slide import Slide


class PresentationRuntimeMixin:
    """Methods for executing commands and managing presentation runtime state."""

    _lib = None

    if TYPE_CHECKING:

        @property
        def slides(self) -> list[Slide]:
            """Cached slide proxy list."""
            ...

    def __init__(self) -> None:
        """Initialize runtime-managed state."""
        super().__init__()
        self._lock = threading.RLock()
        self._handle: int | None = None
        self._batch_active = False
        self._batch_stop_on_error = False
        self._batch_commands: list[dict[str, object]] = []
        self._slides_metadata_cache: list[SlideMetadata] | None = None
        self._metadata_cache: PresentationMetadata | None = None
        self._comment_ref_cache: dict[int, tuple[int, int, int]] = {}
        self._request_counter: int = 0

    @property
    def handle(self) -> int | None:
        """The internal handle to the Go engine deck."""
        return self._handle

    @property
    def batch_active(self) -> bool:
        """Return whether the presentation is currently buffering a batch."""
        return self._batch_active

    def execute_batch(
        self, commands: list[dict[str, object]], *, stop_on_error: bool = False
    ) -> list[BatchItemResult]:
        """Execute multiple bridge commands in one boundary crossing."""
        if not commands:
            return []
        result = self.execute(
            ops.OP_BATCH_EXECUTE, {"commands": commands, "stop_on_error": stop_on_error}
        )
        self.invalidate_cache()
        return cast("list[BatchItemResult]", result.get("results", []))

    def batch(self, *, stop_on_error: bool = False) -> BatchContext:
        """Context manager for buffered mutating operations."""
        return BatchContext(
            cast("PresentationProtocol", self), stop_on_error=stop_on_error
        )

    def begin_batch(self, *, stop_on_error: bool = False) -> None:
        """Begin buffering operations for a batch execute."""
        with self._lock:
            if self._batch_active:
                raise GopptxError(
                    "nested batch() calls are not allowed",
                    code="BATCH_NESTED_NOT_ALLOWED",
                )
            self._batch_active = True
            self._batch_stop_on_error = stop_on_error
            self._batch_commands = []

    def abort_batch(self) -> None:
        """Abort the current batch operation."""
        with self._lock:
            self._batch_active = False
            self._batch_stop_on_error = False
            self._batch_commands = []

    def end_batch(self) -> list[BatchItemResult]:
        """End the batch operation and execute queued commands."""
        with self._lock:
            commands = self._batch_commands
            stop_on_error = self._batch_stop_on_error
            self._batch_active = False
            self._batch_stop_on_error = False
            self._batch_commands = []
        return (
            self.execute_batch(commands, stop_on_error=stop_on_error)
            if commands
            else []
        )

    def execute(
        self, op: str, payload: dict[str, object] | None = None
    ) -> dict[str, object]:
        """Execute a bridge operation."""
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")

            if self._batch_active and op != ops.OP_BATCH_EXECUTE:
                if op in BatchContext.READ_OPS:
                    raise GopptxError(
                        f"read operation {op!r} is not allowed inside batch()",
                        code="BATCH_READ_OP_NOT_ALLOWED",
                    )
                self._batch_commands.append({"op": op, "payload": payload or {}})
                return {"_batched": True}

            self._request_counter += 1
            envelope = {
                "api_version": 1,
                "request_id": str(self._request_counter),
                "op": op,
                "payload": payload or {},
            }
            res_ptr = self._lib.deck_execute_json(self._handle, json_dumps(envelope))  # type: ignore[attr-defined]
            if not res_ptr:
                raise GopptxError("Received null response from Go engine")
            try:
                response = json_loads(ctypes.string_at(cast("int", res_ptr)))
                response_dict = cast("dict[str, object]", response)
                if not response_dict.get("ok", False):
                    err = cast("dict[str, object]", response_dict.get("error", {}))
                    raise GopptxError(
                        str(err.get("message", "Unknown error")),
                        code=str(err.get("code", "UNKNOWN")),
                    )
                result = response_dict.get("result")
                if result is None:
                    return {}
                if not isinstance(result, dict):
                    raise GopptxError("Invalid response payload type")
                return cast("dict[str, object]", result)
            finally:
                self._lib.deck_free_string(res_ptr)  # type: ignore[attr-defined]

    @property
    def slide_count(self) -> int:
        """The number of slides in the presentation."""
        val = self.execute(ops.OP_SLIDE_COUNT, {}).get("count", 0)
        return int(cast("int", val))

    @property
    def metadata(self) -> PresentationMetadata:
        """The presentation metadata."""
        with self._lock:
            if self._metadata_cache is not None:
                return self._metadata_cache
            self._metadata_cache = cast(
                "PresentationMetadata", self.execute(ops.OP_GET_METADATA, {})
            )
            return self._metadata_cache

    def __getitem__(self, index: int | slice) -> Slide | list[Slide]:
        """Return a slide or list of slides by index or slice.

        Args:
            index: Integer index or slice object.

        Returns:
            Slide object if index is int, list of Slides if slice.

        Raises:
            TypeError: If index is not an integer or slice.
            IndexError: If index is out of range.
        """
        count = self.slide_count
        if isinstance(index, int):
            if index < 0:
                index += count
            if not (0 <= index < count):
                raise IndexError("Slide index out of range")
            return self.slides[index]

        if type(index) is slice:
            indices = range(*index.indices(self.slide_count))
            slides = self.slides
            return [slides[i] for i in indices]

        raise TypeError("Slide index must be an integer or a slice")

    @property
    def slides_metadata(self) -> list[SlideMetadata]:
        """List of metadata for all slides in the presentation."""
        with self._lock:
            if self._slides_metadata_cache is not None:
                return self._slides_metadata_cache
            slides = self.execute(ops.OP_LIST_SLIDES, {}).get("slides", [])
            self._slides_metadata_cache = cast(
                "list[SlideMetadata]", with_key_aliases(slides)
            )
            return self._slides_metadata_cache

    def invalidate_cache(self) -> None:
        """Clear all cached data for the presentation."""
        with self._lock:
            self._slides_metadata_cache = None
            self._metadata_cache = None
            self._comment_ref_cache = {}
            invalidate_slide_lookup = getattr(
                self, "_invalidate_slide_lookup_cache", None
            )
            if callable(invalidate_slide_lookup):
                invalidate_slide_lookup()
            invalidate_slide_proxy = getattr(
                self, "_invalidate_slide_proxy_cache", None
            )
            if callable(invalidate_slide_proxy):
                invalidate_slide_proxy()
