package notificationService

import (
	"context"
	"encoding/json"
	"fmt"
	"test-va/internals/Repository/notificationRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/service/validationService"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
)

type NotificationSrv interface {
	RegisterForNotifications(req *notificationEntity.CreateNotification) *ResponseEntity.ServiceError
	SendNotification(token, title, body string, taskId string) error
	SendVaNotification(token, title, body string, taskId string) error
	GetUserVaToken(userId string) (string, error)
	SendNotificationToVA(userId, topic, body string, data interface{})
	GetTasksToExpireToday() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTasksToExpireInAFewHours() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	//GetTaskFromUser(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error)
}

type notificationSrv struct {
	app       *firebase.App
	repo      notificationRepo.NotificationRepository
	validator validationService.ValidationSrv
}

func (n notificationSrv) SendVaNotification(token, title, body string, taskId string) error {
	//TODO implement me
	panic("implement me")
}

//func (n notificationSrv) GetTaskFromUser(userId string)  (*notificationEntity.GetExpiredTasksWithDeviceId, error) {
//	task, err := n.repo.GetTaskDetailsWhenDue(userId)
//	if err != nil {
//		return nil, err
//	}
//	return task, nil
//}

func New(app *firebase.App, repo notificationRepo.NotificationRepository,
	validator validationService.ValidationSrv) NotificationSrv {
	return notificationSrv{

		app:       app,
		repo:      repo,
		validator: validator,
	}
}

func (n notificationSrv) SendNotification(token, title, body string, taskId string) error {
	ctx := context.Background()
	fmcClient, err := n.app.Messaging(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// taskIdsToString, err := json.Marshal(taskIds)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	response, err := fmcClient.Send(ctx, &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"tasks": taskId,
		},
		Webpush: &messaging.WebpushConfig{
			Headers: map[string]string{
				"Urgency": "high",
				"TTL":     "5000",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully Sent Message", response)
	return nil
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

func (n notificationSrv) GetUserVaToken(userId string) (string, error) {
	return n.repo.GetUserVaToken(userId)
}

func (n notificationSrv) SendNotificationToVA(userId, topic, body string, data interface{}) {
	vaDeviceId, err := n.GetUserVaToken(userId)
	if err != nil {
		fmt.Println("Error Getting VA DeviceId For Notifications", err)
	}
	dataToString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error Marshalling in Notifications", err)
	}
	n.SendNotification(vaDeviceId, "Task Created", fmt.Sprintf("%s Just Created a Task", userId), string(dataToString))

}

func (n notificationSrv) GetTasksToExpireToday() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	// Select All The Users with Pending Tasks and Send Notifications to Them
	userTaskMap, err := n.repo.GetTasksToExpireToday("user")
	if err != nil {
		return nil, err
	}

	vaTaskMap, err := n.repo.GetTasksToExpireToday("va")
	if err != nil {
		return nil, err
	}

	taskMap := make(map[string][]notificationEntity.GetExpiredTasksWithDeviceId)

	for k, v := range userTaskMap {
		taskMap[k] = v
	}

	for k, v := range vaTaskMap {
		if _, ok := taskMap[k]; !ok {
			taskMap[k] = v
		} else {
			taskMap[k] = append(taskMap[k], v...)
		}
	}

	return taskMap, nil
}

func (n notificationSrv) GetTasksToExpireInAFewHours() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	// Select All The Users with Pending Tasks and Send Notifications to Them
	userTaskMap, err := n.repo.GetTasksToExpireInAFewHours("user")
	if err != nil {
		return nil, err
	}
	
	// Select All The VA with Users that Have Pending Tasks and Send Notifications to Them
	vaTaskMap, err := n.repo.GetTasksToExpireInAFewHours("va")
	if err != nil {
		return nil, err
	}

	taskMap := make(map[string][]notificationEntity.GetExpiredTasksWithDeviceId)

	for k, v := range userTaskMap {
		taskMap[k] = v
	}

	for k, v := range vaTaskMap {
		if _, ok := taskMap[k]; !ok {
			taskMap[k] = v
		} else {
			taskMap[k] = append(taskMap[k], v...)
		}
	}

	return taskMap, nil
}

