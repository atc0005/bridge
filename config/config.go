// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package config provides types and functions to collect, validate and apply
// user-provided settings.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atc0005/bridge/paths"
)

// ErrInvalidSubcommand represents cases where the user did not pass a valid
// subcommand
var ErrInvalidSubcommand = errors.New("expected 'prune' or 'report' subcommands")

const pruneSubcommand string = "prune"
const reportSubcommand string = "report"

// InputCSVFieldCount represents the number of expected fields when processing
// an input file previously generated by this application for file removal
// decision logic. This value is enforced by the CSV Reader object that
// processes the CSV input file.
// TODO: Find a better place to root this value
const InputCSVFieldCount int = 6

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

	// OutputCSVFile is the fully-qualified path to a CSV file that this application
	// should generate
	OutputCSVFile string

	// OutputCSVFile is the fully-qualified path to a CSV file that this application
	// should use for file removal decisions
	InputCSVFile string

	// ExcelFile is the fully-qualified path to an Excel file that this
	// application should generate
	ExcelFile string

	// BlankLineBetweenSets controls whether a blank line is added between
	// each set of matching files in console and file output.
	BlankLineBetweenSets bool

	// DryRun allows simulation of file removal behavior.
	DryRun bool

	// PruneFiles enables file removal based on provided input CSV file
	PruneFiles bool

	// BackupDirectory is writable directory path where files should be
	// relocated instead of removed
	BackupDirectory string

	// UseFirstRow enables attempts to use the first row from the input CSV
	// file. This should rarely be needed since the input CSV files previously
	// generated by this application contain a header row, but support for
	// overriding this behavior is provided in an effort to support edge cases
	UseFirstRow bool
}

