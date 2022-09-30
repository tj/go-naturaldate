package naturaldate

import (
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
)

// base time.
var base = time.Unix(1574687238, 0).UTC()

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

		// ambiguous weekday inputs
		{`sunday`},
		{`monday`},
		{`tuesday`},
		{`wednesday`},
		{`thursday`},
		{`friday`},
		{`saturday`},

		// ambiguous month inputs
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

		// ambiguous ordinal dates
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

		// ambiguous: 12-hour clock
		{`10am`},
		{`10 am`},
		{`5pm`},
		{`10:25am`},
		{`1:05pm`},
		{`10:25:10am`},
		{`1:05:10pm`},

		// ambiguous: 24-hour clock
		{`10`},
		{`10:25`},
		{`10:25:30`},
		{`17`},
		{`17:25:30`},

		// goofy input
		{`10:am`},
	}
	for _, c := range badCases {
		t.Run(c.input, func(t *testing.T) {
			_, err := Parse(c.input, base)
			if err == nil {
				t.Errorf("err is %v, want nil", err)
			}

			_, err = Parse(c.input, base, WithDirection(Future))
			if err == nil {
				t.Errorf("future: err is %v, want nil", err)
			}
		})
	}
}

// Test parsing on cases that are expected to parse successfully.
func TestParse_good(test *testing.T) {
	var baseTimes = []time.Time{
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local),
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.UTC),
	}

	for _, t := range baseTimes {
		var pastCases = []struct {
			Input    string
			WantTime time.Time
		}{
			// now
			{`now`, t},

			// minutes
			{`next minute`, t.Add(time.Minute)},
			{`last minute`, t.Add(-time.Minute)},
			{`1 minute ago`, t.Add(-time.Minute)},
			{`5 minutes ago`, t.Add(-5 * time.Minute)},
			{`five minutes ago`, t.Add(-5 * time.Minute)},
			{`   5    minutes  ago   `, t.Add(-5 * time.Minute)},
			{`2 minutes from now`, t.Add(2 * time.Minute)},
			{`two minutes from now`, t.Add(2 * time.Minute)},

			// hours
			{`last hour`, t.Add(-time.Hour)},
			{`next hour`, t.Add(time.Hour)},
			{`1 hour ago`, t.Add(-time.Hour)},
			{`6 hours ago`, t.Add(-6 * time.Hour)},
			{`1 hour from now`, t.Add(time.Hour)},

			// days
			{`next day`, t.Add(24 * time.Hour)},
			{`1 day ago`, t.Add(-24 * time.Hour)},
			{`3 days ago`, t.Add(-3 * 24 * time.Hour)},
			{`3 days ago at 11:25am`, dateAtTime(t.Add(-3*24*time.Hour), 11, 25, 0)},
			{`1 day from now`, t.Add(24 * time.Hour)},

			// weeks
			{`1 week ago`, t.Add(-7 * 24 * time.Hour)},
			{`2 weeks ago`, t.Add(-2 * 7 * 24 * time.Hour)},
			{`2 weeks ago at 8am`, dateAtTime(t.Add(-2*7*24*time.Hour), 8, 0, 0)},
			{`next week`, t.Add(7 * 24 * time.Hour)},

			// months
			{`1 month ago`, t.AddDate(0, -1, 0)},
			{`last month`, t.AddDate(0, -1, 0)},
			{`next month`, t.AddDate(0, 1, 0)},
			{`1 month ago at 9:30am`, dateAtTime(t.AddDate(0, -1, 0), 9, 30, 0)},
			{`2 months ago`, t.AddDate(0, -2, 0)},
			{`12 months ago`, t.AddDate(0, -12, 0)},
			{`1 month from now`, t.AddDate(0, 1, 0)},
			{`2 months from now`, t.AddDate(0, 2, 0)},
			{`12 months from now at 6am`, dateAtTime(t.AddDate(0, 12, 0), 6, 0, 0)},

			// years
			{`last year`, t.AddDate(-1, 0, 0)},
			{`next year`, t.AddDate(1, 0, 0)},
			{`one year ago`, t.AddDate(-1, 0, 0)},
			{`one year from now`, t.AddDate(1, 0, 0)},
			{`two years ago`, t.AddDate(-2, 0, 0)},
			{`2 years ago`, t.AddDate(-2, 0, 0)},

			// today
			{`today`, t},
			{`today at 10am`, dateAtTime(t, 10, 0, 0)},

			// yesterday
			{`yesterday`, t.AddDate(0, 0, -1)},
			{`yesterday 10am`, dateAtTime(t.AddDate(0, 0, -1), 10, 0, 0)},
			{`yesterday at 10am`, dateAtTime(t.AddDate(0, 0, -1), 10, 0, 0)},
			{`yesterday at 10:15am`, dateAtTime(t.AddDate(0, 0, -1), 10, 15, 0)},

			// tomorrow
			{`tomorrow`, t.AddDate(0, 0, 1)},
			{`tomorrow 10am`, dateAtTime(t.AddDate(0, 0, 1), 10, 0, 0)},
			{`tomorrow at 10am`, dateAtTime(t.AddDate(0, 0, 1), 10, 0, 0)},
			{`tomorrow at 10:15am`, dateAtTime(t.AddDate(0, 0, 1), 10, 15, 0)},

			// past weekdays
			{`last sunday`, prevWeekday(t, time.Sunday)},
			{`past sunday`, prevWeekday(t, time.Sunday)},
			{`last monday`, prevWeekday(t, time.Monday)},
			{`last tuesday`, prevWeekday(t, time.Tuesday)},
			{`last wednesday`, prevWeekday(t, time.Wednesday)},
			{`last thursday`, prevWeekday(t, time.Thursday)},
			{`last friday`, prevWeekday(t, time.Friday)},
			{`last saturday`, prevWeekday(t, time.Saturday)},

			// future weekdays
			{`next tuesday`, nextWeekday(t, time.Tuesday)},
			{`next wednesday`, nextWeekday(t, time.Wednesday)},
			{`next thursday`, nextWeekday(t, time.Thursday)},
			{`next friday`, nextWeekday(t, time.Friday)},
			{`next saturday`, nextWeekday(t, time.Saturday)},
			{`next sunday`, nextWeekday(t, time.Sunday)},
			{`next monday`, nextWeekday(t, time.Monday)},

			// months
			{`last january`, prevMonth(t, time.January)},
			{`next january`, nextMonth(t, time.January)},

			{"january 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, t.Location())},
			{"january, 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, t.Location())},
			{"april 3 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, t.Location())},
			{"april 3, 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, t.Location())},

			// case sensitivity
			{`next December 23rd AT 5:25 PM`, nextMonthDayTime(t, time.December, 23, 12+5, 25, 0)},

			{`previous tuesday`, prevWeekday(t, time.Tuesday)},
			{`last january`, prevMonth(t, time.January)},
			{`next january`, nextMonth(t, time.January)},
		}

		for _, c := range pastCases {
			test.Run(c.Input, func(t *testing.T) {
				v, err := Parse(c.Input, base)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, c.WantTime, v)
			})
		}
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
