package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/subscribeRepo"
	"test-va/internals/entity/subscribeEntity"
)

type sqlSubscribeRepo struct {
	conn *sql.DB
}
func NewMySqlSubscribeRepo(conn *sql.DB) subscribeRepo.SubscribeRepository {
	return &sqlSubscribeRepo{conn: conn}
}

func (s *sqlSubscribeRepo) PersistEmail(ctx context.Context, req *subscribeEntity.SubscribeReq)  error{
	stmt := fmt.Sprintf(`INSERT INTO Subscribers(
                 email
                  )
	VALUES ('%v')
	`, req.Email)
	_, err := s.conn.Exec(stmt)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
