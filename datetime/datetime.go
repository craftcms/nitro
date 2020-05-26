package datetime

import (
	"strconv"
	"time"
)

// Parse is used to create a Craft style datetime format for
// backups.
func Parse(t time.Time) string {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	y := strconv.Itoa(year)
	m := prepend(strconv.Itoa(int(month)))
	d := prepend(strconv.Itoa(day))
	h := prepend(strconv.Itoa(hour))
	mi := prepend(strconv.Itoa(min))
	s := prepend(strconv.Itoa(sec))

	return y[len(y)-2:] + m + d + "_" + h + mi + s
}

func prepend(s string) string {
	if len(s) == 1 {
		return "0" + s
	}

	return s
}
