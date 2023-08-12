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

	var timeValue = ymdFlag.AsTimeNoCheck()
	assert.Equal(t, time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), timeValue, "should not have a time component")

	var dirPath = formatDirPath(timeValue, '/')
	assert.Equal(t, "2020/01/02", dirPath, "should match given date path")
}
