package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/deathbots/epochtool"
	"strings"
	"testing"
)

var goodParse = []string{"3902432", "4928432432"}
var badParse = []string{"foo", "bar"}

func TestClipboardParseGood(t *testing.T) {
	if clipboard.Unsupported {
		return
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
		t.Errorf("Should have parsed %d numbers but got %d parsed", len(goodParse), len(epochResults))
	}
}

func TestClipboardParseBad(t *testing.T) {
	if clipboard.Unsupported {
		return
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
	fmt.Println(badStrings)
	if len(badStrings) != len(badParse) {
		t.Errorf("Should have failed to parse %d numbers but failed to parse %d", len(badParse), len(badStrings))
	}
	if len(epochResults) != 0 {
		t.Error("There should be no epoch results for ParseBad tests")
	}
}
