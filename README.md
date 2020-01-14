# Bridge

A small CLI utility used to find duplicate files

[![Latest Release](https://img.shields.io/github/release/atc0005/bridge.svg?style=flat-square)](https://github.com/atc0005/bridge/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/bridge?status.svg)](https://godoc.org/github.com/atc0005/bridge)
![Validate Codebase](https://github.com/atc0005/bridge/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/bridge/workflows/Validate%20Docs/badge.svg)

- [Bridge](#bridge)
  - [Project home](#project-home)
  - [Overview](#overview)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
  - [How to install it](#how-to-install-it)
  - [Configuration Options](#configuration-options)
    - [Command-line Arguments](#command-line-arguments)
  - [Examples](#examples)
    - [Single path, recursive](#single-path-recursive)
    - [Multiple paths, non-recursive](#multiple-paths-non-recursive)
    - [Invalid flag](#invalid-flag)
  - [License](#license)
    - [Core project files](#core-project-files)
    - [`ByteCountSI`, `ByteCountIEC` functions](#bytecountsi-bytecountiec-functions)
  - [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/bridge) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

Find duplicate files and report them for (currently manual) resolution.

## Features

- Fast and efficient evaluation of potential duplicates by limiting checksum
  generation to two or more identically sized files
- Support for creating CSV report of all duplicate file matches
- Support for generating (rough) console equivalent of CSV file for
  (potential) quick review
- Support for evaluating one or many paths
- Recursive or shallow directory evaluation
- Go modules (vs classic `GOPATH` setup)

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

- Go 1.12+ (for building)
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

Tested using:

- Go 1.13+
- Windows 10 Version 1903
  - native
  - WSL
- Ubuntu Linux 16.04+

## How to install it

1. [Download](https://golang.org/dl/) Go
1. [Install](https://golang.org/doc/install) Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/bridge`
   1. `cd bridge`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - for current operating system with default `go` build options
     - `go build`
   - for all supported platforms
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems that need to run it
   1. Linux: `/tmp/bridge/bridge`
   1. Windows: `/tmp/bridge/bridge.exe`

## Configuration Options

### Command-line Arguments

| Option          | Required | Default        | Repeat | Possible                            | Description                                                                                                                                  |
| --------------- | -------- | -------------- | ------ | ----------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       | `false`        | No     | `0+`                                | Keep specified number of matching files.                                                                                                     |
| `console`       | No       | `false`        | No     | `true`, `false`                     | Dump CSV file equivalent to console.                                                                                                         |
| `csvfile`       | Yes      | *empty string* | No     | *valid file name characters*        | The fully-qualified path to a CSV file that this application should generate.                                                                |
| `ignore-errors` | No       | `false`        | No     | `true`, `false`                     | Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files. |
| `path`          | Yes      | *empty string* | Yes    | *one or more valid directory paths* | Path to process. This flag may be repeated for each additional path to evaluate.                                                             |
| `recurse`       | No       | `false`        | No     | `true`, `false`                     | Perform recursive search into subdirectories per provided path.                                                                              |

## Examples

### Single path, recursive

This example illustrates using the application to process a single path,
recursively.

```ShellSession
./bridge.exe -recurse -path "/tmp/path1"
```

### Multiple paths, non-recursive

This example illustrates using the application to process multiple paths,
without recursively evaluating any subdirectories.

```ShellSession
./bridge.exe -path "/tmp/path1" -path "/tmp/path2"
```

### Invalid flag

Accidentally typing the wrong flag results in a message like this one:

```ShellSession
flag provided but not defined: -fake-flag
```

## License

### Core project files

From the [LICENSE](LICENSE) file:

```license
MIT License

Copyright (c) 2020 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### `ByteCountSI`, `ByteCountIEC` functions

These utility functions are provided by **Stefan Nilsson** under the
**Attribution 3.0 Unported (CC BY 3.0)** license. See the **References** section
of this document for links to additional information.

## References

- <https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/>
  - <https://yourbasic.org/golang/byte-count.go>
  - <https://creativecommons.org/licenses/by/3.0/>

- <https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang>
- <https://golang.org/pkg/flag/#Value>

- <https://stackoverflow.com/questions/50324612/merge-maps-in-golang/50325337#50325337>
- <https://yourbasic.org/golang/gotcha-change-value-range/>

- <https://www.digitalocean.com/community/tutorials/understanding-defer-in-go>
- <https://golangcode.com/writing-to-file/>
- <https://www.joeshaw.org/dont-defer-close-on-writable-files/>
- <https://golang.org/pkg/os/#File.Sync>
