package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// UserCollection is the name of the collection in mongo database
	UserCollection = "user"
)

// Repo implements Repository
type Repo struct {
	db *mongo.Database
}

func newRepository(db *mongo.Database) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo Repo) create(ctx context.Context, user User) (primitive.ObjectID, error) {

	res, err := repo.db.Collection(UserCollection).InsertOne(ctx, user)
	return res.InsertedID.(primitive.ObjectID), err
}

func (repo Repo) getByID(ctx context.Context, id primitive.ObjectID) (User, error) {

	var user User

	err := repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

func (repo Repo) existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {

	if errors.Is(repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"firstname": firstName, "lastname": lastName}).Err(), mongo.ErrNoDocuments) {
		return false
	}

	return true
}
