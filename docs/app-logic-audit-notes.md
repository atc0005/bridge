# Process Notes

## Summary

Scratch notes from an "audit" of the application logic while attempting to
refactor the codebase. This was done in preparation for splitting the
application logic between "report" and "prune" modes.

## Steps

1. Build configuration
   1. Process flags
   1. Validate
1. Build combined FileSizeIndex using file size as map key
   1. Loop over each provided path
   1. Merge the resulting FileSizeIndex created from each path
   - NOTE: We intentionally do *not* calculate the checksum/hash at this point
     as we would waste time hashing files that do not have a match based on
     file size
1. Prune FileMatches entries from combined FileSizeIndex if below
   duplicates threshold
1. Update checksums for remaining entries in combined FileSizeIndex
1. Build a new FileChecksumIndex using checksums as map key
1. Prune FileMatches entries from FileChecksumIndex if below duplicates
   threshold
1. Create DuplicateFilesSummary based off of FileSizeIndex and
   FileChecksumIndex
1. Optional: Generate Tabwriter output/report for FileChecksumIndex
1. Generate summary to stdout for DuplicateFilesSummary
1. Write CSV file
1. Optional: Write Excel file

## Thoughts

- Split out each "chunk" of work into separate functions?
