# battery

Based on [distatus/battery](https://github.com/distatus/battery)

## Installation

    $ go get -u github.com/LEI/battery
    $ go install $_

## Usage

Default format: `{{.Id}}: {{.State}}, {{.Percent}}%{{if ne .Duration ""}}, {{end}}{{.Duration}}`

    $ battery [flags] [format]
    # -c, --color
    # -t, --tmux
