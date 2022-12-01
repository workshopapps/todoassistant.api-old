package taskHandler

import (
	"log"
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
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error creating Task", errRes, nil))
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
	userId := c.GetString("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No User ID", nil, nil))
		return
	}

	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no task id available", nil, nil))
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

// Handle get all task from a specific user

func (t *taskHandler) GetAllTask(c *gin.Context) {
	log.Println("here")
	userId := c.GetString("userId")
	log.Println(userId)
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	task, errRes := t.srv.GetAllTask(userId)
	if task == nil {
		message := "no Task with id " + userId + " exists"
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

// Handle Delete task by id

func (t *taskHandler) DeleteTaskById(c *gin.Context) {
	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "invalid taskId id", nil, nil))
		return
	}
	userId := c.GetString("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Authentication Error, Invalid UserId", nil, nil))
		return
	}
	_, errRes := t.srv.DeleteTaskByID(taskId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Unable to delete task by id", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Task deleted successfully", nil, nil)
	c.JSON(http.StatusOK, rd)
}

// Handle Delete All Task of a user

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
	rd := ResponseEntity.BuildSuccessResponse(200, "All Task deleted successfully", nil, nil)
	c.JSON(http.StatusOK, rd)
}

// Update user Status

func (t *taskHandler) UpdateUserStatus(c *gin.Context) {
	param := c.Param("taskId")
	if param == "" {
	   c.AbortWithStatusJSON(http.StatusBadRequest,
		  ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no task id available", nil, nil))
	   return
	}
 
	_, errRes := t.srv.UpdateTaskStatusByID(param)
	if errRes != nil {
	   c.AbortWithStatusJSON(http.StatusInternalServerError,
		  ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
			 "Error Setting Task to Done", errRes, nil))
	   return
	}
	rd := ResponseEntity.BuildSuccessResponse(http.StatusOK,
	   "Task status updated successfully", nil, nil)
	c.JSON(http.StatusOK, rd)
 }

// Update task by id

func (t *taskHandler) EditTaskById(c *gin.Context) {
	var req taskEntity.EditTaskReq

	taskId := c.Params.ByName("taskId")
	if taskId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no task id provided", nil, nil))
		return
	}
	//log.Println(taskId)
	err := c.ShouldBind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Input Request", err, nil))
		return
	}
	//log.Println(req)
	EditedTask, errRes := t.srv.EditTaskByID(taskId, &req)

	if errRes != nil {
		message := "no task with ID: " + taskId + " found"
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, message, errRes, nil))
		return
	}

	rd := ResponseEntity.BuildSuccessResponse(200, "Task status updated successfully", EditedTask, nil)
	c.JSON(http.StatusOK, rd)

}
