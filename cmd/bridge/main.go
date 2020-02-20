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
	"path/filepath"

	"log"
	"os"

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/matches"
	"github.com/atc0005/bridge/paths"
)

func main() {
	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		panic(err)
	}

	// DEBUG
	log.Printf("Configuration: %+v\n", appConfig)

	// behavior/logic switch between "prune" and "report" here
	switch os.Args[1] {
	case config.PruneSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.PruneSubcommand)

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

			dfsEntry, err := parseInputRow(record, config.InputCSVFieldCount, rowCounter)
			if err != nil {
				log.Println("Error encountered parsing CSV file:", err)
				if appConfig.IgnoreErrors {
					log.Printf("IgnoringErrors set, ignoring input row %d.\n", rowCounter)
					continue
				}
				log.Fatal("IgnoringErrors NOT set. Exiting.")
			}

			// validate input row before we consider it OK
			if err = validateInputRow(dfsEntry, rowCounter); err != nil {
				log.Println("Error encountered validating CSV row values:", err)
				if appConfig.IgnoreErrors {
					log.Printf("IgnoringErrors set, ignoring input row %d.\n", rowCounter)
					continue
				}
				log.Fatal("IgnoringErrors NOT set. Exiting.")
			}

			// update size details if found missing in CSV row
			if err = dfsEntry.UpdateSizeInfo(); err != nil {
				log.Println("Error encountered while attempting to update file size info:", err)
				if appConfig.IgnoreErrors {
					log.Printf("IgnoringErrors set, ignoring input row %d.\n", rowCounter)
					continue
				}
				log.Fatal("IgnoringErrors NOT set. Exiting.")
			}

			// Start off with collecting all entries in the CSV file that contain
			// all required fields. We'll filter the entries later to just those
			// that have been flagged for removal.
			dfsEntries = append(dfsEntries, dfsEntry)

		}

		// at this point we have parsed the CSV file into dfsEntries, validated
		// their content, regenerated file size details (if applicable) and are
		// now ready to begin work to remove flagged files.

		// DEBUG
		// fmt.Println("Length of dfsEntries:", len(dfsEntries))

		// Print parsed CSV file to the console if user requested it
		// NOTE: This contains ALL CSV file entries, not just those flagged for
		// removal.
		if appConfig.ConsoleReport {
			dfsEntries.Print(appConfig.BlankLineBetweenSets)
		}

		// if there are no files flagged for removal, say so and exit.
		filesToRemove := dfsEntries.FilesToRemove()
		if len(filesToRemove) == 0 {
			fmt.Printf("0 entries out of %d marked for removal in the %q input CSV file.\n",
				len(dfsEntries), appConfig.InputCSVFile)
			fmt.Println("Nothing to do, exiting.")
			return
		}

		// INFO? DEBUG?
		log.Printf("Found %d files to remove in %q", len(filesToRemove), appConfig.InputCSVFile)

		// DEBUG
		filesToRemove.Print(appConfig.BlankLineBetweenSets)

		// Skip backup logic and file removal if running in "dry-run" mode
		if !appConfig.DryRun {

			// DEBUG? INFO?
			fmt.Println("Dry-run not enabled, file removal mode enabled")

			if appConfig.BackupDirectory != "" {
				// DEBUG
				log.Println("Backup directory specified")

				// FIXME: The Config.Validate() method is also performing path checks
				// which is probably outside the normal scope for a config validation
				// function to perform. Because of that, we don't actually make it
				// to this point when a user provides an invalid backup directory path.
				if !paths.PathExists(appConfig.BackupDirectory) {
					// directory doesn't exist, what about the parent directory? do we
					// have permission to create content within the parent directory
					// to create the requested directory?

					// perhaps we should abort if the target directory doesn't exist?
					//
					// For example, we could end up trying to create a directory like
					// /tmp if the app is run as root. Since /tmp requires special
					// permissions, creating it as this application could lead to a
					// lot of problems that we cannot reliably anticipate and prevent

					log.Fatalf(
						"backup directory %q specified, but does not exist",
						appConfig.BackupDirectory,
					)
				}

				// attempt to backup files that the user marked for removal
				for _, file := range filesToRemove {

					fullPathToFile := filepath.Join(file.ParentDirectory, file.Filename)

					// attempt to backup files if user requested that we do so. if backup
					// failure occurs, abort. If file already exists in specified backup
					// directory check to see if they're identical. Report identical status
					// (yeah, nay) and abort unless an override or force option is given
					// (potential future work).

					// DEBUG
					// fmt.Printf("Calling BackupFile(%s, %s)\n", fullPathToFile, appConfig.BackupDirectory)

					err := paths.BackupFile(fullPathToFile, appConfig.BackupDirectory)
					if err != nil {
						// FIXME: Implement check for appconfig.IgnoreErrors
						// extend error message (potentially) to note that the error
						// was encountered when creating a backup
						log.Fatal(err)
					}

				}

			} else {
				// DEBUG
				log.Println("backup directory not set, not backing up files")
			}

			// Once backups complete remove original files. Allow IgnoreErrors setting
			// to apply, but be very noisy about removal failures

			var filesRemovedSuccess int
			var filesRemovedFail int
			for _, dfsEntry := range filesToRemove {

				fullPathToFile := filepath.Join(dfsEntry.ParentDirectory, dfsEntry.Filename)

				err = paths.RemoveFile(fullPathToFile, appConfig.DryRun)
				if err != nil {
					log.Printf("Error encountered while attempting to remove %q: %s\n",
						dfsEntry.Filename, err)
					if appConfig.IgnoreErrors {
						log.Println("IgnoringErrors set, ignoring failed file removal")
						filesRemovedFail++
						continue
					}
					log.Fatal("IgnoringErrors NOT set. Exiting.")
				}

				// note that we have successfully removed a file
				filesRemovedSuccess++

			}

			// print removal results summary
			fmt.Printf("File removal: %d success, %d fail\n",
				filesRemovedSuccess, filesRemovedFail)

		}

		if appConfig.DryRun {
			fmt.Println("Dry-run enabled, no files removed")
		}

	case config.ReportSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.ReportSubcommand)

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

		fmt.Printf("\n\nNext steps:\n\n")
		fmt.Printf("* Open %q\n", appConfig.OutputCSVFile)
		fmt.Printf("* Fill in the %q field with \"true\" for any file that you wish to remove\n",
			matches.CSVRemoveFileColumnHeaderName)
		fmt.Printf("* Run \"%s %s -h\" for a quick list of applicable options\n",
			os.Args[0], config.PruneSubcommand)
		fmt.Println("* Read the README for examples, including optional \"backup first\" behavior.")

	default:
		log.Fatal(config.ErrInvalidSubcommand)
	}

}
