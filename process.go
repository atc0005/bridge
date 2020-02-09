// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

// TODO: Move content here into the appropriate file

// open CSV file
// check first row
// validate header is present, all required fields
// read a row, validate required fields are present
//	if failure
//		note as much
//		check first column
//		if not empty
//			check flag for continuing past errors
//			consider this a failed parse attempt
//		if empty
//			note, consider this an empty row (perhaps CSV library already assists with handling empty rows?)
//			proceed to next row
// 	if required fields are present
//		confirm file exists
//			if failure
//				check flag for continuing past errors
//					record as failed removal (specific error for file not present)
//			if success
//				attempt to remove file
//					if success
//						great, note as much (somehow)
//					if failure
//						check flag for continuing past errors
