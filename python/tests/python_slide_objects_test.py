import os  # noqa: D100
import pathlib
import sys

from gopptx import SHAPE_ROUNDED_RECTANGLE, Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118


# Create output directory
output_dir = os.path.join(project_root, "examples/output")  # noqa: PTH118
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

output_path = os.path.join(output_dir, "python_slide_objects.pptx")  # noqa: PTH118

try:
    with Presentation.new("Object Oriented Test") as pres:
        # 1. Access first slide
        slide = pres.slides[0]

        # 2. Update title via property
        slide.title = "Enhanced via Python"

        # 3. Add shape via slide object
        slide.add_shape(
            SHAPE_ROUNDED_RECTANGLE,
            1000000,
            1000000,
            3000000,
            1500000,
            text="I am a proxy",
        )

        # 4. Set notes via property
        slide.notes = "Proxy objects make API much cleaner."

        # 5. Duplicate slide
        new_slide = slide.duplicate()
        new_slide.title = "I am the copy"

        # 6. Verify newly exposed master/layout object proxies
        masters = pres.slide_masters
        if len(masters) > 0:
            first_master = masters[0]
            master_layouts = first_master.slide_layouts
            if len(master_layouts) > 0:
                first_layout = master_layouts[0]

        # 7. Save
        pres.save(output_path)


except Exception:  # noqa: BLE001
    import traceback

    traceback.print_exc()
