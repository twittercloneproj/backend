package main

import (
	"context"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"net/http"
	"os"
	"os/signal"
	"social_graph_service/data"
	"social_graph_service/handlers"
	"time"
)

func main() {
	//Reading from environment, if not set we will default it to 8080.
	//This allows flexibility in different environments (for eg. when running multiple docker api's and want to override the default port)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8002"
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

	store, err := data.New(log)
	if err != nil {
		log.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()
	//
	socialgraphHandler := handlers.NewMoviesHandler(log, store)
	//
	router := mux.NewRouter()

	postUserNode := router.Methods(http.MethodPost).Subrouter()
	postUserNode.HandleFunc("/user", socialgraphHandler.CreateUser)

	follow := router.Methods(http.MethodPost).Subrouter()
	follow.HandleFunc("/follow/{username}", socialgraphHandler.Follow)

	checkFollow := router.Methods(http.MethodGet).Subrouter()
	checkFollow.HandleFunc("/follow/{username}", socialgraphHandler.CheckFollow)

	canAccessTweet := router.Methods(http.MethodGet).Subrouter()
	canAccessTweet.HandleFunc("/access-tweet/{username}", socialgraphHandler.CanAccessTweet)

	followRequest := router.Methods(http.MethodPost).Subrouter()
	followRequest.HandleFunc("/request/{username}", socialgraphHandler.AcceptRejectRequest)

	removeFollow := router.Methods(http.MethodPost).Subrouter()
	removeFollow.HandleFunc("/unfollow/{username}", socialgraphHandler.RemoveFollow)

	requests := router.Methods(http.MethodGet).Subrouter()
	requests.HandleFunc("/requests", socialgraphHandler.GetFollowRequests)

	followers := router.Methods(http.MethodGet).Subrouter()
	followers.HandleFunc("/followers/{username}", socialgraphHandler.GetFollowersForUser)

	following := router.Methods(http.MethodGet).Subrouter()
	following.HandleFunc("/{username}/following", socialgraphHandler.GetFollowingUsers)

	changePrivacy := router.Methods(http.MethodPost).Subrouter()
	changePrivacy.HandleFunc("/change-privacy", socialgraphHandler.ChangePrivacy)

	profileSuggestion := router.Methods(http.MethodGet).Subrouter()
	profileSuggestion.HandleFunc("/suggestions", socialgraphHandler.GetSuggestedUsers)

	//cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))
	//cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"http://localhost:4200"}),
	//	gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}),
	//	gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
	//	gorillaHandlers.AllowCredentials())

	server := http.Server{
		Addr: ":" + port,
		//Handler:      cors(router),
		Handler: router,
		//IdleTimeout:  120 * time.Second,
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 5 * time.Second,
	}

	log.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
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
