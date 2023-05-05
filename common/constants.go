package common

type VoteValue string

const (
	Upvote   VoteValue = "upvote"
	Downvote VoteValue = "downvote"
	Unvote   VoteValue = "unvote"
)

type VoteType string

const (
	ThreadVote  VoteType = "thread"
	CommentVote VoteType = "comment"
)
