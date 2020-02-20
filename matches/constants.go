// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package matches provides types and functions intended to help with
// collecting and validating file search results against required criteria.
package matches

// CSV header names referenced from both inside and outside of the package
const (
	CSVDirectoryColumnHeaderName            string = "directory"
	CSVFileColumnHeaderName                 string = "file"
	CSVSizeColumnHeaderName                 string = "size"
	CSVSizeInBytesDirectoryColumnHeaderName string = "size_in_bytes"
	CSVChecksumColumnHeaderName             string = "checksum"
	CSVRemoveFileColumnHeaderName           string = "remove_file"
)
