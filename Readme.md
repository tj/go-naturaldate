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
- next December 25th at 7:30am UTC-7
- November 3rd, 1986 at 4:30pm
- january 2017
- january, 2017
- oct 7, 1970
- oct 7 1970
- 7 oct 1970
- 7 oct, 1970
- September 17, 2012 UTC+7
- September 17, 2012
- 03 February 2013
- 2 July 2013
- 2014/3/31
- 2014/3/31 UTC
- 2014/3/31 UTC+1
- 2014/03/31
- 2014/03/31 UTC-1
- 2014-04-26
- 2014-4-26
- 2014-4-6
- 31/3/2014 UTC-8
- 31-3-2014 UTC-8
- 31/3/2014
- 31-3-2014
- January
- See the [tests](./naturaldate_test.go) for more examples
