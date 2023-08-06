package data

type User struct {
	Username string  `json:"username" validate:"required"`
	Sex      string  `json:"sex"`
	Age      float32 `json:"age"`
	Town     string  `json:"town"`
	Privacy  string  `json:"privacy"`
}

type ApproveRequest struct {
	Approved bool `json:"approved"`
}

type UpdatePrivacy struct {
	Privacy string `bson:"privacy,omitempty" json:"privacy"`
}
