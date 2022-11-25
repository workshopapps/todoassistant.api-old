package mySqlRepo

import (
	"database/sql"
	"fmt"
	"test-va/internals/Repository/notificationRepo"
	"test-va/internals/entity/notificationEntity"
)

type mySql struct {
	conn *sql.DB
}

func NewMySqlNotificationRepo(conn *sql.DB) notificationRepo.NotificationRepository {
	return &mySql{conn: conn}
}

func (m *mySql) Persist(req *notificationEntity.CreateNotification) error {
	stmt := fmt.Sprintf(` 
		INSERT INTO Notifications(
        	notification_id,
        	user_id,
			device_id
            ) VALUES ('%v', '%v', '%v')`, req.NotificationId, req.UserId, req.DeviceId)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
