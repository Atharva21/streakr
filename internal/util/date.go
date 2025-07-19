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

func CompareDate(t1, t2 time.Time) int {
	if IsSameDate(t1, t2) {
		return 0
	}

	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	// Compare years first
	if y1 < y2 {
		return 1
	}
	if y1 > y2 {
		return -1
	}

	// Years are equal, compare months
	if m1 < m2 {
		return 1
	}
	if m1 > m2 {
		return -1
	}

	// Years and months are equal, compare days
	if d1 < d2 {
		return 1
	}
	return -1 // d1 > d2
}

func FallInSameMonthYear(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

func AtLeastOneMonthOlder(t1, t2 time.Time) bool {
	if CompareDate(t1, t2) <= 0 {
		return false
	}
	if t1.Year() == t2.Year() && t1.Month() == t2.Month() {
		return false
	}
	return true
}

// func AtLeastOneMonthOlder(t1, t2 time.Time) bool {
// 	if t1.Year() > t2.Year() {
// 		return false
// 	}
// 	if t1.Year() == t2.Year() {
// 		if t1.Month() >= t2.Month() {
// 			return false
// 		}
// 	}
// 	return true
// }
