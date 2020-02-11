// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package matches provides types and functions intended to help with
// collecting and validating file search results against required criteria.
package matches

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/atc0005/bridge/checksums"
	"github.com/atc0005/bridge/paths"
	"github.com/atc0005/bridge/units"

	"github.com/360EntSecGroup-Skylar/excelize"
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

	// Checksum calculated for files meeting the duplicates threshold
	Checksum checksums.SHA256Checksum
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
type FileChecksumIndex map[checksums.SHA256Checksum]FileMatches

// DuplicateFilesSummary is a collection of the metadata calculated from
// evaluating duplicate files. This metadata is displayed via a variety of
// methods, notably just prior to application exit via console and the first
// sheet in the generated workbook.
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
	return units.ByteCountIEC(fm.TotalFileSize())
}

// SizeHR returns a human-readable string of the size of a FileMatch object.
func (fm FileMatch) SizeHR() string {
	return units.ByteCountIEC(fm.Size())
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

// UpdateChecksums acts as a wrapper around the UpdateChecksums method for
// FileMatches objects
func (fi FileSizeIndex) UpdateChecksums(ignoreErrors bool) error {

	//for key, fileMatches := range combinedFileSizeIndex {
	for _, fileMatches := range fi {

		// every key is a file size
		// every value is a slice of files of that file size

		if err := fileMatches.UpdateChecksums(ignoreErrors); err != nil {

			// DEBUG
			log.Println("Error encountered:", err)
			if !ignoreErrors {
				return err
			}
			// DEBUG
			log.Println("Ignoring error as requested")
			continue
		}
	}

	// TODO: Return bool and error instead of just error?
	// This would allow returning true as in success, but also
	// provide the original error that we chose to ignore.
	return nil
}

// UpdateChecksums generates checksum values for each file tracked by a
// FileMatch entry and updates the associated FileMatch.Checksum field value
func (fm FileMatches) UpdateChecksums(ignoreErrors bool) error {

	var err error

	// loop over each FileMatch object and generate a checksum
	// https://yourbasic.org/golang/gotcha-change-value-range/
	for index, file := range fm {

		// DEBUG
		//log.Println("Generating checksum for:", file.FullPath)
		result, err := checksums.GenerateCheckSum(file.FullPath)
		if err != nil {

			if !ignoreErrors {
				return err
			}

			// WARN
			log.Println("Error encountered:", err)
			log.Println("Ignoring error as requested")

			continue

		}

		fm[index].Checksum = result

		// log.Printf("[%d] Checksum for %s: %s",
		// 	index, fullFileName, fm[index].Checksum)

	}

	// Relying on nil Zero value
	return err
}

// GenerateCSVHeaderRow returns a string slice for use with a CSV Writer as a
// header row.
func (fi FileChecksumIndex) GenerateCSVHeaderRow() []string {
	return []string{
		"directory",
		"file",
		"size",
		"size_in_bytes",
		"checksum",
		"remove_file",
	}
}

// GenerateEmptyCSVDataRow returns a string slice for use with a CSV Writer as a
// empty data (non-header) row. This is used as a separator between sets of
// duplicate files.
func (fm FileMatches) GenerateEmptyCSVDataRow() []string {
	return []string{
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// GenerateCSVDataRow returns a string slice for use with a CSV Writer as a
// data (non-header) row
func (fm FileMatch) GenerateCSVDataRow() []string {
	return []string{
		fm.ParentDirectory,
		fm.Name(),
		fm.SizeHR(),
		strconv.FormatInt(fm.Size(), 10),
		fm.Checksum.String(),
		"",
	}
}

// NewFileSizeIndex optionally recursively processes a provided path and returns a
// slice of FileMatch objects
func NewFileSizeIndex(recursiveSearch bool, ignoreErrors bool, fileSizeThreshold int64, dirs ...string) (FileSizeIndex, error) {

	combinedFileSizeIndex := make(FileSizeIndex)

	for _, path := range dirs {

		if !paths.PathExists(path) {
			return nil, fmt.Errorf("provided path %q does not exist", path)
		}

		// DEBUG
		log.Println("Path exists:", path)

		// TODO: Call ProcessPath here
		fileSizeIndex, err := ProcessPath(recursiveSearch, ignoreErrors, fileSizeThreshold, path)
		if err != nil {
			return nil, fmt.Errorf("failed to process path %q: %v", path, err)
		}

		// FIXME: This needs to occur at the end of each loop?
		combinedFileSizeIndex = MergeFileSizeIndexes(combinedFileSizeIndex, fileSizeIndex)

	}

	// TODO: Safe to return err here, relying on it being nil if no errors
	// were caught earlier?
	// return combinedFileSizeIndex, err
	return combinedFileSizeIndex, nil

}

// ProcessPath optionally recursively processes a provided path and returns a
// slice of FileMatch objects
func ProcessPath(recursiveSearch bool, ignoreErrors bool, fileSizeThreshold int64, path string) (FileSizeIndex, error) {

	fileSizeIndex := make(FileSizeIndex)
	var err error

	//log.Println("RecursiveSearch:", recursiveSearch)

	if recursiveSearch {

		// Walk walks the file tree rooted at path, calling the anonymous function
		// for each file or directory in the tree, including path. All errors that
		// arise visiting files and directories are filtered by the anonymous
		// function. The files are walked in lexical order, which makes the output
		// deterministic but means that for very large directories Walk can be
		// inefficient. Walk does not follow symbolic links.
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			// If an error is received, check to see whether we should ignore
			// it or return it. If we return a non-nil error, this will stop
			// the filepath.Walk() function from continuing to walk the path,
			// and your main function will immediately move to the next line.
			if err != nil {
				if !ignoreErrors {
					return err
				}

				// WARN
				log.Println("Error encountered:", err)
				log.Println("Ignoring error as requested")

			}

			// make sure we're not working with the root directory itself
			if path != "." {

				// ignore directories
				if info.IsDir() {
					return nil
				}

				// ignore files below the size threshold
				if info.Size() < fileSizeThreshold {
					return nil
				}

				// If we made it to this point, then we must assume that the file
				// has met all criteria to be evaluated by this application.
				// Let's add the file to our slice of files of the same size
				// using our index based on file size.
				fileSizeIndex[info.Size()] = append(
					fileSizeIndex[info.Size()],
					FileMatch{
						FileInfo:        info,
						FullPath:        path,
						ParentDirectory: filepath.Dir(path),
					})
			}

			return err
		})

	} else {

		// If recursiveSearch is not enabled, process just the provided path

		// err is already declared earlier at a higher scope, so do not
		// redeclare here
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)

		if err != nil {
			// TODO: Wrap error?
			log.Printf("Error from ioutil.ReadDir(): %s", err)

			return fileSizeIndex, err
		}

		// Use []os.FileInfo returned from ioutil.ReadDir() to build slice of
		// FileMatch objects
		for _, file := range files {

			// ignore directories
			if file.IsDir() {
				continue
			}

			// ignore files below the size threshold
			if file.Size() < fileSizeThreshold {
				continue
			}

			// If we made it to this point, then we must assume that the file
			// has met all criteria to be evaluated by this application.
			// Let's add the file to our slice of files of the same size
			// using our index based on file size.
			fileSizeIndex[file.Size()] = append(
				fileSizeIndex[file.Size()],
				FileMatch{
					FileInfo: file,
					FullPath: filepath.Join(path, file.Name()),
					// ParentDirectory: filepath.Dir(path),
					// `path` is a flat directory structure (we are not using
					// recursion), so record it directly as the parent
					// directory for files within
					ParentDirectory: path,
				})
		}
	}

	return fileSizeIndex, err
}

