"""Live text object model facades for shape text editing."""
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportArgumentType=false, reportCallIssue=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from ...text.run_builder import RunBuilder
from ..shapes.shape_text_frame import ShapeTextFrame as BaseShapeTextFrame
from .text_paragraph_model import _ShapeParagraphCollection, _ShapeParagraphProxy
from .text_run_model import _ShapeRunCollection, _ShapeRunProxy


class ShapeTextFrame(BaseShapeTextFrame):
    """Public text-frame facade with paragraph collection access."""

    def __init__(self, slide: object, shape_id: int) -> None:
        """Create a text-frame proxy for a shape."""
        super().__init__(slide, shape_id)
        self._paragraphs: _ShapeParagraphCollection | None = None

    @property
    def paragraphs(self) -> _ShapeParagraphCollection:
        """Return paragraph collection facade."""
        if self._paragraphs is None:
            self._paragraphs = _ShapeParagraphCollection(self)
        return self._paragraphs

    @property
    def text(self) -> str:
        """Return full text content of the frame."""
        return "\n".join(p.text for p in self.paragraphs)

    @text.setter
    def text(self, value: str) -> None:
        """Replace all text in the frame with a single paragraph."""
        self.set_runs([RunBuilder(value)])

    def set_runs(self, builders: list[RunBuilder]) -> None:
        """Replace all runs from a list of :class:`~gopptx.text.RunBuilder` instances.

        Each builder is converted to a ``TextRun`` payload dict and the full
        run list is written to the bridge in a single round-trip.

        Example::

            from gopptx.text import RunBuilder

            shape.text_frame.set_runs([
                RunBuilder("Hello ").bold().color("FF0000"),
                RunBuilder("World").italic(),
            ])
        """
        payloads = [b.to_payload() for b in builders]
        self.replace_runs(payloads)


ShapeRunProxy = _ShapeRunProxy
ShapeRunCollection = _ShapeRunCollection
ShapeParagraphProxy = _ShapeParagraphProxy
ShapeParagraphCollection = _ShapeParagraphCollection
