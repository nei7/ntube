// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: email_verify.sql

package data

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const CreateVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id, email, secret_code)
VALUES ($1, $2, $3)
RETURNING id, email, user_id, secret_code, is_used, created_at, expired_at
`

type CreateVerifyEmailParams struct {
	UserID     pgtype.UUID
	Email      string
	SecretCode string
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, CreateVerifyEmail, arg.UserID, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.UserID,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}
