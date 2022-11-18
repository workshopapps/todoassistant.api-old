package reminderService

import (
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"test-va/internals/data-store/mysql"
	"testing"
	"time"
)

func Test_reminderSrv_SetReminder(t *testing.T) {
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

	gcrn := gocron.NewScheduler(time.UTC)
	srv := NewReminderSrv(gcrn, conn)
	due := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	log.Println(due)
	srv.SetReminder(due)
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

	gcrn := gocron.NewScheduler(time.UTC)
	srv := NewReminderSrv(gcrn, conn)

	srv.SetReminderEveryXMin(30)

}
