"""Lifecycle method mixin for opening, saving, and closing presentations."""
# pyright: reportUnknownMemberType=false, reportAttributeAccessIssue=false, reportUninitializedInstanceVariable=false, reportUnnecessaryTypeIgnoreComment=false, reportUnknownVariableType=false

from __future__ import annotations

import contextlib
import ctypes
from typing import TYPE_CHECKING, cast

from typing_extensions import Self

from ..api_errors import GopptxError
from ._ffi import take_error

if TYPE_CHECKING:
    from types import TracebackType


class PresentationRuntimeLifecycleMixin:
    """Methods for managing presentation lifecycle and context semantics."""

    def _get_last_error(self) -> str:
        """Get the last error message from the Go engine."""
        with self._lock:
            if not self._handle:
                return "No active session"
            err_ptr = self._lib.deck_last_error(self._handle)  # type: ignore[attr-defined]
            if err_ptr:
                err_msg = ctypes.string_at(cast("int", err_ptr)).decode("utf-8")
                self._lib.deck_free_string(err_ptr)  # type: ignore[attr-defined]
                return err_msg
            return "Unknown error"

    def open(self, path: str) -> None:
        """Open an existing presentation file."""
        with self._lock:
            if self._handle:
                self.close()
            err = ctypes.c_void_p()
            handle = cast(
                "int",
                self._lib.deck_open_ex(str(path).encode("utf-8"), ctypes.byref(err)),  # type: ignore[attr-defined]
            )
            if not handle:
                raise GopptxError(
                    f"Failed to open deck: {take_error(cast('object', self._lib), err.value)}"
                )
            self._handle = int(handle)
            self.invalidate_cache()

    def open_bytes(self, data: bytes) -> None:
        """Open a presentation from an in-memory byte string."""
        with self._lock:
            if self._handle:
                self.close()
            buf = (ctypes.c_char * len(data)).from_buffer_copy(data)
            err = ctypes.c_void_p()
            handle = cast(
                "int",
                self._lib.deck_open_bytes_ex(
                    buf, ctypes.c_int(len(data)), ctypes.byref(err)
                ),  # type: ignore[attr-defined]
            )
            if not handle:
                raise GopptxError(
                    f"Failed to open deck from bytes: {take_error(cast('object', self._lib), err.value)}"
                )
            self._handle = int(handle)
            self.invalidate_cache()

    def save(self, path: str) -> None:
        """Save the presentation to a file."""
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")
            status = self._lib.deck_save(self._handle, str(path).encode("utf-8"))  # type: ignore[attr-defined]
            if status != 0:
                raise GopptxError(f"Failed to save deck: {self._get_last_error()}")

    def to_bytes(self) -> bytes:
        """Serialize the presentation to bytes without writing to disk."""
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")
            out_len = ctypes.c_int(0)
            ptr = self._lib.deck_save_bytes(self._handle, ctypes.byref(out_len))  # type: ignore[attr-defined]
            if not ptr:
                raise GopptxError(f"Failed to serialize deck: {self._get_last_error()}")
            try:
                return bytes(ctypes.string_at(ptr, out_len.value))  # type: ignore[arg-type]
            finally:
                self._lib.deck_free_string(ptr)  # type: ignore[attr-defined]

    def close(self) -> None:
        """Close the presentation and release resources."""
        with self._lock:
            if self._handle:
                self._lib.deck_close(self._handle)  # type: ignore[attr-defined]
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

    def __repr__(self) -> str:  # type: ignore[override]
        """Return string representation of the presentation."""
        title = ""
        slide_count = 0
        if self._handle:
            try:
                title = self.metadata.get("title", "")
                slide_count = self.slide_count
            except GopptxError:
                pass
        return f"<Presentation title='{title}' slides={slide_count}>"
