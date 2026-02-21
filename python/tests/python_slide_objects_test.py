import os
import pathlib
import sys

from gopptx import SHAPE_ROUNDED_RECTANGLE, Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../.."))
sys.path.append(os.path.join(project_root, "python"))


# Create output directory
output_dir = os.path.join(project_root, "examples/output")
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

output_path = os.path.join(output_dir, "python_slide_objects.pptx")

try:
    print("Testing object-oriented Slide proxy...")
    with Presentation.new("Object Oriented Test") as pres:
        # 1. Access first slide
        slide = pres.slides[0]
        print(f"Original title: '{slide.title}'")

        # 2. Update title via property
        slide.title = "Enhanced via Python"
        print(f"Updated title: '{slide.title}'")

        # 3. Add shape via slide object
        slide.add_shape(
            SHAPE_ROUNDED_RECTANGLE,
            1000000,
            1000000,
            3000000,
            1500000,
            text="I am a proxy",
        )
        print("Added shape via slide proxy.")

        # 4. Set notes via property
        slide.notes = "Proxy objects make API much cleaner."
        print(f"Slide notes: {slide.notes}")

        # 5. Duplicate slide
        new_slide = slide.duplicate()
        print(f"Duplicated slide title: '{new_slide.title}'")
        new_slide.title = "I am the copy"

        # 6. Save
        pres.save(output_path)
        print(f"Saved to {output_path}")

    print("Success! Slide proxy objects verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
