package epochconv

import (
	"time"
	"fmt"
	"strings"
	"encoding/json"
)

// Holds types of epochs

// The example time formatting string fed to time parse method. This allows custom epochs to be calculated as long as
// the epoch's date format looks precisely like this string - YYYY-MM-DDTHH:MM:SSZ. All epochs in this file must
// follow this format.
const CustomEpochTimeFormatString = "2006-01-02T15:04:05Z"


// Dates which conform to the above formatting string. These must, of course, all be in the past - or there is an
// outside risk the program would crash due to all numbers being positive.
const
(
	dateStringCommonEra = "0001-01-01T00:00:00Z"
	dateStringUnixEpoch = "1970-01-01T00:00:00Z"
	dateStringWindowsEpoch = "1601-01-01T00:00:00Z"
	dateStringVMSEpoch = "1858-11-17T00:00:00Z"
	dateStringMicrosoftCOM = "1899-12-30T00:00:00Z"
	dateStringMicrosoftExcel = "1899-12-31T00:00:00Z"
	dateStringNTP = "1900-01-01T00:00:00Z"
	dateStringMacClassic = "1904-01-01T00:00:00Z"
	dateStringMicrosoftFAT = "1980-01-01T00:00:00Z"
	dateStringGPS = "1980-01-06T00:00:00Z"
	dateStringPostgreSQL = "2000-01-01T00:00:00Z"
	dateStringMacOSX = "2001-01-01T00:00:00Z"
)

type EpochCollection []EpochType

// Skeletal type
type EpochType struct {
	EpochName                   string    `json:"epoch_name"`// Friendly name of epoch
	EpochUses                   []string  `json:"epoch_uses"`// Slice of common uses of this specific epoch
	EpochDateString             string    `json:"-"`		// The date string formatted like CustomEpochTimeFormatString that defines this
	EpochDate                   time.Time `json:"epoch_date"`// The time.Time date representation of the epoch start
	LocalRightNowInSecondsSince int64     `json:"now_local"`// time.Now().Local - Local time in seconds since epoch start.
	UTCRightNowInSecondsSince   int64     `json:"now_utc"`// time.Now().UTC - UTC time in seconds since epoch start.
	Prevalence                  int       `json:"prevalence"`// 0-5, 0 being least common. Helps decide most likely matches when it's close.
}

