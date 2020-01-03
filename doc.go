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

• generate CSV report of all matches

USAGE

Help output is below. See the README for examples.

    Usage of T:\github\bridge\bridge.exe:

    -console
            Dump CSV file equivalent to console.
    -csvfile string
            The fully-qualified path to a CSV file that this application should generate.
    -ignore-errors
            Ignore minor errors whenever possible. This option does not affect handling of fatal errors such as failure to generate output report files.
    -path value
            Path to process. This flag may be repeated for each additional path to evaluate.
    -recurse
            Perform recursive search into subdirectories per provided path.

*/
package main
