package datetime

import (
	"fmt"
	"strconv"
	"time"
)

// Parse is used to create a Craft style datetime format for
// backups.
func Parse(t time.Time) string {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	yr := strconv.Itoa(year)
	mon := prepend(strconv.Itoa(int(month)))
	d := prepend(strconv.Itoa(day))
	hr := prepend(strconv.Itoa(hour))
	mi := prepend(strconv.Itoa(min))
	s := prepend(strconv.Itoa(sec))

	return fmt.Sprintf("%s-%s-%s-%s%s%s", yr, mon, d, hr, mi, s)
}

func prepend(s string) string {
	if len(s) == 1 {
		return "0" + s
	}

	return s
}
