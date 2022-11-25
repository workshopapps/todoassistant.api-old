package notificationRepo

import (
	"test-va/internals/entity/notificationEntity"
)

type NotificationRepository interface {
	Persist(req *notificationEntity.CreateNotification) error
}
