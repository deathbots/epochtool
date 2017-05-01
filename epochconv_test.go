package epochconv

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

// used to determine the strings are parseable because particular errors are swallowed
// at struct init type. Failure here would mean the package wouldn't be able to load.
var timeStringTests = []struct {
	in EpochType
}{
	{EpochCommonEra},
	{EpochWindowsEpoch},
	{EpochVMS},
	{EpochMicrosoftCOM},
	{EpochMicrosoftExcel},
	{EpochNTP},
	{EpochMacClassic},
	{EpochUnix},
	{EpochFAT},
	{EpochGPS},
	{EpochPostgreSQL},
	{EpochMacOSX},
}

// Tests whether the time string constants are parseable. If they are not, when sp and te functions are used to
// initialize the struct literal, you would get Time that prints like 0001-01-01 00:00:00 +0000 UTC and the empty string
// respectively.
func TestTimeStringsParseable(t *testing.T) {
	for _, tt := range timeStringTests {
		if tt.in.EpochDate.IsZero() && tt.in.EpochName != "CommonEra" { // This fails for CommonEra since it's 0
			t.Errorf("Time was not initialized properly for Epoch Date %s, check format string constant used", tt.in.EpochName)
		}
		if tt.in.LocalRightNowInSecondsSince < 1 {
			t.Errorf("Epoch time %d, from Epoch Date %s, was not greater than one", tt.in.LocalRightNowInSecondsSince, tt.in.EpochName)
		}
	}
}

// Tests whether every struct in timeStringTests is included in AllEpochs.
func TestAllEpochsContainsAll(t *testing.T) {
	for _, e := range timeStringTests {
		if !epochInSlice(AllEpochs, e.in) {
			t.Errorf("Epoch Type %s was not in the AllEpochs Slice", e.in.EpochName)
		}

	}
}

// Tests whether top result is correctly picked for a known epoch
func TestFirstSortedEpochFromUnix(t *testing.T) {
	sorted := AllEpochs.OrderedEpochsByClosestMatch(0, time.Unix(0, 0).UTC())
	if sorted[0].EpochName != "Unix" {
		t.Errorf("Epoch Type %s was not incorrect for testing against Unix timestamp", sorted[0].EpochName)
	}
}

// Tests whether epochs sorting is done correctly.
func TestEpochsSortByDistance(t *testing.T) {
	epochStartTime, err := time.Parse(time.RFC3339, dateStringCommonEra)
	if err != nil {
		t.Errorf("Could not parse epoch start time: %s", err)
	}
	sorted := AllEpochs.OrderedEpochsByClosestMatch(0, epochStartTime)

	epochSeconds := make([]int64, len(sorted))
	for i, et := range sorted {
		epochSeconds[i] = et.LocalRightNowInSecondsSince
	}
	// these should be in descending order
	last := int64(math.MaxInt64)
	for _, v := range epochSeconds {
		if v > last {
			t.Errorf("Epoch collection Sorting was done incorrectly: %s", err)
		}
		last = v
	}
	for _, sec := range epochSeconds {
		fmt.Printf("%d\n", sec)
	}
}

func epochInSlice(s []EpochType, e EpochType) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}
