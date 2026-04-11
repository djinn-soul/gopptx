# Core Presentation

Scope: `web_sid/docs.aspose.com/slides/python-net/` pages for presentation setup and the basic presentation lifecycle. Support status below is derived from the repo-wide [`TODO.md`](../TODO.md): `[x]` = Supported, `[ ]` = Not supported, anything not covered there = Unclear.

## Setup

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Installation / environment setup | `installation/`, `system-requirements/`, `getting-started/` | Unclear | Present in the docs tree, but not tracked in `TODO.md`. |
| Licensing / metered licensing | `licensing/`, `metered-licensing/` | Unclear | Present in the docs tree, but not tracked in `TODO.md`. |
| Platform/cloud notes | `slides-on-cloud-platforms/`, `automating-powerpoint-generation-on-cloud-platforms/` | Unclear | Present in the docs tree, but not tracked in `TODO.md`. |
| Evaluate product / known issues | `evaluate-aspose-slides/`, `known-issues/`, `faq/` | Unclear | Present in the docs tree, but not tracked in `TODO.md`. |

## Presentation Lifecycle

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Create blank presentation | `create-presentation/` | Supported | `Presentation.new(title)` is marked `[x]` in `TODO.md`. |
| Create from template | `create-presentation/` | Supported | `Presentation.from_template(path, context)` is marked `[x]`. |
| Open from file | `open-presentation/` | Supported | `Presentation(path)` / `.open(path)` is marked `[x]`. |
| Open from bytes | `open-presentation/` | Supported | `.open_bytes(data)` is marked `[x]`. |
| Context manager support | `open-presentation/` | Supported | `with Presentation(...) as prs:` is marked `[x]`. |
| Save to file | `save-presentation/` | Supported | `.save(path)` is marked `[x]`. |
| Save to bytes | `save-presentation/` | Supported | `.to_bytes()` is marked `[x]`. |
| Merge presentations | `merge-presentation/` | Supported | `merge_from_file(path)` and `merge_from_editor(other)` are marked `[x]`. |
| Presentation properties | `presentation-properties/` | Supported | Core properties, validation, repair, modify password, mark as final, and digital signature checks are marked `[x]`. |
| Presentation security | `presentation-security/` | Supported | Password protection / final-state flows are marked `[x]` in `TODO.md`. |
| Import from ODP / PPT | `import-presentation/`, `ppt-vs-pptx/` | Not supported | Both import items are marked `[ ]` in `TODO.md`. |
| Export to ODP / PPT | `convert-openoffice-odp/`, `convert-ppt-to-pptx/`, `convert-pptx-to-ppt/` | Not supported | ODP/PPT conversion items are marked `[ ]` in `TODO.md`. |
| Convert to image / video / XAML | `convert-powerpoint-to-png/`, `convert-powerpoint-to-jpg/`, `convert-powerpoint-to-video/`, `export-to-xaml/` | Not supported | These export targets are marked `[ ]` in `TODO.md`. |

## Slide Operations

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Add slides | `add-slide-to-presentation/` | Supported | `add_slide` is marked `[x]`. |
| Access slides | `presentation-slide/` | Supported | `slides`, `slides[i]` are marked `[x]`. |
| Remove slides | `remove-slide-from-presentation/` | Supported | `remove_slide` is marked `[x]`. |
| Clone / duplicate slides | `clone-slides/` | Supported | `duplicate_slide`, `duplicate_slide_after` are marked `[x]`. |
| Move slides | `presentation-slide/` | Supported | `move_slide` is marked `[x]`. |
| Hide / unhide slides | `presentation-slide/` | Supported | `set_slide_hidden` is marked `[x]`. |
| Apply / change slide layouts | `slide-layout/`, `slide-master/` | Supported | `SlideLayoutType` and `update_slide` are marked `[x]`. |
| Change slide size | `slide-size/` | Supported | `set_slide_size`, `with_slide_size` are marked `[x]`. |
| Manage slide masters | `slide-master/` | Supported | `slide_masters` is marked `[x]`. |
| Manage slide transitions | `slide-transition/` | Supported | `transitions` module is marked `[x]`. |
| Manage slide sections | `slide-section/` | Supported | `add_section`, `remove_section`, `rename_section` are marked `[x]`. |
| Compare slides | `compare-slides/` | Not supported | Comparison is marked `[ ]` in `TODO.md`. |
| Convert slides to images | `convert-powerpoint-to-png/`, `convert-powerpoint-to-jpg/` | Not supported | Per-slide image export is marked `[ ]` in `TODO.md`. |

## Coverage Summary

- Supported core items are concentrated in presentation creation/open/save/merge, properties/security, and the basic slide container operations.
- Not supported items in this slice are mainly import/export conversions and slide comparison/image rendering.
- Unclear items are mostly setup and installation-adjacent docs that exist in the Aspose tree but are not tracked in `TODO.md`.
