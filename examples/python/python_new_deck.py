import sys
import os

# Add project root to sys.path to find 'gopptx' package
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.."))
sys.path.append(os.path.join(project_root, "python"))

from gopptx import Presentation

# Create output directory
output_dir = os.path.join(project_root, "examples/output")
os.makedirs(output_dir, exist_ok=True)

output_path = os.path.join(output_dir, "python_from_scratch.pptx")

try:
    print("Creating new presentation from scratch...")
    with Presentation.new("Hello from Python Scratch") as pres:
        print(f"Initial slide count: {pres.slide_count}")
        
        # Add another slide
        pres.add_slide("Second Slide")
        print(f"New slide count: {pres.slide_count}")
        
        # Add a shape to the first slide
        pres.add_shape(0, "rect", 1000000, 1000000, 2000000, 1000000, text="Born in Python")
        
        pres.save(output_path)
        print(f"Saved new deck to {output_path}")

    print("Success! Creating decks from scratch verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback
    traceback.print_exc()
