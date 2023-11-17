package mongo

import (
	"context"
	"errors"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// UserCollection is the name of the collection in mongo database
	UserCollection = "user"
)

// UserPersistenceAdapter implements GetUserPort and SaveUserPort interface
type UserPersistenceAdapter struct {
	db *mongo.Database
}

// NewUserPersistenceAdapter creates a new instance of mongo adapter
func NewUserPersistenceAdapter(db *mongo.Database) *UserPersistenceAdapter {
	return &UserPersistenceAdapter{
		db: db,
	}
}

// GetUserById find a user by id from the database if it exists or returns an error if it does not
func (repo UserPersistenceAdapter) GetUserById(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	objectID, _ := primitive.ObjectIDFromHex(id)

	err := repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

// SaveUser saves a user to the database
func (repo UserPersistenceAdapter) SaveUser(ctx context.Context, user domain.User) (primitive.ObjectID, error)  {
	res, err := repo.db.Collection(UserCollection).InsertOne(ctx, user)
	return res.InsertedID.(primitive.ObjectID), err
}

// ExistsByFirstNameAndLastName checks if a user with the given first and last name exists in the database
func (repo UserPersistenceAdapter) ExistsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {

	if errors.Is(repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"firstname": firstName, "lastname": lastName}).Err(), mongo.ErrNoDocuments) {
		return false
	}

	return true
}