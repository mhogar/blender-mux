package controllers

import (
	"log"
	"net/http"

	"github.com/blendermux/server/dependencies"
	"github.com/blendermux/server/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AccountController struct {
	dependencies.UserCRUD
}

// PostAccount handles Post requests to "/account"
func (con *AccountController) PostAccount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		//TODO: handle error
		log.Println(err)
		sendResponse(w, errorResponse{false, "invalid json body"})
		return
	}

	//validate email
	emailValid := models.ValidateEmail(body.Email)
	if !emailValid {
		sendResponse(w, errorResponse{false, "email is not valid"})
	}

	//TODO: validate password meets criteria

	//validate email is unique
	otherUser := con.GetUserByEmail(body.Email)
	if otherUser != nil {
		sendResponse(w, errorResponse{false, "an account with that email already exists"})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		//TODO: handle error
		log.Println(err)
		sendResponse(w, errorResponse{false, "error generating password hash"})
		return
	}

	//save the user
	err = con.CreateUser(&models.User{
		uuid.New(),
		body.Email,
		hash,
	})
	if err != nil {
		//TODO: handle error
	}

	//return success
	sendResponse(w, basicResponse{true})
}
