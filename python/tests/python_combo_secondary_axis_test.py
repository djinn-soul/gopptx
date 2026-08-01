"""Combo charts can put their line series on a secondary value axis (Issue #141)."""

import re
import zipfile

from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches

CATEGORIES = ["Q1", "Q2", "Q3", "Q4"]
# Revenue against a growth percentage: two magnitudes that flatten the line
# against zero when both series share one value axis.
BAR_SERIES = [{"name": "Revenue (k)", "values": [820.0, 932.0, 901.0, 1290.0]}]
LINE_SERIES = [{"name": "Growth %", "values": [4.5, 13.7, -3.3, 43.2]}]
LINE_SECONDARY = [dict(LINE_SERIES[0], secondary_axis=True)]
LINE_TITLED = [dict(LINE_SERIES[0], secondary_axis_title="Growth %")]
BOUNDS = (Inches(0.5), Inches(1.3), Inches(8.9), Inches(5.2))


def _chart_xml(tmp_path, name, line_series=None):
    path = tmp_path / name
    with Presentation.new("Combo Secondary Axis") as pres:
        pres.slides[0].add_combo_chart(
            CATEGORIES,
            BAR_SERIES,
            line_series if line_series is not None else LINE_SERIES,
            title="Revenue vs Growth",
            bounds=BOUNDS,
        )
        pres.save(path)
    with zipfile.ZipFile(path) as archive:
        return archive.read("ppt/charts/chart1.xml").decode("utf-8")


def _plot_axis_ids(xml, plot):
    block = re.search(rf"<c:{plot}>.*?</c:{plot}>", xml, re.DOTALL)
    assert block is not None, f"no <c:{plot}> in chart part"
    return re.findall(r'<c:axId val="(\d+)"/>', block.group(0))


def test_combo_chart_shares_one_axis_pair_by_default(tmp_path):
    xml = _chart_xml(tmp_path, "shared.pptx")
    assert _plot_axis_ids(xml, "barChart") == _plot_axis_ids(xml, "lineChart")
    assert xml.count("<c:valAx>") == 1


def test_secondary_axis_binds_the_line_plot_to_its_own_pair(tmp_path):
    xml = _chart_xml(tmp_path, "secondary.pptx", LINE_SECONDARY)

    bar_ids = _plot_axis_ids(xml, "barChart")
    line_ids = _plot_axis_ids(xml, "lineChart")
    assert len(bar_ids) == 2
    assert len(line_ids) == 2
    assert set(bar_ids).isdisjoint(line_ids)
    assert xml.count("<c:valAx>") == 2
    assert xml.count("<c:catAx>") == 2


def test_every_referenced_axis_is_declared(tmp_path):
    xml = _chart_xml(tmp_path, "declared.pptx", LINE_SECONDARY)

    without_plots = re.sub(
        r"<c:(bar|line)Chart>.*?</c:\1Chart>", "", xml, flags=re.DOTALL
    )
    declared = set(re.findall(r'<c:axId val="(\d+)"/>', without_plots))
    referenced = set(_plot_axis_ids(xml, "barChart") + _plot_axis_ids(xml, "lineChart"))
    assert referenced <= declared
    assert len(referenced) == 4


def test_secondary_axis_is_drawn_opposite_the_primary(tmp_path):
    xml = _chart_xml(tmp_path, "geometry.pptx", LINE_TITLED)

    primary, secondary = re.findall(r"<c:valAx>.*?</c:valAx>", xml, re.DOTALL)
    assert '<c:axPos val="l"/>' in primary
    assert '<c:axPos val="r"/>' in secondary
    assert '<c:crosses val="max"/>' in secondary
    assert "Growth %" in secondary
    assert "Growth %" not in primary

    # The duplicate category axis must be hidden or the chart shows two
    # identical category scales.
    category_axes = re.findall(r"<c:catAx>.*?</c:catAx>", xml, re.DOTALL)
    assert '<c:delete val="0"/>' in category_axes[0]
    assert '<c:delete val="1"/>' in category_axes[1]


def test_naming_the_secondary_axis_enables_it(tmp_path):
    xml = _chart_xml(tmp_path, "implied.pptx", LINE_TITLED)
    assert xml.count("<c:valAx>") == 2


def test_primary_axis_formatting_does_not_reach_the_secondary_axis(tmp_path):
    path = tmp_path / "isolation.pptx"
    with Presentation.new("Axis Isolation") as pres:
        slide = pres.slides[0]
        slide.add_combo_chart(
            CATEGORIES,
            BAR_SERIES,
            LINE_TITLED,
            title="Isolation",
            bounds=BOUNDS,
        )
        slide.charts[0].apply_format({
            "value_axis_has_title": True,
            "value_axis_title": "Revenue (k)",
        })
        pres.save(path)

    with zipfile.ZipFile(path) as archive:
        xml = archive.read("ppt/charts/chart1.xml").decode("utf-8")

    primary, secondary = re.findall(r"<c:valAx>.*?</c:valAx>", xml, re.DOTALL)
    assert "Revenue (k)" in primary
    assert "Revenue (k)" not in secondary
    assert "Growth %" in secondary


def test_unmarked_line_series_stays_on_the_primary_axis(tmp_path):
    """A mixed set is drawn as one line plot per axis pair (Codex review)."""
    mixed = [
        dict(LINE_SERIES[0], secondary_axis=True),
        {"name": "Target (k)", "values": [800.0, 900.0, 950.0, 1200.0]},
    ]
    xml = _chart_xml(tmp_path, "mixed.pptx", mixed)

    plots = re.findall(r"<c:lineChart>.*?</c:lineChart>", xml, re.DOTALL)
    assert len(plots) == 2

    bar_ids = _plot_axis_ids(xml, "barChart")
    primary = [plot for plot in plots if "Target (k)" in plot]
    secondary = [plot for plot in plots if "Growth %" in plot]
    assert len(primary) == 1
    assert len(secondary) == 1
    assert "Growth %" not in primary[0]
    assert "Target (k)" not in secondary[0]
    assert re.findall(r'<c:axId val="(\d+)"/>', primary[0]) == bar_ids
    assert set(re.findall(r'<c:axId val="(\d+)"/>', secondary[0])).isdisjoint(bar_ids)
