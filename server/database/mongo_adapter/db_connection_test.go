package mongoadapter_test

import (
	"blendermux/server/config"
	mongoadapter "blendermux/server/database/mongo_adapter"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type DbConnectionTestSuite struct {
	suite.Suite
	DB *mongoadapter.MongoAdapter
}

func (suite *DbConnectionTestSuite) SetupTest() {
	viper.Reset()
	config.InitConfig()

	suite.DB = &mongoadapter.MongoAdapter{
		DbKey: "integration",
	}
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereEnvironmentIsNotFound_ReturnsError() {
	//arrange
	env := "not a real environment"
	viper.Set("env", env)

	//act
	err := suite.DB.OpenConnection()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), env)
}

func (suite *DbConnectionTestSuite) TestPing_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.DB.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.DB.Ping()

	//assert
	suite.NoError(err)
}

func TestDbConnectionTestSuite(t *testing.T) {
	suite.Run(t, &DbConnectionTestSuite{})
}
