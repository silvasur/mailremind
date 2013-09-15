package schedule

import (
	"testing"
	"time"
)

func mktime(y int, month time.Month, d, h, min int) time.Time {
	return time.Date(y, month, d, h, min, 0, 0, time.UTC)
}

func TestSchedule(t *testing.T) {
	tbl := []struct {
		schedule string
		now      time.Time
		want     time.Time
	}{
		{"1991-04-30 00:00:00 +1 Year", mktime(2013, 8, 26, 13, 37), mktime(2014, 4, 30, 0, 0)},
		{"2013-01-01 00:00:00", mktime(2013, 8, 26, 13, 37), nilTime},
		{"2013-01-01 00:00:00", mktime(2012, 1, 1, 0, 0), mktime(2013, 1, 1, 0, 0)},
		{"1900-12-24 12:34:00 +5 Year", mktime(2013, 8, 26, 13, 37), mktime(2015, 12, 24, 12, 34)},
		{"1900-12-24 12:34:00 +5 Year !2010-01-01 01:01:00", mktime(2013, 8, 26, 13, 37), nilTime},
		{"2013-08-01 04:02:00 +3 Week", mktime(2013, 8, 26, 13, 37), mktime(2013, 9, 12, 4, 2)},
		{"2013-08-26 13:37:00", mktime(2013, 8, 26, 13, 37), mktime(2013, 8, 26, 13, 37)},
		{"2013-08-25 13:37:00 +1 Day", mktime(2013, 8, 26, 13, 37), mktime(2013, 8, 26, 13, 37)},
		{"2012-12-31 23:59:00 +100 Minute", mktime(2013, 1, 1, 0, 0), mktime(2013, 1, 1, 1, 39)},
	}

	for i, e := range tbl {
		c, err := ParseSchedule(e.schedule)
		if err != nil {
			t.Errorf("#%d: Failed parsing \"%s\": %s", i, e.schedule, err)
			continue
		}
		have := c.NextAfter(e.now)
		if !have.Equal(e.want) {
			t.Errorf("#%d: Want: %s, Have: %s", i, e.want, have)
		}

		if s := c.String(); s != e.schedule {
			t.Errorf("#%d: String() failed: \"%s\" != \"%s\"", i, e.schedule, s)
		}

	}
}
