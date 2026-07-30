"""Gridline styling, series lines, series formatting and data-label boxes.

Covers issues #984 (gridline colour and width), #846 (series lines), #872
(series-level marker and line formatting) and #662/#716 (data-label fill and
border).
"""

import pathlib
import zipfile

import pytest
from gopptx import Presentation
from gopptx.presentation.charts.chart_types import ChartType

CATEGORIES = ["Alpha", "Beta", "Gamma"]
VALUES = [3.0, 5.0, 2.0]


def _deck_with_chart(
    path: pathlib.Path, chart_type: ChartType = ChartType.COLUMN
) -> None:
    with Presentation.new(title="Chart Line Styles") as pres:
        pres.slides[0].add_chart(
            chart_type,
            CATEGORIES,
            VALUES,
            title="Styles",
            bounds=(500000, 1200000, 5000000, 3000000),
        )
        pres.save(path)


def test_gridlines_take_a_colour_and_width(tmp_path: pathlib.Path) -> None:
    """A styled gridline keeps its colour, width and dash (issue #984)."""
    source = tmp_path / "gridlines.pptx"
    _deck_with_chart(source)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        chart.value_axis.format_major_gridlines(
            color="D9D9D9", width_emu=9525, dash="sysDot"
        )
        styled = tmp_path / "gridlines_styled.pptx"
        pres.save(styled)

    with Presentation(styled) as pres:
        chart = pres.slides[0].charts[0]
        assert chart.value_axis.major_gridlines_visible is True
        fmt = chart.value_axis.major_gridlines_format
        assert fmt is not None
        assert fmt["color"] == "D9D9D9"
        assert fmt["width_emu"] == 9525
        assert fmt["dash"] == "sysDot"


def test_styling_absent_gridlines_draws_them(tmp_path: pathlib.Path) -> None:
    """Styling a gridline the axis lacks turns it on (issue #984)."""
    source = tmp_path / "minor_gridlines.pptx"
    _deck_with_chart(source)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        chart.value_axis.format_minor_gridlines(color="EEEEEE")
        styled = tmp_path / "minor_gridlines_styled.pptx"
        pres.save(styled)

    with Presentation(styled) as pres:
        chart = pres.slides[0].charts[0]
        assert chart.value_axis.minor_gridlines_visible is True
        fmt = chart.value_axis.minor_gridlines_format
        assert fmt is not None
        assert fmt["color"] == "EEEEEE"


def test_gridline_format_rejects_a_bad_colour(tmp_path: pathlib.Path) -> None:
    """A malformed colour is refused before it reaches the engine."""
    source = tmp_path / "gridline_guard.pptx"
    _deck_with_chart(source)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        with pytest.raises(ValueError, match="hex colour"):
            chart.value_axis.format_major_gridlines(color="not-a-colour")


def test_series_lines_are_drawn_and_removed(tmp_path: pathlib.Path) -> None:
    """A stacked bar chart draws and drops c:serLines (issue #846)."""
    source = tmp_path / "series_lines.pptx"
    _deck_with_chart(source, ChartType.BAR_STACKED)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        chart.set_series_lines(show=True, color="404040", width_emu=12700)
        drawn = tmp_path / "series_lines_drawn.pptx"
        pres.save(drawn)

    with Presentation(drawn) as pres:
        chart = pres.slides[0].charts[0]
        lines = chart.series_lines()
        assert lines is not None
        assert lines["show"] is True
        assert lines["line"]["color"] == "404040"
        chart.set_series_lines(show=False)
        cleared = tmp_path / "series_lines_cleared.pptx"
        pres.save(cleared)

    with Presentation(cleared) as pres:
        assert pres.slides[0].charts[0].series_lines() is None


def test_series_markers_take_a_colour(tmp_path: pathlib.Path) -> None:
    """A line series keeps its marker shape when recoloured (issue #872)."""
    source = tmp_path / "series_format.pptx"
    _deck_with_chart(source, ChartType.LINE_MARKERS)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        chart.format_series(
            0,
            line_color="1F77B4",
            line_width_emu=19050,
            marker_symbol="circle",
            marker_size=7,
            marker_fill_color="FFFFFF",
            marker_line_color="1F77B4",
        )
        styled = tmp_path / "series_format_styled.pptx"
        pres.save(styled)

    with Presentation(styled) as pres:
        chart = pres.slides[0].charts[0]
        fmt = chart.series_format(0)
        assert fmt is not None
        assert fmt["line_color"] == "1F77B4"
        assert fmt["marker_symbol"] == "circle"
        assert fmt["marker_size"] == 7
        # Recolouring only the marker keeps the symbol and size already set.
        chart.format_series(0, marker_fill_color="FFF2CC")
        recoloured = tmp_path / "series_format_recoloured.pptx"
        pres.save(recoloured)

    with Presentation(recoloured) as pres:
        fmt = pres.slides[0].charts[0].series_format(0)
        assert fmt is not None
        assert fmt["marker_symbol"] == "circle"
        assert fmt["marker_size"] == 7


def test_series_format_rejects_an_unknown_marker(tmp_path: pathlib.Path) -> None:
    """A symbol outside CT_MarkerStyle is refused."""
    source = tmp_path / "marker_guard.pptx"
    _deck_with_chart(source, ChartType.LINE_MARKERS)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        with pytest.raises(ValueError, match="marker_symbol"):
            chart.format_series(0, marker_symbol="blob")


def test_data_labels_take_a_fill_and_border(tmp_path: pathlib.Path) -> None:
    """Data labels carry a background and an outline (issues #662, #716)."""
    source = tmp_path / "label_box.pptx"
    _deck_with_chart(source)

    with Presentation(source) as pres:
        chart = pres.slides[0].charts[0]
        labels = chart.plots[0].data_labels
        labels.set_fill("FFF2CC")
        labels.set_border(color="BF8F00", width_emu=9525, dash="dash")
        styled = tmp_path / "label_box_styled.pptx"
        pres.save(styled)

    with zipfile.ZipFile(styled) as archive:
        part = archive.read("ppt/charts/chart1.xml").decode("utf-8")
    assert '<a:srgbClr val="FFF2CC"/>' in part
    assert '<a:srgbClr val="BF8F00"/>' in part
    assert '<a:prstDash val="dash"/>' in part
