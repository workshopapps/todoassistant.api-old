package mySqlCallRepo


import (
	"context"
	"database/sql"
	"test-va/internals/Repository/callRepo"
	"test-va/internals/entity/callEntity"
)

type sqlCallRepo struct {
	conn *sql.DB
}

func NewSqlCallRepo(conn *sql.DB) callRepo.CallRepository {
	return &sqlCallRepo{conn: conn}
}

func (s *sqlCallRepo) GetCalls(ctx context.Context) ([]*callEntity.CallRes, error) {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := "SELECT id, va_id, user_id, call_rating, call_comment FROM Calls"

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calls := []*callEntity.CallRes{}

	for rows.Next() {
		var call callEntity.CallRes
		err := rows.Scan(
			&call.CallId,
			&call.VaId,
			&call.UserId,
			&call.CallRating,
			&call.CallComment,
		)
		if err != nil {
			return nil, err
		}
		calls = append(calls, &call)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return calls, nil
}