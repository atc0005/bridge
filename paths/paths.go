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
	"fmt"
	"log"
	"os"
	"strings"
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

// RemoveFile accepts a filename, a boolean flag indicating whether we are
// actually removing files or performing a dry-run and returns any errors that
// are encountered.
func RemoveFile(filename string, dryRun bool) error {

	if !dryRun {
		log.Printf("File removal not enabled, not removing %q\n", filename)
		return nil
	}

	log.Printf("File removal enabled, attempting to remove %q\n", filename)
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error encountered while removing file: %s", err)
	}

	log.Printf("Successfully removed %q\n", filename)
	return nil
}
