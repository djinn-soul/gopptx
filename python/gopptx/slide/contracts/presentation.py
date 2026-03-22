"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from .base import BaseEngineProtocol
from .charts import ChartOperationsProtocol
from .lifecycle import SlideLifecycleProtocol
from .media import MediaOperationsProtocol
from .notes import NotesOperationsProtocol
from .shapes import ShapeOperationsProtocol
from .tables import TableOperationsProtocol
from .text import TextOperationsProtocol


class SlidePresentationProtocol(
    BaseEngineProtocol,
    SlideLifecycleProtocol,
    TextOperationsProtocol,
    ShapeOperationsProtocol,
    MediaOperationsProtocol,
    ChartOperationsProtocol,
    TableOperationsProtocol,
    NotesOperationsProtocol,
):
    """Composite protocol for the full presentation bridge."""
