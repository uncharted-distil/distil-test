# distil-test

[![Go Report Card](https://goreportcard.com/badge/github.com/uncharted-distil/distil-test)](https://goreportcard.com/report/github.com/uncharted-distil/distil-test)
[![GolangCI](https://golangci.com/badges/github.com/uncharted-distil/distil.svg)](https://golangci.com/r/github.com/uncharted-distil/distil)
[![GolangCI](https://golangci.com/badges/github.com/uncharted-distil/distil-test.svg)](https://golangci.com/r/github.com/uncharted-distil/distil-test)

A headless client that interfaces with a running distil server to perform integration tests.  While it can be run stand-alone, its main purpose is integration into the D3M TA3-TA2 CI environment found at https://dash.datadrivendiscovery.org/ta3ta2.  The CI project itself is at https://gitlab.datadrivendiscovery.org/d3m/ta3ta2-ci/.  Both are non-public and require D3M program access.

To build the test container:
```shell
./build.sh
```

To run:
```
./run.sh
```
