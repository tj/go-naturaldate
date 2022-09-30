package naturaldate

import (
	"log"
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
			now := time.Time{}
			_, err := Parse(c.input, now)
			if err == nil {
				t.Errorf("err is nil")
			}
		})
	}
}

// Test parsing on cases that are expected to parse successfully.
func TestParse_goodTimes(t *testing.T) {
	var baseTimes = []time.Time{
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local),
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.UTC),
	}

	for _, now := range baseTimes {
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
			{`next December 23rd AT 5:25 PM`, nextMonthDayTime(now, time.December, 23, 12+5, 25, 0)},
		}

		for _, c := range cases {
			t.Run(c.Input, func(t *testing.T) {
				v, err := Parse(c.Input, now)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, c.WantTime, v)
			})
		}
	}
}

func TestParse_goodDays(t *testing.T) {
	var baseTimes = []time.Time{
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.Local),
		time.Date(2022, 9, 29, 2, 48, 33, 123, time.UTC),
	}

	for _, now := range baseTimes {
		var cases = []struct {
			Input    string
			WantTime time.Time
		}{
			// days
			{`one day ago`, now.Add(-24 * time.Hour)},
			{`1 day ago`, now.Add(-24 * time.Hour)},
			{`3 days ago`, now.Add(-3 * 24 * time.Hour)},
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

			// years
			{`last year`, truncateYear(now.AddDate(-1, 0, 0))},
			{`next year`, truncateYear(now.AddDate(1, 0, 0))},
			{`one year ago`, truncateYear(now.AddDate(-1, 0, 0))},
			{`one year from now`, truncateYear(now.AddDate(1, 0, 0))},
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

			{"january 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, now.Location())},
			{"january, 2017", time.Date(2017, 1, 1, 0, 0, 0, 0, now.Location())},
			{"april 3 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, now.Location())},
			{"april 3, 2017", time.Date(2017, 4, 3, 0, 0, 0, 0, now.Location())},

			{`previous tuesday`, prevWeekday(now, time.Tuesday)},
			{`last january`, prevMonth(now, time.January)},
			{`next january`, nextMonth(now, time.January)},
		}

		for _, c := range cases {
			t.Run(c.Input, func(t *testing.T) {
				v, err := Parse(c.Input, now)
				if err != nil {
					t.Fatal(err)
				}
				want := truncateDay(c.WantTime)
				assert.Equal(t, want, v)
			})
		}
	}
}

// Benchmark parsing.
func BenchmarkParse(b *testing.B) {
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		_, err := Parse(`december 23rd at 5:25pm`, time.Time{})
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}
}
