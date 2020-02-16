// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"log"

	"github.com/atc0005/bridge/config"
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

}
