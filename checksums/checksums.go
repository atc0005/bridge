// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package checksums provides various functions and types related to processing
// file checksums.
package checksums

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// SHA256Checksum is a 64 character string representing a SHA256 hash
// TODO: How to assign a `string` to `[64]string` ?
// Goal: Set SHA256Checksum as the return type for GenerateCheckSum(), but
// make sure that the length is locked in at the specific character length
// for our chosen file hash.
//type SHA256Checksum [64]string
type SHA256Checksum string

func (cs SHA256Checksum) String() string {
	// convert the value via `string(cs)` before recurring to prevent infinite
	// recursion (per https://golang.org/pkg/fmt/ )
	return string(cs)
}

// Verify takes a path to a file, generates a SHA256 checksum from the file
// and compares against the checksum value already present
func (cs SHA256Checksum) Verify(file string) error {

	checksum, err := GenerateCheckSum(file)
	if err != nil {
		return err
	}

	if checksum.String() != cs.String() {
		return fmt.Errorf(
			"checksum mismatch, file likely modified; got %s, expected %s",
			checksum.String(),
			cs.String(),
		)
	}

	// if we got this far then the checksum is believed to be a match
	return nil

}

// GenerateCheckSum returns a SHA256 hash as the checksum generated from a
// provided fully-qualified path to a file.
func GenerateCheckSum(file string) (SHA256Checksum, error) {

	var checksum SHA256Checksum

	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		//log.Fatal(err)
		return checksum, err
	}

	// Note the duplicate f.Close() call at end of function and why
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf(
				"error occurred closing file %q: %v",
				file,
				err,
			)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		//log.Fatal(err)
		return checksum, err
	}

	// Explicitly convert Sprintf output from string to our type
	checksum = SHA256Checksum(fmt.Sprintf("%x", h.Sum(nil)))

	// defer the call to Close per above, and still report on an error if we
	// encounter one (see "Understanding defer in Go" README reference entry)
	return checksum, f.Close()

}
