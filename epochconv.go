// Copyright 2016 Rory Prendergast. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package epochconv will accept an integer epoch and make a best guess, by looking at the converted year, as
// to which type of epoch type it is.

package epochconv

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"regexp"
)

//todo: JS Epoch (unix in milliseconds) currently unhandled?
// perhaps for each number, divide and see if it's close to unix?

// epochResult is used in an EpochResultBundle
type epochResult struct {
	InputNumber      int64     `json:"input_number"`
	EpochType        EpochType `json:"epoch_type"`
	DateInEpochLocal time.Time `json:"converted_date_local"`
	DateInEpochUTC   time.Time `json:"converted_date_utc"`
}

type EpochResults struct {
	InputNumber    int64           `json:"input_number"`
	EpochTypes     EpochCollection `json:"epoch_types"`
	AllResults     []epochResult   `json:"all_results"`
	MostLikelyType EpochType       `json:"most_likely_epoch"`
}

// Given a slice of strings, return a slice of EpochGuessResults type, each of which is an array of EpochResults along
// with the most likely result. Strings in the input slice are parsed in the following way:
// 1) Strings are stripped of leading and trailing whitespace characters.
// 2) Strings have all data after the first dot character removed. This allows for input of decimal numbers
//    Without needing to convert to floats.
// If one string that seemed to match a number cannot be converted, an Error is returned.
// However, the numbers that were convertible are still returned. Ignore the error and continue, if desired.
func GuessesForStrings(stringsToConvert []string) (epochResults []EpochResults, badStrings []string, err error) {
	epochResults, badStrings, err = createGuesses(stringsToConvert, AllEpochs)
	return epochResults, badStrings, err
}

func createGuesses(stringsToConvert []string, collection EpochCollection) (epochResultsSlice []EpochResults,
	badStrings []string, err error) {

	numbers, badStrings, err := stringSliceToInt64Base10s(stringsToConvert)
	// Results array is as long as parsed numbers
	epochResultsSlice = make([]EpochResults, len(numbers))
	// loop through numbers and create epochs result data structures, which are an epoch type
	// and the date in that epoch.
	for i, n := range numbers {
		epochResultsSlice[i].EpochTypes = collection
		epochResultsSlice[i].InputNumber = n
		for _, et := range collection {
			er := epochResult{InputNumber: n,
				EpochType:        et,
				DateInEpochLocal: et.DateForNumber(n, false),
				DateInEpochUTC:   et.DateForNumber(n, true),
			}
			epochResultsSlice[i].AllResults = append(epochResultsSlice[i].AllResults, er)
		}
		// Run OrderedEpochsByClosestMatch on EC which takes a number and a time to match on.
		for i, _ := range epochResultsSlice {
			epochResultsSlice[i].EpochTypes = epochResultsSlice[i].EpochTypes.OrderedEpochsByClosestMatch(n, time.Now())
		}
		epochResultsSlice[i].MostLikelyType = epochResultsSlice[i].EpochTypes[0]
	}
	return epochResultsSlice, badStrings, err
}

// OrderedEpochsByClosestMatch is a Method on an EpochCollection. Given an EpochCollection, typically AllEpochs,
// return a collection order by closest match of an epoch number given a date to convert to all epoch seconds. Do not
// alter the collection slice order in-place but, instead, return the sorted EpochCollection.
// The first item in the returned Collection is the closest match, and matches are less likely at the end of the slice.
//
// Create your own EpochCollection by hand to add and remove existing or custom epochs.
func (ec EpochCollection) OrderedEpochsByClosestMatch(number int64, matchToTime time.Time) (ecOut EpochCollection) {
	// Do not sort the collection in place, return a new one.
	sorted := make(EpochCollection, len(ec))
	copy(sorted, ec)
	sort.Sort(ByEpochDate(sorted))
	// the return list is just as long as the original list
	ecOut = make(EpochCollection, len(ec))

	// make a slice just containing the epoch seconds for the input date
	// the indices will match the indices of sorted. This is a convenience, and
	// makes it simpler to accomplish ordering the list.
	epochsOnly := make([]int64, len(ec))
	for i, et := range sorted {
		epochsOnly[i] = et.NumberForDate(matchToTime)
	}
	// distance stores how close the number is. Is a 2d array because it will store the original position after
	// it is sorted, which can be examined to determine how to fill the final list.
	epochDistances := make(epochDistances, len(ec))
	for i, _ := range epochDistances {
		// fill the distance slice
		distance := epochsOnly[i] - number
		// get abs this way, math.Abs means lots of float64 conversions.
		if distance < 0 {
			distance = distance * -1
		}
		// preserve original array position
		epochDistances[i] = []int64{distance, int64(i)}
	}

	// sort epoch distances. The second array element holds the original position.
	sort.Sort(preserveSecondEl(epochDistances))
	// insertionIndex := SearchInt64s(epochsOnly, number)
	// at this point, the distances are sorted, and we can use the second dimension - the original position
	// to determine what order the ecOut should be in
	for i, distEntry := range epochDistances {
		indexForEcOut := int(distEntry[1])
		ecOut[i] = sorted[indexForEcOut]
	}

	return ecOut
}

// accepts slice of strings, tries to clean them by removing common characters, and returns a list of int64s.
func stringSliceToInt64Base10s(stringsToConvert []string) (numbers []int64, badStrings []string, err error) {
	for _, s := range stringsToConvert {
		s = strings.Trim(s, " \r\n\t")
		//may have a period at the end, with numbers indicating ms. Strip this off.
		s = strings.SplitN(s, ".", 1)[0]
		num, cErr := getInt64Base10(s)
		if cErr != nil {
			badStrings = append(badStrings, s)
		} else {
			numbers = append(numbers, num)
		}
	}
	if len(badStrings) > 0 {
		err = fmt.Errorf("Some strings not converted, %s", badStrings)
	}
	return numbers, badStrings, err
}

func getInt64Base10(parseMe string) (number int64, err error) {
	number, err = strconv.ParseInt(parseMe, 10, 64)
	return number, err
}

// a strange sort of way to sort a 2d array. It preserves the second array element's position.
type epochDistances [][]int64
type preserveSecondEl epochDistances

func (a preserveSecondEl) Len() int {
	return len(a)
}
func (a preserveSecondEl) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a preserveSecondEl) Less(i, j int) bool {
	return a[i][0] < a[j][0]
}

// Given a slice of strings which could have integer data, create a new slice of only numbers in any of the strings
func NumbersInStrings(stringsToClean []string) (numbersOnly []string){
	re := regexp.MustCompile("[0-9]+")
	numbersOnly = make([]string, 0)

	for _, s := range stringsToClean {
		nums := re.FindAllString(s, -1)
		numbersOnly = append(numbersOnly, nums...)
	}
	return numbersOnly
}


