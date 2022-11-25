package notificationService

import (
	"context"
	"database/sql"
	"fmt"
	"test-va/internals/Repository/notificationRepo"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/service/validationService"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

type NotificationSrv interface {
	RegisterForNotifications(req *notificationEntity.CreateNotification) *ResponseEntity.ServiceError
}

type notificationSrv struct {
	app       *firebase.App
	cron      *gocron.Scheduler
	conn      *sql.DB
	taskRepo  taskRepo.TaskRepository
	repo      notificationRepo.NotificationRepository
	validator validationService.ValidationSrv
}

func New(app *firebase.App, cron *gocron.Scheduler, conn *sql.DB, repo notificationRepo.NotificationRepository, taskRepo taskRepo.TaskRepository, validator validationService.ValidationSrv) notificationSrv {
	return notificationSrv{
		app:       app,
		cron:      cron,
		conn:      conn,
		repo:      repo,
		taskRepo:  taskRepo,
		validator: validator,
	}
}

func (n *notificationSrv) ScheduleNotificationDaily() {
	n.cron.Every(24).Hours().Do(func() {
		tasks, err := getExpiredTasks(n.conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}
		for _, v := range tasks {
			err := sendToToken(n.app, v.DeviceId, "Expired Tasks", v.Description)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
}

func (n *notificationSrv) ScheduleNotificationEverySixHours() {
	n.cron.Every(5).Minutes().Do(func() {
		tasks, err := getTasksToExpireToday(n.conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}

		for _, v := range tasks {
			err := sendToToken(n.app, v.DeviceId, "Pending Tasks", v.Description)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
}

func sendToToken(app *firebase.App, token, title, body string) error {
	ctx := context.Background()
	fmcClient, err := app.Messaging(ctx)

	if err != nil {
		fmt.Println(err)
		return err
	}

	response, err := fmcClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully Sent Message", response)
	return nil
}

func getExpiredTasks(conn *sql.DB) ([]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	stmt := fmt.Sprint(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
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
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &task.DeviceId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func getTasksToExpireToday(conn *sql.DB) ([]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	stmt := fmt.Sprint(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notifications ON Tasks.user_id = Notifications.user_id
		WHERE CAST( Tasks.end_time as DATE ) = CAST( NOW() as DATE ) 
		AND Tasks.status = 'PENDING';
	`)

	var tasks []notificationEntity.GetExpiredTasksWithDeviceId

	query, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task notificationEntity.GetExpiredTasksWithDeviceId
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &task.DeviceId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (n notificationSrv) RegisterForNotifications(req *notificationEntity.CreateNotification) *ResponseEntity.ServiceError {
	err := n.validator.Validate(req)
	if err != nil {
		return ResponseEntity.NewValidatingError(err)
	}
	req.NotificationId = uuid.New().String()
	err = n.repo.Persist(req)
	if err != nil {
		return ResponseEntity.NewInternalServiceError(fmt.Sprintf("Unable to Register Notification Details %v", err))
	}
	return nil
}
