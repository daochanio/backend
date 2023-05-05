package common

type Stream = string

const (
	VoteStream Stream = "votes"
)

type VoteMessage struct {
	Id      int64    `json:"id"`
	Address string   `json:"address"`
	Type    VoteType `json:"type"`
}
