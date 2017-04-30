package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/deathbots/epochtool"
	"github.com/fatih/color"
	"os"
	"time"
)

// This type is used only for JSON marshalling.
type EpochResultsArray struct {
	EpochResultsArray []epochconv.EpochResults `json:"epoch_results_array"`
}

var (
	colorMostLikely = color.New(color.FgHiGreen).SprintFunc()
)

// fatalPrint is a convenience function that will quit the program with the specified Exit Code, print some friendly
// context and a colon, and the error message from golang. Pass a nil error in to avoid printing the error string.
func fatalPrint(exitCode int, friendlyContext string, err error) {
	if err == nil {
		stdErr(fmt.Sprintf("Error: %s", friendlyContext))
	} else {
		stdErr(fmt.Sprintf("Error: %s: %s", friendlyContext, err))
	}
	os.Exit(exitCode)
}

// stdErr is a convenience function to send output to stderr. Automatically adds newline.
func stdErr(m string) {
	fmt.Fprintf(os.Stderr, "%s\n", m)
}

func getClipboardString() (out string, err error) {
	if clipboard.Unsupported {
		err = errors.New("Clipboard functionality unsupported.")
		return out, err
	}
	return clipboard.ReadAll()
}

func deDuplicateStringSlice(sliceToDeDupe *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *sliceToDeDupe {
		if !found[x] {
			found[x] = true
			(*sliceToDeDupe)[j] = (*sliceToDeDupe)[i]
			j++
		}
	}
	*sliceToDeDupe = (*sliceToDeDupe)[:j]
}

func epochResultsAsString(ers epochconv.EpochResults) string {
	var out string
	// for non-string types that are printable via %s, you must turn them to strings first
	// in order to apply a color.
	colorMe := fmt.Sprintf("%s", ers.MostLikelyType)
	out = out + fmt.Sprintf("For Input Number: %d\n"+
		"---------Most Likely Result----\n"+
		"%s"+
		"---------Other Results---------\n", ers.InputNumber, colorMostLikely(colorMe))
	c := color.New(color.Reset).SprintfFunc()
	// ol: //tag outer loop so we can break from within the switch when necessary
	for i, er := range ers.AllResults {
		switch {
		//case i == 0:
		//	continue // already printed most likely result
		case i == 1:
			c = color.New(color.FgHiYellow).SprintfFunc()
		case i == 2:
			c = color.New(color.FgYellow).SprintfFunc()
		default:
			c = color.New(color.Reset).SprintfFunc()
		}
		// If the prevalence is low, do not color it lightly.
		if er.EpochType.Prevalence < 3 {
			c = color.New(color.Faint).SprintfFunc()
		}
		m := fmt.Sprintf("%d in this Epoch:\n"+
			" Local - %s\n"+
			" UTC - %s\n"+
			"%s\n", ers.InputNumber, er.DateInEpochLocal.Format(time.RFC3339),
			er.DateInEpochUTC.Format(time.RFC3339), er.EpochType)
		out = out + fmt.Sprintf("%s", c(m))
	}

	return out
}

// easy conversion of this type made for JSON marshalling to json
func (era EpochResultsArray) ToPrintableJson() (string, error) {
	jsonByteArray, err := json.MarshalIndent(&era, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonByteArray), err
}
