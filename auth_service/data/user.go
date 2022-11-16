package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" `
	Surname  string             `bson:"surname" json:"surname"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Gender   string             `bson:"gender" json:"gender"`
	Age      float32            `bson:"age" json:"age"`
	Town     string             `bson:"town" json:"town"`

	//Price     float32 `json:"price" validate:"gt=0"`
	//SKU       string  `json:"sku" validate:"required,sku"` //the tag "sku" is there so we can add custom validation
	//CreatedOn string  `json:"createdOn"`
	//UpdatedOn string  `json:"updatedOn"`
	//DeletedOn string  `json:"deletedOn"`
	//// NoSQL: Bonus exercise - type field
	//Type string `json:"type" validate:"required"`
}

type Users []*User

func (p *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
