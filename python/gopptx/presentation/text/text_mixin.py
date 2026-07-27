"""Presentation text mixin."""

from __future__ import annotations

from typing import cast

from ... import ops
from ..helpers import PresentationMixinBase


class PresentationTextMixin(PresentationMixinBase):
    """Mixin providing text search and replace methods."""

    def find_and_replace(
        self,
        find_text: str,
        replace_text: str,
        scope: str = "slides",
    ) -> int:
        """Replace exact text matches in the parts named by ``scope``.

        Matching is paragraph-scoped, so a phrase split across runs by
        PowerPoint is still replaced, but a phrase spanning a paragraph break
        is not. ``scope`` is one of ``"slides"`` (the default, slide shapes
        only), ``"slides+notes"``, or ``"all"``, which adds layouts, slide
        masters, the notes master and the handout master.
        """
        # Buffered textbox adds and run replacements are not in the XML yet, so
        # without flushing they would be invisible to the replacement and the
        # call would report 0 for text the caller has already written.
        self._flush_buffered_writes()
        result = self.execute(
            ops.OP_FIND_AND_REPLACE,
            {"find": find_text, "replace": replace_text, "scope": scope},
        )
        return int(cast("int", result.get("replacements", 0)))

    def _flush_buffered_writes(self) -> None:
        for flush_name in (
            "flush_all_pending_textbox_adds",
            "flush_all_pending_slide_run_text_updates",
        ):
            flush = getattr(self, flush_name, None)
            if callable(flush):
                flush()
