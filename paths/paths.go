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
