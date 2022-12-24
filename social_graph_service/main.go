package main

import (
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
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

	logger := log.New(os.Stdout, "[graph-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[social-graph] ", log.LstdFlags)

	store, err := data.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()
	//
	socialgraphHandler := handlers.NewMoviesHandler(logger, store)
	//
	router := mux.NewRouter()

	postUserNode := router.Methods(http.MethodPost).Subrouter()
	postUserNode.HandleFunc("/user", socialgraphHandler.CreateUser)

	//
	//	router.Use(moviesHandler.MiddlewareContentTypeSet)
	//
	//	getMovieByTitle := router.Methods(http.MethodGet).Subrouter()
	//	getMovieByTitle.HandleFunc("/movies/title/{title}", moviesHandler.GetAllMoviesByTitle)
	//
	//	getAllMoviesWithCast := router.Methods(http.MethodGet).Subrouter()
	//	getAllMoviesWithCast.HandleFunc("/movies/cast/{limit}", moviesHandler.GetAllMoviesWithCast)
	//
	//	getAllMovies := router.Methods(http.MethodGet).Subrouter()
	//	getAllMovies.HandleFunc("/movies/{limit}", moviesHandler.GetAllMovies)
	//
	//
	//	getActorRole := router.Methods(http.MethodGet).Subrouter()
	//	getActorRole.HandleFunc("/person/{role}", moviesHandler.GetActorRole)
	//
	//	postMovieNode := router.Methods(http.MethodPost).Subrouter()
	//	postMovieNode.HandleFunc("/movies", moviesHandler.CreateMovie)
	//	postMovieNode.Use(moviesHandler.MiddlewareMovieDeserialization)
	//
	//	getActorProducer := router.Methods(http.MethodGet).Subrouter()
	//	getActorProducer.HandleFunc("/actor/producer", moviesHandler.GetPersonWhoActedAndProducedMovie)
	//
	//	getKeanuMovies := router.Methods(http.MethodGet).Subrouter()
	//	getKeanuMovies.HandleFunc("/actor/keanu/{limit}", moviesHandler.GetKeanuMovies)
	//
	//	getActorsWithMostMovies := router.Methods(http.MethodGet).Subrouter()
	//	getActorsWithMostMovies.HandleFunc("/actor/most-movies/{limit}", moviesHandler.GetActorsWithMostMovies)
	//
	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
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
