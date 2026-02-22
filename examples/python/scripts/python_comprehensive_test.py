import os
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../../..")
).resolve()
sys.path.append(os.path.join(project_root, "python"))


# Create output directory
output_dir = os.path.join(project_root, "examples/output")
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

input_path = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")
output_path = os.path.join(output_dir, "python_comprehensive_output.pptx")

if not pathlib.Path(input_path).exists():
    print(f"Input file {input_path} not found. Skipping test.")
    sys.exit(0)

try:
    print(f"Opening presentation: {input_path}")
    with Presentation(input_path) as pres:
        # 1. Slide Count
        count = pres.slide_count
        print(f"Original slide count: {count}")

        # 2. Metadata
        meta = pres.metadata
        print(f"Metadata: {meta}")

        # 3. Add Slide
        idx = pres.add_slide("New Python Slide")
        print(f"Added slide at index: {idx}")

        # 4. Shapes
        print("Adding shape...")
        shape_id = pres.add_shape(
            idx, "rect", 1000000, 1000000, 2000000, 1000000, text="Hello from Python!"
        )
        print(f"Added shape with ID: {shape_id}")

        print("Listing shapes...")
        shapes = pres.list_shapes(idx)
        for s in shapes:
            print(f"  Shape: {s}")

        print("Updating shape...")
        pres.update_shape(idx, shape_id, {"text": "Updated Text", "x": 1500000})

        # 5. Notes
        print("Setting notes...")
        pres.set_notes(idx, "These are speaker notes added via Python API.")
        notes = pres.get_notes(idx)
        print(f"Retrieved notes: {notes}")

        # 6. Layouts
        print("Listing layouts...")
        layouts = pres.list_slide_layouts()
        if layouts:
            print(f"  First layout: {layouts[0]}")
            # pres.rebind_slide_layout(idx, layouts[0]['part'])

        # 7. Core Properties
        print("Updating core properties...")
        props = pres.core_properties
        props["title"] = "Comprehensive Python Test"
        props["creator"] = "GoPPTX Python API"
        pres.core_properties = props

        # 8. Save
        pres.save(output_path)
        print(f"Saved to {output_path}")

    print("Success! Python API expansion verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
