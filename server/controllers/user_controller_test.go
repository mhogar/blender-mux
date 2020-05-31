package controllers_test

import (
	"blendermux/server/controllers"
	"blendermux/server/controllers/mocks"
	databasemocks "blendermux/server/database/mocks"
	"blendermux/server/models"
	modelmocks "blendermux/server/models/mocks"
	"bytes"
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
	SID                           uuid.UUID
	SessionCookie                 *http.Cookie
	UserCRUDMock                  databasemocks.UserCRUD
	PasswordHasherMock            mocks.PasswordHasher
	PasswordCriteriaValidatorMock modelmocks.PasswordCriteriaValidator
	UserController                controllers.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.SID = uuid.New()
	suite.SessionCookie = &http.Cookie{
		Name:  "session",
		Value: suite.SID.String(),
	}

	suite.UserCRUDMock = databasemocks.UserCRUD{}
	suite.PasswordHasherMock = mocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = modelmocks.PasswordCriteriaValidator{}
	suite.UserController = controllers.UserController{
		UserCRUD:                  &suite.UserCRUDMock,
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
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
		req := createRequestWithJSONBody(&suite.Suite, body)

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

func (suite *UserControllerTestSuite) TestPostUser_WithErrorGettingUserByUsername_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := createRequestWithJSONBody(&suite.Suite, body)

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
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(&models.User{}, nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "username already exists")
}

func (suite *UserControllerTestSuite) TestPostUser_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPostUser_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password does not meet minimum criteria")
}

func (suite *UserControllerTestSuite) TestPostUser_WithErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := createRequestWithJSONBody(&suite.Suite, body)

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
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
	req := createRequestWithJSONBody(&suite.Suite, body)

	hash := []byte("password hash")

	suite.UserCRUDMock.On("GetUserByUsername", body.Username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.UserCRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PostUser(w, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByUsername", body.Username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.Password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.Password)
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

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorGettingUserById_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.DeleteUser(w, nil, params)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_WhereUserIsNotFound_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	id := uuid.New()
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: id.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(nil, nil)

	//act
	suite.UserController.DeleteUser(w, nil, params)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "user not found")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.DeleteUser(w, nil, params)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	user := models.CreateNewUser("username", []byte("password hash"))
	params := httprouter.Params{
		httprouter.Param{Key: "id", Value: user.ID.String()},
	}

	suite.UserCRUDMock.On("GetUserByID", mock.Anything).Return(user, nil)
	suite.UserCRUDMock.On("DeleteUser", mock.Anything).Return(nil)

	//act
	suite.UserController.DeleteUser(w, nil, params)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserByID", user.ID)
	suite.UserCRUDMock.AssertCalled(suite.T(), "DeleteUser", user)

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

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorGettingUserBySessionId_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	req, err := http.NewRequest("", "", nil)
	suite.Require().NoError(err)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(nil, errors.New(""))

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

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(nil, nil)

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

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "invalid json body")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithInvalidBodyFields_ReturnsBadRequest() {
	var body controllers.PatchUserPasswordBody

	testCase := func() {
		//arrange
		w := httptest.NewRecorder()

		req := createRequestWithJSONBody(&suite.Suite, body)
		req.AddCookie(suite.SessionCookie)

		suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)

		//act
		suite.UserController.PatchUserPassword(w, req, nil)

		//assert
		assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password and new password cannot be empty")
	}

	body = controllers.PatchUserPasswordBody{
		OldPassword: "",
		NewPassword: "new password",
	}
	suite.Run("EmptyUsername", testCase)

	body = controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "",
	}
	suite.Run("EmptyPassword", testCase)
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereOldPasswordIsInvalid_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := createRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "old password is invalid")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := createRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertErrorResponse(&suite.Suite, w.Result(), http.StatusBadRequest, "password does not meet minimum criteria")
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorHashingNewPassword_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := createRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithErrorUpdatingUser_ReturnsInternalServerError() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := createRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.UserCRUDMock.On("UpdateUser", mock.Anything).Return(errors.New(""))

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	assertInternalServerErrorResponse(&suite.Suite, w.Result())
}

func (suite *UserControllerTestSuite) TestPatchUserPassword_WithValidRequest_ReturnsOK() {
	//arrange
	w := httptest.NewRecorder()

	body := controllers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}

	req := createRequestWithJSONBody(&suite.Suite, body)
	req.AddCookie(suite.SessionCookie)

	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")
	user := models.CreateNewUser("username", oldPasswordHash)

	suite.UserCRUDMock.On("GetUserBySessionID", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(nil)
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.UserCRUDMock.On("UpdateUser", mock.Anything).Return(nil)

	//act
	suite.UserController.PatchUserPassword(w, req, nil)

	//assert
	suite.UserCRUDMock.AssertCalled(suite.T(), "GetUserBySessionID", suite.SID)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", oldPasswordHash, body.OldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", body.NewPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", body.NewPassword)
	suite.UserCRUDMock.AssertCalled(suite.T(), "UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return bytes.Equal(u.PasswordHash, newPasswordHash)
	}))

	assertSuccessResponse(&suite.Suite, w.Result())
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
