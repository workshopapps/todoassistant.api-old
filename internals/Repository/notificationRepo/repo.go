package notificationRepo

import (
	"test-va/internals/entity/notificationEntity"
)

type NotificationRepository interface {
	Persist(req *notificationEntity.CreateNotification) error
	GetTasksToExpireToday() ([]notificationEntity.GetExpiredTasksWithDeviceId, error)
	GetTaskDetailsWhenDue(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error)
}
