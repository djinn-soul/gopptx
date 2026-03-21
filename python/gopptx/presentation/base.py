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
    def from_template(
        cls,
        path: str,
        context: dict[str, object],
        *,
        undefined: str = "keep",
    ) -> PresentationBase:
        """Open a .pptx file as a Jinja2 template and render tags with context data.

        Finds all Jinja2 expressions (``{{ var }}``, ``{{ var | filter }}``, etc.)
        inside shape text across every slide and replaces them with rendered values,
        preserving the original run formatting.

        Args:
            path: Path to the .pptx template file containing Jinja2 tags.
            context: Variable mapping passed to the Jinja2 renderer.
            undefined: How to handle missing variables.
                ``"keep"``   – leave unresolved tags as-is (default).
                ``"empty"``  – replace unresolved tags with an empty string.
                ``"strict"`` – raise ``UndefinedError`` on missing variables.

        Returns:
            A :class:`Presentation` with all Jinja2 tags replaced by their values.

        Raises:
            ImportError: If ``jinja2`` is not installed.
            jinja2.UndefinedError: When ``undefined="strict"`` and a variable is missing.

        Example::

            prs = Presentation.from_template(
                "template.pptx",
                context={"title": "Q4 Launch", "author": "Alice"},
            )
            prs.save("output.pptx")
        """
        try:
            from jinja2 import DebugUndefined, Environment, StrictUndefined, Undefined
        except ImportError as exc:
            raise ImportError(
                "jinja2 is required for from_template(). "
                "Install it with: pip install jinja2"
            ) from exc

        _undefined_map = {
            "keep": DebugUndefined,
            "empty": Undefined,
            "strict": StrictUndefined,
        }
        env = Environment(undefined=_undefined_map.get(undefined, DebugUndefined))  # type: ignore[arg-type]

        from .. import ops as _ops

        pres = cls(path)

        # Collect all raw text tokens that contain Jinja2 expressions.
        # Multi-paragraph shape text is joined with "\n" in the state, but
        # find_and_replace works on individual XML run text. We therefore split
        # on newlines and process each line (run-level token) separately.
        seen: dict[str, str] = {}  # raw_token -> rendered_token

        for slide_index in range(pres.slide_count):  # type: ignore[attr-defined]
            states: list[dict[str, object]] = cast(
                "list[dict[str, object]]",
                pres.execute(
                    _ops.OP_GET_SLIDE_TEXT_STATES, {"slide_index": slide_index}
                ).get("states", []),  # type: ignore[attr-defined]
            )
            for state in states:
                full_text = str(state.get("text", ""))
                if not full_text:
                    continue
                # Jinja2 vars may span a single run or a whole paragraph.
                # Split on paragraph boundaries (\n) so each token matches
                # what find_and_replace will see in the XML text node.
                for token in full_text.split("\n"):
                    if not token or ("{{" not in token and "{%" not in token):
                        continue
                    if token in seen:
                        continue
                    rendered = env.from_string(token).render(**context)
                    seen[token] = rendered

        # Apply all replacements in a single pass.
        for raw, rendered in seen.items():
            if raw != rendered:
                pres.execute(
                    _ops.OP_FIND_AND_REPLACE, {"find": raw, "replace": rendered}
                )  # type: ignore[attr-defined]

        return pres

    def render_template(self, context: dict[str, object]) -> int:
        """Render all Jinja2 template expressions in slide shapes using context values.

        Supports the full Jinja2 syntax — variables, filters, conditionals, and
        loops — via the Go-side gonja engine.  Operates on the presentation
        *in place*, preserving run-level formatting (bold, colour, etc.).

        Args:
            context: Variable mapping passed to the Jinja2 renderer.

        Returns:
            Number of text-run replacements performed.
        """
        from .. import ops as _ops

        result = self.execute(_ops.OP_RENDER_TEMPLATE, {"context": context})  # type: ignore[attr-defined]
        return int(result.get("replacements", 0))

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
