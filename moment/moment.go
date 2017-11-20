package moment

import (
	"fmt"
	"time"
)

// ForNow returns a human friendly string for describing a time point
// t relative to now. It reads the system's current time.
func ForNow(t time.Time) string {
	return String(t, time.Now())
}

// String returns a human friendly string for describing a time point t
// relative to now.
func String(t, now time.Time) string {
	localT := t.In(now.Location())
	d := now.Sub(localT)
	if d < 0 {
		return "in the future"
	}

	secs := int64(d / time.Second)
	if secs < 60 {
		return "just now"
	}

	mins := int64(d / time.Minute)
	if mins <= 1 {
		return "a minute ago"
	}
	if mins < 10 {
		return fmt.Sprintf("%d minutes ago", mins)
	}
	if mins < 60 {
		mins -= mins % 5
		return fmt.Sprintf("%d minutes ago", mins)
	}

	hours := int64(d / time.Hour)
	if hours <= 1 {
		return "an hour ago"
	}
	if hours <= 12 || localT.Day() == now.Day() {
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int64(d / (time.Hour * 24))
	if days <= 1 {
		return "yesterday"
	}
	if days < 30 {
		return fmt.Sprintf("%d days ago", days)
	}

	months := int64(d / (time.Hour * 24 * 30))
	if months <= 1 {
		return "a month ago"
	}
	if months < 12 {
		return fmt.Sprintf("%d months ago", months)
	}

	years := int64(d / (time.Hour * 24 * 365))
	if years <= 1 {
		return "a year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}
