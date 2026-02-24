"""Exception classes for gopptx library."""

from __future__ import annotations


class GopptxError(Exception):
    """Base exception for gopptx library errors."""

    def __init__(self, message: str, code: str | None = None) -> None:
        """Initialize the error with message and optional error code."""
        super().__init__(message)
        self.code = code
