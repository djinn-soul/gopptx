from __future__ import annotations

from dataclasses import dataclass
from typing import Any

import pytest
from gopptx import ops
from gopptx.builder import PresentationBuilder
from gopptx.constants import ShapeType
from gopptx.shapes.shape_builder import ArrowConfig, ShapeBuilder
from gopptx.templates import (
    Milestone,
    PricingTier,
    ProposalTemplate,
    SimpleTemplate,
    StatusTemplate,
    TechnicalTemplate,
    TrainingTemplate,
    apply_slides,
)
from gopptx.text.run_builder import RunBuilder

# Non-password-named constant so ruff S105 does not trigger on the value.
_DOC_LOCK_VALUE = "doc-protect-test-xyz"


@dataclass
class _FakePresentation:
    title: str
    executed: list[tuple[str, dict[str, object]]]
    slides_added: list[tuple[str, str | None, list[str] | None]]
    core_properties: dict[str, object] | None = None
    theme: object | None = None
    slide_size: tuple[int, int] | None = None
    modify_password: str = ""
    final: bool = False
    removed: list[int] | None = None
    notes_set: list[tuple[int, str]] | None = None
    tables_added: (
        list[tuple[int, list[list[str]], tuple[int, int, int, int], bool, bool]] | None
    ) = None
    closed: bool = False

    def set_core_properties(self, props: dict[str, object]) -> None:
        self.core_properties = props

    def apply_theme(self, theme: object) -> None:
        self.theme = theme

    def set_slide_size(self, width: int, height: int) -> None:
        self.slide_size = (width, height)

    def add_slide(
        self,
        title: str = "",
        *,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None:
        self.slides_added.append((title, layout, bullets))

    def set_modify_password(self, password: str) -> None:
        self.modify_password = password

    def set_mark_as_final(self, *, final: bool) -> None:
        self.final = final

    def save(self, _path: str) -> None:
        return None

    def to_bytes(self) -> bytes:
        return b"fake-pptx"

    def remove_slide(self, index: int) -> None:
        if self.removed is None:
            self.removed = []
        self.removed.append(index)

    def set_notes(self, slide_index: int, notes: str) -> None:
        if self.notes_set is None:
            self.notes_set = []
        self.notes_set.append((slide_index, notes))

    def add_table_from_rows(
        self,
        slide_index: int,
        rows: list[list[str]],
        bounds: tuple[int, int, int, int],
        *,
        first_row: bool,
        band_row: bool,
    ) -> None:
        if self.tables_added is None:
            self.tables_added = []
        self.tables_added.append((slide_index, rows, bounds, first_row, band_row))

    def execute(self, op_name: str, payload: dict[str, object]) -> dict[str, object]:
        self.executed.append((op_name, payload))
        return {
            "slides": [
                {"title": "Generated", "bullets": ["a", "b"]},
                {
                    "title": "With notes+table",
                    "notes": "n",
                    "table": {"rows": [["x"], ["y"]]},
                },
            ]
        }

    def close(self) -> None:
        self.closed = True


@pytest.fixture
def fake_presentation_factory(monkeypatch: pytest.MonkeyPatch):
    created: list[_FakePresentation] = []

    def _new(title: str) -> _FakePresentation:
        prs = _FakePresentation(title=title, executed=[], slides_added=[])
        created.append(prs)
        return prs

    monkeypatch.setattr(
        "gopptx.presentation.presentation.Presentation.new",
        staticmethod(_new),
    )
    return created


def test_presentation_builder_build_save_and_bytes(
    fake_presentation_factory: list[_FakePresentation],
) -> None:
    builder = (
        PresentationBuilder("Quarterly Report")
        .with_author("Alice")
        .with_subject("Q1")
        .with_keywords("finance,quarterly")
        .with_description("Summary")
        .with_theme("corporate")
        .with_slide_size(13.33, 7.5)
        .with_modify_password(_DOC_LOCK_VALUE)
        .with_mark_as_final()
        .add_title_slide("Title", layout="title")
        .add_bullet_slide("Highlights", ["A", "B"], layout="content")
    )
    prs = builder.build()

    assert prs.title == "Quarterly Report"
    assert prs.core_properties == {
        "title": "Quarterly Report",
        "creator": "Alice",
        "subject": "Q1",
        "keywords": "finance,quarterly",
        "description": "Summary",
    }
    assert prs.theme == "corporate"
    assert prs.slide_size == (12188952, 6858000)
    assert prs.modify_password == _DOC_LOCK_VALUE
    assert prs.final is True
    assert prs.slides_added == [
        ("Title", "title", []),
        ("Highlights", "content", ["A", "B"]),
    ]

    builder.save("out.pptx")
    assert builder.to_bytes() == b"fake-pptx"
    assert len(fake_presentation_factory) == 3


def test_presentation_builder_without_optional_properties(
    fake_presentation_factory: list[_FakePresentation],
) -> None:
    prs = PresentationBuilder("Only Title").build()
    assert prs.core_properties is None
    assert prs.slide_size is None
    assert not prs.modify_password
    assert prs.final is False


def test_shape_builder_factories_properties_and_repr() -> None:
    builder = (
        ShapeBuilder
        .rectangle(1.0, 2.0, 3.0, 4.0)
        .with_text("Hello")
        .with_fill("FF0000")
        .with_line(
            "00FF00",
            width_emu=12700,
            dash_style="dash",
            start_arrow=ArrowConfig("triangle", width="sm", length="lg"),
            end_arrow=ArrowConfig("oval", width="med", length="sm"),
        )
        .with_shadow(color="111111", blur_emu=100, distance_emu=200, angle_deg=30)
        .with_rotation(12.5)
        .flip_horizontal()
        .flip_vertical()
    )
    assert builder.shape_type == ShapeType.RECTANGLE
    assert builder.bounds == (914400.0, 1828800.0, 2743200.0, 3657600.0)
    kwargs = builder.to_kwargs()
    assert kwargs["text"] == "Hello"
    assert kwargs["properties"]["line"]["end_arrow"] == "oval"
    assert kwargs["properties"]["shadow"]["distance_emu"] == 200
    assert "ShapeBuilder" in repr(builder)

    assert ShapeBuilder.of(ShapeType.HEART, 0, 0, 1, 1).shape_type == ShapeType.HEART
    assert (
        ShapeBuilder.rounded_rectangle(0, 0, 1, 1).shape_type
        == ShapeType.ROUNDED_RECTANGLE
    )
    assert ShapeBuilder.ellipse(0, 0, 1, 1).shape_type == ShapeType.ELLIPSE
    assert ShapeBuilder.triangle(0, 0, 1, 1).shape_type == ShapeType.TRIANGLE
    assert (
        ShapeBuilder.right_triangle(0, 0, 1, 1).shape_type == ShapeType.RIGHT_TRIANGLE
    )
    assert ShapeBuilder.diamond(0, 0, 1, 1).shape_type == ShapeType.DIAMOND
    assert ShapeBuilder.pentagon(0, 0, 1, 1).shape_type == ShapeType.PENTAGON
    assert ShapeBuilder.hexagon(0, 0, 1, 1).shape_type == ShapeType.HEXAGON
    assert ShapeBuilder.parallelogram(0, 0, 1, 1).shape_type == ShapeType.PARALLELOGRAM
    assert ShapeBuilder.cloud(0, 0, 1, 1).shape_type == ShapeType.CLOUD
    assert ShapeBuilder.heart(0, 0, 1, 1).shape_type == ShapeType.HEART
    assert ShapeBuilder.star5(0, 0, 1, 1).shape_type == ShapeType.STAR_5
    assert ShapeBuilder.star6(0, 0, 1, 1).shape_type == ShapeType.STAR_6

    assert ShapeBuilder.rectangle(0, 0, 1, 1).with_no_fill().to_kwargs()["properties"][
        "fill"
    ] == {"background": True}
    assert ShapeBuilder.rectangle(0, 0, 1, 1).with_no_line().to_kwargs()["properties"][
        "line"
    ] == {
        "color": "FFFFFF",
        "width_emu": 0,
    }


def test_run_builder_payload_and_build() -> None:
    rb = (
        RunBuilder("hello")
        .text("world")
        .bold()
        .italic(value=False)
        .underline("dbl")
        .strikethrough("sng")
        .subscript()
        .superscript(value=False)
        .color("FF00FF")
        .highlight("FFFF00")
        .font("Calibri")
        .size_pt(22)
        .code()
        .all_caps()
        .small_caps(value=False)
        .hyperlink("https://example.com", tooltip="tip")
        .hover_action("ppaction://hlinkshowjump?jump=nextslide")
    )
    payload = rb.to_payload()
    assert payload["text"] == "world"
    assert payload["bold"] is True
    assert payload["italic"] is False
    assert payload["hyperlink"] == {"address": "https://example.com", "tooltip": "tip"}
    assert payload["hover_action"]["address"].startswith("ppaction://")
    assert rb.build()["font"] == "Calibri"
    assert repr(rb) == "RunBuilder('world')"
    assert RunBuilder("x").hyperlink("https://x").to_payload()["hyperlink"] == {
        "address": "https://x"
    }


def test_apply_slides_adds_slides_notes_and_tables() -> None:
    prs = _FakePresentation(title="t", executed=[], slides_added=[])
    apply_slides(
        prs,
        [
            {"title": "Intro", "layout": "Title Slide", "bullets": ["a", "b"]},
            {
                "title": "Data",
                "notes": "speaker",
                "table": {
                    "rows": [["H1", "H2"], ["1", "2"]],
                    "x": 1,
                    "y": 2,
                    "cx": 3,
                    "cy": 4,
                },
            },
        ],
    )
    assert prs.removed == [0]
    assert prs.slides_added == [
        ("Intro", "Title Slide", ["a", "b"]),
        ("Data", None, None),
    ]
    assert prs.notes_set == [(1, "speaker")]
    assert prs.tables_added == [
        (1, [["H1", "H2"], ["1", "2"]], (1, 2, 3, 4), True, True)
    ]

    prs2 = _FakePresentation(title="t2", executed=[], slides_added=[])
    apply_slides(prs2, "not-a-list")
    assert prs2.removed == [0]
    assert prs2.slides_added == []


def test_template_builders_execute_and_apply_theme(
    fake_presentation_factory: list[_FakePresentation],
) -> None:
    templates: list[tuple[str, Any]] = [
        (
            ops.OP_BUILD_STATUS_TEMPLATE,
            StatusTemplate(
                project="Apollo",
                okrs=["A"],
                risks=["B"],
                next_steps=["C"],
                theme="modern",
            ),
        ),
        (
            ops.OP_BUILD_SIMPLE_TEMPLATE,
            SimpleTemplate(title="Simple", content="Body", theme="dark"),
        ),
        (
            ops.OP_BUILD_PROPOSAL_TEMPLATE,
            ProposalTemplate(
                title="Proposal",
                subtitle="Sub",
                context="Ctx",
                solution="Sol",
                pricing=[PricingTier("Basic", "$10", ["f1"])],
                timeline=[Milestone("2026-01-01", "Kickoff", "done")],
                theme="corporate",
            ),
        ),
        (
            ops.OP_BUILD_TRAINING_TEMPLATE,
            TrainingTemplate(
                title="Training",
                agenda=["A1"],
                concepts=["C1"],
                summary="S",
                theme="minimal",
            ),
        ),
        (
            ops.OP_BUILD_TECHNICAL_TEMPLATE,
            TechnicalTemplate(
                title="Tech",
                architecture="Arch",
                deep_dive="Deep",
                benchmarks="Bench",
                theme="clean",
            ),
        ),
    ]

    for op_name, template in templates:
        prs = template.build()
        assert prs.executed and prs.executed[0][0] == op_name
        assert prs.theme is not None
        assert len(prs.slides_added) == 2

    assert len(fake_presentation_factory) == 5


@pytest.mark.parametrize(
    ("template", "message"),
    [
        (StatusTemplate(project=""), "project name cannot be empty"),
        (SimpleTemplate(title=""), "title cannot be empty"),
        (ProposalTemplate(title=""), "title cannot be empty"),
        (TrainingTemplate(title=""), "title cannot be empty"),
        (TechnicalTemplate(title=""), "title cannot be empty"),
    ],
)
def test_template_builders_validate_empty_title_or_project(
    template: Any, message: str
) -> None:
    with pytest.raises(ValueError, match=message):
        template.build()


def test_template_builders_close_on_execute_failure(
    monkeypatch: pytest.MonkeyPatch,
    fake_presentation_factory: list[_FakePresentation],
) -> None:
    _ = fake_presentation_factory

    def _execute_raises(
        self: _FakePresentation, _op: str, _payload: dict[str, object]
    ) -> dict[str, object]:
        raise RuntimeError("boom")

    monkeypatch.setattr(_FakePresentation, "execute", _execute_raises)
    with pytest.raises(RuntimeError, match="boom"):
        SimpleTemplate(title="X").build()

    assert fake_presentation_factory[-1].closed is True
