// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package dupesets collects types, functions and methods for working with
// collections of duplicate files.
package dupesets

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/atc0005/bridge/internal/checksums"
	"github.com/atc0005/bridge/internal/paths"
	"github.com/atc0005/bridge/internal/units"
)

// Tabwriter header names displayed in console output
const (
	TabWriterDirectoryColumnHeaderName  string = "Directory"
	TabWriterFileColumnHeaderName       string = "File"
	TabWriterSizeColumnHeaderName       string = "Size"
	TabWriterChecksumColumnHeaderName   string = "Checksum"
	TabWriterRemoveFileColumnHeaderName string = "Remove"
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

// DuplicateFileSetEntries is a collection of DuplicateFileSetEntry objects.
// These objects represent rows in a CSV file containing metadata for
// duplicate file sets previously detected and reported by this application.
// If flagged (by the user), the files noted by these entries may be
// (optionally) backed up and removed.
type DuplicateFileSetEntries []DuplicateFileSetEntry

// Print writes DuplicateFileSetEntry objects to a provided Writer, falling
// back to stdout if not specified.
func (dfsEntries DuplicateFileSetEntries) Print(addSeparatorLine bool) {

	w := &tabwriter.Writer{}
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	// Format in tab-separated columns
	// w.Init(os.Stdout, 16, 8, 8, '\t', 0)
	w.Init(os.Stdout, 8, 8, 4, '\t', 0)

	// NOTE: Skip outputing size in bytes since this is meant to be reviewed
	// by a human and not programatically acted upon
	headerRow := fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s",
		TabWriterDirectoryColumnHeaderName,
		TabWriterFileColumnHeaderName,
		TabWriterSizeColumnHeaderName,
		TabWriterChecksumColumnHeaderName,
		TabWriterRemoveFileColumnHeaderName,
	)
	_, _ = fmt.Fprintln(w, headerRow)

	var lastChecksum checksums.SHA256Checksum
	var entriesCtr int
	for _, row := range dfsEntries {

		// FIXME: Call row.Checksum.String() for comparison instead of using
		// this counter? We can always pull the length of the entries by
		// using len() builtin.
		entriesCtr++

		_, _ = fmt.Fprintf(w,
			"%v\t%v\t%v\t%v\t%v\n",
			row.ParentDirectory,
			row.Filename,
			row.SizeHR,
			row.Checksum,
			row.RemoveFile,
		)

		// if user requested a blank line between file sets, look at the
		// checksum for the last entry and compare against the current
		// checksum. A match indicates we are still processing files of the
		// same set, so do not add a blank line. Also, skip adding a blank
		// line for the first item.
		if addSeparatorLine && entriesCtr != 1 {
			if lastChecksum != row.Checksum {
				_, _ = fmt.Fprintf(w, "\n")
			}
		}

		// record current checksum for comparison at the top of the next loop
		// iteration.
		lastChecksum = row.Checksum

	}

	_, _ = fmt.Fprintln(w)
	if err := w.Flush(); err != nil {
		log.Printf(
			"error occurred flushing tabwriter: %v",
			err,
		)
	}

}

// FilesToRemove returns a dfsEntries object representing the files that the
// user has flagged for removal
func (dfsEntries DuplicateFileSetEntries) FilesToRemove() DuplicateFileSetEntries {

	var filesToRemove DuplicateFileSetEntries
	for _, entry := range dfsEntries {
		if entry.RemoveFile {
			filesToRemove = append(filesToRemove, entry)
		}
	}

	return filesToRemove
}

// UpdateSizeInfo fills in potentially missing size information for each entry
// in the duplicate file set.
func (dfsEntry *DuplicateFileSetEntry) UpdateSizeInfo() error {

	// How to best handle nil receiver?
	if dfsEntry == nil {
		return fmt.Errorf("nil receiver; nothing to update")
	}

	// the parseInputRow function places a zero value here if it was found to
	// be empty in the CSV input file row. If it wasn't empty, an attempt
	// was made to convert whatever was present into an int64.
	fileFullPath := filepath.Join(dfsEntry.ParentDirectory, dfsEntry.Filename)
	if dfsEntry.SizeInBytes == 0 {
		// Recalculate the size in bytes from a file that has passed
		// checksum validation
		fileInfo, err := os.Stat(fileFullPath)
		if err != nil {
			return fmt.Errorf(
				"unable to stat %q to determine size in bytes: %w",
				fileFullPath,
				err,
			)
		}
		dfsEntry.SizeInBytes = fileInfo.Size()
	}

	// Update human-readable size field if not already set
	if dfsEntry.SizeHR == "" {
		dfsEntry.SizeHR = units.ByteCountIEC(dfsEntry.SizeInBytes)
	}

	return nil

}

