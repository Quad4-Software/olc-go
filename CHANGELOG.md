# Changelog

All notable changes to this project are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Changed

- Minimum Go version is now 1.27.1.
- Standalone staticcheck removed from CI and `make ci`; staticcheck v0.8.1
  does not support Go 1.27. golangci-lint still runs the honnef.co SA checks.

## [0.1.0] - 2026-09-17

### Added

- Dependency-free Open Location Code library: Encode, Decode, IsValid,
  IsShort, IsFull, Shorten, RecoverNearest.
- Zero-allocation variants: EncodeTo, AppendEncode, AppendShorten,
  AppendRecoverNearest, EncodeBytes.
- Detailed validation errors via Check, CheckShort, CheckFull.
- Official OLC CSV corpus tests, round-trip fuzzing, allocation guards.
- Example CLI under example/.
