package social_graph

import (
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"time"
)

type SocialGraphCircuitBreaker struct {
	circuitBreaker *gobreaker.CircuitBreaker
	client         *Client
}

func NewCircuitBreaker(socialGraphClient *Client) *SocialGraphCircuitBreaker {
	return &SocialGraphCircuitBreaker{
		circuitBreaker: CircuitBreaker(),
		client:         socialGraphClient,
	}
}

func CircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "SocialGraph",
			MaxRequests: 1,
			Timeout:     5 * time.Second,
			Interval:    0,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 0
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit Breaker '%s' changed from '%s' to '%s'\n", name, from, to)
			},
		},
	)
}

func (cb *SocialGraphCircuitBreaker) CanAccessTweet(username, token string) (bool, error) {
	execute, err := cb.circuitBreaker.Execute(func() (interface{}, error) {
		response, err := cb.client.CanAccessTweet(username, token)

		log.Println("SG Client Err: ", err)

		return response, err
	})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return execute.(bool), nil
}
