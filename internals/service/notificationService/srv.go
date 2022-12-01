package notificationService

import (
	"context"
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
	//GetTaskFromUser(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error)
}

type notificationSrv struct {
	app       *firebase.App
	repo      notificationRepo.NotificationRepository
	validator validationService.ValidationSrv
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
