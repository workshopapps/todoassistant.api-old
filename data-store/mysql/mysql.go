package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type mySQLServer struct {
	conn *sql.DB
}

func NewMySQLServer(connStr string) (*mySQLServer, error) {
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &mySQLServer{conn: conn}, nil
}

func (m *mySQLServer) Ping() error {
	return m.conn.Ping()
}

func (m *mySQLServer) Close() error {
	return m.conn.Close()
}
