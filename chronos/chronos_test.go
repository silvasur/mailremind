package chronos

import (
	"testing"
	"time"
)

func mktime(y int, month time.Month, d, h, min int) time.Time {
	return time.Date(y, month, d, h, min, 0, 0, time.Local)
}

func TestChronos(t *testing.T) {
	tbl := []struct {
		start time.Time
		end   time.Time
		unit  TimeUnit
		count uint
		now   time.Time
		want  time.Time
	}{
		{mktime(1991, 4, 30, 0, 0), NilTime, Year, 1, mktime(2013, 8, 26, 13, 37), mktime(2014, 4, 30, 0, 0)},
		{mktime(2013, 1, 1, 0, 0), NilTime, Year, 0, mktime(2013, 8, 26, 13, 37), NilTime},
		{mktime(2013, 1, 1, 0, 0), NilTime, Year, 0, mktime(2012, 1, 1, 0, 0), mktime(2013, 1, 1, 0, 0)},
		{mktime(1900, 12, 24, 12, 34), NilTime, Year, 5, mktime(2013, 8, 26, 13, 37), mktime(2015, 12, 24, 12, 34)},
		{mktime(1900, 12, 24, 12, 34), mktime(2010, 1, 1, 1, 1), Year, 5, mktime(2013, 8, 26, 13, 37), NilTime},
		{mktime(2013, 8, 1, 4, 2), NilTime, Week, 3, mktime(2013, 8, 26, 13, 37), mktime(2013, 9, 12, 4, 2)},
		{mktime(2013, 8, 26, 13, 37), NilTime, Year, 0, mktime(2013, 8, 26, 13, 37), mktime(2013, 8, 26, 13, 37)},
		{mktime(2013, 8, 25, 13, 37), NilTime, Day, 1, mktime(2013, 8, 26, 13, 37), mktime(2013, 8, 26, 13, 37)},
	}

	for i, e := range tbl {
		have := Chronos{e.start, e.end, Frequency{e.unit, e.count}}.NextAfter(e.now)
		if !have.Equal(e.want) {
			t.Errorf("#%d: Want: %s, Have: %s", i, e.want, have)
		}
	}
}
