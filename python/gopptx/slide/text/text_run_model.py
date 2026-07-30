"""Run proxies for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .text_run import Run, RunHyperlink

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .text_frame_protocol import ShapeTextFrameProtocol


class _ShapeRunProxy:
    """Live run proxy backed by bridge operations."""

    def __init__(
        self,
        text_frame: ShapeTextFrameProtocol,
        paragraph_index: int,
        run_index: int,
    ) -> None:
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index
        self._run_index = run_index

    def payload(self) -> dict[str, object]:
        runs = self._text_frame.get_paragraph_runs(self._paragraph_index)
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        return runs[self._run_index]

    def set_field(self, name: str, value: object) -> None:
        runs = [
            dict(run)
            for run in self._text_frame.get_paragraph_runs(self._paragraph_index)
        ]
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        runs[self._run_index][name] = value
        self._text_frame.replace_paragraph_runs(self._paragraph_index, runs)

    @property
    def text(self) -> str:
        return str(self.payload().get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self.set_field("text", value)

    @property
    def bold(self) -> bool | None:
        value = self.payload().get("bold")
        return bool(value) if isinstance(value, bool) else None

    @bold.setter
    def bold(self, value: bool | None) -> None:
        self.set_field("bold", value)

    @property
    def italic(self) -> bool | None:
        value = self.payload().get("italic")
        return bool(value) if isinstance(value, bool) else None

    @italic.setter
    def italic(self, value: bool | None) -> None:
        self.set_field("italic", value)

    @property
    def underline(self) -> str | None:
        value = self.payload().get("underline")
        return str(value) if isinstance(value, str) else None

    @underline.setter
    def underline(self, value: str | None) -> None:
        self.set_field("underline", value)

    @property
    def hyperlink(self) -> RunHyperlink:
        payload = self.payload().get("hyperlink")
        if isinstance(payload, dict):
            parsed = RunHyperlink.from_payload(cast("dict[str, object]", payload))
            if parsed is not None:
                return parsed
        return RunHyperlink()

    @hyperlink.setter
    def hyperlink(self, value: RunHyperlink | dict[str, object] | None) -> None:
        if isinstance(value, RunHyperlink):
            self.set_field("hyperlink", value.to_payload())
            return
        self.set_field("hyperlink", value)

    @property
    def font(self) -> _FontProxy:
        """Return non-mutating font facade for run font name, size, and color (Issue #1111)."""
        return _FontProxy(self)

    def remove(self) -> None:
        """Remove this run from the owning shape text frame."""
        runs = [
            dict(run)
            for run in self._text_frame.get_paragraph_runs(self._paragraph_index)
        ]
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        del runs[self._run_index]
        self._text_frame.replace_paragraph_runs(self._paragraph_index, runs)


class _FontColorProxy:
    """Non-mutating font color proxy for text run (Issue #1111)."""

    def __init__(self, run_proxy: _ShapeRunProxy) -> None:
        self._run_proxy = run_proxy

    @property
    def rgb(self) -> str | None:
        """Return the RGB hex color without mutating XML (Issue #1111)."""
        payload = self._run_proxy.payload()
        val = payload.get("color")
        return str(val) if isinstance(val, str) else None

    @rgb.setter
    def rgb(self, value: str | None) -> None:
        self._run_proxy.set_field("color", value)

    @property
    def type(self) -> str | None:
        """Return color type ('RGB' or None) without mutating XML (Issue #1111)."""
        return "RGB" if self.rgb else None


class _FontProxy:
    """Non-mutating font proxy for text run (Issue #1111)."""

    def __init__(self, run_proxy: _ShapeRunProxy) -> None:
        self._run_proxy = run_proxy

    @property
    def color(self) -> _FontColorProxy:
        """Return non-mutating font color facade (Issue #1111)."""
        return _FontColorProxy(self._run_proxy)

    @property
    def name(self) -> str | None:
        val = self._run_proxy.payload().get("font")
        return str(val) if isinstance(val, str) else None

    @name.setter
    def name(self, value: str | None) -> None:
        self._run_proxy.set_field("font", value)

    @property
    def size(self) -> float | None:
        val = self._run_proxy.payload().get("size_pt")
        return float(val) if isinstance(val, (int, float)) else None

    @size.setter
    def size(self, value: float | None) -> None:
        self._run_proxy.set_field("size_pt", value)


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
            proxy = _ShapeRunProxy(self._text_frame, self._paragraph_index, index)
            self._run_proxies[index] = proxy
        return proxy

    def __iter__(self) -> Iterator[_ShapeRunProxy]:
        for index, _ in enumerate(
            self._text_frame.get_paragraph_runs(self._paragraph_index)
        ):
            yield _ShapeRunProxy(self._text_frame, self._paragraph_index, index)

    def add_run(self, text: str = "") -> _ShapeRunProxy:
        runs = self._text_frame.get_paragraph_runs(self._paragraph_index)
        runs.append(Run(text=text).to_payload())
        self._text_frame.replace_paragraph_runs(self._paragraph_index, runs)
        return self[len(self) - 1]

    @staticmethod
    def remove(run: _ShapeRunProxy) -> None:
        """Remove one run proxy from this collection."""
        run.remove()
