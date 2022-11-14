package data

import (
	"encoding/json"
	"io"
	"time"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedOn time.Time `json:"created_on"`
	User      string    `json:"user"`
}
type Tweets struct {
	tweets []Tweet
}

func (p *Tweet) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
