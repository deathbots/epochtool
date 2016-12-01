package epochconv

import (
	"testing"
)

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
		if tt.in.RightNowInSecondsSince < 1 {
			t.Errorf("Epoch time %s, from Epoch Date %s, was not greater than one", tt.in.RightNowInSecondsSince, tt.in.EpochName)
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


func epochInSlice(s []EpochType, e EpochType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
