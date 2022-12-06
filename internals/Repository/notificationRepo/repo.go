package notificationRepo

import (
	"test-va/internals/entity/notificationEntity"
)

type NotificationRepository interface {
	Persist(req *notificationEntity.CreateNotification) error
	GetTasksToExpireToday(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTasksToExpireInAFewHours(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTaskDetailsWhenDue(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetUserVaToken(userId string) (string, error)
}
