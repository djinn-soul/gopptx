"""Chart area and 3D scene proxies."""
# ruff: noqa: D101,D102,D107,SLF001
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from typing import Protocol

    from ...schemas import ChartFormatUpdate, ChartState

    class _ChartLike(Protocol):
        def _snapshot(self) -> ChartState: ...
        def _apply_format(self, fmt: ChartFormatUpdate) -> None: ...


class ChartScene3D:
    def __init__(self, chart: _ChartLike) -> None:
        self._chart = chart

    def _state(self) -> dict[str, object]:
        snapshot = self._chart._snapshot()
        raw = snapshot.get("scene3d", {})
        return cast("dict[str, object]", raw)

    def _apply_partial(self, changes: dict[str, object]) -> None:
        merged = dict(self._state())
        merged.update(changes)
        camera = merged.get("camera_preset")
        rig = merged.get("light_rig")
        direction = merged.get("light_direction")
        if not isinstance(camera, str) or not camera:
            raise ValueError("scene3d requires camera_preset")
        if not isinstance(rig, str) or not rig:
            raise ValueError("scene3d requires light_rig")
        if not isinstance(direction, str) or not direction:
            raise ValueError("scene3d requires light_direction")
        payload: ChartFormatUpdate = {
            "camera_preset": camera,
            "light_rig": rig,
            "light_direction": direction,
        }
        fov = merged.get("camera_field_of_view")
        if isinstance(fov, int):
            payload["camera_field_of_view"] = fov
        revolution = merged.get("light_rig_revolution")
        if isinstance(revolution, bool):
            payload["light_rig_revolution"] = revolution
        self._chart._apply_format(payload)

    def update(
        self,
        *,
        camera_preset: str,
        light_rig: str,
        light_direction: str,
        camera_field_of_view: int | None = None,
        light_rig_revolution: bool | None = None,
    ) -> None:
        changes: dict[str, object] = {
            "camera_preset": camera_preset,
            "light_rig": light_rig,
            "light_direction": light_direction,
        }
        if camera_field_of_view is not None:
            changes["camera_field_of_view"] = camera_field_of_view
        if light_rig_revolution is not None:
            changes["light_rig_revolution"] = light_rig_revolution
        self._apply_partial(changes)

    @property
    def camera_preset(self) -> str | None:
        value = self._state().get("camera_preset")
        return str(value) if isinstance(value, str) else None

    @camera_preset.setter
    def camera_preset(self, value: str) -> None:
        self._apply_partial({"camera_preset": value})

    @property
    def camera_field_of_view(self) -> int | None:
        value = self._state().get("camera_field_of_view")
        return int(value) if isinstance(value, int) else None

    @camera_field_of_view.setter
    def camera_field_of_view(self, value: int) -> None:
        self._apply_partial({"camera_field_of_view": int(value)})

    @property
    def light_rig(self) -> str | None:
        value = self._state().get("light_rig")
        return str(value) if isinstance(value, str) else None

    @light_rig.setter
    def light_rig(self, value: str) -> None:
        self._apply_partial({"light_rig": value})

    @property
    def light_direction(self) -> str | None:
        value = self._state().get("light_direction")
        return str(value) if isinstance(value, str) else None

    @light_direction.setter
    def light_direction(self, value: str) -> None:
        self._apply_partial({"light_direction": value})

    @property
    def light_rig_revolution(self) -> bool | None:
        value = self._state().get("light_rig_revolution")
        return bool(value) if isinstance(value, bool) else None

    @light_rig_revolution.setter
    def light_rig_revolution(self, value: bool) -> None:
        self._apply_partial({"light_rig_revolution": bool(value)})


class ChartArea:
    def __init__(self, chart: _ChartLike) -> None:
        self._scene3d = ChartScene3D(chart)

    @property
    def scene3d(self) -> ChartScene3D:
        return self._scene3d
