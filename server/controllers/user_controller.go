package controllers

import (
	"blendermux/common"
	"log"
	"net/http"

	"blendermux/server/database"
	"blendermux/server/models"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles requests to "/user" endpoints
type UserController struct {
	database.UserCRUD
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
	otherUser, err := con.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		sendInternalErrorResponse(w)
		return
	}

	if otherUser != nil {
		sendResponse(w, http.StatusBadRequest, createErrorResponse("an user with that username already exists"))
		return
	}

	//TODO: validate password meets criteria

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		sendInternalErrorResponse(w)
		return
	}

	//save the user
	user := models.CreateNewUser(body.Username, hash)
	err = con.CreateUser(user)
	if err != nil {
		log.Println(common.ChainError("error saving user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}

// DeleteUser handles DELETE requests to "/user"
func (con *UserController) DeleteUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//return success
	sendSuccessResponse(w)
}
