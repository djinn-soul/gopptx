"""Run proxies for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from ...schemas import RGBColor
from .text_run import RunHyperlink

_EMU_PER_POINT = 12700

if TYPE_CHECKING:
    from .text_frame_protocol import ShapeTextFrameProtocol


class _RunCollectionProtocol(Protocol):
    """The part of the owning run collection a run proxy calls back into.

    Structural typing here keeps the collection module importing this one and
    not the other way round.
    """

    def reindex_after_removal(self, removed_index: int) -> None:
        """Shift cached proxies down after a run is removed."""
        ...


class _ShapeRunProxy:
    """Live run proxy backed by bridge operations."""

    def __init__(
        self,
        text_frame: ShapeTextFrameProtocol,
        paragraph_index: int,
        run_index: int,
        collection: _RunCollectionProtocol | None = None,
    ) -> None:
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index
        self._run_index = run_index
        # The owning collection, when there is one, is told about removals so
        # proxies for later runs shift down with the runs they point at.
        self._collection = collection

    def rebind_index(self, run_index: int) -> None:
        """Point this proxy at ``run_index`` after the run list shifted."""
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
    def underline(self) -> bool | str | None:
        value = self.payload().get("underline")
        if not value:
            return False
        if value in {"single", "sng", True}:
            return True
        return str(value)

    @underline.setter
    def underline(self, value: bool | str | None) -> None:
        if value is True:
            self.set_field("underline", "single")
        elif value is False or value is None:
            self.set_field("underline", "")
        else:
            self.set_field("underline", str(value))

    @property
    def strikethrough(self) -> bool | str | None:
        """Return whether run has strikethrough formatting (Issue #339)."""
        value = self.payload().get("strikethrough")
        if not value or value == "none" or value is False:
            return False
        if value in {"sngStrike", "sng", True}:
            return True
        return str(value)

    @strikethrough.setter
    def strikethrough(self, value: bool | str | None) -> None:
        """Set strikethrough formatting (Issue #339)."""
        if value is True:
            self.set_field("strikethrough", "sngStrike")
        elif value is False or value is None:
            self.set_field("strikethrough", "none")
        else:
            self.set_field("strikethrough", str(value))

    @property
    def hyperlink(self) -> _RunHyperlinkProxy:
        return _RunHyperlinkProxy(self)

    @hyperlink.setter
    def hyperlink(self, value: RunHyperlink | dict[str, object] | str | None) -> None:
        if isinstance(value, str):
            self.set_field("hyperlink", {"address": value})
            return
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
        removed_index = self._run_index
        del runs[removed_index]
        self._text_frame.replace_paragraph_runs(self._paragraph_index, runs)
        # This proxy no longer points at a run; leaving its index in place would
        # silently retarget it at whichever run shifted into that position.
        self._run_index = -1
        if self._collection is not None:
            self._collection.reindex_after_removal(removed_index)

    def delete(self) -> None:
        """Alias for remove (Issue #144)."""
        self.remove()


class _FontColorProxy:
    """Non-mutating font color proxy for text run (Issue #1111)."""

    def __init__(self, run_proxy: _ShapeRunProxy) -> None:
        self._run_proxy = run_proxy

    @property
    def rgb(self) -> RGBColor | None:
        """Return the color as an RGBColor without mutating XML (Issue #1111).

        Matches ``shape.fill.fore_color.rgb`` and ``shape.line.color.rgb``,
        which have always returned RGBColor.
        """
        payload = self._run_proxy.payload()
        val = payload.get("color")
        if not isinstance(val, str) or not val:
            return None
        try:
            return RGBColor.from_string(val.lstrip("#"))
        except ValueError:
            return None

    @rgb.setter
    def rgb(self, value: RGBColor | None) -> None:
        if value is None:
            self._run_proxy.set_field("color", None)
        else:
            val_str = str(value).lstrip("#")
            self._run_proxy.set_field("color", val_str)

    @property
    def type(self) -> str | None:
        """Return color type ('RGB' or None) without mutating XML (Issue #1111)."""
        return "RGB" if self.rgb else None


class _RunHyperlinkProxy:
    """Live proxy for run hyperlink attributes (Issue #151)."""

    def __init__(self, run_proxy: _ShapeRunProxy) -> None:
        self._run_proxy = run_proxy

    def _get_field(self, field: str) -> object:
        payload = self._run_proxy.payload().get("hyperlink")
        if isinstance(payload, dict):
            return cast("dict[str, object]", payload).get(field)
        return None

    def _set_field(self, field: str, value: object) -> None:
        raw_payload = self._run_proxy.payload().get("hyperlink")
        payload = (
            dict(cast("dict[str, object]", raw_payload))
            if isinstance(raw_payload, dict)
            else {}
        )
        if value is None:
            payload.pop(field, None)
        else:
            payload[field] = value
        self._run_proxy.set_field("hyperlink", payload or None)

    @property
    def address(self) -> str | None:
        val = self._get_field("address")
        return str(val) if val else None

    @address.setter
    def address(self, value: str | None) -> None:
        self._set_field("address", value)

    @property
    def url(self) -> str | None:
        return self.address

    @url.setter
    def url(self, value: str | None) -> None:
        self.address = value

    @property
    def tooltip(self) -> str | None:
        val = self._get_field("tooltip")
        return str(val) if val else None

    @tooltip.setter
    def tooltip(self, value: str | None) -> None:
        self._set_field("tooltip", value)

    @property
    def action(self) -> str | None:
        val = self._get_field("action")
        return str(val) if val else None

    @action.setter
    def action(self, value: str | None) -> None:
        self._set_field("action", value)

    @property
    def target_slide(self) -> int | None:
        val = self._get_field("target_slide")
        return int(val) if isinstance(val, (int, float)) else None

    @target_slide.setter
    def target_slide(self, value: int | None) -> None:
        self._set_field("target_slide", value)

    @property
    def jump(self) -> str | None:
        val = self._get_field("jump")
        return str(val) if val else None

    @jump.setter
    def jump(self, value: str | None) -> None:
        self._set_field("jump", value)

    @property
    def macro(self) -> str | None:
        val = self._get_field("macro")
        return str(val) if val else None

    @macro.setter
    def macro(self, value: str | None) -> None:
        self._set_field("macro", value)


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
    def size(self) -> int | None:
        """Font size in EMU, so ``Point(28)`` round-trips (python-pptx parity)."""
        points = self.size_pt
        return None if points is None else int(points * _EMU_PER_POINT)

    @size.setter
    def size(self, value: float | None) -> None:
        self.size_pt = None if value is None else float(value) / _EMU_PER_POINT

    @property
    def size_pt(self) -> float | None:
        """Font size in points."""
        val = self._run_proxy.payload().get("size_pt")
        return float(val) if isinstance(val, (int, float)) else None

    @size_pt.setter
    def size_pt(self, value: float | None) -> None:
        self._run_proxy.set_field("size_pt", None if value is None else float(value))

    @property
    def bold(self) -> bool | None:
        return self._run_proxy.bold

    @bold.setter
    def bold(self, value: bool | None) -> None:
        self._run_proxy.bold = value

    @property
    def italic(self) -> bool | None:
        return self._run_proxy.italic

    @italic.setter
    def italic(self, value: bool | None) -> None:
        self._run_proxy.italic = value

    @property
    def underline(self) -> bool | str | None:
        """True/False, or the OOXML underline style name (e.g. ``"dash"``)."""
        return self._run_proxy.underline

    @underline.setter
    def underline(self, value: bool | str | None) -> None:
        self._run_proxy.underline = value

    @property
    def strikethrough(self) -> bool | str | None:
        """Return whether run has strikethrough formatting (Issue #339)."""
        return self._run_proxy.strikethrough

    @strikethrough.setter
    def strikethrough(self, value: bool | str | None) -> None:
        """Set strikethrough formatting (Issue #339)."""
        self._run_proxy.strikethrough = value
