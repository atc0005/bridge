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

## [v0.5.1] - 2022-07-17

### Overview

- Build & Release improvements
- Bug fixes
- Dependency updates
- built using Go 1.19.11
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-275) Add initial automated release notes config

### Changed

- Dependencies
  - `Go`
    - `1.19.8` to `1.19.11`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.4` to `go-ci-oldstable-build-v0.11.4`
  - `xuri/excelize`
    - `v2.7.0` to `v2.7.1`
  - `xuri/nfp`
    - `v0.0.0-20220409054826-5e722a1d9e22` to
      `v0.0.0-20230503010013-3f38cdbb0b83`
  - `xuri/efp`
    - `v0.0.0-20220603152613-6918739fd470` to
      `v0.0.0-20230422071738-01f4e37c47e9`
  - `golang.org/x/crypto`
    - `v0.8.0` to `v0.11.0`
  - `golang.org/x/net`
    - `v0.9.0` to `v0.12.0`
  - `golang.org/x/text`
    - `v0.9.0` to `v0.11.0`
- (GH-261) Update vuln analysis GHAW to remove on.push hook
- (GH-276) Releases: Add separate section for dependencies
- (GH-277) Releases: Add 'New Features' section

### Fixed

- (GH-258) Disable depguard linter
- (GH-265) Restore local CodeQL workflow
- (GH-271) Remove plugin deploy logic from postinstall script
- (GH-278) Releases: Update 'Bug Fixes' label

## [v0.5.0] - 2022-04-09

### Overview

- Add support for generating DEB, RPM packages
- Build improvements
- Generated binary changes
  - filename patterns
  - compression (~ 66% smaller)
  - executable metadata
- built using Go 1.19.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-246) Generate RPM/DEB packages using nFPM
- (GH-245) Add version details to Windows executables

### Changed

- (GH-247) Switch to semantic versioning (semver) compatible versioning
  pattern
- (GH-244) Makefile: Compress binaries & use fixed filenames
- (GH-243) Makefile: Refresh recipes to add "standard" set, new
  package-related options
- (GH-242) Build dev/stable releases using go-ci Docker image
- (GH-248) Move internal packages to internal path

### Fixed

- (GH-241) Fix v0.4.16 release summary
- (GH-251) Fix errwrap linting errors

## [v0.4.16] - 2022-04-09

### Overview

- Bug fixes
- Build improvements
- GitHub Actions workflows updates
- Dependency updates
- built using Go 1.19.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-220) Add Go Module Validation, Dependency Updates jobs

### Changed

- Dependencies
  - `Go`
    - `1.19.4` to `1.19.8`
  - `golang.org/x/crypto`
    - `v0.3.0` to `v0.8.0`
  - `golang.org/x/net`
    - `v0.4.0` to `v0.9.0`
  - `golang.org/x/text`
    - `v0.5.0` to `v0.9.0`
  - `360EntSecGroup-Skylar/excelize` renamed to `xuri/excelize/v2`
  - `xuri/excelize/v2`
    - `v2.4.0` to `v2.7.0`
  - `xuri/nfp`
    - `v0.0.0-20220409054826-5e722a1d9e22`
- (GH-227) Drop `Push Validation` workflow
- (GH-228) Rework workflow scheduling
- (GH-230) Remove `Push Validation` workflow status badge

### Fixed

- (GH-208) Spurious "file already closed" errors emitted during app execution
- (GH-212) The `github.com/360EntSecGroup-Skylar/excelize/v2` dependency has
  moved?
- (GH-226) Fix breakage from `xuri/excelize` dependency update
- (GH-235) Update vuln analysis GHAW to use on.push hook

## [v0.4.15] - 2022-12-07

### Overview

- Dependency updates
- built using Go 1.19.4
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.9` to `1.19.4`
  - `golang.org/x/crypto`
    - `v0.0.0-20210415154028-4f45737414dc` to `v0.3.0`
  - `golang.org/x/net`
    - `v0.0.0-20210415231046-e915ea6b2b7d` to `v0.4.0`
  - `golang.org/x/sys`
    - `v0.0.0-20210927094055-39ccf1dd6fa6` to `v0.3.0`
  - `golang.org/x/text`
    - `v0.3.6` to `v0.5.0`
  - `github.com/richardlehane/mscfb`
    - `v1.0.3` to `v1.0.4`
  - `github.com/richardlehane/msoleps`
    - `v1.0.1` to `v1.0.3`
  - `github.com/xuri/efp`
    - `v0.0.0-20210322160811-ab561f5b45e3` to `v0.0.0-20220603152613-6918739fd470`
- (GH-192) Update project to Go 1.19
- (GH-195) Update Makefile and GitHub Actions Workflows
- (GH-201) Refactor GitHub Actions workflows to import logic

### Fixed

- (GH-190) Update lintinstall Makefile recipe
- (GH-193) Apply linting fixes for Go 1.19 release
- (GH-194) Add missing cmd doc file
- (GH-205) Fix Makefile Go module base path detection

