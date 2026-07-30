"""Shared text-frame state/cache for shape text proxies."""
# pyright: reportAttributeAccessIssue=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import ShapeTextParagraph, ShapeUpdate, TextRun


class _ShapeTextSlideProto(Protocol):
    def get_shape_text_state(self, shape_id: int) -> dict[str, object]: ...

    def set_shape_runs(self, shape_id: int, runs: list[TextRun]) -> None: ...

    def append_shape_run(self, shape_id: int, run: TextRun) -> None: ...

    def update_shape_run_text(
        self, shape_id: int, run_index: int, text: str
    ) -> None: ...

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None: ...


class ShapeTextFrame:
    """Live text-frame facade for one shape."""

    def __init__(
        self,
        slide: _ShapeTextSlideProto,
        shape_id: int,
    ) -> None:
        """Initialize a cache-backed text-frame facade for one shape."""
        self._slide = slide
        self._shape_id = shape_id
        self._paragraphs_cache: list[ShapeTextParagraph] | None = None

    def load_text_state(self) -> None:
        """Load all paragraph state in one bridge round-trip."""
        state = self._slide.get_shape_text_state(self._shape_id)
        raw_paragraphs = state.get("paragraphs")
        paragraphs: list[ShapeTextParagraph] = []
        if isinstance(raw_paragraphs, list):
            for raw_paragraph in cast("list[object]", raw_paragraphs):
                if not isinstance(raw_paragraph, dict):
                    continue
                paragraph = cast("dict[str, object]", raw_paragraph)
                raw_runs = paragraph.get("runs", [])
                runs = [
                    cast("TextRun", dict(cast("dict[str, object]", raw_run)))
                    for raw_run in cast("list[object]", raw_runs)
                    if isinstance(raw_run, dict)
                ]
                raw_props = paragraph.get("paragraph")
                props = (
                    dict(cast("dict[str, object]", raw_props))
                    if isinstance(raw_props, dict)
                    else {}
                )
                paragraphs.append(
                    cast("ShapeTextParagraph", {"runs": runs, "paragraph": props})
                )
        if not paragraphs:
            paragraphs.append(self._legacy_paragraph_from_state(state))
        self._paragraphs_cache = paragraphs

    @staticmethod
    def _legacy_paragraph_from_state(
        state: dict[str, object],
    ) -> ShapeTextParagraph:
        raw_runs = state.get("runs", [])
        runs = [
            cast("TextRun", dict(cast("dict[str, object]", raw_run)))
            for raw_run in cast("list[object]", raw_runs)
            if isinstance(raw_run, dict)
        ]
        raw_props = state.get("paragraph")
        props = (
            dict(cast("dict[str, object]", raw_props))
            if isinstance(raw_props, dict)
            else {}
        )
        return cast("ShapeTextParagraph", {"runs": runs, "paragraph": props})

    def get_paragraphs(self) -> list[dict[str, object]]:
        """Return cached paragraph payloads, loading on first access."""
        if self._paragraphs_cache is None:
            self.load_text_state()
        return [dict(paragraph) for paragraph in self._paragraphs_cache or []]

    def replace_paragraphs(self, paragraphs: list[dict[str, object]]) -> None:
        """Replace all paragraphs and refresh the local cache."""
        payload = [
            cast(
                "ShapeTextParagraph",
                {
                    "runs": [
                        cast("TextRun", dict(cast("dict[str, object]", run)))
                        for run in cast("list[object]", paragraph.get("runs", []))
                        if isinstance(run, dict)
                    ],
                    "paragraph": dict(
                        cast("dict[str, object]", paragraph.get("paragraph", {}))
                    ),
                },
            )
            for paragraph in paragraphs
        ]
        self._slide.update_shape(
            self._shape_id, cast("ShapeUpdate", {"paragraphs": payload})
        )
        self._paragraphs_cache = payload

    def get_paragraph_runs(self, paragraph_index: int) -> list[dict[str, object]]:
        """Return runs for one paragraph."""
        paragraph = self._paragraph_at(paragraph_index)
        return [
            dict(cast("dict[str, object]", run))
            for run in cast("list[object]", paragraph.get("runs", []))
            if isinstance(run, dict)
        ]

    def replace_paragraph_runs(
        self, paragraph_index: int, runs: list[dict[str, object]]
    ) -> None:
        """Replace runs for one paragraph."""
        paragraphs = self.get_paragraphs()
        self._paragraph_at(paragraph_index)
        paragraphs[paragraph_index]["runs"] = [dict(run) for run in runs]
        self.replace_paragraphs(paragraphs)

    def get_paragraph_payload(self, paragraph_index: int) -> dict[str, object]:
        """Return formatting state for one paragraph."""
        paragraph = self._paragraph_at(paragraph_index)
        raw = paragraph.get("paragraph", {})
        return dict(cast("dict[str, object]", raw))

    def set_paragraph_field(
        self, paragraph_index: int, field: str, value: object
    ) -> None:
        """Update one paragraph field and refresh the cache copy."""
        paragraphs = self.get_paragraphs()
        paragraph = self.get_paragraph_payload(paragraph_index)
        if value is None:
            paragraph.pop(field, None)
        else:
            paragraph[field] = value
        paragraphs[paragraph_index]["paragraph"] = paragraph
        self.replace_paragraphs(paragraphs)

    def _paragraph_at(self, paragraph_index: int) -> ShapeTextParagraph:
        if self._paragraphs_cache is None:
            self.load_text_state()
        paragraphs = self._paragraphs_cache or []
        if paragraph_index < 0 or paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        return paragraphs[paragraph_index]

    def fit_text(self) -> None:
        """Best-effort fit text behavior using bridge-supported controls."""
        self._slide.update_shape(
            self._shape_id,
            cast(
                "ShapeUpdate",
                {"text_frame": {"word_wrap": True, "auto_fit_type": "shape"}},
            ),
        )
