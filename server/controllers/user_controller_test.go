package controllers_test

import (
	"blendermux/server/controllers"
	"blendermux/server/controllers/mocks"
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
	SessionCookie      *http.Cookie
	UserCRUDMock       databasemocks.UserCRUD
	PasswordHasherMock mocks.PasswordHasher
	UserController     controllers.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.SessionCookie = &http.Cookie{
		Name:  "session",
		Value: uuid.New().String(),
	}

	suite.UserCRUDMock = databasemocks.UserCRUD{}
	suite.PasswordHasherMock = mocks.PasswordHasher{}
	suite.UserController = controllers.UserController{
		UserCRUD:       &suite.UserCRUDMock,
		PasswordHasher: &suite.PasswordHasherMock,
	}
}

func (suite *UserControllerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
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

		//assert
		assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username and password cannot be empty")
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

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username already exists")
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
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}

	hash := []byte(body.Username)

	bodyStr, err := json.Marshal(body)
	suite.Require().NoError(err)

	req, err := http.NewRequest("", "", bytes.NewReader(bodyStr))
	suite.Require().NoError(err)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.UserCRUDMock.AssertCalled(suite.T(), "CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == body.Username && bytes.Equal(u.PasswordHash, hash)
	}))

	assertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithoutIdInParams_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	//act
	suite.UserController.DeleteUser(w, nil, make(httprouter.Params, 0))

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id must be present")
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

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "id is in invalid format")
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

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
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

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", id)

	assertErrorResponse(&suite.Suite, w.Result(), http.StatusOK, "could not delete user")
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

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", id)

	assertSuccessResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithNoSessionId_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "token not provided")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithInvalidSessionId_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)

	suite.SessionCookie.Value = "invalid session id"
	req.AddCookie(suite.SessionCookie)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "invalid format")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereGetUserBySessionReturnsError_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionId", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereNoUserIsFound_ReturnsUnauthorized() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionId", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusUnauthorized, "no user")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", strings.NewReader("0"))
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionId", mock.Anything).Return(&models.User{}, nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
