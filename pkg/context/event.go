package context

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	ID   string
	Time time.Time
	Body []string
}

func NewEvent(id string, t time.Time, body ...string) *Event {
	return &Event{
		ID:   id,
		Time: t,
		Body: body,
	}
}

func (event *Event) String() string {
	return fmt.Sprintf("%s %s %s",
		event.Time.Format("15:04"),
		event.ID,
		strings.Join(event.Body, " "))
}
