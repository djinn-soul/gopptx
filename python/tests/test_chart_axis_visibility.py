"""Axis visibility, tick-label angle, tick crossing and major unit.

Covers issues #473 and #852 (hide an axis), #329 (tick label orientation),
#349 (on/between tick marks) and #472 (major unit on an axis).
"""

import math
import pathlib

from gopptx import Presentation
from gopptx.presentation.charts.chart_types import ChartType

CATEGORIES = ["Alpha", "Beta", "Gamma"]
VALUES = [3.0, 5.0, 2.0]


def _deck_with_chart(
    path: pathlib.Path, chart_type: ChartType = ChartType.COLUMN
) -> None:
    with Presentation.new(title="Axis Controls") as pres:
        pres.slides[0].add_chart(
            chart_type,
            CATEGORIES,
            VALUES,
            title="Axis",
            bounds=(500000, 1200000, 5000000, 3000000),
        )
        pres.save(path)


def test_axis_can_be_hidden_and_shown(tmp_path: pathlib.Path) -> None:
    """An axis hides through <c:delete>, and the axis element survives."""
    output_path = tmp_path / "axis_hidden.pptx"
    _deck_with_chart(output_path)

    with Presentation(output_path) as pres:
        pres.update_chart_formatting(0, {"index": 0}, {"value_axis_visible": False})
        hidden_path = tmp_path / "axis_hidden_saved.pptx"
        pres.save(hidden_path)

    with Presentation(hidden_path) as pres:
        state = pres.get_chart_state_by_index(0, 0)
        assert state["value_axis"]["visible"] is False
        # The axis is still there to be turned back on.
        assert state["value_axis"]["present"] is True

        pres.update_chart_formatting(0, {"index": 0}, {"value_axis_visible": True})
        shown_path = tmp_path / "axis_shown.pptx"
        pres.save(shown_path)

    with Presentation(shown_path) as pres:
        state = pres.get_chart_state_by_index(0, 0)
        assert state["value_axis"]["visible"] is True


def test_tick_label_rotation_round_trips(tmp_path: pathlib.Path) -> None:
    """Long category labels can be angled, and the angle reads back."""
    output_path = tmp_path / "axis_rotation.pptx"
    _deck_with_chart(output_path)

    with Presentation(output_path) as pres:
        pres.update_chart_formatting(
            0, {"index": 0}, {"category_axis_tick_label_rotation": -45.0}
        )
        rotated_path = tmp_path / "axis_rotation_saved.pptx"
        pres.save(rotated_path)

    with Presentation(rotated_path) as pres:
        state = pres.get_chart_state_by_index(0, 0)
        assert math.isclose(state["category_axis"]["tick_label_rotation"], -45.0)


def test_cross_between_puts_the_line_on_the_tick_marks(tmp_path: pathlib.Path) -> None:
    """A line chart can start on the axis instead of between tick marks."""
    output_path = tmp_path / "axis_cross_between.pptx"
    _deck_with_chart(output_path, ChartType.LINE)

    with Presentation(output_path) as pres:
        pres.update_chart_formatting(
            0, {"index": 0}, {"value_axis_cross_between": "midCat"}
        )
        saved_path = tmp_path / "axis_cross_between_saved.pptx"
        pres.save(saved_path)

    with Presentation(saved_path) as pres:
        state = pres.get_chart_state_by_index(0, 0)
        assert state["value_axis"]["cross_between"] == "midCat"


def test_major_unit_round_trips(tmp_path: pathlib.Path) -> None:
    """The major tick interval is settable and readable (issue #472)."""
    output_path = tmp_path / "axis_major_unit.pptx"
    _deck_with_chart(output_path)

    with Presentation(output_path) as pres:
        pres.update_chart_formatting(0, {"index": 0}, {"value_axis_major_unit": 2.5})
        saved_path = tmp_path / "axis_major_unit_saved.pptx"
        pres.save(saved_path)

    with Presentation(saved_path) as pres:
        state = pres.get_chart_state_by_index(0, 0)
        assert math.isclose(state["value_axis"]["major_unit"], 2.5)
