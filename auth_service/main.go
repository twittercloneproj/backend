package main

import (
	social_graph "auth_service/client/social-graph"
	"auth_service/data"
	"auth_service/handlers"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8003"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new log file rotator
	logWriter, err := rotatelogs.New(
		"twitterlogs-%Y%m%d.txt", // file name pattern
		rotatelogs.WithLinkName("twitterlogs.txt"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.Fatalf("failed to create rotate logs: %s", err)
	}

	log := &logrus.Logger{
		// Log into f file handler and on os.Stdout
		Out:   io.MultiWriter(logWriter, os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &easy.Formatter{
			LogFormat: "[%lvl%]: - %msg%\n",
		},
	}

	store, err := data.New(timeoutContext, log)

	if err != nil {
		log.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	store.Ping()

	socialGraphClient := social_graph.NewClient("social_graph_service", "8002")

	usersHandler := handlers.NewUsersHandler(log, store, socialGraphClient)

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/all", usersHandler.GetAllUsers)

	//getRouterB := router.Methods(http.MethodGet).Subrouter()
	//getRouterB.HandleFunc("/users/allB", usersHandler.GetAllBUsers)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users/add", usersHandler.PostUsers)
	postRouter.Use(usersHandler.MiddlewarePatientDeserialization)

	loginRouter := router.Methods(http.MethodPost).Subrouter()
	loginRouter.HandleFunc("/users/login", usersHandler.Login)

	//postBRouter := router.Methods(http.MethodPost).Subrouter()
	//postBRouter.HandleFunc("/users/addB", usersHandler.PostBUsers)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"http://localhost:4200"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
		gorillaHandlers.AllowCredentials())

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	//certFile := "twitter.crt"
	//keyFile := "twitter.key"

	log.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		//err := server.ListenAndServeTLS(certFile, keyFile)

		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	log.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		log.Fatal("Cannot gracefully shutdown...")
	}
	log.Println("Server stopped")
}
