# Go Natural Date (ijt fork)

[![CircleCI](https://circleci.com/gh/ijt/go-naturaldate/tree/master.svg?style=shield)](https://circleci.com/gh/ijt/go-naturaldate/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ijt/go-naturaldate)](https://goreportcard.com/report/github.com/ijt/go-naturaldate)
[![GoDoc](https://godoc.org/github.com/tj/go-naturaldate?status.svg)](https://godoc.org/github.com/tj/go-naturaldate)
![](https://img.shields.io/badge/license-MIT-blue.svg)

Natural date time parsing for Go. This package was forked from github.com/tj/go-naturaldate with the following goals:

1. Support use within parser combinator packages such as github.com/ijt/goparsify. As part of this, parsing must be more strict: the date must come at the beginning of the string although there can be additional text after it. Also, the substring matched by the parser must be made available.
2. Minimize the amount of ambiguity accepted. Instead of assuming that "Tuesday" means last Tuesday or next Tuesday as in github.com/tj/go-naturaldate, just require additional context so it's clear.

Both of those goals appear to have been met.

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
