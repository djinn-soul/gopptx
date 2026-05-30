"""Base presentation class with core functionality for gopptx library."""

from __future__ import annotations

import ctypes
import os
import pathlib
import sys
import threading
from typing import TYPE_CHECKING, Protocol, cast

from typing_extensions import Self

from .. import ops as _ops
from ..api_errors import GopptxError
from .helpers import PresentationProtocol
from .runtime import PresentationRuntimeMixin
from .runtime_lifecycle import PresentationRuntimeLifecycleMixin

try:
    import jinja2 as _jinja2
    from jinja2.sandbox import SandboxedEnvironment as _JinjaSandboxEnv
except ImportError:
    _jinja2 = None
    _JinjaSandboxEnv = None

if TYPE_CHECKING:

    class _PresentationTemplateProtocol(Protocol):
        @property
        def slide_count(self) -> int: ...

        def execute(
            self, op: str, payload: dict[str, object] | None = None
        ) -> dict[str, object]: ...

    class _JinjaTemplate(Protocol):
        def render(self, **kwargs: object) -> str: ...

    class _JinjaEnv(Protocol):
        def from_string(self, source: str) -> _JinjaTemplate: ...


from .shapes.shape_write_buffer_mixin import PresentationShapeWriteBufferMixin
from .slides.slide_lookup_mixin import PresentationSlideLookupMixin
from .slides.slide_proxy_mixin import PresentationSlideProxyMixin
from .text.text_write_buffer_mixin import PresentationTextWriteBufferMixin


class _GopptxLibProtocol(Protocol):
    def deck_new(self, title: bytes) -> int: ...

    def deck_global_error(self) -> int: ...

    def deck_free_string(self, ptr: int) -> None: ...


class _CtypesFuncProtocol(Protocol):
    argtypes: list[object]
    restype: object


class _RawGopptxLibProtocol(Protocol):
    deck_open: _CtypesFuncProtocol
    deck_new: _CtypesFuncProtocol
    deck_execute_json: _CtypesFuncProtocol
    deck_open_bytes: _CtypesFuncProtocol
    deck_save_bytes: _CtypesFuncProtocol
    deck_save: _CtypesFuncProtocol
    deck_last_error: _CtypesFuncProtocol
    deck_global_error: _CtypesFuncProtocol
    deck_free_string: _CtypesFuncProtocol
    deck_close: _CtypesFuncProtocol


def _collect_jinja_tokens(
    state: dict[str, object],
    env: _JinjaEnv,
    context: dict[str, object],
    seen: dict[str, str],
) -> None:
    """Collect Jinja2 tokens from a single text state into ``seen``."""
    full_text = str(state.get("text", ""))
    if not full_text:
        return
    # Jinja2 vars may span a single run or a whole paragraph.
    # Split on paragraph boundaries (\n) so each token matches
    # what find_and_replace will see in the XML text node.
    for token in full_text.split("\n"):
        if not token or ("{{" not in token and "{%" not in token):
            continue
        if token in seen:
            continue
        seen[token] = env.from_string(token).render(**context)


