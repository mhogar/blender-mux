package controllers

import (
	"log"
	"net/http"

	"github.com/blendermux/server/dependencies"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type SessionController struct {
	dependencies.UserCRUD
}

//PostLogin handles Post requests to "/login"
func (con SessionController) PostLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	//get the user
	user := con.GetUserByEmail(body.Email)
	if user == nil {
		sendResponse(w, errorResponse{false, "Invalid email and/or password"})
		return
	}

	//validate the password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(body.Password))
	if err != nil {
		log.Println(err)
		sendResponse(w, errorResponse{false, "Invalid email and/or password"})
		return
	}

	//return success
	//TODO: gernerate and return session token
	sendResponse(w, basicResponse{true})
}
