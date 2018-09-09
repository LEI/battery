# battery

Based on [distatus/battery](https://github.com/distatus/battery).

[![GoDoc](https://godoc.org/github.com/LEI/battery?status.svg)](https://godoc.org/github.com/LEI/battery)
[![Travis](https://travis-ci.org/LEI/battery.svg?branch=master)](https://travis-ci.org/LEI/battery)

<!--
[![Codecov](https://codecov.io/gh/LEI/battery/branch/master/graph/badge.svg)](https://codecov.io/gh/LEI/battery)
[![Go Report Card](https://goreportcard.com/badge/github.com/LEI/battery)](https://goreportcard.com/report/github.com/LEI/battery)
-->

## Installation

```console
$ go get -u github.com/LEI/battery
$ cd $GOPATH/src/github.com/LEI/battery
$ dep ensure -vendor-only
$ go install github.com/LEI/battery
```

## Usage

```console
$ battery -h
battery [flags] [format]
    -c, --color   Enable color output
    -s, --spark   Enable sparkline bar
    -t, --tmux    Enable tmux status bar colors
```

Default format:

```console
$ battery '{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration ""}}, {{end}}{{.Duration}}'
BAT0: Full, 100%
```

Example:

```console
$ battery --spark '{{.Bar}}'
█ 100%
```
