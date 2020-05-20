package controllers

import (
	"blendermux/common"
	"log"
	"net/http"

	"github.com/google/uuid"

	"blendermux/server/database"
	"blendermux/server/models"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles requests to "/user" endpoints
type UserController struct {
	UserCRUD database.UserCRUD
}

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (con *UserController) PostUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostUserBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		sendResponse(w, http.StatusBadRequest, createErrorResponse("invalid json body"))
		return
	}

	//validate the body fields
	if body.Username == "" || body.Password == "" {
		sendResponse(w, http.StatusBadRequest, createErrorResponse("username and password cannot be empty"))
		return
	}

	//validate username is unique
	otherUser, err := con.UserCRUD.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		sendInternalErrorResponse(w)
		return
	}

	if otherUser != nil {
		sendResponse(w, http.StatusBadRequest, createErrorResponse("an user with that username already exists"))
		return
	}

	//TODO: add unit test
	//validate password meets criteria
	verr := models.ValidatePassword(body.Password)
	if verr.Status != models.ValidateErrorModelValid {
		log.Println(common.ChainError("error validating password", err))
		sendResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		sendInternalErrorResponse(w)
		return
	}

	//save the user
	user := models.CreateNewUser(body.Username, hash)
	err = con.UserCRUD.CreateUser(user)
	if err != nil {
		log.Println(common.ChainError("error saving user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}

// DeleteUser handles DELETE requests to "/user"
func (con *UserController) DeleteUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	//check for id
	idStr := params.ByName("id")
	if idStr == "" {
		sendResponse(w, http.StatusBadRequest, createErrorResponse("id must be present"))
		return
	}

	//parse the id
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(common.ChainError("error parsing user id", err))
		sendResponse(w, http.StatusBadRequest, createErrorResponse("id is in invalid format"))
		return
	}

	//delete the user
	result, err := con.UserCRUD.DeleteUser(id)
	if err != nil {
		log.Println(common.ChainError("error deleting user", err))
		sendInternalErrorResponse(w)
		return
	}

	//check if user was actually deleted
	if !result {
		sendResponse(w, http.StatusOK, createErrorResponse("could not delete user"))
		return
	}

	//return success
	sendSuccessResponse(w)
}

// PatchUserPasswordBody is the struct the body of requests to PatchUserPassword should be parsed into
type PatchUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// PatchUserPassword handles PATCH requests to "/user/password"
func (con *UserController) PatchUserPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PatchUserPasswordBody

	//get the session
	sID, err := getSessionFromRequest(req)
	if err != nil {
		log.Println(common.ChainError("error getting session id from request", err))
		sendResponse(w, http.StatusBadRequest, "session token not provided or was in invalid format")
		return
	}

	//get the user
	user, err := con.UserCRUD.GetUserBySessionId(sID)
	if err != nil {
		log.Println(common.ChainError("error getting user by session id", err))
		sendInternalErrorResponse(w)
		return
	}

	//check user was found
	if user == nil {
		sendNotAuthorizedResponse(w)
		return
	}

	//parse the body
	err = parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchUserPassword request body", err))
		sendResponse(w, http.StatusBadRequest, createErrorResponse("invalid json body"))
		return
	}

	//validate the body fields
	if body.OldPassword == "" || body.NewPassword == "" {
		sendResponse(w, http.StatusBadRequest, createErrorResponse("old password and new password cannot be empty"))
		return
	}

	//validate old password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(body.OldPassword))
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		sendResponse(w, http.StatusBadRequest, createErrorResponse("old password is invalid"))
		return
	}

	//validate new password meets critera
	verr := models.ValidatePassword(body.NewPassword)
	if verr.Status != models.ValidateErrorModelValid {
		log.Println(common.ChainError("error validating password", err))
		sendResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		sendInternalErrorResponse(w)
		return
	}

	//update the user
	user.PasswordHash = hash
	err = con.UserCRUD.UpdateUser(user)
	if err != nil {
		log.Println(common.ChainError("error updating user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}
