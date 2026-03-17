"""Base presentation class with core functionality for gopptx library."""

from __future__ import annotations

import ctypes
import os
import pathlib
import sys
import threading
from typing import cast

from ..api_errors import GopptxError
from .helpers import PresentationProtocol
from .runtime import PresentationRuntimeMixin
from .shapes.shape_write_buffer_mixin import PresentationShapeWriteBufferMixin
from .slides.slide_lookup_mixin import PresentationSlideLookupMixin
from .slides.slide_proxy_mixin import PresentationSlideProxyMixin
from .text.text_write_buffer_mixin import PresentationTextWriteBufferMixin


class PresentationBase(
    PresentationSlideProxyMixin,
    PresentationSlideLookupMixin,
    PresentationShapeWriteBufferMixin,
    PresentationTextWriteBufferMixin,
    PresentationRuntimeMixin,
):
    """Base class for Presentation with core library loading and execution."""

    _lib = None
    _lib_lock = threading.Lock()

    def __init__(self, path: str | None = None) -> None:
        """Initialize the presentation, optionally opening a file."""
        super().__init__()
        self._load_library()
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
                pkg_dir / "../../../bindings/c/build" / lib_name,
                pkg_dir.parent / lib_name,
            ])
            lib_path = next((c.resolve() for c in search_paths if c.exists()), None)
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
        """Create a new presentation with the given title."""
        pres = cls()
        handle = cast("int", cls._lib.deck_new(title.encode("utf-8")))  # type: ignore[attr-defined]
        if not handle:
            err_ptr = cls._lib.deck_global_error()  # type: ignore[attr-defined]
            msg = (
                ctypes.string_at(cast("int", err_ptr)).decode("utf-8")
                if err_ptr
                else "Unknown error"
            )
            if err_ptr:
                cls._lib.deck_free_string(err_ptr)  # type: ignore[attr-defined]
            raise GopptxError(f"Failed to create new deck: {msg}")
        pres._handle = int(handle)
        return pres


__all__ = ["PresentationBase", "PresentationProtocol"]
