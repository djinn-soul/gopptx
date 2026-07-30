"""Chart, metadata, transform, and action features for shape proxies."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from ..chart.model import Chart
from ..text.text_run import RunHyperlink
from .picture_image import ImagePartProxy
from .shape_proxy_extra_features import ShapeProxyExtraMixin

if TYPE_CHECKING:
    from ...schemas import Shape, ShapeUpdate, SlideChartRef
    from ..chart._protocols import ChartSlideProto
    from ..contracts import SlidePresentationProtocol


class _ShapeFeatureSlideProto(Protocol):
    @property
    def presentation(self) -> SlidePresentationProtocol: ...

    @property
    def index(self) -> int: ...

    def list_charts(self) -> list[SlideChartRef]: ...

    def shape(self, shape_id: int) -> _ShapeFeatureHost: ...


class _ShapeFeatureHost(Protocol):
    @property
    def slide(self) -> _ShapeFeatureSlideProto: ...

    @property
    def id(self) -> int: ...

    @property
    def flip_horizontal(self) -> bool: ...

    @property
    def flip_vertical(self) -> bool: ...

    def shape_record(self) -> Shape:
        """Return the current shape payload."""
        ...

    def apply_update(self, patch: ShapeUpdate) -> None:
        """Apply a patch to the current shape."""
        ...


class ShapeProxyFeatureMixin(ShapeProxyExtraMixin):
    """Features kept separate from the core shape and collection facades."""

    @property
    def image(self: _ShapeFeatureHost) -> ImagePartProxy:
        """Return the image part represented by this shape."""
        return ImagePartProxy(self)

    def make_color_transparent(self: _ShapeFeatureHost, color: str) -> None:
        """Make pixels matching ``color`` transparent in this picture.

        ``color`` is a six-digit RGB value with an optional leading ``#``.
        PowerPoint stores this as an ``a:clrChange`` effect (Issue #165).
        """
        self.apply_update(cast("ShapeUpdate", {"transparent_color": color}))

    @property
    def has_chart(self: _ShapeFeatureHost) -> bool:
        """Return whether this exact shape contains a chart."""
        return any(
            chart["ShapeID"] == self.id
            for chart in self.slide.presentation.list_slide_charts(self.slide.index)
        )

    @property
    def shapes(self: _ShapeFeatureHost) -> list[_ShapeFeatureHost]:
        """Return child proxies when this is a group shape."""
        children = self.shape_record().get("Shapes", [])
        if not isinstance(children, list):
            return []
        proxies: list[_ShapeFeatureHost] = []
        for child in cast("list[object]", children):
            if not isinstance(child, dict):
                continue
            payload = cast("dict[str, object]", child)
            shape_id = payload.get("ID", payload.get("id"))
            if shape_id is not None:
                proxies.append(self.slide.shape(int(str(shape_id))))
        return proxies

    @property
    def chart(self: _ShapeFeatureHost) -> object:
        """Return the chart bound to this graphic-frame shape."""
        for chart in self.slide.presentation.list_slide_charts(self.slide.index):
            if chart["ShapeID"] == self.id:
                return Chart(
                    cast("ChartSlideProto", self.slide),
                    chart["Index"],
                    chart["RelID"],
                    chart["ChartPart"],
                )
        raise AttributeError("shape does not contain a chart")

    @property
    def alt_text(self: _ShapeFeatureHost) -> str:
        """Return the shape alternative-text description."""
        shape = self.shape_record()
        return str(
            shape.get(
                "AltText",
                shape.get(
                    "alt_text", shape.get("Description", shape.get("description", ""))
                ),
            )
        )

    @alt_text.setter
    def alt_text(self: _ShapeFeatureHost, value: str) -> None:
        """Set or clear the shape alternative-text description."""
        self.apply_update(
            cast("ShapeUpdate", {"alt_text": value, "description": value})
        )

    @property
    def description(self: _ShapeFeatureHost) -> str:
        """Return the shape alternative-text description."""
        shape = self.shape_record()
        return str(
            shape.get(
                "AltText",
                shape.get(
                    "alt_text", shape.get("Description", shape.get("description", ""))
                ),
            )
        )

    @description.setter
    def description(self: _ShapeFeatureHost, value: str) -> None:
        """Set or clear the shape alternative-text description."""
        self.apply_update(
            cast("ShapeUpdate", {"alt_text": value, "description": value})
        )

    @property
    def title(self: _ShapeFeatureHost) -> str:
        """Return the shape accessibility title."""
        shape = self.shape_record()
        return str(shape.get("Title", shape.get("title", "")))

    @title.setter
    def title(self: _ShapeFeatureHost, value: str) -> None:
        """Set or clear the shape accessibility title."""
        self.apply_update(cast("ShapeUpdate", {"title": value}))

    @property
    def rotation(self: _ShapeFeatureHost) -> float:
        """Return clockwise rotation in degrees."""
        value = self.shape_record().get("Rotation", 0.0)
        return float(value) if isinstance(value, (int, float)) else 0.0

    @rotation.setter
    def rotation(self: _ShapeFeatureHost, value: float) -> None:
        self.apply_update(cast("ShapeUpdate", {"rotation": float(value)}))

    @property
    def flip_horizontal(self: _ShapeFeatureHost) -> bool:
        """Return whether the shape is flipped horizontally."""
        shape = self.shape_record()
        return bool(shape.get("FlipH", shape.get("flip_h", False)))

    @flip_horizontal.setter
    def flip_horizontal(self: _ShapeFeatureHost, value: bool) -> None:
        self.apply_update(cast("ShapeUpdate", {"flip_h": bool(value)}))

    @property
    def flip_vertical(self: _ShapeFeatureHost) -> bool:
        """Return whether the shape is flipped vertically."""
        shape = self.shape_record()
        return bool(shape.get("FlipV", shape.get("flip_v", False)))

    @flip_vertical.setter
    def flip_vertical(self: _ShapeFeatureHost, value: bool) -> None:
        self.apply_update(cast("ShapeUpdate", {"flip_v": bool(value)}))

    @property
    def flip_x(self: _ShapeFeatureHost) -> bool:
        """Alias for flip_horizontal (Issue #547)."""
        return self.flip_horizontal

    @flip_x.setter
    def flip_x(self: _ShapeFeatureHost, value: bool) -> None:
        self.apply_update(cast("ShapeUpdate", {"flip_h": bool(value)}))

    @property
    def flip_y(self: _ShapeFeatureHost) -> bool:
        """Alias for flip_vertical (Issue #547)."""
        return self.flip_vertical

    @flip_y.setter
    def flip_y(self: _ShapeFeatureHost, value: bool) -> None:
        self.apply_update(cast("ShapeUpdate", {"flip_v": bool(value)}))

    def begin_connect(
        self: _ShapeFeatureHost,
        target_shape: _ShapeFeatureHost | int | str,
        connection_site: int = 0,
    ) -> None:
        """Connect this connector's start to a target shape."""
        target_id = (
            target_shape.id
            if not isinstance(target_shape, (int, str))
            else int(target_shape)
        )
        self.apply_update(
            cast(
                "ShapeUpdate",
                {
                    "start_connection": {
                        "shape_id": target_id,
                        "site_index": int(connection_site),
                    }
                },
            )
        )

    def end_connect(
        self: _ShapeFeatureHost,
        target_shape: _ShapeFeatureHost | int | str,
        connection_site: int = 0,
    ) -> None:
        """Connect this connector's end to a target shape."""
        target_id = (
            target_shape.id
            if not isinstance(target_shape, (int, str))
            else int(target_shape)
        )
        self.apply_update(
            cast(
                "ShapeUpdate",
                {
                    "end_connection": {
                        "shape_id": target_id,
                        "site_index": int(connection_site),
                    }
                },
            )
        )

    @property
    def click_action(
        self: _ShapeFeatureHost,
    ) -> RunHyperlink | dict[str, object] | None:
        """Return the click action for this shape."""
        shape = self.shape_record()
        return _action_from_payload(shape.get("ClickAction", shape.get("click_action")))

    @click_action.setter
    def click_action(
        self: _ShapeFeatureHost,
        value: RunHyperlink | dict[str, object] | None,
    ) -> None:
        payload = value.to_payload() if isinstance(value, RunHyperlink) else value
        self.apply_update(cast("ShapeUpdate", {"click_action": payload}))

    @property
    def hover_action(
        self: _ShapeFeatureHost,
    ) -> RunHyperlink | dict[str, object] | None:
        """Return the hover action for this shape."""
        shape = self.shape_record()
        return _action_from_payload(shape.get("HoverAction", shape.get("hover_action")))

    @hover_action.setter
    def hover_action(
        self: _ShapeFeatureHost,
        value: RunHyperlink | dict[str, object] | None,
    ) -> None:
        payload = value.to_payload() if isinstance(value, RunHyperlink) else value
        self.apply_update(cast("ShapeUpdate", {"hover_action": payload}))


def _action_from_payload(payload: object) -> RunHyperlink:
    if isinstance(payload, dict):
        action = RunHyperlink.from_payload(cast("dict[str, object]", payload))
        if action is not None:
            return action
    return RunHyperlink()
