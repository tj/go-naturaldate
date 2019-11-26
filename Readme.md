# Go Natural Date

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
- sunday at 22:45
- next January
- last February
- December 25th at 7:30am
- 10am
- 10:05pm
- 10:05:22pm

## Direction

A default direction can be applied using `WithDirection()` for ambiguous expressions such as `sunday`, or `september`. By default `naturaldate.Past` is used, so they will be equivalent to `last sunday` and `last september`.

## Notes

This package is does not have the goal of parsing all possible date formats and timestamps, I recommend you attempt parsing those first and then falling back on this package to provide "natural language" support. 

---

[![GoDoc](https://godoc.org/github.com/tj/go-naturaldate?status.svg)](https://godoc.org/github.com/tj/go-naturaldate)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

## Sponsors

This project is [sponsored](https://github.com/sponsors/tj) by [CTO.ai](https://cto.ai/), making it easy for development teams to create and share workflow automations without leaving the command line. 

[![](https://apex-software.imgix.net/github/sponsors/cto.png)](https://cto.ai/)