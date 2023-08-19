// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: copyfrom.go

package bindings

import (
	"context"
)

// iteratorForInsertTransfers implements pgx.CopyFromSource.
type iteratorForInsertTransfers struct {
	rows                 []InsertTransfersParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertTransfers) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertTransfers) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].BlockNumber,
		r.rows[0].TransactionID,
		r.rows[0].LogIndex,
		r.rows[0].FromAddress,
		r.rows[0].ToAddress,
		r.rows[0].Amount,
	}, nil
}

func (r iteratorForInsertTransfers) Err() error {
	return nil
}

func (q *Queries) InsertTransfers(ctx context.Context, arg []InsertTransfersParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"transfers"}, []string{"block_number", "transaction_id", "log_index", "from_address", "to_address", "amount"}, &iteratorForInsertTransfers{rows: arg})
}
