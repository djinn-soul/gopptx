"""Regression coverage for recursive leaf-shape iteration (issue #435)."""

from __future__ import annotations

from gopptx import Presentation
from gopptx.schemas import Inches


def test_iter_leaf_shapes_recurses_through_nested_groups(tmp_path) -> None:
    """Nested groups expose only their non-group descendants in document order."""
    output = tmp_path / "nested-groups.pptx"
    with Presentation.new("Issue 435") as presentation:
        slide = presentation.slides[0]
        slide.shapes.clear()
        first = slide.add_shape(
            "rectangle",
            (Inches(1), Inches(1), Inches(1), Inches(1)),
            text="A",
        )
        second = slide.add_shape(
            "ellipse",
            (Inches(2.5), Inches(1), Inches(1), Inches(1)),
            text="B",
        )
        inner = slide.add_group_shape([first, second])
        third = slide.add_shape(
            "triangle",
            (Inches(4), Inches(1), Inches(1), Inches(1)),
            text="C",
        )
        outer = slide.add_group_shape([inner, third])
        assert [shape.id for shape in slide.shapes.iter_leaf_shapes()] == [
            first,
            second,
            third,
        ]
        assert slide.shape(outer).shapes[0].shapes[1].text == "B"
        presentation.save(output)

    reopened = Presentation()
    reopened.open(output)
    try:
        leaves = list(reopened.slides[0].shapes.iter_leaf_shapes())
        assert [shape.id for shape in leaves] == [first, second, third]
        assert [shape.text for shape in leaves] == ["A", "B", "C"]
        assert all(shape.shape_type != "grpSp" for shape in leaves)
    finally:
        reopened.close()
