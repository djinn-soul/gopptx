"""Notes text-frame paragraph/run traversal helpers."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .notes_slide import NotesShape


class NotesTextFrame:
    """Minimal text-frame proxy for notes text placeholders."""

    def __init__(self, shape: NotesShape) -> None:
        """Initialize a text-frame proxy for one notes shape."""
        super().__init__()
        self._shape = shape

    @property
    def text(self) -> str:
        """Return full text content for the shape."""
        return self._shape.text

    @text.setter
    def text(self, value: str) -> None:
        """Replace full text content for the shape."""
        self._shape.text = value

    @property
    def paragraphs(self) -> NotesParagraphCollection:
        """Return paragraph collection view for this text frame."""
        return NotesParagraphCollection(self)

    def clear(self) -> None:
        """Clear all text from the text frame."""
        self.text = ""

    def add_paragraph(self, text: str = "") -> NotesParagraph:
        """Append a paragraph and return its proxy."""
        return self.paragraphs.add_paragraph(text)

    def paragraph_texts(self) -> list[str]:
        """Return paragraph texts split from the full text value."""
        text = self.text
        if not text:
            return [""]
        return text.split("\n")

    def set_paragraph_texts(self, paragraphs: list[str]) -> None:
        """Persist paragraph texts back to the shape as newline-separated text."""
        if not paragraphs:
            self.text = ""
            return
        self.text = "\n".join(paragraphs)

    def _paragraph_texts(self) -> list[str]:
        return self.paragraph_texts()

    def _set_paragraph_texts(self, paragraphs: list[str]) -> None:
        self.set_paragraph_texts(paragraphs)


class NotesRun:
    """Run proxy for notes text; normalized to a single run per paragraph."""

    def __init__(self, text_frame: NotesTextFrame, paragraph_index: int) -> None:
        """Initialize a run proxy for one paragraph index."""
        super().__init__()
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index

    @property
    def text(self) -> str:
        """Return run text for the underlying paragraph."""
        paragraphs = self._text_frame.paragraph_texts()
        if self._paragraph_index < 0 or self._paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        return paragraphs[self._paragraph_index]

    @text.setter
    def text(self, value: str) -> None:
        """Update run text for the underlying paragraph."""
        paragraphs = self._text_frame.paragraph_texts()
        if self._paragraph_index < 0 or self._paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        paragraphs[self._paragraph_index] = value
        self._text_frame.set_paragraph_texts(paragraphs)


class NotesRunCollection:
    """Run collection for one notes paragraph."""

    def __init__(self, text_frame: NotesTextFrame, paragraph_index: int) -> None:
        """Initialize a run collection bound to one paragraph."""
        super().__init__()
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index

    def __len__(self) -> int:
        """Return normalized run count for notes paragraphs."""
        return 1

    def __iter__(self) -> Iterator[NotesRun]:
        """Iterate runs for this paragraph."""
        yield self[0]

    def __getitem__(self, index: int) -> NotesRun:
        """Return the single supported run by index (0 or -1)."""
        if index not in {0, -1}:
            raise IndexError("run index out of range")
        return NotesRun(self._text_frame, self._paragraph_index)

    def add_run(self, text: str = "") -> NotesRun:
        """Append text to the single run and return it."""
        run = self[0]
        run.text += text
        return run


class NotesParagraph:
    """Paragraph proxy for notes text-frame traversal."""

    def __init__(self, text_frame: NotesTextFrame, index: int) -> None:
        """Initialize a paragraph proxy by index."""
        super().__init__()
        self._text_frame = text_frame
        self._index = index

    @property
    def text(self) -> str:
        """Return paragraph text."""
        paragraphs = self._text_frame.paragraph_texts()
        if self._index < 0 or self._index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        return paragraphs[self._index]

    @text.setter
    def text(self, value: str) -> None:
        """Replace paragraph text."""
        paragraphs = self._text_frame.paragraph_texts()
        if self._index < 0 or self._index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        paragraphs[self._index] = value
        self._text_frame.set_paragraph_texts(paragraphs)

    @property
    def runs(self) -> NotesRunCollection:
        """Return run collection for this paragraph."""
        return NotesRunCollection(self._text_frame, self._index)

    def add_run(self, text: str = "") -> NotesRun:
        """Append text to the paragraph run."""
        return self.runs.add_run(text)


class NotesParagraphCollection:
    """Paragraph collection for notes text-frame."""

    def __init__(self, text_frame: NotesTextFrame) -> None:
        """Initialize paragraph collection for a text frame."""
        super().__init__()
        self._text_frame = text_frame

    def __len__(self) -> int:
        """Return paragraph count."""
        return len(self._text_frame.paragraph_texts())

    def __iter__(self) -> Iterator[NotesParagraph]:
        """Iterate paragraph proxies in order."""
        for index in range(len(self)):
            yield NotesParagraph(self._text_frame, index)

    def __getitem__(self, index: int) -> NotesParagraph:
        """Return one paragraph proxy by index."""
        count = len(self)
        if index < 0:
            index += count
        if index < 0 or index >= count:
            raise IndexError("paragraph index out of range")
        return NotesParagraph(self._text_frame, index)

    def add_paragraph(self, text: str = "") -> NotesParagraph:
        """Append a paragraph and return its proxy."""
        paragraphs = self._text_frame.paragraph_texts()
        if len(paragraphs) == 1 and not paragraphs[0]:
            paragraphs[0] = text
            self._text_frame.set_paragraph_texts(paragraphs)
            return NotesParagraph(self._text_frame, 0)
        paragraphs.append(text)
        self._text_frame.set_paragraph_texts(paragraphs)
        return NotesParagraph(self._text_frame, len(paragraphs) - 1)
