package taskHandler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/taskService"
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
