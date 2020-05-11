package migrationrunner_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	migrationrunner "github.com/blendermux/common/migration_runner"
	"github.com/blendermux/common/migration_runner/mocks"

	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	MigrationRepositoryMock mocks.MigrationRepository
	MigrationCRUDMock       mocks.MigrationCRUD
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.MigrationRepositoryMock = mocks.MigrationRepository{}
	suite.MigrationCRUDMock = mocks.MigrationCRUD{}
}

func (suite *MigrationRunnerTestSuite) TestRunMigrations_WithNoLatestTimestamp_RunsAllMigrations() {
	//arrange
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("100: this can be anything", false, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := migrationrunner.RunMigrations(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.NoError(err)

	migrationMocks[0].AssertCalled(suite.T(), "Up")
	migrationMocks[1].AssertCalled(suite.T(), "Up")
	migrationMocks[2].AssertCalled(suite.T(), "Up")
	migrationMocks[3].AssertCalled(suite.T(), "Up")

	suite.MigrationCRUDMock.AssertNumberOfCalls(suite.T(), "CreateMigration", 4)
}

func (suite *MigrationRunnerTestSuite) TestRunMigrations_RunsAllMigrationsWithTimestampGreaterThanLatest() {
	//arrange
	migrationMocks := createMigrationMocks("01", "04", "08", "10")

	migrations := make([]migrationrunner.Migration, len(migrationMocks))
	for i, _ := range migrationMocks {
		migrations[i] = &migrationMocks[i]
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("05", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(nil)

	//act
	err := migrationrunner.RunMigrations(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.NoError(err)

	migrationMocks[0].AssertNotCalled(suite.T(), "Up")
	migrationMocks[1].AssertNotCalled(suite.T(), "Up")
	migrationMocks[2].AssertCalled(suite.T(), "Up")
	migrationMocks[3].AssertCalled(suite.T(), "Up")

	suite.MigrationCRUDMock.AssertNumberOfCalls(suite.T(), "CreateMigration", 2)
}

func (suite *MigrationRunnerTestSuite) TestRunMigrations_WhereGetLatestTimestampReturnsError_ErrorIsReturned() {
	//arrange
	errMessage := "GetLatestTimestamp mock error"

	suite.MigrationRepositoryMock.On("GetMigrations").Return(nil)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("", false, errors.New(errMessage))

	//act
	err := migrationrunner.RunMigrations(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestRunMigrations_WhereMigrationUpReturnsError_ErrorIsReturned() {
	//arrange
	errMessage := "Up mock error"

	migrationMock := mocks.Migration{}
	migrationMock.On("GetTimestamp").Return("1")
	migrationMock.On("Up").Return(errors.New(errMessage))

	migrations := []migrationrunner.Migration{
		&migrationMock,
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)

	//act
	err := migrationrunner.RunMigrations(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func (suite *MigrationRunnerTestSuite) TestRunMigrations_WhereCreateMigrationReturnsError_ErrorIsReturned() {
	//arrange
	errMessage := "CreateMigration mock error"

	migrationMock := mocks.Migration{}
	migrationMock.On("GetTimestamp").Return("1")
	migrationMock.On("Up").Return(nil)

	migrations := []migrationrunner.Migration{
		&migrationMock,
	}

	suite.MigrationRepositoryMock.On("GetMigrations").Return(migrations)
	suite.MigrationCRUDMock.On("GetLatestTimestamp").Return("0", true, nil)
	suite.MigrationCRUDMock.On("CreateMigration", mock.Anything).Return(errors.New(errMessage))

	//act
	err := migrationrunner.RunMigrations(&suite.MigrationRepositoryMock, &suite.MigrationCRUDMock)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), errMessage)
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}

func createMigrationMocks(timestamps ...string) []mocks.Migration {
	migrationMocks := make([]mocks.Migration, len(timestamps))

	for i, timestamp := range timestamps {
		migrationMocks[i] = mocks.Migration{}
		migrationMocks[i].On("GetTimestamp").Return(timestamp)
		migrationMocks[i].On("Up").Return(nil)
	}

	return migrationMocks
}
