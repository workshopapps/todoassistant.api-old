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
	FROM Tasks join Notification_Tokens N on Tasks.user_id = N.user_id and N.user_id = '%s'
`, userId)
	panic("Impl me")
}

func NewMySqlNotificationRepo(conn *sql.DB) notificationRepo.NotificationRepository {
	return &mySql{conn: conn}
}

func (m *mySql) Persist(req *notificationEntity.CreateNotification) error {
	stmt := fmt.Sprintf(` 
		INSERT INTO Notification_Tokens(
        	notification_token_id,
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

func (n *mySql) getUserDetailsFromUserId(userId string) (string, error) {
	stmt := fmt.Sprintf(`
		SELECT first_name, last_name
		FROM Users 
		WHERE user_id = '%s'
		`, userId)
	var firstname string
	var lastname string

	err := n.conn.QueryRow(stmt).Scan(&firstname, &lastname)
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("%s  %s", firstname, lastname)
	return username, nil
}

func (n *mySql) GetUserVaToken(userId string) ([]string, string, string, error) {
	username, err := n.getUserDetailsFromUserId(userId)
	if err != nil {
		return nil, "", "", err
	}
	stmt := fmt.Sprintf(`
		SELECT virtual_assistant_id
		FROM Users
		WHERE user_id = '%s'
		`, userId)
	var vaId string
	err = n.conn.QueryRow(stmt).Scan(&vaId)
	if err != nil {
		return nil, "", "", err
	}

	stmt = fmt.Sprintf(`
		SELECT device_id
		FROM Users 
		INNER JOIN Notification_Tokens ON virtual_assistant_id = Notification_Tokens.user_id
		WHERE Users.user_id = '%s';
		`, userId)

	var deviceIds []string
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, "", "", err
	}
	for query.Next() {
		var deviceId string
		err = query.Scan(&deviceId)
		if err != nil {
			return nil, "", "", err
		}
		deviceIds = append(deviceIds, deviceId)
	}
	return deviceIds, vaId, username, err
}

func (n *mySql) GetUserToken(userId string) ([]string, string, error) {
	username, err := n.getUserDetailsFromUserId(userId)
	if err != nil {
		return nil, "", err
	}
	stmt := fmt.Sprintf(`
		SELECT device_id
		FROM Notification_Tokens 
        WHERE Notification_Tokens.user_id = '%s';
		`, userId)

	var deviceIds []string
	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, "", err
	}
	for query.Next() {
		var deviceId string
		err = query.Scan(&deviceId)
		if err != nil {
			return nil, "", err
		}
		deviceIds = append(deviceIds, deviceId)
	}
	return deviceIds, username, err
}

func (n *mySql) GetTasksToExpireToday(userClass string) (map[string][]notificationEntity.GetExpiredTasksWithDeviceId, error) {
	str := ""
	if userClass == "va" {
		str = "va_id"
	} else {
		str = "user_id"
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notification_Tokens ON Tasks.%s = Notification_Tokens.user_id
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
	if userClass == "va" {
		str = "va_id"
	} else {
		str = "user_id"
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id
		FROM Tasks
		INNER JOIN Notification_Tokens ON Tasks.%v = Notification_Tokens.user_id
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

func (n *mySql) CreateNotification(notificationId, userId, title, time, content, color, taskId string) error {
	stmt := fmt.Sprintf(`
		INSERT INTO Notifications(
			notification_id,
			user_id,
			title,
			time,
			content,
			color,
			notif_task_id
		) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v')`, notificationId, userId, title, time, content, color, taskId)

	_, err := n.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (n *mySql) DeleteNotifications(userId string) error {
	stmt := fmt.Sprintf(`
		DELETE FROM Notifications
		WHERE user_id = '%s'
		`, userId)

	_, err := n.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (n *mySql) GetNotifications(userId string) ([]notificationEntity.GetNotifcationsRes, error) {
	stmt := fmt.Sprintf(`
		SELECT notification_id, user_id, title, time, content, color, notif_task_id, read_status
		FROM Notifications
		WHERE user_id = '%s'`, userId)

	query, err := n.conn.Query(stmt)
	if err != nil {
		return nil, err
	}

	var notifs []notificationEntity.GetNotifcationsRes
	for query.Next() {
		var notif notificationEntity.GetNotifcationsRes
		err = query.Scan(&notif.NotificationId, &notif.UserId, &notif.Title, &notif.Time, &notif.Content, &notif.Color, &notif.TaskId, &notif.ReadStatus)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, notif)
	}
	return notifs, nil
}

func (n *mySql) UpdateNotification(notificationId string) error {
	stmt := fmt.Sprintf(`
		UPDATE Notifications SET read_status = '1' WHERE notification_id = '%s'`, notificationId)
	_, err := n.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
