"""VBA macro mixin for gopptx library."""

from __future__ import annotations

import base64

from ..ops import OP_ADD_VBA
from .helpers import PresentationProtocol


class PresentationVBAMixin(PresentationProtocol):
    """Mixin providing VBA macro manipulation methods for Presentation."""

    def add_vba_project(self, data: bytes) -> None:
        """Add a VBA project binary blob to the presentation.

        Args:
            data: The raw CFB binary data for vbaProject.bin.
        """
        self.execute(
            OP_ADD_VBA,
            {"data": base64.b64encode(data).decode("ascii")},
        )
