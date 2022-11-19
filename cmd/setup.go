package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-va/cmd/handlers/callHandler"
	"test-va/cmd/handlers/taskHandler"
	"test-va/cmd/middlewares"
	mySqlCallRepo "test-va/internals/Repository/callRepo/mySqlRepo"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	"test-va/internals/service/callService"
	log_4_go "test-va/internals/service/loggerService/log-4-go"
	"test-va/internals/service/taskService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Setup() {
	dsn := os.Getenv("dsn")
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()

	// create cron tasks for checking if time is due

	// repo service
	repo := mySqlRepo.NewSqlRepo(conn)

	callRepo := mySqlCallRepo.NewSqlCallRepo(conn)

	// time service
	timeSrv := timeSrv.NewTimeStruct()

	//validation service
	validationSrv := validationService.NewValidationStruct()
	//logger service
	logger := log_4_go.NewLogger()

	// create service
	srv := taskService.NewTaskSrv(repo, timeSrv, validationSrv, logger)

	callSrv := callService.NewCallSrv(callRepo,timeSrv, validationSrv, logger)

	handler := taskHandler.NewTaskHandler(srv)

	callHandler := callHandler.NewCallHandler(callSrv)
	port := os.Getenv("PORT")
	if port == "" {
		port = "2022"
	}

	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	//r.Use(middlewares.Logger())

	r.Use(gin.Recovery())
	r.Use(middlewares.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.POST("/task", handler.CreateTask)
	r.GET("/calls", callHandler.GetCalls)
	r.GET("/task/pending/:userId", handler.GetPendingTasks)
	// get task by id
	r.GET("/task/:taskId", handler.GetTaskByID)
	// search route
	r.GET("/search", handler.SearchTask)

	// Get expired task
	r.GET("/task/expired/:userId", handler.GetExpiredTasks)


	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	srvDetails := http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		log.Println("SERVER STARTING ON PORT:", port)
		err := srvDetails.ListenAndServe()
		if err != nil {
			log.Printf("ERROR STARTING SERVER: %v", err)
			os.Exit(1)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Printf("Closing now, We've gotten signal: %v", sig)

	ctx := context.Background()
	srvDetails.Shutdown(ctx)
}
