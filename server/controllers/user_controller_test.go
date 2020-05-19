package controllers_test

import (
	"blendermux/common"
	"blendermux/server/controllers"
	databasemocks "blendermux/server/database/mocks"
	"blendermux/server/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

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
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)

	//act
	suite.UserController.PostUser(w, req, nil)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidBodyFields_ReturnsBadRequest() {
	var body controllers.PostUserBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()

		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
		suite.Require().NoError(err)

		//act
		suite.UserController.PostUser(w, req, nil)

		var res controllers.ErrorResponse
		status := common.ParseResponse(&suite.Suite, w.Result(), &res)

		//assert
		suite.Equal(http.StatusBadRequest, status)
		suite.False(res.Success)
		suite.Contains(res.Error, "username and password cannot be empty")
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
	w := httptest.NewRecorder()

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
	suite.UserController.PostUser(w, req, nil)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)

	suite.Equal(http.StatusInternalServerError, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "an internal error occurred")
}

func (suite *UserControllerTestSuite) TestPostUser_WithNonUniqueUsername_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

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
	suite.UserController.PostUser(w, req, nil)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)

	suite.Equal(http.StatusBadRequest, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "username already exists")
}

func (suite *UserControllerTestSuite) TestPostUser_WhereCreateUserReturnsError_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

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
	suite.UserController.PostUser(w, req, nil)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.UserCRUDMock.AssertCalled(suite.T(), "CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username
	}))

	suite.Equal(http.StatusInternalServerError, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "an internal error occurred")
}

func (suite *UserControllerTestSuite) TestPostUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

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
	suite.UserController.PostUser(w, req, nil)

	var res controllers.BasicResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.UserCRUDMock.AssertCalled(suite.T(), "CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username
	}))

	suite.Equal(http.StatusOK, status)
	suite.True(res.Success)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithoutIdInParams_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	//act
	suite.UserController.DeleteUser(w, nil, make(httprouter.Params, 0))

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "id must be present")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithIdInInvalidFormat_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	id := 0
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: string(id)},
	}

	//act
	suite.UserController.DeleteUser(w, nil, params)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "id is in invalid format")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WhereDeleteUserReturnsError_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(false, errors.New(""))

	//act
	suite.UserController.DeleteUser(w, nil, params)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", id)

	suite.Equal(http.StatusInternalServerError, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "an internal error occurred")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WhereDeleteUserReturnsFalseResult_ReturnsOKWithError() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(false, nil)

	//act
	suite.UserController.DeleteUser(w, nil, params)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", id)

	suite.Equal(http.StatusOK, status)
	suite.False(res.Success)
	suite.Contains(res.Error, "could not delete user")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(true, nil)

	//act
	suite.UserController.DeleteUser(w, nil, params)

	var res controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, w.Result(), &res)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", id)

	suite.Equal(http.StatusOK, status)
	suite.True(res.Success)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
