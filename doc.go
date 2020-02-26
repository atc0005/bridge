/*

bridge is a small CLI utility used to find duplicate files.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/bridge) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

PURPOSE

bridge is intended for locating duplicate files across one or many paths. The
inspiration was managing digital photos collected in date-based folders as
well as event-based folders. Often we would gather folders together in
"collections" for processing and unintentionally duplicate the existing images
already sorted by date.

FEATURES

• single binary, no outside dependencies

• minimal configuration

• very few build dependencies

• shallow or recursive processing across one or more paths

• matches initially based on file size, confirm via file hash

• generate CSV report of all duplicate file matches

• (optionally) generate Microsoft Excel workbook of all duplicate file matches

• generate (rough) console equivalent of CSV file for quick review

• (optionally) remove user-flagged duplicate files using a previously
  generated CSV report

USAGE

Help output is below. See the README for examples.

$ ./bridge -h

    bridge x.y.z
    https://github.com/atc0005/bridge

    Available subcommands:
            prune
            report

    See available options for each subcommand by running:
            bridge prune -h
            bridge report -h

    DEBUG: Help flags used

    ERROR: missing subcommand


$ ./bridge report -h

    DEBUG: subcommand 'report'

    bridge x.y.z
    https://github.com/atc0005/bridge

    Usage of "bridge report":
    -console
            Dump (approximate) CSV file equivalent to console.
    -csvfile string
            The (required) fully-qualified path to a CSV file that this application should generate.
    -duplicates int
            Number of files of the same file size needed before duplicate validation logic is applied. (default 2)
    -excelfile string
            The (optional) fully-qualified path to an Excel file that this application should generate.
    -ignore-errors
            Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.
    -path value
            Path to process. This flag may be repeated for each additional path to evaluate.
    -recurse
            Perform recursive search into subdirectories per provided path.
    -size int
            File size limit (in bytes) for evaluation. Files smaller than this will be skipped. (default 1)
    DEBUG: err returned from reportCmd.Parse(): flag: help requested

$ ./bridge prune -h

    DEBUG: subcommand 'prune'

    bridge x.y.z
    https://github.com/atc0005/bridge

    Usage of "bridge prune":
    -backup-dir string
            The writable directory path where files should be relocated instead of removing them. The original path structure will be created starting with the specified path as the root.
    -blank-line
            Add a blank line between sets of matching files in console and file output.
    -console
            Dump (approximate) CSV file equivalent to console.
    -dry-run
            Don't actually remove files. Echo what would have been done to stdout.
    -ignore-errors
            Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.
    -input-csvfile string
            The fully-qualified path to a CSV file that this application should use for file removal decisions.
    -use-first-row
            Attempt to use the first row of the input file. Normally this row is skipped since it is usually the header row and not duplicate file data.
    DEBUG: err returned from pruneCmd.Parse(): flag: help requested
*/
package main
