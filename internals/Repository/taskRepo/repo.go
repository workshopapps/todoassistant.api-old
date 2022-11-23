package taskRepo

import (
	"context"
	"test-va/internals/entity/taskEntity"
)

type TaskRepository interface {
	Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error
	GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error)
	GetTaskByID(taskId string, ctx context.Context) (*taskEntity.GetTasksByIdRes, error)
	SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes, error)
	GetListOfExpiredTasks(ctx context.Context) ([]*taskEntity.GetAllExpiredRes, error)
	SetTaskToExpired(id string) error
	GetAllTasks(ctx context.Context) ([]*taskEntity.GetAllTaskRes, error)
	DeleteTaskByID(id string, ctx context.Context) error
	DeleteAllTask(ctx context.Context) error
	UpdateTaskStatusByID(taskId string, status string, ctx context.Context) error
	EditTaskById(taskId string, req *taskEntity.CreateTaskReq, ctx context.Context) error
}
