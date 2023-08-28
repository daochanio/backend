package postgres

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/gateways/postgres/bindings"
)

func (g *postgresGateway) CreateVote(ctx context.Context, vote entities.Vote) error {
	switch vote.Type() {
	case entities.ThreadVote:
		switch vote.Value() {
		case entities.Upvote:
			return g.queries.CreateThreadUpVote(ctx, bindings.CreateThreadUpVoteParams{
				ThreadID:    vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		case entities.Downvote:
			return g.queries.CreateThreadDownVote(ctx, bindings.CreateThreadDownVoteParams{
				ThreadID:    vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		case entities.Unvote:
			return g.queries.CreateThreadUnVote(ctx, bindings.CreateThreadUnVoteParams{
				ThreadID:    vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		default:
			return fmt.Errorf("invalid vote value %v", vote.Value())
		}
	case entities.CommentVote:
		switch vote.Value() {
		case entities.Upvote:
			return g.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
				CommentID:   vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		case entities.Downvote:
			return g.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
				CommentID:   vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		case entities.Unvote:
			return g.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
				CommentID:   vote.Id(),
				Address:     vote.Address(),
				ToTimestamp: float64(vote.UpdatedAt()),
			})
		default:
			return fmt.Errorf("invalid vote value %v", vote.Value())
		}
	default:
		return fmt.Errorf("invalid vote type: %v", vote.Type())
	}
}

func (g *postgresGateway) AggregateVotes(ctx context.Context, id int64, voteType entities.VoteType) error {
	switch voteType {
	case entities.ThreadVote:
		return g.queries.AggregateThreadVotes(ctx, id)
	case entities.CommentVote:
		return g.queries.AggregateCommentVotes(ctx, id)
	default:
		return fmt.Errorf("invalid vote type: %v", voteType)
	}
}
