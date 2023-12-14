package integration_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	mongoImage  = "mongo:7.0.2"
	exposedPort = "27017"
	dbName      = "usertest"
)

type mongodbContainer struct {
	testcontainers.Container
}

func newMongoContainer(ctx context.Context) (*mongodbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        mongoImage,
		ExposedPorts: []string{exposedPort},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort(exposedPort),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mongodbContainer{Container: container}, nil
}

// RepositoryTestSuite provides a reusable suite.Suite that spins up a MongoDB container, cleans up the database
// after each test case execution. This can be embedded in other suites easily.
type RepositoryTestSuite struct {
	suite.Suite
	container *mongodbContainer
	db        *mongo.Database
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.TODO()

	var err error
	s.container, err = newMongoContainer(ctx)
	require.NoError(s.T(), err)

	endpoint, err := s.container.Endpoint(ctx, "mongodb")
	require.NoError(s.T(), err)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	require.NoError(s.T(), err)

	require.Eventually(s.T(), func() bool {
		err = mongoClient.Ping(ctx, nil)
		return err == nil
	}, time.Second*10, time.Millisecond*100)

	s.db = mongoClient.Database(dbName)
}

// The TearDownTest method will be run after every test in the suite.
func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.TODO()

	_ = s.db.Drop(ctx)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (s *RepositoryTestSuite) TearDownSuite() {
	ctx := context.TODO()

	_ = s.container.Terminate(ctx)
}
