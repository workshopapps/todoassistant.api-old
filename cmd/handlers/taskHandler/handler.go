package taskHandler

import (
	"log"
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
	value := c.GetString("userId")
	log.Println("value is: ", value)
	if value == "" {
		log.Println("112")
		c.AbortWithStatusJSON(http.StatusUnauthorized, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "you are not allowed to access this resource", nil, nil))
		return
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}

	req.UserId = value
	task, errRes := t.srv.PersistTask(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "error creating Task", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, task)

}

func (t *taskHandler) GetPendingTasks(c *gin.Context) {
	userId := c.Params.ByName("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}

	tasks, errRes := t.srv.GetPendingTasks(userId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Error Finding Pending Tasks", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (t *taskHandler) GetListOfExpiredTasks(c *gin.Context) {

	tasks, errRes := t.srv.GetListOfExpiredTasks()
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Error finding Expired Tasks", errRes, nil))
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error searching for tasks", errRes, nil))
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
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	task, errRes := t.srv.GetTaskByID(taskId)

	if task == nil {
		message := "no Task with id " + taskId + " exists"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, task, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find Task By Id", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, task)
}

