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
	"os"
	"sort"
	"text/tabwriter"
)

// FileMatch represents a superset of statistics (including os.FileInfo) for a
// file matched by provided search criteria. This allows us to record the
// original full path while also recording file metadata used in later
// calculations.
type FileMatch struct {

	// File metadata used in various calculations
	os.FileInfo

	// The full path to the file
	FullPath string

	// Directory containing the file; analogue to Name() method
	ParentDirectory string

	// Checksum calculated for files meeting DuplicatesThreshold
	Checksum SHA256Checksum
}

// FileMatches is a slice of FileMatch objects that represents the search
// results based on user-specified criteria.
type FileMatches []FileMatch

// FileSizeIndex is an index of files based on their size (in bytes) to
// FileMatches. This data structure represents search results for duplicate
// files based on user-specified criteria before we confirm that multiple
// files of the same size are in fact duplicates. In many cases (e.g., a
// multi-part archive), they may not be.
type FileSizeIndex map[int64]FileMatches

// FileChecksumIndex is an index of files based on their checksums (SHA256
// hash) to FileMatches. This data structure is created from a pruned
// FileSizeIndex. After additional pruning to remove any single-entry
// FileMatches "values", this data structure represents confirmed duplicate
// files.
type FileChecksumIndex map[SHA256Checksum]FileMatches

// DuplicateFilesSummary is a collection of the metadata calculated from
// evaluating duplicate files. This metadata is displayed via a variety of
// methods, notably just prior to application exit via console and (in the
// future) the first sheet in the generated workbook.
type DuplicateFilesSummary struct {
	TotalEvaluatedFiles int

	// Number of sets based on identical file size
	FileSizeMatchSets int

	// Number of sets based on identical file hash
	FileHashMatchSets int

	// Identical files count based on file size
	FileSizeMatches int

	// Identical files count based on file hash
	FileHashMatches int

	// Wasted space for duplicate file sets in bytes
	WastedSpace int64
}

// TotalFileSize returns the cumulative size of all files in the slice in bytes
func (fm FileMatches) TotalFileSize() int64 {

	var totalSize int64

	for _, file := range fm {

		totalSize += file.Size()
	}

	return totalSize

}

// TotalFileSizeHR returns a human-readable string of the cumulative size of
// all files in the slice of bytes
func (fm FileMatches) TotalFileSizeHR() string {
	return ByteCountIEC(fm.TotalFileSize())
}

// SizeHR returns a human-readable string of the size of a FileMatch object.
func (fm FileMatch) SizeHR() string {
	return ByteCountIEC(fm.Size())
}

// InList is a helper function to emulate Python's `if "x"
// in list:` functionality
func InList(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// SortByModTimeAsc sorts slice of FileMatch objects in ascending order with
// older values listed first.
func (fm FileMatches) SortByModTimeAsc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().Before(fm[j].ModTime())
	})
}

// SortByModTimeDesc sorts slice of FileMatch objects in descending order with
// newer values listed first.
func (fm FileMatches) SortByModTimeDesc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().After(fm[j].ModTime())
	})
}

// MergeFileSizeIndexes receives one or more FileSizeIndex objects and merges entries
// between these objects, returning a combined FileSizeIndex object
func MergeFileSizeIndexes(fileSizeIndexes ...FileSizeIndex) FileSizeIndex {

	mergedFileSizeIndex := make(FileSizeIndex)

	//log.Printf("Received %d FileSizeIndex objects", len(fileSizeIndexes))

	// loop over all received FileSizeIndex objects, then out of each FileSizeIndex
	// object loop over each attached FileMatches object in order to append
	// each FileMatch in the FileMatches (slice) to our combined object
	//for counter, fileSizeIndex := range fileSizeIndexes {
	for _, fileSizeIndex := range fileSizeIndexes {

		//log.Printf("length of FileSizeIndex %d: %d", counter, len(fileSizeIndex))

		for fileSize, fileMatches := range fileSizeIndex {

			//log.Printf("length of FileMatches for key %d: %d", fileSize, len(fileMatches))

			// From golangci-lint:
			// matches.go:150:4: should replace loop with mergedFileSizeIndex[fileSize] = append(mergedFileSizeIndex[fileSize], fileMatches...) (S1011)
			mergedFileSizeIndex[fileSize] = append(mergedFileSizeIndex[fileSize], fileMatches...)
			// for _, fileMatch := range fileMatches {
			// 	mergedFileSizeIndex[fileSize] = append(mergedFileSizeIndex[fileSize], fileMatch)
			// }
		}
	}

	//log.Printf("mergedFileSizeIndex length: %d", len(mergedFileSizeIndex))

	return mergedFileSizeIndex
}

// UpdateChecksums generates checksum values for each file tracked by a
// FileMatch entry and updates the associated FileMatch.Checksum field value
func (fm FileMatches) UpdateChecksums() error {

	var err error

	// loop over each FileMatch object and generate a checksum
	// https://yourbasic.org/golang/gotcha-change-value-range/
	for index, file := range fm {

		// DEBUG
		//log.Println("Generating checksum for:", file.FullPath)
		fm[index].Checksum, err = GenerateCheckSum(file.FullPath)
		if err != nil {
			return err
		}

		// log.Printf("[%d] Checksum for %s: %s",
		// 	index, fullFileName, fm[index].Checksum)

	}

	// Relying on nil Zero value
	return err
}

