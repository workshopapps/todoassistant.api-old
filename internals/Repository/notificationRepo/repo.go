package notificationRepo

import (
	"test-va/internals/entity/notificationEntity"
)

type NotificationRepository interface {
	Persist(req *notificationEntity.CreateNotification) error
	GetTasksToExpireToday(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTasksToExpireInAFewHours(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTaskDetailsWhenDue(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetUserVaToken(userId string) ([]string, string, string, error)
	GetUserToken(userId string) ([]string, string, error)
	CreateNotification(notificationId, userId, title, time, content, color, taskId string) error
	GetNotifications(userId string) ([]notificationEntity.GetNotifcationsRes, error)
	DeleteNotifications(userId string) error
	UpdateNotification(notificationId string) error
}
