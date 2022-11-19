package taskHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/taskService"

	"github.com/gin-gonic/gin"
)

type taskHandler struct {
	srv taskService.TaskService
}

func NewTaskHandler(srv taskService.TaskService) *taskHandler {
	return &taskHandler{srv: srv}
}

func (t *taskHandler) CreateTask(c *gin.Context) {
	var req taskEntity.CreateTaskReq

	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}

	task, errRes := t.srv.PersistTask(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error saving into db", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, task)

}

func (t *taskHandler) GetPendingTasks(c *gin.Context) {
	userId := c.Params.ByName("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	tasks, errRes := t.srv.GetPendingTasks(userId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "internal server error", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// handle search function
func (t *taskHandler) SearchTask(c *gin.Context) {

	name := c.Query("q")

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad params provided", "", nil))
		return
	}

	title := taskEntity.SearchTitleParams{
		SearchQuery: name,
	}

	searchedTasks, errRes := t.srv.SearchTask(&title)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error getting tasks", errRes, nil))
		return
	}

	length := len(searchedTasks)

	if length == 0 {
		message := "no Task with title " + title.SearchQuery + " found"
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, searchedTasks, nil))
		return
	}
	message := "successfully fetched Tasks with title " + title.SearchQuery + " and details"

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, searchedTasks, nil))
}

// handle get by ID
func (t *taskHandler) GetTaskByID(c *gin.Context) {
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	task, errRes := t.srv.GetTaskByID(taskId)

	if task == nil {
		message := "no Task with id " + taskId + " exists"
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, task, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "internal server error", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, task)
}

// handle Get expired Tasks

func (t *taskHandler) GetExpiredTasks(c *gin.Context) {
	userId := c.Params.ByName("userId")


	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}

	

	task, errRes := t.srv.GetExpiredTasks(userId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error saving into db", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, task)
}
