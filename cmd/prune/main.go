// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"encoding/csv"

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

}
