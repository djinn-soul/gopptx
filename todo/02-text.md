# 02 - Text

Source: `TODO.md` plus `web_sid/docs.aspose.com/slides/python-net/{manage-text,manage-textbox,manage-paragraph,text-formatting,wordart}`.

## Supported

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Placeholders | `manage-text/` | Supported | List placeholders, set placeholder content, inspect placeholder metadata, and read slide layout refs are marked `[x]`. |
| Text boxes | `manage-textbox/` | Supported | Add and update textbox text, word wrap, margins, and text direction are marked `[x]`. |
| AutoFit | `manage-text/`, `text-formatting/` | Supported | Shape, normal, and none autofit modes are marked `[x]`. |
| Paragraphs | `manage-paragraph/` | Supported | Alignment, spacing, indent level, and default run formatting are marked `[x]`. |
| Run formatting | `text-formatting/` | Supported | Bold, italic, underline, strikethrough, font name/size/color, all-caps, small-caps, spacing, superscript, subscript, highlight, and run hyperlinks are marked `[x]`. |
| Text animation | `manage-text/` | Supported | Slide-level animation effects and the core entrance/exit/emphasis constants are marked `[x]`. |
| Bullets | `manage-text/`, `manage-paragraph/` | Supported | Bullet text and nested indent levels are marked `[x]`. |
| Superscript / subscript | `manage-text/`, `text-formatting/` | Supported | Run-level superscript and subscript plus the builder API are marked `[x]`. |
| Text extraction and search | `manage-text/` | Supported | Slide text states plus find-and-replace are marked `[x]`. |
| Localization | `manage-text/` | Supported | Presentation language property is marked `[x]`. |
| WordArt vertical direction | `wordart/` | Supported | Vertical and RTL vertical text direction are marked `[x]`. |

## Not Supported

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Placeholder mutation | `manage-text/` | Not supported | Remove, add, and reorder placeholder flows remain `[ ]`. |
| Text box borders / rotation | `manage-textbox/` | Not supported | Independent border/line styling and textbox rotation remain `[ ]`. |
| Read autofit | `text-formatting/` | Not supported | Reading the current autofit state remains `[ ]`. |
| Paragraph management | `manage-paragraph/` | Not supported | Paragraph count, deletion, and reordering remain `[ ]`. |
| Theme font colors / underline styling | `text-formatting/` | Not supported | Theme colors and underline color/style remain `[ ]`. |
| Text-by-word / letter animation | `manage-text/` | Not supported | Trigger and timing controls remain `[ ]`. |
| Bullet and numbering variants | `manage-paragraph/` | Not supported | Custom bullets, numbered lists, list start number, and mixed lists remain `[ ]`. |
| Read superscript / subscript state | `text-formatting/` | Not supported | Reading existing run state remains `[ ]`. |
| Rich text extraction | `manage-text/` | Not supported | Full text, formatting metadata, notes text, and table text extraction remain `[ ]`. |
| Localization extras | `manage-text/` | Not supported | Per-run language tags, RTL per paragraph, bulk translation workflow helpers, and translation export remain `[ ]`. |
| WordArt extras | `wordart/` | Not supported | Transforms, glow/shadow/reflection, and preset styles remain `[ ]`. |
