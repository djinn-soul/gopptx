"""Base presentation class with core functionality for gopptx library."""

from __future__ import annotations

import contextlib
import ctypes
import json
import os
import pathlib
import re
import sys
import threading
import uuid
from typing import TYPE_CHECKING, Any, cast

from typing_extensions import Self

from . import ops
from .api_batch import BatchContext
from .api_errors import GopptxError
from .api_slide import Slide

if TYPE_CHECKING:
    from types import TracebackType

    from .schemas import (
        BatchCommand,
        BatchItemResult,
        PresentationMetadata,
        SlideMetadata,
    )

try:
    import orjson as _orjson  # type: ignore[import-not-found]
except ImportError:
    _orjson = None


def _json_dumps(payload: dict[str, Any]) -> bytes:
    """Serialize a dictionary to JSON bytes."""
    if _orjson is not None:
        return _orjson.dumps(payload)
    return json.dumps(payload, separators=(",", ":")).encode("utf-8")


def _json_loads(raw: bytes) -> Any:  # noqa: ANN401 - JSON return type is intentionally dynamic
    """Deserialize JSON bytes to Python objects."""
    if _orjson is not None:
        return _orjson.loads(raw)
    return json.loads(raw.decode("utf-8"))


def _snake_case(name: str) -> str:
    """Convert CamelCase to snake_case."""
    s1 = re.sub(r"(.)([A-Z][a-z]+)", r"\1_\2", name)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", s1).lower()


def _with_key_aliases(obj: Any) -> Any:  # noqa: ANN401 - Recursive JSON transform is intentionally dynamic
    """Add lowercase and snake_case aliases for all keys in a nested structure."""
    if isinstance(obj, list):
        return [_with_key_aliases(item) for item in obj]
    if not isinstance(obj, dict):
        return obj
    out: dict[str, Any] = {}
    for k, v in obj.items():
        out[k] = _with_key_aliases(v)
        out[k.lower()] = out[k]
        out[_snake_case(k)] = out[k]
    return out


