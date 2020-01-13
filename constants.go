// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

// DuplicatesThreshold is the number of files of the same file size needed
// before duplicate validation logic is applied.
const DuplicatesThreshold int = 2

// SizeThreshold is the minimum size in bytes that a file must be before
// it is added to our FileSizeIndex. This is an attempt to limit index entries
// to just the files that are most relevant; the assumption is that zero-byte
// files are not relevant.
const SizeThreshold int64 = 1
