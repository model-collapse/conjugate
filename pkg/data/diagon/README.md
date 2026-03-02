# Diagon CGO Bindings

This package provides Go bindings to the Diagon C++ search engine library via CGO.

## Architecture

**Diagon Repository** (100% C++):
- Location: `../../../src/3rdparty/diagon` (in CONJUGATE main repository)
- Contains: C++ implementation + C API
- C API headers: `../../../src/3rdparty/diagon/src/core/include/diagon/c_api/`
- Built library: `build/libdiagon.so` (local build cache)

**CONJUGATE Repository** (Go + CGO):
- Go bindings: `*.go` files in this directory
- CGO configuration: Imports C API from `../../../src/3rdparty/diagon/src/core/include/`
- No C++ bridge code (all C API in Diagon upstream)

## Directory Structure

```
pkg/data/diagon/
├── README.md              # This file
├── build/                 # Local build cache
│   └── libdiagon.so       # Built library (linked from src/3rdparty/diagon/build)
├── analysis_go.go         # Pure Go text analysis (default)
├── analysis_cgo.go        # CGO text analysis (optional, requires cgo_analysis tag)
├── analysis_test.go       # Tests for analysis
├── bridge.go              # Main Go-to-C++ bridge (CGO directives point to src/3rdparty/diagon)
├── bridge_test.go         # Bridge tests
└── ...

src/3rdparty/diagon/       # Diagon C++ library (in CONJUGATE repository)
├── src/core/
│   ├── include/diagon/c_api/
│   │   └── diagon_c_api.h       # Main C API
│   └── src/
│       └── ...                  # C++ implementation
└── build/
    └── libdiagon.so             # Built library

NO C++ CODE IN pkg/data/diagon - All C++ code is in src/3rdparty/diagon
```

## Building

The C++ library is built via CMake in the Diagon directory:

```bash
cd src/3rdparty/diagon
cmake -B build -S . -DCMAKE_BUILD_TYPE=Release
cmake --build build -j$(nproc)
```

The Go bindings automatically link against the built library via CGO directives:

```bash
go build ./cmd/data
```

CGO configuration in `bridge.go` points directly to `src/3rdparty/diagon`:
```go
/*
#cgo CFLAGS: -I${SRCDIR}/../../../src/3rdparty/diagon/src/core/include
#cgo LDFLAGS: -L${SRCDIR}/build -ldiagon -L${SRCDIR}/../../../src/3rdparty/diagon/build/src/core -ldiagon_core ...
#include "diagon/c_api/diagon_c_api.h"
*/
```

**Key Points:**
- CGO directives use `${SRCDIR}` which expands to `pkg/data/diagon`
- Relative path `../../../src/3rdparty/diagon` navigates from `pkg/data/diagon` to `src/3rdparty/diagon`
- No symlinks needed - direct path references in CGO configuration

## Design Principles

1. **No C++ in pkg/data/diagon**: All C++ code belongs in src/3rdparty/diagon
2. **C API Boundary**: Go never calls C++ directly, always via C API
3. **Minimal Bridge**: Keep Go bindings thin - just type conversion
4. **Proper Layering**: Diagon = library, CONJUGATE = application
5. **Direct Path References**: CGO directives use relative paths, no symlinks needed

## Migration History

**February 1, 2026**: Direct path references (no symlinks)
- Removed `upstream/` symlink (was pointing to `src/3rdparty/diagon`)
- Updated CGO directives to use direct relative paths: `../../../src/3rdparty/diagon`
- Diagon now in CONJUGATE main repository at `src/3rdparty/diagon`
- No external dependencies or symlink workarounds needed

**January 27, 2026**: Completed architecture cleanup
- Moved `diagon_c_api.{h,cpp}` from Quidditch to Diagon upstream
- Moved `MatchAllQuery` to Diagon core
- Removed obsolete stub implementations (`minimal_wrapper`)
- Removed `c_api_src/` directory (empty after cleanup)
- Bridge layer reduced from 2,263 lines to 0 lines (100% reduction)
- All C API functionality now properly in Diagon library

**Result**: Clean separation - Diagon = 100% C++, CONJUGATE = Go + CGO bindings

## See Also

- [Repository Architecture](../../../REPOSITORY_ARCHITECTURE.md) - Defines boundaries between Diagon and CONJUGATE
- [Architecture Cleanup Plan](../../../ARCHITECTURE_CLEANUP_PLAN.md) - Migration plan executed
- [Diagon README](../../../src/3rdparty/diagon/README.md) - Diagon C++ library documentation
- [CLAUDE.md](../../../CLAUDE.md) - Build instructions and project overview
