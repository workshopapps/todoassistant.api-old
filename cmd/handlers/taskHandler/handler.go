package taskHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/userEntity"
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
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}

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
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	task, errRes := t.srv.GetTaskByID(taskId, userId.(userEntity.CreateUserReq).UserId)

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

//Handle get all task from a specific user
func (t *taskHandler) GetAllTask(c *gin.Context) {
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	task, errRes := t.srv.GetAllTask(userId.(userEntity.CreateUserReq).UserId)

	if task == nil {
		message := "no Task with id " + userId.(userEntity.CreateUserReq).UserId + " exists"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, task, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find all task", errRes, nil))
		return
	}
	c.JSON(http.StatusOK, task)

}

//Handle Delete task by id
func (t *taskHandler) DeleteTaskById(c *gin.Context) {
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	_, errRes := t.srv.DeleteTaskByID(taskId, userId.(userEntity.CreateUserReq).UserId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Unable to delete task by id", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Task deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

//Handle Delete All Task of a user
func (t *taskHandler) DeleteAllTask(c *gin.Context) {
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	_, errRes := t.srv.DeleteAllTask(userId.(userEntity.CreateUserReq).UserId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Unable to delete task by id", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "All Task deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

//Update user Status
func (t *taskHandler) UpdateUserStatus(c *gin.Context) {
	var status string
	err := c.ShouldBind(&status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	updatedStatusTask, errRes := t.srv.UpdateTaskStatusByID(taskId, userId.(userEntity.CreateUserReq).UserId, status)
	if updatedStatusTask == nil {
		message := "no Task with id " + userId.(userEntity.CreateUserReq).UserId + " updated"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, updatedStatusTask, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find all task", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Task status updated successfully", nil)
	c.JSON(http.StatusOK, rd)
}

//Update task by id
func (t *taskHandler) EditTaskById(c *gin.Context) {
	var req *taskEntity.CreateTaskReq
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user id available", nil, nil))
		return
	}
	userId, exists := c.Get("userId")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}

	EditedTask, errRes := t.srv.EditTaskByID(taskId, userId.(userEntity.CreateUserReq).UserId, req)
	if EditedTask == nil {
		message := "no Task with id " + userId.(userEntity.CreateUserReq).UserId + " updated"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, EditedTask, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find all task", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Task status updated successfully", nil)
	c.JSON(http.StatusOK, rd)

}
