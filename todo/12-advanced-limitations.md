# 12 - Advanced Limitations

Scope: document the runtime constraints, legacy surfaces, and explicit unsupported behavior called out in the Aspose.Slides Python docs subset.

## Supported

- Standalone engine: Microsoft PowerPoint is not required for create, modify, convert, or render flows.
- Platforms: Windows 32/64-bit, macOS, and 64-bit Linux are supported.
- Python: the docs call out Python 3.5+.
- Linux/macOS runtime prerequisites: GCC 6 runtime libraries or later, `libgdiplus`, .NET Core runtime dependencies, and a usable `libpython` shared library.
- Packaging: Aspose.Slides extensions are standard Python C extensions and can be frozen with PyInstaller or cx_Freeze.
- PyInstaller: bundled hooks detect and copy dependencies when the import is visible; hidden imports may still need explicit inclusion.
- cx_Freeze: include the root package so the extension and dependent modules are copied into the build.
- Multithreading: parallel work is acceptable when each thread/process uses its own `Presentation` instance.
- ActiveX: existing controls can be accessed, modified, and removed through `IControlCollection`.
- VBA: VBA projects can be added, removed, extracted, and queried for password protection.
- Modern API migration path: thumbnail/image APIs moved from `get_thumbnail` / `get_thumbnails` to `get_image` / `get_images`.
- AI: AI features are exposed through `SlidesAIAgent` and `IAIWebClient` implementations, including the built-in `OpenAIWebClient`.

## Explicitly Not Supported

- A single `Presentation` instance must not be loaded, saved, or cloned from multiple threads.
- Passing live `Presentation` or `Slide` objects between threads is not recommended.
- The library does not execute VBA macros.
- Macros are not transferred when converting a macro-enabled presentation to non-VBA formats such as PDF or HTML.
- Saving a macro-enabled file as PPTX removes macros because PPTX does not support VBA.
- Modern API migration deprecates `aspose.pydrawing`-based imaging APIs; `Graphics`-based members are slated for removal.
- The API does not allow overriding the `Application` metadata in PPTX or the `Creator` / `Producer` metadata in PDF exports.
- Image generation is not yet supported by the AI presentation generator.

## Unclear / Needs Confirmation

- Cloud-platform guidance appears to describe usage examples and deployment context, but it does not define a separate hosted/cloud runtime product in this subset.
- The local `known-issues` page is present, but this dump does not expose concrete issue entries beyond navigation metadata.
- The `api-reference` page confirms module coverage, but it does not itself define feature parity or unsupported calls.

## Source Pages Used

- `TODO.md`
- `web_sid/docs.aspose.com/slides/python-net/api-limitations/index.html`
- `web_sid/docs.aspose.com/slides/python-net/modern-api/index.html`
- `web_sid/docs.aspose.com/slides/python-net/compatibility-with-pyinstaller/index.html`
- `web_sid/docs.aspose.com/slides/python-net/multithreading/index.html`
- `web_sid/docs.aspose.com/slides/python-net/ai/index.xml`
- `web_sid/docs.aspose.com/slides/python-net/activex/index.html`
- `web_sid/docs.aspose.com/slides/python-net/developer-guide/index.html`
- `web_sid/docs.aspose.com/slides/python-net/presentation-via-vba/index.html`
- `web_sid/docs.aspose.com/slides/python-net/slides-on-cloud-platforms/index.html`
- `web_sid/docs.aspose.com/slides/python-net/automating-powerpoint-generation-on-cloud-platforms/index.html`
- `web_sid/docs.aspose.com/slides/python-net/product-support/index.html`
- `web_sid/docs.aspose.com/slides/python-net/system-requirements/index.html`
- `web_sid/docs.aspose.com/slides/python-net/technical-articles/index.html`
- `web_sid/docs.aspose.com/slides/python-net/api-reference/index.html`
