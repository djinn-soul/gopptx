"""Generate a PPTX from the Python API and validate it with the Go smoke validator."""

from __future__ import annotations

import pathlib
import subprocess
import sys


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

    validate_cmd = [
        "go",
        "run",
        "./scripts/smoke/validate_smoke_outputs",
        "-file",
        str(output_path),
    ]
    completed = subprocess.run(
        validate_cmd,
        cwd=root,
        capture_output=True,
        text=True,
        check=False,
    )
    if completed.returncode != 0:
        raise RuntimeError(
            "validator failed\n"
            f"stdout:\n{completed.stdout}\n"
            f"stderr:\n{completed.stderr}"
        )

    print(f"Generated: {output_path}")
    print(completed.stdout.strip())


if __name__ == "__main__":
    run()
