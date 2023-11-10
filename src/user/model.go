package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user model
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Age       int                `json:"age"`
}
