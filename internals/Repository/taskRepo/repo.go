package taskRepo

import (
	"context"
	"test-va/internals/entity/taskEntity"
)

type TaskRepository interface {
	Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error
	GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error)
}
