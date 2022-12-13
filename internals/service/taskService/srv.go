package taskService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
	"test-va/internals/service/loggerService"
	"test-va/internals/service/notificationService"
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
	UpdateTaskStatusByID(taskId string, req *taskEntity.UpdateTaskStatus) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	EditTaskByID(taskId string, req *taskEntity.EditTaskReq) (*taskEntity.EditTaskRes, *ResponseEntity.ServiceError)

	GetVADetails(userId string) (string, *ResponseEntity.ServiceError)
	AssignVAToTask(req *taskEntity.AssignReq) *ResponseEntity.ServiceError
	GetTaskAssignedToVA(vaId string) ([]*vaEntity.VATask, *ResponseEntity.ServiceError)
	GetAllTaskForVA() ([]*vaEntity.VATaskAll, *ResponseEntity.ServiceError)

	//comments
	PersistComment(req *taskEntity.CreateCommentReq) (*taskEntity.CreateCommentRes, *ResponseEntity.ServiceError)
	GetAllComments(taskId string) ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError)
	DeleteCommentByID(commentId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	GetComments() ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError)
}

type taskSrv struct {
	repo          taskRepo.TaskRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
	logger        loggerService.LogSrv
	remindSrv     reminderService.ReminderSrv
	nSrv          notificationService.NotificationSrv
}

func NewTaskSrv(repo taskRepo.TaskRepository, timeSrv timeSrv.TimeService,
	srv validationService.ValidationSrv, logSrv loggerService.LogSrv,
	reminderSrv reminderService.ReminderSrv,
	notificationSrv notificationService.NotificationSrv) TaskService {
	return &taskSrv{repo: repo, timeSrv: timeSrv, validationSrv: srv,
		logger: logSrv, remindSrv: reminderSrv, nSrv: notificationSrv}
}

func (t *taskSrv) GetTaskAssignedToVA(vaId string) ([]*vaEntity.VATask, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	va, err := t.repo.GetAllTaskAssignedToVA(ctx, vaId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return va, nil
}

func (t *taskSrv) GetAllTaskForVA() ([]*vaEntity.VATaskAll, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	va, err := t.repo.GetAllTaskForVA(ctx)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return va, nil
}

func (t *taskSrv) AssignVAToTask(req *taskEntity.AssignReq) *ResponseEntity.ServiceError {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	vaID, serviceError := t.GetVADetails(req.UserId)
	if serviceError != nil {
		log.Println(" error here", serviceError)
		return serviceError
	}

	err := t.repo.AssignTaskToVa(ctx, vaID, req.TaskId)
	if err != nil {
		log.Println(" error here 2", err)
		return ResponseEntity.NewInternalServiceError(err)
	}

	// t.nSrv.SendNotificationToVA(req.UserId, "Task Assigned", fmt.Sprintf("%s Just Assigned a Task to You", req.UserId), data)

	return nil
}

func (t *taskSrv) GetVADetails(userId string) (string, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	vaId, err := t.repo.GetVADetails(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ResponseEntity.NewInternalServiceError("No VA assigned yet")
		}
		return "", ResponseEntity.NewInternalServiceError(err)
	}

	return vaId, nil
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
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
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
		err = t.remindSrv.SetReminder(req)

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

	data := taskEntity.CreateTaskRes{
		TaskId:      req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		VAOption:    req.VAOption,
		Repeat:      req.Repeat,
	}

	if err != nil {
		fmt.Println("Error Uploading Notification to DB", err)
	}
	body := []notificationEntity.NotificationBody{
		{
			Content: fmt.Sprintf("%s Just Created a Task", req.UserId),
			Color:   notificationEntity.CreatedColor,
			Time:    time.Now().String(),
		},
	}

	tokens, vaId, err := t.nSrv.GetUserVaToken(req.UserId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vaId)
	if vaId != "" {
		err = t.nSrv.CreateNotification(vaId, "Task Created", time.Now().String(), fmt.Sprintf("%s just created a new task", req.UserId), notificationEntity.CreatedColor, req.TaskId)
		if err != nil {
			fmt.Println("Error Uploading Notification to DB", err)
		}
	}
	if len(tokens) > 0 {
		err := t.nSrv.SendBatchNotifications(tokens, "Task Created", body, data)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("User Has Not VA or VA Has Not Registered For Notifications")
	}

	switch req.Assigned {
	case "assigned":
		err = t.repo.PersistAndAssign(ctx, req)
		if err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), `"virtual_Assistant_id": converting NULL to string is unsupported`) {
				return nil, ResponseEntity.NewInternalServiceError("YOU DON'T HAVE A VA. GET YA MONEY UP. BROKE BOY.")
			}

			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	default:
		// insert into db
		err = t.repo.Persist(ctx, req)
		if err != nil {
			log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}

	}

	return &data, nil

}

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

// Get Task godoc
// @Summary	Get a single task
// @Description	Get a particular task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Success	200  {object}  taskEntity.GetTasksByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security BasicAuth
// @Router	/task/{taskId} [get]
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

// Get Expired Tasks godoc
// @Summary	Get all expired tasks
// @Description	Get all expired task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []taskEntity.GetAllExpiredRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security BasicAuth
// @Router	/task/expired [get]
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

func (t *taskSrv) GetAllTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
	log.Println("inside Fn")
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetAllTasks(ctx, userId)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// get all tasks of user assigned to va
func (t *taskSrv) GetAllVaTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
	log.Println("inside Fn")
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetAllTasks(ctx, userId)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

func (t *taskSrv) DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteTaskByID(ctx, taskId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Deleted successfully", nil, nil), nil
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
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "deleted user successfully", nil, nil), nil

}

func (t *taskSrv) UpdateTaskStatusByID(taskId string, req *taskEntity.UpdateTaskStatus) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.UpdateTaskStatusByID(ctx, taskId, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Updated status successfully", nil, nil), nil

}

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
		err = t.remindSrv.SetReminder(ndate)

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

func (t *taskSrv) PersistComment(req *taskEntity.CreateCommentReq) (*taskEntity.CreateCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	// implement validation for struct

	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	//set time
	req.CreatedAt = t.timeSrv.CurrentTime().Format(time.RFC3339)

	// insert into db
	err = t.repo.PersistComment(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := taskEntity.CreateCommentRes{
		TaskId:  req.TaskId,
		Comment: req.Comment,
	}

	return &data, nil

}

// get all comments
func (t *taskSrv) GetAllComments(taskId string) ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	comments, err := t.repo.GetAllComments(ctx, taskId)

	if comments == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return comments, nil

}

// get all comments
func (t *taskSrv) GetComments() ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	comments, err := t.repo.GetComments(ctx)

	if comments == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return comments, nil

}

// delete comment By ID
func (t *taskSrv) DeleteCommentByID(commentId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteCommentByID(ctx, commentId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Deleted successfully", nil, nil), nil
}