class PresentationBase:
    """Base class for Presentation with core library loading and execution."""

    _lib = None
    _lib_lock = threading.Lock()

    def __init__(self, path: str | None = None) -> None:
        """Initialize the presentation, optionally opening a file.

        Args:
            path: Optional path to an existing presentation file to open.
        """
        self._load_library()
        self._lock = threading.RLock()
        self._handle: int | None = None
        self._slides_metadata_cache: list[SlideMetadata] | None = None
        self._metadata_cache: PresentationMetadata | None = None
        self._batch_active = False
        self._batch_stop_on_error = False
        self._batch_commands: list[dict[str, Any]] = []
        self._comment_ref_cache: dict[int, tuple[int, int, int]] = {}
        if path:
            self.open(path)

    @classmethod
    def _load_library(cls) -> None:
        with cls._lib_lock:
            if cls._lib:
                return
            pkg_dir = pathlib.Path(__file__).parent
            lib_name = (
                "gopptx.dll"
                if sys.platform == "win32"
                else ("libgopptx.dylib" if sys.platform == "darwin" else "libgopptx.so")
            )
            search_paths: list[pathlib.Path] = []
            env_path = os.environ.get("GOPPTX_LIB_PATH")
            if env_path:
                env_as_path = pathlib.Path(env_path)
                if env_as_path.is_dir():
                    search_paths.append(env_as_path / lib_name)
                else:
                    search_paths.append(env_as_path)
            search_paths.extend([
                pkg_dir / "../../bindings/c/build" / lib_name,
                pkg_dir / lib_name,
            ])
            lib_path = next(
                (c.resolve() for c in search_paths if c.exists()),
                None,
            )
            if not lib_path:
                raise GopptxError(
                    f"Could not find shared library {lib_name}. Please build it first."
                )

            cls._lib = ctypes.CDLL(lib_path)
            cls._lib.deck_open.argtypes = [ctypes.c_char_p]
            cls._lib.deck_open.restype = ctypes.c_void_p
            cls._lib.deck_new.argtypes = [ctypes.c_char_p]
            cls._lib.deck_new.restype = ctypes.c_void_p
            cls._lib.deck_execute_json.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
            cls._lib.deck_execute_json.restype = ctypes.c_void_p
            cls._lib.deck_save.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
            cls._lib.deck_save.restype = ctypes.c_int
            cls._lib.deck_last_error.argtypes = [ctypes.c_void_p]
            cls._lib.deck_last_error.restype = ctypes.c_void_p
            cls._lib.deck_global_error.argtypes = []
            cls._lib.deck_global_error.restype = ctypes.c_void_p
            cls._lib.deck_free_string.argtypes = [ctypes.c_void_p]
            cls._lib.deck_free_string.restype = None
            cls._lib.deck_close.argtypes = [ctypes.c_void_p]
            cls._lib.deck_close.restype = None

    @classmethod
    def new(cls, title: str) -> PresentationBase:
        """Create a new presentation with the given title.

        Args:
            title: The title for the new presentation.

        Returns:
            A new PresentationBase instance.
        """
        pres = cls()
        handle = cls._lib.deck_new(title.encode("utf-8"))
        if not handle:
            err_ptr = cls._lib.deck_global_error()
            msg = (
                ctypes.string_at(err_ptr).decode("utf-8")
                if err_ptr
                else "Unknown error"
            )
            if err_ptr:
                cls._lib.deck_free_string(err_ptr)
            raise GopptxError(f"Failed to create new deck: {msg}")
        pres._handle = int(handle)
        return pres

    def execute(self, op: str, payload: dict[str, Any] | None = None) -> dict[str, Any]:
        """Execute a bridge operation.

        Args:
            op: The operation name.
            payload: Optional payload for the operation.

        Returns:
            The result dictionary from the operation.
        """
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

            envelope = {
                "api_version": 1,
                "request_id": str(uuid.uuid4()),
                "op": op,
                "payload": payload or {},
            }
            res_ptr = self._lib.deck_execute_json(self._handle, _json_dumps(envelope))
            if not res_ptr:
                raise GopptxError("Received null response from Go engine")
            try:
                response = _json_loads(ctypes.string_at(res_ptr))
                if not response.get("ok", False):
                    err = response.get("error", {})
                    raise GopptxError(
                        err.get("message", "Unknown engine error"), code=err.get("code")
                    )
                result = response.get("result")
                if result is None:
                    return {}
                if not isinstance(result, dict):
                    raise GopptxError("Invalid response payload type")
                return cast("dict[str, Any]", result)
            finally:
                self._lib.deck_free_string(res_ptr)

    def execute_batch(
        self, commands: list[BatchCommand], *, stop_on_error: bool = False
    ) -> list[BatchItemResult]:
        """Execute multiple bridge commands in one boundary crossing.

        Returns ordered per-command results. Each result has `ok` plus either
        `result` or `error` fields from the Go bridge.

        Args:
            commands: List of batch commands to execute.
            stop_on_error: If True, stop execution on first error.

        Returns:
            List of results for each command.
        """
        if not commands:
            return []
        result = self.execute(
            ops.OP_BATCH_EXECUTE, {"commands": commands, "stop_on_error": stop_on_error}
        )
        self.invalidate_cache()
        return cast("list[BatchItemResult]", result.get("results", []))

    def batch(self, *, stop_on_error: bool = False) -> BatchContext:
        """Context manager for buffered mutating operations executed as one batch.

        Args:
            stop_on_error: If True, stop batch execution on first error.

        Returns:
            A batch context manager.
        """
        return BatchContext(self, stop_on_error=stop_on_error)

    def begin_batch(self, *, stop_on_error: bool = False) -> None:
        """Begin a batch operation context.

        Args:
            stop_on_error: If True, stop batch execution on first error.
        """
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
        """End the batch operation and execute all queued commands.

        Returns:
            List of results for each batched command.
        """
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

    @property
    def slide_count(self) -> int:
        """The number of slides in the presentation."""
        return int(self.execute(ops.OP_SLIDE_COUNT, {}).get("count", 0))

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

    @property
    def slides(self) -> list[Slide]:
        """List of slide proxies for all slides in the presentation."""
        return [Slide(self, m) for m in self.slides_metadata]

    @property
    def slides_metadata(self) -> list[SlideMetadata]:
        """List of metadata for all slides in the presentation."""
        with self._lock:
            if self._slides_metadata_cache is not None:
                return self._slides_metadata_cache
            slides = self.execute(ops.OP_LIST_SLIDES, {}).get("slides", [])
            self._slides_metadata_cache = cast(
                "list[SlideMetadata]", _with_key_aliases(slides)
            )
            return self._slides_metadata_cache

    def invalidate_cache(self) -> None:
        """Clear all cached data for the presentation."""
        with self._lock:
            self._slides_metadata_cache = None
            self._metadata_cache = None
            self._comment_ref_cache = {}

    def _get_last_error(self) -> str:
        """Get the last error message from the Go engine."""
        with self._lock:
            if not self._handle:
                return "No active session"
            err_ptr = self._lib.deck_last_error(self._handle)
            if err_ptr:
                err_msg = ctypes.string_at(err_ptr).decode("utf-8")
                self._lib.deck_free_string(err_ptr)
                return err_msg
            return "Unknown error"

    def open(self, path: str) -> None:
        """Open an existing presentation file.

        Args:
            path: Path to the presentation file.
        """
        with self._lock:
            if self._handle:
                self.close()
            handle = self._lib.deck_open(path.encode("utf-8"))
            if not handle:
                err_ptr = self._lib.deck_global_error()
                msg = (
                    ctypes.string_at(err_ptr).decode("utf-8")
                    if err_ptr
                    else "Unknown error"
                )
                if err_ptr:
                    self._lib.deck_free_string(err_ptr)
                raise GopptxError(f"Failed to open deck: {msg}")
            self._handle = int(handle)
            self.invalidate_cache()

    def save(self, path: str) -> None:
        """Save the presentation to a file.

        Args:
            path: Path to save the presentation to.
        """
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")
            if self._lib.deck_save(self._handle, path.encode("utf-8")) != 0:
                raise GopptxError(f"Failed to save deck: {self._get_last_error()}")

    def close(self) -> None:
        """Close the presentation and release resources."""
        with self._lock:
            if self._handle:
                self._lib.deck_close(self._handle)
                self._handle = None
            self.invalidate_cache()

    def __enter__(self) -> Self:
        """Enter context manager."""
        return self

    def __exit__(
        self,
        exc_type: type[BaseException] | None,
        exc_val: BaseException | None,
        exc_tb: TracebackType | None,
    ) -> None:
        """Exit context manager and close presentation."""
        self.close()

    def __del__(self) -> None:
        """Clean up resources on deletion."""
        with contextlib.suppress(Exception):
            self.close()

    def __repr__(self) -> str:
        """Return string representation of the presentation."""
        title = self.metadata.get("title", "") if self._handle else ""
        return f"<Presentation title='{title}' slides={self.slide_count}>"
