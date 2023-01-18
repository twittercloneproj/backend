package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	social_graph "tweet_service/client/social-graph"
	"tweet_service/data"
	"tweet_service/handlers"
)

func main() {
	//Reading from environment, if not set we will default it to 8080.
	//This allows flexibility in different environments (for eg. when running multiple docker api's and want to override the default port)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8001"
	}

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[tweet-api] ", log.LstdFlags)

	// NoSQL: Initialize Product Repository store
	store, err := data.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	socialGraphClient := social_graph.NewClient("social_graph_service", "8002")

	//Initialize the handler and inject said logger
	tweetsHandler := handlers.NewTweetsHandler(logger, store, socialGraphClient)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()

	getAllRouter := router.Methods(http.MethodGet).Subrouter()
	getAllRouter.HandleFunc("/all", tweetsHandler.GetAllTweets)

	getTweetList := router.Methods(http.MethodGet).Subrouter()
	getTweetList.HandleFunc("/about/{username}", tweetsHandler.GetAllUserTweets)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/tweets", tweetsHandler.PostTweet)

	likeRouter := router.Methods(http.MethodPost).Subrouter()
	likeRouter.HandleFunc("/like/{id}", tweetsHandler.LikeTweet)

	unlikeRouter := router.Methods(http.MethodPost).Subrouter()
	unlikeRouter.HandleFunc("/unlike/{id}", tweetsHandler.UnlikeTweet)

	getUsersWhoLikedTweet := router.Methods(http.MethodGet).Subrouter()
	getUsersWhoLikedTweet.HandleFunc("/likes/{id}", tweetsHandler.GetUsersWhoLikedTweet)

	home := router.Methods(http.MethodGet).Subrouter()
	home.HandleFunc("/home", tweetsHandler.HomeFeed)

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,        // Addr optionally specifies the TCP address for the server to listen on, in the form "host:port". If empty, ":http" (port 80) is used.
		Handler:      router,            // handler to invoke, http.DefaultServeMux if nil
		IdleTimeout:  120 * time.Second, // IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
		ReadTimeout:  1 * time.Second,   // ReadTimeout is the maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout.
		WriteTimeout: 1 * time.Second,   // WriteTimeout is the maximum duration before timing out writes of the response.
	}

	//certFile := "twitter.crt"
	//keyFile := "twitter.key"

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		//err := server.ListenAndServeTLS(certFile, keyFile)
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	//When we receive an interrupt or kill, if we don't have any current connections the code will terminate.
	//But if we do the code will stop receiving any new connections and wait for maximum of 30 seconds to finish all current requests.
	//After that the code will terminate.
	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
