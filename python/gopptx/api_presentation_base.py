from __future__ import annotations

import ctypes
import json
import os
import pathlib
import re
import sys
import threading
import uuid
from typing import Any, Dict, Optional, cast

from . import ops
from .api_batch import _BatchContext
from .api_errors import GopptxError
from .api_slide import Slide
from .types import PresentationMetadata, SlideMetadata

try:
    import orjson as _orjson  # type: ignore[import-not-found]
except ImportError:
    _orjson = None


def _json_dumps(payload: Dict[str, Any]) -> bytes:
    if _orjson is not None:
        return _orjson.dumps(payload)
    return json.dumps(payload, separators=(",", ":")).encode("utf-8")


def _json_loads(raw: bytes) -> Any:
    if _orjson is not None:
        return _orjson.loads(raw)
    return json.loads(raw.decode("utf-8"))


def _snake_case(name: str) -> str:
    s1 = re.sub(r"(.)([A-Z][a-z]+)", r"\1_\2", name)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", s1).lower()


def _with_key_aliases(obj: Any) -> Any:
    if isinstance(obj, list):
        return [_with_key_aliases(item) for item in obj]
    if not isinstance(obj, dict):
        return obj
    out: Dict[str, Any] = {}
    for k, v in obj.items():
        out[k] = _with_key_aliases(v)
        out[k.lower()] = out[k]
        out[_snake_case(k)] = out[k]
    return out


class PresentationBase:
    _lib = None
    _lib_lock = threading.Lock()

    def __init__(self, path: Optional[str] = None):
        self._load_library()
        self._lock = threading.RLock()
        self._handle: Optional[int] = None
        self._slides_metadata_cache: Optional[list[SlideMetadata]] = None
        self._metadata_cache: Optional[PresentationMetadata] = None
        self._batch_active = False
        self._batch_stop_on_error = False
        self._batch_commands: list[dict] = []
        self._comment_ref_cache: Dict[int, tuple[int, int, int]] = {}
        if path:
            self.open(path)

    @classmethod
    def _load_library(cls) -> None:
        with cls._lib_lock:
            if cls._lib:
                return
            pkg_dir = os.path.dirname(__file__)
            lib_name = (
                "gopptx.dll"
                if sys.platform == "win32"
                else ("libgopptx.dylib" if sys.platform == "darwin" else "libgopptx.so")
            )
            search_paths: list[str] = []
            env_path = os.environ.get("GOPPTX_LIB_PATH")
            if env_path:
                if pathlib.Path(env_path).is_dir():
                    search_paths.append(os.path.join(env_path, lib_name))
                else:
                    search_paths.append(env_path)
            search_paths.extend([
                os.path.join(pkg_dir, lib_name),
                os.path.join(pkg_dir, "../../bindings/c/build", lib_name),
            ])
            lib_path = next(
                (os.path.abspath(c) for c in search_paths if pathlib.Path(c).exists()), None
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

    def execute(
        self, op: str, payload: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")
            if self._batch_active and op != ops.OP_BATCH_EXECUTE:
                if op in _BatchContext._READ_OPS:
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
                return cast(Dict[str, Any], result)
            finally:
                self._lib.deck_free_string(res_ptr)

    def execute_batch(
        self, commands: list[dict], stop_on_error: bool = False
    ) -> list[Dict[str, Any]]:
        """Execute multiple bridge commands in one boundary crossing.

        Returns ordered per-command results. Each result has `ok` plus either
        `result` or `error` fields from the Go bridge.
        """
        if not commands:
            return []
        result = self.execute(
            ops.OP_BATCH_EXECUTE, {"commands": commands, "stop_on_error": stop_on_error}
        )
        self.invalidate_cache()
        return cast(list[Dict[str, Any]], result.get("results", []))

    def batch(self, stop_on_error: bool = False) -> _BatchContext:
        """Context manager for buffered mutating operations executed as one batch."""
        return _BatchContext(self, stop_on_error=stop_on_error)

    def _begin_batch(self, stop_on_error: bool) -> None:
        with self._lock:
            if self._batch_active:
                raise GopptxError(
                    "nested batch() calls are not allowed",
                    code="BATCH_NESTED_NOT_ALLOWED",
                )
            self._batch_active = True
            self._batch_stop_on_error = stop_on_error
            self._batch_commands = []

    def _abort_batch(self) -> None:
        with self._lock:
            self._batch_active = False
            self._batch_stop_on_error = False
            self._batch_commands = []

    def _end_batch(self) -> list[Dict[str, Any]]:
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
        return int(self.execute(ops.OP_SLIDE_COUNT, {}).get("count", 0))

    @property
    def metadata(self) -> PresentationMetadata:
        with self._lock:
            if self._metadata_cache is not None:
                return self._metadata_cache
            self._metadata_cache = cast(
                PresentationMetadata, self.execute(ops.OP_GET_METADATA, {})
            )
            return self._metadata_cache

    @property
    def slides(self) -> list[Slide]:
        return [Slide(self, m) for m in self.slides_metadata]

    @property
    def slides_metadata(self) -> list[SlideMetadata]:
        with self._lock:
            if self._slides_metadata_cache is not None:
                return self._slides_metadata_cache
            slides = self.execute(ops.OP_LIST_SLIDES, {}).get("slides", [])
            self._slides_metadata_cache = cast(
                list[SlideMetadata], _with_key_aliases(slides)
            )
            return self._slides_metadata_cache

    def invalidate_cache(self) -> None:
        with self._lock:
            self._slides_metadata_cache = None
            self._metadata_cache = None
            self._comment_ref_cache = {}

    def _get_last_error(self) -> str:
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
        with self._lock:
            if not self._handle:
                raise GopptxError("Presentation is not open.")
            if self._lib.deck_save(self._handle, path.encode("utf-8")) != 0:
                raise GopptxError(f"Failed to save deck: {self._get_last_error()}")

    def close(self) -> None:
        with self._lock:
            if self._handle:
                self._lib.deck_close(self._handle)
                self._handle = None
            self.invalidate_cache()

    def __enter__(self) -> PresentationBase:
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.close()

    def __del__(self) -> None:
        try:
            self.close()
        except Exception:
            pass

    def __repr__(self) -> str:
        return f"<Presentation title='{self.title}' slides={self.slide_count}>"
