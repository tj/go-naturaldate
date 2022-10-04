package naturaldate

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"
)

func dateAtTime(dateFrom time.Time, hour int, min int, sec int) time.Time {
	t := dateFrom
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, sec, 0, t.Location())
}

// Test parsing with inputs that are expected to result in errors.
func TestParse_bad(t *testing.T) {
	var badCases = []struct {
		input string
	}{
		{``},
		{`a`},
		{`not a date or a time`},
		{`right now`},
		{`  right  now  `},
		{`Message me in 2 minutes`},
		{`Message me in 2 minutes from now`},
		{`Remind me in 1 hour`},
		{`Remind me in 1 hour from now`},
		{`Remind me in 1 hour and 3 minutes from now`},
		{`Remind me in an hour`},
		{`Remind me in an hour from now`},
		{`Remind me one day from now`},
		{`Remind me in a day`},
		{`Remind me in one day`},
		{`Remind me in one day from now`},
		{`Message me in a week`},
		{`Message me in one week`},
		{`Message me in one week from now`},
		{`Message me in two weeks from now`},
		{`Message me two weeks from now`},
		{`Message me in two weeks`},
		{`Remind me in 12 months from now at 6am`},
		{`Remind me in a month`},
		{`Remind me in 2 months`},
		{`Remind me in a month from now`},
		{`Remind me in 2 months from now`},
		{`Remind me in one year from now`},
		{`Remind me in a year`},
		{`Remind me in a year from now`},
		{`Restart the server in 2 days from now`},
		{`Remind me on the 5th of next month`},
		{`Remind me on the 5th of next month at 7am`},
		{`Remind me at 7am on the 5th of next month`},
		{`Remind me in one month from now`},
		{`Remind me in one month from now at 7am`},
		{`Remind me on the December 25th at 7am`},
		{`Remind me at 7am on December 25th`},
		{`Remind me on the 25th of December at 7am`},
		{`Check logs in the past 5 minutes`},

		// "1 minute" is a duration, not a time.
		{`1 minute`},

		// "one minute" is also a duration.
		{`one minute`},

		// "1 hour" is also a duration.
		{`1 hour`},

		// "1 day" is also a duration.
		{`1 day`},

		// "1 week" is also a duration.
		{`1 week`},

		// "1 month" is also a duration.
		{`1 month`},

		// "next 2 months" is a date range, not a time or a date.
		{`next 2 months`},

		// Ambiguous weekday inputs:
		// These are ambiguous because they don't tell whether it's the
		// previous, next or, in some cases, current instance of the weekday.
		{`sunday`},
		{`monday`},
		{`tuesday`},
		{`wednesday`},
		{`thursday`},
		{`friday`},
		{`saturday`},

		// Ambiguous month inputs:
		// These are ambiguous because they don't include the year.
		{`january`},
		{`february`},
		{`march`},
		{`april`},
		{`may`},
		{`june`},
		{`july`},
		{`august`},
		{`september`},
		{`october`},
		{`november`},

		// Ambiguous ordinal dates:
		// These are ambiguous because they don't include the year.
		{`november 15th`},
		{`december 1st`},
		{`december 2nd`},
		{`december 3rd`},
		{`december 4th`},
		{`december 15th`},
		{`december 23rd`},
		{`december 23rd 5pm`},
		{`december 23rd at 5pm`},
		{`december 23rd at 5:25pm`},
		{`December 23rd AT 5:25 PM`},
		{`December 25th at 7am`},
		{`7am on December 25th`},
		{`25th of December at 7am`},

		// Ambiguous 12-hour times:
		// These are ambiguous because they don't include the date.
		{`10am`},
		{`10 am`},
		{`5pm`},
		{`10:25am`},
		{`1:05pm`},
		{`10:25:10am`},
		{`1:05:10pm`},

		// Ambiguous 24-hour times:
		// These are ambiguous because they don't include the date.
		{`10`},
		{`10:25`},
		{`10:25:30`},
		{`17`},
		{`17:25:30`},

		// Goofy input:
		{`10:am`},
	}
	for _, c := range badCases {
		t.Run(c.input, func(t *testing.T) {
			now := time.Time{}
			v, _, err := Parse(c.input, now)
			if err == nil {
				t.Errorf("err is nil, result is %v", v)
			}
		})
	}
}

