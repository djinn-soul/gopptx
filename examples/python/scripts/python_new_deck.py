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

output_path = os.path.join(output_dir, "python_from_scratch.pptx")

try:
    print("Creating new presentation from scratch...")
    with Presentation.new("Hello from Python Scratch") as pres:
        print(f"Initial slide count: {pres.slide_count}")

        # Add another slide
        pres.add_slide("Second Slide")
        print(f"New slide count: {pres.slide_count}")

        # Add a shape to the first slide
        pres.add_shape(
            0, "rect", 1000000, 1000000, 2000000, 1000000, text="Born in Python"
        )

        pres.save(output_path)
        print(f"Saved new deck to {output_path}")

    print("Success! Creating decks from scratch verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
