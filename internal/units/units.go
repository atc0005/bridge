// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge

// Package units provides helper functions to perform conversions between
// various units of measurement.
package units

import "fmt"

// ByteCountSI converts a size in bytes to a human-readable string in SI
// (decimal) format.
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
// https://creativecommons.org/licenses/by/3.0/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountIEC converts a size in bytes to a human-readable string in IEC
// (binary) format.
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
// https://creativecommons.org/licenses/by/3.0/
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