// Test parsing on cases that are expected to parse successfully.
func TestParse_goodTimes(t *testing.T) {
	now := time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local)
	var cases = []struct {
		Input    string
		WantTime time.Time
	}{
		// now
		{`now`, now},

		// minutes
		{`a minute from now`, now.Add(time.Minute)},
		{`a minute ago`, now.Add(-time.Minute)},
		{`next minute`, now.Add(time.Minute)},
		{`last minute`, now.Add(-time.Minute)},
		{`1 minute ago`, now.Add(-time.Minute)},
		{`5 minutes ago`, now.Add(-5 * time.Minute)},
		{`five minutes ago`, now.Add(-5 * time.Minute)},
		{`   5    minutes  ago   `, now.Add(-5 * time.Minute)},
		{`2 minutes from now`, now.Add(2 * time.Minute)},
		{`two minutes from now`, now.Add(2 * time.Minute)},

		// hours
		{`an hour from now`, now.Add(time.Hour)},
		{`an hour ago`, now.Add(-time.Hour)},
		{`last hour`, now.Add(-time.Hour)},
		{`next hour`, now.Add(time.Hour)},
		{`1 hour ago`, now.Add(-time.Hour)},
		{`6 hours ago`, now.Add(-6 * time.Hour)},
		{`1 hour from now`, now.Add(time.Hour)},

		// dates with times
		{`3 days ago at 11:25am`, dateAtTime(now.Add(-3*24*time.Hour), 11, 25, 0)},
		{`2 weeks ago at 8am`, dateAtTime(now.Add(-2*7*24*time.Hour), 8, 0, 0)},
		{`today at 10am`, dateAtTime(now, 10, 0, 0)},
		{`yesterday 10am`, dateAtTime(now.AddDate(0, 0, -1), 10, 0, 0)},
		{`yesterday at 10am`, dateAtTime(now.AddDate(0, 0, -1), 10, 0, 0)},
		{`yesterday at 10:15am`, dateAtTime(now.AddDate(0, 0, -1), 10, 15, 0)},
		{`tomorrow 10am`, dateAtTime(now.AddDate(0, 0, 1), 10, 0, 0)},
		{`10am tomorrow`, dateAtTime(now.AddDate(0, 0, 1), 10, 0, 0)},
		{`tomorrow at 10am`, dateAtTime(now.AddDate(0, 0, 1), 10, 0, 0)},
		{`tomorrow at 10:15am`, dateAtTime(now.AddDate(0, 0, 1), 10, 15, 0)},
		{"next December 25th at 7:30am UTC-7", timeInLocation(nextMonthDayTime(now, time.December, 25, 7, 30, 0), fixedZone(-7))},
		{`next December 23rd AT 5:25 PM`, nextMonthDayTime(now, time.December, 23, 12+5, 25, 0)},
		{`last sunday at 5:30pm`, dateAtTime(prevWeekday(now, time.Sunday), 12+5, 30, 0)},
		{`next sunday at 22:45`, dateAtTime(nextWeekday(now, time.Sunday), 22, 45, 0)},
		{`next sunday at 22:45`, dateAtTime(nextWeekday(now, time.Sunday), 22, 45, 0)},
		{`November 3rd, 1986 at 4:30pm`, time.Date(1986, 11, 3, 12+4, 30, 0, 0, now.Location())},
		{"September 17, 2012 at 10:09am UTC", time.Date(2012, 9, 17, 10, 9, 0, 0, time.UTC)},
		{"September 17, 2012 at 10:09am UTC-8", time.Date(2012, 9, 17, 10, 9, 0, 0, fixedZone(-8))},
		{"September 17, 2012 at 10:09am UTC+8", time.Date(2012, 9, 17, 10, 9, 0, 0, fixedZone(8))},
		{"September 17, 2012, 10:11:09", time.Date(2012, 9, 17, 10, 11, 9, 0, now.Location())},
		{"September 17, 2012, 10:11", time.Date(2012, 9, 17, 10, 11, 0, 0, now.Location())},
		{"September 17, 2012 10:11", time.Date(2012, 9, 17, 10, 11, 0, 0, now.Location())},
		{"September 17 2012 10:11", time.Date(2012, 9, 17, 10, 11, 0, 0, now.Location())},
		{"September 17 2012 at 10:11", time.Date(2012, 9, 17, 10, 11, 0, 0, now.Location())},

		// formats from the Go time package:
		// ANSIC
		{"Mon Jan _2 15:04:05 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, now.Location())},
		// UnixDate
		{"Mon Jan _2 15:04:05 MST 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, location("MST"))},
		// RubyDate
		{"Mon Jan 02 15:04:05 -0700 2006", time.Date(2006, 1, 2, 15, 4, 5, 0, fixedZone(-7))},
		// RFC1123
		{"Mon, 02 Jan 2006 15:04:05 MST", time.Date(2006, 1, 2, 15, 4, 5, 0, location("MST"))},
		// RFC1123Z
		{"Mon, 02 Jan 2006 15:04:05 -0700", time.Date(2006, 1, 2, 15, 4, 5, 0, fixedZone(-7))},
		// RFC3339
		{"2006-01-02T15:04:05Z07:00", time.Date(2006, 1, 2, 15, 4, 5, 0, fixedZone(7))},
		// RFC3339Nano
		{"2006-01-02T15:04:05.999999999Z07:00", time.Date(2006, 1, 2, 15, 4, 5, 999999999, fixedZone(7))},
	}

	for _, c := range cases {
		t.Run(c.Input, func(t *testing.T) {
			v, match, err := Parse(c.Input, now)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, c.WantTime, v)
			assert.Equal(t, strings.ToLower(c.Input), match)
		})
	}
}

