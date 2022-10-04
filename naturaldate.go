//go:generate peg -inline -switch grammar.peg

// Package naturaldate provides natural date time parsing.
package naturaldate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	gp "github.com/ijt/goparsify"
)

// day duration.
var day = time.Hour * 24

// week duration.
var week = time.Hour * 24 * 7

// Parse query string.
func Parse(s string, ref time.Time) (time.Time, error) {
	s = strings.ToLower(s)

	now := gp.Bind("now", ref)
	lastMonth := gp.Seq("last", "month").Map(func(n *gp.Result) {
		n.Result = truncateDay(ref.AddDate(0, -1, 0))
	})
	weekday := gp.AnyWithName("weekday", "mon", "tue", "wed", "thu", "fri", "sat", "sun")
	month := gp.AnyWithName("month", "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec").Map(func(n *gp.Result) {
		t, err := time.Parse("Jan", n.Token)
		if err != nil {
			panic(fmt.Sprintf("identifying month: %v", err))
		}
		n.Result = t.Month()
	})
	dayOfMonth := gp.Regex(`[0-3]?\d`).Map(func(n *gp.Result) {
		d, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing day of month: %v", err))
		}
		n.Result = d
	})
	hour := gp.Regex(`[0-2]?\d`).Map(func(n *gp.Result) {
		h, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing hour: %v", err))
		}
		n.Result = h
	})
	minute := gp.Regex(`[0-5]?\d`).Map(func(n *gp.Result) {
		m, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing minute: %v", err))
		}
		n.Result = m
	})
	second := gp.Regex(`[0-5]?\d`).Map(func(n *gp.Result) {
		s, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing second: %v", err))
		}
		n.Result = s
	})
	hourMinuteSecond := gp.Seq(hour, ":", minute, ":", second).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		m := n.Child[2].Result.(int)
		s := n.Child[4].Result.(int)
		n.Result = time.Date(1, 1, 1, h, m, s, 0, ref.Location())
	})
	zoneHour := gp.Regex(`[-+][01]\d`).Map(func(n *gp.Result) {
		h, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing time zone hour: %v", err))
		}
		n.Result = h
	})
	zone := gp.Seq(zoneHour, minute).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		m := n.Child[1].Result.(int)
		n.Result = fixedZoneHM(h, m)
	})
	year := gp.Regex(`[12]\d{3}`).Map(func(n *gp.Result) {
		y, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing year: %v", err))
		}
		n.Result = y
	})
	ansiC := gp.Seq(weekday, month, dayOfMonth, hourMinuteSecond, year).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		d := n.Child[2].Result.(int)
		t := n.Child[3].Result.(time.Time)
		y := n.Child[4].Result.(int)
		n.Result = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), 0, ref.Location())
	})
	rubyDate := gp.Seq(weekday, month, dayOfMonth, hourMinuteSecond, zone, year).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		d := n.Child[2].Result.(int)
		t := n.Child[3].Result.(time.Time)
		z := n.Child[4].Result.(*time.Location)
		y := n.Child[5].Result.(int)
		n.Result = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), 0, z)
	})
	rfc1123Z := gp.Seq(weekday, gp.Maybe(","), dayOfMonth, month, year, hourMinuteSecond, zone).Map(func(n *gp.Result) {
		d := n.Child[2].Result.(int)
		m := n.Child[3].Result.(time.Month)
		y := n.Child[4].Result.(int)
		t := n.Child[5].Result.(time.Time)
		z := n.Child[6].Result.(*time.Location)
		n.Result = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), 0, z)
	})
	p := gp.AnyWithName("datetime",
		now, lastMonth, ansiC, rubyDate, rfc1123Z)
	result, err := gp.Run(p, s, gp.UnicodeWhitespace)
	if err != nil {
		return time.Time{}, fmt.Errorf("running parser: %w", err)
	}
	t := result.(time.Time)
	return t, nil
}

func fixedZoneHM(h, m int) *time.Location {
	offset := h*60*60 + m*60
	sign := "+"
	if h < 0 {
		sign = "-"
		h = -h
	}
	name := fmt.Sprintf("%s%02d:%02d", sign, h, m)
	return time.FixedZone(name, offset)
}

func fixedZone(offsetHours int) *time.Location {
	return fixedZoneHM(offsetHours, 0)
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
