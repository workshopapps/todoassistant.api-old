package mysql

import (
	"log"
	"testing"
)

func TestNewMySQLServer(t *testing.T) {
	conn, err := NewMySQLServer("hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todo")
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}

	err = conn.conn.Ping()
	if err != nil {
		return
	}
}
