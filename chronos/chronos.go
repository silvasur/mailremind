package chronos

import (
	"math"
	"time"
)

type TimeUnit int

const (
	Minute TimeUnit = iota
	Hour
	Day
	Week
	Month
	Year
)

var NilTime time.Time

func (tu TimeUnit) String() string {
	switch tu {
	case Minute:
		return "Minute"
	case Hour:
		return "Hour"
	case Day:
		return "Day"
	case Week:
		return "Week"
	case Month:
		return "Month"
	case Year:
		return "Year"
	default:
		return "(Unknown TimeUnit)"
	}
}

func (tu TimeUnit) minApprox() time.Duration {
	const (
		maMinute = time.Minute
		maHour   = time.Hour
		maDay    = 24*time.Hour - time.Second
		maWeek   = 7 * maDay
		maMonth  = 28 * maDay
		maYear   = 365 * maDay
	)

	switch tu {
	case Minute:
		return maMinute
	case Hour:
		return maHour
	case Day:
		return maDay
	case Week:
		return maWeek
	case Month:
		return maMonth
	case Year:
		return maYear
	default:
		return 0
	}
}

func (tu TimeUnit) maxApprox() time.Duration {
	const (
		maMinute = time.Minute
		maHour   = time.Hour
		maDay    = 24*time.Hour + time.Second
		maWeek   = 7 * maDay
		maMonth  = 31 * maDay
		maYear   = 366 * maDay
	)

	switch tu {
	case Minute:
		return maMinute
	case Hour:
		return maHour
	case Day:
		return maDay
	case Week:
		return maWeek
	case Month:
		return maMonth
	case Year:
		return maYear
	default:
		return 0
	}
}

type Frequency struct {
	Unit  TimeUnit
	Count uint
}

func (f Frequency) addTo(t time.Time, mul uint) time.Time {
	sec := t.Second()
	min := t.Minute()
	hour := t.Hour()
	day := t.Day()
	month := t.Month()
	year := t.Year()
	loc := t.Location()

	fq := int(f.Count * mul)

	switch f.Unit {
	case Minute:
		return t.Add(time.Minute * time.Duration(fq))
	case Hour:
		return t.Add(time.Hour * time.Duration(fq))
	case Day:
		return time.Date(year, month, day+fq, hour, min, sec, 0, loc)
	case Week:
		return time.Date(year, month, day+fq*7, hour, min, sec, 0, loc)
	case Month:
		return time.Date(year, month+time.Month(fq), day, hour, min, sec, 0, loc)
	case Year:
		return time.Date(year+fq, month, day, hour, min, sec, 0, loc)
	default:
		return NilTime
	}
}

func (f Frequency) minApprox() time.Duration { return time.Duration(f.Count) * f.Unit.minApprox() }
func (f Frequency) maxApprox() time.Duration { return time.Duration(f.Count) * f.Unit.maxApprox() }

type Chronos struct {
	Start, End time.Time
	Freq       Frequency
}

func (c Chronos) NextAfter(t time.Time) time.Time {
	if !t.After(c.Start) {
		return c.Start
	}
	if c.Freq.Count == 0 {
		return NilTime
	}

	d := t.Sub(c.Start)

	fmin := uint(math.Floor(float64(d) / float64(c.Freq.maxApprox())))
	fmax := uint(math.Ceil(float64(d) / float64(c.Freq.minApprox())))

	for f := fmin; f <= fmax; f++ {
		t2 := c.Freq.addTo(c.Start, f)
		if t2.Before(c.Start) || t2.Before(t) {
			continue
		}
		if (!c.End.IsZero()) && t2.After(c.End) {
			return NilTime
		}
		return t2
	}

	return NilTime // Should actually never happen...
}
