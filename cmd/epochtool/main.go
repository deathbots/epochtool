// Copyright 2016 Rory Prendergast. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command line interface to epochconv. Accepts string arguments, but can also read from stdin and
// the clipboard.
// Makes a best guess as to the type of epoch being used based on its relation to the current time.

// The epochconv package is currently more flexible than the command line interface, allowing
// time ranges.

package main

import (
	"flag"
	"fmt"
	"github.com/deathbots/epochtool"
	"github.com/fatih/color"
	"io"
	"os"
	"strings"
)

//todo: flag to decide whether to print the time in UTC or in local timezone
// Use Semver always
var Version = "1.0.0-alpha.1"

type options struct {
	// The epoch date which is either given as a flag, or taken from clipboard
	printVersionFlag   bool
	epochsIn           []string
	useStdIn           bool
	useClipboard       bool
	noColor            bool
	emitJson           bool
	showAllConversions bool
}

// Some globals
var (
	opts             = new(options)
	progFriendlyName = "epochtool"
)

// Exit codes are stored here. Never remove one or change ordering, ever.
// Add new exit code const names here and they will automatically get a number.
// To look up which exit code corresponds to a number, count downward starting at
// exitBadFlags, which is -1. All subsequent are -2, -3, etc...
const (
	exitNoError  = iota // 0 value, all other values are -1 down
	exitBadFlags = -1 * iota
	exitNoEpochStringsError
	exitClipboardError
	exitNoNumbersParseableError
	exitStdinError
	exitJSONMarshallingError
)

const (
	// For testing if a required flag is not set. Make any required flags have this
	// default value.
	defFlagString = "REQUIRED"
)

func init() {
	flag.BoolVar(&opts.printVersionFlag, "version", false, "Print the version and quit")
	flag.BoolVar(&opts.useClipboard, "clipboard", false, "Parse data from the clipboard")
	flag.BoolVar(&opts.noColor, "no-color", false, "Disable color output")
	flag.BoolVar(&opts.emitJson, "json", true, "Print json output")
	flag.BoolVar(&opts.showAllConversions, "top-match", false, "Show only the closest match for each parsed epoch, "+
		"instead of the default case which is to show all matches.")
}

func main() {
	// todo: accept hex values as 0xXXXXXX and convert
	// todo: try to parse out using regex any part of the clipboard string.
	err := parseArgs()
	if err != nil {
		fatalPrint(exitBadFlags, "Unable to parse arguments", err)
	}
	// Add any items from stdin
	if opts.useStdIn {
		err = epochStringsFromStdin(&opts.epochsIn)
		if err != nil {
			fatalPrint(exitStdinError, "Unable to read data sent from stdin", err)
		}
	}
	// Add any items from os.args
	epochStringsFromCommandLine(&opts.epochsIn)

	if opts.useClipboard {
		err = epochStringsFromClipboard(&opts.epochsIn)
		if err != nil {
			fatalPrint(exitClipboardError, "Unable to read data from clipboard", err)
		}
	}
	if len(opts.epochsIn) == 0 {
		fatalPrint(exitNoEpochStringsError, "No data from command line, clipboard, or stdin", nil)
	}
	deDuplicateStringSlice(&opts.epochsIn)
	epochResults, badStrings, err := epochconv.GuessesForStrings(opts.epochsIn)
	if err != nil {
		stdErr("Could not parse the following input strings")
		for _, badString := range badStrings {
			stdErr(fmt.Sprintf("%s\n", badString))
		}
	}
	if len(epochResults) == 0 {
		fatalPrint(exitNoNumbersParseableError, "Found no numbers in input, cannot produce results\n", nil)
	}
	if opts.emitJson {
		era := EpochResultsArray{EpochResultsArray: epochResults}
		outJson, err := era.ToPrintableJson()
		if err != nil {
			fatalPrint(exitJSONMarshallingError, "Could not convert epoch results to JSON", err)
		}
		fmt.Println(outJson)
	} else {
		for _, er := range epochResults {
			// color output - Windows requires color.Output as the FPrint arg.
			fmt.Fprintf(color.Output, "%s\n", epochResultsAsString(er))
		}
		if err != nil {
			if opts.useClipboard {
				stdErr("Some strings could not be parsed, but they will remain hidden in clipboard mode.")
			}
			stdErr("Could not parse the following input strings:")
			for _, badString := range badStrings {
				stdErr(fmt.Sprintf("%s", badString))
			}
		}
	}
}

