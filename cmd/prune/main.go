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

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/paths"
)

func main() {

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		panic(err)
	}

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

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
		fmt.Println(record[0])
		fmt.Println(record[1])
		fmt.Println(record[2])
		fmt.Println(record[3])
		fmt.Println(record[4])
		fmt.Println(record[5])
	}

}
