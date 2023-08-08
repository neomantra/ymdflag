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
	locationName := "Antarctica/Syowa"
	location, err := time.LoadLocation(locationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query location: %s err: %s\n", locationName, err.Error())
	}

	ymdWithLocation := ymdflag.NewYMDFlagWithLocation(location)

	pflag.VarP(&ymdWithLocation, "date", "d", fmt.Sprintf("YYYYMMDD date in %s", locationName))
	pflag.Parse()

	fmt.Fprintf(os.Stdout, "location: %s   date: %s  time:   timeUTC: %s\n",
		locationName,
		ymdWithLocation.AsTime().String(),
		ymdWithLocation.AsTime().In(time.UTC).String())
}
