"""Shape text-run operations for the presentation facade."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...slide.text.text_run import serialize_runs_for_payload
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import TextRun


class PresentationShapeTextRunMixin(PresentationMixinBase):
    """Methods for reading and mutating text-run state on shape text."""

    def get_shape_text_state(
        self, slide_index: int, shape_id: int
    ) -> dict[str, object]:
        """Get text/runs/text-frame/paragraph state for a shape."""
        return self.execute(
            ops.OP_GET_SHAPE_TEXT_STATE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )

    def get_slide_text_states(self, slide_index: int) -> list[dict[str, object]]:
        """Get text/runs/text-frame/paragraph state for all shapes on a slide."""
        result = self.execute(
            ops.OP_GET_SLIDE_TEXT_STATES,
            {"slide_index": slide_index},
        )
        return cast("list[dict[str, object]]", result.get("states", []))

    def get_shape_runs(self, slide_index: int, shape_id: int) -> list[TextRun]:
        """Get text runs for a shape."""
        result = self.execute(
            ops.OP_GET_SHAPE_RUNS,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("list[TextRun]", result.get("runs", []))

    def set_shape_runs(
        self, slide_index: int, shape_id: int, runs: list[TextRun]
    ) -> None:
        """Replace all text runs on a shape."""
        self.execute(
            ops.OP_SET_SHAPE_RUNS,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "runs": serialize_runs_for_payload(cast("object", runs)),
            },
        )

    def update_shape_run_text(
        self,
        slide_index: int,
        shape_id: int,
        run_index: int,
        text: str,
    ) -> None:
        """Update text for one run by run index."""
        self.execute(
            ops.OP_UPDATE_SHAPE_RUN_TEXT,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "run_index": run_index,
                "text": text,
            },
        )

    def append_shape_run(
        self,
        slide_index: int,
        shape_id: int,
        run: TextRun,
    ) -> None:
        """Append a run to a shape."""
        payload = serialize_runs_for_payload([cast("object", run)])
        run_payload = cast("dict[str, object]", cast("list[object]", payload)[0])
        self.execute(
            ops.OP_APPEND_SHAPE_RUN,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "run": run_payload,
            },
        )
