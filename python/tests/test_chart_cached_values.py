"""Refreshing chart cached values without replacing the link (Issue #115)."""

import pathlib

from gopptx import Presentation
from gopptx.presentation.charts.chart_types import ChartType


def _build_chart_deck(path: pathlib.Path) -> None:
    with Presentation.new(title="Chart Cache") as pres:
        slide = pres.slides[0]
        slide.add_chart(
            ChartType.COLUMN,
            ["Q1", "Q2"],
            [1.0, 2.0],
            title="Revenue",
            bounds=(1000000, 1000000, 5000000, 3000000),
        )
        pres.save(path)


def test_cached_values_refresh_without_touching_the_workbook(
    tmp_path: pathlib.Path,
) -> None:
    """Cached values update in place, and the data source is unchanged."""
    output_path = tmp_path / "chart_cache.pptx"
    _build_chart_deck(output_path)

    with Presentation(output_path) as pres:
        before = pres.get_chart_data_source(0, {"index": 0})
        assert before["kind"] in {"embedded", "external"}

        pres.update_chart_cached_values(
            0,
            {"index": 0},
            {
                "categories": ["Q1", "Q2"],
                "series": [{"name": "Revenue", "values": [11.0, 22.0]}],
            },
        )
        refreshed_path = tmp_path / "chart_cache_refreshed.pptx"
        pres.save(refreshed_path)

    with Presentation(refreshed_path) as pres:
        after = pres.get_chart_data_source(0, {"index": 0})
        assert after["kind"] == before["kind"]
        assert after["part_path"] == before.get("part_path", "")

        state = pres.get_chart_state_by_index(0, 0)
        series = state["series"]
        assert series[0]["values"] == [11.0, 22.0]


def test_chart_data_source_reports_the_workbook_part(tmp_path: pathlib.Path) -> None:
    """A chart built by gopptx reports its embedded workbook part."""
    output_path = tmp_path / "chart_source.pptx"
    _build_chart_deck(output_path)

    with Presentation(output_path) as pres:
        source = pres.get_chart_data_source(0, {"index": 0})
        assert source["chart_part"].startswith("ppt/charts/")
        if source["kind"] == "embedded":
            assert source["part_path"].startswith("ppt/embeddings/")
