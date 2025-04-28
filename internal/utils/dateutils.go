package utils

import "time"

func Now() time.Time {
	return CleanDate(time.Now())
}

func CleanDate(date time.Time) time.Time {
	year, month, day := date.Date()
	clean := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return clean
}

func FormatDate(date time.Time) string {
	return date.Format("02-01-2006")
}

func ParseDate(dateStr string) (time.Time, error) {
	var date time.Time
	var err error

	dateLayouts := []string{
		"02.01.2006",
		"02/01/2006",
		"02-01-2006",
		"02.01.06",
		"02/01/06",
		"02-01-06",
	}
	for _, layout := range dateLayouts {
		date, err = time.Parse(layout, dateStr)
		if err == nil {
			break
		}
	}
	return date, err
}
