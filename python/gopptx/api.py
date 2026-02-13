from __future__ import annotations

import ctypes
import json
import os
import sys
import uuid
from typing import Any, Dict, Optional, cast

from . import ops
from .types import PresentationMetadata

class GopptxError(Exception):
    """Base exception for gopptx library errors."""

    def __init__(self, message: str, code: Optional[str] = None):
        super().__init__(message)
        self.code = code


class Presentation:
    """High-level wrapper for a PowerPoint presentation handled by the Go engine."""

    _lib = None

    def __init__(self, path: Optional[str] = None):
        self._load_library()
        self._handle: Optional[int] = None
        if path:
            self.open(path)

    @classmethod
    def _load_library(cls) -> None:
        if cls._lib:
            return

        pkg_dir = os.path.dirname(__file__)
        if sys.platform == "win32":
            lib_name = "gopptx.dll"
        elif sys.platform == "darwin":
            lib_name = "libgopptx.dylib"
        else:
            lib_name = "libgopptx.so"

        search_paths = [
            os.path.join(pkg_dir, lib_name),
            os.path.join(pkg_dir, "../../bindings/c/build", lib_name),
        ]

        lib_path = None
        for candidate in search_paths:
            if os.path.exists(candidate):
                lib_path = os.path.abspath(candidate)
                break

        if not lib_path:
            raise GopptxError(
                f"Could not find shared library {lib_name}. Please build it first."
            )

        cls._lib = ctypes.CDLL(lib_path)
        cls._lib.deck_open.argtypes = [ctypes.c_char_p]
        cls._lib.deck_open.restype = ctypes.c_void_p

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

    def execute(self, op: str, payload: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Execute a command operation against the Go engine."""
        if not self._handle:
            raise GopptxError("Presentation is not open.")

        envelope = {
            "api_version": 1,
            "request_id": str(uuid.uuid4()),
            "op": op,
            "payload": payload or {},
        }
        json_input = json.dumps(envelope).encode("utf-8")
        res_ptr = self._lib.deck_execute_json(self._handle, json_input)
        if not res_ptr:
            raise GopptxError("Received null response from Go engine")

        try:
            res_str = ctypes.string_at(res_ptr).decode("utf-8")
            response = json.loads(res_str)
            if not response.get("ok", False):
                err = response.get("error", {})
                raise GopptxError(
                    err.get("message", "Unknown engine error"),
                    code=err.get("code"),
                )
            result = response.get("result")
            if result is None:
                return {}
            if not isinstance(result, dict):
                raise GopptxError("Invalid response payload type")
            return cast(Dict[str, Any], result)
        finally:
            self._lib.deck_free_string(res_ptr)

    @property
    def slide_count(self) -> int:
        """Returns the total number of slides in the presentation."""
        result = self.execute(ops.OP_SLIDE_COUNT, {})
        return int(result.get("count", 0))

    @property
    def metadata(self) -> PresentationMetadata:
        """Returns presentation metadata as a dictionary with keys 'title', 'slide_count', and 'size'."""
        result = self.execute(ops.OP_GET_METADATA, {})
        return cast(PresentationMetadata, result)

    def add_slide(self, title: str) -> int:
        """Adds a new slide and returns its index."""
        result = self.execute(ops.OP_ADD_SLIDE, {"title": title})
        return int(result.get("index", -1))

    def remove_slide(self, index: int) -> None:
        """Removes the slide at the specified zero-based index."""
        result = self.execute(ops.OP_REMOVE_SLIDE, {"index": index})
        return None

    def move_slide(self, from_index: int, to_index: int) -> None:
        """Moves a slide from one position to another."""
        result = self.execute(ops.OP_MOVE_SLIDE, {"from": from_index, "to": to_index})
        return None

    def duplicate_slide(self, index: int, insert_at: Optional[int] = None) -> int:
        """Duplicates a slide and returns the new slide index."""
        if insert_at is None:
            insert_at = index + 1
        result = self.execute(ops.OP_DUPLICATE_SLIDE, {"index": index, "insert_at": insert_at})
        return int(result.get("new_index", -1))


    def _get_last_error(self) -> str:
        if not self._handle:
            return "No active session"
        err_ptr = self._lib.deck_last_error(self._handle)
        if err_ptr:
            err_msg = ctypes.string_at(err_ptr).decode("utf-8")
            self._lib.deck_free_string(err_ptr)
            return err_msg
        return "Unknown error"

    def open(self, path: str) -> None:
        """Open an existing PPTX file."""
        if self._handle:
            self.close()
        handle = self._lib.deck_open(path.encode("utf-8"))
        if not handle:
            err_ptr = self._lib.deck_global_error()
            msg = ctypes.string_at(err_ptr).decode("utf-8") if err_ptr else "Unknown error"
            if err_ptr:
                self._lib.deck_free_string(err_ptr)
            raise GopptxError(f"Failed to open deck: {msg}")
        self._handle = int(handle)

    def save(self, path: str) -> None:
        """Save the presentation to the specified path."""
        if not self._handle:
            raise GopptxError("Presentation is not open.")
        rc = self._lib.deck_save(self._handle, path.encode("utf-8"))
        if rc != 0:
            raise GopptxError(f"Failed to save deck: {self._get_last_error()}")

    def close(self) -> None:
        """Close the presentation and release resources."""
        if self._handle:
            self._lib.deck_close(self._handle)
            self._handle = None

    def __enter__(self) -> "Presentation":
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.close()

    def __del__(self) -> None:
        self.close()
