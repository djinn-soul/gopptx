"""Text utilities for gopptx — fluent run builder and helpers."""

from __future__ import annotations

from .lang_script import (
    SCRIPT_COMPLEX,
    SCRIPT_EAST_ASIAN,
    SCRIPT_LATIN,
    lang_to_script,
    script_kind_for_language,
)
from .run_builder import RunBuilder

__all__ = [
    "SCRIPT_COMPLEX",
    "SCRIPT_EAST_ASIAN",
    "SCRIPT_LATIN",
    "RunBuilder",
    "lang_to_script",
    "script_kind_for_language",
]
