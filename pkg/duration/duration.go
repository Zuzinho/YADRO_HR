package duration

import (
	"time"
)

const MinutesPerHour = 60

type Duration struct {
	Start time.Time
	End   time.Time
}

func NewDuration(start time.Time) *Duration {
	return &Duration{
		Start: start,
	}
}

func NewDurationWithEnd(start, end time.Time) *Duration {
	return &Duration{
		Start: start,
		End:   end,
	}
}

func (duration *Duration) Hours() int {
	diff := duration.End.Sub(duration.Start)

	minutes := int(diff.Minutes())

	hours := minutes / MinutesPerHour

	if minutes%MinutesPerHour > 0 {
		return hours + 1
	}

	return hours
}

func (duration *Duration) Minutes() int {
	diff := duration.End.Sub(duration.Start)

	return int(diff.Minutes())
}

func (duration *Duration) InTime(t time.Time) bool {
	return (duration.Start.Before(t) && duration.End.After(t)) || duration.Start.Equal(t) || duration.End.Equal(t)
}

type Durations []*Duration

func NewDurations() *Durations {
	durs := make(Durations, 0)

	return &durs
}

func (durations *Durations) Append(duration *Duration) {
	*durations = append(*durations, duration)
}
