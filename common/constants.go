package common

type VoteType string

const (
	UpVote   VoteType = "up"
	DownVote VoteType = "down"
	UnVote   VoteType = "undo"
)
