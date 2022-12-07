package taskRepo

import (
	"context"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
)

type TaskRepository interface {
	Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error
	PersistAndAssign(ctx context.Context, req *taskEntity.CreateTaskReq) error
	GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error)
	GetTaskByID(ctx context.Context, taskId string) (*taskEntity.GetTasksByIdRes, error)
	SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes, error)
	GetListOfExpiredTasks(ctx context.Context) ([]*taskEntity.GetAllExpiredRes, error)
	GetListOfPendingTasks(ctx context.Context) ([]*taskEntity.GetAllPendingRes, error)
	SetTaskToExpired(id string) error

	GetAllTasks(ctx context.Context, userId string) ([]*taskEntity.GetAllTaskRes, error)
	DeleteTaskByID(ctx context.Context, taskId string) error
	DeleteAllTask(ctx context.Context, userId string) error
	UpdateTaskStatusByID(ctx context.Context, taskId string) error
	EditTaskById(ctx context.Context, taskId string, req *taskEntity.EditTaskReq) error
	SetNewEvent(req *taskEntity.CreateTaskReq) error

	//VA
	GetAllTaskAssignedToVA(ctx context.Context, vaId string) ([]*vaEntity.VATask, error)
	GetVADetails(ctx context.Context, userId string) (string, error)
	AssignTaskToVa(ctx context.Context, vaId, taskId string) error

	//Comment
	PersistComment(ctx context.Context, req *taskEntity.CreateCommentReq) error
	GetAllComments(ctx context.Context, taskId string) ([]*taskEntity.GetCommentRes, error)
}
