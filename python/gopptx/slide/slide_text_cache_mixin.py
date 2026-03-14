"""Helpers for maintaining slide-local text caches during buffered mutations."""
# pyright: reportUninitializedInstanceVariable=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from collections.abc import Sequence

    from ..schemas import Shape


class SlideTextCacheMixin:
    """Update slide-local text and shape caches without forcing bridge flushes."""

    if TYPE_CHECKING:
        _shape_records_cache: list[Shape] | None
        _shape_record_map: dict[int, Shape] | None
        _shape_text_state_cache: dict[int, dict[str, object]] | None

    def _apply_cached_run_text_update(
        self,
        shape_id: int,
        run_index: int,
        text: str,
    ) -> None:
        state = (
            self._shape_text_state_cache.get(shape_id)
            if self._shape_text_state_cache
            else None
        )
        if state is None:
            return
        raw_runs = state.get("runs")
        if not isinstance(raw_runs, list):
            return
        runs = cast("list[object]", raw_runs)
        if run_index < 0 or run_index >= len(runs):
            return
        run_payload = runs[run_index]
        if not isinstance(run_payload, dict):
            return
        cast("dict[str, object]", run_payload)["text"] = text
        self._patch_cached_shape_record_text(shape_id, self._joined_run_text(runs))

    def _replace_cached_runs(
        self,
        shape_id: int,
        runs: list[dict[str, object]],
    ) -> None:
        runs_copy = [dict(run) for run in runs]
        if (
            self._shape_text_state_cache is not None
            and shape_id in self._shape_text_state_cache
        ):
            self._shape_text_state_cache[shape_id]["runs"] = runs_copy
        self._patch_cached_shape_record_text(shape_id, self._joined_run_text(runs_copy))

    def _append_cached_run(
        self,
        shape_id: int,
        run: dict[str, object],
    ) -> None:
        state = (
            self._shape_text_state_cache.get(shape_id)
            if self._shape_text_state_cache
            else None
        )
        if state is not None:
            raw_runs = state.setdefault("runs", [])
            if isinstance(raw_runs, list):
                runs = cast("list[object]", raw_runs)
                runs.append(dict(run))
                self._patch_cached_shape_record_text(
                    shape_id, self._joined_run_text(runs)
                )
                return
        record_text = str(run.get("text", ""))
        if record_text:
            self._patch_cached_shape_record_text(
                shape_id,
                f"{self._cached_shape_record_text(shape_id)}{record_text}",
            )

    def _cached_shape_record_text(self, shape_id: int) -> str:
        record = (
            self._shape_record_map.get(shape_id) if self._shape_record_map else None
        )
        if record is None and self._shape_records_cache is not None:
            for shape in self._shape_records_cache:
                shape_payload = cast("dict[str, object]", shape)
                raw_id = shape_payload.get("ID", shape_payload.get("id"))
                if raw_id is not None and int(str(raw_id)) == shape_id:
                    record = shape
                    break
        if record is None:
            return ""
        record_payload = cast("dict[str, object]", record)
        value = record_payload.get("Text", record_payload.get("text", ""))
        return str(value)

    def _patch_cached_shape_record_text(self, shape_id: int, text: str) -> None:
        if self._shape_record_map and shape_id in self._shape_record_map:
            record_payload = cast("dict[str, object]", self._shape_record_map[shape_id])
            record_payload["Text"] = text
            record_payload["text"] = text
        if self._shape_records_cache is None:
            return
        for shape in self._shape_records_cache:
            shape_payload = cast("dict[str, object]", shape)
            raw_id = shape_payload.get("ID", shape_payload.get("id"))
            if raw_id is None or int(str(raw_id)) != shape_id:
                continue
            shape_payload["Text"] = text
            shape_payload["text"] = text
            return

    @staticmethod
    def _joined_run_text(runs: Sequence[object]) -> str:
        return "".join(
            str(cast("dict[str, object]", run).get("text", ""))
            for run in runs
            if isinstance(run, dict)
        )
