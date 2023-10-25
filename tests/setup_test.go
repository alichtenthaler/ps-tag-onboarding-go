package main

import (
	"context"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/config"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/database"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/tests/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

var httpTool tools.HTTPTool
var testDBConnection *mongo.Database

func TestMain(m *testing.M) {
	httpTool.Init("http://localhost", "8080", 30*time.Second)
	err := databaseSetup()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func databaseSetup() error {
	var err error
	config.Load()
	fmt.Println("Connecting to the database...")
	testDBConnection, err = createTestDBConnection()
	if err != nil {
		fmt.Printf("Failed to connect database: %s\n", err)
		return err
	}
	fmt.Println("Connected to the database...")
	fmt.Println("Cleaning database before running tests...")
	_, err = testDBConnection.Collection(user.UserCollection).DeleteMany(context.Background(), bson.M{})
	if err != nil {
		fmt.Printf("Database setup failed: %s\n", err)
		return err
	}

	return err
}

func createTestDBConnection() (*mongo.Database, error) {
	db, err := database.Connect(fmt.Sprintf("mongodb://%s:%s@%s:%s",
			config.DBUser,
			config.DBPass,
			os.Getenv("DB_HOST_LOCAL"),
			config.DBPort,
		))
	if err != nil {
		return nil, err
	}
	return db, nil
}
