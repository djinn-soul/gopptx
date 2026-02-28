# ruff: noqa: D100, D103, S101
import pathlib

import pytest
from gopptx.presentation.presentation import Presentation

# Add project root to sys.path
project_root = pathlib.Path(__file__).parent.parent.parent.resolve()


@pytest.fixture
def presentation() -> Presentation:
    p = Presentation()
    p.open(str(project_root / "testdata" / "placeholders.pptx"))
    return p


def test_list_placeholders(presentation: Presentation) -> None:
    # placeholders.pptx slide 0 is a title/body slide
    slide = presentation.slides[0]
    placeholders = slide.placeholders()

    assert len(placeholders) > 0

    types = [ph.placeholder_format for ph in placeholders]
    assert "title" in types or "ctrTitle" in types

    # Check that we can lookup by idx
    first_idx = placeholders[0].idx
    found = slide.get_placeholder(first_idx)
    assert found is not None
    assert found.idx == first_idx


def test_insert_placeholder_text(presentation: Presentation) -> None:
    slide = presentation.slides[0]
    placeholders = slide.placeholders()

    # Find any placeholder we can type into
    ph = placeholders[0]

    # Insert text
    ph.insert_text("Injected Title Text")

    # We can verify by looking at the shapes collection to see if the text was written
    shapes = slide.list_shapes()

    found_text = False
    for s in shapes:
        if s.get("Text") == "Injected Title Text":
            found_text = True
            break

    assert found_text, "Failed to find injected text in shapes list"


def test_insert_placeholder_picture(
    presentation: Presentation, tmp_path: pathlib.Path
) -> None:
    # We just need to insert a picture and make sure it doesn't crash,
    # since we don't have a specific picture placeholder fixture here.
    slide = presentation.slides[0]
    ph = slide.placeholders()[0]

    img_path = tmp_path / "test_pic.jpg"
    pathlib.Path(img_path).write_bytes(b"fake image data")

    # For now, we'll just test that the API call executes without error
    ph.insert_picture(str(img_path))

    # Verify shape count increased or at least it didn't throw
    shapes = slide.list_shapes()
    pics = [s for s in shapes if s.get("Type") == "pic"]
    assert len(pics) > 0


def test_placeholder_repr(presentation: Presentation) -> None:
    slide = presentation.slides[0]
    ph = slide.placeholders()[0]

    r = repr(ph)
    assert "Placeholder" in r
    assert f"idx={ph.idx}" in r
