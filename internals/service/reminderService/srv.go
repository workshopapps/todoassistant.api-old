package reminderService

import (
	"database/sql"
	"errors"
	"encoding/json"
	"fmt"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/notificationService"
	"time"

	"github.com/go-co-op/gocron"
)

type ReminderSrv interface {
	SetReminder(dueDate, taskId string) error
	SetReminderEvery30Min()
	SetReminderEvery5Min()
	SetDailyReminder(data *taskEntity.CreateTaskReq) error
	SetWeeklyReminder(data *taskEntity.CreateTaskReq) error
	SetBiWeeklyReminder(data *taskEntity.CreateTaskReq) error
	SetMonthlyReminder(data *taskEntity.CreateTaskReq) error
	SetYearlyReminder(data *taskEntity.CreateTaskReq) error
	ScheduleNotificationEverySixHours()
	ScheduleNotificationDaily()
}

type reminderSrv struct {
	cron *gocron.Scheduler
	conn *sql.DB
	repo taskRepo.TaskRepository
	nSrv notificationService.NotificationSrv
}

func (r *reminderSrv) SetBiWeeklyReminder(data *taskEntity.CreateTaskReq) error {
	s := gocron.NewScheduler(time.UTC)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	s.Every(14).Weeks().StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 14).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.SetNewEvent(data)
		if err != nil {
			return err
		}
		return nil
	})
	s.StartAsync()
	log.Println("created new event.")
	return nil
}

func (r *reminderSrv) SetYearlyReminder(data *taskEntity.CreateTaskReq) error {
	s := gocron.NewScheduler(time.UTC)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	s.Every(12).Months().StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(1, 0, 0).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.SetNewEvent(data)
		if err != nil {
			return err
		}
		return nil
	})
	s.StartAsync()
	log.Println("created new event.")
	return nil
}

func (r *reminderSrv) SetMonthlyReminder(data *taskEntity.CreateTaskReq) error {
	s := gocron.NewScheduler(time.UTC)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	s.Every(1).Months().StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 1).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.SetNewEvent(data)
		if err != nil {
			return err
		}
		s.StartAsync()
		return nil
	})
	log.Println("created new event.")
	return nil
}

func (r *reminderSrv) SetWeeklyReminder(data *taskEntity.CreateTaskReq) error {
	s := gocron.NewScheduler(time.UTC)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	s.Every(7).Day().StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 7).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.SetNewEvent(data)
		if err != nil {
			return err
		}
		s.StartAsync()
		return nil
	})
	log.Println("created new event.")
	return nil
}

func (r *reminderSrv) SetDailyReminder(data *taskEntity.CreateTaskReq) error {
	s := gocron.NewScheduler(time.UTC)

	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	if dDate.Before(time.Now()) {
		return errors.New("invalid Time, try again")
	}

	s.Every(1).Day().StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		log.Printf("\n")
		r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 1).Format(time.RFC3339)
		data.Status = "PENDING"

		err = r.repo.SetNewEvent(data)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("created new event.")

		//// send notification
		//task, err := r.noSrv.GetTaskFromUser(data.UserId)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println("notification sent out")
		//r.noSrv.SendNotification(task.DeviceId,
		//	"Your Notification is about to expire",
		//	"your Task is due in 5 miutes",
		//	[]string{task.TaskId},
		//)

		return nil
	})
	s.StartAsync()
	return nil
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
			fmt.Println("notification sent out")
			r.nSrv.SendNotification(task.DeviceId,
				"Your Notification is about to expire",
				"your Task is due in 5 minutes",
				task.TaskId,
			)
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
			fmt.Println("notification sent out")
			r.nSrv.SendNotification(task.DeviceId,
				"Your Notification is about to expire",
				"your Task is due in 30 miutes",
				task.TaskId,
			)
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

// Everyday By 12:00am you get Notifications For All Tasks That are Due That Day
func (r *reminderSrv) ScheduleNotificationDaily() {
	fmt.Println("Daily Notifications Setup")
	r.cron.Every(1).Days().At("00:00").Do(func() {
		tasks, err := r.nSrv.GetTasksToExpireToday()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Daily")

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}
		for k, v := range tasks {
			tasksToString, err := json.Marshal(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = r.nSrv.SendNotification(k, "Due Today", "Rise and Shine, These Are Due Today", string(tasksToString))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
}

// Your Pending Tasks are Checked On Six Hour Intervals to Get Tasks That Are Just About To Expire
func (r *reminderSrv) ScheduleNotificationEverySixHours() {
	fmt.Println("Six Hour Notifications Setup")
	r.cron.Every(6).Hours().Do(func() {
		tasks, err := r.nSrv.GetTasksToExpireInAFewHours()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Six Hourly")

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}

		for k, v := range tasks {
			tasksToString, err := json.Marshal(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = r.nSrv.SendNotification(k, "Due Shortly", "Warning: Few Hours to Go", string(tasksToString))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
}


func getPendingTasks(conn *sql.DB) ([]taskEntity.GetPendingTasks, error) {
	stmt := fmt.Sprint(`
		SELECT T.task_id, T.user_id, T.title,T.description, T.end_time, N.device_id
		FROM Tasks T join Notifications N on T.user_id = N.user_id
		WHERE status = 'PENDING';
	`)
	var tasks []taskEntity.GetPendingTasks
	query, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task taskEntity.GetPendingTasks
		var deviceId string
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func getExpiredTasks(conn *sql.DB) ([]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	stmt := fmt.Sprint(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id, va_id
		FROM Tasks
		INNER JOIN Notifications ON Tasks.user_id = Notifications.user_id
		WHERE Tasks.status = 'EXPIRED';
	`)

	var tasks []notificationEntity.GetExpiredTasksWithDeviceId

	query, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task notificationEntity.GetExpiredTasksWithDeviceId
		var deviceId string
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func NewReminderSrv(s *gocron.Scheduler, conn *sql.DB, taskrepo taskRepo.TaskRepository, nSrv notificationService.NotificationSrv) ReminderSrv {
	return &reminderSrv{cron: s, conn: conn, repo: taskrepo, nSrv: nSrv}
}
