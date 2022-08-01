package jobs

import (
	"testing"
	"time"
)

type testTime struct {
	date     time.Time
	excepted bool
}

func TestIsHoliday(t *testing.T) {
	tests := []testTime{
		{time.Date(2022, time.January, 3, 0, 0, 0, 0, time.Local), true},
		{time.Date(2022, time.May, 1, 0, 0, 0, 0, time.Local), true},
		{time.Date(2022, time.July, 31, 0, 0, 0, 0, time.Local), true},
		{time.Date(2022, time.October, 7, 0, 0, 0, 0, time.Local), true},
		{time.Date(2022, time.October, 8, 0, 0, 0, 0, time.Local), false},
	}

	h := Holiday{}
	h.Init()

	for _, v := range tests {
		if h.IsHoliday(v.date) != v.excepted {
			t.Errorf("%s is holiday, but excepted %v", v.date, v.excepted)
		}
	}
}
