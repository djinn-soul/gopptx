import os  # noqa: D100
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118


# Create output directory
output_dir = os.path.join(project_root, "examples/output")  # noqa: PTH118
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

image_path = os.path.join(  # noqa: PTH118
    project_root, "examples/assets/55/repository-open-graph-template.png"
)
output_path = os.path.join(output_dir, "python_image_test.pptx")  # noqa: PTH118

if not pathlib.Path(image_path).exists():
    sys.exit(0)

try:
    with Presentation.new("Image Test") as pres:
        slide = pres.slides[0]

        # Add image to the first slide
        shape_id = slide.add_image(image_path, 1000000, 2000000, 4000000, 2000000)

        pres.save(output_path)


except Exception:  # noqa: BLE001
    import traceback

    traceback.print_exc()
