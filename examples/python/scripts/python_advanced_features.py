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

input_path = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")
output_path = os.path.join(output_dir, "python_advanced_output.pptx")

if not pathlib.Path(input_path).exists():
    print(f"Input file {input_path} not found. Skipping test.")
    exit(0)

try:
    print(f"Opening presentation: {input_path}")
    with Presentation(input_path) as pres:
        # 1. List Slides
        print("Listing slides...")
        for slide in pres.slides:
            print(f"  - Index {slide.index}: {slide.title} (ID: {slide.slide_id})")

        # 2. Find and Replace
        print("Testing Find and Replace...")
        count = pres.find_and_replace("Basic", "SUPER")
        print(f"  Made {count} replacements.")

        # 3. Shape Search
        print("Searching for shapes...")
        results = pres.search_shapes({"text_contains": "SUPER"})
        for res in results:
            print(
                f"  Found shape '{res['Shape']['Name']}' on slide {res['SlideIndex']} with text: '{res['Shape']['Text']}'"
            )

        # 4. Comments and Authors
        print("Testing Comments...")
        # Add an author
        author_id = pres.add_author("Jane Doe", "JD")
        print(f"  Added author 'Jane Doe' with ID: {author_id}")

        # Add a comment
        pres.add_comment(
            0, author_id, "This is a comment from Python!", x=100000, y=100000
        )
        print("  Added comment to slide 0.")

        # Get authors
        authors = pres.get_authors()
        print("  Current authors:")
        for a in authors:
            print(f"    - {a['Name']} ({a['Initials']}) ID: {a['ID']}")

        # Get comments
        comments = pres.get_comments(0)
        print("  Comments on slide 0:")
        for c in comments:
            print(f"    - Author {c['AuthorID']}: {c['Text']} (Index: {c['Index']})")

        # 5. Save
        pres.save(output_path)
        print(f"Saved to {output_path}")

    print("\nSuccess! Advanced Python features verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
