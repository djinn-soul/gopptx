"""VBA macro mixin for gopptx library."""

from __future__ import annotations

import base64
from typing import TYPE_CHECKING

from ..ops import OP_ADD_VBA

if TYPE_CHECKING:
    from .base import PresentationBase


class PresentationVBAMixin:
    """Mixin providing VBA macro manipulation methods for Presentation."""

    def add_vba_project(self: PresentationBase, data: bytes) -> None:
        """Add a VBA project binary blob to the presentation.

        Args:
            data: The raw CFB binary data for vbaProject.bin.
        """
        self.execute(
            OP_ADD_VBA,
            {"data": base64.b64encode(data).decode("ascii")},
        )
