package models_test

import (
	"testing"

	"github.com/blendermux/server/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	User *models.User
}

func (suite *UserTestSuite) SetupTest() {
	suite.User = &models.User{
		uuid.New(),
		"email@email.com",
		[]byte("password"),
	}
}

func (suite *UserTestSuite) TestValidate_WithValidUser_ReturnsNilError() {
	//act
	err := suite.User.Validate()

	//assert
	suite.Nil(err)
}

func (suite *UserTestSuite) TestValidate_WithNullUserID_ReturnsError() {
	//arrange
	suite.User.ID = uuid.Nil

	//act
	err := suite.User.Validate()

	//assert
	suite.Require().NotNil(err)
	suite.Equal(err.Status, models.UserInvalidID)
}

func (suite *UserTestSuite) TestValidate_WithVariousInvalidEmails_ReturnsError() {
	var email string
	testCase := func() {
		//arrange
		suite.User.Email = email

		//act
		err := suite.User.Validate()

		//assert
		suite.Require().NotNil(err)
		suite.Equal(err.Status, models.UserInvalidEmail)
	}

	email = "@domain.ca"
	suite.Run("NoUser", testCase)

	email = "test?@domain.ca"
	suite.Run("UserContainsInvalidChars", testCase)

	email = "domain.ca"
	suite.Run("No@", testCase)

	email = "test@"
	suite.Run("NoDomain", testCase)

	email = "test@domain?.ca"
	suite.Run("DomainContainsInvalidChars", testCase)

	email = "test@domain"
	suite.Run("NoTopLevelDomain", testCase)

	email = "test@domain.a"
	suite.Run("TopLevelDomainTooShort", testCase)
}

func (suite *UserTestSuite) TestValidate_WithEmptyPasswordHash_ReturnsError() {
	//arrange
	suite.User.PasswordHash = make([]byte, 0)

	//act
	err := suite.User.Validate()

	//assert
	suite.Require().NotNil(err)
	suite.Equal(err.Status, models.UserInvalidPasswordHash)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{})
}
