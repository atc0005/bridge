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
// var ErrInvalidSubcommand = fmt.Errorf(
// 	"expected '%s' or '%s' subcommands",
// 	PruneSubcommand,
// 	ReportSubcommand,
// )
var ErrInvalidSubcommand = errors.New("invalid subcommand")

// ErrMissingSubcommand represents cases where the user did not pass a
// subcommand
var ErrMissingSubcommand = errors.New("missing subcommand")

// PruneSubcommand is meant as a label to be easily used/referenced in place
// of the subcommand of the same name.
const PruneSubcommand string = "prune"

// ReportSubcommand is meant as a label to be easily used/referenced in place
// of the subcommand of the same name.
const ReportSubcommand string = "report"

// version is updated via Makefile builds by referencing the fully-qualified
// path to this variable, including the package. We set a placeholder value so
// that something resembling a version string will be provided for
// non-Makefile builds.
var version string = "x.y.z"

const myAppName string = "bridge"
const myAppURL string = "https://github.com/atc0005/bridge"

// TODO: Needed?
var validSubcommands = []string{PruneSubcommand, ReportSubcommand}

// activeFlagSet represents the matching flagset for the options the user
// chose. This is referenced later from Validate() in order to print the
// "default options" for that chosen flagset if/when validation of settings
// fails.
var activeFlagSet *flag.FlagSet

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

// Branding is responsible for emitting application name, version and origin
func Branding() {
	fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s\n%s\n\n", myAppName, version, myAppURL)
}

// MainCmdUsage is meant to be called whenever a valid subcommand is missing
// (not provided or invalid)
func MainCmdUsage(subCmds ...string) func() {

	return func() {

		myBinaryName := filepath.Base(os.Args[0])

		Branding()

		fmt.Fprintln(flag.CommandLine.Output(), "Available subcommands:")
		for _, subCmd := range subCmds {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s\n", subCmd)
		}
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "See available options for each subcommand by running:")
		for _, subCmd := range subCmds {
			fmt.Fprintf(flag.CommandLine.Output(), "\t%s %s -h\n", myBinaryName, subCmd)
		}
		fmt.Fprintln(flag.CommandLine.Output(), "")
	}

}

// SubcommandUsage is a custom override for the default Help text provided by
// the flag package. Here we prepend some additional metadata to the existing
// output.
func SubcommandUsage(flagSet *flag.FlagSet) func() {

	return func() {

		myBinaryName := filepath.Base(os.Args[0])

		Branding()

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"%s %s\":\n",
			myBinaryName,
			flagSet.Name(),
		)
		flagSet.PrintDefaults()

	}
}

// Config represents the application configuration as specified via
// command-line flags
type Config struct {

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

	// BlankLineBetweenSets controls whether a blank line is added between
	// each set of matching files in console and file output.
	BlankLineBetweenSets bool

	// DryRun allows simulation of file removal behavior.
	DryRun bool

	// UseFirstRow enables attempts to use the first row from the input CSV
	// file. This should rarely be needed since the input CSV files previously
	// generated by this application contain a header row, but support for
	// overriding this behavior is provided in an effort to support edge cases
	UseFirstRow bool

	// FileDuplicatesThreshold is the number of files of the same file size
	// needed before duplicate validation logic is applied.
	FileDuplicatesThreshold int

	// FileSizeThreshold is the minimum size in bytes that a file must be
	// before it is added to our FileSizeIndex. This is an attempt to limit
	// index entries to just the files that are most relevant; the assumption
	// is that zero-byte files are not relevant, but the user may wish to
	// limit the threshold to a specific size (e.g., DVD ISO images)
	FileSizeThreshold int64

	// OutputCSVFile is the fully-qualified path to a CSV file that this application
	// should generate
	OutputCSVFile string

	// OutputCSVFile is the fully-qualified path to a CSV file that this application
	// should use for file removal decisions
	InputCSVFile string

	// ExcelFile is the fully-qualified path to an Excel file that this
	// application should generate
	ExcelFile string

	// BackupDirectory is writable directory path where files should be
	// relocated instead of removed
	BackupDirectory string

	// Paths represents the various paths checked for duplicate files
	Paths multiValueFlag
}

