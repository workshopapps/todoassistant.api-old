package taskEntity

type TaskFile struct {
	FileLink string `json:"file_link"`
	FileType string `json:"file_type"`
}

type CreateTaskReq struct {
	TaskId      string     `json:"task_id"`
	UserId      string     `json:"user_id"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	Repeat      string     `json:"repeat"`
	Files       []TaskFile `json:"files"`
	StartTime   string     `json:"start_time" validate:"required"`
	EndTime     string     `json:"end_time" validate:"required"`
	VAOption    string     `json:"va_option"`
	Status      string     `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

type CreateTaskRes struct {
	TaskId      string `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	VAOption    string `json:"va_option"`
	Repeat      string `json:"repeat"`
}

type GetPendingTasksRes struct {
	TaskId      string `json:"task_id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
}

type GetPendingTasks struct {
	TaskId      string `json:"task_id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndTime     string `json:"end_time"`
	// request for searched task
}

type GetTasksByIdRes struct {
	TaskId      string     `json:"task_id"`
	UserId      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Files       []TaskFile `json:"files"`
	StartTime   string     `json:"start_time"`
	EndTime     string     `json:"end_time"`
	Status      string     `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

type SearchTitleReq struct {
	Q string `form:"q" validate:"required"`
}

// params for searched task
type SearchTitleParams struct {
	SearchQuery string `json:"search_query"`
}

// response for searched task
type SearchTaskRes struct {
	TaskId    string `json:"task_id"`
	Title     string `json:"title"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type GetAllExpiredRes struct {
	TaskId    string `json:"task_id"`
	Title     string `json:"title"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
