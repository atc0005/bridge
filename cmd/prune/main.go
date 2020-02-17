// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/atc0005/bridge/checksums"
	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/paths"
)

// DuplicateFileSetEntry represents a duplicate file set entry recorded as a
// row within an input CSV file. This row is expected to contain data
// originally generated by the `report` subcommand and a user-provided flag
// indicating whether one or more files from the duplicate file set are to be
// removed.
// TODO: Any logical way to merge this type definition with FileMatch? Is
// there a thin interface that could be created between the types? This
// type is very close to the existing FileMatch type.
type DuplicateFileSetEntry struct {

	// ParentDirectory represents the directory containing a file from a
	// duplicate file sets
	ParentDirectory string

	// Filename is the name of a file from a duplicate file set
	Filename string

	// SizeHR is the size of a file from a duplicate file set in
	// human-readable text format (e.g., 1 GB, 500 MB)
	SizeHR string

	// SizeInBytes is the size of a file from a duplicate file set in bytes
	SizeInBytes int64

	// Checksum is the file hash for a file from a duplicate file set
	Checksum checksums.SHA256Checksum

	// RemoveFile is a flag indicating whether a file from a duplicate file
	// set is to be removed
	RemoveFile bool
}

// DuplicateFileSetEntries is a collection of DuplicateFileSetEntry objects
// representing rows that should be acted upon in some way, usually file
// pruning actions.
type DuplicateFileSetEntries []DuplicateFileSetEntry

// Print writes DuplicateFileSetEntry objects to a provided Writer, falling
// back to stdout if not specified.
func (dfsEntries DuplicateFileSetEntries) Print() {

	w := &tabwriter.Writer{}
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	// Format in tab-separated columns
	//w.Init(os.Stdout, 16, 8, 8, '\t', 0)
	w.Init(os.Stdout, 8, 8, 4, '\t', 0)

	// NOTE: Skip outputing size in bytes since this is meant to be reviewed
	// by a human and not programatically acted upon

	// Header row in output
	fmt.Fprintln(w,
		"Directory\tFile\tSize\tChecksum\tRemove File")

	for _, row := range dfsEntries {

		fmt.Fprintf(w,
			"%v\t%v\t%v\t%v\t%v\n",
			row.ParentDirectory,
			row.Filename,
			row.SizeHR,
			row.Checksum,
			row.RemoveFile,
		)

	}

	fmt.Fprintln(w)
	w.Flush()
}

// parseInputRow evaluates each row returned from the CSV Reader returning a
// DuplicateFileSetEntry object if parsing succeeds, otherwise returning nil.
func parseInputRow(row []string, fieldCount int, rowNum int) (DuplicateFileSetEntry, error) {

	dfsEntry := DuplicateFileSetEntry{}

	// The CSV Reader already performs field count validation, but let's be
	// paranoid and recheck to help ensure that we didn't make a mistake and
	// pass a different slice than expected.
	if len(row) != fieldCount {
		return dfsEntry, fmt.Errorf(
			"unexpected number of fields received. got %d, expected %d",
			len(row),
			fieldCount,
		)
	}

	// validate parent directory
	parentDirectory := strings.TrimSpace(row[0])
	switch {
	case parentDirectory == "":
		return dfsEntry,
			fmt.Errorf("row %d, field %d has empty parent directory path", rowNum, 0)
	case !paths.PathExists(parentDirectory):
		return dfsEntry,
			fmt.Errorf("row %d, field %d has invalid parent directory path", rowNum, 0)
	}

	sizeInBytes, err := strconv.ParseInt(row[3], 10, 64)
	if err != nil {
		log.Printf("DEBUG | CSV row %d, field 4: %q\n", rowNum, row[3])
		// TODO: Use error wrapping here
		return dfsEntry, fmt.Errorf("failed to convert CSV sizeInBytes field %v", err)
	}

	removeFile, err := strconv.ParseBool(row[5])
	if err != nil {
		log.Printf("DEBUG | CSV row %d, field 6: %q\n", rowNum, row[5])
		// TODO: Use error wrapping here
		return dfsEntry, fmt.Errorf("failed to convert CSV remove_file field %v", err)

	}

	// convert a CSV row into an object representing the various named
	// fields found in that row
	dfsEntry = DuplicateFileSetEntry{
		ParentDirectory: row[0],
		Filename:        row[1],
		SizeHR:          row[2],
		SizeInBytes:     sizeInBytes,
		Checksum:        checksums.SHA256Checksum(row[4]),
		RemoveFile:      removeFile,
	}

	// everything went well
	return dfsEntry, nil

}

func main() {

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		panic(err)
	}

	// DEBUG
	log.Printf("Configuration: %+v\n", appConfig)

	// Values we can work with:
	//
	// config.IgnoreErrors
	// config.DryRun
	// config.BlankLineBetweenSets
	// config.InputCSVFile
	// config.ConsoleReport
	// config.BackupDirectory
	// config.PruneFiles
	// config.SkipFirstRow

	// -------------------------------------------------------------------- //

	/*

		STEPS

		Parse config options
		Open CSV file
			? Buffered reader?
		Create CSV Reader object
		Apply CSV parsing requirements
			- Require specific number of fields
			- Skip blank lines
		Loop over rows
		Validate row fields
			- Field content?
				- e.g., checksum field has a length expected of current
				  hash algorithm
		Verify files exist
		Verify checksum for each file removal candidate
		Backup file removal candidate (if option is set)

	*/

	if !paths.PathExists(appConfig.InputCSVFile) {
		log.Fatal("specified CSV input file does not exist")
	}

	file, err := os.Open(appConfig.InputCSVFile)
	if err != nil {
		log.Fatal(err)
	}
	// NOTE: We're not manipulating contents for this file, so relying solely
	// on a defer statement to close the file should be sufficient?
	defer file.Close()

	csvReader := csv.NewReader(file)

	// Require that the number of fields found matches what we expect to find
	csvReader.FieldsPerRecord = config.InputCSVFieldCount

	// TODO: Even with this set, we should probably still trim whitespace
	// ourselves so that we can be assured that leading AND trailing
	// whitespace has been removed
	csvReader.TrimLeadingSpace = true

	var dfsEntries DuplicateFileSetEntries
	var rowCounter int = 0
	for {

		// Go ahead and bump the counter to reflect that humans start counting
		// CSV rows from 1 and not 0
		rowCounter++

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// If we are currently evaluating the very first line of the CSV file
		// and the user did not override the default option of skipping the
		// first row (due to it usually being the header row)
		if rowCounter == 1 {
			if !appConfig.UseFirstRow {
				// DEBUG
				log.Println("Skipping first row in input file to avoid processing column headers")
				continue
			}
			log.Println("Attempting to parse row 1 from input CSV file as requested")
		}

		// ###########################################################################################################
		// TODO: "row" is not a good representation of the duplicate file set
		// entries recorded in the CSV file. This also applies to csvRows,
		// DuplicateFileSetEntry, and DuplicateFileSetEntries.
		// ###########################################################################################################
		dfsEntry, err := parseInputRow(record, config.InputCSVFieldCount, rowCounter)
		if err != nil {
			if appConfig.IgnoreErrors {
				log.Println("IgnoringErrors set, ignoring input row and continuing with the next one.")
				continue
			}
			log.Fatal("IgnoringErrors NOT set. Exiting.")
		}
		dfsEntries = append(dfsEntries, dfsEntry)

		//printCSVRow(row)

	}

	dfsEntries.Print()

}
