package usecases

import (
	"context"
	"fmt"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type Buffer map[string]AggregateVotesInput

type AggregateVotes struct {
	logger   common.Logger
	database Database
}

func NewAggregateVotesUseCase(logger common.Logger, database Database) *AggregateVotes {
	return &AggregateVotes{
		logger,
		database,
	}
}

type AggregateVotesInput struct {
	Votes []entities.Vote
}

// votes are deduped in several ways:
//   - first we could have many votes for the same thread/comment, so we mark them as dirty based on votes for that id and only aggregate once
//   - secondly, there is no guarantee that the votes in the slice are in-order or that the vote we are processing is the latest vote for that thread/comment.
//     so we check the vote's updatedAt timestamp against the latest vote for that thread/comment in the db and only aggregate if the vote is newer
//
// TODO: Could do some batch processing in the db here
func (u *AggregateVotes) Execute(ctx context.Context, input AggregateVotesInput) {
	latestVotes := map[string]entities.Vote{}
	dirtyThreads := map[int64]bool{}
	dirtyComments := map[int64]bool{}
	for _, vote := range input.Votes {
		voteKey := fmt.Sprintf("%v:%v", vote.Type(), vote.Id())
		if latestVote, ok := latestVotes[voteKey]; !ok || (ok && latestVote.UpdatedAt() < vote.UpdatedAt()) {
			latestVotes[voteKey] = vote
		}
		switch vote.Type() {
		case common.ThreadVote:
			dirtyThreads[vote.Id()] = true
		case common.CommentVote:
			dirtyComments[vote.Id()] = true
		default:
			u.logger.Error(ctx).Msgf("error aggregating invalid vote type %v", vote.Type())
		}
	}

	if len(latestVotes) > 0 {
		u.logger.Info(ctx).Msgf("creating %v votes", len(latestVotes))
	}
	for _, vote := range latestVotes {
		if err := u.database.CreateVote(ctx, vote); err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error creating vote %v", vote)
		}
	}

	if len(dirtyThreads) > 0 {
		u.logger.Info(ctx).Msgf("updating %v thread votes", len(dirtyThreads))
	}
	for id := range dirtyThreads {
		if err := u.database.AggregateVotes(ctx, id, common.ThreadVote); err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error aggregating thread votes for %v", id)
		}
	}

	if len(dirtyComments) > 0 {
		u.logger.Info(ctx).Msgf("updating %v comment votes", len(dirtyComments))
	}
	for id := range dirtyComments {
		if err := u.database.AggregateVotes(ctx, id, common.CommentVote); err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error aggregating comment votes for %v", id)
		}
	}
}