func TestParse_goodDays(t *testing.T) {
	now := time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local)
	var cases = []struct {
		Input    string
		WantTime time.Time
	}{
		// days
		{`one day ago`, now.Add(-24 * time.Hour)},
		{`1 day ago`, now.Add(-24 * time.Hour)},
		{`3 days ago`, now.Add(-3 * 24 * time.Hour)},
		{`three days ago`, now.Add(-3 * 24 * time.Hour)},
		{`1 day from now`, now.Add(24 * time.Hour)},

		// weeks
		{`1 week ago`, now.Add(-7 * 24 * time.Hour)},
		{`2 weeks ago`, now.Add(-2 * 7 * 24 * time.Hour)},
		{`next week`, now.Add(7 * 24 * time.Hour)},
		{`a week from now`, now.Add(7 * 24 * time.Hour)},
		{`a week from today`, now.Add(7 * 24 * time.Hour)},

		// months
		{`a month ago`, now.AddDate(0, -1, 0)},
		{`1 month ago`, now.AddDate(0, -1, 0)},
		{`last month`, now.AddDate(0, -1, 0)},
		{`next month`, now.AddDate(0, 1, 0)},
		{`2 months ago`, now.AddDate(0, -2, 0)},
		{`12 months ago`, now.AddDate(0, -12, 0)},
		{`a month from now`, now.AddDate(0, 1, 0)},
		{`1 month from now`, now.AddDate(0, 1, 0)},
		{`2 months from now`, now.AddDate(0, 2, 0)},
		{`last january`, prevMonth(now, time.January)},
		{`next january`, nextMonth(now, time.January)},

		// years
		{`last year`, truncateYear(now.AddDate(-1, 0, 0))},
		{`next year`, truncateYear(now.AddDate(1, 0, 0))},
		{`one year ago`, truncateYear(now.AddDate(-1, 0, 0))},
		{`one year from now`, truncateYear(now.AddDate(1, 0, 0))},
		{`one year from today`, truncateYear(now.AddDate(1, 0, 0))},
		{`two years ago`, truncateYear(now.AddDate(-2, 0, 0))},
		{`2 years ago`, truncateYear(now.AddDate(-2, 0, 0))},

		// today
		{`today`, now},

		// yesterday
		{`yesterday`, now.AddDate(0, 0, -1)},

		// tomorrow
		{`tomorrow`, now.AddDate(0, 0, 1)},

		// past weekdays
		{`last sunday`, prevWeekday(now, time.Sunday)},
		{`past sunday`, prevWeekday(now, time.Sunday)},
		{`last monday`, prevWeekday(now, time.Monday)},
		{`last tuesday`, prevWeekday(now, time.Tuesday)},
		{`previous tuesday`, prevWeekday(now, time.Tuesday)},
		{`last wednesday`, prevWeekday(now, time.Wednesday)},
		{`last thursday`, prevWeekday(now, time.Thursday)},
		{`last friday`, prevWeekday(now, time.Friday)},
		{`last saturday`, prevWeekday(now, time.Saturday)},

		// future weekdays
		{`next tuesday`, nextWeekday(now, time.Tuesday)},
		{`next wednesday`, nextWeekday(now, time.Wednesday)},
		{`next thursday`, nextWeekday(now, time.Thursday)},
		{`next friday`, nextWeekday(now, time.Friday)},
		{`next saturday`, nextWeekday(now, time.Saturday)},
		{`next sunday`, nextWeekday(now, time.Sunday)},
		{`next monday`, nextWeekday(now, time.Monday)},

		// months
		{`last january`, prevMonth(now, time.January)},
		{`next january`, nextMonth(now, time.January)},

		// absolute dates
		{"january 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, now.Location())},
		{"january, 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, now.Location())},
		{"april 3 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, now.Location())},
		{"april 3, 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, now.Location())},
		{"oct 7, 1970", time.Date(1970, 10, 7, 0, 0, 0, 0, now.Location())},
		{"oct 7 1970", time.Date(1970, 10, 7, 0, 0, 0, 0, now.Location())},
		{"oct. 7, 1970", time.Date(1970, 10, 7, 0, 0, 0, 0, now.Location())},
		{"September 17, 2012 UTC+7", time.Date(2012, 9, 17, 10, 9, 0, 0, fixedZone(7))},
		{"September 17, 2012", time.Date(2012, 9, 17, 10, 9, 0, 0, now.Location())},
		{"7 oct 1970", time.Date(1970, 10, 7, 0, 0, 0, 0, now.Location())},
		{"7 oct, 1970", time.Date(1970, 10, 7, 0, 0, 0, 0, now.Location())},
		{"03 February 2013", time.Date(2013, 2, 3, 0, 0, 0, 0, now.Location())},
		{"2 July 2013", time.Date(2013, 7, 2, 0, 0, 0, 0, now.Location())},
		// yyyy/mm/dd
		{"2014/3/31", time.Date(2014, 3, 31, 0, 0, 0, 0, now.Location())},
		{"2014/3/31 UTC", time.Date(2014, 3, 31, 0, 0, 0, 0, location("UTC"))},
		{"2014/3/31 UTC+1", time.Date(2014, 3, 31, 0, 0, 0, 0, fixedZone(1))},
		{"2014/03/31", time.Date(2014, 3, 31, 0, 0, 0, 0, now.Location())},
		{"2014/03/31 UTC-1", time.Date(2014, 3, 31, 0, 0, 0, 0, fixedZone(-1))},
		{"2014-04-26", time.Date(2014, 4, 26, 0, 0, 0, 0, now.Location())},
		{"2014-4-26", time.Date(2014, 4, 26, 0, 0, 0, 0, now.Location())},
		{"2014-4-6", time.Date(2014, 4, 6, 0, 0, 0, 0, now.Location())},
		{"31/3/2014 UTC-8", time.Date(2014, 3, 31, 0, 0, 0, 0, fixedZone(-8))},
		{"31-3-2014 UTC-8", time.Date(2014, 3, 31, 0, 0, 0, 0, fixedZone(-8))},
		{"31/3/2014", time.Date(2014, 3, 31, 0, 0, 0, 0, now.Location())},
		{"31-3-2014", time.Date(2014, 3, 31, 0, 0, 0, 0, now.Location())},
	}

	for _, c := range cases {
		t.Run(c.Input, func(t *testing.T) {
			v, match, err := Parse(c.Input, now)
			if err != nil {
				t.Fatal(err)
			}
			want := truncateDay(c.WantTime)
			assert.Equal(t, want, v)
			assert.Equal(t, strings.ToLower(c.Input), match)
		})
	}
}

func fixedZone(offset int) *time.Location {
	name := fmt.Sprintf("UTC+%d", offset)
	if offset < 0 {
		name = fmt.Sprintf("UTC-%d", -offset)
	}
	return time.FixedZone(name, offset)
}

func location(locStr string) *time.Location {
	l, err := time.LoadLocation(locStr)
	if err != nil {
		panic(fmt.Sprintf("loading location %q: %v", locStr, err))
	}
	return l
}

func timeInLocation(t time.Time, l *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), l)
}

func TestParse_withStuffAtEnd(t *testing.T) {
	now := time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local)
	var cases = []struct {
		Input     string
		WantMatch string
		WantTime  time.Time
	}{
		{`last year I moved to a new location`, "last year ", truncateYear(now.AddDate(-1, 0, 0))},
		{`today I'm going out of town`, "today ", truncateDay(now)},
		{`next Monday is an important meeting`, "next monday ", truncateDay(nextWeekday(now, time.Monday))},
	}
	for _, c := range cases {
		t.Run(c.Input, func(t *testing.T) {
			v, match, err := Parse(c.Input, now)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, c.WantTime, v)
			assert.Equal(t, c.WantMatch, match)
		})
	}
}

// Benchmark parsing.
func BenchmarkParse(b *testing.B) {
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(`december 23rd 2022 at 5:25pm`, time.Time{})
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}
}