// ValidateInputRow performs basic validation steps against fields in a
// DuplicateFileSetEntry to determine whether an input CSV row will be
// processed further
func ValidateInputRow(dfsEntry DuplicateFileSetEntry, rowNum int) error {

	if !paths.PathExists(dfsEntry.ParentDirectory) {
		return fmt.Errorf(
			"row %d, field %d has invalid parent directory path", rowNum, 0)
	}

	// Filename field
	// TODO: What to check here? We have already enforced non-empty field value
	// during parsing.

	// join ParentDirectory and Filename and check whether the fully-qualified
	// path to the file exists
	fileFullPath := filepath.Join(dfsEntry.ParentDirectory, dfsEntry.Filename)
	if !paths.PathExists(fileFullPath) {
		return fmt.Errorf(
			"row %d, has invalid path to file: %q", rowNum, fileFullPath)
	}

	// Now that we know the parent directory exists and the full path to the
	// file exists, verify checksum before proceeding further
	if err := dfsEntry.Checksum.Verify(fileFullPath); err != nil {
		return fmt.Errorf(
			"checksum validation failed for %q: %w", fileFullPath, err)
	}

	// TODO: Any validation needed against the RemoveFile field? At this point
	// we are not trying to decide whether the file should be removed, just
	// whether the DuplicateFileSetEntry object is properly constructed.

	// Optimism!
	return nil
}

// ParseInputRow evaluates each row returned from the CSV Reader returning a
// DuplicateFileSetEntry object if parsing succeeds, otherwise returning nil.
func ParseInputRow(row []string, fieldCount int, rowNum int) (DuplicateFileSetEntry, error) {

	// TODO: Use error wrapping extensively in this function

	dfsEntry := DuplicateFileSetEntry{}
	var err error

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

	// Go ahead and trim space from all fields
	for index, field := range row {
		row[index] = strings.TrimSpace(field)
	}

	// ParentDirectory
	if row[0] == "" {
		return dfsEntry,
			fmt.Errorf("row %d, field %d has empty parent directory path", rowNum, 1)
	}

	// Filename
	if row[1] == "" {
		return dfsEntry,
			fmt.Errorf("row %d, field %d has empty filename", rowNum, 2)
	}

	// Do not require that this field be populated. We do not have an
	// immediate use for the value in this field and it was mostly used to
	// display the size value as a human-readable value for the report.
	if row[2] == "" {
		log.Printf("DEBUG | CSV row %d, field %d: %q\n", rowNum, 3, row[2])
	}

	// This field is optional; we can regenerate the value later when needed
	// if the row field is empty, we end up with the zero value that we can
	// later check against.
	var sizeInBytes int64
	if row[3] != "" {
		sizeInBytes, err = strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			log.Printf("DEBUG | CSV row %d, field %d: %q\n", rowNum, 4, row[3])
			return dfsEntry, fmt.Errorf("failed to convert CSV sizeInBytes field: %w", err)
		}
	}

	// Checksum
	// Required. We require the checksum to be present so that we can confirm
	// that the file to be removed matches the original checksum recorded for
	// it. If we allowed an empty checksum here, we remove any protection
	// against removing non-duplicate files.
	if row[4] == "" {
		return dfsEntry,
			fmt.Errorf("row %d, field %d has empty checksum", rowNum, 5)
	}

	// Optional field, use default zero value of false if not set
	var removeFile bool
	if row[5] != "" {
		removeFile, err = strconv.ParseBool(row[5])
		if err != nil {
			log.Printf("DEBUG | CSV row %d, field %d: %q\n", rowNum, 6, row[5])
			return dfsEntry, fmt.Errorf("failed to convert CSV remove_file field: %w", err)
		}
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
