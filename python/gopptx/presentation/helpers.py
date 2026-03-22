"""Shared helpers and protocol contracts for presentation APIs."""

from __future__ import annotations

import json
import pathlib
import re
from typing import TYPE_CHECKING, cast

from typing_extensions import Protocol

if TYPE_CHECKING:
    from ..schemas import BatchItemResult

try:
    import orjson as _orjson
except ImportError:  # pragma: no cover - optional dependency
    _orjson = None


class PresentationProtocol(Protocol):
    """Typed contract for presentation runtime methods."""

    _handle: int | None
    _comment_ref_cache: dict[int, tuple[int, int, int]]

    @property
    def handle(self) -> int | None:
        """Return the opaque integer handle for this presentation instance."""
        ...

    def execute(
        self, op: str, payload: dict[str, object] | None = None
    ) -> dict[str, object]:
        """Execute a bridge operation and return the decoded response payload."""
        ...

    def invalidate_cache(self) -> None:
        """Invalidate cached presentation-derived state."""
        ...

    def begin_batch(self, *, stop_on_error: bool = False) -> None:
        """Start buffering operations into a batch."""
        ...

    def end_batch(self) -> list[BatchItemResult]:
        """Flush the active batch and return per-item execution results."""
        ...

    def abort_batch(self) -> None:
        """Discard buffered batch operations without executing them."""
        ...


if TYPE_CHECKING:

    class PresentationMixinBase(PresentationProtocol, Protocol):
        """Typed contract for mixins without altering runtime MRO behavior."""

else:

    class PresentationMixinBase:
        """Runtime marker base for presentation mixins."""


def json_dumps(payload: dict[str, object]) -> bytes:
    """Serialize a dictionary to JSON bytes."""

    def default(obj: object) -> str:
        if isinstance(obj, pathlib.Path):
            return str(obj)
        raise TypeError(f"Type {type(obj)} not serializable")

    if _orjson is not None:
        return _orjson.dumps(payload, default=default)
    return json.dumps(payload, separators=(",", ":"), default=default).encode("utf-8")


def json_loads(raw: bytes) -> object:
    """Deserialize JSON bytes to Python objects."""
    if _orjson is not None:
        return cast("object", _orjson.loads(raw))
    return cast("object", json.loads(raw.decode("utf-8")))


def snake_case(name: str) -> str:
    """Convert CamelCase to snake_case."""
    s1 = re.sub(r"(.)([A-Z][a-z]+)", r"\1_\2", name)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", s1).lower()


def get_required_int(result: dict[str, object], key: str) -> int:
    """Get an integer value from the result dict, raising TypeError if it's not an int."""
    value = result.get(key)
    if not isinstance(value, int):
        raise TypeError(f"bridge response {key} must be an int")
    return value


def with_key_aliases(obj: object) -> object:
    """Add lowercase and snake_case aliases for nested dictionary keys."""
    if isinstance(obj, list):
        list_obj = cast("list[object]", obj)
        return [with_key_aliases(item) for item in list_obj]
    if not isinstance(obj, dict):
        return obj

    dict_obj = cast("dict[str, object]", obj)
    out: dict[str, object] = {}
    for key, value in dict_obj.items():
        transformed = with_key_aliases(value)
        out[key] = transformed
        lower_key = str(key).lower()
        if lower_key not in out:
            out[lower_key] = transformed
        snake_key = snake_case(str(key))
        if snake_key not in out:
            out[snake_key] = transformed
    return out
