package schedule

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type TimeUnit int

const timefmt = "2006-01-02 15:04:05"

const (
	Minute TimeUnit = iota
	Hour
	Day
	Week
	Month
	Year
)

var nilTime time.Time

var tuLookup = map[string]TimeUnit{
	"Minute": Minute,
	"Hour":   Hour,
	"Day":    Day,
	"Week":   Week,
	"Month":  Month,
	"Year":   Year,
}

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

func (f Frequency) String() string {
	return fmt.Sprintf("%d %s", f.Count, f.Unit)
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
		return nilTime
	}
}

func (f Frequency) minApprox() time.Duration { return time.Duration(f.Count) * f.Unit.minApprox() }
func (f Frequency) maxApprox() time.Duration { return time.Duration(f.Count) * f.Unit.maxApprox() }

// Schedule describes a time schedule. It has a start and optional end point and an optional frequency.
type Schedule struct {
	Start, End time.Time
	Freq       Frequency
}

// NextAfter calculates the next time in the schedule after t. If no such time exists, nil is returned (test with Time.IsZero()).
func (c Schedule) NextAfter(t time.Time) time.Time {
	if !t.After(c.Start) {
		return c.Start
	}
	if c.Freq.Count == 0 {
		return nilTime
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
			return nilTime
		}
		return t2
	}

	return nilTime // Should actually never happen...
}

func (c Schedule) String() string {
	s := c.Start.UTC().Format(timefmt)
	if c.Freq.Count > 0 {
		s += " +" + c.Freq.String()
		if !c.End.IsZero() {
			s += " !" + c.End.UTC().Format(timefmt)
		}
	}
	return s
}

func ParseSchedule(s string) (c Schedule, err error) {
	elems := strings.Split(s, " ")

	switch len(elems) {
	case 6: // Everything specified
		_end := elems[4] + " " + elems[5]
		if c.End, err = time.ParseInLocation(timefmt, _end[1:], time.UTC); err != nil {
			return
		}
		fallthrough
	case 4: // start time and frequency
		var count uint64
		if count, err = strconv.ParseUint(elems[2][1:], 10, 32); err != nil {
			return
		}
		c.Freq.Count = uint(count)

		var ok bool
		if c.Freq.Unit, ok = tuLookup[elems[3]]; !ok {
			err = fmt.Errorf("Unknown timeunit %s", elems[3])
			return
		}
		fallthrough
	case 2: // Only start time
		if c.Start, err = time.ParseInLocation(timefmt, elems[0]+" "+elems[1], time.UTC); err != nil {
			return
		}
	default:
		err = errors.New("Unknown schedule format")
	}

	return
}

type MultiSchedule []Schedule

func (mc MultiSchedule) NextAfter(t time.Time) time.Time {
	var nearest time.Time

	for _, c := range mc {
		next := c.NextAfter(t)
		if next.IsZero() {
			continue
		}

		if nearest.IsZero() {
			nearest = next
		} else if next.Before(nearest) {
			nearest = next
		}
	}

	return nearest
}

func (mc MultiSchedule) String() (s string) {
	sep := ""

	for _, c := range mc {
		s += sep + c.String()
		sep = "\n"
	}

	return
}

func ParseMultiSchedule(s string) (mc MultiSchedule, err error) {
	parts := strings.Split(s, "\n")
	for l, _part := range parts {
		part := strings.TrimSpace(_part)
		if part == "" {
			continue
		}

		c, err := ParseSchedule(part)
		if err != nil {
			return nil, fmt.Errorf("Line %d: %s", l+1, err)
		}

		mc = append(mc, c)
	}

	return
}
