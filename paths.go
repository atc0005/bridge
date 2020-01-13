// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// PathExists confirms that the specified path exists
func PathExists(path string) bool {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
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

// ProcessPath optionally recursively processes a provided path and returns a
// slice of FileMatch objects
func ProcessPath(recursiveSearch bool, ignoreErrors bool, path string) (FileSizeIndex, error) {

	fileSizeIndex := make(FileSizeIndex)
	var err error

	//log.Println("RecursiveSearch:", recursiveSearch)

	if recursiveSearch {

		// Walk walks the file tree rooted at path, calling the anonymous function
		// for each file or directory in the tree, including path. All errors that
		// arise visiting files and directories are filtered by the anonymous
		// function. The files are walked in lexical order, which makes the output
		// deterministic but means that for very large directories Walk can be
		// inefficient. Walk does not follow symbolic links.
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			// If an error is received, return it. If we return a non-nil error, this
			// will stop the filepath.Walk() function from continuing to walk the
			// path, and your main function will immediately move to the next line.
			if err != nil {
				if !ignoreErrors {
					return err
				}

				// WARN
				log.Println("Error encountered:", err)
				log.Println("Ignoring error as requested")

			}

			// make sure we're not working with the root directory itself
			if path != "." {

				// ignore directories
				if info.IsDir() {
					return nil
				}

				// ignore files below the size threshold
				if info.Size() < SizeThreshold {
					return nil
				}

				// If we made it to this point, then we must assume that the file
				// has met all criteria to be evaluated by this application.
				// Let's add the file to our slice of files of the same size
				// using our index based on file size.
				fileSizeIndex[info.Size()] = append(
					fileSizeIndex[info.Size()],
					FileMatch{
						FileInfo:        info,
						FullPath:        path,
						ParentDirectory: filepath.Dir(path),
					})
			}

			return err
		})

	} else {

		// If recursiveSearch is not enabled, process just the provided path

		// err is already declared earlier at a higher scope, so do not
		// redeclare here
		var files []os.FileInfo
		files, err = ioutil.ReadDir(path)

		if err != nil {
			// TODO: Wrap error?
			log.Printf("Error from ioutil.ReadDir(): %s", err)

			return fileSizeIndex, err
		}

		// Use []os.FileInfo returned from ioutil.ReadDir() to build slice of
		// FileMatch objects
		for _, file := range files {

			// ignore directories
			if file.IsDir() {
				continue
			}

			// ignore files below the size threshold
			if file.Size() < SizeThreshold {
				continue
			}

			// If we made it to this point, then we must assume that the file
			// has met all criteria to be evaluated by this application.
			// Let's add the file to our slice of files of the same size
			// using our index based on file size.
			fileSizeIndex[file.Size()] = append(
				fileSizeIndex[file.Size()],
				FileMatch{
					FileInfo: file,
					FullPath: filepath.Join(path, file.Name()),
					// ParentDirectory: filepath.Dir(path),
					// `path` is a flat directory structure (we are not using
					// recursion), so record it directly as the parent
					// directory for files within
					ParentDirectory: path,
				})
		}
	}

	return fileSizeIndex, err
}
