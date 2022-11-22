package reminderService

import (
	"fmt"
	"log"
	"os"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
)

func Test_reminderSrv_SetReminder(t *testing.T) {

	taskId := "ccfe6ddf-a3c5-40ed-9d00-a1fce7369e82"
	dsn := os.Getenv("dsn")
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()
	repo := mySqlRepo.NewSqlRepo(conn)

	gcrn := gocron.NewScheduler(time.UTC)
	srv := NewReminderSrv(gcrn, conn, repo)
	due := time.Now().Add(2 * time.Minute).Format(time.RFC3339)
	log.Println(due)
	srv.SetReminder(due, taskId)

	time.Sleep(3 * time.Minute)
}

func Test_reminderSrv_SetReminderEveryXMin(t *testing.T) {
	dsn := os.Getenv("dsn")
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()
	conn.Ping()
	fmt.Println(time.Now().Format(time.RFC3339), time.Minute*2)

	//gcrn := gocron.NewScheduler(time.UTC)
	//srv := NewReminderSrv(gcrn, conn)
	//
	//srv.SetReminderEveryXMin(30)

}
