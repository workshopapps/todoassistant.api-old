package notificationEntity

type GetExpiredTasksWithDeviceId struct {
	TaskId      string `json:"task_id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndTime     string `json:"end_time"`
	DeviceId    string `json:"device_id"`
}

type CreateNotification struct {
	NotificationId string `json:"notification_id"`
	UserId         string `json:"user_id"`
	DeviceId       string `json:"device_id"`
}

type NotificationRegisterReq struct {
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
}

type NotificationRegisterRes struct {
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
	Message  string `json:"message"`
}
