"""Shared text-frame state/cache for shape text proxies."""
# pyright: reportAttributeAccessIssue=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate, TextRun
    from ..slide import Slide


class ShapeTextFrame:
    """Live text-frame facade for one shape."""

    def __init__(
        self,
        slide: Slide,
        shape_id: int,
    ) -> None:
        """Initialize a cache-backed text-frame facade for one shape."""
        self._slide = slide
        self._shape_id = shape_id
        self._runs_cache: list[TextRun] | None = None
        self._paragraph_cache: dict[str, object] | None = None

    def load_text_state(self) -> None:
        """Load run and paragraph state in one bridge round-trip."""
        state = self._slide.get_shape_text_state(self._shape_id)
        raw_runs = state.get("runs")
        runs: list[TextRun] = []
        if isinstance(raw_runs, list):
            typed_runs = cast("list[object]", raw_runs)
            runs.extend(
                cast("TextRun", dict(cast("dict[str, object]", raw_run)))
                for raw_run in typed_runs
                if isinstance(raw_run, dict)
            )
        self._runs_cache = runs
        raw_paragraph = state.get("paragraph")
        self._paragraph_cache = (
            dict(cast("dict[str, object]", raw_paragraph))
            if isinstance(raw_paragraph, dict)
            else {}
        )

    def get_runs(self) -> list[TextRun]:
        """Return cached runs, loading state on first access."""
        if self._runs_cache is None:
            self.load_text_state()
        return self._runs_cache or []

    def replace_runs(self, runs: list[dict[str, object]]) -> None:
        """Replace all runs and refresh the local cache."""
        payload = [cast("TextRun", dict(run)) for run in runs]
        self._slide.set_shape_runs(self._shape_id, payload)
        self._runs_cache = payload

    def append_run(self, run: dict[str, object]) -> None:
        """Append a run and update the local cache when present."""
        payload = cast("TextRun", dict(run))
        self._slide.append_shape_run(self._shape_id, payload)
        if self._runs_cache is not None:
            self._runs_cache.append(payload)

    def update_run_text(self, run_index: int, value: str) -> None:
        """Update one run's text and mirror the change in cache."""
        self._slide.update_shape_run_text(self._shape_id, run_index, value)
        if self._runs_cache is None:
            return
        if run_index < 0 or run_index >= len(self._runs_cache):
            raise IndexError("run index out of range")
        self._runs_cache[run_index]["text"] = value

    def get_paragraph_payload(self) -> dict[str, object]:
        """Return cached paragraph state, loading text state on first access."""
        if self._paragraph_cache is None:
            self.load_text_state()
        return dict(self._paragraph_cache or {})

    def set_paragraph_field(self, field: str, value: object) -> None:
        """Update one paragraph field and refresh the cache copy."""
        paragraph = self.get_paragraph_payload()
        if value is None:
            paragraph.pop(field, None)
        else:
            paragraph[field] = value
        self._slide.update_shape(
            self._shape_id, cast("ShapeUpdate", {"paragraph": paragraph})
        )
        self._paragraph_cache = paragraph

    def fit_text(self) -> None:
        """Best-effort fit text behavior using bridge-supported controls."""
        self._slide.update_shape(
            self._shape_id,
            cast(
                "ShapeUpdate",
                {"text_frame": {"word_wrap": True, "auto_fit_type": "shape"}},
            ),
        )
