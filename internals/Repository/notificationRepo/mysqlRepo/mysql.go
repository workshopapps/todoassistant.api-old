package mySqlRepo

import (
	"database/sql"
	"fmt"
	"strings"
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
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}
		return err
	}
	return nil
}

func (n *mySql) GetUserVaToken(userId string) ([]string, error) {
	stmt := fmt.Sprintf(`
		SELECT device_id 
		FROM Users 
		INNER JOIN Notifications ON virtual_assistant_id = Notifications.user_id
		WHERE Users.user_id = '%s';
		`, userId)

	var deviceIds []string
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var deviceId string
		err = query.Scan(&deviceId)
		if err != nil {
			return nil, err
		}
		deviceIds = append(deviceIds, deviceId)
	}
	return deviceIds, err
}

func (n *mySql) GetUserToken(userId string) ([]string, error) {
	stmt := fmt.Sprintf(`
		SELECT * 
		FROM Notifications 
        WHERE Notifications.user_id = '%s';
		`, userId)

	var deviceIds []string
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var deviceId string
		err = query.Scan(&deviceId)
		if err != nil {
			return nil, err
		}
		deviceIds = append(deviceIds, deviceId)
	}
	return deviceIds, err
}

func (n *mySql) GetTasksToExpireToday(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	str := ""
	if (userClass == "va") {
		str = "va_id"
	} else {
		str = "user_id"
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notifications ON Tasks.%s = Notifications.user_id
		WHERE CAST( Tasks.end_time as DATE ) = CAST( NOW() as DATE )
		AND Tasks.status = 'PENDING';
	`, str)

	taskMap := make(map[string][]notificationEntity.GetExpiredTasksWithDeviceId)
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task notificationEntity.GetExpiredTasksWithDeviceId
		var deviceId string
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
		if err != nil {
			return nil, err
		}

		if mTask, ok := taskMap[deviceId]; !ok {
			taskMap[deviceId] = []notificationEntity.GetExpiredTasksWithDeviceId{task}
		} else {
			taskMap[deviceId] = append(mTask, task)
		}
	}
	return taskMap, nil
}

func (n *mySql) GetTasksToExpireInAFewHours(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	str := ""
	if (userClass == "va") {
		str = "va_id"
	} else {
		str = "user_id"
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notifications ON Tasks.%v = Notifications.user_id
		WHERE CAST( Tasks.end_time as DATE ) = CAST( NOW() as DATE ) 
		AND  
		CAST(Tasks.end_time as TIME ) < CAST(NOW() + INTERVAL 7 HOUR as TIME )
		AND Tasks.status = 'PENDING';
	`, str)

	taskMap := make(map[string][]notificationEntity.GetExpiredTasksWithDeviceId)
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task notificationEntity.GetExpiredTasksWithDeviceId
		var deviceId string
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
		if err != nil {
			return nil, err
		}

		if mTask, ok := taskMap[deviceId]; !ok {
			taskMap[deviceId] = []notificationEntity.GetExpiredTasksWithDeviceId{task}
		} else {
			taskMap[deviceId] = append(mTask, task)
		}
	}
	return taskMap, nil
}
