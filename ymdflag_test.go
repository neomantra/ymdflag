package ymdflag

// Copyright (c) 2023 Neomantra BV

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYMDFlag(t *testing.T) {
	var ymdFlag YMDFlag

	assert.True(t, ymdFlag.IsZero())
	assert.NotPanics(t, func() { _ = ymdFlag.AsTime() })
	assert.NotPanics(t, func() { _ = ymdFlag.AsYMD() })
}

func Test_uninitializaed_flag_becomes_today_when_accessed(t *testing.T) {
	var ymdFlag YMDFlag
	var expected = time.Now()

	var result = ymdFlag.AsTime()

	assert.Equal(t, expected.Year(), result.Year())
	assert.Equal(t, expected.Month(), result.Month())
	assert.Equal(t, expected.Day(), result.Day())
}

func TestNonMutatingMethods(t *testing.T) {
	ymdFlag := NewYMDFlag(time.Date(2020, time.January, 2, 1, 2, 3, 4, time.UTC))

	var timeValue = ymdFlag.AsTime()
	assert.Equal(t, time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local), timeValue, "should not have a time component")

	var dirPath = FormatDirPath(ymdFlag, '/')
	assert.Equal(t, "2020/01/02", dirPath, "should match given date path")
}

func TestValidateYMD(t *testing.T) {
	err := ValidateYMD(20220101)
	assert.NoError(t, err, "valid date should not return an error")

	err = ValidateYMD(0)
	assert.NoError(t, err, "zero is ok")

	err = ValidateYMD(20240229)
	assert.NoError(t, err, "leap day is ok")

	err = ValidateYMD(20230229)
	assert.Error(t, err, "leap day on wrong year is not ok")

	err = ValidateYMD(209901231)
	assert.Error(t, err, "has more than 8 digits")

	err = ValidateYMD(010101)
	assert.Error(t, err, "too few digits")

	err = ValidateYMD(20221301)
	assert.Error(t, err, "invalid month")

	err = ValidateYMD(20221241)
	assert.Error(t, err, "invalid day")

	err = ValidateYMD(-1)
	assert.Error(t, err, "negative date")
}

func TestAsYearMonthDay(t *testing.T) {

	// default is zero
	flag := YMDFlag{}
	year, month, day := flag.AsYearMonthDay()
	assert.Equal(t, 0, year)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, day)

	flag = NewYMDFlag(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC))
	year, month, day = flag.AsYearMonthDay()
	assert.Equal(t, 2022, year)
	assert.Equal(t, 1, month)
	assert.Equal(t, 2, day)
}

func TestNewFlagFromInt(t *testing.T) {
	flag, err := NewYMDFlagFromInt(0)
	assert.NoError(t, err, "zero is ok")

	year, month, day := flag.AsYearMonthDay()
	assert.Equal(t, 0, year)
	assert.Equal(t, 0, month)
	assert.Equal(t, 0, day)

	flag, err = NewYMDFlagFromInt(20220102)
	assert.NoError(t, err, "valid date is ok")
	assert.Equal(t, 20220102, flag.AsYMD())

	_, err = NewYMDFlagFromInt(20221301)
	assert.Error(t, err, "invalid month should return an error")

	_, err = NewYMDFlagFromInt(-1)
	assert.Error(t, err, "negative should return an error")
}

func TestFlagValueRoundTrip(t *testing.T) {
	var flag YMDFlag

	// Type is used by pflag to describe the flag
	assert.Equal(t, "YMDFlag", flag.Type())

	// zero value stringifies to empty
	assert.Equal(t, "", flag.String())

	err := flag.Set("20220102")
	assert.NoError(t, err, "valid date should Set without error")
	assert.Equal(t, 20220102, flag.AsYMD())
	assert.Equal(t, "20220102", flag.String())

	err = flag.Set("")
	assert.NoError(t, err, "empty string clears the flag")
	assert.True(t, flag.IsZero())

	err = flag.Set("2022-01-02")
	assert.Error(t, err, "invalid format should return an error")
}

func TestGetYMD(t *testing.T) {
	flag := NewYMDFlag(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC))
	assert.Equal(t, 20220102, flag.GetYMD())
	assert.Equal(t, flag.AsYMD(), flag.GetYMD())
}

func TestFormatDirPathZero(t *testing.T) {
	assert.Equal(t, "", FormatDirPath(YMDFlag{}, '/'), "zero flag yields empty string")
	flag := NewYMDFlag(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, "2022-01-02", FormatDirPath(flag, '-'), "honors the separator")
}

func TestYMDToTime(t *testing.T) {
	assert.True(t, YMDToTime(0, time.UTC).IsZero(), "zero yields the zero time")

	tm := YMDToTime(20220102, time.UTC)
	assert.Equal(t, time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC), tm)

	// nil location implies local time
	tm = YMDToTime(20220102, nil)
	assert.Equal(t, time.Local, tm.Location())
}

func TestTimeToYMD(t *testing.T) {
	assert.Equal(t, 0, TimeToYMD(time.Time{}), "zero time yields 0")
	assert.Equal(t, 20220102, TimeToYMD(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)))
}

func TestUpdateNilToNow(t *testing.T) {
	// zero flag becomes today
	var flag YMDFlag
	flag.UpdateNilToNow(time.UTC)
	assert.False(t, flag.IsZero())
	assert.Equal(t, TimeToYMD(time.Now().In(time.UTC)), flag.AsYMD())

	// a set flag is left untouched
	set := NewYMDFlag(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC))
	set.UpdateNilToNow(time.UTC)
	assert.Equal(t, 20220102, set.AsYMD(), "already-set flag should not change")
}

func TestAsTimeRawWithLoc(t *testing.T) {
	// zero flag stays zero — no auto-update, unlike AsTime
	var flag YMDFlag
	tm := flag.AsTimeRawWithLoc(time.UTC)
	assert.True(t, tm.IsZero(), "raw accessor does not auto-populate")
	assert.True(t, flag.IsZero(), "raw accessor does not mutate the flag")

	set := NewYMDFlag(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC))
	tm = set.AsTimeRawWithLoc(time.UTC)
	assert.Equal(t, time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC), tm)

	// nil location implies local time
	tm = set.AsTimeRawWithLoc(nil)
	assert.Equal(t, time.Local, tm.Location())
}

func TestStringToYMD(t *testing.T) {

	yyyymmdd, err := StringToYMD("20220101")
	assert.NoError(t, err, "valid date should not return an error")
	assert.Equal(t, 20220101, yyyymmdd)

	_, err = StringToYMD("2022-01-01")
	assert.Error(t, err, "date string should not have separators")

	_, err = StringToYMD("hello world")
	assert.Error(t, err, "non numeric string should return an error")

	_, err = StringToYMD("123456789")
	assert.Error(t, err, "too long string")

	yyyymmdd, err = StringToYMD("")
	assert.NoError(t, err, "empty string should not return an error")
	assert.Equal(t, 0, yyyymmdd)
}
