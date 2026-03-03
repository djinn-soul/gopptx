import os
import pathlib
import sys

from gopptx import SHAPE_ROUNDED_RECTANGLE, Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
sys.path.append(os.path.join(project_root, "python"))  # noqa: PTH118


def test_python_slide_objects(tmp_path: pathlib.Path) -> None:
    output_path = tmp_path / "python_slide_objects.pptx"

    with Presentation.new("Object Oriented Test") as pres:
        # 1. Access first slide
        slide = pres.slides[0]

        # 2. Update title via property
        slide.title = "Enhanced via Python"

        # 3. Add shape via slide object
        shape_id = slide.add_shape(
            SHAPE_ROUNDED_RECTANGLE,
            (1000000, 1000000, 3000000, 1500000),
            text="I am a proxy",
        )
        assert shape_id > 0

        # 4. Set notes via property
        slide.notes = "Proxy objects make API much cleaner."

        # 5. Duplicate slide
        new_slide = slide.duplicate()
        new_slide.title = "I am the copy"

        # 6. Verify newly exposed master/layout object proxies
        masters = pres.slide_masters
        if masters:
            first_master = masters[0]
            master_layouts = first_master.slide_layouts
            if master_layouts:
                _ = master_layouts[0]

        # 7. Save
        pres.save(output_path)
        assert pres.slide_count >= 2

    assert output_path.exists()
