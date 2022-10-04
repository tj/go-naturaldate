//go:generate peg -inline -switch grammar.peg

// Package naturaldate provides natural date time parsing.
package naturaldate

import (
	"fmt"
	"time"

	gp "github.com/ijt/goparsify"
)

// day duration.
var day = time.Hour * 24

// week duration.
var week = time.Hour * 24 * 7

// Parse query string.
func Parse(s string, ref time.Time) (time.Time, error) {
	// lastMonth := gp.Bind("last month", lm)
	lastMonth := gp.Seq("last", "month").Map(func(n *gp.Result) {
		n.Result = truncateDay(ref.AddDate(0, -1, 0))
	})
	now := gp.Bind("now", ref)
	p := gp.AnyWithName("datetime",
		now, lastMonth)
	result, err := gp.Run(p, s, gp.UnicodeWhitespace)
	if err != nil {
		return time.Time{}, fmt.Errorf("running parser: %w", err)
	}
	t := result.(time.Time)
	return t, nil
}

// prevWeekday returns the previous week day relative to time t.
// TODO: test this with t = some sunday, day = time.Sunday.
func prevWeekday(t time.Time, day time.Weekday) time.Time {
	d := t.Weekday() - day
	if d <= 0 {
		d += 7
	}
	return t.Add(-time.Hour * 24 * time.Duration(d))
}

// nextWeekday returns the next week day relative to time t.
// TODO: test this with t = some sunday, day = time.Sunday.
func nextWeekday(t time.Time, day time.Weekday) time.Time {
	d := day - t.Weekday()
	if d <= 0 {
		d += 7
	}
	return t.Add(time.Hour * 24 * time.Duration(d))
}

// nextMonthDayTime returns the next month relative to time t, with given day of month and time of day.
func nextMonthDayTime(t time.Time, month time.Month, day int, hour int, min int, sec int) time.Time {
	t = nextMonth(t, month)
	return time.Date(t.Year(), t.Month(), day, hour, min, sec, 0, t.Location())
}

// prevMonthDayTime returns the previous month relative to time t, with given day of month and time of day.
func prevMonthDayTime(t time.Time, month time.Month, day int, hour int, min int, sec int) time.Time {
	t = prevMonth(t, month)
	return time.Date(t.Year(), t.Month(), day, hour, min, sec, 0, t.Location())
}

// nextMonth returns the next month relative to time t.
func nextMonth(t time.Time, month time.Month) time.Time {
	y := t.Year()
	if month-t.Month() <= 0 {
		y++
	}
	_, _, day := t.Date()
	return time.Date(y, month, day, 0, 0, 0, 0, t.Location())
}

// prevMonth returns the next month relative to time t.
func prevMonth(t time.Time, month time.Month) time.Time {
	y := t.Year()
	if t.Month()-month <= 0 {
		y--
	}
	_, _, day := t.Date()
	return time.Date(y, month, day, 0, 0, 0, 0, t.Location())
}

// truncateDay returns a date truncated to the day.
func truncateDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// truncateYear returns a date truncated to the year.
func truncateYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}
