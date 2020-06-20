package utils

import (
	"time"
)

// ToDateStarted coverts the date to the beginning of the date
func ToDateStarted(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// ToDateEnded coverts the date to the beginning of the date
func ToDateEnded(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, date.Location())
}

// ToMonthStarted coverts the date to the beginning of the date
func ToMonthStarted(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

// FirstDateOfWeek returns the beginning datetime of the week according to the week start day
func FirstDateOfWeek(date time.Time, ws time.Weekday) time.Time {
	return ToDateStarted(date.AddDate(0, 0, int(ws)-int(date.Weekday())))
}
