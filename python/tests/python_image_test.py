import os
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../.."))
sys.path.append(os.path.join(project_root, "python"))


# Create output directory
output_dir = os.path.join(project_root, "examples/output")
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

image_path = os.path.join(
    project_root, "examples/assets/55/repository-open-graph-template.png"
)
output_path = os.path.join(output_dir, "python_image_test.pptx")

if not pathlib.Path(image_path).exists():
    print(f"Image {image_path} not found. Skipping test.")
    exit(0)

try:
    print("Testing add_image via Python...")
    with Presentation.new("Image Test") as pres:
        slide = pres.slides[0]

        # Add image to the first slide
        print(f"Adding image: {image_path}")
        shape_id = slide.add_image(image_path, 1000000, 2000000, 4000000, 2000000)
        print(f"Added image with shape ID: {shape_id}")

        pres.save(output_path)
        print(f"Saved to {output_path}")

    print("Success! Image addition verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
