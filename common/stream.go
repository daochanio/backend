package common

type Stream = string

const (
	SigninStream Stream = "signin"
	VoteStream   Stream = "vote"
)

type VoteMessage struct {
	Id      int64    `json:"id"`
	Address string   `json:"address"`
	Type    VoteType `json:"type"`
}

type SigninMessage struct {
	Address string `json:"address"`
}
