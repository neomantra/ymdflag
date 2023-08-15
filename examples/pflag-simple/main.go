// Copyright (c) 2023 Neomantra BV

package main

import (
	"github.com/neomantra/ymdflag"
	"github.com/spf13/pflag"
)

func main() {
	var ymd ymdflag.YMDFlag
	pflag.VarP(&ymd, "date", "d", "YYYYMMDD date; defaults to today in local time")
	pflag.Parse()
	println("time of date:", ymd.AsTime().String())
}
