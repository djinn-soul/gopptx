"""Notes domain package for slide APIs."""

from .notes_shape import NotesShape
from .notes_shape_collection import NotesShapeCollection
from .notes_slide import NotesSlide
from .notes_slide_style_mixin import NotesSlideStyleMixin
from .notes_text_model import (
    NotesParagraph,
    NotesParagraphCollection,
    NotesRun,
    NotesRunCollection,
    NotesTextFrame,
)

__all__ = [
    "NotesParagraph",
    "NotesParagraphCollection",
    "NotesRun",
    "NotesRunCollection",
    "NotesShape",
    "NotesShapeCollection",
    "NotesSlide",
    "NotesSlideStyleMixin",
    "NotesTextFrame",
]
