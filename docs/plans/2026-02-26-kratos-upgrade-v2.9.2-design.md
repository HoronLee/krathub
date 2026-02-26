# Design Doc: Kratos Upgrade to v2.9.2

## 1. Introduction
This document outlines the rationale and process for ensuring Krathub is fully aligned with go-kratos/kratos v2.9.2.

## 2. Current Status
- **Current Version**: v2.9.2 (verified in `go.mod`)
- **Target Version**: v2.9.2 (Stable)

## 3. Rationale for Alignment
The upgrade to v2.9.2 brings several critical improvements:
- **Stability**: Fixes gRPC StreamMiddleware initialization issues.
- **Accuracy**: Fixes `google.protobuf.Empty` type generation.
- **Reliability**: Improves HTTP error reporting and Metadata cloning safety.
- **Modernization**: Full support for Go 1.25 and improved Consul integration.

## 4. Design Sections

### 4.1 Dependency Management
We will ensure all `contrib` packages and the core framework are pinned to v2.9.2 to avoid version mismatch issues.
- Update `go.mod` core: `github.com/go-kratos/kratos/v2 v2.9.2`
- Update `contrib` components: Consul, Kubernetes, Registry, Config.

### 4.2 Code Generation
Since v2.9.0 introduced deprecation markers in HTTP code generation and v2.9.2 fixed `Empty` type handling, we must re-run code generation to benefit from these fixes.
- Execute: `make gen`

### 4.3 Validation
- Run unit tests: `make test`
- Verify gRPC/HTTP servers start correctly.

## 5. Implementation Steps
1. Synchronize `go.mod` via `go get -u github.com/go-kratos/kratos/v2@v2.9.2`.
2. Run `go mod tidy`.
3. Re-generate all Protobuf and Wire code: `make gen`.
4. Run all service tests to ensure no regressions.

## 6. Conclusion
Krathub is already on v2.9.2. This design doc serves as a record of the migration validation and ensures future consistency.