var (
	EpochCommonEra = EpochType{
		EpochName:"CommonEra",
		EpochUses: []string{"Common Era", "ISO 2014", "RFC 3339", "Microsoft .NET", "Go", "REXX", "Rata Die"},
		EpochDateString: dateStringCommonEra,
		EpochDate: sp(dateStringCommonEra),
		LocalRightNowInSecondsSince: te(dateStringCommonEra, false),
		UTCRightNowInSecondsSince: te(dateStringCommonEra, true),
		Prevalence: 1,
	}

	EpochUnix = EpochType{
		EpochName:"Unix",
		EpochUses: []string{"Unix", "Unix Variants (Linux, MacOS, Solaris, BSD, etc...)", "POSIX"},
		EpochDateString: dateStringUnixEpoch,
		EpochDate: sp(dateStringUnixEpoch),
		LocalRightNowInSecondsSince: te(dateStringUnixEpoch, false),
		UTCRightNowInSecondsSince: te(dateStringUnixEpoch, true),
		Prevalence: 5,
	}

	EpochWindowsEpoch = EpochType{
		EpochName:"Windows",
		EpochUses: []string{"Windows", "NTFS", "COBOL"},
		EpochDateString: dateStringWindowsEpoch,
		EpochDate: sp(dateStringWindowsEpoch),
		LocalRightNowInSecondsSince: te(dateStringWindowsEpoch, false),
		UTCRightNowInSecondsSince: te(dateStringWindowsEpoch, true),
		Prevalence: 5,
	}

	EpochVMS = EpochType{
		EpochName:"VMS",
		EpochUses: []string{"VMS", "United States Naval Observatory", "DVB SI 16-bit day stamps", "Astronomy-related"},
		EpochDateString: dateStringVMSEpoch,
		EpochDate: sp(dateStringVMSEpoch),
		LocalRightNowInSecondsSince: te(dateStringVMSEpoch, false),
		UTCRightNowInSecondsSince: te(dateStringVMSEpoch, true),
		Prevalence: 3,
	}

	EpochMicrosoftCOM = EpochType{
		EpochName:"Microsoft COM",
		EpochUses: []string{"Microsoft COM DATE", "Object Pascal", "LibreOffice Calc", "Google Sheets", "Technical internal value used by Microsoft Excel"},
		EpochDateString: dateStringMicrosoftCOM,
		EpochDate: sp(dateStringMicrosoftCOM),
		LocalRightNowInSecondsSince: te(dateStringMicrosoftCOM, false),
		UTCRightNowInSecondsSince: te(dateStringMicrosoftCOM, true),
		Prevalence: 4,
	}

	EpochMicrosoftExcel = EpochType{
		EpochName:"Microsoft Excel",
		EpochUses: []string{"Microsoft Excel", "Lotus 1-2-3"},
		EpochDateString: dateStringMicrosoftExcel,
		EpochDate: sp(dateStringMicrosoftExcel),
		LocalRightNowInSecondsSince: te(dateStringMicrosoftExcel, false),
		UTCRightNowInSecondsSince: te(dateStringMicrosoftExcel, true),
		Prevalence: 3,
	}

	EpochNTP = EpochType{
		EpochName:"NTP",
		EpochUses: []string{"Network Time Protocol", "IBM CICS", "Mathematica", "RISC OS", "VME", "Common Lisp", "Michigan Terminal System"},
		EpochDateString: dateStringNTP,
		EpochDate: sp(dateStringNTP),
		LocalRightNowInSecondsSince: te(dateStringNTP, false),
		UTCRightNowInSecondsSince: te(dateStringNTP, true),
		Prevalence: 2,
	}

	EpochMacClassic = EpochType{
		EpochName:"Mac Classic",
		EpochUses: []string{"Apple Inc.'s classic Mac OS, LabVIEW, Palm OS, MP4, Microsoft Excel (optionally), IGOR Pro"},
		EpochDateString: dateStringMacClassic,
		EpochDate: sp(dateStringMacClassic),
		LocalRightNowInSecondsSince: te(dateStringMacClassic, false),
		UTCRightNowInSecondsSince: te(dateStringMacClassic, true),
		Prevalence: 2,
	}

	EpochFAT = EpochType{
		EpochName:"FAT",
		EpochUses: []string{"FAT12", "FAT16", "FAT32", "exFAT filesystems", "IBM BIOS", "INT 1Ah", "DOS", "OS/2", },
		EpochDateString: dateStringMicrosoftFAT,
		EpochDate: sp(dateStringMicrosoftFAT),
		LocalRightNowInSecondsSince: te(dateStringMicrosoftFAT, false),
		UTCRightNowInSecondsSince: te(dateStringMicrosoftFAT, true),
		Prevalence: 5,
	}

	// This is very close to FAT
	EpochGPS = EpochType{
		EpochName:"GPS",
		EpochUses: []string{"Qualcomm BREW", "GPS", "ATSC 32-bit time stamps"},
		EpochDateString: dateStringGPS,
		EpochDate: sp(dateStringGPS),
		LocalRightNowInSecondsSince: te(dateStringGPS, false),
		UTCRightNowInSecondsSince: te(dateStringGPS, true),
		Prevalence: 2,
	}
	// This epoch is very close to OS X epoch
	EpochPostgreSQL = EpochType{
		EpochName:"PostgreSQL",
		EpochUses: []string{"PostgreSQL", "AppleSingle", "AppleDouble", "ZigBee UTCTime"},
		EpochDateString: dateStringPostgreSQL,
		EpochDate: sp(dateStringPostgreSQL),
		LocalRightNowInSecondsSince: te(dateStringPostgreSQL, false),
		UTCRightNowInSecondsSince: te(dateStringPostgreSQL, true),
		Prevalence: 3,
	}

	EpochMacOSX = EpochType{
		EpochName:"Mac OS X",
		EpochUses: []string{"OS X, Apple Cocoa"},
		EpochDateString: dateStringMacOSX,
		EpochDate: sp(dateStringMacOSX),
		LocalRightNowInSecondsSince: te(dateStringMacOSX, false),
		UTCRightNowInSecondsSince: te(dateStringMacOSX, true),
		Prevalence: 5,
	}
	AllEpochs = EpochCollection{EpochCommonEra, EpochWindowsEpoch, EpochVMS, EpochMicrosoftCOM, EpochMicrosoftExcel,
		EpochNTP, EpochMacClassic, EpochUnix, EpochFAT, EpochGPS, EpochPostgreSQL, EpochMacOSX}
)


