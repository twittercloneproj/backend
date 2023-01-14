package social_graph

import (
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	address string
}

func NewClient(host, port string) Client {
	return Client{
		address: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (client Client) CanAccessProfileData(username, token string) (bool, error) {
	bearerToken := fmt.Sprintf("Bearer %s", token)

	requestURL := client.address + fmt.Sprintf("/follow/%s", username)
	httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)
	httpReq.Header.Add("Authorization", bearerToken)

	if err != nil {
		return false, errors.New("req err")
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return false, errors.New("error getting info from social graph")
	}

	return res.StatusCode == http.StatusOK, nil
}
