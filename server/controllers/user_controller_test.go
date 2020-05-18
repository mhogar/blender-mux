package controllers_test

import (
	"blendermux/server/controllers"
	databasemocks "blendermux/server/database/mocks"
	"blendermux/server/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	UserCRUDMock   databasemocks.UserCRUD
	UserController controllers.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.UserCRUDMock = databasemocks.UserCRUD{}
	suite.UserController = controllers.UserController{UserCRUD: &suite.UserCRUDMock}
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)

	//act
	status, res := suite.UserController.PostUser(nil, req, nil)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	suite.Require().IsType(controllers.ErrorResponse{}, res)

	errRes := res.(controllers.ErrorResponse)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidBodyFields_ReturnsBadRequest() {
	var body controllers.PostUserBody

	testCase := func() {
		//arrange
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
		suite.Require().NoError(err)

		//act
		status, res := suite.UserController.PostUser(nil, req, nil)

		//assert
		suite.Equal(http.StatusBadRequest, status)
		suite.Require().IsType(controllers.ErrorResponse{}, res)

		errRes := res.(controllers.ErrorResponse)
		suite.False(errRes.Success)
		suite.Contains(errRes.Error, "username and password cannot be empty")
	}

	body = controllers.PostUserBody{
		Username: "",
		Password: "password",
	}
	suite.Run("EmptyUsername", testCase)

	body = controllers.PostUserBody{
		Username: "username",
		Password: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPostUser_WhereGetUserByUsernameReturnsError_ReturnsInternalServerError() {
	//arrange
	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}

	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, errors.New(""))

	//act
	status, res := suite.UserController.PostUser(nil, req, nil)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	suite.Require().IsType(controllers.ErrorResponse{}, res)

	errRes := res.(controllers.ErrorResponse)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, "an internal error occurred")
}

func (suite *UserControllerTestSuite) TestPostUser_WithNonUniqueUsername_ReturnsBadRequest() {
	//arrange
	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}

	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(&models.User{}, nil)

	//act
	status, res := suite.UserController.PostUser(nil, req, nil)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	suite.Require().IsType(controllers.ErrorResponse{}, res)

	errRes := res.(controllers.ErrorResponse)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, "username already exists")
}

func (suite *UserControllerTestSuite) TestPostUser_WhereCreateUserReturnsError_ReturnsInternalServerError() {
	//arrange
	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}

	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(errors.New(""))

	//act
	status, res := suite.UserController.PostUser(nil, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)

	suite.Equal(http.StatusInternalServerError, status)
	suite.Require().IsType(controllers.ErrorResponse{}, res)

	errRes := res.(controllers.ErrorResponse)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, "an internal error occurred")
}

func (suite *UserControllerTestSuite) TestPostUser_WithValidRequest_ReturnsOK() {
	//arrange
	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}

	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	status, res := suite.UserController.PostUser(nil, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.UserCRUDMock.AssertCalled(suite.T(), "CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username
	}))

	suite.Equal(http.StatusOK, status)
	suite.Require().IsType(controllers.BasicResponse{}, res)
	suite.True(res.(controllers.BasicResponse).Success)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