// GuessesForStrings is a method on any EpochCollection, which can be constructed to pick and choose relevant or custom
// EpochTypes.
// Given a slice of strings, return an EpochGuessResults type, which is an array of EpochResults along with the most
// likely result. Strings in the input slice are parsed in the following way:
// 1) Strings are stripped of leading and trailing whitespace characters.
// 2) Strings have all data after the first dot character removed. This allows for input of decimal numbers
//    Without needing to convert to floats.
// If one string cannot be converted, an Error is created indicating at least one string could not be converted. These strings
// are returned in the badStrings slice.
// This can, of course, be ignored - and may be in a typical use case.
func (ec EpochCollection) GuessesForStrings(stringsToConvert []string) (epochResults []EpochResults, badStrings []string, err error) {
	epochResults, badStrings, err = createGuesses(stringsToConvert, ec)
	return epochResults, badStrings, err
}

// String satisfies the Stringer interface, so this is printed when %s is used in a formatting string for this type.
func (e EpochType) String() string {
	return fmt.Sprintf("Name of Epoch: %s\n" +
		"Used for: %s\n" +
		"Started On (UTC): %s\n" +
		"Current UTC Time in Epoch Seconds: %d\n" +
		"Current Local Time in Epoch Seconds: %d\n", e.EpochName, strings.Join(e.EpochUses, ", "),
		e.EpochDate.Format(time.RFC3339), e.UTCRightNowInSecondsSince, e.LocalRightNowInSecondsSince)
}

// String satisfies the Stringer interface, so this is printed when %s is used in a formatting string for this type.
func (e EpochType) ToJson() (string, error) {
	jsonByteArray, err := json.Marshal(&e)
	if err != nil {
		return "", err
	}
	//n := bytes.IndexByte(jsonByteArray, 0)
	//return string(jsonByteArray[:n]), err
	return string(jsonByteArray), err
}

// String satisfies the Stringer interface, so this is printed when %s is used in a formatting string for this type.
func (ec EpochCollection) String() (out string) {
	for _, e := range ec {
		out = out + fmt.Sprintf("%s\n----------\n", e)
	}
	return out
}

func (ec EpochCollection) ToJson() (string, error) {
	jsonByteArray, err := json.Marshal(&ec)
	if err != nil {
		return "", err
	}
	//n := bytes.IndexByte(jsonByteArray, 0)
	//return string(jsonByteArray[:n]), err
	return string(jsonByteArray), err
}

// Satisfy the sort.Interface for the collection type so it can be sort.Sort'ed - like sort.Sort(ByEpochDate(ec))
type ByEpochDate EpochCollection

func (a ByEpochDate) Len() int {
	return len(a)
}
func (a ByEpochDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByEpochDate) Less(i, j int) bool {
	return a[i].EpochDate.Second() < a[j].EpochDate.Second()
}

type ByNearestDate EpochCollection

func (a ByNearestDate) Len() int {
	return len(a)
}
func (a ByNearestDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByNearestDate) Less(i, j int) bool {
	return a[i].EpochDate.Second() < a[j].EpochDate.Second()
}


// DateForNumber is a method on an EpochType. Given a number (in seconds), return the date (as time.Time) for the epoch
func (e *EpochType) DateForNumber(epochSeconds int64, utcFlag bool) (timeInEpoch time.Time) {
	if utcFlag {
		timeInEpoch = e.EpochDate.Add(time.Second * time.Duration(epochSeconds))
	} else {
		// local time
		_, offsetSeconds := time.Now().In(time.Local).Zone()
		timeInEpoch = e.EpochDate.Add(time.Second * time.Duration(epochSeconds+int64(offsetSeconds)))
	}
	return timeInEpoch
}

// NumberForDate is a method on an EpochType. Given a date (as time.Time), return the seconds since that epoch.
func (e *EpochType) NumberForDate(date time.Time) int64 {
	return int64(date.Sub(e.EpochDate).Seconds())
}

// secondsForEpochString returns a specific date in the epoch const formatting string.
func secondsForEpochString(epochConstFormatString string, specificTime time.Time, utcFlag bool) (int64, error) {
	// get epoch date
	epochStart, err := time.Parse(CustomEpochTimeFormatString, epochConstFormatString)
	if err != nil {
		return 0, err
	}
	var dur int64
	if utcFlag {
		dur = int64(specificTime.UTC().Sub(epochStart).Seconds())
	} else {
		dur = int64(specificTime.Local().Sub(epochStart).Seconds())
	}
	if !utcFlag {
		_, offsetSeconds := time.Now().In(time.Local).Zone()
		dur += int64(offsetSeconds)
	}
	return dur, err
}

