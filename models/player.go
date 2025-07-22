package models

type Player struct {
	ID      string  `json:"id" bson:"_id,omitempty" db:"id"`
	Name    string  `json:"name" bson:"name" db:"name"`
	Surname string  `json:"surname" bson:"surname" db:"surname"`
	Balance float64 `json:"balance" bson:"balance" db:"balance"`
}
