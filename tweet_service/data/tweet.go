package data

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

type Tweet struct {
	ID               gocql.UUID `json:"id"`
	Text             string     `json:"text"`
	PostedBy         string     `json:"posted_by"`
	Retweet          bool       `json:"retweet"`
	OriginalPostedBy string     `json:"original_posted_by"`
}

type Likes struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
}
type Tweets struct {
	tweets []Tweet
}

type User struct {
	Username string  `json:"username" validate:"required"`
	Sex      string  `json:"sex"`
	Age      float32 `json:"age"`
	Town     string  `json:"town"`
	Privacy  string  `json:"privacy"`
}

func (p *Tweet) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (o *Tweet) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
