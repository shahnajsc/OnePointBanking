package repo

import (
	"context"
	"database/sql"
)

type OPConnectRepo struct {
	db *sql.DB
}

func NewOPConnectRepo(db *sql.DB) *OPConnectRepo {
	return &OPConnectRepo{db: db}
}

func (r *OPConnectRepo) SavePending(ctx context.Context, state, userID, authorizationID, nonce string) error {
	const q = `
		INSERT INTO op_authorizations (state, user_id, authorization_id, nonce)
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.ExecContext(ctx, q, state, userID, authorizationID, nonce)
	return err
}