## [v0.4.14] - 2022-05-06

### Overview

- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.7` to `1.17.9`

## [v0.4.13] - 2022-03-03

### Overview

- Dependency updates
- CI / linting improvements
- built using Go 1.17.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.6` to `1.17.7`
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

- (GH-175) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-176) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)

### Fixed

- (GH-177) revive, gosec linting errors surfaced by GHAWs refresh

## [v0.4.12] - 2022-01-25

### Overview

- Dependency updates
- Built using Go 1.17.6
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.12` to `1.17.6`
    - (GH-170) Update go.mod file, canary Dockerfile to reflect current
      dependencies

## [v0.4.11] - 2021-12-29

### Overview

- Dependency updates
- Built using Go 1.16.12
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.10` to `1.16.12`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.1`

## [v0.4.10] - 2021-11-10

### Overview

- Dependency updates
- Built using Go 1.16.10
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.8` to `1.16.10`
  - `actions/checkout`
    - `v2.3.4` to `v2.4.0`
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`

### Fixed

- (GH-161) False positive `G307: Deferring unsafe method "Close" on type
  "*os.File" (gosec)` linting error

## [v0.4.9] - 2021-09-27

### Overview

- Dependency updates
- Built using Go 1.16.8
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.8`

## [v0.4.8] - 2021-08-09

### Overview

- Dependency updates
- Built using Go 1.16.7
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - `v2.2.0` to `v2.4.0`

## [v0.4.7] - 2021-07-19

### Overview

- Dependency updates
- Minor fixes
- Built using Go 1.16.6
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- Add "canary" Dockerfile to track stable Go releases, serve as a reminder to
  generate fresh binaries

### Changed

- Dependencies
  - `Go`
    - `1.15.8` to `1.16.6`
  - `360EntSecGroup-Skylar/excelize`
    - `v2.3.2` to `v2.4.0`
  - `actions/setup-node`
    - `v2.1.4` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

### Fixed

- ST1023: should omit type int from declaration; it will be inferred from the
  right-hand side (stylecheck)

## [v0.4.6] - 2021-02-21

### Overview

- Dependency updates
- Minor fixes
- Built using Go 1.15.8

### Changed

- Swap out GoDoc badge for pkg.go.dev badge

- Dependencies
  - Built using Go 1.15.8
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `360EntSecGroup-Skylar/excelize`
    - `v2.3.1` to `v2.3.2`
  - `actions/checkout`
    - `v2.3.3` to `v2.3.4`
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

### Fixed

- Fix explicit exit code handling

## [v0.4.5] - 2020-10-11

### Added

- Binary release
  - Built using Go 1.15.2
  - **Statically linked**
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

### Changed

- Dependencies
  - `360EntSecGroup-Skylar/excelize`
    - `v2.3.0` to `v2.3.1`
  - `actions/checkout`
    - `v2.3.2` to `v2.3.3`
  - `actions/setup-node`
    - `v2.1.1` to `v2.1.2`

- Add `-trimpath` build flag

### Fixed

- Makefile build options do not generate static binaries
- Misc linting errors raised by latest `gocritic` release included with
  `golangci-lint` `v1.31.0`
- Makefile generates checksums with qualified path

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

[Unreleased]: https://github.com/atc0005/bridge/compare/v0.5.1...HEAD
[v0.5.1]: https://github.com/atc0005/bridge/releases/tag/v0.5.1
[v0.5.0]: https://github.com/atc0005/bridge/releases/tag/v0.5.0
[v0.4.16]: https://github.com/atc0005/bridge/releases/tag/v0.4.16
[v0.4.15]: https://github.com/atc0005/bridge/releases/tag/v0.4.15
[v0.4.14]: https://github.com/atc0005/bridge/releases/tag/v0.4.14
[v0.4.13]: https://github.com/atc0005/bridge/releases/tag/v0.4.13
[v0.4.12]: https://github.com/atc0005/bridge/releases/tag/v0.4.12
[v0.4.11]: https://github.com/atc0005/bridge/releases/tag/v0.4.11
[v0.4.10]: https://github.com/atc0005/bridge/releases/tag/v0.4.10
[v0.4.9]: https://github.com/atc0005/bridge/releases/tag/v0.4.9
[v0.4.8]: https://github.com/atc0005/bridge/releases/tag/v0.4.8
[v0.4.7]: https://github.com/atc0005/bridge/releases/tag/v0.4.7
[v0.4.6]: https://github.com/atc0005/bridge/releases/tag/v0.4.6
[v0.4.5]: https://github.com/atc0005/bridge/releases/tag/v0.4.5
[v0.4.4]: https://github.com/atc0005/bridge/releases/tag/v0.4.4
[v0.4.3]: https://github.com/atc0005/bridge/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/bridge/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/bridge/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/bridge/releases/tag/v0.4.0
[v0.3.0]: https://github.com/atc0005/bridge/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/bridge/releases/tag/v0.2.0
[v0.1.1]: https://github.com/atc0005/bridge/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/bridge/releases/tag/v0.1.0
