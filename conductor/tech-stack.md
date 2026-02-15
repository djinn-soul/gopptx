# Technology Stack - gopptx

## Core Engine (Go)
- **Version:** Go 1.25
- **Key Dependencies:**
  - `github.com/richardlehane/mscfb`: For handling Compound File Binary Format (CFBF).
  - `golang.org/x/text`: For character encoding and text processing.
- **Role:** High-performance PPTX generation, XML processing, and concurrency management.

## Interop Layer (C Bridge)
- **Mechanism:** Go `cgo` and `-buildmode=c-shared`.
- **Artifacts:** Platform-specific shared libraries (`.dll`, `.so`, `.dylib`).
- **Role:** Providing a stable, low-level interface for other languages to communicate with the Go engine.

## Python Bindings
- **Version:** Python 3.7+ (Recommended 3.12 for development)
- **Build System:** `setuptools`
- **Type Checking:** `basedpyright`
- **Role:** Providing an idiomatic, high-level API for Python developers to utilize the high-performance Go engine.

## Infrastructure & Tooling
- **CI/CD:** GitHub Actions (inferred from `.github/workflows`).
- **Automation:** PowerShell scripts (e.g., `scripts/build_python.ps1`) for cross-platform build automation.
- **Environment:** Container-ready (via `Dockerfile`).
