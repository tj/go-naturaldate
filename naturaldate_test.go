package naturaldate

import (
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
)

// TODO: CI tests
// TODO: noon / midnight
// TODO: last week / next week

// base time.
var base = time.Unix(1574687238, 0).UTC()

// pastCases are test cases for the past direction.
var pastCases = []struct {
	Input  string
	Output string
}{
	{`now`, `2019-11-25 13:07:18 +0000 UTC`},

	// minutes
	{`1 minute`, `2019-11-25 13:06:18 +0000 UTC`},
	{`one minute`, `2019-11-25 13:06:18 +0000 UTC`},
	{`1 minute ago`, `2019-11-25 13:06:18 +0000 UTC`},
	{`5 minutes ago`, `2019-11-25 13:02:18 +0000 UTC`},
	{`five minutes ago`, `2019-11-25 13:02:18 +0000 UTC`},
	{`   5    minutes  ago   `, `2019-11-25 13:02:18 +0000 UTC`},

	// hours
	{`1 hour`, `2019-11-25 12:07:18 +0000 UTC`},
	{`1 hour ago`, `2019-11-25 12:07:18 +0000 UTC`},
	{`6 hours ago`, `2019-11-25 07:07:18 +0000 UTC`},

	// days
	{`1 day`, `2019-11-24 00:00:00 +0000 UTC`},
	{`1 day ago`, `2019-11-24 00:00:00 +0000 UTC`},
	{`3 days ago`, `2019-11-22 00:00:00 +0000 UTC`},

	// weeks
	{`1 week`, `2019-11-18 00:00:00 +0000 UTC`},
	{`1 week ago`, `2019-11-18 00:00:00 +0000 UTC`},
	{`2 weeks ago`, `2019-11-11 00:00:00 +0000 UTC`},

	// months
	{`1 month ago`, `2019-10-25 00:00:00 +0000 UTC`},
	{`1 month ago at 9:30am`, `2019-10-25 09:30:00 +0000 UTC`},
	{`2 months ago`, `2019-09-25 00:00:00 +0000 UTC`},
	{`12 months ago`, `2018-11-25 00:00:00 +0000 UTC`},
	{`1 month from now`, `2019-12-25 00:00:00 +0000 UTC`},
	{`2 months from now`, `2020-01-25 00:00:00 +0000 UTC`},
	{`12 months from now at 6am`, `2020-11-25 06:00:00 +0000 UTC`},

	// years
	{`last year`, `2018-01-01 00:00:00 +0000 UTC`},
	{`next year`, `2020-01-01 00:00:00 +0000 UTC`},
	{`one year ago`, `2018-11-25 00:00:00 +0000 UTC`},
	{`one year from now`, `2020-11-25 00:00:00 +0000 UTC`},
	{`two years ago`, `2017-11-25 00:00:00 +0000 UTC`},
	{`2 years ago`, `2017-11-25 00:00:00 +0000 UTC`},

	// today
	{`today`, `2019-11-25 00:00:00 +0000 UTC`},
	{`today at 10am`, `2019-11-25 10:00:00 +0000 UTC`},

	// yesterday
	{`yesterday`, `2019-11-24 00:00:00 +0000 UTC`},
	{`yesterday 10am`, `2019-11-24 10:00:00 +0000 UTC`},
	{`yesterday at 10am`, `2019-11-24 10:00:00 +0000 UTC`},
	{`yesterday at 10:15am`, `2019-11-24 10:15:00 +0000 UTC`},

	// past weekdays
	{`sunday`, `2019-11-24 00:00:00 +0000 UTC`},
	{`monday`, `2019-11-18 00:00:00 +0000 UTC`},
	{`tuesday`, `2019-11-19 00:00:00 +0000 UTC`},
	{`wednesday`, `2019-11-20 00:00:00 +0000 UTC`},
	{`thursday`, `2019-11-21 00:00:00 +0000 UTC`},
	{`friday`, `2019-11-22 00:00:00 +0000 UTC`},
	{`saturday`, `2019-11-23 00:00:00 +0000 UTC`},

	{`last sunday`, `2019-11-24 00:00:00 +0000 UTC`},
	{`past sunday`, `2019-11-24 00:00:00 +0000 UTC`},
	{`last monday`, `2019-11-18 00:00:00 +0000 UTC`},
	{`last tuesday`, `2019-11-19 00:00:00 +0000 UTC`},
	{`last wednesday`, `2019-11-20 00:00:00 +0000 UTC`},
	{`last thursday`, `2019-11-21 00:00:00 +0000 UTC`},
	{`last friday`, `2019-11-22 00:00:00 +0000 UTC`},
	{`last saturday`, `2019-11-23 00:00:00 +0000 UTC`},

	// future weekdays
	{`next tuesday`, `2019-11-26 00:00:00 +0000 UTC`},
	{`next wednesday`, `2019-11-27 00:00:00 +0000 UTC`},
	{`next thursday`, `2019-11-28 00:00:00 +0000 UTC`},
	{`next friday`, `2019-11-29 00:00:00 +0000 UTC`},
	{`next saturday`, `2019-11-30 00:00:00 +0000 UTC`},
	{`next sunday`, `2019-12-01 00:00:00 +0000 UTC`},
	{`next monday`, `2019-12-02 00:00:00 +0000 UTC`},

	// months
	{`next january`, `2020-01-01 00:00:00 +0000 UTC`},
	{`last january`, `2019-01-01 00:00:00 +0000 UTC`},
	{`january`, `2019-01-01 00:00:00 +0000 UTC`},
	{`february`, `2019-02-01 00:00:00 +0000 UTC`},
	{`march`, `2019-03-01 00:00:00 +0000 UTC`},
	{`april`, `2019-04-01 00:00:00 +0000 UTC`},
	{`may`, `2019-05-01 00:00:00 +0000 UTC`},
	{`june`, `2019-06-01 00:00:00 +0000 UTC`},
	{`july`, `2019-07-01 00:00:00 +0000 UTC`},
	{`august`, `2019-08-01 00:00:00 +0000 UTC`},
	{`september`, `2019-09-01 00:00:00 +0000 UTC`},
	{`october`, `2019-10-01 00:00:00 +0000 UTC`},
	{`november`, `2018-11-01 00:00:00 +0000 UTC`},

	// ordinal dates
	{`december 1`, `2018-12-01 00:00:00 +0000 UTC`},
	{`december 15`, `2018-12-15 00:00:00 +0000 UTC`},
	{`december 1st`, `2018-12-01 00:00:00 +0000 UTC`},
	{`december 2nd`, `2018-12-02 00:00:00 +0000 UTC`},
	{`december 3rd`, `2018-12-03 00:00:00 +0000 UTC`},
	{`december 4th`, `2018-12-04 00:00:00 +0000 UTC`},
	{`december 15th`, `2018-12-15 00:00:00 +0000 UTC`},
	{`december 23rd`, `2018-12-23 00:00:00 +0000 UTC`},
	{`december 23rd 5pm`, `2018-12-23 17:00:00 +0000 UTC`},
	{`december 23rd at 5pm`, `2018-12-23 17:00:00 +0000 UTC`},
	{`december 23rd at 5:25pm`, `2018-12-23 17:25:00 +0000 UTC`},

	// 12-hour clock
	{`10am`, `2019-11-25 10:00:00 +0000 UTC`},
	{`10 am`, `2019-11-25 10:00:00 +0000 UTC`},
	{`5pm`, `2019-11-25 17:00:00 +0000 UTC`},
	{`10:25am`, `2019-11-25 10:25:00 +0000 UTC`},
	{`1:05pm`, `2019-11-25 13:05:00 +0000 UTC`},
	{`10:25:10am`, `2019-11-25 10:25:10 +0000 UTC`},
	{`1:05:10pm`, `2019-11-25 13:05:10 +0000 UTC`},

	// 24-hour clock
	{`10`, `2019-11-25 10:00:00 +0000 UTC`},
	{`10:25`, `2019-11-25 10:25:00 +0000 UTC`},
	{`10:25:30`, `2019-11-25 10:25:30 +0000 UTC`},
	{`17`, `2019-11-25 17:00:00 +0000 UTC`},
	{`17:25:30`, `2019-11-25 17:25:30 +0000 UTC`},

	// case sensitivity
	{`December 23rd AT 5:25 PM`, `2018-12-23 17:25:00 +0000 UTC`},
	{`next December 23rd AT 5:25 PM`, `2019-12-23 17:25:00 +0000 UTC`},

	// errors
	{`10:am`, "\nparse error near PegText (line 1 symbol 1 - line 1 symbol 3):\n\"10\"\n"},
	{`today at`, "\nparse error near AT (line 1 symbol 7 - line 1 symbol 9):\n\"at\"\n"},
	{`yesterday at`, "\nparse error near AT (line 1 symbol 11 - line 1 symbol 13):\n\"at\"\n"},
}

