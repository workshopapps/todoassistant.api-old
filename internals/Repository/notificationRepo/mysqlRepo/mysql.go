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

func (m *mySql) GetTaskDetailsWhenDue(userId string) (*notificationEntity.GetExpiredTasksWithDeviceId, error) {
	fmt.Sprintf(`SELECT 
    task_id,
	user_id, 
	title,
	description,
	end_time,
	device_id
	FROM Tasks join Notifications N on Tasks.user_id = N.user_id and N.user_id = '%s'
`, userId)
	panic("Impl me")
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

func (n *mySql) GetTasksToExpireToday() ([]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	stmt := fmt.Sprint(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notifications ON Tasks.user_id = Notifications.user_id
		WHERE CAST( Tasks.end_time as DATE ) = CAST( NOW() as DATE ) 
		AND Tasks.status = 'PENDING';
	`)

	var tasks []notificationEntity.GetExpiredTasksWithDeviceId
	query, err := n.conn.Query(stmt)
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
