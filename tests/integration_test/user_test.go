package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/adapter/out/mongo"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/service"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserIntegrationTestSuite struct {
	RepositoryTestSuite
}

func (s *UserIntegrationTestSuite) SetupSuite() {
	s.RepositoryTestSuite.SetupSuite()
}

func TestUserIntegrationRun(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
}

func (s *UserIntegrationTestSuite) TestUserCreationOK() {
	u := domain.User{
		FirstName: "Ann",
		LastName:  "Peterson",
		Email:     "a@p.com",
		Age:       20,
	}

	userByte, err := json.Marshal(u)
	if err != nil {
		log.Fatal(s.T(), err, userByte)
	}

	userPersistenceAdapter := mongo.NewUserPersistenceAdapter(s.db)
	createUserService := service.NewCreateUserService(userPersistenceAdapter)
	createUserHandler := web.NewCreateUserHandler(createUserService)

	req := httptest.NewRequest(http.MethodPost, "/fake-path", bytes.NewBuffer(userByte))
	res := httptest.NewRecorder()
	createUserHandler.HandleCreteUser(res, req)

	assert.Exactly(s.T(), http.StatusCreated, res.Code)

	var respUser domain.User
	err = json.Unmarshal(res.Body.Bytes(), &respUser)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	var userSaved domain.User
	err = s.db.Collection(mongo.UserCollection).FindOne(context.Background(), primitive.M{"_id": respUser.ID}).Decode(&userSaved)
	if err != nil {
		s.T().Errorf("error getting user from the database: %s", err.Error())
	}

	assert.Exactly(s.T(), respUser, userSaved)
	assert.Exactly(s.T(), u.FirstName, respUser.FirstName)
	assert.Exactly(s.T(), u.LastName, respUser.LastName)
	assert.Exactly(s.T(), u.Email, respUser.Email)
	assert.Exactly(s.T(), u.Age, respUser.Age)
}

func (s *UserIntegrationTestSuite) TestUserCreationValidationFails() {
	ctx := context.Background()

	firstNameNotUnique := "TestUserCreationValidationFails"
	lastNameNotUnique := "ErrorNameUnique"

	setUpTestInsertUser(ctx, s, domain.User{
		FirstName: firstNameNotUnique,
		LastName:  lastNameNotUnique,
		Email:     "email@email.com",
		Age:       20,
	})

	testCases := []struct {
		name            string
		user            domain.User
		validationError errs.ValidationError
	}{
		{
			name: "Missing user first name",
			user: domain.User{
				LastName: "ann",
				Email:    "s@s.com",
				Age:      22,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorNameRequired.Message}},
		},
		{
			name: "Missing user last name",
			user: domain.User{
				FirstName: "ann",
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorNameRequired.Message}},
		},
		{
			name: "User minimum age not reached",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "s@s.com",
				Age:       12,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorAgeMinimum.Message}},
		},
		{
			name: "Missing user email",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Age:       22,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorEmailRequired.Message}},
		},
		{
			name: "User wrong email format",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "ss.com",
				Age:       22,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorEmailFormat.Message}},
		},
		{
			name: "First and lastname are not unique",
			user: domain.User{
				FirstName: firstNameNotUnique,
				LastName:  lastNameNotUnique,
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{errs.ErrorNameUnique.Message}},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.name, func(tt *testing.T) {

			userByte, err := json.Marshal(testCase.user)
			if err != nil {
				log.Fatal(s.T(), err, userByte)
			}

			userPersistenceAdapter := mongo.NewUserPersistenceAdapter(s.db)
			createUserService := service.NewCreateUserService(userPersistenceAdapter)
			createUserHandler := web.NewCreateUserHandler(createUserService)

			req := httptest.NewRequest(http.MethodPost, "/fake-path", bytes.NewBuffer(userByte))
			res := httptest.NewRecorder()
			createUserHandler.HandleCreteUser(res, req)

			assert.Exactly(s.T(), http.StatusBadRequest, res.Code)

			var errResp errs.ValidationError
			err = json.Unmarshal(res.Body.Bytes(), &errResp)
			if err != nil {
				log.Fatal(s.T(), err, res.Body.String())
			}

			assert.Equal(tt, testCase.validationError, errResp)
		})
	}
}

func (s *UserIntegrationTestSuite) TestUserGetExistingID() {
	ctx := context.Background()

	existingID := setUpTestInsertUser(ctx, s, domain.User{
		FirstName: "TestUserGetExistingID",
		LastName:  "Lastname",
		Email:     "t@l.com",
		Age:       20,
	})

	userPersistenceAdapter := mongo.NewUserPersistenceAdapter(s.db)
	getUserService := service.NewGetUserService(userPersistenceAdapter)
	getUserHandler := web.NewGetUserHandler(getUserService)

	req := httptest.NewRequest(http.MethodGet, "/fake-path", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": existingID.Hex()})
	res := httptest.NewRecorder()
	getUserHandler.HandleGetUser(res, req)

	assert.Exactly(s.T(), http.StatusOK, res.Code)

	var respUser domain.User
	err := json.Unmarshal(res.Body.Bytes(), &respUser)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	assert.Exactly(s.T(), existingID, respUser.ID)
}

func (s *UserIntegrationTestSuite) TestUserGetNotExistingID() {
	notExistingID := "a"

	userPersistenceAdapter := mongo.NewUserPersistenceAdapter(s.db)
	getUserService := service.NewGetUserService(userPersistenceAdapter)
	getUserHandler := web.NewGetUserHandler(getUserService)

	req := httptest.NewRequest(http.MethodGet, "/fake-path", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": notExistingID})
	res := httptest.NewRecorder()
	getUserHandler.HandleGetUser(res, req)

	var respError errs.GenericError
	err := json.Unmarshal(res.Body.Bytes(), &respError)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	assert.Exactly(s.T(), http.StatusNotFound, res.Code)
	assert.Equal(s.T(), errs.ResponseUserNotFound.Error(), respError.Error())
}

func setUpTestInsertUser(ctx context.Context, s *UserIntegrationTestSuite, existingUser domain.User) primitive.ObjectID {
	res, err := s.db.Collection(mongo.UserCollection).InsertOne(ctx, existingUser)
	if err != nil {
		s.T().Errorf("error inserting user in the database: %s", err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}