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
	"path/filepath"

	"io"
	"log"
	"os"

	"github.com/atc0005/bridge/config"
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

		dfsEntries = append(dfsEntries, dfsEntry)

	}

	// Print parsed CSV file to the console if user requested it
	if appConfig.ConsoleReport {
		dfsEntries.Print(appConfig.BlankLineBetweenSets)
	}

	// at this point we have parsed the CSV file into dfsEntries, validated
	// their content, regenerated file size details (if applicable) and are
	// now ready to begin work to remove flagged files.

	// if there are no files flagged for removal, say so and exit.
	filesToRemove := dfsEntries.FilesToRemove()
	if filesToRemove == 0 {
		fmt.Printf("0 entries out of %d marked for removal in the %q input CSV file.\n",
			len(dfsEntries), appConfig.InputCSVFile)
		fmt.Println("Nothing to do, exiting.")
		return
	}

	// INFO? DEBUG?
	log.Printf("Found %d files to remove in %q", filesToRemove, appConfig.InputCSVFile)

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

			// attempt to backup file
			// NOTE: at this point the only files that would be removed (or backed
			// up) are those that were flagged for removal in the CSV file
			for _, file := range dfsEntries {

				fullPathToFile := filepath.Join(file.ParentDirectory, file.Filename)
				// paths.CreateBackupDirectoryTree(fullPathToFile, appConfig.BackupDirectory)

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
		for _, dfsEntry := range dfsEntries {

			err = paths.RemoveFile(dfsEntry.Filename, appConfig.DryRun)
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
