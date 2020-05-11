package controllers

import (
	"log"
	"net/http"

	"github.com/blendermux/server/database"
	"github.com/blendermux/server/models"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AccountController struct {
	database.UserCRUD
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
	verr := models.ValidateUserEmail(body.Email)
	if verr.Status != models.ModelValid {
		log.Println(verr)
		sendResponse(w, errorResponse{false, "email is not valid"})
	}

	//TODO: validate password meets criteria

	//validate email is unique
	otherUser, _ := con.GetUserByEmail(body.Email)
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
	user := models.CreateNewUser(body.Email, hash)
	err = con.CreateUser(user)
	if err != nil {
		//TODO: handle error
	}

	//return success
	sendResponse(w, basicResponse{true})
}