class PresentationBase(
    PresentationSlideProxyMixin,
    PresentationSlideLookupMixin,
    PresentationShapeWriteBufferMixin,
    PresentationTextWriteBufferMixin,
    PresentationRuntimeMixin,
    PresentationRuntimeLifecycleMixin,
):
    """Base class for Presentation with core library loading and execution."""

    _lib: object | None = None
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

            lib = cast("_RawGopptxLibProtocol", ctypes.CDLL(lib_path))
            lib.deck_open.argtypes = [ctypes.c_char_p]
            lib.deck_open.restype = ctypes.c_void_p
            lib.deck_new.argtypes = [ctypes.c_char_p]
            lib.deck_new.restype = ctypes.c_void_p
            lib.deck_execute_json.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
            lib.deck_execute_json.restype = ctypes.c_void_p
            lib.deck_open_bytes.argtypes = [ctypes.c_char_p, ctypes.c_int]
            lib.deck_open_bytes.restype = ctypes.c_void_p
            lib.deck_save_bytes.argtypes = [
                ctypes.c_void_p,
                ctypes.POINTER(ctypes.c_int),
            ]
            lib.deck_save_bytes.restype = ctypes.c_void_p
            lib.deck_save.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
            lib.deck_save.restype = ctypes.c_int
            lib.deck_last_error.argtypes = [ctypes.c_void_p]
            lib.deck_last_error.restype = ctypes.c_void_p
            lib.deck_global_error.argtypes = []
            lib.deck_global_error.restype = ctypes.c_void_p
            lib.deck_free_string.argtypes = [ctypes.c_void_p]
            lib.deck_free_string.restype = None
            lib.deck_close.argtypes = [ctypes.c_void_p]
            lib.deck_close.restype = None
            cls._lib = lib

    @classmethod
    def from_template(
        cls,
        path: str,
        context: dict[str, object],
        *,
        undefined: str = "keep",
    ) -> Self:
        """Open a .pptx file as a Jinja2 template and render tags with context data.

        Finds all Jinja2 expressions (``{{ var }}``, ``{{ var | filter }}``, etc.)
        inside shape text across every slide and replaces them with rendered values,
        preserving the original run formatting.

        Args:
            path: Path to the .pptx template file containing Jinja2 tags.
            context: Variable mapping passed to the Jinja2 renderer.
            undefined: How to handle missing variables.
                ``"keep"``   - leave unresolved tags as-is (default).
                ``"empty"``  - replace unresolved tags with an empty string.
                ``"strict"`` - raise ``UndefinedError`` on missing variables.

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
        if _jinja2 is None:
            raise ImportError(
                "jinja2 is required for from_template(). Install it with: pip install jinja2"
            )

        jinja2 = _jinja2
        undefined_map: dict[str, type[object]] = {
            "keep": jinja2.DebugUndefined,
            "empty": jinja2.Undefined,
            "strict": jinja2.StrictUndefined,
        }
        if _JinjaSandboxEnv is None:
            raise ImportError(
                "jinja2 sandbox support is unavailable. Reinstall jinja2 to use from_template()."
            )
        # Autoescape disabled for all extensions and strings to match original behavior.
        # SandboxedEnvironment limits template capabilities when rendering untrusted input.
        env = _JinjaSandboxEnv(
            undefined=undefined_map.get(undefined, jinja2.DebugUndefined),
            autoescape=jinja2.select_autoescape(
                enabled_extensions=(), default=False, default_for_string=False
            ),
        )

        pres = cls(path)
        pres_api = cast("_PresentationTemplateProtocol", pres)

        # Collect all raw text tokens that contain Jinja2 expressions.
        # Multi-paragraph shape text is joined with "\n" in the state, but
        # find_and_replace works on individual XML run text. We therefore split
        # on newlines and process each line (run-level token) separately.
        seen: dict[str, str] = {}  # raw_token -> rendered_token

        for slide_index in range(pres_api.slide_count):
            states: list[dict[str, object]] = cast(
                "list[dict[str, object]]",
                pres_api.execute(
                    _ops.OP_GET_SLIDE_TEXT_STATES, {"slide_index": slide_index}
                ).get("states", []),
            )
            for state in states:
                _collect_jinja_tokens(state, env, context, seen)

        # Apply all replacements in a single pass.
        for raw, rendered in seen.items():
            if raw != rendered:
                pres_api.execute(
                    _ops.OP_FIND_AND_REPLACE, {"find": raw, "replace": rendered}
                )

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
        result = self.execute(_ops.OP_RENDER_TEMPLATE, {"context": context})
        return cast("int", result.get("replacements", 0))

    @classmethod
    def new(cls, title: str) -> Self:
        """Create a new presentation with the given title."""
        pres = cls()
        if cls._lib is None:
            raise GopptxError("Presentation library is not loaded")
        lib = cast("_GopptxLibProtocol", cls._lib)
        handle = lib.deck_new(title.encode("utf-8"))
        if not handle:
            err_ptr = lib.deck_global_error()
            msg = (
                ctypes.string_at(err_ptr).decode("utf-8")
                if err_ptr
                else "Unknown error"
            )
            if err_ptr:
                lib.deck_free_string(err_ptr)
            raise GopptxError(f"Failed to create new deck: {msg}")
        pres._handle = int(handle)
        return pres


__all__ = ["PresentationBase", "PresentationProtocol"]
