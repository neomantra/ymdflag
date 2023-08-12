package ymdflag

// Copyright (c) 2023 Neomantra BV

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode"
)

// YMDFlag represents a Golang flag.Value for `YYYYMMDD`-specified dates.
//
// To faciliates the use of YMD dates in command line flags, it implements the
// [flag.Value interface], making it compatible with the [flag] and [pflag] packages.
//
// It has a `yyyymmddd` integral part and a `loc` location part.
// if the `yyyymmdd` part is 0, that implies a date fetch on the first request for the time.
// If the location is nil, then the local timezone is used.  Otherwise, it is used when
// extracting times from the YMDFlag.
//
// [flag.Value interface]: https://pkg.go.dev/flag#Value
// [flag]: https://pkg.go.dev/flag
// [pflag]: https://pkg.go.dev/github.com/spf13/pflag
type YMDFlag struct {
	yyyymmdd int            // internal yyyymmdd value, nil values might be mutated
	loc      *time.Location // internal location value, nil value means local time
}

///////////////////////////////////////////////////////////////////////////////

// TODO: internal error consts how?

// YMDtoTime returns the Time corresponding to the YYYYMMDD in the specified location.
// A nil location implies local time.
func YMDToTime(yyyymmdd int, loc *time.Location) time.Time {
	var year int = yyyymmdd / 10000
	var month int = (yyyymmdd % 10000) / 100
	var day int = yyyymmdd % 100
	if loc == nil {
		loc = time.Local
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

// TimeToYMD returns the YYYYMMDD for the Time in its location.
// A zero time returns a 0 value.
func TimeToYMD(t time.Time) int {
	if t.IsZero() {
		return 0
	} else {
		return 10000*t.Year() + 100*int(t.Month()) + t.Day()
	}
}

// ValidateYMD returns nil if the passed `yyyymmdd` is of a proper YYYYMMDD form.  Zero is a valid value, indicating auto-detection.
// Otherwise, returns an error.
// This function is not forgiving like `time.Date`, e.g. 10/32 (Oct 32) is not considered 11/01 (Nov 1).
func ValidateYMD(yyyymmdd int) error {
	if yyyymmdd == 0 {
		return nil
	} else if yyyymmdd < 0 {
		return fmt.Errorf("yyyymmdd is negative")
	} else if yyyymmdd > 99999999 {
		return fmt.Errorf("yyyymmdd is more than 8 digits")
	}
	var year int = yyyymmdd / 10000
	var month int = (yyyymmdd % 10000) / 100
	var day int = yyyymmdd % 100
	dt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if year != dt.Year() || month != int(dt.Month()) || day != dt.Day() {
		return fmt.Errorf("yyyymmdd is bad or unnormalized")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// flag.Value interface

// Type implements pflag.Value.Type.  Returns "YMDFlag".
func (*YMDFlag) Type() string {
	return "YMDFlag"
}

// String implements the flag.Value interface.
// If the YMDFlag is nil, then a date fetch occurs,
// updating it to the current local date.
func (ymd *YMDFlag) String() string {
	ymd.UpdateNilToNow()
	return strconv.Itoa(ymd.yyyymmdd)
}

// Set implements the flag.Value interface.
// The default value of empty string `""` is the current local date.
func (ymd *YMDFlag) Set(value string) error {
	// default value (empty string) is today
	if len(value) == 0 {
		ymd.yyyymmdd = 0
		ymd.UpdateNilToNow()
		return nil
	}
	if len(value) != 8 || !isInt(value) {
		return fmt.Errorf("expect string of format YYYYMMDD")
	}
	loc := ymd.loc
	if loc == nil {
		loc = time.Local
	}
	t, err := time.ParseInLocation("20060102", value, loc)
	if err != nil {
		return err
	}
	ymd.yyyymmdd = TimeToYMD(t)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// YMDFlag implementation

// NewYMDFlag creates a new YMDFlag for the given time.Time's date and location.
func NewYMDFlag(t time.Time) YMDFlag {
	var ymd YMDFlag
	ymd.yyyymmdd = TimeToYMD(t)
	ymd.loc = t.Location()
	return ymd
}

// NewYMDFlagWithLocation creates a new nil YMDFlag with the given location.
// This allows preparing a YMDFlag for a specific location before using in a `pflag` function call.
func NewYMDFlagWithLocation(loc *time.Location) YMDFlag {
	return YMDFlag{yyyymmdd: 0, loc: loc}
}

// NewYMDFlag creates a new YMDFlag for the given integral `YYYYMMDD` value, for example `20230704`.
// Returns a non-nil error if YMDFlag is malformed.  `0` is a valid value.
func NewYMDFlagFromInt(i int, loc *time.Location) (YMDFlag, error) {
	if err := ValidateYMD(i); err != nil {
		return YMDFlag{}, err
	}
	return YMDFlag{yyyymmdd: i, loc: loc}, nil
}

// GetYMD returns the YMDFlag as integer `YYYYMMDD`.  It may be zero.
func (ymd YMDFlag) GetYMD() int {
	return ymd.yyyymmdd
}

// GetLocation returns the location of the YMDFlag.  It may be nil.
func (ymd YMDFlag) GetLocation() *time.Location {
	return ymd.loc
}

// SetLocation sets the location of the YMDFlag, which affects future calls to AsTime.
func (ymd *YMDFlag) SetLocation(loc *time.Location) {
	ymd.loc = loc
}

// IsZero returns true if the YMDFlag is nil.  The location is ignored in this case.
func (ymd YMDFlag) IsZero() bool {
	return (ymd.yyyymmdd == 0)
}

// AsYMD returns the YMDFlag as integer `YYYYMMDD`.
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current date according to the YMDFlag timezone.
func (ymd *YMDFlag) AsYMD() int {
	ymd.UpdateNilToNow()
	return ymd.yyyymmdd
}

// AsYMDString returns the YMDFlag as string `"YYYYMMDD"`
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current date according to the YMDFlag timezone.
func (ymd *YMDFlag) AsYMDString() string {
	return strconv.Itoa(ymd.AsYMD())
}

// AsDirPath returns the YMDFlag as `"YYYY/MM/DD"` using the OS path seperator
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current date according to the YMDFlag timezone.
func (ymd *YMDFlag) AsDirPath() string {
	return formatDirPath(ymd.AsTime(), os.PathSeparator)
}

// AsDirPath returns the YMDFlag as `"YYYY/MM/DD"` using given path seperator
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current date according to the YMDFlag timezone.
func (ymd *YMDFlag) AsDirPathSep(separator rune) string {
	return formatDirPath(ymd.AsTime(), separator)
}

// AsTime returns the YMDFlag as a time.Time in the YMDFlag's location.
// If the YMDFlag's location is nil, then the local timezone is used.
// If the YMDFlag's YMD is 0, then a date fetch occurs, updating it to the current local date.
func (ymd *YMDFlag) AsTime() time.Time {
	ymd.UpdateNilToNow()
	return YMDToTime(ymd.yyyymmdd, ymd.loc)
}

// AsTimeNoCheck returns the YMDFlag as time.Time in the YMDFlag's location.
// If the YMDFlag's location is nil, then the local timezone is used.
// NOTE: This method does not check if zeroed.  Ensure you call it with a non-zero YMDFlag.
func (ymd YMDFlag) AsTimeNoCheck() time.Time {
	return YMDToTime(ymd.yyyymmdd, ymd.loc)
}

// UpdateNilToNow updates a nil YMDFlag to the current local date.
func (ymd *YMDFlag) UpdateNilToNow() {
	if ymd.yyyymmdd == 0 {
		now := time.Now()
		if ymd.loc != nil {
			now = now.In(ymd.loc)
		}
		ymd.yyyymmdd = TimeToYMD(now)
	}
}

//////////////////////////////////////////////////////////////////////////////

// formatDirPath returns the `time` as `"YYYY/MM/DD"` using the given path seperator.
func formatDirPath(time time.Time, sep rune) string {
	return time.Format(fmt.Sprintf("2006%c01%c02", sep, sep))
}

// isInt checks if a string can be converted safely to an int
func isInt(value string) bool {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
