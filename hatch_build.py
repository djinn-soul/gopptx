"""Custom build hook for building Go bridge during wheel build."""

import os
import pathlib
import subprocess  # noqa: S404
import sys
from typing import Any

from hatchling.builders.hooks.plugin.interface import BuildHookInterface


class CustomBuildHook(BuildHookInterface):
    """Build hook to compile Go shared library for Python bindings."""

    _BINARY_NAMES = ("gopptx.dll", "libgopptx.so", "libgopptx.dylib")

    @staticmethod
    def _release_build_enabled() -> bool:
        raw = os.getenv("GOPPTX_RELEASE_BUILD", "").strip().lower()
        return raw in {"1", "true", "yes", "on"}

    @staticmethod
    def _target_lib_name(platform_name: str) -> str:
        if platform_name == "win32":
            return "gopptx.dll"
        if platform_name == "darwin":
            return "libgopptx.dylib"
        return "libgopptx.so"

    @classmethod
    def _hide_non_target_binaries(
        cls,
        project_root: pathlib.Path,
        pkg_dir: pathlib.Path,
        target_lib_name: str,
    ) -> list[tuple[pathlib.Path, pathlib.Path]]:
        hidden_dir = project_root / ".hatch_build_hidden_bins"
        hidden_dir.mkdir(parents=True, exist_ok=True)
        hidden: list[tuple[pathlib.Path, pathlib.Path]] = []
        for candidate in cls._BINARY_NAMES:
            if candidate == target_lib_name:
                continue
            src = pkg_dir / candidate
            if not src.exists():
                continue
            dst = hidden_dir / candidate
            if dst.exists():
                dst.unlink()
            src.replace(dst)
            hidden.append((dst, src))
        return hidden

    @staticmethod
    def _restore_hidden_binaries(hidden: list[tuple[pathlib.Path, pathlib.Path]]) -> None:
        for dst, src in hidden:
            if dst.exists():
                if src.exists():
                    src.unlink()
                dst.replace(src)

    def initialize(self, version: str, build_data: dict[str, Any]) -> None:
        """Initialize the build hook and compile Go shared library.

        Args:
            version: The version being built (required by interface, may be unused).
            build_data: Build data dictionary to modify.
        """
        _ = version  # Acknowledge unused parameter required by interface
        if self.target_name != "wheel":
            return

        project_root = pathlib.Path(self.root)
        go_bridge_src = project_root / "bindings" / "c" / "bridge.go"
        pkg_dir = project_root / "python" / "gopptx"

        lib_name = self._target_lib_name(sys.platform)

        hidden = self._hide_non_target_binaries(project_root, pkg_dir, lib_name)
        build_data["gopptx_hidden_bins"] = [(str(dst), str(src)) for dst, src in hidden]

        out_path = pkg_dir / lib_name

        cmd = [
            "go",
            "build",
            "-o",
            str(out_path),
            "-buildmode=c-shared",
            str(go_bridge_src),
        ]
        if self._release_build_enabled():
            cmd[2:2] = ["-trimpath", "-buildvcs=false", "-ldflags=-s -w"]

        try:
            subprocess.run(cmd, check=True, capture_output=True, text=True)
        except subprocess.CalledProcessError:
            sys.exit(1)

        build_data["artifacts"].append(f"python/gopptx/{lib_name}")
        import sysconfig
        plat = sysconfig.get_platform().replace("-", "_").replace(".", "_")
        build_data["tag"] = f"py3-none-{plat}"
        build_data["pure_python"] = False

    def finalize(
        self,
        version: str,  # noqa: ARG002
        build_data: dict[str, Any],
        artifact_path: str,  # noqa: ARG002
    ) -> None:
        """Restore any non-target binaries hidden during wheel build."""
        hidden = build_data.pop("gopptx_hidden_bins", [])
        parsed_hidden = [(pathlib.Path(dst_raw), pathlib.Path(src_raw)) for dst_raw, src_raw in hidden]
        self._restore_hidden_binaries(parsed_hidden)
        hidden_dir = pathlib.Path(self.root) / ".hatch_build_hidden_bins"
        if hidden_dir.exists() and not any(hidden_dir.iterdir()):
            hidden_dir.rmdir()
