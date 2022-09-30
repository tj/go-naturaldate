# Go Natural Date

[![CircleCI](https://circleci.com/gh/ijt/go-naturaldate/tree/master.svg?style=shield)](https://circleci.com/gh/ijt/go-naturaldate/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ijt/go-naturaldate)](https://goreportcard.com/report/github.com/ijt/go-naturaldate)
[![GoDoc](https://godoc.org/github.com/tj/go-naturaldate?status.svg)](https://godoc.org/github.com/tj/go-naturaldate)
![](https://img.shields.io/badge/license-MIT-blue.svg)

Natural date time parsing for Go. This package was designed for parsing human-friendly relative date/time ranges in [Apex Logs](https://apex.sh/logs/)' command-line log search.

## Examples

Here are some examples of the types of expressions currently supported:

- now
- today
- yesterday
- 5 minutes ago
- three days ago
- last month
- next month
- one year from now
- yesterday at 10am
- last sunday at 5:30pm
- next sunday at 22:45
- next January
- last February
- next December 25th at 7:30am
- November 3rd, 1986 at 4:30pm
- See the [tests](./naturaldate_test.go) for more examples
