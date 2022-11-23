package reminderService

import (
	"database/sql"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/taskEntity"
	"time"
)

type ReminderSrv interface {
	SetReminder(dueDate, taskId string) error
	SetReminderEvery30Min()
	SetReminderEvery5Min()
}

type reminderSrv struct {
	cron *gocron.Scheduler
	conn *sql.DB
	repo taskRepo.TaskRepository
}

func (r *reminderSrv) SetReminderEvery5Min() {
	tasks, err := getPendingTasks(r.conn)
	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range tasks {
		//check if time until is <=30 minutes or 5 minutes
		yes := checkIfTimeElapsed5Minutes(task.EndTime)

		if yes {
			fmt.Println("notification sent")
			// send a notification
			continue
		}
	}
}

func (r *reminderSrv) SetReminderEvery30Min() {
	tasks, err := getPendingTasks(r.conn)
	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range tasks {
		//check if time until is <=30 minutes and greater than 5 minutes
		yes := checkIfTimeElapsed30Minutes(task.EndTime)

		if yes {
			fmt.Println("notification sent")
			// send a notification
			continue
		}
	}
}

func checkIfTimeElapsed30Minutes(endTime string) bool {
	dueTime, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		log.Println(err)
		return false
	}

	if dueTime.Before(time.Now()) {
		return false
	}

	minutes := time.Until(dueTime).Minutes()
	if minutes <= float64(30) && minutes > float64(5) {
		return true
	} else {
		return false
	}
}

func checkIfTimeElapsed5Minutes(due string) bool {

	dueTime, err := time.Parse(time.RFC3339, due)
	if err != nil {
		log.Println(err)
		return false
	}

	if dueTime.Before(time.Now()) {
		return false
	}

	minutes := time.Until(dueTime).Minutes()
	if minutes <= float64(5) {
		return true
	} else {
		return false
	}
}

func (r *reminderSrv) SetReminder(dueDate, taskId string) error {
	s := gocron.NewScheduler(time.UTC)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return err
	}
	s.Every(1).StartAt(dDate).Do(func() {
		log.Println("setting status to expired")
		r.repo.SetTaskToExpired(taskId)
	})

	s.LimitRunsTo(1)
	s.StartAsync()
	return nil
}

func NewReminderSrv(s *gocron.Scheduler, conn *sql.DB, taskrepo taskRepo.TaskRepository) ReminderSrv {
	return &reminderSrv{cron: s, conn: conn, repo: taskrepo}
}

func getPendingTasks(conn *sql.DB) ([]taskEntity.GetPendingTasks, error) {
	stmt := fmt.Sprint(`
		SELECT task_id, user_id, title,description, end_time
		FROM Tasks
		WHERE status = 'PENDING';
`)
	var tasks []taskEntity.GetPendingTasks
	query, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task taskEntity.GetPendingTasks
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
