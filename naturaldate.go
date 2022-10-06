// Package naturaldate provides natural date time parsing.
package naturaldate

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	gp "github.com/ijt/goparsify"
)

type direction int

const (
	future = iota
	past
)

type opts struct {
	defaultDirection direction
}

// DefaultToFuture sets the option to default to the future in case of
// ambiguous dates.
func DefaultToFuture(o *opts) {
	o.defaultDirection = future
}

// DefaultToPast sets the option to default to the past in case of
// ambiguous dates.
func DefaultToPast(o *opts) {
	o.defaultDirection = past
}

// Parse parses a string assumed to contain a date and possibly a time
// in one of various formats.
func Parse(s string, ref time.Time, opts ...func(o *opts)) (time.Time, error) {
	gp.EnableLogging(os.Stdout)
	s = strings.ToLower(s)
	p := Parser(ref, opts...)
	result, err := gp.Run(p, s, gp.UnicodeWhitespace)
	if err != nil {
		return time.Time{}, fmt.Errorf("running parser: %w", err)
	}
	t := result.(time.Time)
	return t, nil
}

// Parser returns a parser of dates with a given reference time called ref.
func Parser(ref time.Time, options ...func(o *opts)) gp.Parser {
	var o opts
	for _, optFunc := range options {
		optFunc(&o)
	}

	now := gp.Bind("now", ref)

	prevMo := gp.Seq("last", "month").Map(func(n *gp.Result) {
		n.Result = truncateMonth(ref.AddDate(0, -1, 0))
	})

	nextMo := gp.Seq("next", "month").Map(func(n *gp.Result) {
		n.Result = truncateMonth(ref.AddDate(0, 1, 0))
	})

	lastWeek := gp.Seq("last", "week").Map(func(n *gp.Result) {
		n.Result = truncateWeek(ref.AddDate(0, 0, -7))
	})

	nextWeek := gp.Seq("next", "week").Map(func(n *gp.Result) {
		n.Result = truncateWeek(ref.AddDate(0, 0, 7))
	})

	one := gp.Bind("one", 1)
	a := gp.Bind("a", 1)
	an := gp.Bind("an", 1)
	two := gp.Bind("two", 2)
	three := gp.Bind("three", 3)
	four := gp.Bind("four", 4)
	five := gp.Bind("five", 5)
	six := gp.Bind("six", 6)
	seven := gp.Bind("seven", 7)
	eight := gp.Bind("eight", 8)
	nine := gp.Bind("nine", 9)
	ten := gp.Bind("ten", 10)
	eleven := gp.Bind("eleven", 11)
	twelve := gp.Bind("twelve", 12)
	numeral := gp.Regex(`\d+`).Map(func(n *gp.Result) {
		num, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing numeral: %v", err))
		}
		n.Result = num
	})

	number := gp.AnyWithName("number", one, an, a, two, three, four, five, six, seven, eight, nine, ten, eleven, twelve, numeral).Map(func(n *gp.Result) {
		pass()
	})

	months := gp.Regex(`months?`)

	monthsAgo := gp.Seq(number, months, "ago").Map(func(n *gp.Result) {
		num := n.Child[0].Result.(int)
		n.Result = ref.AddDate(0, -num, 0)
	})

	monthsFromNow := gp.Seq(number, months, gp.Any(gp.Seq("from", "now"), "hence")).Map(func(n *gp.Result) {
		num := n.Child[0].Result.(int)
		n.Result = ref.AddDate(0, num, 0)
	})

	shortWeekday := gp.AnyWithName("short weekday", "mon", "tue", "wed", "thu", "fri", "sat", "sun")

	longWeekday := gp.AnyWithName("long weekday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday")

	weekday := gp.AnyWithName("weekday", longWeekday, shortWeekday).Map(func(n *gp.Result) {
		m := map[string]time.Weekday{
			"sun": time.Sunday,
			"mon": time.Monday,
			"tue": time.Tuesday,
			"wed": time.Wednesday,
			"thu": time.Thursday,
			"fri": time.Friday,
			"sat": time.Saturday,
		}
		shortName := n.Token[:3]
		day := m[shortName]
		n.Result = day
	})

	longMonth := gp.AnyWithName("long month",
		"january", "february", "march", "april",
		/* may is already short */ "june", "july", "august", "september",
		"october", "november", "december").Map(func(n *gp.Result) {
		t, err := time.Parse("January", n.Token)
		if err != nil {
			panic(fmt.Sprintf("identifying month (long): %v", err))
		}
		n.Result = t.Month()
	})

	shortMonth := gp.AnyWithName("month", "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec").Map(func(n *gp.Result) {
		t, err := time.Parse("Jan", n.Token)
		if err != nil {
			panic(fmt.Sprintf("identifying month: %v", err))
		}
		n.Result = t.Month()
	})

	shortMonthMaybeDot := gp.Seq(shortMonth, gp.Maybe(".")).Map(func(n *gp.Result) {
		n.Result = n.Child[0].Result
	})

	month := gp.AnyWithName("month", longMonth, shortMonthMaybeDot)

	lastSpecificMonth := gp.Seq("last", month).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		n.Result = prevMonth(ref, m)
	})

	nextSpecificMonth := gp.Seq("next", month).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		n.Result = nextMonth(ref, m)
	})

	monthNum := gp.Regex(`[01]?\d\b`).Map(func(n *gp.Result) {
		m, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing month number: %v", err))
		}
		n.Result = time.Month(m)
	})

	dayOfMonthNum := gp.Regex(`[0-3]?\d`).Map(func(n *gp.Result) {
		d, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing day of month: %v", err))
		}
		n.Result = d
	})

	dayOfMonthEnding := gp.Regex(`(st|nd|rd|th)`).Map(func(n *gp.Result) {
		pass()
	})

	dayOfMonth := gp.Seq(dayOfMonthNum, gp.Maybe(dayOfMonthEnding)).Map(func(n *gp.Result) {
		n.Result = n.Child[0].Result
	})

	hour12 := gp.Regex(`[0-1]?\d`).Map(func(n *gp.Result) {
		h, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing hour (12h clock): %v", err))
		}
		n.Result = h
	})

	hour24 := gp.Regex(`[0-2]?\d\b`).Map(func(n *gp.Result) {
		h, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing hour (24h clock): %v", err))
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

	// Second can go up to 60 because of leap seconds, for example
	// 1990-12-31T15:59:60-08:00.
	second := gp.Regex(`[0-6]?\d\b`).Map(func(n *gp.Result) {
		s, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing second: %v", err))
		}
		n.Result = s
	})

	amPM := gp.AnyWithName("AM or PM", "am", "pm")

	colonSecond := gp.Seq(":", second).Map(func(n *gp.Result) {
		n.Result = n.Child[1].Result
	})

	colonMinute := gp.Seq(gp.Maybe(":"), minute).Map(func(n *gp.Result) {
		n.Result = n.Child[1].Result
	})

	colonMinuteColonSecond := gp.Seq(colonMinute, gp.Maybe(colonSecond)).Map(func(n *gp.Result) {
		m := n.Child[0].Result.(int)
		c1 := n.Child[1].Result
		s := 0
		if c1 != nil {
			s = c1.(int)
		}
		n.Result = time.Date(1, 1, 1, 0, m, s, 0, ref.Location())
	})

	hour12MinuteSecond := gp.Seq(hour12, gp.Maybe(colonMinuteColonSecond), amPM).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		c1 := n.Child[1].Result
		m := 0
		s := 0
		if c1 != nil {
			ms := n.Child[1].Result.(time.Time)
			m = ms.Minute()
			s = ms.Second()
		}
		if n.Child[2].Token == "pm" {
			h += 12
		}
		n.Result = time.Date(1, 1, 1, h, m, s, 0, ref.Location())
	})

	hour24MinuteSecond := gp.Seq(hour24, colonMinute, gp.Maybe(colonSecond)).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		m := n.Child[1].Result.(int)
		s := 0
		c2 := n.Child[2].Result
		if c2 != nil {
			s = c2.(int)
		}
		n.Result = time.Date(1, 1, 1, h, m, s, 0, ref.Location())
	})

	hourMinuteSecond := gp.AnyWithName("h:m:s", hour12MinuteSecond, hour24MinuteSecond)

	zoneHour := gp.Regex(`[-+][01]?\d`).Map(func(n *gp.Result) {
		h, err := strconv.Atoi(n.Token)
		if err != nil {
			panic(fmt.Sprintf("parsing time zone hour: %v", err))
		}
		n.Result = h
	})

	zoneOffset := gp.Seq(zoneHour, gp.Maybe(colonMinute)).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		c1 := n.Child[1].Result
		m := 0
		if c1 != nil {
			m = c1.(int)
		}
		n.Result = fixedZoneHM(h, m)
	})

	zoneUTC := gp.Seq("utc", gp.Maybe(zoneOffset)).Map(func(n *gp.Result) {
		c1 := n.Child[1].Result
		z := time.UTC
		if c1 != nil {
			z = c1.(*time.Location)
		}
		n.Result = z
	})

	zoneZ := gp.Bind("z", time.UTC)

	zone := gp.AnyWithName("time zone", zoneUTC, zoneOffset, zoneZ).Map(func(n *gp.Result) {
		pass()
	})

	year := gp.Regex(`[12]\d{3}\b`).Map(func(n *gp.Result) {
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

	rfc1123Z := gp.Seq(weekday, gp.Maybe(","), dayOfMonth, month, year, hourMinuteSecond, gp.Cut(), zone).Map(func(n *gp.Result) {
		d := n.Child[2].Result.(int)
		m := n.Child[3].Result.(time.Month)
		y := n.Child[4].Result.(int)
		t := n.Child[5].Result.(time.Time)
		z := n.Child[7].Result.(*time.Location)
		n.Result = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), 0, z)
	})

	rfc3339 := gp.Regex(`[12]\d{3}-[01]\d-[0-3]\dt[0-2]\d:[0-5]\d:[0-6]\d(z|[-+][01]\d:\d\d)`).Map(func(n *gp.Result) {
		t, err := time.Parse(time.RFC3339, strings.ToUpper(n.Token))
		if err != nil {
			panic(fmt.Sprintf("parsing time in RFC3339 format: %v", err))
		}
		n.Result = t
	})

	dmyDate := gp.Seq(dayOfMonth, gp.Maybe("of"), month, gp.Maybe(","), year).Map(func(n *gp.Result) {
		d := n.Child[0].Result.(int)
		m := n.Child[2].Result.(time.Month)
		y := n.Child[4].Result.(int)
		n.Result = time.Date(y, m, d, 0, 0, 0, 0, ref.Location())
	})

	// "my" here stands for "month, year"
	myDate := gp.Seq(month, gp.Maybe(","), year).Map(func(n *gp.Result) {
		m := n.Child[0].Result.(time.Month)
		y := n.Child[2].Result.(int)
		n.Result = time.Date(y, m, 1, 0, 0, 0, 0, ref.Location())
	})

	mdyDate := gp.Seq(month, gp.Maybe(dayOfMonth), gp.Maybe(","), year).Map(func(n *gp.Result) {
		m := n.Child[0].Result.(time.Month)
		d := 1
		c1 := n.Child[1].Result
		if c1 != nil {
			d = c1.(int)
		}
		y := n.Child[3].Result.(int)
		n.Result = time.Date(y, m, d, 0, 0, 0, 0, ref.Location())
	})

	date := gp.AnyWithName("date", dmyDate, myDate, mdyDate)

	atTimeWithMaybeZone := gp.Seq(gp.Maybe("at"), hourMinuteSecond, gp.Maybe(zone)).Map(func(n *gp.Result) {
		t := n.Child[1].Result.(time.Time)
		z := ref.Location()
		c2 := n.Child[2].Result
		if c2 != nil {
			z = c2.(*time.Location)
		}
		n.Result = time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, z)
	})

	lastWeekday := gp.Seq("last", weekday, gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		day := n.Child[1].Result.(time.Weekday)
		d := prevWeekdayFrom(ref, day)
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	nextWeekday := gp.Seq("next", weekday, gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		day := n.Child[1].Result.(time.Weekday)
		d := nextWeekdayFrom(ref, day)
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	lastSpecificMonthDay := gp.Seq("last", month, dayOfMonth, gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		d := n.Child[2].Result.(int)
		t := prevMonth(ref, m)
		t = time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, ref.Location())
		n.Result = setTimeMaybe(t, n.Child[3].Result)
	})

	nextSpecificMonthDay := gp.Seq("next", month, dayOfMonth, gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		m := n.Child[1].Result.(time.Month)
		d := n.Child[2].Result.(int)
		t := nextMonth(ref, m)
		t = time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, ref.Location())
		n.Result = setTimeMaybe(t, n.Child[3].Result)
	})

	todayTime := gp.Seq("today", gp.Cut(), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		d := truncateDay(ref)
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	timeToday := gp.Seq(atTimeWithMaybeZone, "today").Map(func(n *gp.Result) {
		d := truncateDay(ref)
		n.Result = setTimeMaybe(d, n.Child[0].Result)
	})

	today := gp.Any(timeToday, todayTime)

	yesterdayTime := gp.Seq("yesterday", gp.Cut(), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		d := truncateDay(ref.AddDate(0, 0, -1))
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	timeYesterday := gp.Seq(atTimeWithMaybeZone, "yesterday").Map(func(n *gp.Result) {
		d := truncateDay(ref.AddDate(0, 0, -1))
		n.Result = setTimeMaybe(d, n.Child[0].Result)
	})

	yesterday := gp.Any(timeYesterday, yesterdayTime)

	tomorrowTime := gp.Seq("tomorrow", gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		d := truncateDay(ref.AddDate(0, 0, 1))
		n.Result = setTimeMaybe(d, n.Child[1].Result)
	})

	timeTomorrow := gp.Seq(atTimeWithMaybeZone, "tomorrow").Map(func(n *gp.Result) {
		d := truncateDay(ref.AddDate(0, 0, 1))
		n.Result = setTimeMaybe(d, n.Child[0].Result)
	})

	tomorrow := gp.Any(timeTomorrow, tomorrowTime)

	dateTime := gp.Seq(date, gp.Maybe(","), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		d := n.Child[0].Result.(time.Time)
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	slash := gp.Any("/", "-")

	ymdNumbers := gp.Seq(year, slash, monthNum, slash, dayOfMonthNum).Map(func(n *gp.Result) {
		y := n.Child[0].Result.(int)
		m := n.Child[2].Result.(time.Month)
		d := n.Child[4].Result.(int)
		n.Result = time.Date(y, m, d, 0, 0, 0, 0, ref.Location())
	})

	dmyNumbers := gp.Seq(dayOfMonthNum, slash, monthNum, slash, year).Map(func(n *gp.Result) {
		d := n.Child[0].Result.(int)
		m := n.Child[2].Result.(time.Month)
		y := n.Child[4].Result.(int)
		n.Result = time.Date(y, m, d, 0, 0, 0, 0, ref.Location())
	})

	dateAsNumbers := gp.AnyWithName("date as numbers", ymdNumbers, dmyNumbers)

	dateAsNumbersMaybeZone := gp.Seq(dateAsNumbers, gp.Maybe(zone)).Map(func(n *gp.Result) {
		d := n.Child[0].Result.(time.Time)
		c1 := n.Child[1].Result
		z := d.Location()
		if c1 != nil {
			z = c1.(*time.Location)
		}
		n.Result = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, z)
	})

	dateZone := gp.Seq(date, zone).Map(func(n *gp.Result) {
		d := n.Child[0].Result.(time.Time)
		z := n.Child[1].Result.(*time.Location)
		n.Result = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, z)
	})

	lastYear := gp.Seq("last", "year").Map(func(n *gp.Result) {
		n.Result = truncateYear(ref.AddDate(-1, 0, 0))
	})

	nextYear := gp.Seq("next", "year").Map(func(n *gp.Result) {
		n.Result = truncateYear(ref.AddDate(1, 0, 0))
	})

	yearsLabel := gp.Regex(`years?`)

	xYearsAgo := gp.Seq(number, yearsLabel, "ago").Map(func(n *gp.Result) {
		y := n.Child[0].Result.(int)
		n.Result = ref.AddDate(-y, 0, 0)
	})

	fromNowOrToday := gp.Any("hence", gp.Seq("from", gp.Any("now", "today")))

	xYearsFromToday := gp.Seq(number, yearsLabel, fromNowOrToday).Map(func(n *gp.Result) {
		y := n.Child[0].Result.(int)
		n.Result = ref.AddDate(y, 0, 0)
	})

	daysLabel := gp.Regex(`days?`)

	xDaysAgo := gp.Seq(number, daysLabel, "ago", gp.Cut(), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		delta := n.Child[0].Result.(int)
		d := ref.AddDate(0, 0, -delta)
		n.Result = setTimeMaybe(d, n.Child[4].Result)
	})

	xDaysFromNow := gp.Seq(number, daysLabel, fromNowOrToday, gp.Cut(), gp.Maybe(atTimeWithMaybeZone), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		delta := n.Child[0].Result.(int)
		d := ref.AddDate(0, 0, delta)
		n.Result = setTimeMaybe(d, n.Child[4].Result)
	})

	weeksLabel := gp.Regex(`weeks?`)

	xWeeksAgo := gp.Seq(number, weeksLabel, "ago", gp.Cut(), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		delta := n.Child[0].Result.(int)
		d := ref.AddDate(0, 0, -7*delta)
		n.Result = setTimeMaybe(d, n.Child[4].Result)
	})

	xWeeksFromNow := gp.Seq(number, weeksLabel, fromNowOrToday, gp.Cut(), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		delta := n.Child[0].Result.(int)
		d := ref.AddDate(0, 0, 7*delta)
		n.Result = setTimeMaybe(d, n.Child[4].Result)
	})

	minutesLabel := gp.Regex(`minutes?`)

	xMinutesAgo := gp.Seq(number, minutesLabel, "ago").Map(func(n *gp.Result) {
		m := n.Child[0].Result.(int)
		n.Result = ref.Add(-time.Duration(m) * time.Minute)
	})

	fromNow := gp.Any("hence", gp.Seq("from", "now"))

	xMinutesFromNow := gp.Seq(number, minutesLabel, fromNow).Map(func(n *gp.Result) {
		m := n.Child[0].Result.(int)
		n.Result = ref.Add(time.Duration(m) * time.Minute)
	})

	hoursLabel := gp.Regex(`hours?`)

	xHoursAgo := gp.Seq(number, hoursLabel, "ago").Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		n.Result = ref.Add(-time.Duration(h) * time.Hour)
	})

	xHoursFromNow := gp.Seq(number, hoursLabel, fromNow).Map(func(n *gp.Result) {
		h := n.Child[0].Result.(int)
		n.Result = ref.Add(time.Duration(h) * time.Hour)
	})

	color := gp.AnyWithName("color",
		"white", "red", "green", "blue", "gold", "purple", "orange", "pink",
		"silver", "copper")

	colorMonth := gp.Seq(color, month).Map(func(n *gp.Result) {
		c := n.Child[0].Token
		m := n.Child[1].Result.(time.Month)
		color2delta := map[string]int{
			"white":  0,
			"red":    1,
			"green":  2,
			"blue":   3,
			"gold":   4,
			"purple": 5,
			"orange": 6,
			"pink":   7,
			"silver": 8,
			"copper": 9,
		}
		delta := color2delta[c]
		t := nextMonth(ref, m)
		n.Result = t.AddDate(delta, 0, 0)
	})

	monthNoYear := gp.Seq(month, gp.Maybe(dayOfMonth), gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		var d time.Time
		m := n.Child[0].Result.(time.Month)
		switch o.defaultDirection {
		case future:
			d = nextMonth(ref, m)
		case past:
			d = prevMonth(ref, m)
		default:
			panic(fmt.Sprintf("invalid default direction: %q", o.defaultDirection))
		}
		d = setDayMaybe(d, n.Child[1].Result)
		n.Result = setTimeMaybe(d, n.Child[2].Result)
	})

	weekdayNoDirection := gp.Seq(weekday, gp.Maybe(atTimeWithMaybeZone)).Map(func(n *gp.Result) {
		w := n.Child[0].Result.(time.Weekday)
		var t time.Time
		switch o.defaultDirection {
		case future:
			t = nextWeekdayFrom(ref, w)
		case past:
			t = prevWeekdayFrom(ref, w)
		default:
			panic(fmt.Sprintf("invalid default direction: %q", o.defaultDirection))
		}
		n.Result = setTimeMaybe(t, n.Child[1].Result)
	})

	return gp.AnyWithName("natural date",
		now, today, yesterday, tomorrow,
		ansiC, rubyDate, rfc1123Z, rfc3339,
		dateZone, dateTime,
		dateAsNumbersMaybeZone,
		xMinutesAgo, xMinutesFromNow,
		xHoursAgo, xHoursFromNow,
		xDaysAgo, xDaysFromNow,
		xWeeksAgo, xWeeksFromNow,
		monthsAgo, monthsFromNow,
		xYearsAgo, xYearsFromToday,
		lastSpecificMonthDay, nextSpecificMonthDay,
		lastSpecificMonth, nextSpecificMonth,
		lastYear, nextYear,
		nextMo, prevMo,
		lastWeekday, nextWeekday,
		lastWeek, nextWeek,
		colorMonth, monthNoYear,
		weekdayNoDirection)
}

