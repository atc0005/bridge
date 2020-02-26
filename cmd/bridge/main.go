// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
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
		// TODO: Replace this with Go 1.13 error equality check once 1.12 goes
		// EOL *and* we update CI to no longer use Go 1.12
		// if errors.Is(err, config.ErrMissingSubcommand) || errors.Is(err, flag.ErrHelp) {
		if (err == config.ErrMissingSubcommand) || (err == flag.ErrHelp) {
			os.Exit(0)
		}
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

	// We should not be able to reach this section
	default:
		log.Fatalf("invalid subcommand: %s", os.Args[1])
	}

}
