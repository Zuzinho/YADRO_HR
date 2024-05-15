package table

import (
	"time"
	"yadro/pkg/duration"
	"yadro/pkg/username"
)

type Table struct {
	ID             int
	Owner          *username.UserName
	CurrentSession *duration.Duration
	History        *duration.Durations
}

func NewTable(id int) *Table {
	return &Table{
		ID:      id,
		Owner:   nil,
		History: duration.NewDurations(),
	}
}

func (table *Table) Hours() int {
	count := 0
	for _, dur := range *table.History {
		count += dur.Hours()
	}

	return count
}

func (table *Table) OverCurrentSession(end time.Time) {
	table.CurrentSession.End = end
	table.History.Append(table.CurrentSession)
	table.CurrentSession = nil

	table.Owner = nil
}

func (table *Table) CommonTime() string {
	commonMinutes := 0
	for _, dur := range *table.History {
		commonMinutes += dur.Minutes()
	}

	return time.Date(0, 0, 0,
		commonMinutes/duration.MinutesPerHour,
		commonMinutes%duration.MinutesPerHour,
		0, 0, time.UTC).Format("15:04")
}

type Tables []*Table

func NewTablesWithLen(len int) *Tables {
	tbls := make(Tables, len)
	for i := range tbls {
		tbls[i] = NewTable(i + 1)
	}

	return &tbls
}
