// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/atc0005/bridge/units"
)

func main() {

	config := Config{}

	flag.Var(&config.Paths, "path", "Path to process. This flag may be repeated for each additional path to evaluate.")
	flag.Int64Var(&config.FileSizeThreshold, "size", 1, "File size limit for evaluation. Files smaller than this will be skipped.")
	flag.IntVar(&config.FileDuplicatesThreshold, "duplicates", 2, "number of files of the same file size needed before duplicate validation logic is applied.")
	flag.BoolVar(&config.RecursiveSearch, "recurse", false, "Perform recursive search into subdirectories per provided path.")
	flag.BoolVar(&config.ConsoleReport, "console", false, "Dump CSV file equivalent to console.")
	flag.BoolVar(&config.IgnoreErrors, "ignore-errors", false, "Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.")
	flag.StringVar(&config.CSVFile, "csvfile", "", "The fully-qualified path to a CSV file that this application should generate.")
	flag.StringVar(&config.ExcelFile, "excelfile", "", "The fully-qualified path to an Excel file that this application should generate.")

	// parse flag definitions from the argument list
	flag.Parse()

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Configuration: %+v\n", config)

	// evaluate all paths building a combined index of all files based on size
	combinedFileSizeIndex := make(FileSizeIndex)
	for _, path := range config.Paths {
		if PathExists(path) {
			log.Println("Path exists:", path)

			fileSizeIndex, err := ProcessPath(config.RecursiveSearch, config.IgnoreErrors, config.FileSizeThreshold, path)
			if err != nil {
				log.Println("Error encountered:", err)
				if config.IgnoreErrors {
					log.Println("Ignoring error as requested")
					continue
				}
			}

			combinedFileSizeIndex = MergeFileSizeIndexes(combinedFileSizeIndex, fileSizeIndex)
		}
	}

	var duplicateFiles DuplicateFilesSummary

	duplicateFiles.TotalEvaluatedFiles = len(combinedFileSizeIndex)
	//log.Println("combinedFileSizeIndex before pruning:", duplicateFiles.TotalEvaluatedFiles)

	// Prune FileMatches entries from map if below our file duplicates threshold
	combinedFileSizeIndex.PruneFileSizeIndex(config.FileDuplicatesThreshold)

	// Potential duplicate files going off of file size only (inconclusive)
	duplicateFiles.FileSizeMatchSets = len(combinedFileSizeIndex)
	//log.Println("combinedFileSizeIndex after pruning:", duplicateFiles.PotentialDuplicates)

	duplicateFiles.FileSizeMatches = combinedFileSizeIndex.GetTotalFilesCount()

	//for key, fileMatches := range combinedFileSizeIndex {
	for _, fileMatches := range combinedFileSizeIndex {

		// every key is a file size
		// every value is a slice of files of that file size

		if err := fileMatches.UpdateChecksums(config.IgnoreErrors); err != nil {
			log.Println("Error encountered:", err)
			if config.IgnoreErrors {
				log.Println("Ignoring error as requested")
				continue
			}
			log.Println("Exiting; error encountered, option to ignore (minor) errors not provided.")
			os.Exit(1)

		}

	}

	// At this point checksums have been calculated. We can use those
	// checksums to build a FileChecksumIndex in order to map checksums to
	// specific FileMatches objects.
	fileChecksumIndex := make(FileChecksumIndex)
	for _, fileMatches := range combinedFileSizeIndex {
		for _, fileMatch := range fileMatches {
			fileChecksumIndex[fileMatch.Checksum] = append(
				fileChecksumIndex[fileMatch.Checksum],
				fileMatch)
		}
	}

	// Remove FileMatches objects not meeting our file duplicates threshold
	// value. Remaining FileMatches that meet our file duplicates value are
	// composed entirely of duplicate files (based on file hash).
	//log.Println("fileChecksumIndex before pruning:", len(fileChecksumIndex))
	fileChecksumIndex.PruneFileChecksumIndex(config.FileDuplicatesThreshold)
	duplicateFiles.FileHashMatchSets = len(fileChecksumIndex)
	//log.Println("fileChecksumIndex after pruning:", len(fileChecksumIndex))

	duplicateFiles.FileHashMatches = fileChecksumIndex.GetTotalFilesCount()

	// TODO: Clean up redundant variables
	wastedSpace, err := fileChecksumIndex.GetWastedSpace()
	duplicateFiles.WastedSpace = wastedSpace
	if err != nil {
		// TODO: This shouldn't occur; worth testing?
		log.Fatal(err)
	}

	// Use text/tabwriter to dump results of the calculations directly to the
	// console. This is primarily intended for troubleshooting purposes.
	if config.ConsoleReport {
		fileChecksumIndex.PrintFileMatches()
	}

	// TODO: Use tabwriter to generate summary report?
	log.Printf("%d evaluated files in specified paths", duplicateFiles.TotalEvaluatedFiles)
	log.Printf("%d potential duplicate file sets found using file size", duplicateFiles.FileSizeMatchSets)
	log.Printf("%d confirmed duplicate file sets found using file hash", duplicateFiles.FileHashMatchSets)
	log.Printf("%d files with identical file size", duplicateFiles.FileSizeMatches)
	log.Printf("%d files with identical file hash", duplicateFiles.FileHashMatches)
	log.Printf("%d duplicate files", fileChecksumIndex.GetDuplicateFilesCount())
	log.Printf("%s wasted space for duplicate file sets", units.ByteCountIEC(duplicateFiles.WastedSpace))

	// Use CSV writer to generate an input file in order to take action
	// TODO: Implement better error handling
	if err := fileChecksumIndex.WriteFileMatchesCSV(config.CSVFile); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully created CSV file: %q", config.CSVFile)

	// Generate Excel workbook for review
	// TODO: Implement better error handling
	if err := fileChecksumIndex.WriteFileMatchesWorkbook(config.ExcelFile, duplicateFiles); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully created workbook file: %q", config.ExcelFile)

}