// NewConfig is a factory function that produces a new Config object based
// on user provided flag values.
func NewConfig() (*Config, error) {

	// https://gobyexample.com/command-line-subcommands

	// The subcommand is expected as the first argument to the program.
	if len(os.Args) < 2 {
		return nil, ErrInvalidSubcommand
	}

	config := Config{}

	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
	reportCmd.Var(&config.Paths, "path", "Path to process. This flag may be repeated for each additional path to evaluate.")
	reportCmd.Int64Var(&config.FileSizeThreshold, "size", 1, "File size limit (in bytes) for evaluation. Files smaller than this will be skipped.")
	reportCmd.IntVar(&config.FileDuplicatesThreshold, "duplicates", 2, "Number of files of the same file size needed before duplicate validation logic is applied.")
	reportCmd.BoolVar(&config.RecursiveSearch, "recurse", false, "Perform recursive search into subdirectories per provided path.")
	reportCmd.BoolVar(&config.ConsoleReport, "console", false, "Dump (approximate) CSV file equivalent to console.")
	reportCmd.BoolVar(&config.IgnoreErrors, "ignore-errors", false, "Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.")
	reportCmd.StringVar(&config.OutputCSVFile, "csvfile", "", "The (required) fully-qualified path to a CSV file that this application should generate.")
	reportCmd.StringVar(&config.ExcelFile, "excelfile", "", "The (optional) fully-qualified path to an Excel file that this application should generate.")

	pruneCmd := flag.NewFlagSet("prune", flag.ExitOnError)
	pruneCmd.BoolVar(&config.DryRun, "dry-run", false, "Pretend to remove files, echo what would have been done to stdout. Setting this false does not enable file removal.")
	pruneCmd.BoolVar(&config.PruneFiles, "prune", false, "Enable file removal behavior. This option requires that the input CSV file be specified.")
	pruneCmd.BoolVar(&config.BlankLineBetweenSets, "blank-line", false, "Add a blank line between sets of matching files in console and file output.")
	pruneCmd.StringVar(&config.InputCSVFile, "input-csvfile", "", "The fully-qualified path to a CSV file that this application should use for file removal decisions.")
	pruneCmd.StringVar(&config.BackupDirectory, "backup-dir", "", "The writable directory path where files should be relocated instead of removing them. The original path structure will be created starting with the specified path as the root.")
	pruneCmd.BoolVar(&config.ConsoleReport, "console", false, "Dump (approximate) CSV file equivalent to console.")
	pruneCmd.BoolVar(&config.IgnoreErrors, "ignore-errors", false, "Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.")
	pruneCmd.BoolVar(&config.UseFirstRow, "use-first-row", false, "Attempt to use the first row of the input file. Normally this row is skipped since it is usually the header row and not duplicate file data.")

	// For every subcommand, we parse its own flags and have access to trailing positional arguments.
	switch os.Args[1] {
	case pruneSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s'\n", pruneSubcommand)
		pruneCmd.Parse(os.Args[2:])
	case reportSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s'\n", reportSubcommand)
		reportCmd.Parse(os.Args[2:])
	default:
		return nil, ErrInvalidSubcommand
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate verifies all struct fields have been provided acceptable values
func (c Config) Validate() error {

	// We enforce some common validation requirements for all subcommands
	// and then specific requirements as applicable to each subcommand.
	switch os.Args[1] {

	case pruneSubcommand:

		// DEBUG
		fmt.Printf("validating subcommand '%s'\n", pruneSubcommand)

		// FIXME: The PathExists checks are currently duplicated here and within
		// matches package
		// NOTE: Checking at this point is (potentially) cheaper than waiting
		// until later and then attempting to read in the file. Optional flag, but
		// if set we require that the path actually exist
		if c.InputCSVFile != "" {
			if !paths.PathExists(c.InputCSVFile) {
				return fmt.Errorf("specified CSV file to process does not exist")
			}
		}

		// FIXME: The PathExists checks are currently duplicated here and within
		// matches package
		// NOTE: Checking at this point is cheaper than waiting until later and
		// then attempting to write out the file.
		// Optional flag, optional file backups
		if c.BackupDirectory != "" {
			if !paths.PathExists(c.BackupDirectory) {
				return fmt.Errorf("directory for backup files does not exist: %q", c.BackupDirectory)
			}
		}

	case reportSubcommand:

		// DEBUG
		fmt.Printf("validating subcommand '%s'\n", reportSubcommand)

		if c.Paths == nil {
			return fmt.Errorf("one or more paths not provided")
		}

		if c.FileSizeThreshold < 0 {
			return fmt.Errorf("0 bytes is the minimum size for evaluated files")
		}

		if c.FileDuplicatesThreshold < 2 {
			return fmt.Errorf("2 is the minimum duplicates number for evaluated files")
		}

		// FIXME: The PathExists checks are currently duplicated here and within
		// matches package
		// NOTE: Checking at this point is cheaper than waiting until later and
		// then attempting to write out the file.
		switch {
		case c.OutputCSVFile == "":
			return fmt.Errorf("missing fully-qualified path to CSV file to create")
		case !paths.PathExists(filepath.Dir(c.OutputCSVFile)):
			return fmt.Errorf("parent directory for specified CSV file to create does not exist")
		}

		// FIXME: The PathExists checks are currently duplicated here and within
		// matches package
		// NOTE: Checking at this point is cheaper than waiting until later and
		// then attempting to write out the file.
		// Optional flag, optional file generation
		if c.ExcelFile != "" {
			if !paths.PathExists(filepath.Dir(c.ExcelFile)) {
				return fmt.Errorf("parent directory for specified Excel file to create does not exist")
			}
		}

	default:
		// NOTE: This default case statement should not be reached due to
		// NewConfig() applying the same set of subcommand checks, but
		// providing this step for completeness.
		return ErrInvalidSubcommand
	}

	// TODO: Examine boolean flags for illogical groupings
	// Contrived example:
	//
	// * Enable logging
	// * Disable all logging

	// Optimist
	return nil

}
