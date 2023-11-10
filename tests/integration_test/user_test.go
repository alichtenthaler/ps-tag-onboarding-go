package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/response"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
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
	u := user.User{
		FirstName: "Ann",
		LastName:  "Peterson",
		Email:     "a@p.com",
		Age:       20,
	}

	userByte, err := json.Marshal(u)
	if err != nil {
		log.Fatal(s.T(), err, userByte)
	}

	userService := user.New(s.db)
	req := httptest.NewRequest(http.MethodPost, "/fake-path", bytes.NewBuffer(userByte))
	res := httptest.NewRecorder()
	userService.CreateUser(res, req)

	assert.Exactly(s.T(), http.StatusCreated, res.Code)

	var respUser user.User
	err = json.Unmarshal(res.Body.Bytes(), &respUser)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	assert.Exactly(s.T(), u.FirstName, respUser.FirstName)
	assert.Exactly(s.T(), u.LastName, respUser.LastName)
	assert.Exactly(s.T(), u.Email, respUser.Email)
	assert.Exactly(s.T(), u.Age, respUser.Age)
}

func (s *UserIntegrationTestSuite) TestUserCreationValidationFails() {
	ctx := context.Background()

	firstNameNotUnique := "TestUserCreationValidationFails"
	lastNameNotUnique := "ErrorNameUnique"

	setUpTestInsertUser(ctx, s, user.User{
		FirstName: firstNameNotUnique,
		LastName:  lastNameNotUnique,
		Email:     "email@email.com",
		Age:       20,
	})

	testCases := []struct {
		name            string
		user            user.User
		validationError response.ValidationError
	}{
		{
			name: "Missing user first name",
			user: user.User{
				LastName: "ann",
				Email:    "s@s.com",
				Age:      22,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorNameRequired}},
		},
		{
			name: "Missing user last name",
			user: user.User{
				FirstName: "ann",
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorNameRequired}},
		},
		{
			name: "User minimum age not reached",
			user: user.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "s@s.com",
				Age:       12,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorAgeMinimum}},
		},
		{
			name: "Missing user email",
			user: user.User{
				FirstName: "ann",
				LastName:  "peterson",
				Age:       22,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorEmailRequired}},
		},
		{
			name: "User wrong email format",
			user: user.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "ss.com",
				Age:       22,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorEmailFormat}},
		},
		{
			name: "First and lastname are not unique",
			user: user.User{
				FirstName: firstNameNotUnique,
				LastName:  lastNameNotUnique,
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: response.ValidationError{Error: user.ResponseValidationFailed, Details: []string{user.ErrorNameUnique}},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.name, func(tt *testing.T) {

			userByte, err := json.Marshal(testCase.user)
			if err != nil {
				log.Fatal(s.T(), err, userByte)
			}

			userService := user.New(s.db)
			req := httptest.NewRequest(http.MethodPost, "/fake-path", bytes.NewBuffer(userByte))
			res := httptest.NewRecorder()
			userService.CreateUser(res, req)

			assert.Exactly(s.T(), http.StatusBadRequest, res.Code)

			var errResp response.ValidationError
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

	existingID := setUpTestInsertUser(ctx, s, user.User{
		FirstName: "TestUserGetExistingID",
		LastName:  "Lastname",
		Email:     "t@l.com",
		Age:       20,
	})

	userIdURLParam := existingID.Hex()

	userService := user.New(s.db)
	req := httptest.NewRequest(http.MethodGet, "/fake-path", nil)
	res := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"userId": userIdURLParam})
	userService.FindUserById(res, req)

	assert.Exactly(s.T(), http.StatusOK, res.Code)

	var respUser user.User
	err := json.Unmarshal(res.Body.Bytes(), &respUser)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	assert.Exactly(s.T(), existingID, respUser.ID)
}

func (s *UserIntegrationTestSuite) TestUserGetNotExistingID() {
	notExistingID := "a"

	userService := user.New(s.db)
	req := httptest.NewRequest(http.MethodGet, "/fake-path", nil)
	res := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"userId": notExistingID})
	userService.FindUserById(res, req)

	var respError response.GenericError
	err := json.Unmarshal(res.Body.Bytes(), &respError)
	if err != nil {
		log.Fatal(s.T(), err, res.Body.String())
	}

	assert.Exactly(s.T(), http.StatusNotFound, res.Code)
	assert.Equal(s.T(), user.ResponseUserNotFound, respError.Error)
}

func setUpTestInsertUser(ctx context.Context, s *UserIntegrationTestSuite, existingUser user.User) primitive.ObjectID {
	res, err := s.db.Collection(user.UserCollection).InsertOne(ctx, existingUser)
	if err != nil {
		s.T().Errorf("error inserting user in the database: %s", err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}
