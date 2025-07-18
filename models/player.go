package models

type Player struct {
	ID      string  `json:"id" bson:"_id,omitempty"`
	Name    string  `json:"name" bson:"name"`
	Surname string  `json:"surname" bson:"surname"`
	Balance float64 `json:"balance" bson:"balance"`
}
