"""Shared pytest fixtures for gopptx Python tests."""

import pytest
from gopptx.presentation.presentation import Presentation


@pytest.fixture
def presentation(tmp_path) -> Presentation:
    """Generate a temporary PPTX with title/body placeholder slides."""
    path = str(tmp_path / "placeholders.pptx")
    with Presentation.new("Placeholder Test Deck") as pres:
        pres.add_slide("Title Slide", layout="title_and_content", bullets=["Body text"])
        pres.remove_slide(0)  # Remove the auto-created initial slide (title-only)
        pres.save(path)
    p = Presentation()
    p.open(path)
    return p
