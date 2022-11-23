package taskService

import (
	"context"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/loggerService"
	"test-va/internals/service/reminderService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError)
	GetPendingTasks(userId string) ([]*taskEntity.GetPendingTasksRes, *ResponseEntity.ServiceError)
	SearchTask(req *taskEntity.SearchTitleParams) ([]*taskEntity.SearchTaskRes, *ResponseEntity.ServiceError)
	GetTaskByID(taskId string) (*taskEntity.GetTasksByIdRes, *ResponseEntity.ServiceError)
	GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError)
	DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	GetAllTask() ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError)
	DeleteAllTask() (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
}

type taskSrv struct {
	repo          taskRepo.TaskRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
	logger        loggerService.LogSrv
	remindSrv     reminderService.ReminderSrv
}

func NewTaskSrv(repo taskRepo.TaskRepository, timeSrv timeSrv.TimeService, srv validationService.ValidationSrv, logSrv loggerService.LogSrv, reminderSrv reminderService.ReminderSrv) TaskService {
	return &taskSrv{repo: repo, timeSrv: timeSrv, validationSrv: srv, logger: logSrv, remindSrv: reminderSrv}
}

func (t *taskSrv) GetPendingTasks(userId string) ([]*taskEntity.GetPendingTasksRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	tasks, err := t.repo.GetPendingTasks(userId, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return tasks, nil
}

func (t *taskSrv) PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	// implement validation for struct

	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	//check if timeDueDate and StartDate is valid
	err = t.timeSrv.CheckFor339Format(req.EndTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	}

	err = t.timeSrv.CheckFor339Format(req.StartTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	}

	//set time
	req.CreatedAt = t.timeSrv.CurrentTime().Format(time.RFC3339)
	//set id
	req.TaskId = uuid.New().String()
	req.Status = "PENDING"
	// insert into db
	err = t.repo.Persist(ctx, req)
	if err != nil {
		log.Println(err, "rrrr")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := taskEntity.CreateTaskRes{
		TaskId:      req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	// create a reminder
	err = t.remindSrv.SetReminder(req.EndTime, req.TaskId)

	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return &data, nil

}

// search task by name func
func (t *taskSrv) SearchTask(title *taskEntity.SearchTitleParams) ([]*taskEntity.SearchTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.validationSrv.Validate(title)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError(err)
	}
	tasks, err := t.repo.SearchTasks(title, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return tasks, nil
}

func (t *taskSrv) GetTaskByID(taskId string) (*taskEntity.GetTasksByIdRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetTaskByID(taskId, ctx)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

func (t *taskSrv) GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	task, err := t.repo.GetListOfExpiredTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

//Get all task service
func (t *taskSrv) GetAllTask() ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetAllTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

//Delete task by Id
func (t *taskSrv) DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteTaskByID(taskId, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Deleted successfully", nil), nil
}

//Delete All task
func (t *taskSrv) DeleteAllTask() (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteAllTask(ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Deleted successfully", nil), nil

}

//Update task status
func (t *taskSrv) UpdateTaskStatusByID(taskId string, status string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.UpdateTaskStatusByID(taskId, status, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Updated successfully successfully", nil), nil

}

// Edit task by Id
func (t *taskSrv) EditTaskByID(taskId string, req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	//validating the struct
	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}
	//Get the task by the ID
	task, err := t.repo.GetTaskByID(taskId, ctx)

	if req.TaskId == "" {
		task.TaskId = taskId
	} else {
		task.TaskId = req.TaskId

	}
	if req.Title == "" {
		task.Title = task.Title
	} else {
		task.Title = req.Title
	}

	if req.Description == "" {
		task.Description = task.Description
	} else {
		task.Description = req.Description
	}

	if req.EndTime == "" {
		task.EndTime = task.EndTime
	} else {
		task.EndTime = req.EndTime
	}

	task.UpdatedAt = t.timeSrv.CurrentTime().Format(time.RFC3339)

	//Update Task
	err = t.repo.EditTaskById(taskId, req, ctx)

	if err != nil {
		log.Println(err, "error creating data")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	//Returning Data
	data := taskEntity.CreateTaskRes{
		TaskId:      req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	// create a reminder
	err = t.remindSrv.SetReminder(req.EndTime, req.TaskId)

	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return &data, nil

}
