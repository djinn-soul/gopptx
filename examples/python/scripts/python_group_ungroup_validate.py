"""Generate a PPTX from the Python API and validate basic PPTX structure."""

from __future__ import annotations

import pathlib
import sys
import zipfile


def project_root_from_here() -> pathlib.Path:
    return pathlib.Path(__file__).resolve().parents[3]


def run() -> None:
    root = project_root_from_here()
    sys.path.append(str(root / "python"))

    from gopptx import Presentation  # noqa: PLC0415

    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "python_generated_validated.pptx"

    with Presentation.new("Python Generated Validation") as pres:
        slide = pres.add_slide("Python API Slide")
        shape_id_1 = pres.add_shape(
            slide.index,
            "rect",
            bounds=(1_000_000, 1_000_000, 2_000_000, 900_000),
            text="Shape A",
        )
        shape_id_2 = pres.add_shape(
            slide.index,
            "ellipse",
            bounds=(3_300_000, 1_000_000, 2_000_000, 900_000),
            text="Shape B",
        )

        shapes = pres.list_shapes(slide.index)
        found_shape_ids = {int(shape["ID"]) for shape in shapes}
        if shape_id_1 not in found_shape_ids or shape_id_2 not in found_shape_ids:
            raise RuntimeError("generated slide does not contain expected shapes")

        pres.save(str(output_path))

    with zipfile.ZipFile(output_path) as archive:
        required_parts = {
            "[Content_Types].xml",
            "_rels/.rels",
            "ppt/presentation.xml",
            "ppt/slides/slide1.xml",
        }
        names = set(archive.namelist())
        missing_parts = sorted(required_parts.difference(names))
        if missing_parts:
            raise RuntimeError(f"missing required PPTX parts: {missing_parts}")

    print(f"Generated: {output_path}")
    print("PPTX structure validation passed")


if __name__ == "__main__":
    run()
