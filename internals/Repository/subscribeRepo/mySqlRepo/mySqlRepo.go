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

func (s *sqlSubscribeRepo) CheckEmail(ctx context.Context, req *subscribeEntity.SubscribeReq)  ( *subscribeEntity.SubscribeRes, error){
	var res subscribeEntity.SubscribeRes
	stmt := fmt.Sprintf(`SELECT email FROM Subscribers WHERE email ='%v'`, req.Email)

	row :=  s.conn.QueryRow(stmt)

	if err := row.Scan(
        &res.Email,
    ); err != nil {
        return nil, err
    }
	return &res, nil
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
