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
	"path/filepath"

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/dupesets"
	"github.com/atc0005/bridge/paths"
)

// pruneSubcommand is a wrapper around the "prune" subcommand logic
func pruneSubcommand(appConfig *config.Config) {

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

	var dfsEntries dupesets.DuplicateFileSetEntries
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

		dfsEntry, err := dupesets.ParseInputRow(record, config.InputCSVFieldCount, rowCounter)
		if err != nil {
			log.Println("Error encountered parsing CSV file:", err)
			if appConfig.IgnoreErrors {
				log.Printf("IgnoringErrors set, ignoring input row %d.\n", rowCounter)
				continue
			}
			log.Fatal("IgnoringErrors NOT set. Exiting.")
		}

		// validate input row before we consider it OK
		if err = dupesets.ValidateInputRow(dfsEntry, rowCounter); err != nil {
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
}
