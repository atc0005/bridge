// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/atc0005/bridge/config"
)

func main() {

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		// if err is config.ErrMissingSubcommand or flag.ErrHelp we can skip
		// emitting err since the Help output shown by Parse() should be
		// sufficient enough
		// if errors.Is(err, config.ErrMissingSubcommand) || errors.Is(err, flag.ErrHelp) {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Printf("\nERROR: %s\n", err)
		os.Exit(1)
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

	// We should not be able to reach this section
	default:
		log.Fatalf("invalid subcommand: %s", os.Args[1])
	}

}
