# battery

Based on [distatus/battery](https://github.com/distatus/battery)

## Installation

    $ go get -u github.com/LEI/battery
    $ go install $_

## Usage

    $ battery -h
    battery [flags] [format]
      -c, --color   Enable color output
      -s, --spark   Enable sparkline bar
      -t, --tmux    Enable tmux status bar colors

Default format

    $ battery '{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration ""}}, {{end}}{{.Duration}}'

Example

    $ battery --spark '{{.Bar}}'
