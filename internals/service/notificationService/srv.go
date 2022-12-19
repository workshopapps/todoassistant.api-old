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
	SendNotification(token, title string, body, data interface{}) error
	SendBatchNotifications(tokens []string, title string, body, data interface{}) error
	SendVaNotification(token, title, body string, taskId string) error
	GetUserVaToken(userId string) ([]string, string, string, error)
	GetUserToken(userId string) ([]string, string, error)
	GetNotifications(userId string) ([]notificationEntity.GetNotifcationsRes, *ResponseEntity.ServiceError)
	GetTasksToExpireToday() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTasksToExpireInAFewHours() (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	CreateNotification(userId, title, time, content, color, taskId string) error
	DeleteNotifications(userId string) error
	UpdateNotification(notificationId string) error
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

func (n notificationSrv) SendBatchNotifications(tokens []string, title string, body, data interface{}) error {
	ctx := context.Background()
	if n.app == nil {
		return fmt.Errorf("could not initialize firebase app")
	}
	fmcClient, err := n.app.Messaging(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	bodyToString, err := json.Marshal(body)
	if err != nil {
		return err
	}

	dataToString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := fmcClient.SendMulticast(ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  string(bodyToString),
		},
		Data: map[string]string{
			"tasks": string(dataToString),
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

func (n notificationSrv) SendNotification(token, title string, body, data interface{}) error {
	ctx := context.Background()
	if n.app == nil {
		return fmt.Errorf("could not initialize firebase app")
	}
	fmcClient, err := n.app.Messaging(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	bodyToString, err := json.Marshal(body)
	if err != nil {
		return err
	}

	dataToString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := fmcClient.Send(ctx, &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  string(bodyToString),
		},
		Data: map[string]string{
			"tasks": string(dataToString),
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

// Register for notifications godoc
// @Summary	Register to get notifications service
// @Description	Notification registration route
// @Tags	Notifications
// @Accept	json
// @Produce	json
// @Param	request	body	notificationEntity.CreateNotification	true "Notifications Details"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/notification [post]
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

func (n notificationSrv) GetUserVaToken(userId string) ([]string, string, string, error) {
	return n.repo.GetUserVaToken(userId)
}

func (n notificationSrv) GetUserToken(userId string) ([]string, string, error) {
	return n.repo.GetUserToken(userId)
}

// Delete notifications godoc
// @Summary	Delete notifications service
// @Description	Delete notification route
// @Tags	Notifications
// @Accept	json
// @Produce	json
// @Success	200  {string}  string    "Deleted Successful"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/notification [delete]
func (n notificationSrv) DeleteNotifications(userId string) error {
	return n.repo.DeleteNotifications(userId)
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

func (n notificationSrv) CreateNotification(userId, title, time, content, color, taskId string) error {
	notificationId := uuid.New().String()
	return n.repo.CreateNotification(notificationId, userId, title, time, content, color, taskId)
}

// Get notifications godoc
// @Summary	Get all your notifications
// @Description	Get notification route
// @Tags	Notifications
// @Accept	json
// @Produce	json
// @Success	200  {object}	[]notificationEntity.GetNotifcationsRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/notification [get]
func (n notificationSrv) GetNotifications(userId string) ([]notificationEntity.GetNotifcationsRes, *ResponseEntity.ServiceError) {
	notifications, err := n.repo.GetNotifications(userId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return notifications, nil
}

// Update notification godoc
// @Summary	Update a specific notification
// @Description	Update notification route
// @Tags	Notifications
// @Accept	json
// @Produce	json
// @Param	notificationId	path	string	true	"Notification Id"
// @Success	200  {string}	string	"Updated Successful"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/notification/{notificationId} [patch]
func (n notificationSrv) UpdateNotification(notificationId string) error {
	return n.repo.UpdateNotification(notificationId)
}