func parseArgs() (err error) {
	printVersion := func() {
		fmt.Printf("%s version %s\n", progFriendlyName, Version)
	}
	usage := func() {
		fmt.Printf("%s\nAccepts data to parse on command line, to stdin, or from the clipboard.", progFriendlyName)
		fmt.Println("Command line parsing:")
		fmt.Printf("\tUsage: %s -flags data1, data2 data3 \n", progFriendlyName)
		fmt.Println("Stdin parsing:")
		fmt.Printf("\tUsage: %s - < *.txt\n", progFriendlyName)
		fmt.Println("Clipboard parsing:")
		fmt.Printf("\tUsage: %s -clipboard\n", progFriendlyName)
		flag.PrintDefaults()
		fmt.Println("Unparseable strings are sent to stderr, except when -clipboard is specified.")
	}
	flag.Usage = usage
	defaultsChecker := func(a *flag.Flag) {
		if a.Value.String() == defFlagString {
			err = fmt.Errorf("A required flag -%s was not set", a.Name)
		}
	}
	opts.emitJson = false
	flag.Parse()
	if opts.noColor {
		color.NoColor = true // disables colorized output
	}
	if len(os.Args) == 1 {
		usage()
		os.Exit(exitNoError)
	}
	if opts.printVersionFlag {
		// Version specified in code, or may be set at build time with a linker flag to set the version based
		// on git tags.
		printVersion()
		os.Exit(exitNoError)
	}
	flag.VisitAll(defaultsChecker)
	if err != nil {
		usage()
		fmt.Printf("%s\n", err)
		return err
	}
	for _, arg := range os.Args[1:] {
		// Check if stdin is specified with -
		if arg == "-" {
			opts.useStdIn = true
		}
	}
	return err
}

// epochStringsFromCommandLine collects strings from the os.args, before any start with -,
// and adds to the collected strings list - passed by reference.
func epochStringsFromCommandLine(sliceToFill *[]string) {
	for _, arg := range flag.Args() {
		// Can we split this arg by comma? Could be 3098432,318401,4190
		// note - spaces after commas are OK, as are newlines, cr's, and tabs.
		// those are handled in the conversion to int64 function.

		//multiEpochsInSingleArg := strings.Split(arg, ",")
		//if len(multiEpochsInSingleArg) > 1 {
		//	// was split out
		//	fmt.Printf("Parsed out multiple strings in single comma separated")
		//	*sliceToFill = append(*sliceToFill, multiEpochsInSingleArg...)
		//} else {
		//	*sliceToFill = append(*sliceToFill, arg)
		//}
		*sliceToFill = append(*sliceToFill, arg)
	}
}

// epochStringsFromStdin takes a slice of strings and adds items from stdin using fmt.Scan
// which adds space-separated or newline separated values as successive items.
func epochStringsFromStdin(sliceToFill *[]string) (err error) {
	var s string
	for {
		_, err = fmt.Scan(&s)
		if err != nil {
			if err != io.EOF {
				return err
			}
			// before we break, clear err which would simply be io.EOF
			err = nil
			break
		}
		*sliceToFill = append(*sliceToFill, s)
	}
	return err
}

func epochStringsFromClipboard(sliceToFill *[]string) (err error) {
	s, err := getClipboardString()
	if err != nil {
		return err
	}
	*sliceToFill = append(*sliceToFill, strings.Split(s, "\n")...)
	*sliceToFill = append(*sliceToFill, strings.Split(s, "\r\n")...)
	*sliceToFill = append(*sliceToFill, strings.Split(s, ",")...)
	return err

}
