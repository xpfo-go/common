# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-03-31

### Added
- GitHub Actions CI workflow (`go test`, `go vet`, `golangci-lint`).
- GitHub Actions release workflow for tag-based release publishing.
- Repository governance documents: `CONTRIBUTING.md`, `SECURITY.md`, issue templates, PR template.
- Redis integration test gating via environment variables:
  - `COMMON_RUN_REDIS_INTEGRATION`
  - `COMMON_REDIS_ADDR`
  - `COMMON_REDIS_PASSWORD`

### Changed
- Module path upgraded to `github.com/xpfo-go/common/v2` for semantic import versioning.
- `limiter` distributed constructors now return explicit errors:
  - `NewDTokenBucket(...) (*DTokenBucket, error)`
  - `NewDLeakyBucket(...) (*DLeakyBucket, error)`
- `limiter` distributed Lua scripts are now loaded with `go:embed` (removed init-time filesystem dependency).
- `limiter` local bucket time accounting upgraded to nanosecond precision to improve deterministic behavior.
- `rsa` APIs now return explicit errors:
  - `Generate(bits int) error`
  - `Encrypt(plainText []byte, publicKeyPath string) ([]byte, error)`
  - `Decrypt(cipherText []byte, privateKeyPath string) ([]byte, error)`
- `rsa` tests now use temporary directories and no longer mutate repository key files.

### Fixed
- Removed init-time panic risk in distributed limiter initialization under non-package working directories.
- Removed flaky limiter tests based on wall-clock sleep timing.

### Breaking Changes
- Import path changed to `/v2`.
- Distributed limiter constructors and RSA APIs now require explicit error handling.

## [1.x]
- Historical versions before semantic import versioning.
