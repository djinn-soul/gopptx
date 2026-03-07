"""Tests for chart/placeholder/shape object-model parity helpers."""

from __future__ import annotations

from gopptx import CategoryChartData, Presentation


def test_chart_object_model_updates() -> None:
    with Presentation.new("Chart Object Model") as prs:
        slide = prs.add_slide("Chart")
        data = CategoryChartData(categories=["A", "B"])
        data.add_series("S1", [1.0, 2.0])
        _shape_id = slide.add_chart("bar", data)

        chart = slide.charts[0]
        chart.title_visible = True
        chart.chart_title.text = "Updated Title"
        chart.legend.visible = True
        chart.legend.position = "r"
        chart.category_axis.tick_label_position = "low"
        chart.category_axis.major_gridlines_visible = True
        chart.category_axis.crosses = "max"
        chart.value_axis.tick_label_position = "high"
        chart.value_axis.major_gridlines_visible = False
        chart.value_axis.crosses = "min"
        chart.chart_area.scene3d.update(
            camera_preset="orthographicFront",
            light_rig="threePt",
            light_direction="t",
            camera_field_of_view=45,
            light_rig_revolution=True,
        )
        chart.plots[0].data_labels_visible = True
        chart.plots[0].data_labels.show_value = True
        assert chart.category_axis.tick_label_position == "low"
        assert chart.category_axis.crosses == "max"
        assert chart.value_axis.tick_label_position == "high"
        assert chart.value_axis.crosses == "min"
        assert chart.chart_area.scene3d.camera_preset == "orthographicFront"
        assert chart.chart_area.scene3d.light_rig == "threePt"

        new_data = CategoryChartData(categories=["X", "Y"])
        new_data.add_series("S2", [3.0, 4.0])
        chart.replace_data(new_data)


def test_placeholder_inheritance_navigation() -> None:
    with Presentation.new("Placeholder Inheritance") as prs:
        slide = prs.add_slide("Placeholder")
        placeholders = list(slide.placeholders)
        if not placeholders:
            return

        ph = placeholders[0]
        assert ph.layout_placeholder is not None
        assert ph.master_placeholder is not None


def test_shape_fill_line_shadow_proxies() -> None:
    with Presentation.new("Shape Formatting") as prs:
        slide = prs.add_slide("Shape")
        shape_id = slide.add_shape(
            "rect", (1000000, 1000000, 2000000, 1000000), text="X"
        )
        shape = slide.shape(shape_id)

        shape.fill.solid_color = "FF0000"
        shape.line.color = "00FF00"
        shape.line.width = 25400
        shape.line.dash_style = "dash"
        shape.shadow.color = "333333"
        shape.shadow.blur_radius = 38100
        shape.shadow.distance = 25400
        shape.shadow.angle = 45.0

        _refreshed = slide.shape(shape_id)
