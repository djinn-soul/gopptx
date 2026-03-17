"""Presentation text mixin."""

from __future__ import annotations

from typing import cast

from ... import ops
from ..helpers import PresentationMixinBase


class PresentationTextMixin(PresentationMixinBase):
    """Mixin providing text search and replace methods."""

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        """Replace all exact text matches across the presentation."""
        result = self.execute(
            ops.OP_FIND_AND_REPLACE,
            {"find": find_text, "replace": replace_text},
        )
        return int(cast("int", result.get("replacements", 0)))
