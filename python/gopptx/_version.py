"""Package version helpers."""

from importlib.metadata import PackageNotFoundError, version

try:
    __version__ = version("gopptx")
except PackageNotFoundError:
    # Source-tree imports can happen before package metadata exists.
    __version__ = "0.1.0"
