import os
import subprocess
import sys

from hatchling.builders.hooks.plugin.interface import BuildHookInterface


class CustomBuildHook(BuildHookInterface):
    def initialize(self, version, build_data) -> None:
        if self.target_name != "wheel":
            return

        project_root = self.root
        go_bridge_src = os.path.join(project_root, "bindings/c/bridge.go")
        pkg_dir = os.path.join(project_root, "python/gopptx")

        if sys.platform == "win32":
            lib_name = "gopptx.dll"
        elif sys.platform == "darwin":
            lib_name = "libgopptx.dylib"
        else:
            lib_name = "libgopptx.so"

        out_path = os.path.join(pkg_dir, lib_name)

        cmd = ["go", "build", "-o", out_path, "-buildmode=c-shared", go_bridge_src]

        try:
            subprocess.run(cmd, check=True, capture_output=True, text=True)
        except subprocess.CalledProcessError:
            sys.exit(1)

        build_data["artifacts"].append(f"python/gopptx/{lib_name}")