// sp and te functions are used only for initializing a struct literal, which can only handle a single return value.
// They are unsafe as they swallow errors, and are unexported.
// Tests will check whether any code here cannot parse out - essentially looking for typos in the code.

// sp - SilentParse is used to fill the readable epoch date string in as a time.Time datetime object.
func sp(parseMe string) time.Time {
	parsedTime, _ := time.Parse(CustomEpochTimeFormatString, parseMe)
	return parsedTime
}

// te - TodayEpoch uses dateForEpochString, but swallows its errors. Used to fill Epoch Example time value.
func te(timeInRFC3339ZuluFormat string, utcFlag bool) (epoch int64) {
	epoch, _ = secondsForEpochString(timeInRFC3339ZuluFormat, time.Now().UTC(), utcFlag)
	return epoch
}


/* =========== Epoch Data from Wikipedia! =============

-- indicates it's not going to be used.

Epoch date	Notable uses	Rationale for selection
-- (too rare, and bad college memories) January 0, 1 BC[10]	MATLAB[11]
January 1, AD 1[10]	Microsoft .NET,[12][13] Go,[14] REXX,[15] Rata Die[16]	Common Era, ISO 2014,[17] RFC 3339[18]
January 1, 1601	NTFS, COBOL, Win32/Win64	1601 was the first year of the 400-year Gregorian calendar cycle at the time Windows NT was made.[19]
-- (too rare) December 31, 1840	MUMPS programming language	1841 was a non-leap year several years before the birth year of the oldest living US citizen when the language was designed.[20]
November 17, 1858	VMS, United States Naval Observatory, DVB SI 16-bit day stamps, other astronomy-related computations[21]	November 17, 1858, 00:00:00 UT is the zero of the Modified Julian Day (MJD) equivalent to Julian day 2400000.5[22]
December 30, 1899	Microsoft COM DATE, Object Pascal, LibreOffice Calc, Google Sheets[23]	Technical internal value used by Microsoft Excel; for compatibility with Lotus 1-2-3.[24]
-- (too rare) December 31, 1899	Microsoft C/C++ 7.0[25]	A change in Microsoft’s last version of non-Visual C/C++ that was subsequently reverted.
-- (implemented as Dec 31,1899) January 0, 1900	Microsoft Excel,[24] Lotus 1-2-3[26]	While logically January 0, 1900 is equivalent to December 31, 1899, these systems do not allow users to specify the latter date.
January 1, 1900	Network Time Protocol, IBM CICS, Mathematica, RISC OS, VME, Common Lisp, Michigan Terminal System
January 1, 1904	LabVIEW, Apple Inc.'s classic Mac OS, Palm OS, MP4, Microsoft Excel (optionally),[27] IGOR Pro	1904 is the first leap year of the 20th century.[28]
-- (too rare) December 31, 1967	Pick OS and variants (jBASE, Universe, Unidata, Revelation, Reality)	Chosen so that (date mod 7) would produce 0=Sunday, 1=Monday, 2=Tuesday, 3=Wednesday, 4=Thursday, 5=Friday, and 6=Saturday.[29]
January 1, 1970	Unix Epoch aka POSIX time, used by Unix and Unix-like systems (Linux, macOS), and programming languages: most C/C++ implementations,[30] Java, JavaScript, Perl, PHP, Python, Ruby, Tcl, ActionScript. Also used by Precision Time Protocol.
January 1, 1980	IBM BIOS INT 1Ah, DOS, OS/2, FAT12, FAT16, FAT32, exFAT filesystems	The IBM PC with its BIOS as well as 86-DOS, MS-DOS and PC DOS with their FAT12 file system were developed and introduced between 1980 and 1981
January 6, 1980	Qualcomm BREW, GPS, ATSC 32-bit time stamps	GPS counts weeks (a week is defined to start on Sunday) and January 6 is the first Sunday of 1980.[31][32]
January 1, 2000	AppleSingle, AppleDouble,[33] PostgreSQL,[34] ZigBee UTCTime[35]
January 1, 2001	Apple's Cocoa framework	2001 is the year of the release of Mac OS X 10.0 (but NSDate for Apple's EOF 1.0 was developed in 1994).

*/

