# Changelog

All notable changes to this project are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Apache-2.0 attribution and license text for upstream implementation portions.
- `DefaultCodeLen` and `MaxCodeLen` constants.
- Pinned gosec and deadcode checks in `make ci`.
- Semantic versioning compatibility notes in the README.

### Changed

- Require Go 1.27.4 and update the pinned lint tools. Staticcheck runs again
  with Go 1.27 support.
- Move the module to `git.quad4.io/libraries/olc-go`.
- `make ci` now runs formatting, lint, staticcheck, gosec, deadcode, race,
  bench, and fuzz checks.
- `EncodeTo` accepts buffers as small as `EncodedLen(codeLen)` instead of
  requiring `MaxEncodedLen`.

### Fixed

- Valid overlong short codes recover without panicking. Full overlong codes
  pass through `RecoverNearest` unchanged.
- Encode boundary snapping covers only two floating-point steps, so values
  more than two floating-point steps below a boundary stay in the lower cell.
- The example CLI decode action returns after decoding.
- `Shorten` and `RecoverNearest` measure longitude across the antimeridian, so
  nearby references on opposite sides of 180/-180 degrees recover correctly.

### Removed

- Remove the GitHub-only `.github` metadata and workflows.

## [0.1.0] - 2026-09-17

### Added

- Dependency-free Open Location Code library: Encode, Decode, IsValid,
  IsShort, IsFull, Shorten, RecoverNearest.
- Zero-allocation variants: EncodeTo, AppendEncode, AppendShorten,
  AppendRecoverNearest, EncodeBytes.
- Detailed validation errors via Check, CheckShort, CheckFull.
- Official OLC CSV corpus tests, round-trip fuzzing, allocation guards.
- Example CLI under example/.
- Minimum Go version is 1.27.1.
