"""Tests for bridge response safety checks."""

import pytest
from gopptx.presentation.comments.comment_mixin import PresentationCommentMixin
from gopptx.presentation.shapes.shape_media_mixin import PresentationShapeMediaMixin
from gopptx.presentation.shapes.shapes_tables import PresentationShapeMixin
from gopptx.presentation.tables.table_mixin import PresentationTableMixin
from gopptx.slide.shapes.shape_mixin import SlideShapeMixin
from gopptx.slide.shapes.smartart_anim_mixin import SlideSmartArtAnimMixin


class _MockPresentation(
    PresentationTableMixin,
    PresentationShapeMixin,
    PresentationShapeMediaMixin,
    PresentationCommentMixin,
    SlideShapeMixin,
    SlideSmartArtAnimMixin,
):
    def __init__(self, response=None):
        self.response = response or {}
        self.calls = []
        self._index = 0

    @property
    def index(self) -> int:
        return self._index

    @property
    def _presentation(self):
        return self

    def execute(self, op: str, payload: dict[str, object]) -> dict[str, object]:
        self.calls.append((op, payload))
        return self.response

    def _invalidate_shape_cache_if_present(self):
        pass

    def _invalidate_text_state_cache_if_present(self):
        pass


def test_add_table_id_safety():
    """add_table should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": "not-an-int"})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_table(slide_index=0, rows=2, cols=2, x=0, y=0, cx=0, cy=0)


def test_add_shape_id_safety():
    """add_shape should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": None})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_shape(0, "rect", (0, 0, 100, 100))


def test_add_textbox_id_safety():
    """add_textbox should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": []})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_textbox(0, 0, 0, 100, 100)


def test_group_id_safety():
    """group_shapes should raise TypeError if group_id is not an int."""
    prs = _MockPresentation({"group_id": "invalid"})
    with pytest.raises(TypeError, match="bridge response group_id must be an int"):
        prs.group_shapes(0, [1, 2])


def test_ungroup_id_safety():
    """ungroup_shapes should raise TypeError if group_id is not an int."""
    prs = _MockPresentation({"group_id": 1.5})
    with pytest.raises(TypeError, match="bridge response group_id must be an int"):
        prs.ungroup_shapes(0, 1)


def test_add_connector_id_safety():
    """add_connector should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": "conn"})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_connector(0, "line", 0, 0, 100, 100)


def test_add_group_shape_id_safety():
    """add_group_shape should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": (1,)})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_group_shape(0, [1, 2])


def test_add_image_id_safety():
    """add_image should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": {}})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_image(0, source="test.png")


def test_add_video_id_safety():
    """add_video should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": "vid"})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_video(0, source="test.mp4", bounds=(0, 0, 1, 1))


def test_add_audio_id_safety():
    """add_audio should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": []})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_audio(0, source="test.mp3", bounds=(0, 0, 1, 1))


def test_add_ole_object_id_safety():
    """add_ole_object should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": None})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_ole_object(0, source="test.xlsx", bounds=(0, 0, 1, 1))


def test_add_author_id_safety():
    """add_author should raise TypeError if author_id is not an int."""
    prs = _MockPresentation({"author_id": "0"})
    with pytest.raises(TypeError, match="bridge response author_id must be an int"):
        prs.add_author("Name", "Initials")


def test_add_smartart_id_safety():
    """add_smartart should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": 1.1})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.add_smartart("layout", ["item"], (0, 0, 1, 1))


def test_commit_freeform_id_safety():
    """commit_freeform should raise TypeError if shape_id is not an int."""
    prs = _MockPresentation({"shape_id": 1.0})
    with pytest.raises(TypeError, match="bridge response shape_id must be an int"):
        prs.commit_freeform(0, [(0, 0), (1, 1)], close=True)


def test_add_mermaid_shape_count_safety():
    """add_mermaid should raise TypeError if shape_count is not an int."""
    prs = _MockPresentation({"shape_count": "many", "connector_count": 0})
    with pytest.raises(TypeError, match="bridge response shape_count must be an int"):
        prs.add_mermaid("graph TD; A-->B")


def test_add_mermaid_connector_count_safety():
    """add_mermaid should raise TypeError if connector_count is not an int."""
    prs = _MockPresentation({"shape_count": 2, "connector_count": None})
    with pytest.raises(
        TypeError, match="bridge response connector_count must be an int"
    ):
        prs.add_mermaid("graph TD; A-->B")
