"""Lifecycle method mixin for opening, saving, and closing presentations."""
# pyright: reportUnknownMemberType=false, reportAttributeAccessIssue=false, reportUninitializedInstanceVariable=false, reportUnnecessaryTypeIgnoreComment=false, reportUnknownVariableType=false

from __future__ import annotations

import contextlib
import ctypes
from typing import TYPE_CHECKING, cast

from typing_extensions import Self

from ..api_errors import GopptxError

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
            handle = cast("int", self._lib.deck_open(str(path).encode("utf-8")))  # type: ignore[attr-defined]
            if not handle:
                err_ptr = self._lib.deck_global_error()  # type: ignore[attr-defined]
                msg = (
                    ctypes.string_at(cast("int", err_ptr)).decode("utf-8")
                    if err_ptr
                    else "Unknown error"
                )
                if err_ptr:
                    self._lib.deck_free_string(err_ptr)  # type: ignore[attr-defined]
                raise GopptxError(f"Failed to open deck: {msg}")
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
