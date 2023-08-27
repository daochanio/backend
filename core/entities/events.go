package entities

type Events struct {
	transfers []Transfer
}

func NewEvents(transfers []Transfer) Events {
	return Events{
		transfers: transfers,
	}
}

func (e Events) Transfers() []Transfer {
	return e.transfers
}
