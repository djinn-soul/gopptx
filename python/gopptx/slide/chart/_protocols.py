"""Structural protocols for chart-domain proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol

if TYPE_CHECKING:
    from ...schemas import SlideChartRef
    from ..contracts.presentation import SlidePresentationProtocol


class ChartSlideProto(Protocol):
    @property
    def presentation(self) -> SlidePresentationProtocol: ...

    @property
    def index(self) -> int: ...

    def list_charts(self) -> list[SlideChartRef]: ...
