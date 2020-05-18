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

// UserController handles request to /user endpoints
type UserController struct {
	database.UserCRUD
}

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (con *UserController) PostUser(_ http.ResponseWriter, req *http.Request, _ httprouter.Params) (int, interface{}) {
	var body PostUserBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		return http.StatusBadRequest, createErrorResponse("invalid json body")
	}

	//validate the body fields
	if body.Username == "" || body.Password == "" {
		return http.StatusBadRequest, createErrorResponse("username and password cannot be empty")
	}

	//validate username is unique
	otherUser, err := con.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return createInternalErrorResponse()
	}

	if otherUser != nil {
		return http.StatusBadRequest, createErrorResponse("an user with that username already exists")
	}

	//TODO: validate password meets criteria

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return createInternalErrorResponse()
	}

	//save the user
	user := models.CreateNewUser(body.Username, hash)
	err = con.CreateUser(user)
	if err != nil {
		log.Println(common.ChainError("error saving user", err))
		return createInternalErrorResponse()
	}

	//return success
	return createSuccessResponse()
}
