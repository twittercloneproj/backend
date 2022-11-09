package data

import (
	"encoding/json"
	"io"
)

type Tweet struct {
	// NoSQL: we don't want to keep track of our int ID, so we use UUID
	ID        string `json:"id"`                       //specifies that in the incoming Body the field to map to this will be called "id"
	Text      string `json:"text" validate:"required"` //there are some integrated validation, for eg. this specifies that a value for name must be provided, otherwise it will not be valid
	CreatedOn string `json:"createdOn"`
}

type Tweets []*Tweet

func (p *Tweets) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Tweet) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Tweet) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
