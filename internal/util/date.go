package util

import "time"

func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetDayDiff(older, newer time.Time) int {
	if IsSameDate(older, newer) {
		return 0
	}
	days := int(newer.Truncate(24*time.Hour).Sub(older.Truncate(24*time.Hour)).Hours() / 24)
	return days
}

func GetNextDayOf(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func GetPrevDayOf(t time.Time) time.Time {
	return t.AddDate(0, 0, -1)
}

func GetDateWithDaysDiff(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}
