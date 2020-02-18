// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package paths provides various functions and types related to processing
// paths in the filesystem.
package paths

import (
	"log"
	"os"
	"strings"

	"github.com/atc0005/bridge/config"
	"github.com/atc0005/bridge/matches"
	"github.com/sirupsen/logrus"
)

// PathExists confirms that the specified path exists
func PathExists(path string) bool {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		// DEBUG?
		// WARN?
		// ERROR?
		log.Println("path is empty string")
		return false
	}

	// https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//log.Println("path found")
		return true
	}

	return false

}

// CleanPath receives a slice of FileMatch objects and removes each file. Any
// errors encountered while removing files may optionally be ignored via
// command-line flag(default is to return immediately upon first error). The
// total number of files successfully removed is returned along with an error
// code (nil if no errors were encountered).
func CleanPath(files matches.FileMatches, config *config.Config) (PathPruningResults, error) {

	log := config.GetLogger()

	for _, file := range files {
		log.WithFields(logrus.Fields{
			"fullpath":        strings.TrimSpace(file.Path),
			"shortpath":       file.Name(),
			"size":            file.Size(),
			"modified":        file.ModTime().Format("2006-01-02 15:04:05"),
			"removal_enabled": config.GetRemove(),
		}).Debug("Matching file")
	}

	var removalResults PathPruningResults

	if !config.GetRemove() {

		log.Info("File removal not enabled, not removing files")

		// Nothing to show for this yet, but since the initial state reflects
		// that we can return it as-is
		return removalResults, nil
	}

	for _, file := range files {

		log.WithFields(logrus.Fields{
			"removal_enabled": config.GetRemove(),

			// fully-qualified path to the file
			"file": file.Path,
		}).Debug("Removing file")

		// We need to reference the full path here, not the short name since
		// the current working directory may not be the same directory
		// where the file is located
		err := os.Remove(file.Path)
		if err != nil {
			log.WithFields(logrus.Fields{

				// Include full details for troubleshooting purposes
				"file": file,
			}).Errorf("Error encountered while removing file: %s", err)

			// Record failed removal, proceed to the next file
			removalResults.FailedRemovals = append(removalResults.FailedRemovals, file)

			// Confirm that we should ignore errors (likely enabled)
			if !config.GetIgnoreErrors() {
				remainingFiles := len(files) - len(removalResults.FailedRemovals) - len(removalResults.SuccessfulRemovals)
				log.Debugf("Abandoning removal of %d remaining files", remainingFiles)
				break
			}

			log.Debug("Ignoring error as requested")
			continue
		}

		// Record successful removal
		removalResults.SuccessfulRemovals = append(removalResults.SuccessfulRemovals, file)
	}

	return removalResults, nil

}
