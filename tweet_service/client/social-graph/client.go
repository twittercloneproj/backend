package social_graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tweet_service/data"
)

type Client struct {
	address string
}

func NewClient(host, port string) Client {
	return Client{
		address: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (client Client) GetFollowers(username string) ([]string, error) {
	requestURL := client.address + fmt.Sprintf("/followers/%s", username)
	httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		return []string{}, errors.New("req err")
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return []string{}, errors.New("error getting info from social graph")
	}
	defer res.Body.Close()

	var followers []data.User
	var usernames []string

	err = json.NewDecoder(res.Body).Decode(&followers)
	if err != nil {
		return []string{}, errors.New("error decoding body")
	}

	for _, user := range followers {
		usernames = append(usernames, user.Username)
	}

	return usernames, nil
}
