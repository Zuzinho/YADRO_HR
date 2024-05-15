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

func (ctx *Context) SitAtTable(event *Event, client username.UserName, table *table.Table) {
	table.Owner = &client
	ctx.FreeTableCount -= 1
	ctx.Users[client] = table
	table.CurrentSession = duration.NewDuration(event.Time)
	ctx.Queue.Delete(client)
}

func (ctx *Context) GoAway(event *Event, client username.UserName) {
	oldTable := ctx.Users[client]
	delete(ctx.Users, client)

	if oldTable == nil {
		ctx.Queue.Delete(client)
		return
	}

	oldTable.OverCurrentSession(event.Time)
	ctx.FreeTableCount += 1
}