// GetCSVRow returns a string slice for use with a CSV Writer
func (fm FileMatch) GetCSVRow() []string {
	return []string{fm.ParentDirectory, fm.Name(), fm.SizeHR(), fm.Checksum.String()}
}

// PruneFileSizeIndex removes map entries with single-entry slices which do
// not reflect potential duplicate files (i.e., duplicate file size !=
// duplicate files)
func (fi FileSizeIndex) PruneFileSizeIndex() {

	for key, fileMatches := range fi {

		// every key is a file size
		// every value is a slice of files of that file size

		// Remove any FileMatches objects that would not contain duplicates
		if len(fileMatches) < DuplicatesThreshold {
			delete(fi, key)
		}
	}
}

// GetTotalFilesCount returns the total number of files in a
// checksum-based file index
func (fi FileSizeIndex) GetTotalFilesCount() int {

	var files int

	for _, fileMatches := range fi {
		files += len(fileMatches)
	}

	return files
}

// PruneFileChecksumIndex removes map entries with single-entry slices which
// do not reflect duplicate files.
func (fi FileChecksumIndex) PruneFileChecksumIndex() {

	for key, fileMatches := range fi {

		// every key is a file checksum
		// every value is a slice of files of that file checksum

		// Remove any FileMatches objects not composed of duplicate checksums
		if len(fileMatches) < DuplicatesThreshold {

			// DEBUG level troubleshooting
			//
			// fmt.Println("Removing key:", key)
			//
			// for _, fileMatch := range fileMatches {
			// 	fmt.Println(fileMatch.GetCSVRow())
			// }

			delete(fi, key)
		}
	}
}

// GetTotalFilesCount returns the total number of files in a
// checksum-based file index
func (fi FileChecksumIndex) GetTotalFilesCount() int {

	var files int

	for _, fileMatches := range fi {
		files += len(fileMatches)
	}

	return files
}

// GetWastedSpace calculates the wasted space from all confirmed duplicate
// files
func (fi FileChecksumIndex) GetWastedSpace() (int64, error) {
	var wastedSpace int64

	// Loop over each duplicate file set in the file checksum index
	// Get count of duplicate file set, minus 1 for the original
	// Get file size in bytes of first entry in that duplicate file set
	// Multiply file size by earlier count of duplicate file set
	// Append cumulative file size of the set (minus original file)
	for _, fileMatches := range fi {

		duplicateFileMatchEntries := (len(fileMatches) - 1)

		// FIXME: This shouldn't be reachable
		if len(fileMatches) == 0 {
			return 0, fmt.Errorf("attempted to calculate wasted space of empty duplicate file set")
		}

		fileSize := fileMatches[0].Size()
		wastedSpace += int64(duplicateFileMatchEntries) * fileSize
	}

	return wastedSpace, nil
}

// GetDuplicateFilesCount returns the number of non-original files in a
// checksum-based file index
func (fi FileChecksumIndex) GetDuplicateFilesCount() int {

	var duplicateFiles int

	for _, fileMatches := range fi {
		duplicateFiles += (len(fileMatches) - 1)
	}

	return duplicateFiles
}

// WriteFileMatches writes duplicate files recorded in a FileChecksumIndex for
// to the specified CSV file. See also PrintFileMatches.
func (fi FileChecksumIndex) WriteFileMatches(filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//w := csv.NewWriter(os.Stdout)
	w := csv.NewWriter(file)

	csvHeader := []string{"directory", "file", "size", "checksum"}

	if err := w.Write(csvHeader); err != nil {
		// at this point we're still trying to write to a non-flushed buffer,
		// so any failures are highly unexpected
		// TODO: Wrap error
		return err
	}

	//for key, fileMatches := range fi {
	for _, fileMatches := range fi {

		for _, file := range fileMatches {
			if err := w.Write(file.GetCSVRow()); err != nil {
				// TODO: Use error wrapping instead?
				return fmt.Errorf("error writing record to csv: %v", err)
			}
		}

	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		// TODO: Wrap the error here with the context?
		// TODO: We accept that CSV Write() or Flush() errors are returned
		// here and not file closure or Sync errors?
		return err
	}

	// TODO: How to return errors from CSV package AND any potential errors
	// from attempting to close the file handle?
	return file.Sync()
}

// PrintFileMatches prints duplicate files recorded in a FileChecksumIndex to
// stdout for development or troubleshooting purposes. See also
// WriteFileMatches for the expected production output method.
func (fi FileChecksumIndex) PrintFileMatches() {

	w := new(tabwriter.Writer)
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	// Format in tab-separated columns
	w.Init(os.Stdout, 8, 8, 5, '\t', 0)

	for _, fileMatches := range fi {

		// Header row in output
		fmt.Fprintln(w, "Directory\tFile\tSize\tChecksum\t")

		for _, file := range fileMatches {

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", file.ParentDirectory, file.Name(), file.SizeHR(), file.Checksum)
			fmt.Fprintln(w)
			w.Flush()
		}

	}

}
