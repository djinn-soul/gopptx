import os  # noqa: D100
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118

# Ensure smoke sample exists
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118
if not pathlib.Path(input_deck).exists():
    sys.exit(1)

try:
    with Presentation(input_deck) as pres:
        # 1. Slide Count
        initial_count = pres.slide_count

        # 2. Metadata
        meta = pres.metadata

        # 3. Duplicate Slide
        new_idx = pres.duplicate_slide(0, 1)

        # 4. Move Slide
        pres.move_slide(0, pres.slide_count - 1)

        # 5. Remove Slide
        pres.remove_slide(1)

        # 6. Save
        output_dir = os.path.join(project_root, "examples/output")  # noqa: PTH118
        pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)
        out_path = os.path.join(output_dir, "python_management_output.pptx")  # noqa: PTH118
        pres.save(out_path)


except Exception as e:  # noqa: BLE001
    if hasattr(e, "code"):
        pass
    sys.exit(1)
