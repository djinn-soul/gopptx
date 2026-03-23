"""Structural protocols for notes-domain proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate


class NotesSlideProto(Protocol):
    @property
    def text(self) -> str: ...

    @text.setter
    def text(self, value: str) -> None: ...

    @property
    def presentation(self) -> object: ...

    @property
    def index(self) -> int: ...

    def shape_payloads(self) -> list[dict[str, object]]: ...

    def _set_shape_props(self, shape_id: int, updates: ShapeUpdate) -> None: ...

    def _set_shape_text(self, shape_id: int, text: str) -> None: ...
