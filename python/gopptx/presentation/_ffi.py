"""Low-level helpers shared by the ctypes bridge call sites."""

from __future__ import annotations

import ctypes
from typing import Protocol, cast


class _StringFreeingLib(Protocol):
    """The subset of the bridge needed to release an owned error string."""

    def deck_free_string(self, ptr: int) -> None: ...


def take_error(lib: object, err_ptr: int | None) -> str:
    """Decode and free an error string returned through a bridge out-parameter.

    The bridge allocates the message with C.CString, so ownership transfers to
    us and it must be released with deck_free_string.
    """
    if not err_ptr:
        return "Unknown error"
    msg = ctypes.string_at(err_ptr).decode("utf-8")
    cast("_StringFreeingLib", lib).deck_free_string(err_ptr)
    return msg


__all__ = ["take_error"]
