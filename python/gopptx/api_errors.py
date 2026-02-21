from __future__ import annotations

from typing import Optional


class GopptxError(Exception):
    """Base exception for gopptx library errors."""

    def __init__(self, message: str, code: Optional[str] = None):
        super().__init__(message)
        self.code = code
