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
	assert.NotPanics(t, func() { ymdFlag.AsTime() })
	assert.NotPanics(t, func() { ymdFlag.AsYMD() })
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
