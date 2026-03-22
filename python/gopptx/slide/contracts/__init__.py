"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from .base import BaseEngineProtocol as BaseEngineProtocol
from .charts import ChartOperationsProtocol as ChartOperationsProtocol
from .lifecycle import SlideLifecycleProtocol as SlideLifecycleProtocol
from .media import MediaOperationsProtocol as MediaOperationsProtocol
from .notes import NotesOperationsProtocol as NotesOperationsProtocol
from .presentation import SlidePresentationProtocol as SlidePresentationProtocol
from .shapes import ShapeOperationsProtocol as ShapeOperationsProtocol
from .tables import TableOperationsProtocol as TableOperationsProtocol
from .text import TextOperationsProtocol as TextOperationsProtocol
