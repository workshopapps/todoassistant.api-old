package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/cmd/handlers/taskHandler"
	"test-va/cmd/middlewares"
	"test-va/internals/service/taskService"
)

func TaskRoutes(v1 *gin.RouterGroup, service taskService.TaskService) {
	handler := taskHandler.NewTaskHandler(service)
	task := v1.Group("/task")
	task.Use(middlewares.ValidateJWT())
	{
		task.POST("", handler.CreateTask)
		task.GET("/:taskId", handler.GetTaskByID)
		task.GET("/pending/:userId", handler.GetPendingTasks)
		task.GET("/expired", handler.GetListOfExpiredTasks)
		task.GET("/", handler.GetAllTask)               //Get all task by a user
		task.DELETE("/:taskId", handler.DeleteTaskById) //Delete Task By ID
		//task.DELETE("/", handler.DeleteAllTask)               //Delete all task of a user
		task.POST("/status/:taskId", handler.UpdateUserStatus) //Update User Status
		task.PUT("/:taskId", handler.EditTaskById) //EditTaskById
		task.GET("/search", handler.SearchTask)

	}

}
