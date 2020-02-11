// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"log"
	"os"

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/matches"
	"github.com/atc0005/bridge/units"
)

func main() {

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		panic(err)
	}

	log.Printf("Configuration: %+v\n", appConfig)

	// evaluate all paths building a combined index of all files based on size
	combinedFileSizeIndex, err := matches.NewFileSizeIndex(
		appConfig.RecursiveSearch,
		appConfig.IgnoreErrors,
		appConfig.FileSizeThreshold,
		appConfig.Paths...,
	)

	if err != nil {
		if !appConfig.IgnoreErrors {
			log.Fatalf("Failed to process paths %q: %v", appConfig.Paths.String(), err)
		}
		log.Println("Error encountered:", err)
		log.Println("Attempting to ignore errors as requested")
	}

	// Prune FileMatches entries from map if below our file duplicates threshold
	combinedFileSizeIndex.PruneFileSizeIndex(appConfig.FileDuplicatesThreshold)

	if err := combinedFileSizeIndex.UpdateChecksums(appConfig.IgnoreErrors); err != nil {
		log.Println("Exiting; error encountered, option to ignore (minor) errors not provided.")
		os.Exit(1)
	}

	// TODO: Move this to matches package
	//
	// At this point checksums have been calculated. We can use those
	// checksums to build a FileChecksumIndex in order to map checksums to
	// specific FileMatches objects.
	fileChecksumIndex := make(matches.FileChecksumIndex)
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
	fileChecksumIndex.PruneFileChecksumIndex(appConfig.FileDuplicatesThreshold)

	// TODO: Move this into a separate package
	var duplicateFiles matches.DuplicateFilesSummary

	duplicateFiles.TotalEvaluatedFiles = len(combinedFileSizeIndex)
	//log.Println("combinedFileSizeIndex before pruning:", duplicateFiles.TotalEvaluatedFiles)

	// Potential duplicate files going off of file size only (inconclusive)
	duplicateFiles.FileSizeMatchSets = len(combinedFileSizeIndex)
	//log.Println("combinedFileSizeIndex after pruning:", duplicateFiles.PotentialDuplicates)

	duplicateFiles.FileSizeMatches = combinedFileSizeIndex.GetTotalFilesCount()
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
	if appConfig.ConsoleReport {
		fileChecksumIndex.PrintFileMatches()
	}

	// TODO: Move this into a separate package
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
	if err := fileChecksumIndex.WriteFileMatchesCSV(appConfig.OutputCSVFile); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully created CSV file: %q", appConfig.OutputCSVFile)

	// Generate Excel workbook for review IF user requested it
	if appConfig.ExcelFile != "" {
		// TODO: Implement better error handling
		if err := fileChecksumIndex.WriteFileMatchesWorkbook(appConfig.ExcelFile, duplicateFiles); err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully created workbook file: %q", appConfig.ExcelFile)
	}

}
