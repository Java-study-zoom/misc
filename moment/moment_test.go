package moment

import (
	"testing"

	"time"
)

func TestString(t *testing.T) {
	now := time.Now()
	for _, test := range []struct {
		dur  time.Duration
		want string
	}{
		{0, "just now"},
		{time.Second, "just now"},
		{3 * time.Second, "just now"},
		{59 * time.Second, "just now"},
		{time.Minute, "a minute ago"},
		{1*time.Minute + 30*time.Second, "a minute ago"},
		{3 * time.Minute, "3 minutes ago"},
		{5*time.Minute + 27*time.Second, "5 minutes ago"},
		{12 * time.Minute, "10 minutes ago"},
		{59 * time.Minute, "55 minutes ago"},
		{time.Hour, "an hour ago"},
		{time.Hour + 59*time.Minute, "an hour ago"},
		{time.Hour * 2, "2 hours ago"},
		{time.Hour * 24, "yesterday"},
		{time.Hour * 47, "yesterday"}, // this might not be so correct
	} {
		got := String(now.Add(-test.dur), now)
		if got != test.want {
			t.Errorf(
				"moment string for %s: got %q, want %q",
				test.dur, got, test.want,
			)
		}
	}

	for _, test := range []struct {
		hour int
		want string
	}{
		{10, "13 hours ago"},
		{1, "13 hours ago"},
		{23, "yesterday"},
		{12, "yesterday"},
	} {
		dummyTime := time.Date(1, 1, 1, test.hour, 0, 0, 0, time.UTC)
		dur := 13 * time.Hour
		got := String(dummyTime, dummyTime.Add(dur))
		if got != test.want {
			t.Errorf(
				"moment string for 13 hours afer %q: got %q, want %q",
				dummyTime, got, test.want,
			)
		}
	}
}
