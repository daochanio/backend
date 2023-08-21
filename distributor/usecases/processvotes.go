package usecases

import "context"

type ProcessVotes struct {
}

func NewProcessVotes() *ProcessVotes {
	return &ProcessVotes{}
}

// Make decisions on whether the vote should be counted towards a distribution or not and create a vote record for it.
// I.e if the vote is on a comment/thread that is older than a certain cuttoff.
// We always record a record regardless of whether it is counted or not with some kind of accepted/discarded flag
func (u *ProcessVotes) Execute(ctx context.Context) error {
	return nil
}
