package social_graph

import (
	"auth_service/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (client Client) CreateUser(user *data.User) error {
	user.Privacy = "Private"
	reqBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(reqBytes)
	requestURL := client.address + "/user"
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if err != nil {
		log.Println(err)
		return errors.New("error creating user")
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		log.Println(res.StatusCode)
		return errors.New("error creating user")
	}
	return nil
}
