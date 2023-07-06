# ymdflag #

[![Units tests](https://github.com/neomantra/ymdflag/actions/workflows/unit_tests.yaml/badge.svg)](https://github.com/neomantra/ymdflag/actions/workflows/unit_tests.yaml)
[![Coverage Status](https://coveralls.io/repos/neomantra/ymdflag/badge.svg?branch=main&service=github)](https://coveralls.io/github/neomantra/ymdflag?branch=main)
[![golangci-lint](https://github.com/neomantra/ymdflag/actions/workflows/golangci-lint.yaml/badge.svg)](https://github.com/neomantra/ymdflag/actions/workflows/golangci-lint.yaml)
[![CodeQL](https://github.com/neomantra/ymdflag/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/neomantra/ymdflag/actions/workflows/codeql-analysis.yml)
[![Go ReportCard](https://goreportcard.com/badge/neomantra/ymdflag)](http://goreportcard.com/report/neomantra/ymdflag)
[![Go Reference](https://pkg.go.dev/badge/github.com/neomantra/ymdflag.svg)](https://pkg.go.dev/github.com/neomantra/ymdflag)

[`YMDFlag`](https://github.com/neomantra/ymdflag) implements a Golang [`flag.Value`](https://pkg.go.dev/flag#Value) interface for `YYYYMMDD`-specified dates.   This facilitiates command-line argument handling of date parameters such  `-start-date=20210101`.

## Documentation ##

```
go get github.com/neomantra/ymdflag
```

https://pkg.go.dev/github.com/neomantra/ymdflag

----

## Credits and License

Copyright (c) 2022-2023 Neomantra BV.  Authored by Evan Wies.

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).