// NewConfig is a factory function that produces a new Config object based
// on user provided flag values.
func NewConfig() (*Config, error) {

	config := Config{}

	// Note: We define our flagsets and attempt to parse them, printing usage
	// details to the user and returning an error to main() for handling there
	// instead of forcing an early exit here (credit: Miki Tebeka)

	// Setup an entirely new flagset and set it as the primary for the main
	// application itself (subcommand flagsets defined separately). The intent
	// (as noted in "Learning CoreDNS") is to protect against other imported
	// packages overriding flags set for this application.
	mainFlagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	mainFlagSet.Usage = MainCmdUsage(validSubcommands...)
	flag.CommandLine = mainFlagSet
	flag.Parse()

	// The subcommand is expected as the first argument to the program.
	if len(os.Args) < 2 {
		MainCmdUsage(validSubcommands...)()
		return nil, ErrMissingSubcommand
	}

	reportCmd := flag.NewFlagSet("report", flag.ContinueOnError)
	reportCmd.Var(&config.Paths, "path", "Path to process. This flag may be repeated for each additional path to evaluate.")
	reportCmd.Int64Var(&config.FileSizeThreshold, "size", 1, "File size limit (in bytes) for evaluation. Files smaller than this will be skipped.")
	reportCmd.IntVar(&config.FileDuplicatesThreshold, "duplicates", 2, "Number of files of the same file size needed before duplicate validation logic is applied.")
	reportCmd.BoolVar(&config.RecursiveSearch, "recurse", false, "Perform recursive search into subdirectories per provided path.")
	reportCmd.BoolVar(&config.ConsoleReport, "console", false, "Dump (approximate) CSV file equivalent to console.")
	reportCmd.BoolVar(&config.IgnoreErrors, "ignore-errors", false, "Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.")
	reportCmd.StringVar(&config.OutputCSVFile, "csvfile", "", "The (required) fully-qualified path to a CSV file that this application should generate.")
	reportCmd.StringVar(&config.ExcelFile, "excelfile", "", "The (optional) fully-qualified path to an Excel file that this application should generate.")

	pruneCmd := flag.NewFlagSet("prune", flag.ContinueOnError)
	pruneCmd.BoolVar(&config.DryRun, "dry-run", false, "Don't actually remove files. Echo what would have been done to stdout.")
	pruneCmd.BoolVar(&config.BlankLineBetweenSets, "blank-line", false, "Add a blank line between sets of matching files in console and file output.")
	pruneCmd.StringVar(&config.InputCSVFile, "input-csvfile", "", "The fully-qualified path to a CSV file that this application should use for file removal decisions.")
	pruneCmd.StringVar(&config.BackupDirectory, "backup-dir", "", "The writable directory path where files should be relocated instead of removing them. The original path structure will be created starting with the specified path as the root.")
	pruneCmd.BoolVar(&config.ConsoleReport, "console", false, "Dump (approximate) CSV file equivalent to console.")
	pruneCmd.BoolVar(&config.IgnoreErrors, "ignore-errors", false, "Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.")
	pruneCmd.BoolVar(&config.UseFirstRow, "use-first-row", false, "Attempt to use the first row of the input file. Normally this row is skipped since it is usually the header row and not duplicate file data.")

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand (found at os.Args[1])

	// FIXME: How can we have "-h" and "-help" *not* caught by this switch
	// statement?
	switch os.Args[1] {
	case PruneSubcommand:
		// DEBUG
		fmt.Printf("DEBUG: subcommand '%s'\n", PruneSubcommand)
		pruneCmd.Usage = SubcommandUsage(pruneCmd)
		if err := pruneCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("DEBUG: err returned from pruneCmd.Parse():", err)
			return nil, err
		}
		activeFlagSet = pruneCmd

	case ReportSubcommand:
		// DEBUG
		fmt.Printf("DEBUG: subcommand '%s'\n", ReportSubcommand)
		reportCmd.Usage = SubcommandUsage(reportCmd)
		if err := reportCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("DEBUG: err returned from reportCmd.Parse():", err)
			return nil, err
		}
		activeFlagSet = reportCmd

	// TODO: How can we allow the flag package to deal with this instead of
	// explicitly matching against the flags here? Otherwise the default case
	// statement is used ...
	case "-h", "-help":
		fmt.Println("DEBUG: Help flags used")
		mainFlagSet.PrintDefaults()
		activeFlagSet = nil
		return nil, ErrMissingSubcommand
	default:
		// TODO: Confirm whether MainCmdUsage() is used here automatically or
		// whether we have to call it explicitly
		mainFlagSet.PrintDefaults()
		activeFlagSet = nil
		fmt.Println("DEBUG: default case statement for subcommand")
		return nil, ErrInvalidSubcommand
	}

	// Fallback to our base flagset if a subcommand was not specified
	if activeFlagSet == nil {
		activeFlagSet = mainFlagSet
	}

	if err := config.Validate(activeFlagSet); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate verifies all struct fields have been provided acceptable values
func (c Config) Validate(flagset *flag.FlagSet) error {

	// We enforce some common validation requirements for all subcommands
	// and then specific requirements as applicable to each subcommand.
	switch os.Args[1] {

	case PruneSubcommand:

		// DEBUG
		fmt.Printf("DEBUG: validating subcommand '%s'\n", PruneSubcommand)

		if strings.TrimSpace(c.InputCSVFile) == "" {
			flagset.Usage()
			return fmt.Errorf("required input CSV file to process not specified")
		}

		// c.BackupDirectory is optional; applying length checks here
		// if user provides value would be unreliable. Path exist check
		// is applied later at use point, so not duplicating here as it
		// would be outsid the intent/scope of this function's purpose.

	case ReportSubcommand:

		// DEBUG
		fmt.Printf("DEBUG: validating subcommand '%s'\n", ReportSubcommand)

		if c.Paths == nil {
			flagset.Usage()
			return fmt.Errorf("one or more paths not provided")
		}

		if c.FileSizeThreshold < 0 {
			flagset.Usage()
			return fmt.Errorf("0 bytes is the minimum size for evaluated files")
		}

		if c.FileDuplicatesThreshold < 2 {
			flagset.Usage()
			return fmt.Errorf("2 is the minimum duplicates number for evaluated files")
		}

		// FIXME: The PathExists checks are currently duplicated here and within
		// matches package
		// NOTE: Checking at this point is cheaper than waiting until later and
		// then attempting to write out the file.
		switch {
		case c.OutputCSVFile == "":
			flagset.Usage()
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
