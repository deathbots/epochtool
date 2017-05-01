package main

import (
	"encoding/json"
	"github.com/atotto/clipboard"
	"github.com/deathbots/epochtool"
	"strings"
	"testing"
	"runtime"
	"fmt"
)

var goodParse = []string{"3902432", "4928432432"}
// this number won't fit in an int64
var badParse = []string{"3452543252352353253253252"}

func TestClipboardParseGood(t *testing.T) {
	if clipboard.Unsupported {
		t.Skipf("Clipboard not supported on OS %s", runtime.GOOS)
	}
	strs := make([]string, 0)
	clipboard.WriteAll(strings.Join(goodParse, ", \t"))
	err := epochStringsFromClipboard(&strs)
	if err != nil {
		t.Errorf(err.Error())
	}
	epochResults, badStrings, err := epochconv.GuessesForStrings(strs)
	if err != nil {
		t.Errorf("Failure to parse known good numbers - these failed:%v with error: %s", badStrings, err)
	}
	if len(epochResults) != len(goodParse) {
		t.Errorf("Should have parsed %d numbers but got %d parsed " +
			"- in %v, out %v", len(goodParse), len(epochResults), goodParse, epochResults)
	}
}

func TestClipboardParseBad(t *testing.T) {
	if clipboard.Unsupported {
		t.Skipf("Clipboard not supported on OS %s", runtime.GOOS)
	}
	strs := make([]string, 0)
	clipboard.WriteAll(strings.Join(badParse, ", \t"))
	err := epochStringsFromClipboard(&strs)
	if err != nil {
		t.Errorf(err.Error())
	}
	epochResults, badStrings, err := epochconv.GuessesForStrings(strs)
	if err == nil {
		t.Errorf("Should have received error parsing bad string: %s", err)
	}
	if len(badStrings) != len(badParse) {
		t.Errorf("Should have failed to parse %d numbers but failed to parse %d", len(badParse), len(badStrings))
	}
	if len(epochResults) != 0 {
		t.Error("There should be no epoch results for ParseBad tests")
	}
}


func TestCmdLineParseGood(t *testing.T) {
	strs := make([]string, 0)
	epochStringsFromCommandLine(&strs, goodParse)
	epochResults, badStrings, err := epochconv.GuessesForStrings(strs)
	if err != nil {
		t.Errorf("Failure to parse known good numbers - these failed:%v with error: %s", badStrings, err)
	}
	if len(epochResults) != len(goodParse) {
		t.Errorf("Should have parsed %d numbers but got %d parsed " +
			"- in %v, out %v", len(goodParse), len(epochResults), goodParse, epochResults)
	}
}

func TestCmdLineParseBad(t *testing.T) {
	strs := make([]string, 0)
	epochStringsFromCommandLine(&strs, badParse)
	epochResults, badStrings, err := epochconv.GuessesForStrings(strs)
	if err == nil {
		t.Errorf("Should have received error parsing bad string: %s", err)
	}
	if len(badStrings) != len(badParse) {
		t.Errorf("Should have failed to parse %d numbers but failed to parse %d", len(badParse), len(badStrings))
	}
	if len(epochResults) != 0 {
		t.Error("There should be no epoch results for ParseBad tests")
	}
}

// this number is so high, it will always produce common era, the oldest epoch
var commonEraHighInt = []string{"9223346836"}

func TestJsonParse(t *testing.T) {
	epochResults, badStrings, err := epochconv.GuessesForStrings(commonEraHighInt)
	if err != nil {
		stdErr("Could not parse the following input strings")
		for _, badString := range badStrings {
			stdErr(fmt.Sprintf("%s\n", badString))
		}
	}
	era := EpochResultsArray{EpochResultsArray: epochResults}
	outJson, err := era.ToPrintableJson()
	err = json.Unmarshal([]byte(outJson), &era)
	if err != nil {
		t.Errorf("Could not convert json back into data structure, %s", err)
	}
	if era.EpochResultsArray[0].MostLikelyType.EpochName != "CommonEra" {
		t.Error("JSON result was incorrect")
	}
}

