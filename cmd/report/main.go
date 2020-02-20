// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/matches"
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
			log.Fatalf("Failed to build file size index from paths (%q): %v", appConfig.Paths.String(), err)
		}
		log.Println("Error encountered:", err)
		log.Println("Attempting to ignore errors as requested")
	}

	// TODO: Refactor this; merge into NewFileSizeIndex? NewFileChecksumIndex?
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
	fileChecksumIndex := matches.NewFileChecksumIndex(combinedFileSizeIndex)

	// Remove FileMatches objects not meeting our file duplicates threshold
	// value. Remaining FileMatches that meet our file duplicates value are
	// composed entirely of duplicate files (based on file hash).
	//log.Println("fileChecksumIndex before pruning:", len(fileChecksumIndex))
	fileChecksumIndex.PruneFileChecksumIndex(appConfig.FileDuplicatesThreshold)

	// Use text/tabwriter to dump results of the calculations directly to the
	// console. This is primarily intended for troubleshooting purposes.
	if appConfig.ConsoleReport {
		fileChecksumIndex.PrintFileMatches(appConfig.BlankLineBetweenSets)
	}

	// TODO: Move this into a separate package?
	// Note: FileSizeMatchSets represents *potential* duplicate files going
	// off of file size only (inconclusive)
	duplicateFiles := matches.DuplicateFilesSummary{
		TotalEvaluatedFiles: len(combinedFileSizeIndex),
		FileSizeMatches:     combinedFileSizeIndex.GetTotalFilesCount(),
		FileSizeMatchSets:   len(combinedFileSizeIndex),
		FileHashMatches:     fileChecksumIndex.GetTotalFilesCount(),
		FileHashMatchSets:   len(fileChecksumIndex),
		WastedSpace:         fileChecksumIndex.GetWastedSpace(),
		DuplicateCount:      fileChecksumIndex.GetDuplicateFilesCount(),
	}

	duplicateFiles.PrintSummary()

	// Use CSV writer to generate an input file in order to take action
	// TODO: Implement better error handling
	if err := fileChecksumIndex.WriteFileMatchesCSV(
		appConfig.OutputCSVFile, appConfig.BlankLineBetweenSets); err != nil {
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

	fmt.Printf("\n\nNext steps:\n")
	fmt.Printf("Open %q\n", appConfig.OutputCSVFile)
	fmt.Printf("Fill in the %q field with \"true\" for any file that you wish to remove\n", csvRemoveFileColumnHeaderName)
	fmt.Printf("Run \"%s %s\" with applicable options\n", os.Args[0], config.pruneSubcommand)
	fmt.Println("Read the README for specific options, including optional \"backup first\" behavior.")

}
