package reminderService

import (
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/taskEntity"
	"time"

	"github.com/go-co-op/gocron"
)

type ReminderSrv interface {
	SetReminder(dueDate, taskId string) error
	SetReminderEveryXMin(x int)
}

type reminderSrv struct {
	cron *gocron.Scheduler
	conn *sql.DB
	repo taskRepo.TaskRepository
}

func (r *reminderSrv) SetReminderEveryXMin(x int) {
	tasks, err := getPendingTasks(r.conn)
	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range tasks {
		//check if time until is <=30 minutes
		yes := checkIfTimeElapsedXMinutes(task.EndTime, x)
		if yes {
			fmt.Println("notification sent.....")
			continue
		}
	}
}

func (r *reminderSrv) SetReminder(dueDate, taskId string) error {

	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return err
	}

	// find time till time is expired
	fmt.Println(dDate)

	duration := time.Until(dDate)

	// convert to minutes
	minutes := duration.Minutes()
	ss := fmt.Sprintf("%vm", minutes)
	log.Println(ss)

	r.cron.Every(2).Minutes().Do(func() {
		log.Println("Doing... set task status to expired")
		r.repo.SetTaskToExpired(taskId)
	})

	r.cron.LimitRunsTo(1)
	r.cron.StartAsync()

	return nil
}

func NewReminderSrv(scheduler *gocron.Scheduler, conn *sql.DB, taskrepo taskRepo.TaskRepository) ReminderSrv {
	return &reminderSrv{cron: scheduler, conn: conn, repo: taskrepo}
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

func checkIfTimeElapsedXMinutes(due string, x int) bool {
	dueTime, err := time.Parse(time.RFC3339, due)
	if err != nil {
		log.Println(err)
		return false
	}

	minutes := time.Until(dueTime).Minutes()
	log.Println(minutes)
	if minutes <= float64(x) {
		return true
	} else {
		return false
	}
}
