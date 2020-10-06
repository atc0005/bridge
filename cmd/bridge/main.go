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

	// Emulate returning exit code from main function by "queuing up" a
	// default exit code that matches expectations, but allow explicitly
	// setting the exit code in such a way that is compatible with using
	// deferred function calls throughout the application.
	var appExitCode *int
	defer func(code *int) {
		var exitCode int
		if code != nil {
			exitCode = *code
		}
		os.Exit(exitCode)
	}(appExitCode)

	var appConfig *config.Config
	var err error

	if appConfig, err = config.NewConfig(); err != nil {
		// if err is config.ErrMissingSubcommand or flag.ErrHelp we can skip
		// emitting err since the Help output shown by Parse() should be
		// sufficient enough
		// if errors.Is(err, config.ErrMissingSubcommand) || errors.Is(err, flag.ErrHelp) {
		if errors.Is(err, flag.ErrHelp) {
			*appExitCode = 0
			return
		}
		fmt.Printf("\nERROR: %s\n", err)
		*appExitCode = 1
		return
	}

	// DEBUG
	log.Printf("Configuration: %+v\n", appConfig)

	// behavior/logic switch between "prune" and "report" here
	switch os.Args[1] {
	case config.PruneSubcommand:

		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.PruneSubcommand)

		if err := pruneSubcommand(appConfig); err != nil {
			*appExitCode = 1
			fmt.Println(err)
			return
		}

	case config.ReportSubcommand:
		// DEBUG
		fmt.Printf("subcommand '%s' called\n", config.ReportSubcommand)

		if err := reportSubcommand(appConfig); err != nil {
			*appExitCode = 1
			fmt.Println(err)
			return
		}

	// We should not be able to reach this section
	default:
		log.Printf("invalid subcommand: %s", os.Args[1])
		*appExitCode = 1
		return
	}

}