// futureCases are test cases for the future direction.
var futureCases = []struct {
	Input  string
	Output string
}{
	{`now`, `2019-11-25 13:07:18 +0000 UTC`},
	{`1 minute`, `2019-11-25 13:08:18 +0000 UTC`},
	{`1 hour`, `2019-11-25 14:07:18 +0000 UTC`},
	{`1 day`, `2019-11-26 00:00:00 +0000 UTC`},
	{`1 week`, `2019-12-02 00:00:00 +0000 UTC`},
	{`previous tuesday`, `2019-11-19 00:00:00 +0000 UTC`},
	{`tuesday`, `2019-11-26 00:00:00 +0000 UTC`},
	{`wednesday`, `2019-11-27 00:00:00 +0000 UTC`},
	{`thursday`, `2019-11-28 00:00:00 +0000 UTC`},
	{`friday`, `2019-11-29 00:00:00 +0000 UTC`},
	{`saturday`, `2019-11-30 00:00:00 +0000 UTC`},
	{`sunday`, `2019-12-01 00:00:00 +0000 UTC`},
	{`monday`, `2019-12-02 00:00:00 +0000 UTC`},
	{`last january`, `2019-01-01 00:00:00 +0000 UTC`},
	{`january`, `2020-01-01 00:00:00 +0000 UTC`},
	{`next january`, `2020-01-01 00:00:00 +0000 UTC`},
}

// Test parsing with past direction.
func TestParse_past(t *testing.T) {
	for _, c := range pastCases {
		t.Run(c.Input, func(t *testing.T) {
			v, err := Parse(c.Input, base)
			if err != nil {
				assert.Equal(t, c.Output, err.Error())
				return
			}
			assert.Equal(t, c.Output, v.UTC().String())
		})
	}
}

// Test parsing with future direction.
func TestParse_future(t *testing.T) {
	for _, c := range futureCases {
		t.Run(c.Input, func(t *testing.T) {
			v, err := Parse(c.Input, base, WithDirection(Future))
			if err != nil {
				assert.Equal(t, c.Output, err.Error())
				return
			}
			assert.Equal(t, c.Output, v.UTC().String())
		})
	}
}

// Benchmark parsing.
func BenchmarkParse(b *testing.B) {
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		_, err := Parse(`december 23rd at 5:25pm`, base)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}
}
