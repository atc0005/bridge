// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// multiValueFlag is a custom type that satisfies the flag.Value interface in
// order to accept multiple values for some of our flags
type multiValueFlag []string

// String returns a comma separated string consisting of all slice elements
func (i *multiValueFlag) String() string {

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if i == nil {
		return ""
	}

	return strings.Join(*i, ",")
}

// Set is called once by the flag package, in command line order, for each
// flag present
func (i *multiValueFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Config represents the application configuration as specified via
// command-line flags
type Config struct {

	// Paths represents the various paths checked for duplicate files
	Paths multiValueFlag

	// RecursiveSearch indicates whether paths are crawled recursively or
	// treated as single level directories
	RecursiveSearch bool

	// ConsoleReport indicates whether the (rough) equivalent of our output
	// CSV file is dumped to console. This is primarily intended for
	// troubleshooting.
	ConsoleReport bool

	// IgnoreErrors indicates whether the application should proceed with
	// execution whenever possible by ignoring minor errors. This flag does
	// not affect handling of fatal errors such as failure to generate output
	// report files.
	IgnoreErrors bool

	// FileSizeThreshold is the minimum size in bytes that a file must be
	// before it is added to our FileSizeIndex. This is an attempt to limit
	// index entries to just the files that are most relevant; the assumption
	// is that zero-byte files are not relevant, but the user may wish to
	// limit the threshold to a specific size (e.g., DVD ISO images)
	FileSizeThreshold int64

	// FileDuplicatesThreshold is the number of files of the same file size
	// needed before duplicate validation logic is applied.
	FileDuplicatesThreshold int

	// CSVFile is the fully-qualified path to a CSV file that this application
	// should generate
	CSVFile string

	// ExcelFile is the fully-qualified path to an Excel file that this
	// application should generate
	ExcelFile string
}

// Validate verifies all struct fields have been provided acceptable values
func (c Config) Validate() error {

	if c.Paths == nil {
		return fmt.Errorf("one or more paths not provided")
	}

	if c.FileSizeThreshold < 0 {
		return fmt.Errorf("0 bytes is the minimum size for evaluated files")
	}

	switch {
	case c.CSVFile == "":
		return fmt.Errorf("missing fully-qualified path to CSV file to create")
	case !PathExists(filepath.Dir(c.CSVFile)):
		return fmt.Errorf("parent directory for specified CSV file to create does not exist")
	}

	// Optional flag, optional file generation
	if c.ExcelFile != "" {
		if !PathExists(filepath.Dir(c.ExcelFile)) {
			return fmt.Errorf("parent directory for specified Excel file to create does not exist")
		}
	}

	// RecursiveSearch is a boolean flag. The flag package takes care of
	// assigning a usable default value, so nothing to do here.
	//
	// ConsoleReport is another boolean flag.

	// Optimist
	return nil

}