func setTimeMaybe(datePart time.Time, timePart any) time.Time {
	d := datePart
	if timePart == nil {
		return d
	}
	t := timePart.(time.Time)
	return time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func setDayMaybe(t time.Time, dayAsAny any) time.Time {
	if dayAsAny == nil {
		return t
	}
	d := dayAsAny.(int)
	return time.Date(t.Year(), t.Month(), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
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

// prevWeekdayFrom returns the previous week day relative to time t.
// TODO: test this with t = some sunday, day = time.Sunday.
func prevWeekdayFrom(t time.Time, day time.Weekday) time.Time {
	d := t.Weekday() - day
	if d <= 0 {
		d += 7
	}
	return truncateDay(t.AddDate(0, 0, -int(d)))
}

// nextWeekdayFrom returns the next week day relative to time t.
// TODO: test this with t = some sunday, day = time.Sunday.
func nextWeekdayFrom(t time.Time, day time.Weekday) time.Time {
	d := day - t.Weekday()
	if d <= 0 {
		d += 7
	}
	return truncateDay(t.AddDate(0, 0, int(d)))
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

// truncateWeek returns a date truncated to the week.
func truncateWeek(t time.Time) time.Time {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// truncateMonth returns a date truncated to the month.
func truncateMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// truncateYear returns a date truncated to the year.
func truncateYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// setTime takes the date from d and the time from the remaining args and
// returns the combined result.
func setTime(d time.Time, h, m, s, ns int) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), h, m, s, ns, d.Location())
}

// pass is something to call so you can put a breakpoint in an empty func.
func pass() {
}
