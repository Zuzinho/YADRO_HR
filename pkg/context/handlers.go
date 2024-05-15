package context

import (
	"sort"
	"strconv"
	"yadro/pkg/username"
)

type HandlerFunc func(event *Event, ctx *Context) (*Event, error)

type Handler map[string]HandlerFunc

func NewHandler() *Handler {

	handlers := Handler{
		onClientComeStatus:       onClientCome,
		onClientSitAtTableStatus: onClientSitAtTable,
		onClientWaitStatus:       onClientWait,
		onClientGoAwayStatus:     onClientGoAway,
	}

	return &handlers
}

func (handler *Handler) Handle(event *Event, ctx *Context) (*Event, error) {
	return (*handler)[event.ID](event, ctx)
}

func (handler *Handler) Close(ctx *Context) ([]*Event, error) {
	userNames := make([]username.UserName, len(ctx.Users))

	index := 0
	for userName := range ctx.Users {
		userNames[index] = userName
		index++
	}

	userNames = append(userNames, ctx.Queue.Array()...)

	sort.Slice(userNames, func(i, j int) bool {
		return userNames[j] > userNames[i]
	})

	events := make([]*Event, len(userNames))
	for i, userName := range userNames {
		e, err := onSystemClientGoAway(
			NewEvent(systemOnClientGoAwayStatus, ctx.WorkDuration.End, string(userName)),
			ctx)
		if err != nil {
			return nil, err
		}

		events[i] = e
	}

	return events, nil
}

func onClientCome(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, onClientComeBodyLen); err != nil {
		return nil, err
	}

	client := username.UserName(event.Body[0])

	if _, ok := ctx.Users[client]; ok {
		return NewEvent(systemOnErrorStatus, event.Time, youShallNotPassError), nil
	}

	if !ctx.WorkDuration.InTime(event.Time) {
		return NewEvent(systemOnErrorStatus, event.Time, notOpenYetError), nil
	}

	ctx.Users[client] = nil

	return nil, nil
}

func onClientSitAtTable(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, onClientSitAtTableBodyLen); err != nil {
		return nil, err
	}

	client := username.UserName(event.Body[0])
	tableID, err := strconv.Atoi(event.Body[1])
	if err != nil {
		return nil, err
	}

	table := (*ctx.Tables)[tableID-1]
	if table.Owner != nil {
		return NewEvent(systemOnErrorStatus, event.Time, placeIsBusyError), nil
	}

	if oldTable, ok := ctx.Users[client]; !ok {
		return NewEvent(systemOnErrorStatus, event.Time, clientUnknownError), nil
	} else if oldTable != nil {
		oldTable.OverCurrentSession(event.Time)
		ctx.FreeTableCount += 1
	}

	ctx.SitAtTable(event, client, table)

	return nil, nil
}

func onClientWait(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, onClientWaitBodyLen); err != nil {
		return nil, err
	}

	if ctx.FreeTableCount > 0 {
		return NewEvent(systemOnErrorStatus, event.Time, iCanWaitNoLongerError), nil
	}

	if ctx.Queue.Len() > len(*ctx.Tables)*2 {
		return onSystemClientGoAway(event, ctx)
	}

	client := username.UserName(event.Body[0])
	ctx.Queue.Append(client)

	return nil, nil
}

func onClientGoAway(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, onClientGoAwayBodyLen); err != nil {
		return nil, err
	}

	client := username.UserName(event.Body[0])
	oldTable, ok := ctx.Users[client]
	if !ok {
		return NewEvent(systemOnErrorStatus, event.Time, clientUnknownError), nil
	}

	ctx.GoAway(event, client)

	return onSystemClientSitAtTable(
		NewEvent(onClientSitAtTableStatus, event.Time, append(event.Body, strconv.Itoa(oldTable.ID))...),
		ctx)
}

func onSystemClientGoAway(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, systemOnClientGoAwayBodyLen); err != nil {
		return nil, err
	}

	client := username.UserName(event.Body[0])

	ctx.GoAway(event, client)

	return NewEvent(systemOnClientGoAwayStatus, event.Time, event.Body...), nil
}

func onSystemClientSitAtTable(event *Event, ctx *Context) (*Event, error) {
	if err := checkBodyLen(&event.Body, systemOnClientSitAtTableBodyLen); err != nil {
		return nil, err
	}

	client := ctx.Queue.Get()
	if client == "" {
		return nil, nil
	}

	tableID, err := strconv.Atoi(event.Body[1])
	if err != nil {
		return nil, err
	}
	table := (*ctx.Tables)[tableID-1]

	ctx.SitAtTable(event, client, table)

	return NewEvent(systemOnClientSitAtTableStatus, event.Time, event.Body...), nil
}

func checkBodyLen(body *[]string, wait int) error {
	if ln := len(*body); ln != wait {
		return newNoEnoughArgsError(wait, ln)
	}

	return nil
}
