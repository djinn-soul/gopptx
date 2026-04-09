from __future__ import annotations

import importlib.util
import sys
import types
from pathlib import Path


def _load_custom_build_hook() -> type:
    interface_module = types.ModuleType("hatchling.builders.hooks.plugin.interface")

    class _BuildHookInterface:
        """Minimal stub for test import."""

    interface_module.BuildHookInterface = _BuildHookInterface
    sys.modules["hatchling"] = types.ModuleType("hatchling")
    sys.modules["hatchling.builders"] = types.ModuleType("hatchling.builders")
    sys.modules["hatchling.builders.hooks"] = types.ModuleType(
        "hatchling.builders.hooks"
    )
    sys.modules["hatchling.builders.hooks.plugin"] = types.ModuleType(
        "hatchling.builders.hooks.plugin"
    )
    sys.modules["hatchling.builders.hooks.plugin.interface"] = interface_module

    module_path = Path(__file__).resolve().parents[2] / "hatch_build.py"
    spec = importlib.util.spec_from_file_location("hatch_build", module_path)
    if spec is None or spec.loader is None:
        raise AssertionError("failed to load hatch_build.py module spec")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module.CustomBuildHook


CustomBuildHook = _load_custom_build_hook()


def _write_dummy_binaries(pkg_dir: Path) -> None:
    pkg_dir.mkdir(parents=True, exist_ok=True)
    (pkg_dir / "gopptx.dll").write_bytes(b"dll")
    (pkg_dir / "libgopptx.so").write_bytes(b"so")
    (pkg_dir / "libgopptx.dylib").write_bytes(b"dylib")


def test_target_lib_name_for_platforms() -> None:
    target_name_attr = "_target_lib_name"
    target_lib_name = getattr(CustomBuildHook, target_name_attr)
    assert target_lib_name("win32") == "gopptx.dll"
    assert target_lib_name("darwin") == "libgopptx.dylib"
    assert target_lib_name("linux") == "libgopptx.so"
    assert target_lib_name("linux2") == "libgopptx.so"


def test_hide_non_target_binaries_for_linux(tmp_path: Path) -> None:
    project_root = tmp_path / "repo"
    pkg_dir = project_root / "python" / "gopptx"
    _write_dummy_binaries(pkg_dir)

    hide_name_attr = "_hide_non_target_binaries"
    hide_non_target_binaries = getattr(CustomBuildHook, hide_name_attr)
    hidden = hide_non_target_binaries(
        project_root=project_root,
        pkg_dir=pkg_dir,
        target_lib_name="libgopptx.so",
    )

    assert (pkg_dir / "libgopptx.so").exists()
    assert not (pkg_dir / "gopptx.dll").exists()
    assert not (pkg_dir / "libgopptx.dylib").exists()
    assert len(hidden) == 2

    restore_name_attr = "_restore_hidden_binaries"
    restore_hidden_binaries = getattr(CustomBuildHook, restore_name_attr)
    restore_hidden_binaries(hidden)

    assert (pkg_dir / "gopptx.dll").exists()
    assert (pkg_dir / "libgopptx.so").exists()
    assert (pkg_dir / "libgopptx.dylib").exists()


def test_hide_non_target_binaries_for_windows(tmp_path: Path) -> None:
    project_root = tmp_path / "repo"
    pkg_dir = project_root / "python" / "gopptx"
    _write_dummy_binaries(pkg_dir)

    hide_name_attr = "_hide_non_target_binaries"
    hide_non_target_binaries = getattr(CustomBuildHook, hide_name_attr)
    hidden = hide_non_target_binaries(
        project_root=project_root,
        pkg_dir=pkg_dir,
        target_lib_name="gopptx.dll",
    )

    assert (pkg_dir / "gopptx.dll").exists()
    assert not (pkg_dir / "libgopptx.so").exists()
    assert not (pkg_dir / "libgopptx.dylib").exists()
    assert len(hidden) == 2

    restore_name_attr = "_restore_hidden_binaries"
    restore_hidden_binaries = getattr(CustomBuildHook, restore_name_attr)
    restore_hidden_binaries(hidden)

    assert (pkg_dir / "gopptx.dll").exists()
    assert (pkg_dir / "libgopptx.so").exists()
    assert (pkg_dir / "libgopptx.dylib").exists()


def test_hide_non_target_binaries_for_macos(tmp_path: Path) -> None:
    project_root = tmp_path / "repo"
    pkg_dir = project_root / "python" / "gopptx"
    _write_dummy_binaries(pkg_dir)

    hide_name_attr = "_hide_non_target_binaries"
    hide_non_target_binaries = getattr(CustomBuildHook, hide_name_attr)
    hidden = hide_non_target_binaries(
        project_root=project_root,
        pkg_dir=pkg_dir,
        target_lib_name="libgopptx.dylib",
    )

    assert (pkg_dir / "libgopptx.dylib").exists()
    assert not (pkg_dir / "gopptx.dll").exists()
    assert not (pkg_dir / "libgopptx.so").exists()
    assert len(hidden) == 2

    restore_name_attr = "_restore_hidden_binaries"
    restore_hidden_binaries = getattr(CustomBuildHook, restore_name_attr)
    restore_hidden_binaries(hidden)

    assert (pkg_dir / "gopptx.dll").exists()
    assert (pkg_dir / "libgopptx.so").exists()
    assert (pkg_dir / "libgopptx.dylib").exists()
