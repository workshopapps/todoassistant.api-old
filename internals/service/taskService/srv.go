package taskService

import (
	"context"
	"github.com/google/uuid"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/errorEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"
)

type TaskService interface {
	PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *errorEntity.ErrorRes)
}

type taskSrv struct {
	repo          taskRepo.TaskRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
}

func (t taskSrv) PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *errorEntity.ErrorRes) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	// implement validation for struct

	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, errorEntity.NewCustomError(400, "Bad system input")
	}

	//set time
	req.CreatedAt = t.timeSrv.CurrentTime().Format(time.RFC3339)
	//set id
	req.TaskId = uuid.New().String()
	// insert into db
	err = t.repo.Persist(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, errorEntity.NewCustomError(500, "Error Saving to Database")
	}
	data := taskEntity.CreateTaskRes{
		TaskId:      req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	// create a reminder

	return &data, nil

}

func NewTaskSrv(repo taskRepo.TaskRepository, timeSrv timeSrv.TimeService) TaskService {
	return &taskSrv{repo: repo, timeSrv: timeSrv}
}
