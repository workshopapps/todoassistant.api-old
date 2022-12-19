package routes

import (
	"test-va/cmd/handlers/taskHandler"
	"test-va/cmd/middlewares"
	"test-va/cmd/middlewares/vaMiddleware"

	"test-va/internals/service/taskService"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(v1 *gin.RouterGroup, service taskService.TaskService, srv tokenservice.TokenSrv) {
	mWare := vaMiddleware.NewVaMiddleWare(srv)
	jwtMWare := middlewares.NewJWTMiddleWare(srv)

	handler := taskHandler.NewTaskHandler(service)
	task := v1.Group("/task")
	task2 := v1.Group("/task")

	task.Use(jwtMWare.ValidateJWT())
	{
		task.POST("", handler.CreateTask)
		task.GET("/:taskId", handler.GetTaskByID)
		task.GET("/pending/:userId", handler.GetPendingTasks)
		task.GET("/expired", handler.GetListOfExpiredTasks)
		task.GET("/", handler.GetAllTask)               //Get all task by a user
		task.DELETE("/:taskId", handler.DeleteTaskById) //Delete Task By ID
		//task.DELETE("/", handler.DeleteAllTask)               //Delete all task of a user
		task.PUT("/status/:taskId", handler.UpdateTaskStatus) //Update task status

		//comments
		task.POST("/comment", handler.CreateComment)              //comment on task
		task.GET("/comment/:taskId", handler.GetComments)         //get all comment on task
		task.GET("/comment/all", handler.GetAllComments)          //get all comment available
		task.DELETE("/comment/:commentId", handler.DeleteComment) //delete comment

		//task.PUT("/comment", handler.CreateComment)   //edit comment
		task.PUT("/:taskId", handler.EditTaskById) //EditTaskById
		task.GET("/search", handler.SearchTask)

		//assign task to VA
		task.POST("/assign/:taskId", handler.AssignTaskToVA)
	}

	task2.Use(mWare.MapVAToReq)
	{
		//list of all task assigned to VA
		task2.GET("/all/va", handler.GetTasksAssignedToVa)
		// get alllll task
		task2.GET("/all", handler.GetAllTasksAssignedForVa)
		// Get list of user all Pending task for Va
		task2.GET("/all/pendingtasks", handler.GetListOfPendingTasks)
	}

}
