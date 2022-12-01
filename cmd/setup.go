package cmd

import (
	"context"
	"fmt"

	"github.com/gin-contrib/cors"

	"log"
	"net/http"
	"os"
	"os/signal"
	"test-va/cmd/handlers/callHandler"
	"test-va/cmd/handlers/notificationHandler"
	"test-va/cmd/handlers/taskHandler"
	"test-va/cmd/handlers/userHandler"
	"test-va/cmd/middlewares"
	mySqlCallRepo "test-va/internals/Repository/callRepo/mySqlRepo"
	mySqlNotifRepo "test-va/internals/Repository/notificationRepo/mysqlRepo"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	mySqlRepo2 "test-va/internals/Repository/userRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	firebaseinit "test-va/internals/firebase-init"
	"test-va/internals/service/callService"
	"test-va/internals/service/cryptoService"
	log_4_go "test-va/internals/service/loggerService/log-4-go"
	"test-va/internals/service/notificationService"
	"test-va/internals/service/reminderService"
	"test-va/internals/service/taskService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/userService"
	"test-va/internals/service/validationService"
	"test-va/utils"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
)

func Setup() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	dsn := config.DataSourceName
	log.Println(dsn)
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

	// repo service
	repo := mySqlRepo.NewSqlRepo(conn)
	callRepo := mySqlCallRepo.NewSqlCallRepo(conn)
	userRepo := mySqlRepo2.NewMySqlUserRepo(conn)
	notificationRepo := mySqlNotifRepo.NewMySqlNotificationRepo(conn)
	// time service
	timeSrv := timeSrv.NewTimeStruct()

	//validation service
	validationSrv := validationService.NewValidationStruct()

	//Notification Service
	//Note Handle Unable to Connect to Firebase
	firebaseApp, err := firebaseinit.SetupFirebase()
	if err != nil {
		fmt.Println("UNABLE TO CONNECT TO FIREBASE", err)
	}
	notificationSrv := notificationService.New(firebaseApp, notificationRepo, validationSrv)
	err = notificationSrv.SendNotification("ckh2hTktbwD5VWfHUqIiH6:APA91bGtAyfluuCsR_-eCkDdwYBRZlRv9a6BBQGwumzttGV64H4OhMy6KILyRWy1bN1EvKQ6K131yS8oy4sR11ofTgSFPSpeviXQPYdt_PMhXI8a1RJm8I8lemh-iU8uFym3TPOSPspn", "Notification", "notification", "hello")
	if err != nil {
		fmt.Println("Could Not Send Message", err)
	}

	// create cron tasks for checking if time is due

	s := gocron.NewScheduler(time.UTC)
	reminderSrv := reminderService.NewReminderSrv(s, conn, repo, notificationSrv)
	reminderSrv.ScheduleNotificationEverySixHours()
	reminderSrv.ScheduleNotificationDaily()

	s.Every(5).Minutes().Do(func() {
		log.Println("checking for 5 minutes reminders")
		reminderSrv.SetReminderEvery5Min()
	})

	s.Every(30).Minutes().Do(func() {
		log.Println("checking for 30 minutes reminders")
		reminderSrv.SetReminderEvery30Min()
	})
	s.StartAsync()

	//logger service
	logger := log_4_go.NewLogger()
	//crypto service
	cryptoSrv := cryptoService.NewCryptoSrv()

	// create service
	taskSrv := taskService.NewTaskSrv(repo, timeSrv, validationSrv, logger, reminderSrv)
	userSrv := userService.NewUserSrv(userRepo, validationSrv, timeSrv, cryptoSrv)

	callSrv := callService.NewCallSrv(callRepo, timeSrv, validationSrv, logger)

	handler := taskHandler.NewTaskHandler(taskSrv)

	userHandler := userHandler.NewUserHandler(userSrv)

	callHandler := callHandler.NewCallHandler(callSrv)
	notificationHandler := notificationHandler.NewNotificationHandler(notificationSrv)

	port := config.SeverAddress
	log.Println(port)
	if port == "" {
		port = "2022"
	}

	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	//r.Use(middlewares.Logger())

	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "POST", "GET"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	v1 := r.Group("/api/v1")
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
		//task.POST("/:taskId", handler.UpdateUserStatus) //Update User Status
		task.PUT("/:taskId", handler.EditTaskById) //EditTaskById

	}

	//r.POST("/task", handler.CreateTask)
	v1.GET("/calls", callHandler.GetCalls)
	//r.GET("/task/pending/:userId", handler.GetPendingTasks)
	//get list of pending tasks belonging to a user
	//r.GET("/task/expired/", handler.GetListOfExpiredTasks)
	// get task by id
	//r.GET("/task/:taskId", handler.GetTaskByID)
	// search route
	v1.GET("/search", handler.SearchTask)

	//chat service connection

	pusherClient := pusher.Client{
		AppID:   "1512808",
		Key:     "f79030d90753a91854e6",
		Secret:  "06b8abef8713abd21cc9",
		Cluster: "eu",
		Secure:  true,
	}

	v1.POST("dashboard/assistant", func(c *gin.Context) {
		// var data map[string]string
		var data map[string]string

		if err := c.BindJSON(&data); err != nil {
			return
		}
		pusherClient.Trigger("vachat", "message", data)

		c.JSON(http.StatusOK, []string{})
	})

	// USER
	//create user
	// Register a user

	v1.POST("/user", userHandler.CreateUser)
	// Login into the user account
	v1.POST("/user/login", userHandler.Login)
	users := v1.Group("/user")

	users.Use(middlewares.ValidateJWT())
	{
		// Get all users
		users.GET("", userHandler.GetUsers)
		// Get a specific user
		users.GET("/:user_id", userHandler.GetUser)
		// Update a specific user
		users.PUT("/:user_id", userHandler.UpdateUser)
		// Change user password

		users.PUT("/:user_id/change-password", userHandler.ChangePassword)
		// Delete a user
		users.DELETE("/:user_id", userHandler.DeleteUser)
	}

	// Notifications
	// Register to Recieve Notifications
	v1.POST("/notification", notificationHandler.RegisterForNotifications)

	v1.GET("/ping", func(c *gin.Context) {

		c.String(http.StatusOK, "pong")
	})

	v1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Ticked Backend Server - V1.0.0")
	})

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
