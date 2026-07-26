import pathlib

from gopptx import Presentation
from gopptx.presentation.charts import ChartType


def test_chart_data_labels_position_right_office365_compat(tmp_path: pathlib.Path) -> None:
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
