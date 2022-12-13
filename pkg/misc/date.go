package misc

import (
	"fmt"
	"time"
)

const (
	hoursInDay    = 24
	minutesInHour = 60
	tenDays       = time.Hour * hoursInDay * 10
)

func RelativeDate(d time.Time) string {
	diff := time.Since(d)
	if diff > tenDays {
		return d.Format("Mon 2 Jan 2006")
	}
	if diff.Hours() > hoursInDay {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/hoursInDay))
	}
	if diff.Minutes() > minutesInHour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}
	return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
}
