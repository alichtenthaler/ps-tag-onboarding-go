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

type UserRepositoryI interface {
	create(ctx context.Context, user User) (primitive.ObjectID, error)
	getByID(ctx context.Context, id primitive.ObjectID) (User, error)
	existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool
}

type UserRepository struct {
	db *mongo.Database
}

func newRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo UserRepository) create(ctx context.Context, user User) (primitive.ObjectID, error) {

	res, err := repo.db.Collection(UserCollection).InsertOne(ctx, user)
	return res.InsertedID.(primitive.ObjectID), err
}

func (repo UserRepository) getByID(ctx context.Context, id primitive.ObjectID) (User, error) {

	var user User

	err := repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

func (repo UserRepository) existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {

	if errors.Is(repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"firstname": firstName, "lastname": lastName}).Err(), mongo.ErrNoDocuments) {
		return false
	}

	return true
}
