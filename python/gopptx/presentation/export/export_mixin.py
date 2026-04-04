"""Export mixin for gopptx library — PDF and HTML export."""

from __future__ import annotations

import warnings
from dataclasses import dataclass, field
from typing import Literal

from ... import ops
from ..helpers import PresentationMixinBase


@dataclass
class PDFOptions:
    """Options for PDF export.

    Attributes:
        driver: Which PDF backend to use.
            ``"auto"``  — try LibreOffice/PowerPoint first, then native fallback
            ``"native"`` — built-in Go PDF renderer (experimental fidelity)
            ``"libreoffice"`` — requires LibreOffice on PATH
            ``"powerpoint"`` — requires Microsoft PowerPoint (Windows only)
        font_paths: Extra font directories for the native renderer.
    """

    driver: Literal["auto", "native", "libreoffice", "powerpoint"] = "auto"
    font_paths: list[str] = field(default_factory=list)


@dataclass
class HTMLOptions:
    """Options for HTML export.

    Attributes:
        embed_images: Inline images as base64 data URIs (default True).
        include_navigation: Add prev/next JS navigation buttons (default True).
        include_slide_numbers: Overlay slide number on each slide (default True).
        base_url: Asset path prefix when ``embed_images`` is False.
    """

    embed_images: bool = True
    include_navigation: bool = True
    include_slide_numbers: bool = True
    base_url: str = ""


class PresentationExportMixin(PresentationMixinBase):
    """Mixin providing PDF and HTML export methods."""

    def save_as_pdf(
        self,
        output_path: str | None = None,
        options: PDFOptions | None = None,
    ) -> str:
        """Export the presentation to a PDF file.

        Args:
            output_path: Destination path for the PDF file. If omitted,
                defaults to ``presentation.pdf`` in the current working directory.
            options: Optional :class:`PDFOptions` controlling the export driver.

        Returns:
            Absolute path to the written PDF file.
        """
        opts = options or PDFOptions()
        if opts.driver == "native":
            warnings.warn(
                "PDF driver 'native' is experimental and may not match PowerPoint rendering for all layouts.",
                UserWarning,
                stacklevel=2,
            )
        payload: dict[str, object] = {
            "driver": opts.driver,
        }
        if output_path is not None:
            payload["output_path"] = output_path
        result = self.execute(ops.OP_EXPORT_PDF, payload)
        return str(result.get("output_path", output_path or "presentation.pdf"))

    def export_html(
        self,
        output_path: str | None = None,
        options: HTMLOptions | None = None,
    ) -> str:
        """Export the presentation to an HTML document.

        Args:
            output_path: If given, write the HTML to this file and return the
                absolute path.  If omitted, return the HTML string directly.
            options: Optional :class:`HTMLOptions` controlling the output.

        Returns:
            Absolute file path when ``output_path`` is provided, otherwise the
            raw HTML string.
        """
        opts = options or HTMLOptions()
        payload: dict[str, object] = {
            "embed_images": opts.embed_images,
            "include_navigation": opts.include_navigation,
            "include_slide_numbers": opts.include_slide_numbers,
            "base_url": opts.base_url,
        }
        if output_path is not None:
            payload["output_path"] = output_path

        result = self.execute(ops.OP_EXPORT_HTML, payload)

        if output_path is not None:
            return str(result.get("output_path", output_path))
        return str(result.get("html", ""))
