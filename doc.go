/*
bridge is a small CLI utility used to find duplicate files.

# Project Home

See our GitHub repo (https://github.com/atc0005/bridge) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

# Purpose

bridge is intended for locating duplicate files across one or many paths. The
inspiration was managing digital photos collected in date-based folders as
well as event-based folders. Often we would gather folders together in
"collections" for processing and unintentionally duplicate the existing images
already sorted by date.

# Features

  - single binary, no outside dependencies
  - minimal configuration
  - very few build dependencies
  - shallow or recursive processing across one or more paths
  - matches initially based on file size, confirm via file hash
  - generate CSV report of all duplicate file matches
  - (optionally) generate Microsoft Excel workbook of all duplicate file matches
  - generate (rough) console equivalent of CSV file for quick review
  - (optionally) remove user-flagged duplicate files using a previously generated CSV report

# Usage

See our main README for supported settings and examples.
*/
package main
