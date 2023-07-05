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
type YMDFlag struct {
	yyyymmdd int // internal yyyymmdd value, nil values might be mutated
}

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
	t, err := time.Parse("20060102", value)
	if err != nil {
		return err
	}
	ymd.yyyymmdd = 10000*t.Year() + 100*int(t.Month()) + t.Day()
	return nil
}

// NewYMDFlag creates a new YMDFlag for the given time.Time's date.
func NewYMDFlag(t time.Time) YMDFlag {
	var ymd YMDFlag
	if t.IsZero() {
		ymd.yyyymmdd = 0
	} else {
		ymd.yyyymmdd = 10000*t.Year() + 100*int(t.Month()) + t.Day()
	}
	return ymd
}

// NewYMDFlag creates a new YMDFlag for the given integral `YYYYMMDD` value, for example `20230704`.
// No validation is performed.
func NewYMDFlagFromInt(i int) YMDFlag {
	return YMDFlag{yyyymmdd: i}
}

// IsZero returns true if the YMDFlag is nil.
func (ymd YMDFlag) IsZero() bool {
	return (ymd.yyyymmdd == 0)
}

// AsYMD returns the YMDFlag as integer `YYYYMMDD`.
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current local date.
func (ymd *YMDFlag) AsYMD() int {
	ymd.UpdateNilToNow()
	return ymd.yyyymmdd
}

// AsYMDString returns the YMDFlag as string `"YYYYMMDD"`
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current local date.
func (ymd *YMDFlag) AsYMDString() string {
	return strconv.Itoa(ymd.AsYMD())
}

// AsDirPath returns the YMDFlag as `"YYYY/MM/DD"` using the OS path seperator
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current local date.
func (ymd *YMDFlag) AsDirPath() string {
	return formatDirPath(ymd.AsTime())
}

// AsDirPathNoCheck returns the YMDFlag as `"YYYY/MM/DD"` using the OS path seperator
// Note: This method does not check if zeroed. Please ensure you call it with a non-zero YMDFlag
func (ymd YMDFlag) AsDirPathNoCheck() string {
	return formatDirPath(ymd.AsLocalTimeNoCheck())
}

// AsLocalTime returns the YMDFlag as a time.Time.
// If the YMDFlag is nil, then a date fetch occurs, updating it to the current local date.
func (ymd *YMDFlag) AsTime() time.Time {
	ymd.UpdateNilToNow()
	return ymd.AsLocalTimeNoCheck()
}

// AsTimeNoCheck returns the YMDFlag as time.Time in Local time
// NOTE: This method does not check if zeroed.
// Please ensure you call it with a non-zero YMDFlag.
func (ymd YMDFlag) AsLocalTimeNoCheck() time.Time {
	var year int = ymd.yyyymmdd / 10000
	var month int = (ymd.yyyymmdd % 10000) / 100
	var day int = ymd.yyyymmdd % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

// AsTimeNoCheck returns the YMDFlag as time.Time in UTC time
// NOTE: This method does not check if zeroed.
// Please ensure you call it with a non-zero YMDFlag.
func (ymd YMDFlag) AsUTCTimeNoCheck() time.Time {
	var year int = ymd.yyyymmdd / 10000
	var month int = (ymd.yyyymmdd % 10000) / 100
	var day int = ymd.yyyymmdd % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// UpdateNilToNow updates a nil YMDFlag to the current local date.
func (ymd *YMDFlag) UpdateNilToNow() {
	if ymd.yyyymmdd == 0 {
		ymd.yyyymmdd = NewYMDFlag(time.Now()).yyyymmdd
	}
}

//////////////////////////////////////////////////////////////////////////////

// formatDirPath returns the `time` as `"YYYY/MM/DD"` using the OS path seperator.
func formatDirPath(time time.Time) string {
	return time.Format(fmt.Sprintf(
		"2006%c01%c02", os.PathSeparator, os.PathSeparator))
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
