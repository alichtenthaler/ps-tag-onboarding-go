package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/out/mongo"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestUserCreationOK(t *testing.T) {

	u := domain.User{
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

	var respUser domain.User
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

	setUpTestInsertUser(t, domain.User{
		FirstName: firstNameNotUnique,
		LastName: lastNameNotUnique,
		Email: "email@email.com",
		Age: 20,
	})

	testCases := []struct {
		name            string
		user            domain.User
		validationError domain.ValidationError
	}{
		{
			name: "Missing user first name",
			user: domain.User{
				LastName: "ann",
				Email:    "s@s.com",
				Age:      22,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorNameRequired}},
		},
		{
			name: "Missing user last name",
			user: domain.User{
				FirstName: "ann",
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorNameRequired}},
		},
		{
			name: "User minimum age not reached",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "s@s.com",
				Age:       12,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorAgeMinimum}},
		},
		{
			name: "Missing user email",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Age:       22,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorEmailRequired}},
		},
		{
			name: "User wrong email format",
			user: domain.User{
				FirstName: "ann",
				LastName:  "peterson",
				Email:     "ss.com",
				Age:       22,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorEmailFormat}},
		},
		{
			name: "First and lastname are not unique",
			user: domain.User{
				FirstName: firstNameNotUnique,
				LastName:  lastNameNotUnique,
				Email:     "s@s.com",
				Age:       22,
			},
			validationError: domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{domain.ErrorNameUnique}},
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

			var errResp domain.ValidationError
			err = json.Unmarshal(resp, &errResp)
			if err != nil {
				log.Fatal(t, err, resp)
			}

			assert.Equal(tt, testCase.validationError, errResp)
		})
	}
}

func TestUserGetExistingID(t *testing.T) {

	existingID := setUpTestInsertUser(t, domain.User{
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

	var respUser domain.User
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
	assert.Equal(t, domain.ResponseUserNotFound, respError.Error)
}

func setUpTestInsertUser(t *testing.T, existingUser domain.User) primitive.ObjectID {
	res, err := testDBConnection.Collection(mongo.UserCollection).InsertOne(context.Background(), existingUser)
	if err != nil {
		t.Errorf("error inserting user in the database: %s", err.Error())
	}
	return res.InsertedID.(primitive.ObjectID)
}