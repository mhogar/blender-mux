package models_test

import (
	"testing"

	"blendermux/server/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	User *models.User
}

func (suite *UserTestSuite) SetupTest() {
	suite.User = models.CreateNewUser(
		"email@email.com",
		[]byte("password"),
	)
}

func (suite *UserTestSuite) TestCreateNewUser_CreatesUserWithSuppliedFields() {
	//arrange
	email := "this is a test email"
	hash := []byte("this is a password")

	//act
	user := models.CreateNewUser(email, hash)

	//assert
	suite.Require().NotNil(user)
	suite.NotEqual(user.ID, uuid.Nil)
	suite.EqualValues(user.Email, email)
	suite.EqualValues(user.PasswordHash, hash)
}

func (suite *UserTestSuite) TestValidate_WithValidUser_ReturnsModelValid() {
	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(err.Status, models.ModelValid)
}

func (suite *UserTestSuite) TestValidate_WithNilUserID_ReturnsUserInvalidID() {
	//arrange
	suite.User.ID = uuid.Nil

	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(err.Status, models.UserInvalidID)
}

func (suite *UserTestSuite) TestValidate_WithVariousInvalidEmails_ReturnsUserInvalidEmail() {
	var email string
	testCase := func() {
		//arrange
		suite.User.Email = email

		//act
		err := suite.User.Validate()

		//assert
		suite.EqualValues(err.Status, models.UserInvalidEmail)
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

func (suite *UserTestSuite) TestValidate_WithEmptyPasswordHash_ReturnsUserInvalidPasswordHash() {
	//arrange
	suite.User.PasswordHash = make([]byte, 0)

	//act
	err := suite.User.Validate()

	//assert
	suite.EqualValues(err.Status, models.UserInvalidPasswordHash)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{})
}
