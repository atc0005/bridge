# Bridge

A small CLI utility used to find duplicate files

[![Latest Release](https://img.shields.io/github/release/atc0005/bridge.svg?style=flat-square)](https://github.com/atc0005/bridge/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/bridge?status.svg)](https://godoc.org/github.com/atc0005/bridge)
![Validate Codebase](https://github.com/atc0005/bridge/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/bridge/workflows/Validate%20Docs/badge.svg)

- [Bridge](#bridge)
  - [Project home](#project-home)
  - [Overview](#overview)
    - [Generate report](#generate-report)
    - [Prune duplicate files](#prune-duplicate-files)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
  - [How to install it](#how-to-install-it)
  - [Configuration Options](#configuration-options)
    - [Command-line Arguments](#command-line-arguments)
      - [`report` subcommand](#report-subcommand)
      - [`prune` subcommand](#prune-subcommand)
  - [Examples](#examples)
    - [Generating a report](#generating-a-report)
      - [Single path, recursive](#single-path-recursive)
      - [Multiple paths, non-recursive](#multiple-paths-non-recursive)
      - [Invalid flag](#invalid-flag)
    - [Pruning duplicate files](#pruning-duplicate-files)
      - [Dry-run (minimal)](#dry-run-minimal)
      - [Dry-run (verbose)](#dry-run-verbose)
      - [Backup files before removing them](#backup-files-before-removing-them)
  - [License](#license)
    - [Core project files](#core-project-files)
    - [`ByteCountSI`, `ByteCountIEC` functions](#bytecountsi-bytecountiec-functions)
  - [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/bridge) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

1. Generate report
   - Find duplicate files and report them via console-only output or an output
     CSV file
1. Remove flagged files
   - Process CSV file report generated earlier: if flag is set,
     (optionally) backup and then remove marked files

### Generate report

Generating a report is the first step towards indicating which files from a
duplicate file set that you wish to remove (specified explicitly) and which
you wish to keep (default behavior).

### Prune duplicate files

Pruning duplicate files is an optional second step following the generation of
a duplicate files report (via the `report` subcommand).

You first open the CSV file using an application like Microsoft Excel or
LibreOffice Calc and then mark each file (`remove_file` column) that you wish
to remove with either `true` or `false`; the default is `false`, so marking an
entry with `false` is not strictly necessary.

Once marked, you are then able to remove those files by specifying the full
path to the CSV file (via the `prune` subcommand). See the
[Examples](#examples) section for details.

## Features

- Efficient evaluation of potential duplicates by limiting checksum generation
  to two or more identically sized files
- Support for creating CSV report of all duplicate file matches
- Support for generating (rough) console equivalent of CSV file for
  (potential) quick review
- Support for creating Microsoft Excel workbook of all duplicate file matches
- Support for evaluating one or many paths
- Recursive or shallow directory evaluation
- Optional removal of (user-flagged) duplicate files from a previously
  generated CSV report
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

#### `report` subcommand

| Option          | Required | Default        | Repeat | Possible                            | Description                                                                                                                                  |
| --------------- | -------- | -------------- | ------ | ----------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       | `false`        | No     | `h`, `help`                         | Show Help text along with the list of supported flags.                                                                                       |
| `console`       | No       | `false`        | No     | `true`, `false`                     | Dump (approximate) CSV file equivalent to console.                                                                                           |
| `csvfile`       | Yes      | *empty string* | No     | *valid file name characters*        | The fully-qualified path to a CSV file that this application should generate.                                                                |
| `excelfile`     | No       | *empty string* | No     | *valid file name characters*        | The fully-qualified path to a Microsoft Excel file that this application should generate.                                                    |
| `size`          | No       | `1` (byte)     | No     | `0+`                                | File size limit for evaluation. Files smaller than this will be skipped.                                                                     |
| `duplicates`    | No       | `2`            | No     | `2+`                                | Number of files of the same file size needed before duplicate validation logic is applied.                                                   |
| `ignore-errors` | No       | `false`        | No     | `true`, `false`                     | Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files. |
| `path`          | Yes      | *empty string* | Yes    | *one or more valid directory paths* | Path to process. This flag may be repeated for each additional path to evaluate.                                                             |
| `recurse`       | No       | `false`        | No     | `true`, `false`                     | Perform recursive search into subdirectories per provided path.                                                                              |

#### `prune` subcommand

| Option          | Required | Default        | Repeat | Possible                     | Description                                                                                                                                                                     |
| --------------- | -------- | -------------- | ------ | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       | `false`        | No     | `h`, `help`                  | Show Help text along with the list of supported flags.                                                                                                                          |
| `console`       | No       | `false`        | No     | `true`, `false`              | Dump (approximate) CSV file equivalent to console.                                                                                                                              |
| `dry-run`       | No       | `false`        | No     | `true`, `false`              | Don't actually remove files. Echo what would have been done to stdout.                                                                                                          |
| `ignore-errors` | No       | `false`        | No     | `true`, `false`              | Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.                                    |
| `input-csvfile` | Yes      | *empty string* | No     | *valid file name characters* | The fully-qualified path to a CSV file that this application should use for file removal decisions.                                                                             |
| `backup-dir`    | No       | *empty string* | No     | *valid directory path*       | The writable directory path where files should be relocated instead of removing them. The original path structure will be created starting with the specified path as the root. |
| `blank-line`    | No       | `false`        | No     | `true`, `false`              | Add a blank line between sets of matching files in console and file output.                                                                                                     |
| `use-first-row` | No       | `false`        | No     | `true`, `false`              | Attempt to use the first row of the input file. Normally this row is skipped since it is usually the header row and not duplicate file data.                                    |

## Examples

### Generating a report

#### Single path, recursive

This example illustrates using the application to process a single path,
recursively.

```ShellSession
./bridge.exe report -recurse -path "/tmp/path1" -csvfile "path1-report.csv"
```

#### Multiple paths, non-recursive

This example illustrates using the application to process multiple paths,
without recursively evaluating any subdirectories.

```ShellSession
./bridge.exe report -path "/tmp/path1" -path "/tmp/path2"  -csvfile "report.csv"
```

#### Invalid flag

Accidentally typing the wrong flag results in a message like this one:

```ShellSession
flag provided but not defined: -fake-flag
```

### Pruning duplicate files

#### Dry-run (minimal)

```ShellSession
./bridge.exe prune -input-csvfile "report.csv" -dry-run -ignore-errors
```

Here we specify:

- Don't actually remove files, just simulate the process
- input CSV file (file previously generated by the `report` subcommand)
- ignore (minor) errors

Because the `console` flag wasn't specified, the output is minimal.

#### Dry-run (verbose)

```ShellSession
./bridge.exe prune -input-csvfile "report.csv" -dry-run -ignore-errors -console
```

Here we specify:

- Don't actually remove files, just simulate the process
- input CSV file (file previously generated by the `report` subcommand)
- ignore (minor) errors
- `console` flag
  - enables printing table of parsed CSV contents
  - enables printing table of file removal candidates

Because the `console` flag *was* specified, the output is more verbose.

#### Backup files before removing them

```ShellSession
./bridge.exe prune -input-csvfile "report.csv" -backup-dir /tmp/tacos -dry-run -ignore-errors -console
```

Here we specify:

- the input CSV file (file previously generated by the `report` subcommand)
- the backup directory that should be used to copy files to (just before a
  file removal operation is attempted)
- ignore (minor) errors
- `console` flag
  - enables printing table of parsed CSV contents
  - enables printing table of file removal candidates

Because the `console` flag *was* specified, the output is more verbose. This
can make the removal process easier to troubleshoot due to the explicit
listing of what *would* be removed and what actually occurred.

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
- <https://gobyexample.com/command-line-subcommands>

- <https://stackoverflow.com/questions/50324612/merge-maps-in-golang/50325337#50325337>
- <https://yourbasic.org/golang/gotcha-change-value-range/>

- <https://www.digitalocean.com/community/tutorials/understanding-defer-in-go>
- <https://golangcode.com/writing-to-file/>
- <https://www.joeshaw.org/dont-defer-close-on-writable-files/>
- <https://golang.org/pkg/os/#File.Sync>
- <https://www.linode.com/docs/development/go/creating-reading-and-writing-files-in-go-a-tutorial/>

- <https://medium.com/@sebassegros/golang-dealing-with-maligned-structs-9b77bacf4b97>

- <https://github.com/360EntSecGroup-Skylar/excelize>
