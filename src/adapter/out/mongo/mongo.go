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
	UserCollection = "user"
)
type UserPersistenceAdapter struct {
	db *mongo.Database
}

func NewUserPersistenceAdapter(db *mongo.Database) *UserPersistenceAdapter {
	return &UserPersistenceAdapter{
		db: db,
	}
}

func (repo UserPersistenceAdapter) GetUserById(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	objectID, _ := primitive.ObjectIDFromHex(id)

	err := repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

func (repo UserPersistenceAdapter) SaveUser(ctx context.Context, user domain.User) (primitive.ObjectID, error)  {
	res, err := repo.db.Collection(UserCollection).InsertOne(ctx, user)
	return res.InsertedID.(primitive.ObjectID), err
}

func (repo UserPersistenceAdapter) ExistsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {

	if errors.Is(repo.db.Collection(UserCollection).FindOne(ctx, bson.M{"firstname": firstName, "lastname": lastName}).Err(), mongo.ErrNoDocuments) {
		return false
	}

	return true
}