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
	GetExpiredTasks(userId string, ctx context.Context) ([]*taskEntity.GetExpiredTaskRes, error)
}