// PruneFileSizeIndex removes map entries with single-entry slices which do
// not reflect potential duplicate files (i.e., duplicate file size !=
// duplicate files)
func (fi FileSizeIndex) PruneFileSizeIndex(duplicatesThreshold int) {

	for key, fileMatches := range fi {

		// every key is a file size
		// every value is a slice of files of that file size

		// Remove any FileMatches objects that do not contain a number of
		// duplicate checksums meething our threshold
		if len(fileMatches) < duplicatesThreshold {
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

// NewFileChecksumIndex takes in a FileSizeIndex, generates checksums for
// FileMatch objects and then returns a FileChecksumIndex and an error, if
// one was encountered.
func NewFileChecksumIndex(fi FileSizeIndex) FileChecksumIndex {
	fileChecksumIndex := make(FileChecksumIndex)
	for _, fileMatches := range fi {
		for _, fileMatch := range fileMatches {
			fileChecksumIndex[fileMatch.Checksum] = append(
				fileChecksumIndex[fileMatch.Checksum],
				fileMatch)
		}
	}

	return fileChecksumIndex
}

// PruneFileChecksumIndex removes map entries with single-entry slices which
// do not reflect duplicate files.
func (fi FileChecksumIndex) PruneFileChecksumIndex(duplicatesThreshold int) {

	for key, fileMatches := range fi {

		// every key is a file checksum
		// every value is a slice of files of that file checksum

		// Remove any FileMatches objects that do not contain a number of
		// duplicate checksums meething our threshold
		if len(fileMatches) < duplicatesThreshold {

			// DEBUG level troubleshooting
			//
			// fmt.Println("Removing key:", key)
			//
			// for _, fileMatch := range fileMatches {
			// 	fmt.Println(fileMatch.GenerateCSVDataRow())
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
		// subtract one so that we don't count the original as a duplicate
		duplicateFiles += (len(fileMatches) - 1)
	}

	return duplicateFiles
}

// WriteFileMatchesWorkbook is a prototype method to generate an Excel
// workbook from duplicate file details
func (fi FileChecksumIndex) WriteFileMatchesWorkbook(filename string, summary DuplicateFilesSummary) error {

	if !paths.PathExists(filepath.Dir(filename)) {
		return fmt.Errorf("parent directory for specified CSV file to create does not exist")
	}

	f := excelize.NewFile()

	summarySheet := "Summary"
	defaultSheet := "Sheet1"

	// Create a new sheet for duplicate file metadata
	summarySheetIndex := f.NewSheet(summarySheet)

	f.DeleteSheet(defaultSheet)

	// Write out the summary sheet labels
	f.SetCellValue(summarySheet, "A1", "Evaluated Files")
	f.SetCellValue(summarySheet, "A2", "Sets of files with identical size")
	f.SetCellValue(summarySheet, "A3", "Sets of files with identical fingerprint")
	f.SetCellValue(summarySheet, "A4", "Files with identical size")
	f.SetCellValue(summarySheet, "A5", "Files with identical fingerprint")
	// blank link
	f.SetCellValue(summarySheet, "A7", "Duplicate Files")
	f.SetCellValue(summarySheet, "A8", "Wasted Space")

	// Write out the summary sheet values
	f.SetCellValue(summarySheet, "B1", summary.TotalEvaluatedFiles)
	f.SetCellValue(summarySheet, "B2", summary.FileSizeMatchSets)
	f.SetCellValue(summarySheet, "B3", summary.FileHashMatchSets)
	f.SetCellValue(summarySheet, "B4", summary.FileSizeMatches)
	f.SetCellValue(summarySheet, "B5", summary.FileHashMatches)
	// blank line

	f.SetCellValue(summarySheet, "B7", fi.GetDuplicateFilesCount())
	f.SetCellValue(summarySheet, "B8", units.ByteCountIEC(summary.WastedSpace))

	for duplicateFileSetIndex, fileMatches := range fi {

		//sheetHeader := []string{"directory", "file", "size", "checksum"}

		// Create a new sheet for duplicate file metadata
		f.NewSheet(duplicateFileSetIndex.String())

		// Write out the sheet header
		f.SetCellValue(duplicateFileSetIndex.String(), "A1", "directory")
		f.SetCellValue(duplicateFileSetIndex.String(), "B1", "file")
		f.SetCellValue(duplicateFileSetIndex.String(), "C1", "size")
		f.SetCellValue(duplicateFileSetIndex.String(), "D1", "size in bytes")
		f.SetCellValue(duplicateFileSetIndex.String(), "E1", "checksum")

		for index, file := range fileMatches {

			// Excel starts at 1, but our header occupies row 1, so increment
			// by +2 to account for that
			row := index + 2

			f.SetCellValue(duplicateFileSetIndex.String(), fmt.Sprintf("A%d", row), file.ParentDirectory)
			f.SetCellValue(duplicateFileSetIndex.String(), fmt.Sprintf("B%d", row), file.Name())
			f.SetCellValue(duplicateFileSetIndex.String(), fmt.Sprintf("C%d", row), file.SizeHR())
			f.SetCellValue(duplicateFileSetIndex.String(), fmt.Sprintf("D%d", row), file.Size())
			f.SetCellValue(duplicateFileSetIndex.String(), fmt.Sprintf("E%d", row), file.Checksum.String())

		}

	}

	f.SetActiveSheet(summarySheetIndex)

	// Save xlsx file by the given path.
	return f.SaveAs(filename)
	// if err := f.SaveAs("Book1.xlsx"); err != nil {
	// 	println(err.Error())
	// }
}

// WriteFileMatchesCSV writes duplicate files recorded in a FileChecksumIndex
// to the specified CSV file.
func (fi FileChecksumIndex) WriteFileMatchesCSV(filename string) error {

	if !paths.PathExists(filepath.Dir(filename)) {
		return fmt.Errorf("parent directory for specified CSV file to create does not exist")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//w := csv.NewWriter(os.Stdout)
	w := csv.NewWriter(file)

	if err := w.Write(fi.GenerateCSVHeaderRow()); err != nil {
		// at this point we're still trying to write to a non-flushed buffer,
		// so any failures are highly unexpected
		// TODO: Wrap error
		return err
	}

	//for key, fileMatches := range fi {
	for _, fileMatches := range fi {

		// TODO: Make this an external setting, either constant or via flag?
		if err := w.Write(fileMatches.GenerateEmptyCSVDataRow()); err != nil {
			// TODO: Use error wrapping instead?
			return fmt.Errorf("error writing record to csv: %v", err)
		}

		for _, file := range fileMatches {
			if err := w.Write(file.GenerateCSVDataRow()); err != nil {
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
		fmt.Fprintln(w,
			"Directory\tFile\tSize\tSize in bytes\tChecksum\t")

		for _, file := range fileMatches {

			// TODO: Confirm that newline between file sets is useful
			fmt.Fprintf(w,
				"\n%s\t%s\t%s\t%s\t%s",
				file.ParentDirectory,
				file.Name(),
				file.SizeHR(),
				strconv.FormatInt(file.Size(), 10),
				file.Checksum)
			fmt.Fprintln(w)
			w.Flush()
		}

	}

}
