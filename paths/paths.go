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

	if dryRun {
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
			"provided path %q is not an existing directory to use for backups",
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

	// Strip off filename
	filenameRemoved := filepath.Dir(volRemoved)

	// fully-qualified path to place file backup which approximates
	// the original fully-qualified path
	targetBackupDirPath := filepath.Join(fullPathToBackupDir, filenameRemoved)

	return targetBackupDirPath, err

}

// CreateBackupDirectoryTree receives a full path to a file that should be
// backed up and the full path to a base/target directory where the file
// should be copied. This function recreates the original directory structure
// as a means of avoiding name collisions when potentially backing up files
// with the same name. If successful, this new path is returned, otherwise an
// empty string and the error is returned.
func CreateBackupDirectoryTree(filename string, fullPathToBackupDir string) (string, error) {

	// DEBUG
	// fmt.Printf("Calling GetBackupTargetDir(%s, %s)\n", filename, fullPathToBackupDir)

	targetBackupDirPath, err := GetBackupTargetDir(filename, fullPathToBackupDir)
	if err != nil {
		return "", err
	}

	// DEBUG
	//fmt.Printf("Calling os.MkdirAll(%v, %v)\n", targetBackupDirPath, defaultDirectoryPerms)

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

	// DEBUG
	// fmt.Printf("Calling CreateBackupDirectoryTree(%s, %s)\n", sourceFilename, destinationDirectory)

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

	// DEBUG
	fmt.Printf("sourceFilename: %q, baseFilename: %q, destinationFile: %q\n",
		sourceFilename, baseFileName, destinationFile)

	// verify that destinationFile does not already exist before calling
	// os.Create(), otherwise we will end up truncating the existing file
	if PathExists(destinationFile) {
		return fmt.Errorf(
			"destination file %q already exists; skipping backup of %q to prevent overwriting existing file",
			destinationFile,
			sourceFilename,
		)
	}

	destinationFileHandle, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("unable to create new backup file %q: %s",
			destinationFile, err)
	}
	defer func() {
		if err := destinationFileHandle.Close(); err != nil {
			log.Printf(
				"error occurred closing file %q: %v",
				destinationFile,
				err,
			)
		}
	}()

	// guard against invalid source files
	// NOTE: This shouldn't be possible since we only add files to a list to
	// be removed once we verify checksum and collect file size details, but
	// to be complete (e.g, perhaps if/when this function is moved to a
	// standalone package for use by other applications one day) we go ahead
	// and verify again that the source file is a valid backup source
	sourceFileStat, err := os.Stat(sourceFilename)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%q is not a regular file", sourceFileStat)
	}

	sourceFileHandle, err := os.Open(filepath.Clean(sourceFilename))
	if err != nil {
		return fmt.Errorf("unable to open source file %q in order to create backup copy: %s",
			sourceFilename, err)
	}
	defer func() {
		if err := sourceFileHandle.Close(); err != nil {
			log.Printf(
				"error occurred closing file %q: %v",
				sourceFilename,
				err,
			)
		}
	}()

	sizeCopied, err := io.Copy(destinationFileHandle, sourceFileHandle)
	if err != nil {
		// copy failed, we should cleanup here
		log.Printf("failed to copy %q to %q: %s\n", sourceFilename, destinationFile, err)
		return fmt.Errorf("failed to copy %q to %q: %s", sourceFilename, destinationFile, err)
	}

	// I encountered this when I unintentionally switched the dest/source
	// values for io.Copy() (I keep thinking of copying source to destination,
	// not copy to destination from source)
	if sizeCopied != sourceFileStat.Size() {
		// no content copied failed, we should consider this a failure
		sizeCopiedMismatchMsg := fmt.Sprintf(
			"failed to copy %q to %q: %d of %d bytes copied\n",
			sourceFilename,
			destinationFile,
			sizeCopied,
			sourceFileStat.Size(),
		)
		log.Println(sizeCopiedMismatchMsg)
		return fmt.Errorf(sizeCopiedMismatchMsg)
	}

	// copy was successful, we should cleanup and log (DEBUG) how much data
	// was written (in case we need that for troubleshooting later)

	// DEBUG
	log.Printf("File %q successfully copied to %q (%s)",
		sourceFilename,
		destinationFile,
		units.ByteCountIEC(sizeCopied),
	)

	if err := destinationFileHandle.Sync(); err != nil {
		return fmt.Errorf(
			"failed to explicitly sync file %q after backup attempt: %s",
			destinationFile,
			err,
		)
	}

	if err := sourceFileHandle.Close(); err != nil {
		return fmt.Errorf(
			"failed to close original file %q after backup attempt: %s",
			sourceFilename,
			err,
		)
	}

	return destinationFileHandle.Close()

}
