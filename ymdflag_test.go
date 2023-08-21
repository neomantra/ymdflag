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

func TestAsDirPath(t *testing.T) {

	flag := NewYMDFlag(time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC))
	path := flag.AsDirPath()

	assert.Equal(t, "2022/01/02", path)
}
