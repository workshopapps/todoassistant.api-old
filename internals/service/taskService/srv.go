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
	GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError)
	GetListOfPendingTasks() ([]*taskEntity.GetAllPendingRes, *ResponseEntity.ServiceError)
	DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	GetAllTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError)
	GetTaskByID(taskId string) (*taskEntity.GetTasksByIdRes, *ResponseEntity.ServiceError)
	DeleteAllTask(userId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	UpdateTaskStatusByID(taskId string, userId string, status string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	EditTaskByID(taskId string, req *taskEntity.EditTaskReq) (*taskEntity.EditTaskRes, *ResponseEntity.ServiceError)
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

	// create a reminder
	switch req.Repeat {
	case "never":
		err = t.remindSrv.SetReminder(req.EndTime, req.TaskId)

		if err != nil {
			log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	case "daily":
		err = t.remindSrv.SetDailyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Daily Input")
		}
	case "weekly":
		err = t.remindSrv.SetWeeklyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Weekly Input")
		}
	case "bi-weekly":
		err = t.remindSrv.SetBiWeeklyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Bi Weekly Input")
		}
	case "monthly":
		err = t.remindSrv.SetMonthlyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Monthly Input")
		}
	case "yearly":
		err = t.remindSrv.SetYearlyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Yearly Input")
		}
	default:
		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Input(check enum data)")
	}

	// insert into db
	err = t.repo.Persist(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := taskEntity.CreateTaskRes{
		TaskId:      req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		VAOption:    req.VAOption,
		Repeat:      req.Repeat,
	}

	return &data, nil

}

// Create Task
func (t *taskSrv) CreateTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError) {
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

	task, err := t.repo.GetTaskByID(ctx, taskId)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	log.Println("From getByID", task)
	return task, nil

}

func (t *taskSrv) GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetListOfExpiredTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError("No Task")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

func (t *taskSrv) GetListOfPendingTasks() ([]*taskEntity.GetAllPendingRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetListOfPendingTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError("No Task")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// Get all task service
func (t *taskSrv) GetAllTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	task, err := t.repo.GetAllTasks(ctx, userId)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// Delete task by Id

func (t *taskSrv) DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteTaskByID(ctx, taskId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Deleted successfully", nil), nil
}

// Delete All task

func (t *taskSrv) DeleteAllTask(userId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteAllTask(ctx, userId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Deleted successfully", nil), nil

}

// Update task status
func (t *taskSrv) UpdateTaskStatusByID(taskId string, userId string, status string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.UpdateTaskStatusByID(ctx, taskId, userId, status)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(200, "Updated successfully successfully", nil), nil

}

// Edit task by Id

func (t *taskSrv) EditTaskByID(taskId string, req *taskEntity.EditTaskReq) (*taskEntity.EditTaskRes, *ResponseEntity.ServiceError) {

	// create context of 1 minute

	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	//validating the struct
	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	// Get task by ID
	task, err := t.repo.GetTaskByID(ctx, taskId)
	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	log.Println(req)
	//Update Task
	//err = t.repo.EditTaskById(ctx, taskId, req)

	// if err != nil {
	// 	log.Println(err, "error creating data")
	// 	return nil, ResponseEntity.NewInternalServiceError(err)
	// }

	//Returning Data
	data := taskEntity.EditTaskRes{
		Title:       req.Title,
		Description: req.Description,
		Repeat:      req.Repeat,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		VAOption:    req.VAOption,
		Status:      req.Status,
	}
	updateAt := t.timeSrv.CurrentTime().Format(time.RFC3339)
	ndate := &taskEntity.CreateTaskReq{
		TaskId:      taskId,
		UserId:      task.UserId,
		Title:       data.Title,
		Description: data.Description,
		Repeat:      data.Repeat,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		VAOption:    data.VAOption,
		Status:      data.Status,
		UpdatedAt:   updateAt,
		CreatedAt:   task.CreatedAt,
	}

	// delete former task
	_, err2 := t.DeleteTaskByID(taskId)
	if err2 != nil {
		log.Println(err2)
		return nil, err2
	}
	log.Println("Deleted task line 381")

	// create new task

	//check if timeDueDate and StartDate is valid
	err = t.timeSrv.CheckFor339Format(ndate.EndTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	}

	err = t.timeSrv.CheckFor339Format(ndate.StartTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	}

	// create a reminder
	switch req.Repeat {
	case "never":
		err = t.remindSrv.SetReminder(ndate.EndTime, ndate.TaskId)

		if err != nil {
			log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	case "daily":
		err = t.remindSrv.SetDailyReminder(ndate)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Daily Input")
		}
	case "weekly":
		err = t.remindSrv.SetWeeklyReminder(ndate)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Weekly Input")
		}
	case "bi-weekly":
		err = t.remindSrv.SetBiWeeklyReminder(ndate)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Bi Weekly Input")
		}
	case "monthly":
		err = t.remindSrv.SetMonthlyReminder(ndate)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Monthly Input")
		}
	case "yearly":
		err = t.remindSrv.SetYearlyReminder(ndate)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Yearly Input")
		}
	default:
		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Input(check enum data)")
	}

	// insert into db
	err = t.repo.Persist(ctx, ndate)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	log.Println("update complete")

	return &data, nil
}
