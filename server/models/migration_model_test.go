package models_test

import (
	"testing"

	"blendermux/server/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MigrationTestSuite struct {
	suite.Suite
	Migration *models.Migration
}

func (suite *MigrationTestSuite) SetupTest() {
	suite.Migration = models.CreateNewMigration(
		"00010101000000",
	)
}

func (suite *MigrationTestSuite) TestCreateNewMigration_CreatesMigrationWithSuppliedFields() {
	//arrange
	timestamp := "this is a timestamp"

	//act
	migration := models.CreateNewMigration(timestamp)

	//assert
	suite.Require().NotNil(migration)
	suite.NotEqual(migration.ID, primitive.NilObjectID)
	suite.EqualValues(migration.Timestamp, timestamp)
}

func (suite *MigrationTestSuite) TestValidate_WithValidMigration_ReturnsModelValid() {
	//act
	err := suite.Migration.Validate()

	//assert
	suite.EqualValues(err.Status, models.ModelValid)
}

func (suite *MigrationTestSuite) TestValidate_WithNilMigrationID_ReturnsMigrationInvalidID() {
	//arrange
	suite.Migration.ID = uuid.Nil

	//act
	err := suite.Migration.Validate()

	//assert
	suite.EqualValues(err.Status, models.MigrationInvalidID)
}

func (suite *MigrationTestSuite) TestValidate_WithVariousInvalidTimestamps_ReturnsError() {
	var timestamp string
	testCase := func() {
		//arrange
		suite.Migration.Timestamp = timestamp

		//act
		err := suite.Migration.Validate()

		//assert
		suite.EqualValues(err.Status, models.MigrationInvalidTimestamp)
	}

	timestamp = "0001010100000"
	suite.Run("TooFewDigits", testCase)

	timestamp = "000101010000000"
	suite.Run("TooManyDigits", testCase)

	timestamp = "000101010a0000"
	suite.Run("ContainsNonDigit", testCase)
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, &MigrationTestSuite{})
}
