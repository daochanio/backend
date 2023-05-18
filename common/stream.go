package common

type Stream = string

const (
	SigninStream Stream = "signin"
	VoteStream   Stream = "vote"
)

type VoteMessage struct {
	Id      int64     `json:"id"`
	Address string    `json:"address"`
	Type    VoteType  `json:"type"`
	Value   VoteValue `json:"value"`
	// the timestamp at which the vote was accepted
	// this is used to discard old votes that may be processed out of order
	UpdatedAt int64 `json:"updated_at"`
}

type SigninMessage struct {
	Address string `json:"address"`
}
