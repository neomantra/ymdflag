package ymdflag

// Copyright (c) 2023 Neomantra BV

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

// YMDFlag represents a Golang flag.Value for `YYYYMMDD`-specified dates.
//
// To faciliates the use of YMD dates in command line flags, it implements the
// [flag.Value interface], making it compatible with the [flag] and [pflag] packages.
//
// It stores an integral `yyyymmddd`.  The special value of 0 indicates that the value
// is indeterminate and may be may be auto-populated by `UpdateNilToNow`, `AsTime`, or `AsTimeWithLoc`.
//
// [flag.Value interface]: https://pkg.go.dev/flag#Value
// [flag]: https://pkg.go.dev/flag
// [pflag]: https://pkg.go.dev/github.com/spf13/pflag
type YMDFlag struct {
	yyyymmdd int // internal yyyymmdd value, nil values might be mutated
}

///////////////////////////////////////////////////////////////////////////////

// TODO: internal error consts how?

// YMDtoTime returns the Time corresponding to the YYYYMMDD in the specified location, without validating the argument.`
// A value of 0 returns a Zero Time, independent of location.
// A nil location implies local time.
func YMDToTime(yyyymmdd int, loc *time.Location) time.Time {
	if yyyymmdd == 0 {
		return time.Time{}
	}
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

// StringToYMD returns an integral YYYYMMDD value or 0 for an empty string.
// If the string is invalid, an error is returned.
func StringToYMD(str string) (int, error) {
	// default value (empty string) is 0
	if str == "" {
		return 0, nil
	}

	if len(str) != 8 || !isInt(str) {
		return 0, fmt.Errorf("expect string of format YYYYMMDD")
	}

	yyyymmdd, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string %w", err)
	}

	if err := ValidateYMD(yyyymmdd); err != nil {
		return 0, fmt.Errorf("failed to validate string %w", err)
	}
	return yyyymmdd, nil
}

// ValidateYMD returns nil if the passed `yyyymmdd` is of a proper YYYYMMDD form.
// Zero is a valid value, meaning indeindicating potential auto-detection.
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

// AsDirPath returns the YMDFlag as `"YYYY/MM/DD"` using given path seperator
// If the YMDFlag is nil, then an empty string is returned.
func FormatDirPath(ymd YMDFlag, separator rune) string {
	if ymd.IsZero() {
		return ""
	}
	year, month, day := ymd.AsYearMonthDay()
	return fmt.Sprintf("%04d%c%02d%c%02d", year, separator, month, separator, day)
}

///////////////////////////////////////////////////////////////////////////////
// flag.Value interface

// Type implements pflag.Value.Type.  Returns "YMDFlag".
func (*YMDFlag) Type() string {
	return "YMDFlag"
}

// String implements the flag.Value interface.
// If the YMDFlag is nil, then an empty string is returned.
func (ymd *YMDFlag) String() string {
	return ymd.AsYMDString()
}

// Set implements the flag.Value interface.
// The default value of empty string `""` implies it is unset
// and may be auto-filled by some methods.
func (ymd *YMDFlag) Set(value string) error {
	// convert value to YMD int
	yyyymmdd, err := StringToYMD(value)
	if err != nil {
		return err
	}
	ymd.yyyymmdd = yyyymmdd
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// YMDFlag implementation

// NewYMDFlag creates a new YMDFlag for the given time.Time's date and location.
func NewYMDFlag(t time.Time) YMDFlag {
	var ymd YMDFlag
	ymd.yyyymmdd = TimeToYMD(t)
	return ymd
}

// NewYMDFlag creates a new YMDFlag for the given integral `YYYYMMDD` value, for example `20230704`.
// Returns a non-nil error if YMDFlag is malformed.  `0` is a valid value.
func NewYMDFlagFromInt(i int) (YMDFlag, error) {
	if err := ValidateYMD(i); err != nil {
		return YMDFlag{}, err
	}
	return YMDFlag{yyyymmdd: i}, nil
}

// GetYMD returns the YMDFlag as integer `YYYYMMDD`.  It may be zero.
func (ymd YMDFlag) GetYMD() int {
	return ymd.yyyymmdd
}

// IsZero returns true if the YMDFlag is nil.  The location is ignored in this case.
func (ymd YMDFlag) IsZero() bool {
	return (ymd.yyyymmdd == 0)
}

// AsYMD returns the YMDFlag as integer `YYYYMMDD`.  Returns 0 if the YMDFlag is nil.
func (ymd YMDFlag) AsYMD() int {
	return ymd.yyyymmdd
}

// AsYMDString returns the YMDFlag as string `"YYYYMMDD"`.  If the YMDFlag is nil, it returns the empty string.
func (ymd YMDFlag) AsYMDString() string {
	if ymd.yyyymmdd == 0 {
		return ""
	}
	return strconv.Itoa(ymd.yyyymmdd)
}

// AsYearMonthDay returns the YMDFlag decomposed into Year, Month, and Day.
// All values of 0 will be returned if the YMDFlag is 0
func (ymd YMDFlag) AsYearMonthDay() (int, int, int) {
	if ymd.IsZero() {
		return 0, 0, 0
	}
	var year int = ymd.yyyymmdd / 10000
	var month int = (ymd.yyyymmdd % 10000) / 100
	var day int = ymd.yyyymmdd % 100
	return year, month, day
}

// UpdateNilToNow updates a nil YMDFlag (with `yyyymmdd` == 0) to the current date in the specified location.
// If location is nil, local time is used.
// If `yyyymmdd` is not nil, then this method does nothing.
func (ymd *YMDFlag) UpdateNilToNow(location *time.Location) {
	if ymd.yyyymmdd != 0 {
		return
	}
	now := time.Now()
	if location != nil {
		now = now.In(location)
	}
	ymd.yyyymmdd = TimeToYMD(now)
}

// AsTime returns the YMDFlag as a `time.Time“ in local time.  Use `AsTimeWithLoc` to specify a location.
// If the YMDFlag's `yyyymmdd` is 0, then a zero time in that location is returned.
func (ymd *YMDFlag) AsTime() time.Time {
	return ymd.AsTimeWithLoc(nil)
}

// AsTime returns the YMDFlag as a `time.Time` in the specified location.
// If the YMDFlag's `yyyymmdd` is 0, then a zero time in that location is returned.
// If `location“ is nil, then `time.Local` is used.
func (ymd *YMDFlag) AsTimeWithLoc(location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	ymd.UpdateNilToNow(location)
	return YMDToTime(ymd.yyyymmdd, location)
}

//////////////////////////////////////////////////////////////////////////////

// isInt checks if a string can be converted safely to an int
func isInt(value string) bool {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
