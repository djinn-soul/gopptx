"""Notes text-frame paragraph/run traversal helpers."""
# ruff: noqa: D102, D105, D107, SLF001
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .notes_slide import NotesShape


class NotesTextFrame:
    """Minimal text-frame proxy for notes text placeholders."""

    def __init__(self, shape: NotesShape) -> None:
        super().__init__()
        self._shape = shape

    @property
    def text(self) -> str:
        return self._shape.text

    @text.setter
    def text(self, value: str) -> None:
        self._shape.text = value

    @property
    def paragraphs(self) -> NotesParagraphCollection:
        return NotesParagraphCollection(self)

    def clear(self) -> None:
        self.text = ""

    def add_paragraph(self, text: str = "") -> NotesParagraph:
        return self.paragraphs.add_paragraph(text)

    def _paragraph_texts(self) -> list[str]:
        text = self.text
        if not text:
            return [""]
        return text.split("\n")

    def _set_paragraph_texts(self, paragraphs: list[str]) -> None:
        if not paragraphs:
            self.text = ""
            return
        self.text = "\n".join(paragraphs)


class NotesRun:
    """Run proxy for notes text; normalized to a single run per paragraph."""

    def __init__(self, text_frame: NotesTextFrame, paragraph_index: int) -> None:
        super().__init__()
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index

    @property
    def text(self) -> str:
        paragraphs = self._text_frame._paragraph_texts()
        if self._paragraph_index < 0 or self._paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        return paragraphs[self._paragraph_index]

    @text.setter
    def text(self, value: str) -> None:
        paragraphs = self._text_frame._paragraph_texts()
        if self._paragraph_index < 0 or self._paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        paragraphs[self._paragraph_index] = value
        self._text_frame._set_paragraph_texts(paragraphs)


class NotesRunCollection:
    """Run collection for one notes paragraph."""

    def __init__(self, text_frame: NotesTextFrame, paragraph_index: int) -> None:
        super().__init__()
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index

    def __len__(self) -> int:
        return 1

    def __iter__(self) -> Iterator[NotesRun]:
        yield self[0]

    def __getitem__(self, index: int) -> NotesRun:
        if index not in {0, -1}:
            raise IndexError("run index out of range")
        return NotesRun(self._text_frame, self._paragraph_index)

    def add_run(self, text: str = "") -> NotesRun:
        run = self[0]
        run.text += text
        return run


class NotesParagraph:
    """Paragraph proxy for notes text-frame traversal."""

    def __init__(self, text_frame: NotesTextFrame, index: int) -> None:
        super().__init__()
        self._text_frame = text_frame
        self._index = index

    @property
    def text(self) -> str:
        paragraphs = self._text_frame._paragraph_texts()
        if self._index < 0 or self._index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        return paragraphs[self._index]

    @text.setter
    def text(self, value: str) -> None:
        paragraphs = self._text_frame._paragraph_texts()
        if self._index < 0 or self._index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        paragraphs[self._index] = value
        self._text_frame._set_paragraph_texts(paragraphs)

    @property
    def runs(self) -> NotesRunCollection:
        return NotesRunCollection(self._text_frame, self._index)

    def add_run(self, text: str = "") -> NotesRun:
        return self.runs.add_run(text)


class NotesParagraphCollection:
    """Paragraph collection for notes text-frame."""

    def __init__(self, text_frame: NotesTextFrame) -> None:
        super().__init__()
        self._text_frame = text_frame

    def __len__(self) -> int:
        return len(self._text_frame._paragraph_texts())

    def __iter__(self) -> Iterator[NotesParagraph]:
        for index in range(len(self)):
            yield NotesParagraph(self._text_frame, index)

    def __getitem__(self, index: int) -> NotesParagraph:
        count = len(self)
        if index < 0:
            index += count
        if index < 0 or index >= count:
            raise IndexError("paragraph index out of range")
        return NotesParagraph(self._text_frame, index)

    def add_paragraph(self, text: str = "") -> NotesParagraph:
        paragraphs = self._text_frame._paragraph_texts()
        if len(paragraphs) == 1 and not paragraphs[0]:
            paragraphs[0] = text
            self._text_frame._set_paragraph_texts(paragraphs)
            return NotesParagraph(self._text_frame, 0)
        paragraphs.append(text)
        self._text_frame._set_paragraph_texts(paragraphs)
        return NotesParagraph(self._text_frame, len(paragraphs) - 1)
