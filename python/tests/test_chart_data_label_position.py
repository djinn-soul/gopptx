import pathlib
import zipfile

from gopptx import Presentation
from gopptx.presentation.charts import ChartType


def test_chart_data_labels_position_right_office365_compat(
    tmp_path: pathlib.Path,
) -> None:
    """Setting data_labels position = 'RIGHT' on a bar chart generates valid OOXML without breaking Office 365 web."""
    output_path = tmp_path / "chart_position_right.pptx"

    with Presentation.new(title="Chart Data Labels Test") as pres:
        slide = pres.slides[0]

        # Add a bar chart using ChartType enum
        chart_id = slide.add_chart(
            ChartType.BAR,
            ["East", "West", "Midwest"],
            [{"name": "Series 1", "values": [19.2, 21.4, 16.7]}],
        )
        assert chart_id > 0

        chart = slide.charts[0]
        plot = chart.plots[0]
        plot.data_labels.position = "RIGHT"

        pres.save(output_path)

    # Reopen presentation deck to verify shape and state
    with Presentation(output_path) as pres:
        shapes = pres.slides[0].list_shapes()
        assert len(shapes) > 0

    with zipfile.ZipFile(output_path) as archive:
        chart_xml = archive.read("ppt/charts/chart1.xml").decode("utf-8")
    assert '<c:dLblPos val="outEnd"/>' in chart_xml
    assert '<c:dLblPos val="r"/>' not in chart_xml


def test_chart_data_label_state_survives_reopen(tmp_path: pathlib.Path) -> None:
    """Data-label properties are read from persisted OOXML after reopening."""
    output_path = tmp_path / "chart_data_label_state.pptx"
    with Presentation.new(title="Chart Data Label State") as pres:
        slide = pres.slides[0]
        _ = slide.add_chart(ChartType.BAR, ["A", "B"], [1.0, 2.0])
        labels = slide.charts[0].plots[0].data_labels
        labels.position = "outEnd"
        labels.show_value = True
        labels.show_category_name = True
        labels.show_series_name = True
        labels.number_format = "$#,##0"
        labels.word_wrap = False
        pres.save(output_path)

    with zipfile.ZipFile(output_path) as archive:
        chart_xml = archive.read("ppt/charts/chart1.xml").decode("utf-8")
    assert (
        '<c:txPr><a:bodyPr wrap="none"/><a:lstStyle/>'
        '<a:p><a:endParaRPr lang="en-US"/></a:p></c:txPr>'
    ) in chart_xml

    with Presentation(output_path) as pres:
        plot = pres.slides[0].charts[0].plots[0]
        labels = plot.data_labels
        assert plot.data_labels_visible is True
        assert labels.position == "outEnd"
        assert labels.show_value is True
        assert labels.show_category_name is True
        assert labels.show_series_name is True
        assert labels.number_format == "$#,##0"
        assert labels.word_wrap is False
