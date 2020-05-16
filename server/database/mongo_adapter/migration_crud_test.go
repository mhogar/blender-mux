package mongoadapter_test

import (
	"log"
	"testing"

	"blendermux/common"
	"blendermux/server/config"
	mongoadapter "blendermux/server/database/mongo_adapter"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	suite.Suite
	DB *mongoadapter.MongoAdapter
}

func (suite *MigrationCRUDTestSuite) SetupSuite() {
	config.InitConfig()

	suite.DB = &mongoadapter.MongoAdapter{
		DbKey: "integration",
	}

	err := suite.DB.OpenConnection()
	if err != nil {
		log.Fatal(common.ChainError("error openning database connection", err))
	}
}

func (suite *MigrationCRUDTestSuite) SetupTest() {
	//drop the db to start with a fresh one before every test
	ctx, cancel := suite.DB.CreateStandardTimeoutContext()
	err := suite.DB.Database.Drop(ctx)
	cancel()

	if err != nil {
		log.Fatal(common.ChainError("error droping database", err))
	}
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//arrange
	timestamp := "invalid"

	//act
	err := suite.DB.CreateMigration(timestamp)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "timestamp is invalid")
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_MigrationIsAddedToDatabase() {
	timestamp := "10000000000000"

	//verify no matching timestamps to start with
	ctx, cancel := suite.DB.CreateStandardTimeoutContext()
	count, err := suite.DB.Migrations.CountDocuments(ctx, bson.D{{Key: "timestamp", Value: timestamp}})
	cancel()

	suite.Require().NoError(err)
	suite.EqualValues(count, 0)

	//insert migration and verify no err
	err = suite.DB.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//verify a matching timestamp can now be found after insert
	ctx, cancel = suite.DB.CreateStandardTimeoutContext()
	count, err = suite.DB.Migrations.CountDocuments(ctx, bson.D{{Key: "timestamp", Value: timestamp}})
	cancel()

	suite.Require().NoError(err)
	suite.EqualValues(count, 1)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_MigrationIsRemovedFromDatabase() {
	timestamp := "10000000000000"

	//populate the database with a migration
	err := suite.DB.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//get the number matching timestamps for comparison later
	ctx, cancel := suite.DB.CreateStandardTimeoutContext()
	oldCount, err := suite.DB.Migrations.CountDocuments(ctx, bson.D{{Key: "timestamp", Value: timestamp}})
	cancel()
	suite.Require().NoError(err)

	//delete migration and verify no err
	err = suite.DB.DeleteMigrationByTimestamp(timestamp)
	suite.Require().NoError(err)

	//verify one less matching timestamp after delete
	ctx, cancel = suite.DB.CreateStandardTimeoutContext()
	newCount, err := suite.DB.Migrations.CountDocuments(ctx, bson.D{{Key: "timestamp", Value: timestamp}})
	cancel()

	suite.Require().NoError(err)
	suite.EqualValues(newCount, oldCount-1)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_WithNoTimestampsInDatabase_ReturnsNoHasLatest() {
	//verify database has no timestamps in it
	ctx, cancel := suite.DB.CreateStandardTimeoutContext()
	count, err := suite.DB.Migrations.CountDocuments(ctx, bson.D{})
	cancel()

	suite.Require().NoError(err)
	suite.EqualValues(count, 0)

	//act
	_, hasLatest, err := suite.DB.GetLatestTimestamp()

	//assert
	suite.Require().NoError(err)
	suite.False(hasLatest)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_LatestTimestampIsReturned() {
	//arrange
	timestamps := []string{"10000000000000", "30000000000000", "20000000000000"}

	//populate the database with a few migrations
	for _, timestamp := range timestamps {
		err := suite.DB.CreateMigration(timestamp)
		suite.Require().NoError(err)
	}

	//act
	timestamp, hasLatest, err := suite.DB.GetLatestTimestamp()

	//assert
	suite.Require().NoError(err)
	suite.EqualValues(timestamp, timestamps[1])
	suite.True(hasLatest)
}

func (suite *MigrationCRUDTestSuite) TearDownSuite() {
	suite.DB.CloseConnection()
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
