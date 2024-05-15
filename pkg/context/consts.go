package context

const (
	notOpenYetError       = "NotOpenYet"
	youShallNotPassError  = "YouShallNotPass"
	placeIsBusyError      = "PlaceIsBusy"
	clientUnknownError    = "ClientUnknown"
	iCanWaitNoLongerError = "ICanWaitNoLonger!"
)

const (
	onClientComeStatus       = "1"
	onClientSitAtTableStatus = "2"
	onClientWaitStatus       = "3"
	onClientGoAwayStatus     = "4"

	systemOnClientGoAwayStatus     = "11"
	systemOnClientSitAtTableStatus = "12"
	systemOnErrorStatus            = "13"
)

const (
	onClientComeBodyLen       = 1
	onClientSitAtTableBodyLen = 2
	onClientWaitBodyLen       = 1
	onClientGoAwayBodyLen     = 1

	systemOnClientGoAwayBodyLen     = 1
	systemOnClientSitAtTableBodyLen = 2
)
