package usecases

import "context"

type ProcessVote struct {
}

func NewProcessVote() *ProcessVote {
	return &ProcessVote{}
}

type ProcessVoteInput struct {
}

// Make decisions on whether the vote should be counted towards a distribution or not and create a vote record for it.
// We always record a record regardless of whether it is counted or not with some kind of accepted/discarded flag
// hydrate the vote along with its thread/comment
// accept or reject the vote based on criteria
// - address must have a certain amount of reputation to have their votes counted
// - if the vote is on a comment/thread that is older than a certain cuttoff, it is rejected
// - if the vote is from the same address as the comment/thread author, it is rejected
// - if the vote is from a blacklisted address, it is rejected
// - if the vote is on a deleted thread/comment, it is rejected
func (u *ProcessVote) Execute(ctx context.Context) error {
	return nil
}
