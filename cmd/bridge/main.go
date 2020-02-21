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
)

func main() {

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		log.Fatal(err)
	}

	// DEBUG
	log.Printf("Configuration: %+v\n", appConfig)

	// behavior/logic switch between "prune" and "report" here
	switch os.Args[1] {
	case config.PruneSubcommand:

		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.PruneSubcommand)

		pruneSubcommand(appConfig)

	case config.ReportSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.ReportSubcommand)

		reportSubcommand(appConfig)

	default:
		log.Fatal(config.ErrInvalidSubcommand)
	}

}
