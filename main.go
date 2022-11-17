package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test-va/cmd/handlers/taskHandler"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	"test-va/internals/service/taskService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"
)

func main() {

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

	// time service
	timeSrv := timeSrv.NewTimeStruct()

	//validation service
	validationSrv := validationService.NewValidationStruct()

	// create service
	srv := taskService.NewTaskSrv(repo, timeSrv, validationSrv)

	handler := taskHandler.NewTaskHandler(srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2022"
	}

	router := chi.NewRouter()

	router.Use(middleware.DefaultLogger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	})

	router.Post("/task", handler.CreateTask)

	srvDetails := http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     router,
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
