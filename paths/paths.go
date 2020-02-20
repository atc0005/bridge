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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atc0005/bridge/units"
)

const defaultDirectoryPerms os.FileMode = 0700

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

// GetBackupTargetDir returns a fully-qualified directory path based off of  a
// source file name and a destination "base" directory which will be used to
// hold a nested directory structure which attempts to approximate the
// current source file location.
func GetBackupTargetDir(filename string, fullPathToBackupDir string) (string, error) {

	fullPathToFile, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("unable to determine absolute path to %q: %s",
			filename,
			err,
		)
	}

	if !PathExists(fullPathToFile) {
		return "", fmt.Errorf("file %q does not exist or is inaccessible", fullPathToFile)
	}

	fileInfo, err := os.Stat(fullPathToFile)
	if err != nil {
		return "", err
	}

	if fileInfo.IsDir() {
		return "", fmt.Errorf(
			"provided path to file %q is a directory, not a fully-qualified filename",
			fullPathToFile,
		)
	}

	if !PathExists(fullPathToBackupDir) {
		return "", fmt.Errorf("directory %q does not exist or is inaccessible", fullPathToBackupDir)
	}

	dirInfo, err := os.Stat(fullPathToBackupDir)
	if err != nil {
		return "", err
	}

	if !dirInfo.IsDir() {
		return "", fmt.Errorf(
			"provided path %q is a file, not an existing directory to use for backups",
			fullPathToBackupDir,
		)
	}

	// at this point the source file has been confirmed and the target backup
	// directory has also been confirmed. We have yet to try actually reading
	// or writing from the provided locations, so access has (AFAIK) not been
	// confirmed.

	// create a sanitized version of the source filename path
	slashConvertedSourcePath := filepath.ToSlash(fullPathToFile)
	volumeName := filepath.VolumeName(slashConvertedSourcePath)
	volRemoved := strings.ReplaceAll(slashConvertedSourcePath, volumeName+"/", "")

	// fully-qualified path to place file backup which approximates
	// the original fully-qualified path
	targetBackupDirPath := filepath.Join(fullPathToBackupDir, volRemoved)

	return targetBackupDirPath, err

}

// CreateBackupDirectoryTree receives a full path to a file that should be
// backed up and the full path to a base/target directory where the file
// should be copied. This function recreates the original directory structure
// as a means of avoiding name collisions when potentially backing up files
// with the same name. If successful, this new path is returned, otherwise an
// empty string and the error is returned.
func CreateBackupDirectoryTree(filename string, fullPathToBackupDir string) (string, error) {

	targetBackupDirPath, err := GetBackupTargetDir(filename, fullPathToBackupDir)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(targetBackupDirPath, defaultDirectoryPerms); err != nil {
		return "", fmt.Errorf("failed to create fully-qualified backup path %q for %q: %s",
			targetBackupDirPath,
			filename,
			err,
		)
	}

	return targetBackupDirPath, err

}

// BackupFile accepts a path to a file and a destination directory where the
// file should be placed. The destination directory structure serves as a base
// directory for a nested structure that approximates the source file
// directory structure, omitting any OS-specific volume names (e.g., "C:\" on
// Windows).
func BackupFile(sourceFilename string, destinationDirectory string) error {

	targetBackupDirPath, err := CreateBackupDirectoryTree(sourceFilename, destinationDirectory)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory %q in order to backup %q: %s",
			sourceFilename,
			targetBackupDirPath,
			err,
		)
	}

	baseFileName := filepath.Base(sourceFilename)
	destinationFile := filepath.Join(targetBackupDirPath, baseFileName)

	// open source file for reading
	// open destination file for reading
	destinationFileHandle, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("unable to create new backup file %q: %s",
			destinationFile, err)
	}

	// TODO: Add wrapper for catching potential errors
	defer destinationFileHandle.Close()

	sourceFileHandle, err := os.Open(sourceFilename)
	if err != nil {
		return fmt.Errorf("unable to open source file %q in order to create backup copy: %s",
			sourceFilename, err)
	}
	defer sourceFileHandle.Close()

	sizeCopied, err := io.Copy(sourceFileHandle, destinationFileHandle)
	if err != nil {
		// copy failed, we should cleanup here
		log.Printf("failed to copy %q to %q: %s\n", sourceFilename, destinationFile, err)
		return fmt.Errorf("failed to copy %q to %q: %s", sourceFilename, destinationFile, err)
	}

	// DEBUG
	log.Printf("File %q successfully copied to %q (%s)",
		sourceFilename,
		destinationFile,
		units.ByteCountIEC(sizeCopied),
	)

	// copy was successful, we should still cleanup here, but should also log
	// (DEBUG) how much data was written in case we need that for
	// troubleshooting later

	// TODO:
	//
	// Sync file contents
	// Safely close destination file, catching as many errors as possible
	// Safely close source file
	return destinationFileHandle.Sync()

}
