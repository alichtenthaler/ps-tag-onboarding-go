package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UserCollection = "user"
)

func (up *Processor) create(user User) (primitive.ObjectID, error) {

	res, err := up.db.Collection(UserCollection).InsertOne(context.TODO(), user)
	return res.InsertedID.(primitive.ObjectID), err
}

func (up *Processor) getByID(ctx context.Context, id primitive.ObjectID) (User, error) {

	var user User

	err := up.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

func (up *Processor) existsByFirstNameAndLastName(firstName, lastName string) bool {

	if errors.Is(up.db.Collection(UserCollection).FindOne(context.TODO(), bson.M{"firstname": firstName, "lastname": lastName}).Err(), mongo.ErrNoDocuments) {
		return false
	}

	return true
}
