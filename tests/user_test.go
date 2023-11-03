package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	user2 "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestUserCreationOK(t *testing.T) {

	u := user2.User{
		FirstName: "Ann",
		LastName: "Peterson",
		Email: "a@p.com",
		Age: 20,
	}

	userByte, err := json.Marshal(u)
	if err != nil {
		log.Fatal(t, err, userByte)
	}

	code, resp, err := httpTool.POST("save", userByte)
	if err != nil {
		panic(err)
	}

	assert.Exactly(t, http.StatusCreated, code)

	var respUser user2.User
	err = json.Unmarshal(resp, &respUser)
	if err != nil {
		log.Fatal(t, err, resp)
	}

	assert.Exactly(t, u.FirstName, respUser.FirstName)
	assert.Exactly(t, u.LastName, respUser.LastName)
	assert.Exactly(t, u.Email, respUser.Email)
	assert.Exactly(t, u.Age, respUser.Age)
}

func TestUserCreationValidationFails(t *testing.T) {

	firstNameNotUnique := "TestUserCreationValidationFails"
	lastNameNotUnique := "ErrorNameUnique"

	setUpTestInsertUser(t, user2.User{
		FirstName: firstNameNotUnique,
		LastName: lastNameNotUnique,
		Email: "email@email.com",
		Age: 20,
	})

	testCases := []struct {
		name            string
		user            user2.User
		validationError web.ValidationError
	}{
		{
			name: "Missing user first name",
			user: user2.User{
				LastName: "ann",
				Email:    "s@s.com",
				Age:      22,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorNameRequired}},
		},
		{
			name: "Missing user last name",
			user: user2.User{
				FirstName: "ann",
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorNameRequired}},
		},
		{
			name: "User minimum age not reached",
			user: user2.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "s@s.com",
				Age:       12,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorAgeMinimum}},
		},
		{
			name: "Missing user email",
			user: user2.User{
				FirstName: "ann",
				LastName:  "peterson",
				Age:       22,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorEmailRequired}},
		},
		{
			name: "User wrong email format",
			user: user2.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "ss.com",
				Age:       22,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorEmailFormat}},
		},
		{
			name: "First and lastname are not unique",
			user: user2.User{
				FirstName: firstNameNotUnique,
				LastName:  lastNameNotUnique,
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: web.ValidationError{Error: user2.ResponseValidationFailed, Details: []string{user2.ErrorNameUnique}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {

			userByte, err := json.Marshal(testCase.user)
			if err != nil {
				log.Fatal(t, err, userByte)
			}

			code, resp, err := httpTool.POST("save", userByte)
			if err != nil {
				panic(err)
			}

			assert.Exactly(t, http.StatusBadRequest, code)

			var errResp web.ValidationError
			err = json.Unmarshal(resp, &errResp)
			if err != nil {
				log.Fatal(t, err, resp)
			}

			assert.Equal(tt, testCase.validationError, errResp)
		})
	}
}

func TestUserGetExistingID(t *testing.T) {

	existingID := setUpTestInsertUser(t, user2.User{
			FirstName: "TestUserGetExistingID",
			LastName: "Lastname",
			Email: "t@l.com",
			Age: 20,
		})

	code, resp, err := httpTool.GET(fmt.Sprintf("find/%s", existingID.Hex()))
	if err != nil {
		panic(err)
	}

	assert.Exactly(t, http.StatusOK, code)

	var respUser user2.User
	err = json.Unmarshal(resp, &respUser)
	if err != nil {
		log.Fatal(t, err, resp)
	}

	assert.Exactly(t, existingID, respUser.ID)
}

func TestUserGetNotExistingID(t *testing.T) {
	notExistingID := "a"

	code, resp, err := httpTool.GET(fmt.Sprintf("find/%s", notExistingID))
	if err != nil {
		panic(err)
	}

	var respError web.GenericError
	err = json.Unmarshal(resp, &respError)
	if err != nil {
		log.Fatal(t, err, resp)
	}

	assert.Exactly(t, http.StatusNotFound, code)
	assert.Equal(t, user2.ResponseUserNotFound, respError.Error)
}

func setUpTestInsertUser(t *testing.T, existingUser user2.User) primitive.ObjectID {
	res, err := testDBConnection.Collection(user.UserCollection).InsertOne(context.Background(), existingUser)
	if err != nil {
		t.Errorf("error inserting user in the database: %s", err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}