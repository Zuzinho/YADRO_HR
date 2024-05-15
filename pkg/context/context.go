package context

import (
	"fmt"
	"time"
	"yadro/pkg/duration"
	"yadro/pkg/table"
	"yadro/pkg/username"
)

type Context struct {
	Price          int
	Tables         *table.Tables
	WorkDuration   *duration.Duration
	Users          map[username.UserName]*table.Table
	Queue          *username.Queue
	FreeTableCount int
}

func NewContext(tablesCount int, price int, start, end time.Time) *Context {
	return &Context{
		Price:          price,
		Tables:         table.NewTablesWithLen(tablesCount),
		WorkDuration:   duration.NewDurationWithEnd(start, end),
		Users:          make(map[username.UserName]*table.Table),
		Queue:          username.NewQueue(),
		FreeTableCount: tablesCount,
	}
}

func (ctx *Context) TablesMoney() []string {
	results := make([]string, len(*ctx.Tables))

	for i, t := range *ctx.Tables {
		results[i] = fmt.Sprintf("%d %d %s",
			t.ID,
			t.Hours()*ctx.Price,
			t.CommonTime())
	}

	return results
}
