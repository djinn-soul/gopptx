"""Run proxies for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .text_run import Run, RunHyperlink

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .text_model import ShapeTextFrame


class _ShapeRunProxy:
    """Live run proxy backed by bridge operations."""

    def __init__(self, text_frame: ShapeTextFrame, run_index: int) -> None:
        self._text_frame = text_frame
        self._run_index = run_index

    def _run_payload(self) -> dict[str, object]:
        runs = self._text_frame.get_runs()
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        return cast("dict[str, object]", runs[self._run_index])

    def _set_field(self, name: str, value: object) -> None:
        runs = [dict(run) for run in self._text_frame.get_runs()]
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        runs[self._run_index][name] = value
        self._text_frame.replace_runs(runs)

    @property
    def text(self) -> str:
        return str(self._run_payload().get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._text_frame.update_run_text(self._run_index, value)

    @property
    def bold(self) -> bool | None:
        value = self._run_payload().get("bold")
        return bool(value) if isinstance(value, bool) else None

    @bold.setter
    def bold(self, value: bool | None) -> None:
        self._set_field("bold", value)

    @property
    def italic(self) -> bool | None:
        value = self._run_payload().get("italic")
        return bool(value) if isinstance(value, bool) else None

    @italic.setter
    def italic(self, value: bool | None) -> None:
        self._set_field("italic", value)

    @property
    def underline(self) -> str | None:
        value = self._run_payload().get("underline")
        return str(value) if isinstance(value, str) else None

    @underline.setter
    def underline(self, value: str | None) -> None:
        self._set_field("underline", value)

    @property
    def hyperlink(self) -> RunHyperlink:
        payload = self._run_payload().get("hyperlink")
        if isinstance(payload, dict):
            parsed = RunHyperlink.from_payload(cast("dict[str, object]", payload))
            if parsed is not None:
                return parsed
        return RunHyperlink()

    @hyperlink.setter
    def hyperlink(self, value: RunHyperlink | dict[str, object] | None) -> None:
        if isinstance(value, RunHyperlink):
            self._set_field("hyperlink", value.to_payload())
            return
        self._set_field("hyperlink", value)

    def remove(self) -> None:
        """Remove this run from the owning shape text frame."""
        runs = [dict(run) for run in self._text_frame.get_runs()]
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        del runs[self._run_index]
        self._text_frame.replace_runs(runs)


class _ShapeRunCollection:
    """Live run collection for a paragraph proxy."""

    def __init__(self, text_frame: ShapeTextFrame) -> None:
        self._text_frame = text_frame
        self._run_proxies: dict[int, _ShapeRunProxy] = {}

    def __len__(self) -> int:
        return len(self._text_frame.get_runs())

    def __getitem__(self, index: int) -> _ShapeRunProxy:
        if index < 0:
            index += len(self)
        if index < 0:
            raise IndexError("run index out of range")
        proxy = self._run_proxies.get(index)
        if proxy is None:
            proxy = _ShapeRunProxy(self._text_frame, index)
            self._run_proxies[index] = proxy
        return proxy

    def __iter__(self) -> Iterator[_ShapeRunProxy]:
        for index, _ in enumerate(self._text_frame.get_runs()):
            yield _ShapeRunProxy(self._text_frame, index)

    def add_run(self, text: str = "") -> _ShapeRunProxy:
        self._text_frame.append_run(Run(text=text).to_payload())
        return self[len(self) - 1]

    @staticmethod
    def remove(run: _ShapeRunProxy) -> None:
        """Remove one run proxy from this collection."""
        run.remove()
