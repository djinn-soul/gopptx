"""Custom build hook for building Go bridge during wheel build."""

import pathlib
import subprocess
import sys
from typing import Any

from hatchling.builders.hooks.plugin.interface import BuildHookInterface


class CustomBuildHook(BuildHookInterface):
    """Build hook to compile Go shared library for Python bindings."""

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

        out_path = pkg_dir / lib_name

        cmd = [
            "go",
            "build",
            "-o",
            str(out_path),
            "-buildmode=c-shared",
            str(go_bridge_src),
        ]

        try:
            subprocess.run(cmd, check=True, capture_output=True, text=True)
        except subprocess.CalledProcessError:
            sys.exit(1)

        build_data["artifacts"].append(f"python/gopptx/{lib_name}")
