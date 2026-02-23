import pathlib
import sys
import traceback

from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / ".." / ".." / "..").resolve()
sys.path.append(str(project_root / "python"))
output_dir = project_root / "examples" / "output"
output_dir.mkdir(exist_ok=True, parents=True)
input_path = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"
output_path = output_dir / "python_advanced_output.pptx"
if not input_path.exists():
    print(f"Input file {input_path} not found. Skipping test.")
    sys.exit(0)
try:
    print(f"Opening presentation: {input_path}")
    with Presentation(str(input_path)) as pres:
        print("Listing slides...")
        for slide in pres.slides:
            print(f"  - Index {slide.index}: {slide.title} (ID: {slide.slide_id})")
        print("Testing Find and Replace...")
        count = pres.find_and_replace("Basic", "SUPER")
        print(f"  Made {count} replacements.")
        print("Searching for shapes...")
        results = pres.search_shapes({"text_contains": "SUPER"})
        for res in results:
            print(
                f"  Found shape '{res['Shape']['Name']}' on slide {res['SlideIndex']} with text: '{res['Shape']['Text']}'"
            )
        print("Testing Comments...")
        author_id = pres.add_author("Jane Doe", "JD")
        print(f"  Added author 'Jane Doe' with ID: {author_id}")
        pres.add_comment(
            0, author_id, "This is a comment from Python!", x=100000, y=100000
        )
        print("  Added comment to slide 0.")
        authors = pres.get_authors()
        print("  Current authors:")
        for a in authors:
            print(f"    - {a['Name']} ({a['Initials']}) ID: {a['ID']}")
        comments = pres.get_comments(0)
        print("  Comments on slide 0:")
        for c in comments:
            print(f"    - Author {c['AuthorID']}: {c['Text']} (Index: {c['Index']})")
        pres.save(str(output_path))
        print(f"Saved to {output_path}")
    print("\nSuccess! Advanced Python features verified.")
except (OSError, RuntimeError, ValueError) as e:
    print(f"Error during verification: {e}")
    traceback.print_exc()
