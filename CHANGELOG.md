# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/bridge/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.4.4] - 2020-08-22

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable"
      - "stable"
      - "unstable"
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

- Add new README badges for additional CI workflows
  - each badge also links to the associated workflow results

### Changed

- Disable `golangci-lint` default exclusions

- dependencies
  - `go.mod` Go version
    - updated from `1.13` to `1.14`
  - upgrade `360EntSecGroup-Skylar/excelize`
    - `v2.2.0` to `v2.3.0`
  - upgrade `actions/setup-go`
    - `v2.1.1` to `v2.1.2`
      - since replaced with Docker containers
  - upgrade `actions/setup-node`
    - `v2.1.0` to `v2.1.1`
  - upgrade `actions/checkout`
    - `v2.3.1` to `v2.3.2`

- README
  - Link badges to applicable GitHub Actions workflows results

- Linting
  - Local
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
          version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

### Fixed

- Multiple linting issues exposed when disabling `exclude-use-default` setting

## [v0.4.3] - 2020-07-19

### Added

- Dependabot
  - enable version updates
  - enable GitHub Actions updates

### Changed

- Dependencies
  - upgrade `360EntSecGroup-Skylar/excelize`
    - `v2.1.0` to `v2.2.0`
  - upgrade `actions/setup-go`
    - `v1` to `v2.1.0`
  - upgrade `actions/checkout`
    - `v1` to `v2.3.1`
  - upgrade `actions/setup-node`
    - `v1` to `v2.1.0`

### Fixed

- Fix unintentional license change
  - copy/paste config file from another project with a different license
- Fix CHANGELOG section order

## [v0.4.2] - 2020-04-30

### Changed

- README
  - Update README to list accurate build/deploy steps based on recent
    restructuring work

- Update dependencies
  - `stretchr/testify`
    - `v1.4.0` to `v1.5.1`
  - `gopkg.in/yaml.v2`
    - `v2.2.4` to `v2.2.8`
  - `360EntSecGroup-Skylar/excelize`
    - `v1.4.1` to `v2.1.0`
    - Worth noting: the API changed from v1 to v2, so our use of the library
      changed slightly to accommodate those changes

- Vendor dependencies

- Makefile
  - include `-mod=vendor` flag force builds to use new `vendor`
    top-level directory
  - replace two external shell scripts with equivalent embedded commands
  - borrow heavily from existing `Makefile` for `atc0005/elbow` project
  - dynamically determine go module path for version tag use

- Update GitHub Actions Workflows
  - Disable running `go get` after checking out code
  - Exclude `vendor` folder from ...
    - Markdown linting checks
    - tests
    - basic build
  - include `-mod=vendor` flag force builds to use new `vendor`
    top-level directory

- Linting
  - golangci-lint
    - Install and use specific binary version instead of building from  master
    - Move linters/settings to external config file
    - Enable `gofmt` linter
    - Enable `scopelint` linter
    - Enable `dogsled` linter

### Fixed

- GoDoc formatting issue due to forced line wrapping

## [v0.4.1] - 2020-03-02

### Fixed

- (GH-55) `Makefile` builds failed to set version information

## [v0.4.0] - 2020-02-27

### Added

- Add support for pruning flagged/marked items in the input CSV file (previously
  generated by the `report` subcommand)
- Split application logic (and flags) into subcommands
  - `report` for existing behavior and set of flags
  - `prune` for new behavior and new set of flags

### Changed

- GitHub Actions Workflow: `Validate Codebase`
  - `Build with default options` step updated to run `go build` against cmd
    dir path
  - Go 1.12.x removed from build matrix
  - Go 1.14.x added to build matrix
- Move related chunks of code into subpackages
  - e.g., `matches`, `paths`, ...
- Help/Usage output
  - Emit extended Help/Usage information for each subcommand
  - Emit overall summary of subcommands when binary is called without
    subcommands or with `-h` or `-help` flags
  - Emit branding details (`App Name`, `Version`, `Repo URL`)

### Fixed

- README coverage for help flags

## [v0.3.0] - 2020-02-09

### Added

- Echo Go version used in CI workflows so that it is saved in CI output logs
- Flag for duplicates threshold
- Flag for size threshold

### Fixed

- Add missing (and required) `csvfile` flag in README examples
- Add missing guard against creation of Microsoft Excel file when user did not
  request it
- Emphasize that the `csvfile` flag is required, `excelfile` flag is optional
- Miscellaneous docs cleanup

## [v0.2.0] - 2020-01-15

### Added

- Support for creating Microsoft Excel workbook of all duplicate file matches
- README
  - CI badges to indicate current linting and build results
  - GoDoc badge
  - Latest release badge

### Fixed

- Ignore release assets generated by Makefile builds

## [v0.1.1] - 2020-01-13

### Fixed

- Missing support in multiple locations for `IgnoreErrors` option

## [v0.1.0] - 2020-01-13

### Added

This initial prototype supports/provides:

- Fast and efficient evaluation of potential duplicates by limiting checksum
  generation to two or more identically sized files
- Support for creating CSV report of duplicate file matches
- Support for generating (rough) console equivalent of CSV file for
  (potential) quick review
- Support for evaluating one or many paths
- Recursive or single-level directory evaluation

Worth noting (in no particular order):

- Command-line flags support via `flag` standard library package
- Go modules (vs classic `GOPATH` setup)
- GitHub Actions linting and build checks
- Makefile for general use cases
- No external, non-standard library packages

[Unreleased]: https://github.com/atc0005/bridge/compare/v0.4.4...HEAD
[v0.4.4]: https://github.com/atc0005/bridge/releases/tag/v0.4.4
[v0.4.3]: https://github.com/atc0005/bridge/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/bridge/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/bridge/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/bridge/releases/tag/v0.4.0
[v0.3.0]: https://github.com/atc0005/bridge/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/bridge/releases/tag/v0.2.0
[v0.1.1]: https://github.com/atc0005/bridge/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/bridge/releases/tag/v0.1.0
