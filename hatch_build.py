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

        if sys.platform == "win32":
            lib_name = "gopptx.dll"
        elif sys.platform == "darwin":
            lib_name = "libgopptx.dylib"
        else:
            lib_name = "libgopptx.so"

        hidden_dir = project_root / ".hatch_build_hidden_bins"
        hidden_dir.mkdir(parents=True, exist_ok=True)
        hidden: list[tuple[pathlib.Path, pathlib.Path]] = []
        for candidate in self._BINARY_NAMES:
            if candidate == lib_name:
                continue
            src = pkg_dir / candidate
            if not src.exists():
                continue
            dst = hidden_dir / candidate
            if dst.exists():
                dst.unlink()
            src.replace(dst)
            hidden.append((dst, src))
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
        for dst_raw, src_raw in hidden:
            dst = pathlib.Path(dst_raw)
            src = pathlib.Path(src_raw)
            if dst.exists():
                if src.exists():
                    src.unlink()
                dst.replace(src)
        hidden_dir = pathlib.Path(self.root) / ".hatch_build_hidden_bins"
        if hidden_dir.exists() and not any(hidden_dir.iterdir()):
            hidden_dir.rmdir()
