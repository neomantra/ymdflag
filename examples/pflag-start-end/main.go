// Copyright (c) 2023 Neomantra BV

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/neomantra/ymdflag"
	"github.com/spf13/pflag"
)

func main() {
	var startDate ymdflag.YMDFlag
	var endDate ymdflag.YMDFlag

	pflag.VarP(&startDate, "start", "s", "YYYYMMDD start date; defaults to end date")
	pflag.VarP(&endDate, "end", "e", "YYYYMMDD end date; defaults to today (local time)")
	pflag.Parse()

	// set up start/end times
	if startDate.IsZero() {
		// if startDate is not set, default to endDate
		startDate = endDate
	}

	var startTime, endTime time.Time
	if st, et := startDate.AsTime(), endDate.AsTime(); et.Equal(st) || et.After(st) {
		startTime = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, time.UTC)
		endTime = time.Date(et.Year(), et.Month(), et.Day(), 23, 59, 59, 0, time.UTC)
	} else {
		fmt.Fprint(os.Stderr, "--start must be before --end\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "startTime: %s   endTime: %s\n", startTime.String(), endTime.String())
}
