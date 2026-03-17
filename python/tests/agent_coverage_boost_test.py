import base64

import pytest
from gopptx.presentation.batch import BatchContext
from gopptx.presentation.presentation import Presentation

# 1x1 Transparent PNG
PNG_DATA = base64.b64decode(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="
)


@pytest.fixture
def prs():
    p = Presentation.new("Title")
    yield p
    p.close()


def test_batch_context_getattr_error(prs):
    with BatchContext(prs) as batch, pytest.raises(AttributeError):
        _ = batch.non_existent_attribute


def test_batch_context_exception_abort(prs):
    class MockPresentation:
        def __init__(self):
            self.aborted = False

        def begin_batch(self, **kwargs):
            pass

        def abort_batch(self):
            self.aborted = True

        def end_batch(self):
            return []

    mock_prs = MockPresentation()
    try:
        with BatchContext(mock_prs):
            raise ValueError("Abort")
    except ValueError:
        pass
    assert mock_prs.aborted


def test_shape_mixin_extra(prs, tmp_path):
    slide = prs.add_slide("S1")
    # add_connector
    slide.add_connector("straight", 0, 0, 10, 10)
    # add_video
    v = tmp_path / "v.mp4"
    v.write_bytes(b"fake video")
    slide.add_video(str(v), (0, 0, 10, 10))
    # add_audio
    a = tmp_path / "a.mp3"
    a.write_bytes(b"fake audio")
    slide.add_audio(str(a), (0, 0, 10, 10))
    # add_image & metadata (use real PNG)
    img = tmp_path / "i.png"
    img.write_bytes(PNG_DATA)
    sid = slide.add_image(str(img), (0, 0, 10, 10))
    meta = slide.get_image_metadata(sid)
    assert meta is not None
    # group/ungroup
    b1 = slide.add_textbox(0, 0, 5, 5)
    b2 = slide.add_textbox(5, 5, 5, 5)
    gid = slide.group_shapes([b1, b2])
    slide.ungroup_shapes(gid)
    # Z-order
    slide.move_shape_to_front(b1)
    slide.move_shape_to_back(b1)


def test_slide_extra(prs):
    # Presentation.new creates 1 slide by default.
    initial_count = len(prs.slides)
    assert initial_count == 1
    slide = prs.add_slide("S1")
    assert len(prs.slides) == 2
    # update
    slide.update(title="New Title", layout="title_content")
    assert slide.title == "New Title"
    # duplicate
    new_slide = slide.duplicate()
    assert len(prs.slides) == 3
    # remove
    new_slide.remove()
    assert len(prs.slides) == 2


def test_table_collections_extra(prs):
    slide = prs.add_slide("S1")
    tid = slide.add_table(2, 2, (0, 0, 10, 10))
    table = slide.shape(tid).table
    row = table.rows[0]
    row.height = 100
    col = table.columns[0]
    col.width = 100
    assert int(row.height) == 100
    assert int(col.width) == 100


def test_text_cache_mixin_extra(prs):
    slide = prs.add_slide("S1")
    sid = slide.add_textbox(0, 0, 10, 10, text="T")
    slide.update_shape_run_text(sid, 0, "U")
    runs = slide.get_shape_runs(sid)
    assert runs[0]["text"] == "U"
    slide.set_shape_runs(sid, [{"text": "R"}])
    slide.append_shape_run(sid, {"text": "A"})
    updated = slide.get_shape_runs(sid)
    assert "".join(run["text"] for run in updated) == "RA"
