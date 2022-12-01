package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-va/cmd/routes"
	mySqlNotifRepo "test-va/internals/Repository/notificationRepo/mysqlRepo"
	mySqlRepo4 "test-va/internals/Repository/subscribeRepo/mySqlRepo"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	mySqlRepo2 "test-va/internals/Repository/userRepo/mySqlRepo"
	mySqlRepo3 "test-va/internals/Repository/vaRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	firebaseinit "test-va/internals/firebase-init"
	"test-va/internals/service/cryptoService"
	log_4_go "test-va/internals/service/loggerService/log-4-go"
	"test-va/internals/service/notificationService"
	"test-va/internals/service/reminderService"
	"test-va/internals/service/subscribeService"
	"test-va/internals/service/taskService"
	"test-va/internals/service/timeSrv"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/userService"
	"test-va/internals/service/vaService"
	"test-va/internals/service/validationService"
	"test-va/utils"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/go-co-op/gocron"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
)

func Setup() {

	//Load configurations
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	dsn := config.DataSourceName
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	port := config.SeverAddress
	if port == "" {
		port = "2022"
	}

	secret := config.TokenSecret
	if secret == "" {
		log.Fatal("secret key not found")
	}

	//Repo

	//db service
	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()

	// task repo service
	repo := mySqlRepo.NewSqlRepo(conn)

	//user repo service
	userRepo := mySqlRepo2.NewMySqlUserRepo(conn)

	//notification repo service
	notificationRepo := mySqlNotifRepo.NewMySqlNotificationRepo(conn)

	//va repo service
	vaRepo := mySqlRepo3.NewVASqlRepo(conn)

	// subscribe repo
	subRepo := mySqlRepo4.NewMySqlSubscribeRepo(conn)

	//SERVICES

	//time service
	timeSrv := timeSrv.NewTimeStruct()

	//callRepo := mySqlCallRepo.NewSqlCallRepo(conn)

	// cron service
	s := gocron.NewScheduler(time.UTC)

	// reminder service and implementation
	reminderSrv := reminderService.NewReminderSrv(s, conn, repo)

	s.Every(5).Minutes().Do(func() {
		log.Println("checking for 5 minutes reminders")
		reminderSrv.SetReminderEvery5Min()
	})

	s.Every(30).Minutes().Do(func() {
		log.Println("checking for 30 minutes reminders")
		reminderSrv.SetReminderEvery30Min()
	})

	// run cron jobs
	s.StartAsync()

	//validation service
	validationSrv := validationService.NewValidationStruct()

	// token service
	srv := tokenservice.NewTokenSrv(secret)

	//logger service
	logger := log_4_go.NewLogger()

	//crypto service
	cryptoSrv := cryptoService.NewCryptoSrv()

	//Notification Service
	//Note Handle Unable to Connect to Firebase
	firebaseApp, err := firebaseinit.SetupFirebase()
	if err != nil {
		fmt.Println("UNABLE TO CONNECT TO FIREBASE", err)
	}
	notificationSrv := notificationService.New(firebaseApp, notificationRepo, validationSrv)
	err = notificationSrv.SendNotification("ckh2hTktbwD5VWfHUqIiH6:APA91bGtAyfluuCsR_"+
		"-eCkDdwYBRZlRv9a6BBQGwumzttGV64H4OhMy6KILyRWy1bN1EvKQ6K131yS8oy4sR11ofTgSFPSpeviXQPYdt"+
		"_PMhXI8a1RJm8I8lemh-iU8uFym3TPOSPspn",
		"Notification", "notification",
		[]string{"hello"})
	if err != nil {
		fmt.Println("Could Not Send Message", err)
	}

	// task service
	taskSrv := taskService.NewTaskSrv(repo, timeSrv, validationSrv, logger, reminderSrv)

	// user service
	userSrv := userService.NewUserSrv(userRepo, validationSrv, timeSrv, cryptoSrv)

	// va service
	vaSrv := vaService.NewVaService(vaRepo, validationSrv, timeSrv, cryptoSrv)

	// subscribe service
	subscribeSrv := subscribeService.NewSubscribeSrv(subRepo)

	//router setup
	r := gin.New()
	v1 := r.Group("/api/v1")

	// Middlewares
	v1.Use(gin.Logger())
	v1.Use(gin.Recovery())
	v1.Use(gzip.Gzip(gzip.DefaultCompression))
	v1.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	// routes

	//ping route
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//welcome message route
	v1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Ticked Backend Server - V1.0.0")
	})

	//handle user routes
	routes.UserRoutes(v1, userSrv)

	//handle task routes
	routes.TaskRoutes(v1, taskSrv, srv)

	//handle Notifications
	routes.NotificationRoutes(v1, notificationSrv)

	//handle VA
	routes.VARoutes(v1, vaSrv, srv)

	//handle subscribe route
	routes.SubscribeRoutes(v1, subscribeSrv)

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

	// Notifications
	// Register to Receive Notifications
	//v1.POST("/notification", notificationHandler.RegisterForNotifications)

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
