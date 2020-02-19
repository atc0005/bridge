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
		return fmt.Errorf("error encountered while removing %q: %s", filename, err)
	}

	log.Printf("Successfully removed %q\n", filename)
	return nil
}

// CreateBackupDirectoryTree receives a full path to a file that should be
// backed up and the full path to a base/target directory where the file
// should be copied. This function recreates the original directory structure
// as a means of avoiding name collisions when potentially backing up files
// with the same name.
func CreateBackupDirectoryTree(fullPathToFile string, fullPathToBackupDir string) error {

	// verify source file exists
	// verify target "base" directory exists
	// create source directory replicate underneath "base" directory

	// TODO:
	//
	// Need to confirm that fullPathToFile is a file
	// Need to get parent directory for fullPathToFile
	// If on Windows, need to strip off "C:\" as mkdir with that pattern present
	//  will likely not work well

	if !PathExists(fullPathToFile) {
		return fmt.Errorf("file %q does not exist or is inaccessible", fullPathToFile)
	}

	fileInfo, err := os.Stat(fullPathToFile)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf(
			"provided path to file %q is a directory, not a fully-qualified filename",
			fullPathToFile,
		)
	}

	if !PathExists(fullPathToBackupDir) {
		return fmt.Errorf("directory %q does not exist or is inaccessible", fullPathToBackupDir)
	}

	fileInfo, err := os.Stat(fullPathToBackupDir)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf(
			"provided path %q is a file, not an existing directory to use for backups",
			fullPathToBackupDir,
		)
	}

	// origString := `C:\Users\testing\Desktop\pictures\cat.jpg`

	// fmt.Println(filepath.Dir(origString))
	// fmt.Println(filepath.Base(origString))
	// fmt.Println(filepath.VolumeName(origString))

	// // replace volume name in directory path string
	// cleanString := filepath.ToSlash(origString)
	// fmt.Println(cleanString)

	// volRemoved := strings.ReplaceAll(cleanString, filepath.VolumeName(cleanString), "")
	// fmt.Println(volRemoved)

	// leadingSlashRemoved := strings.TrimPrefix(volRemoved, "/")
	// fmt.Println(leadingSlashRemoved)

	// replace backslashes with slashes
	//filepath.ToSlash()

	// if err := os.MkdirAll(); err != nil {

	// }

	return nil

}
