"""Run collection for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING

from .text_run import Run
from .text_run_model import _ShapeRunProxy

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .text_frame_protocol import ShapeTextFrameProtocol


class _ShapeRunCollection:
    """Live run collection for a paragraph proxy."""

    def __init__(
        self, text_frame: ShapeTextFrameProtocol, paragraph_index: int
    ) -> None:
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index
        self._run_proxies: dict[int, _ShapeRunProxy] = {}

    def __len__(self) -> int:
        return len(self._text_frame.get_paragraph_runs(self._paragraph_index))

    def __getitem__(self, index: int) -> _ShapeRunProxy:
        if index < 0:
            index += len(self)
        if index < 0:
            raise IndexError("run index out of range")
        proxy = self._run_proxies.get(index)
        if proxy is None:
            proxy = _ShapeRunProxy(self._text_frame, self._paragraph_index, index, self)
            self._run_proxies[index] = proxy
        return proxy

    def __iter__(self) -> Iterator[_ShapeRunProxy]:
        for index, _ in enumerate(
            self._text_frame.get_paragraph_runs(self._paragraph_index)
        ):
            yield self[index]

    def reindex_after_removal(self, removed_index: int) -> None:
        """Shift cached proxies down after the run at ``removed_index`` is gone.

        Run proxies address runs by position, so deleting an earlier run would
        otherwise leave retained proxies pointing one run past their own.
        """
        rebound: dict[int, _ShapeRunProxy] = {}
        for index, proxy in self._run_proxies.items():
            if index == removed_index:
                continue
            new_index = index - 1 if index > removed_index else index
            proxy.rebind_index(new_index)
            rebound[new_index] = proxy
        self._run_proxies = rebound

    def add_run(self, text: str = "") -> _ShapeRunProxy:
        runs = self._text_frame.get_paragraph_runs(self._paragraph_index)
        runs.append(Run(text=text).to_payload())
        self._text_frame.replace_paragraph_runs(self._paragraph_index, runs)
        return self[len(self) - 1]

    @staticmethod
    def remove(run: _ShapeRunProxy) -> None:
        """Remove one run proxy from this collection."""
        run.remove()

    def __delitem__(self, index: int) -> None:
        """Delete run at index (Issue #144)."""
        self[index].remove()
