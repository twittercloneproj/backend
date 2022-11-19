package main

import (
	"auth_service/data"
	"auth_service/handlers"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
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

	logger := log.New(os.Stdout, "[user-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[user-store] ", log.LstdFlags)

	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	store.Ping()

	usersHandler := handlers.NewUsersHandler(logger, store)

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/all", usersHandler.GetAllUsers)

	//getRouterB := router.Methods(http.MethodGet).Subrouter()
	//getRouterB.HandleFunc("/users/allB", usersHandler.GetAllBUsers)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users/add", usersHandler.PostUsers)
	postRouter.Use(usersHandler.MiddlewarePatientDeserialization)

	//postBRouter := router.Methods(http.MethodPost).Subrouter()
	//postBRouter.HandleFunc("/users/addB", usersHandler.PostBUsers)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServeTLS("sertifikat/twitter.cer", "sertifikat/twitterSSL.key")
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
