package postgres

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
	"github.com/daochanio/backend/db/bindings"
)

func (g *postgresGateway) CreateVote(ctx context.Context, vote entities.Vote) error {
	switch vote.Type() {
	case common.ThreadVote:
		switch vote.Value() {
		case common.Upvote:
			return g.queries.CreateThreadUpVote(ctx, bindings.CreateThreadUpVoteParams{
				ThreadID: vote.Id(),
				Address:  vote.Address(),
			})
		case common.Downvote:
			return g.queries.CreateThreadDownVote(ctx, bindings.CreateThreadDownVoteParams{
				ThreadID: vote.Id(),
				Address:  vote.Address(),
			})
		case common.Unvote:
			return g.queries.CreateThreadUnVote(ctx, bindings.CreateThreadUnVoteParams{
				ThreadID: vote.Id(),
				Address:  vote.Address(),
			})
		default:
			return fmt.Errorf("invalid vote value %v", vote.Value())
		}
	case common.CommentVote:
		switch vote.Value() {
		case common.Upvote:
			return g.queries.CreateCommentUpVote(ctx, bindings.CreateCommentUpVoteParams{
				CommentID: vote.Id(),
				Address:   vote.Address(),
			})
		case common.Downvote:
			return g.queries.CreateCommentDownVote(ctx, bindings.CreateCommentDownVoteParams{
				CommentID: vote.Id(),
				Address:   vote.Address(),
			})
		case common.Unvote:
			return g.queries.CreateCommentUnVote(ctx, bindings.CreateCommentUnVoteParams{
				CommentID: vote.Id(),
				Address:   vote.Address(),
			})
		default:
			return fmt.Errorf("invalid vote value %v", vote.Value())
		}
	default:
		return fmt.Errorf("invalid vote type: %v", vote.Type())
	}
}

func (g *postgresGateway) AggregateVotes(ctx context.Context, id int64, voteType common.VoteType) error {
	switch voteType {
	case common.ThreadVote:
		return g.queries.AggregateThreadVotes(ctx, id)
	case common.CommentVote:
		return g.queries.AggregateCommentVotes(ctx, id)
	default:
		return fmt.Errorf("invalid vote type: %v", voteType)
	}
}
